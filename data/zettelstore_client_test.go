package data

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// TestGetZettel verifies that IDs, titles, and tags are retrieved from a
// Zettelstore metadata response (NFA-06).
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
  (zettel "20261111100000" (meta (tags "#go #uni") (title "First Zettel")) (rights read))
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
	if zettels[0].ID != "20261111100000" {
		t.Errorf("unexpected first zettel ID: %q", zettels[0].ID)
	}
	if zettels[1].Title != "Second Zettel" {
		t.Errorf("unexpected second zettel title: %q", zettels[1].Title)
	}
	if zettels[2].Title != "20261111120000" {
		t.Errorf("expected ID as fallback title, got %q", zettels[2].Title)
	}
	if len(zettels[0].Tags) != 2 || zettels[0].Tags[0] != "#go" || zettels[0].Tags[1] != "#uni" {
		t.Errorf("unexpected tags: %#v", zettels[0].Tags)
	}
}

// TestGetZettelWithoutZettelstore verifies that connection errors are returned.
func TestGetZettelWithoutZettelstore(t *testing.T) {
	_, err := NewZettelstoreClient("http://127.0.0.1:1").GetZettel()
	if err == nil {
		t.Fatal("expected a connection error")
	}
}

// TestRealZettelstore retrieves zettels from a running Zettelstore instance
// and skips the test if unavailable.
func TestRealZettelstore(t *testing.T) {
	zettels, err := NewZettelstoreClient("http://127.0.0.1:23123").GetZettel()
	if err != nil {
		t.Skipf("Zettelstore is unavailable: %v", err)
	}

	t.Logf("Retrieved %d zettels", len(zettels))
	for index, zettel := range zettels {
		t.Logf("%3d. %s %s", index+1, zettel.ID, strings.TrimSpace(zettel.Title))
	}
}
