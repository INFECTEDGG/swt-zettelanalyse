package application

import (
	"fmt"

	"Zettelanalyse/model"
)

func ZettelperMonat(zettel []model.Zettel, year string) []int {
	count := []int{}
	for month := 1; month <= 12; month++ {
		// Monat zweistellig machen: 1 -> "01", 11 -> "11"
		month := year + fmt.Sprintf("%02d", count)
		result := sumZettelMonth(zettel, month)
		count = append(count, len(result))
	}
	return count
}
