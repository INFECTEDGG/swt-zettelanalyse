package application

import (
	"reflect"
	"testing"

	"Zettelanalyse/model"
)

func TestHaeufigsteTags(t *testing.T) {
	zettel := []model.Zettel{
		{Tags: []string{"#go", "#uni", "#go"}},
		{Tags: []string{"#go", "#projekt"}},
		{Tags: []string{"#uni"}},
		{},
	}
	erwartet := []model.TagAnzahl{
		{Tag: "#go", Anzahl: 2},
		{Tag: "#uni", Anzahl: 2},
		{Tag: "#projekt", Anzahl: 1},
	}
	if erhalten := HaeufigsteTags(zettel); !reflect.DeepEqual(erhalten, erwartet) {
		t.Fatalf("HaeufigsteTags() = %#v, erwartet %#v", erhalten, erwartet)
	}
}
