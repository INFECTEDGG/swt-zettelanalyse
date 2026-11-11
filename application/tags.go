package application

import (
	"sort"
	"strings"

	"github.com/INFECTEDGG/swt-zettelanalyse/model"
)

// HaeufigsteTags zaehlt jeden Tag hoechstens einmal pro Zettel. Sortiert wird
// absteigend nach Anzahl und bei Gleichstand alphabetisch nach Tag.
func HaeufigsteTags(zettel []model.Zettel) []model.TagAnzahl {
	zaehler := make(map[string]int)
	for _, z := range zettel {
		imZettel := make(map[string]struct{}, len(z.Tags))
		for _, tag := range z.Tags {
			tag = strings.TrimSpace(tag)
			if tag == "" {
				continue
			}
			if _, vorhanden := imZettel[tag]; vorhanden {
				continue
			}
			imZettel[tag] = struct{}{}
			zaehler[tag]++
		}
	}

	ergebnis := make([]model.TagAnzahl, 0, len(zaehler))
	for tag, anzahl := range zaehler {
		ergebnis = append(ergebnis, model.TagAnzahl{Tag: tag, Anzahl: anzahl})
	}
	sort.Slice(ergebnis, func(i, j int) bool {
		if ergebnis[i].Anzahl == ergebnis[j].Anzahl {
			return ergebnis[i].Tag < ergebnis[j].Tag
		}
		return ergebnis[i].Anzahl > ergebnis[j].Anzahl
	})
	return ergebnis
}
