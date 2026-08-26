package analysis

import (
	"sort"

	"github.com/jb843051627/foram-bench/internal/model"
)

type Statistics struct {
	TotalObservations int     `json:"total_observations"`
	TotalCount        int     `json:"total_count"`
	MeanConfidence    float64 `json:"mean_confidence"`
	UniqueTaxa        int     `json:"unique_taxa"`
	ActiveSections    int     `json:"active_sections"`
	DamagedEvidence   int     `json:"damaged_evidence"`
}

func CalculateStatistics(sections []model.ThinSection, observations []model.Observation) Statistics {
	result := Statistics{TotalObservations: len(observations), ActiveSections: len(sections)}
	taxa := make(map[string]struct{})
	for _, observation := range observations {
		result.TotalCount += observation.Count
		result.MeanConfidence += observation.Confidence
		if observation.Taxon != "" {
			taxa[observation.Taxon] = struct{}{}
		}
		if observation.Preservation == "poor" || observation.Preservation == "fragmented" {
			result.DamagedEvidence++
		}
	}
	if result.TotalObservations > 0 {
		result.MeanConfidence /= float64(result.TotalObservations)
	}
	result.UniqueTaxa = len(taxa)
	return result
}

func Quantiles(observations []model.Observation) (float64, float64, float64) {
	values := make([]int, 0, len(observations))
	for _, observation := range observations {
		values = append(values, observation.Count)
	}
	if len(values) == 0 {
		return 0, 0, 0
	}
	sort.Ints(values)
	return percentile(values, 0.25), percentile(values, 0.5), percentile(values, 0.75)
}

func percentile(values []int, fraction float64) float64 {
	if len(values) == 1 {
		return float64(values[0])
	}
	position := fraction * float64(len(values)-1)
	lower := int(position)
	upper := lower + 1
	if upper >= len(values) {
		return float64(values[lower])
	}
	weight := position - float64(lower)
	return float64(values[lower])*(1-weight) + float64(values[upper])*weight
}
