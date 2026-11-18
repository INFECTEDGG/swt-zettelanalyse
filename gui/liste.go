// Package gui enthält das Programmfenster. Nur dieses Paket importiert tk9.0 (ADR 01).
package gui

import (
	"fmt"
	"net/url"
	"strings"
	"time"

	"Zettelanalyse/application"

	. "modernc.org/tk9.0"
)

// loadZettel leert beide Tabellen, lädt die Zettel neu und gibt den Statustext zurück.
func loadZettel(table, tagTable *TTreeviewWidget, source application.ZettelQuelle) string {
	// alte Zeilen der Zettel-Tabelle entfernen
	for _, item := range table.Children("") {
		table.Delete(item)
	}
	// alte Zeilen der Tag-Tabelle entfernen
	for _, item := range tagTable.Children("") {
		tagTable.Delete(item)
	}

	// Zettel über die Anwendungslogik laden (ohne Systemzettel)
	zettel, err := application.LadeZettel(source)
	if err != nil {
		return "Fehler: " + err.Error()
	}

	// Zettel-Tabelle füllen
	for _, z := range zettel {
		table.Insert("", "end", Values([]string{z.ID, z.Title, strings.Join(z.Tags, " ")}))
	}
	// Tag-Tabelle mit den häufigsten Tags füllen
	for _, tag := range application.MostFrequentTags(zettel) {
		tagTable.Insert("", "end", Values([]string{tag.Tag, fmt.Sprint(tag.Count)}))
	}

	// Zettel des aktuellen Monats zählen (Funktion aus der Anwendungslogik)
	now := time.Now()
	monthZettel := application.ZettelInMonth(zettel, now.Year(), int(now.Month()))

	return fmt.Sprintf("%d Zettel geladen, davon %d im %02d/%d", len(zettel), len(monthZettel), int(now.Month()), now.Year())
}

// StartWindow öffnet das Fenster und zeigt die Zettel aus dem Zettelstore an.
// address -> aktuelle Adresse des Zettelstores
// newSource baut aus einer Adresse eine neue ZettelQuelle
func StartWindow(source application.ZettelQuelle, address string, newSource func(string) application.ZettelQuelle) {
	App.WmTitle("Zettelanalyse")
	WmGeometry(App, "1100x700")

	table := TTreeview(Columns("id title tags"), Show("headings"))
	table.Heading("id", Txt("ID"))
	table.Heading("title", Txt("Titel"))
	table.Heading("tags", Txt("Tags"))

	tagTable := TTreeview(Columns("tag count"), Show("headings"), Height(8))
	tagTable.Heading("tag", Txt("Tag"))
	tagTable.Heading("count", Txt("Anzahl"))

	// statusLabel zeigt das Ergebnis des ersten Ladens
	statusLabel := Label(Txt(loadZettel(table, tagTable, source)))

	// bar enthält Adressfeld und Button in einer Zeile
	bar := TFrame()
	entry := bar.TEntry(Textvariable(address))

	refresh := bar.TButton(Txt("Aktualisieren"), Command(func() {
		text := strings.TrimSpace(entry.Textvariable())
		// ungültige Adresse: Meldung zeigen und abbrechen
		if _, err := url.ParseRequestURI(text); err != nil {
			statusLabel.Configure(Txt("Ungültige Adresse: " + text))
			return
		}
		// nur bei geänderter Adresse eine neue Quelle bauen
		if text != address {
			address = text
			source = newSource(text)
		}
		statusLabel.Configure(Txt(loadZettel(table, tagTable, source)))
	}))

	Pack(entry, Side("left"), Expand(true), Fill("x"), Padx("2m"))
	Pack(refresh, Side("left"), Padx("2m"))
	Pack(bar, Fill("x"), Pady("1m"))
	Pack(table, Expand(true), Fill("both"))
	Pack(Label(Txt("Häufigste Tags")), Anchor("w"), Padx("2m"), Pady("2m 0"))
	Pack(tagTable, Fill("x"), Padx("2m"))
	Pack(statusLabel, Anchor("w"), Padx("2m"), Pady("1m"))

	App.Center().Wait()
}
