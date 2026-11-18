package data

import (
	"errors"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestGetZettel verifies that IDs, titles, tags, and timestamps are
// retrieved from a Zettelstore metadata response (NFA-06).
func TestGetZettel(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet || r.URL.Path != "/z" || r.URL.Query().Get("enc") != "data" {
			http.Error(w, "unexpected request", http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		_, _ = w.Write([]byte(`(meta-list
 (query "")
 (human "")
 (zettel "20261111100000"
  (meta (tags "#go #uni") (title "First Zettel")
        (created "20261111100000") (modified "20261111103000"))
  (rights read))
 (zettel "20261111110000" (meta (tags "#go") (title "Second Zettel")) (rights read))
 (zettel "20261111120000" (meta) (rights read)))`))
	}))
	defer server.Close()

	zettels, err := NewZettelstoreClient(server.URL).GetZettel()
	if err != nil {
		t.Fatalf("failed to retrieve zettels: %v", err)
	}
	if len(zettels) != 3 {
		t.Fatalf("expected 3 zettels, got %d", len(zettels))
	}

	first := zettels[0]
	if first.ID != "20261111100000" {
		t.Errorf("unexpected first zettel ID: %q", first.ID)
	}
	if first.Title != "First Zettel" {
		t.Errorf("unexpected first zettel title: %q", first.Title)
	}
	if len(first.Tags) != 2 || first.Tags[0] != "#go" || first.Tags[1] != "#uni" {
		t.Errorf("unexpected first zettel tags: %q", first.Tags)
	}
	if first.Created != "20261111100000" {
		t.Errorf("Created = %q, want %q", first.Created, "20261111100000")
	}
	if first.Modified != "20261111103000" {
		t.Errorf("Modified = %q, want %q", first.Modified, "20261111103000")
	}

	if zettels[1].Title != "Second Zettel" {
		t.Errorf("unexpected second zettel title: %q", zettels[1].Title)
	}

	third := zettels[2]
	if third.Title != "20261111120000" {
		t.Errorf("expected ID as fallback title, got %q", third.Title)
	}
	if third.Created != "" || third.Modified != "" {
		t.Errorf("expected empty timestamps, got Created = %q and Modified = %q",
			third.Created, third.Modified)
	}
}

// TestGetZettelWithoutZettelstore verifies that dial errors are returned.
func TestGetZettelWithoutZettelstore(t *testing.T) {
	server := httptest.NewServer(http.NotFoundHandler())
	address := server.URL
	server.Close()

	_, err := NewZettelstoreClient(address).GetZettel()

	var opErr *net.OpError
	if !errors.As(err, &opErr) || opErr.Op != "dial" {
		t.Fatalf("expected a dial error, got %v", err)
	}
}

// TestRealZettelstore logs metadata from a running local Zettelstore.
func TestRealZettelstore(t *testing.T) {
	zettels, err := NewZettelstoreClient("http://127.0.0.1:23123").GetZettel()
	if err != nil {
		var opErr *net.OpError
		if errors.As(err, &opErr) && opErr.Op == "dial" {
			t.Skipf("could not connect to the local Zettelstore: %v", err)
		}
		t.Fatalf("could not retrieve zettels from the local Zettelstore: %v", err)
	}

	t.Logf("Retrieved %d zettels", len(zettels))
	for index, zettel := range zettels {
		t.Logf("%3d. %s %s | Tags: %q | Created: %q | Modified: %q",
			index+1, zettel.ID, zettel.Title, zettel.Tags, zettel.Created, zettel.Modified)
	}
}
