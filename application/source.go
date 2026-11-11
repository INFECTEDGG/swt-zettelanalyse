package application

import (
	"fmt"
	"sort"

	"github.com/INFECTEDGG/swt-zettelanalyse/model"
)

// ZettelQuelle wird von der Datenschicht und in Tests von einem Mock
// implementiert. Die Anwendungsschicht kennt die konkrete Implementierung
// nicht.
type ZettelQuelle interface {
	GetZettel() ([]model.Zettel, error)
}

// AdressierbareZettelQuelle ist die optionale Faehigkeit einer Quelle, ihr
// Ziel zur Laufzeit zu wechseln.
type AdressierbareZettelQuelle interface {
	ZettelQuelle
	SetAdresse(string) error
}

// LadeZettel holt alle Zettel aus der Quelle, entfernt Systemzettel und
// Duplikate und liefert eine deterministisch sortierte Liste.
func LadeZettel(q ZettelQuelle) ([]model.Zettel, error) {
	if q == nil {
		return nil, fmt.Errorf("keine Zettelquelle konfiguriert")
	}
	alle, err := q.GetZettel()
	if err != nil {
		return nil, fmt.Errorf("Zettel konnten nicht geladen werden: %w", err)
	}

	ergebnis := make([]model.Zettel, 0, len(alle))
	gesehen := make(map[string]struct{}, len(alle))
	for _, z := range alle {
		if IstSystemzettel(z.ID) || z.Rolle == "manual" {
			continue
		}
		if _, vorhanden := gesehen[z.ID]; vorhanden {
			continue
		}
		gesehen[z.ID] = struct{}{}
		z.Tags = append([]string(nil), z.Tags...)
		ergebnis = append(ergebnis, z)
	}

	sort.SliceStable(ergebnis, func(i, j int) bool {
		if ergebnis[i].Erstellt.Equal(ergebnis[j].Erstellt) {
			return ergebnis[i].ID < ergebnis[j].ID
		}
		if ergebnis[i].Erstellt.IsZero() {
			return false
		}
		if ergebnis[j].Erstellt.IsZero() {
			return true
		}
		return ergebnis[i].Erstellt.Before(ergebnis[j].Erstellt)
	})
	return ergebnis, nil
}

// IstSystemzettel erhaelt die bereits im Projekt verwendete Grenze. Die
// reservierten Zettel und das mitgelieferte Handbuch liegen in diesem
// ID-Bereich; regulaere zeitbasierte Zettel-IDs liegen darueber.
func IstSystemzettel(id string) bool {
	return id <= "00010000000000"
}
