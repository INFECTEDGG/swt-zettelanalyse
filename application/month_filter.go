package application

import (
	"sort"
	"time"

	"github.com/INFECTEDGG/swt-zettelanalyse/model"
)

// FilterNachMonat liefert eine Kopie aller Zettel des angegebenen
// Kalendermonats, chronologisch und anschliessend nach ID sortiert.
func FilterNachMonat(zettel []model.Zettel, jahr int, monat time.Month) []model.Zettel {
	if monat < time.January || monat > time.December {
		return []model.Zettel{}
	}
	ergebnis := make([]model.Zettel, 0)
	for _, z := range zettel {
		if !z.Erstellt.IsZero() && z.Erstellt.Year() == jahr && z.Erstellt.Month() == monat {
			z.Tags = append([]string(nil), z.Tags...)
			ergebnis = append(ergebnis, z)
		}
	}
	sort.SliceStable(ergebnis, func(i, j int) bool {
		if ergebnis[i].Erstellt.Equal(ergebnis[j].Erstellt) {
			return ergebnis[i].ID < ergebnis[j].ID
		}
		return ergebnis[i].Erstellt.Before(ergebnis[j].Erstellt)
	})
	return ergebnis
}
