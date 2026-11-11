package model

import "time"

// Zettel speichert die Daten eines Zettels aus dem Zettelstore
type Zettel struct {
	ID    string
	Titel string
	Tags  []string
}

// MonatsAnzahl speichert, wie viele Zettel in einem Monat erstellt wurden
type MonatsAnzahl struct {
	Jahr   int
	Monat  time.Month
	Anzahl int
}

// Erstellt gibt das Datum zurück, das in der ID steckt.
// Die ID ist so aufgebaut: JahrMonatTagStundeMinuteSekunde (z. B. 20251021093000).
// Wenn die ID kein gültiges Datum ist, wird false zurückgegeben.
func (z Zettel) Erstellt() (time.Time, bool) {
	datum, err := time.Parse("20060102150405", z.ID)
	if err != nil {
		return time.Time{}, false
	}
	return datum, true
}
