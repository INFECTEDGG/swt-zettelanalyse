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

	// Initialize the official ZSC client.
	zscClient := client.NewClient(serverURL)

	return &ZettelstoreClient{
		client: zscClient,
	}
}

// GetZettel retrieves all zettel IDs and titles.
func (c *ZettelstoreClient) GetZettel() ([]model.Zettel, error) {
	ctx, cancel := context.WithTimeout(
		context.Background(),
		10*time.Second,
	)
	defer cancel()

	// Retrieve zettel using the official ZSC client.
	results, err := c.client.QueryZettel(ctx, "")
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve zettel: %w", err)
	}

	var zettels []model.Zettel

	// Convert the API response into application models.
	for _, entry := range results {
		line := strings.TrimSpace(string(entry))

		if line == "" {
			continue
		}

		parts := strings.SplitN(line, " ", 2)

		id := parts[0]
		title := id

		if len(parts) == 2 {
			title = strings.TrimSpace(parts[1])
		}

		zettels = append(zettels, model.Zettel{
			ID:    id,
			Title: title,
		})
	}

	return zettels, nil
}
