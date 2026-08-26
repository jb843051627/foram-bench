package analysis

import (
	"math"
	"sort"

	"github.com/jb843051627/foram-bench/internal/model"
)

type Abundance struct {
	Taxon    string  `json:"taxon"`
	Count    int     `json:"count"`
	Relative float64 `json:"relative"`
	Rank     int     `json:"rank"`
}

func Abundances(observations []model.Observation) []Abundance {
	counts := make(map[string]int)
	total := 0
	for _, observation := range observations {
		if observation.Count <= 0 || observation.Taxon == "" {
			continue
		}
		counts[observation.Taxon] += observation.Count
		total += observation.Count
	}
	result := make([]Abundance, 0, len(counts))
	for taxon, count := range counts {
		relative := 0.0
		if total > 0 {
			relative = float64(count) / float64(total) * 100
		}
		result = append(result, Abundance{Taxon: taxon, Count: count, Relative: relative})
	}
	sort.Slice(result, func(i, j int) bool {
		if result[i].Count == result[j].Count {
			return result[i].Taxon < result[j].Taxon
		}
		return result[i].Count > result[j].Count
	})
	for index := range result {
		result[index].Rank = index + 1
	}
	return result
}

func Shannon(observations []model.Observation) float64 {
	items := Abundances(observations)
	result := 0.0
	for _, item := range items {
		if item.Relative <= 0 {
			continue
		}
		p := item.Relative / 100
		result -= p * math.Log(p)
	}
	return result
}

func Dominant(observations []model.Observation) (Abundance, bool) {
	items := Abundances(observations)
	if len(items) == 0 {
		return Abundance{}, false
	}
	return items[0], true
}

func Evenness(observations []model.Observation) float64 {
	items := Abundances(observations)
	if len(items) <= 1 {
		return 0
	}
	return Shannon(observations) / math.Log(float64(len(items)))
}
