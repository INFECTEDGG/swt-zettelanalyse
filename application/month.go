package application

import (
	"strconv"

	"Zettelanalyse/model"
)

// ZettelInMonth gibt alle Zettel zurück, die im angegebenen Jahr und Monat erstellt wurden (Issue #7).
// month geht von 1 (Januar) bis 12 (Dezember).
func ZettelInMonth(zettel []model.Zettel, year int, month int) []model.Zettel {
	var result []model.Zettel
	for _, z := range zettel {
		timestamp := z.Published
		if timestamp == "" {
			timestamp = z.ID // Fallback auf die ID, falls Created leer ist
		}
		if len(timestamp) < 6 {
			continue
		}
		y, errYear := strconv.Atoi(timestamp[0:4])
		m, errMonth := strconv.Atoi(timestamp[4:6])
		if errYear == nil && errMonth == nil && y == year && m == month {
			result = append(result, z)
		}
	}
	return result
}
