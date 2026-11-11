package data

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

const gueltigeAntwort = `(meta-list
  (query "")
  (human "")
  (zettel "20260115083000"
    (meta
      (created "20260115083000")
      (modified "20260116120000")
      (role "zettel")
      (tags "#go #uni")
      (title "Go \x26 Tk \"Desktop\""))
    (rights read))
  (zettel "20260201090000"
    (meta
      (role "zettel")
      (title "Ohne computed created \U01F642"))
    (rights read)))`

func TestZettelstoreClientErfolgreicheAPIAnfrage(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet || r.URL.Path != "/z" || r.URL.Query().Get("enc") != "data" {
			http.Error(w, fmt.Sprintf("unerwartete Anfrage: %s %s", r.Method, r.URL.String()), http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		_, _ = w.Write([]byte(gueltigeAntwort))
	}))
	defer server.Close()

	client := NewZettelstoreClient(server.URL)
	zettel, err := client.GetZettel()
	if err != nil {
		t.Fatalf("GetZettel() Fehler: %v", err)
	}
	if len(zettel) != 2 {
		t.Fatalf("GetZettel() lieferte %d Zettel", len(zettel))
	}
	if zettel[0].Titel != `Go & Tk "Desktop"` || len(zettel[0].Tags) != 2 {
		t.Fatalf("Metadaten falsch dekodiert: %#v", zettel[0])
	}
	if zettel[0].Erstellt != time.Date(2026, time.January, 15, 8, 30, 0, 0, time.UTC) {
		t.Fatalf("created falsch dekodiert: %v", zettel[0].Erstellt)
	}
	if zettel[1].Erstellt.IsZero() || !strings.Contains(zettel[1].Titel, "🙂") {
		t.Fatalf("Fallback oder Unicode-Escape falsch: %#v", zettel[1])
	}
}

func TestZettelstoreClientUnterstuetztURLPrefix(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/mein-prefix/z" {
			http.Error(w, "falscher Pfad", http.StatusNotFound)
			return
		}
		_, _ = w.Write([]byte(`(meta-list (query "") (human ""))`))
	}))
	defer server.Close()

	client := NewZettelstoreClient(server.URL + "/mein-prefix/")
	if _, err := client.GetZettel(); err != nil {
		t.Fatalf("GetZettel() mit URL-Prefix: %v", err)
	}
}

func TestZettelstoreClientHTTPFehler(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		http.Error(w, "Dienst nicht bereit", http.StatusServiceUnavailable)
	}))
	defer server.Close()

	_, err := NewZettelstoreClient(server.URL).GetZettel()
	if err == nil || !strings.Contains(err.Error(), "503") {
		t.Fatalf("erwarteter HTTP-Fehler, erhalten: %v", err)
	}
}

func TestZettelstoreClientUngueltigeAntwort(t *testing.T) {
	tests := []string{
		`kein-ausdruck`,
		`(meta-list (zettel "20260101120000" (meta (created "kaputt"))))`,
		`(meta-list (zettel "keine-id" (meta)))`,
		`(meta-list`,
	}
	for _, antwort := range tests {
		t.Run(antwort, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				_, _ = w.Write([]byte(antwort))
			}))
			defer server.Close()

			_, err := NewZettelstoreClient(server.URL).GetZettel()
			if err == nil || !strings.Contains(err.Error(), "ungueltige Zettelstore-Antwort") {
				t.Fatalf("erwarteter Parserfehler, erhalten: %v", err)
			}
		})
	}
}

func TestZettelstoreClientNichtErreichbar(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {}))
	adresse := server.URL
	server.Close()

	_, err := NewZettelstoreClient(adresse).GetZettel()
	if err == nil || !strings.Contains(err.Error(), "nicht erreichbar") {
		t.Fatalf("erwarteter Verbindungsfehler, erhalten: %v", err)
	}
}

func TestZettelstoreClientLehntUngueltigeAdresseAb(t *testing.T) {
	client := NewZettelstoreClient("file:///tmp/zettel")
	if _, err := client.GetZettel(); err == nil || !strings.Contains(err.Error(), "nur http und https") {
		t.Fatalf("erwarteter Adressfehler, erhalten: %v", err)
	}
	if err := client.SetAdresse("http://user:pass@127.0.0.1:23123"); err == nil {
		t.Fatal("Adresse mit Zugangsdaten wurde akzeptiert")
	}
	if err := client.SetAdresse("http://127.0.0.1:23123?fremd=1"); err == nil {
		t.Fatal("Adresse mit Query wurde akzeptiert")
	}
}
