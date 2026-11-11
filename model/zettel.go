package model

import "time"

// Zettel enthaelt die fuer die Auswertung benoetigten Metadaten eines Zettels.
// Fehlende Zeitstempel werden als time.Time-Nullwert dargestellt.
type Zettel struct {
	ID       string
	Titel    string
	Tags     []string
	Erstellt time.Time
	Rolle    string
}

// JahrMonat beschreibt einen Monatswert innerhalb eines Kalenderjahres.
type JahrMonat struct {
	Jahr   int
	Monat  time.Month
	Anzahl int
}

// TagAnzahl enthaelt einen Tag und die Zahl der Zettel, die ihn verwenden.
type TagAnzahl struct {
	Tag    string
	Anzahl int
}
