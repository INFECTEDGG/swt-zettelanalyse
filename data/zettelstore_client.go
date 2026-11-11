// Package data ist das einzige Paket, das mit dem Zettelstore spricht (ADR 01).
// Es importiert nur model und greift nur lesend zu (NFA-01).
package data

import (
	"bufio"
	"fmt"
	"net/http"
	"strings"
	"time"

	"Zettelanalyse/model"
)

// ZettelstoreClient fragt einen Zettelstore über dessen HTTP-API ab.
// Er erfüllt das Interface ZettelQuelle aus dem Paket application (ADR 01).
type ZettelstoreClient struct {
	baseURL string
	client  *http.Client
}

// NewZettelstoreClient erstellt einen Client für die angegebene Adresse (Issue #1),
// z. B. "http://127.0.0.1:23123". Der Zettelstore kann auch auf einem anderen Rechner laufen.
func NewZettelstoreClient(address string) *ZettelstoreClient {
	return &ZettelstoreClient{
		baseURL: strings.TrimRight(strings.TrimSpace(address), "/"),
		client:  &http.Client{Timeout: 10 * time.Second},
	}
}

// GetZettel ruft alle Zettel über GET /z ab.
// Der Zettelstore antwortet mit einer Textliste: je Zeile "<ID> <Titel>".
func (c *ZettelstoreClient) GetZettel() ([]model.Zettel, error) {
	resp, err := c.client.Get(c.baseURL + "/z")
	if err != nil {
		return nil, fmt.Errorf("Zettelstore unter %s nicht erreichbar: %w", c.baseURL, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API-Fehler: %s", resp.Status)
	}

	var zettel []model.Zettel
	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		parts := strings.SplitN(line, " ", 2)
		id, title := parts[0], parts[0]
		if len(parts) > 1 {
			title = strings.TrimSpace(parts[1])
		}
		zettel = append(zettel, model.Zettel{ID: id, Title: title})
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("Lesefehler: %w", err)
	}
	return zettel, nil
}
