package application

import (
	"fmt"

	"Zettelanalyse/model"
)

func ZettelProMonat(zettel []model.Zettel, jahr string) []int {
	anzahl := []int{}
	for monat := 1; monat <= 12; monat++ {
		// Monat zweistellig machen: 1 -> "01", 11 -> "11"
		suche := jahr + fmt.Sprintf("%02d", monat)
		imMonat := sumZettelMonth(zettel, suche)
		anzahl = append(anzahl, len(imMonat))
	}
	return anzahl
}
