package analysis

import (
	"strings"
	"unicode"

	"github.com/jb843051627/foram-bench/internal/model"
)

func NormalizeTaxon(value string) string {
	words := strings.FieldsFunc(strings.TrimSpace(value), func(r rune) bool {
		return unicode.IsSpace(r) || r == '_' || r == '-'
	})
	if len(words) == 0 {
		return ""
	}
	for index, word := range words {
		if word == "" {
			continue
		}
		word = strings.ToLower(word)
		words[index] = strings.ToUpper(word[:1]) + word[1:]
	}
	return strings.Join(words, " ")
}

func NormalizeObservations(items []model.Observation) []model.Observation {
	result := make([]model.Observation, len(items))
	copy(result, items)
	for index := range result {
		result[index].Taxon = NormalizeTaxon(result[index].Taxon)
		result[index].Preservation = strings.ToLower(strings.TrimSpace(result[index].Preservation))
		result[index].Notes = strings.TrimSpace(result[index].Notes)
	}
	return result
}

func MergeDuplicateTaxa(items []model.Observation) []model.Observation {
	result := make(map[string]model.Observation)
	for _, item := range NormalizeObservations(items) {
		key := item.SectionID + "\x00" + item.Taxon
		current, ok := result[key]
		if !ok {
			result[key] = item
			continue
		}
		current.Count += item.Count
		if item.Confidence > current.Confidence {
			current.Confidence = item.Confidence
		}
		result[key] = current
	}
	output := make([]model.Observation, 0, len(result))
	for _, item := range result {
		output = append(output, item)
	}
	return output
}
