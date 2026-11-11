package model

// Zettel speichert die Daten eines Zettels aus dem Zettelstore
type Zettel struct {
	ID       string
	Titel    string
	Tags     []string
	Created  string // Metadaten "created": wann der Zettel erstellt wurde, z. B. 20251021093000
	Modified string // Metadaten "modified": wann der Zettel zuletzt geändert wurde (leer, wenn nie geändert)
}
