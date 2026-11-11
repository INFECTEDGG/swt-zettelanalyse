package window

import (
	//Go Standardbibliothek für Hilfsfunktionen Texte
	"strings"

	//TK Bibliothek für Fenster und GUI Elemente
	. "modernc.org/tk9.0"
)

// Platzhalter für zettel weil noch kein api gemacht

type Zettel struct {
	Titel string
	Tags  []string
}

// Beispiel Zettel mit Inhalt
var placeholder = []Zettel{
	{"Projektauftrag Zettelanalyse", []string{"#uni", "#swt"}},
	{"Notizen zur Vorlesung", []string{"#uni"}},
	{"Rezept Linsencurry", []string{"#essen", "#rezept"}},
}

// Öffnet das Fenster
func StartWindow() {

	//Window manager + initial Größe, Fenstereinstellungen
	App.WmTitle("Zettelanalyse")
	WmGeometry(App, "1100x700")
	table := TTreeview(Columns("title tags"), Show("headings"))

	table.Heading("title", Txt("Titel"))
	table.Heading("tags", Txt("Tags"))

	//geht die Platzhalterliste durch und fügt die Zettel in die Tabelle ein
	for _, z := range placeholder {
		tagText := strings.Join(z.Tags, " ")
		table.Insert("", "end", Values([]string{z.Titel, tagText}))
	}

	//Packt die Tabelle in das Fenster und macht sie sichtbar
	Pack(table, Expand(true), Fill("both"))

	//Startet die GUI und wartet auf Eingaben
	App.Wait()
	App.Center().Wait()
}
