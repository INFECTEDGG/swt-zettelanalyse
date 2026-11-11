package application

import (
	"fmt"
	"sync"
	"time"

	"github.com/INFECTEDGG/swt-zettelanalyse/model"
)

// App koordiniert die Anwendungsfaelle und haelt den zuletzt erfolgreich
// geladenen, bereits um Systemzettel bereinigten Datenbestand.
type App struct {
	quelle ZettelQuelle

	ladenMu sync.Mutex
	mu      sync.RWMutex
	zettel  []model.Zettel
}

func NewApp(quelle ZettelQuelle) *App {
	return &App{quelle: quelle}
}

// Aktualisieren konfiguriert bei Bedarf die Adresse der Quelle und liest den
// Datenbestand neu ein. Bei einem Fehler bleibt der letzte erfolgreiche Stand
// erhalten.
func (a *App) Aktualisieren(adresse string) (int, error) {
	a.ladenMu.Lock()
	defer a.ladenMu.Unlock()

	if a.quelle == nil {
		return 0, fmt.Errorf("keine Zettelquelle konfiguriert")
	}
	adressierbar, istAdressierbar := a.quelle.(AdressierbareZettelQuelle)
	if istAdressierbar {
		if err := adressierbar.SetAdresse(adresse); err != nil {
			return 0, err
		}
	} else if adresse != "" {
		return 0, fmt.Errorf("die konfigurierte Zettelquelle unterstuetzt keine Adressaenderung")
	}

	zettel, err := LadeZettel(a.quelle)
	if err != nil {
		return 0, err
	}

	a.mu.Lock()
	a.zettel = zettel
	a.mu.Unlock()
	return len(zettel), nil
}

func (a *App) zettelKopie() []model.Zettel {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return kopiereZettel(a.zettel)
}

func (a *App) Monatsstatistik(jahr int) [12]int {
	return ZettelProMonat(a.zettelKopie(), jahr)
}

func (a *App) AktivsterMonat(jahr int) (model.JahrMonat, bool) {
	return AktivsterMonat(a.zettelKopie(), jahr)
}

func (a *App) ZettelDesMonats(jahr int, monat time.Month) []model.Zettel {
	return FilterNachMonat(a.zettelKopie(), jahr, monat)
}

func (a *App) HaeufigsteTags() []model.TagAnzahl {
	return HaeufigsteTags(a.zettelKopie())
}

// Jahre liefert alle Kalenderjahre mit auswertbaren Zeitstempeln, neuestes
// Jahr zuerst.
func (a *App) Jahre() []int {
	return VerfuegbareJahre(a.zettelKopie())
}

func kopiereZettel(eingabe []model.Zettel) []model.Zettel {
	ausgabe := make([]model.Zettel, len(eingabe))
	copy(ausgabe, eingabe)
	for i := range ausgabe {
		ausgabe[i].Tags = append([]string(nil), eingabe[i].Tags...)
	}
	return ausgabe
}
