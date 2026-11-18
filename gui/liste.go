// Package gui enthält das Programmfenster. Nur dieses Paket importiert tk9.0 (ADR 01).
package gui

import (
	"fmt"
	"strings"

	"Zettelanalyse/application"

	. "modernc.org/tk9.0"
)

// StartWindow öffnet das Fenster und zeigt die Zettel aus dem Zettelstore an.
func StartWindow(source application.ZettelQuelle) {
	App.WmTitle("Zettelanalyse")
	WmGeometry(App, "1100x700")

	table := TTreeview(Columns("id title tags"), Show("headings"))
	table.Heading("id", Txt("ID"))
	table.Heading("title", Txt("Titel"))
	table.Heading("tags", Txt("Tags"))
	tagTable := TTreeview(Columns("tag count"), Show("headings"), Height(8))
	tagTable.Heading("tag", Txt("Tag"))
	tagTable.Heading("count", Txt("Anzahl"))

	// Zettel über die Anwendungslogik laden (ohne Systemzettel)
	zettels, err := application.LadeZettel(source)
	status := fmt.Sprintf("%d Zettel geladen", len(zettels))
	if err != nil {
		status = "Fehler: " + err.Error()
	}
	for _, note := range zettels {
		table.Insert("", "end", Values([]string{note.ID, note.Title, strings.Join(note.Tags, " ")}))
	}
	for _, tag := range application.MostFrequentTags(zettels) {
		tagTable.Insert("", "end", Values([]string{tag.Tag, fmt.Sprint(tag.Count)}))
	}

	Pack(table, Expand(true), Fill("both"))
	Pack(Label(Txt("Häufigste Tags")), Anchor("w"), Padx("2m"), Pady("2m 0"))
	Pack(tagTable, Fill("x"), Padx("2m"))
	Pack(Label(Txt(status)), Anchor("w"), Padx("2m"), Pady("1m"))

	App.Center().Wait()
}
