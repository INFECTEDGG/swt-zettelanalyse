package application

import "Zettelanalyse/model"

func sumZettelMonth(zettel []model.Zettel, month string) []model.Zettel {
	result := []model.Zettel{}
	for _, z := range zettel {

		timestamp := z.Published
		if timestamp == "" {
			timestamp = z.ID
		}

		if len(timestamp) < 6 {
			continue
		}

		if timestamp[0:6] == month {
			result = append(result, z)
		}
	}
	return result
}
