package analysis

import (
	"math"
	"sort"

	"github.com/jb843051627/foram-bench/internal/model"
)

type Comparison struct {
	Taxon       string  `json:"taxon"`
	LeftCount   int     `json:"left_count"`
	RightCount  int     `json:"right_count"`
	Difference  int     `json:"difference"`
	RelativeGap float64 `json:"relative_gap"`
}

func Compare(left, right []model.Observation) []Comparison {
	leftCounts := counts(left)
	rightCounts := counts(right)
	seen := make(map[string]struct{}, len(leftCounts)+len(rightCounts))
	for taxon := range leftCounts {
		seen[taxon] = struct{}{}
	}
	for taxon := range rightCounts {
		seen[taxon] = struct{}{}
	}
	result := make([]Comparison, 0, len(seen))
	for taxon := range seen {
		l, r := leftCounts[taxon], rightCounts[taxon]
		gap := 0.0
		if l+r > 0 {
			gap = math.Abs(float64(l-r)) / float64(l+r) * 100
		}
		result = append(result, Comparison{Taxon: taxon, LeftCount: l, RightCount: r, Difference: r - l, RelativeGap: gap})
	}
	sort.Slice(result, func(i, j int) bool { return result[i].Taxon < result[j].Taxon })
	return result
}

func counts(observations []model.Observation) map[string]int {
	result := make(map[string]int)
	for _, observation := range observations {
		result[observation.Taxon] += observation.Count
	}
	return result
}

func SignificantChange(comparisons []Comparison, threshold float64) []Comparison {
	result := make([]Comparison, 0)
	for _, item := range comparisons {
		if item.RelativeGap >= threshold {
			result = append(result, item)
		}
	}
	return result
}
