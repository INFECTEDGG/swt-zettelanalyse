package data

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"
	"unicode"
	"unicode/utf8"

	"github.com/INFECTEDGG/swt-zettelanalyse/model"
)

const (
	Standardadresse        = "http://127.0.0.1:23123"
	standardHTTPTimeout    = 10 * time.Second
	maximaleAntwortgroesse = 32 << 20
)

// ZettelstoreClient liest die Metadatenliste des Zettelstores. Er stellt
// absichtlich nur GET-basierte Operationen bereit.
type ZettelstoreClient struct {
	mu        sync.RWMutex
	basisURL  *url.URL
	konfigErr error
	http      *http.Client
}

func NewZettelstoreClient(adresse string) *ZettelstoreClient {
	client := &ZettelstoreClient{
		http: &http.Client{
			Timeout: standardHTTPTimeout,
			CheckRedirect: func(_ *http.Request, _ []*http.Request) error {
				return http.ErrUseLastResponse
			},
		},
	}
	_ = client.SetAdresse(adresse)
	return client
}

// SetAdresse validiert die Basisadresse, bevor sie fuer Anfragen verwendet
// wird. Zugangsdaten, Query-Parameter und Fragmente sind nicht erlaubt.
func (c *ZettelstoreClient) SetAdresse(adresse string) error {
	parsed, err := validiereAdresse(adresse)
	c.mu.Lock()
	defer c.mu.Unlock()
	if err != nil {
		c.konfigErr = err
		return err
	}
	c.basisURL = parsed
	c.konfigErr = nil
	return nil
}

func validiereAdresse(adresse string) (*url.URL, error) {
	adresse = strings.TrimSpace(adresse)
	parsed, err := url.Parse(adresse)
	if err != nil {
		return nil, fmt.Errorf("ungueltige Zettelstore-Adresse: %w", err)
	}
	if parsed.Scheme != "http" && parsed.Scheme != "https" {
		return nil, fmt.Errorf("ungueltige Zettelstore-Adresse: nur http und https werden unterstuetzt")
	}
	if parsed.Host == "" || parsed.Hostname() == "" {
		return nil, fmt.Errorf("ungueltige Zettelstore-Adresse: Hostname fehlt")
	}
	if parsed.User != nil {
		return nil, fmt.Errorf("ungueltige Zettelstore-Adresse: eingebettete Zugangsdaten sind nicht erlaubt")
	}
	if parsed.RawQuery != "" || parsed.Fragment != "" {
		return nil, fmt.Errorf("ungueltige Zettelstore-Adresse: Query und Fragment sind nicht erlaubt")
	}
	parsed.Path = strings.TrimRight(parsed.Path, "/")
	parsed.RawPath = ""
	return parsed, nil
}

// GetZettel ruft GET /z?enc=data auf. Das Format ist die von Zettelstore
// dokumentierte symbolische Metadatenliste.
func (c *ZettelstoreClient) GetZettel() ([]model.Zettel, error) {
	c.mu.RLock()
	konfigErr := c.konfigErr
	var basis url.URL
	if c.basisURL != nil {
		basis = *c.basisURL
	}
	httpClient := c.http
	c.mu.RUnlock()

	if konfigErr != nil {
		return nil, konfigErr
	}
	if basis.Host == "" {
		return nil, fmt.Errorf("keine Zettelstore-Adresse konfiguriert")
	}

	basis.Path = strings.TrimRight(basis.Path, "/") + "/z"
	basis.RawPath = ""
	query := basis.Query()
	query.Set("enc", "data")
	basis.RawQuery = query.Encode()

	request, err := http.NewRequest(http.MethodGet, basis.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("Zettelstore-Anfrage konnte nicht erstellt werden: %w", err)
	}
	request.Header.Set("Accept", "text/plain")
	request.Header.Set("User-Agent", "swt-zettelanalyse/1")

	response, err := httpClient.Do(request)
	if err != nil {
		return nil, fmt.Errorf("Zettelstore ist nicht erreichbar: %w", err)
	}
	defer response.Body.Close()

	body, err := io.ReadAll(io.LimitReader(response.Body, maximaleAntwortgroesse+1))
	if err != nil {
		return nil, fmt.Errorf("Zettelstore-Antwort konnte nicht gelesen werden: %w", err)
	}
	if len(body) > maximaleAntwortgroesse {
		return nil, fmt.Errorf("Zettelstore-Antwort ist groesser als %d MiB", maximaleAntwortgroesse>>20)
	}
	if response.StatusCode != http.StatusOK {
		detail := strings.TrimSpace(string(body))
		if len(detail) > 300 {
			detail = detail[:300] + "..."
		}
		if detail == "" {
			return nil, fmt.Errorf("Zettelstore meldet HTTP-Fehler %s", response.Status)
		}
		return nil, fmt.Errorf("Zettelstore meldet HTTP-Fehler %s: %s", response.Status, detail)
	}

	zettel, err := parseZettelListe(body)
	if err != nil {
		return nil, fmt.Errorf("ungueltige Zettelstore-Antwort: %w", err)
	}
	return zettel, nil
}

