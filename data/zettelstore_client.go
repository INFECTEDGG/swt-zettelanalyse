package data

import (
	"context"
	"fmt"
	"net/url"
	"strings"
	"time"

	"Zettelanalyse/model"

	"t73f.de/r/zsc/client"
)

// BaseURL ist die Adresse des lokalen Zettelstores.
const BaseURL = "http://localhost:23123"

// ZettelstoreClient verwaltet den offiziellen ZSC-Client.
type ZettelstoreClient struct {
	client *client.Client
}

// NewZettelstoreClient erstellt einen neuen ZSC-Client.
func NewZettelstoreClient() *ZettelstoreClient {

	// BaseURL in ein URL-Objekt umwandeln.
	adresse, err := url.Parse(BaseURL)

	if err != nil {
		panic(fmt.Errorf("Ungültige BaseURL: %w", err))
	}

	// Offiziellen Zettelstore-Client erstellen.
	zscClient := client.NewClient(adresse)

	return &ZettelstoreClient{
		client: zscClient,
	}
}

// GetZettel ruft alle Zettel über die offizielle API ab.
func (c *ZettelstoreClient) GetZettel() ([]model.Zettel, error) {

	// Timeout für die Anfrage festlegen.
	ctx, cancel := context.WithTimeout(
		context.Background(),
		10*time.Second,
	)
	defer cancel()

	// Alle Zettel über ZSC abfragen.
	ergebnis, err := c.client.QueryZettel(ctx, "")

	if err != nil {
		return nil, fmt.Errorf(
			"Zettel konnten nicht geladen werden: %w",
			err,
		)
	}

	var zettel []model.Zettel

	// Ergebnis in unser eigenes Datenmodell umwandeln.
	for _, eintrag := range ergebnis {

		zeile := strings.TrimSpace(string(eintrag))

		if zeile == "" {
			continue
		}

		// Format: ID + Leerzeichen + Titel
		parts := strings.SplitN(zeile, " ", 2)

		id := parts[0]
		titel := id

		if len(parts) == 2 {
			titel = strings.TrimSpace(parts[1])
		}

		zettel = append(zettel, model.Zettel{
			ID:    id,
			Title: titel,
		})
	}

	return zettel, nil
}
