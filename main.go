package main

import (
	"flag"

	"Zettelanalyse/application"
	"Zettelanalyse/data"
	"Zettelanalyse/gui"
)

func main() {
	// Allow the Zettelstore address to be configured at startup.
	address := flag.String(
		"adresse",
		"http://127.0.0.1:23123",
		"Zettelstore server address",
	)
	flag.Parse()

	// Create the Zettelstore client.
	source := data.NewZettelstoreClient(*address)

	// Start the graphical user interface.
	gui.StartWindow(source, *address, func(newAddress string) application.ZettelQuelle {
		return data.NewZettelstoreClient(newAddress)
	})
}
