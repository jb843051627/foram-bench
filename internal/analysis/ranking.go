package analysis

import (
	"sort"

	"github.com/jb843051627/foram-bench/internal/model"
)

type RankedSection struct {
	Section          model.ThinSection
	ObservationCount int
	TaxaCount        int
	MeanConfidence   float64
	Priority         int
}

func RankSections(sections []model.ThinSection, observations []model.Observation) []RankedSection {
	bySection := make(map[string][]model.Observation)
	for _, observation := range observations {
		bySection[observation.SectionID] = append(bySection[observation.SectionID], observation)
	}
	result := make([]RankedSection, 0, len(sections))
	for _, section := range sections {
		items := bySection[section.ID]
		taxa := make(map[string]struct{})
		confidence := 0.0
		for _, item := range items {
			taxa[item.Taxon] = struct{}{}
			confidence += item.Confidence
		}
		if len(items) > 0 {
			confidence /= float64(len(items))
		}
		priority := len(taxa)*10 + len(items)
		if confidence < 0.6 {
			priority += 20
		}
		result = append(result, RankedSection{Section: section, ObservationCount: len(items), TaxaCount: len(taxa), MeanConfidence: confidence, Priority: priority})
	}
	sort.SliceStable(result, func(i, j int) bool {
		if result[i].Priority == result[j].Priority {
			return result[i].Section.Label < result[j].Section.Label
		}
		return result[i].Priority > result[j].Priority
	})
	return result
}

func TopTaxa(observations []model.Observation, limit int) []Abundance {
	items := Abundances(observations)
	if limit <= 0 || limit >= len(items) {
		return items
	}
	return append([]Abundance(nil), items[:limit]...)
}
