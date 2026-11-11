package application

import "Zettelanalyse/model"

// ZettelQuelle wird von der Datenschicht implementiert (data.ZettelstoreClient).
type ZettelQuelle interface {
	GetZettel() ([]model.Zettel, error)
}

// LadeZettel holt alle Zettel aus der Quelle und entfernt die Systemzettel.
func LadeZettel(q ZettelQuelle) ([]model.Zettel, error) {
	alle, err := q.GetZettel()
	if err != nil {
		return nil, err
	}
	var ergebnis []model.Zettel
	for _, z := range alle {
		if !IstSystemzettel(z.id) {
			ergebnis = append(ergebnis, z)
		}
	}
	return ergebnis, nil
}

// Systemzettel und Handbuch haben IDs mit 0000..., die Startseite 00010000000000.
func IstSystemzettel(id string) bool {
	return id <= "00010000000000"
}