type sxArt uint8

const (
	sxSymbol sxArt = iota
	sxString
	sxListe
)

type sxWert struct {
	art   sxArt
	text  string
	liste []sxWert
}

type sxParser struct {
	daten []byte
	pos   int
}

func parseSExpression(daten []byte) (sxWert, error) {
	if !utf8.Valid(daten) {
		return sxWert{}, fmt.Errorf("Antwort ist kein gueltiges UTF-8")
	}
	p := sxParser{daten: daten}
	p.ueberspringeLeerraum()
	if p.pos == len(p.daten) {
		return sxWert{}, fmt.Errorf("leere Antwort")
	}
	wert, err := p.parseWert(0)
	if err != nil {
		return sxWert{}, err
	}
	p.ueberspringeLeerraum()
	if p.pos != len(p.daten) {
		return sxWert{}, fmt.Errorf("unerwartete Daten ab Position %d", p.pos)
	}
	return wert, nil
}

func (p *sxParser) parseWert(tiefe int) (sxWert, error) {
	if tiefe > 64 {
		return sxWert{}, fmt.Errorf("symbolischer Ausdruck ist zu tief verschachtelt")
	}
	p.ueberspringeLeerraum()
	if p.pos >= len(p.daten) {
		return sxWert{}, fmt.Errorf("unerwartetes Ende der Antwort")
	}
	switch p.daten[p.pos] {
	case '(':
		p.pos++
		liste := make([]sxWert, 0)
		for {
			p.ueberspringeLeerraum()
			if p.pos >= len(p.daten) {
				return sxWert{}, fmt.Errorf("nicht geschlossene Liste")
			}
			if p.daten[p.pos] == ')' {
				p.pos++
				return sxWert{art: sxListe, liste: liste}, nil
			}
			wert, err := p.parseWert(tiefe + 1)
			if err != nil {
				return sxWert{}, err
			}
			liste = append(liste, wert)
		}
	case ')':
		return sxWert{}, fmt.Errorf("unerwartete schliessende Klammer an Position %d", p.pos)
	case '"':
		text, err := p.parseString()
		if err != nil {
			return sxWert{}, err
		}
		return sxWert{art: sxString, text: text}, nil
	default:
		start := p.pos
		for p.pos < len(p.daten) {
			r, groesse := utf8.DecodeRune(p.daten[p.pos:])
			if unicode.IsSpace(r) || r == '(' || r == ')' {
				break
			}
			p.pos += groesse
		}
		if start == p.pos {
			return sxWert{}, fmt.Errorf("ungueltiges Symbol an Position %d", p.pos)
		}
		return sxWert{art: sxSymbol, text: string(p.daten[start:p.pos])}, nil
	}
}

func (p *sxParser) parseString() (string, error) {
	p.pos++
	var ergebnis strings.Builder
	for p.pos < len(p.daten) {
		r, groesse := utf8.DecodeRune(p.daten[p.pos:])
		p.pos += groesse
		switch r {
		case '"':
			return ergebnis.String(), nil
		case '\\':
			if p.pos >= len(p.daten) {
				return "", fmt.Errorf("unvollstaendige Escape-Sequenz")
			}
			escape := p.daten[p.pos]
			p.pos++
			switch escape {
			case '"', '\\':
				ergebnis.WriteByte(escape)
			case 't':
				ergebnis.WriteByte('\t')
			case 'n':
				ergebnis.WriteByte('\n')
			case 'x':
				wert, err := p.parseHexRune(2)
				if err != nil {
					return "", err
				}
				ergebnis.WriteRune(wert)
			case 'u':
				wert, err := p.parseHexRune(4)
				if err != nil {
					return "", err
				}
				ergebnis.WriteRune(wert)
			case 'U':
				wert, err := p.parseHexRune(6)
				if err != nil {
					return "", err
				}
				ergebnis.WriteRune(wert)
			default:
				return "", fmt.Errorf("unbekannte Escape-Sequenz \\%c", escape)
			}
		default:
			ergebnis.WriteRune(r)
		}
	}
	return "", fmt.Errorf("nicht geschlossene Zeichenkette")
}

