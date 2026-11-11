package data

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// Test mit nachgebautem Zettelstore, braucht keinen laufenden Zettelstore (NFA-06).
func TestGetZettel(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("20261111100000 Erster Zettel\n20261111110000 Zweiter Zettel\n20261111120000\n"))
	}))
	defer server.Close()

	zettel, err := NewZettelstoreClient(server.URL).GetZettel()
	if err != nil {
		t.Fatal(err)
	}
	if len(zettel) != 3 {
		t.Fatalf("3 Zettel erwartet, %d erhalten", len(zettel))
	}
	if zettel[1].Title != "Zweiter Zettel" {
		t.Errorf("falscher Titel: %q", zettel[1].Title)
	}
	if zettel[2].Title != "20261111120000" {
		t.Errorf("ohne Titel sollte die ID stehen, war: %q", zettel[2].Title)
	}
}

// Fehler, wenn kein Zettelstore erreichbar ist.
func TestGetZettelWithoutZettelstore(t *testing.T) {
	if _, err := NewZettelstoreClient("http://127.0.0.1:1").GetZettel(); err == nil {
		t.Fatal("Fehler erwartet")
	}
}

// Test mit dem echten Zettelstore. Wird übersprungen, wenn er nicht läuft.
func TestRealZettelstore(t *testing.T) {
	zettel, err := NewZettelstoreClient("http://127.0.0.1:23123").GetZettel()
	if err != nil {
		t.Skip("Zettelstore läuft nicht:", err)
	}
	t.Logf("%d Zettel geladen:", len(zettel))
	for i, z := range zettel {
		t.Logf("%3d. %s  %s", i+1, z.ID, z.Title)
	}
}
