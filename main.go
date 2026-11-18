// main erzeugt die Komponenten und verbindet sie (ADR 01).
package main

import (
	"Zettelanalyse/data"
	"Zettelanalyse/gui"
)

func main() {
	// Verbindung zum lokalen Zettelstore herstellen
	quelle := data.NewZettelstoreClient()

	// Benutzeroberfläche starten
	gui.StartWindow(quelle)
}
