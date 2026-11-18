package application

import (
	"sort"
	"strings"

	"Zettelanalyse/model"
)

// MostFrequentTags counts each tag at most once per note. Results are sorted
// by count in descending order and alphabetically when counts are equal.
func MostFrequentTags(notes []model.Zettel) []model.TagCount {
	counts := make(map[string]int)
	for _, note := range notes {
		seenTags := make(map[string]struct{}, len(note.Tags))
		for _, tag := range note.Tags {
			tag = strings.TrimSpace(tag)
			if tag == "" {
				continue
			}
			if _, exists := seenTags[tag]; exists {
				continue
			}
			seenTags[tag] = struct{}{}
			counts[tag]++
		}
	}

	result := make([]model.TagCount, 0, len(counts))
	for tag, count := range counts {
		result = append(result, model.TagCount{Tag: tag, Count: count})
	}
	sort.Slice(result, func(i, j int) bool {
		if result[i].Count == result[j].Count {
			return result[i].Tag < result[j].Tag
		}
		return result[i].Count > result[j].Count
	})
	return result
}
