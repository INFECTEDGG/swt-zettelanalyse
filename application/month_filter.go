package application

import "Zettelanalyse/model"

func ZettelImMonat(zettel []model.Zettel, monat string) []model.Zettel {
	ergebnis := []model.Zettel{}
	for _, z := range zettel {

		zeitpunkt := z.Created
		if zeitpunkt == "" {
			zeitpunkt = z.ID
		}

		if len(zeitpunkt) < 6 {
			continue
		}

		if zeitpunkt[0:6] == monat {
			ergebnis = append(ergebnis, z)
		}
	}
	return ergebnis
}
