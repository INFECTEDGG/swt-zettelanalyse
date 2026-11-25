package application

import (
	"fmt"
	"os"
	"path/filepath"

	"Zettelanalyse/model"
)

func Wtemp(q []model.Zettel) {
	var latest *model.Zettel
	var latestDate string
	for i := range q {
		date := q[i].Published
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

	path := filepath.Join(os.TempDir(), "zettelanalyse_last.txt")
	err := os.WriteFile(path, []byte(latest.Published), 0644)

	print(err)
}
