// Package gui enthält das Programmfenster. Nur dieses Paket importiert tk9.0 (ADR 01).
package gui

import (
	"fmt"
	"strings"

	"Zettelanalyse/application"

	. "modernc.org/tk9.0"
)

// StartWindow öffnet das Fenster und zeigt die Zettel aus dem Zettelstore an.
func StartWindow(quelle application.ZettelQuelle) {
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
	zettel, err := application.LadeZettel(quelle)
	status := fmt.Sprintf("%d Zettel geladen", len(zettel))
	if err != nil {
		status = "Fehler: " + err.Error()
	}
	for _, z := range zettel {
		table.Insert("", "end", Values([]string{z.ID, z.Title, strings.Join(z.Tags, " ")}))
	}
	for _, tag := range application.HaeufigsteTags(zettel) {
		tagTable.Insert("", "end", Values([]string{tag.Tag, fmt.Sprint(tag.Anzahl)}))
	}

	Pack(table, Expand(true), Fill("both"))
	Pack(Label(Txt("Häufigste Tags")), Anchor("w"), Padx("2m"), Pady("2m 0"))
	Pack(tagTable, Fill("x"), Padx("2m"))
	Pack(Label(Txt(status)), Anchor("w"), Padx("2m"), Pady("1m"))

	App.Center().Wait()
}
