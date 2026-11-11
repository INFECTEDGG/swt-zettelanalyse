package main

import (
	"github.com/INFECTEDGG/swt-zettelanalyse/application"
	"github.com/INFECTEDGG/swt-zettelanalyse/data"
	"github.com/INFECTEDGG/swt-zettelanalyse/gui"
)

func main() {
	client := data.NewZettelstoreClient(data.Standardadresse)
	app := application.NewApp(client)
	gui.Start(app, data.Standardadresse)
}
