// main erzeugt die Komponenten und verbindet sie (ADR 01).
package main

import (
	"flag"

	"Zettelanalyse/data"
	"Zettelanalyse/gui"
)

func main() {
	// Adresse des Zettelstores, änderbar beim Start (Issue #1),
	// z. B. go run . -adresse http://192.168.0.10:23123
	adresse := flag.String("adresse", "http://127.0.0.1:23123", "Adresse des Zettelstores")
	flag.Parse()

	quelle := data.NewZettelstoreClient(*adresse)
	gui.StartWindow(quelle)
}
