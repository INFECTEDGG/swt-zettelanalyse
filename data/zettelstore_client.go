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

// ZettelstoreClient provides read-only access to the Zettelstore.
type ZettelstoreClient struct {
	client *client.Client
}

// NewZettelstoreClient creates a client for the given server address.
func NewZettelstoreClient(address string) *ZettelstoreClient {
	address = strings.TrimRight(strings.TrimSpace(address), "/")

	serverURL, err := url.Parse(address)
	if err != nil {
		panic(fmt.Errorf("invalid Zettelstore address: %w", err))
	}

	return &ZettelstoreClient{
		client: client.NewClient(serverURL),
	}
}

// GetZettel retrieves all zettel IDs, titles, and tags.
func (c *ZettelstoreClient) GetZettel() ([]model.Zettel, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, _, results, err := c.client.QueryZettelData(ctx, "")
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve zettel list: %w", err)
	}

	zettels := make([]model.Zettel, 0, len(results))
	for _, entry := range results {
		id := entry.ID.String()
		title := strings.TrimSpace(entry.Meta["title"])
		if title == "" {
			title = id
		}

		zettels = append(zettels, model.Zettel{
			ID:       id,
			Title:    title,
			Tags:     strings.Fields(entry.Meta["tags"]),
			Created:  entry.Meta["created"],
			Modified: entry.Meta["modified"],
		})
	}

	return zettels, nil
}
