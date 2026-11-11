package application

import (
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/INFECTEDGG/swt-zettelanalyse/model"
)

type mockQuelle struct {
	zettel  []model.Zettel
	err     error
	adresse string
}

func (m *mockQuelle) GetZettel() ([]model.Zettel, error) {
	return kopiereZettel(m.zettel), m.err
}

func (m *mockQuelle) SetAdresse(adresse string) error {
	if adresse == "ungueltig" {
		return errors.New("ungueltige Testadresse")
	}
	m.adresse = adresse
	return nil
}

func zeit(wert string) time.Time {
	zeitpunkt, err := time.Parse("20060102150405", wert)
	if err != nil {
		panic(err)
	}
	return zeitpunkt
}

func TestIstSystemzettel(t *testing.T) {
	tests := []struct {
		id   string
		will bool
	}{
		{"00000000000001", true},
		{"00001012051200", true},
		{"00010000000000", true},
		{"00010000000001", false},
		{"20261111120000", false},
	}
	for _, test := range tests {
		if erhalten := IstSystemzettel(test.id); erhalten != test.will {
			t.Errorf("IstSystemzettel(%q) = %v, erwartet %v", test.id, erhalten, test.will)
		}
	}
}

func TestLadeZettelFiltertSystemzettelUndDuplikate(t *testing.T) {
	quelle := &mockQuelle{zettel: []model.Zettel{
		{ID: "00001012051200", Titel: "Handbuch"},
		{ID: "20260102090000", Titel: "Externes Handbuch", Rolle: "manual", Erstellt: zeit("20260102090000")},
		{ID: "20260201090000", Titel: "Zweiter", Erstellt: zeit("20260201090000")},
		{ID: "20260101090000", Titel: "Erster", Erstellt: zeit("20260101090000")},
		{ID: "20260101090000", Titel: "Duplikat", Erstellt: zeit("20260101090000")},
	}}

	erhalten, err := LadeZettel(quelle)
	if err != nil {
		t.Fatalf("LadeZettel() Fehler: %v", err)
	}
	if len(erhalten) != 2 {
		t.Fatalf("LadeZettel() lieferte %d statt 2 Zettel: %#v", len(erhalten), erhalten)
	}
	if erhalten[0].ID != "20260101090000" || erhalten[1].ID != "20260201090000" {
		t.Fatalf("unerwartete Reihenfolge: %#v", erhalten)
	}
}

func TestZettelProMonatBeruecksichtigtAlleMonate(t *testing.T) {
	zettel := []model.Zettel{
		{ID: "1", Erstellt: zeit("20260101120000")},
		{ID: "2", Erstellt: zeit("20260131120000")},
		{ID: "3", Erstellt: zeit("20260310120000")},
		{ID: "4", Erstellt: zeit("20250310120000")},
		{ID: "5"},
	}
	erwartet := [12]int{2, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	if erhalten := ZettelProMonat(zettel, 2026); erhalten != erwartet {
		t.Fatalf("ZettelProMonat() = %v, erwartet %v", erhalten, erwartet)
	}
}

func TestAktivsterMonat(t *testing.T) {
	zettel := []model.Zettel{
		{Erstellt: zeit("20260101120000")},
		{Erstellt: zeit("20260301120000")},
	}
	erhalten, ok := AktivsterMonat(zettel, 2026)
	if !ok || erhalten.Monat != time.January || erhalten.Anzahl != 1 {
		t.Fatalf("Gleichstand muss den fruehesten Monat liefern: %#v, ok=%v", erhalten, ok)
	}
	if wert, ok := AktivsterMonat(nil, 2026); ok || wert.Anzahl != 0 {
		t.Fatalf("leere Daten sollten ok=false liefern: %#v, ok=%v", wert, ok)
	}
}

func TestHaeufigsteTagsSortiertDeterministisch(t *testing.T) {
	zettel := []model.Zettel{
		{Tags: []string{"#go", "#uni", "#go", ""}},
		{Tags: []string{"#go", "#projekt"}},
		{Tags: []string{"#uni", "#projekt"}},
	}
	erwartet := []model.TagAnzahl{
		{Tag: "#go", Anzahl: 2},
		{Tag: "#projekt", Anzahl: 2},
		{Tag: "#uni", Anzahl: 2},
	}
	if erhalten := HaeufigsteTags(zettel); !reflect.DeepEqual(erhalten, erwartet) {
		t.Fatalf("HaeufigsteTags() = %#v, erwartet %#v", erhalten, erwartet)
	}
}

func TestFilterNachMonat(t *testing.T) {
	zettel := []model.Zettel{
		{ID: "b", Erstellt: zeit("20260620120000")},
		{ID: "a", Erstellt: zeit("20260601120000")},
		{ID: "c", Erstellt: zeit("20260701120000")},
		{ID: "d", Erstellt: zeit("20250601120000")},
	}
	erhalten := FilterNachMonat(zettel, 2026, time.June)
	if len(erhalten) != 2 || erhalten[0].ID != "a" || erhalten[1].ID != "b" {
		t.Fatalf("FilterNachMonat() = %#v", erhalten)
	}
	if erhalten := FilterNachMonat(zettel, 2026, 13); len(erhalten) != 0 {
		t.Fatalf("ungueltiger Monat lieferte Zettel: %#v", erhalten)
	}
}

func TestAppAktualisierenVerwendetMockQuelle(t *testing.T) {
	quelle := &mockQuelle{zettel: []model.Zettel{{ID: "20260101120000", Erstellt: zeit("20260101120000")}}}
	app := NewApp(quelle)
	anzahl, err := app.Aktualisieren("http://127.0.0.1:23123")
	if err != nil || anzahl != 1 {
		t.Fatalf("Aktualisieren() = %d, %v", anzahl, err)
	}
	if quelle.adresse != "http://127.0.0.1:23123" {
		t.Fatalf("Adresse wurde nicht gesetzt: %q", quelle.adresse)
	}
}