func (p *sxParser) parseHexRune(anzahl int) (rune, error) {
	if p.pos+anzahl > len(p.daten) {
		return 0, fmt.Errorf("unvollstaendige hexadezimale Escape-Sequenz")
	}
	text := string(p.daten[p.pos : p.pos+anzahl])
	p.pos += anzahl
	wert, err := strconv.ParseUint(text, 16, 32)
	if err != nil {
		return 0, fmt.Errorf("ungueltige hexadezimale Escape-Sequenz %q", text)
	}
	r := rune(wert)
	if !utf8.ValidRune(r) {
		return 0, fmt.Errorf("ungueltiger Unicode-Codepoint U+%X", wert)
	}
	return r, nil
}

func (p *sxParser) ueberspringeLeerraum() {
	for p.pos < len(p.daten) {
		r, groesse := utf8.DecodeRune(p.daten[p.pos:])
		if !unicode.IsSpace(r) {
			return
		}
		p.pos += groesse
	}
}

func parseZettelListe(daten []byte) ([]model.Zettel, error) {
	wurzel, err := parseSExpression(daten)
	if err != nil {
		return nil, err
	}
	if !istListeMitKopf(wurzel, "meta-list") {
		return nil, fmt.Errorf("erwartet wurde eine meta-list")
	}

	ergebnis := make([]model.Zettel, 0)
	for _, eintrag := range wurzel.liste[1:] {
		if !istListeMitKopf(eintrag, "zettel") {
			continue
		}
		zettel, err := parseZettel(eintrag)
		if err != nil {
			return nil, err
		}
		ergebnis = append(ergebnis, zettel)
	}

	return ergebnis, nil
}

func parseZettel(wert sxWert) (model.Zettel, error) {
	if len(wert.liste) < 3 || wert.liste[1].art != sxString {
		return model.Zettel{}, fmt.Errorf("Zettel-Eintrag ohne gueltige ID")
	}
	id := wert.liste[1].text
	if !istZettelID(id) {
		return model.Zettel{}, fmt.Errorf("ungueltige Zettel-ID %q", id)
	}

	var meta *sxWert
	for i := 2; i < len(wert.liste); i++ {
		if istListeMitKopf(wert.liste[i], "meta") {
			meta = &wert.liste[i]
			break
		}
	}
	if meta == nil {
		return model.Zettel{}, fmt.Errorf("Zettel %s enthaelt keine Metadatenliste", id)
	}

	metadaten := make(map[string]string, len(meta.liste)-1)
	for _, paar := range meta.liste[1:] {
		if paar.art != sxListe || len(paar.liste) != 2 || paar.liste[0].art != sxSymbol || paar.liste[1].art != sxString {
			return model.Zettel{}, fmt.Errorf("Zettel %s enthaelt ungueltige Metadaten", id)
		}
		metadaten[paar.liste[0].text] = paar.liste[1].text
	}

	erstellt, err := parseZeitstempel(metadaten["created"])
	if err != nil {
		return model.Zettel{}, fmt.Errorf("Zettel %s hat ungueltiges created-Metadatum: %w", id, err)
	}
	if erstellt.IsZero() {
		// Bei fehlendem computed-Metadatum ist eine zeitbasierte Zettel-ID ein
		// kontrollierter Fallback. Nicht zeitbasierte IDs bleiben ohne Datum.
		erstellt, _ = parseZeitstempel(id)
	}
	return model.Zettel{
		ID:       id,
		Titel:    metadaten["title"],
		Tags:     strings.Fields(metadaten["tags"]),
		Erstellt: erstellt,
		Rolle:    metadaten["role"],
	}, nil
}

func istListeMitKopf(wert sxWert, kopf string) bool {
	return wert.art == sxListe && len(wert.liste) > 0 && wert.liste[0].art == sxSymbol && wert.liste[0].text == kopf
}

func istZettelID(id string) bool {
	if len(id) != 14 {
		return false
	}
	for _, zeichen := range id {
		if zeichen < '0' || zeichen > '9' {
			return false
		}
	}
	return true
}

func parseZeitstempel(wert string) (time.Time, error) {
	if wert == "" {
		return time.Time{}, nil
	}
	zeitpunkt, err := time.Parse("20060102150405", wert)
	if err != nil {
		return time.Time{}, errors.New("erwartet wird JJJJMMTThhmmss")
	}
	return zeitpunkt, nil
}
