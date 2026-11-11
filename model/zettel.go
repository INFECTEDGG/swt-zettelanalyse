package model

// Zettel speichert die Daten eines Zettels aus dem Zettelstore
type Zettel struct {
	id       string
	title    string
	tags     []string
	created  string // Metadaten "created": wann der Zettel erstellt wurde, z. B. 20251021093000
	modified string // Metadaten "modified": wann der Zettel zuletzt geändert wurde (leer, wenn nie geändert)
}
