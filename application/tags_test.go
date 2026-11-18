package application

import (
	"reflect"
	"testing"

	"Zettelanalyse/model"
)

func TestMostFrequentTags(t *testing.T) {
	notes := []model.Zettel{
		{Tags: []string{"#go", "#uni", "#go"}},
		{Tags: []string{"#go", "#projekt"}},
		{Tags: []string{"#uni"}},
		{},
	}
	expected := []model.TagCount{
		{Tag: "#go", Count: 2},
		{Tag: "#uni", Count: 2},
		{Tag: "#projekt", Count: 1},
	}
	if actual := MostFrequentTags(notes); !reflect.DeepEqual(actual, expected) {
		t.Fatalf("MostFrequentTags() = %#v, expected %#v", actual, expected)
	}
}
