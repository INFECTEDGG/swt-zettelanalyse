package application

import (
	"fmt"
	"os"

	"Zettelanalyse/model"
)

func Wtemp(q []model.Zettel) {
	var latest *model.Zettel
	var latestDate string
	for i := range q {
		date := q[i].Modified
		if date == "" {
			date = q[i].Created
		}
		if date == "" {
			continue
		}
		if latest == nil || date > latestDate {
			latest = &q[i]
			latestDate = date
		}
	}

	if latest == nil {
		fmt.Println("Kein Zettel mit Modified- oder Created-Datum gefunden.")
		return
	}

	f, err := os.CreateTemp("", "zettelanalyse")
	if err != nil {
		panic(err)
	}

	_, err = f.Write([]byte(latest.Modified))

}
