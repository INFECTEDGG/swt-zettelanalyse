package application

import (
	"sort"
	"time"

	"github.com/INFECTEDGG/swt-zettelanalyse/model"
)

// ZettelProMonat zaehlt Zettel eines Kalenderjahres. Der Array-Index 0 steht
// fuer Januar; dadurch sind auch Monate ohne Zettel immer enthalten.
func ZettelProMonat(zettel []model.Zettel, jahr int) [12]int {
	var ergebnis [12]int
	for _, z := range zettel {
		if z.Erstellt.IsZero() || z.Erstellt.Year() != jahr {
			continue
		}
		ergebnis[int(z.Erstellt.Month())-1]++
	}
	return ergebnis
}

// AktivsterMonat liefert bei Gleichstand den fruehesten Monat. Fuer ein Jahr
// ohne Zettel wird ok=false geliefert.
func AktivsterMonat(zettel []model.Zettel, jahr int) (wert model.JahrMonat, ok bool) {
	werte := ZettelProMonat(zettel, jahr)
	for i, anzahl := range werte {
		if anzahl > wert.Anzahl {
			wert = model.JahrMonat{Jahr: jahr, Monat: time.Month(i + 1), Anzahl: anzahl}
			ok = true
		}
	}
	return wert, ok
}

// VerfuegbareJahre liefert alle Jahre mit gueltigem Erstellungszeitpunkt in
// absteigender Reihenfolge.
func VerfuegbareJahre(zettel []model.Zettel) []int {
	menge := make(map[int]struct{})
	for _, z := range zettel {
		if !z.Erstellt.IsZero() {
			menge[z.Erstellt.Year()] = struct{}{}
		}
	}
	jahre := make([]int, 0, len(menge))
	for jahr := range menge {
		jahre = append(jahre, jahr)
	}
	sort.Sort(sort.Reverse(sort.IntSlice(jahre)))
	return jahre
}
