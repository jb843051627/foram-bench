package analysis

import (
	"sort"
	"strings"

	"github.com/jb843051627/foram-bench/internal/model"
)

type Evidence struct {
	SectionID    string
	Taxon        string
	Count        int
	Confidence   float64
	Preservation string
	ObservedAt   string
}

type EvidenceSet struct {
	Items []Evidence
}

func FromObservations(items []model.Observation) EvidenceSet {
	result := EvidenceSet{Items: make([]Evidence, 0, len(items))}
	for _, item := range items {
		result.Items = append(result.Items, Evidence{SectionID: item.SectionID, Taxon: strings.TrimSpace(item.Taxon), Count: item.Count,
			Confidence: item.Confidence, Preservation: strings.TrimSpace(item.Preservation), ObservedAt: item.ObservedAt.Format("2006-01-02T15:04:05Z07:00")})
	}
	return result
}

func (s EvidenceSet) Sorted() []Evidence {
	result := make([]Evidence, len(s.Items))
	copy(result, s.Items)
	sort.SliceStable(result, func(i, j int) bool {
		if result[i].Taxon == result[j].Taxon {
			return result[i].Count > result[j].Count
		}
		return result[i].Taxon < result[j].Taxon
	})
	return result
}

func (s EvidenceSet) TotalCount() int {
	total := 0
	for _, item := range s.Items {
		total += item.Count
	}
	return total
}

func (s EvidenceSet) SectionCount() int {
	seen := make(map[string]struct{})
	for _, item := range s.Items {
		seen[item.SectionID] = struct{}{}
	}
	return len(seen)
}

func (s EvidenceSet) Taxa() []string {
	seen := make(map[string]struct{})
	for _, item := range s.Items {
		if item.Taxon != "" {
			seen[item.Taxon] = struct{}{}
		}
	}
	result := make([]string, 0, len(seen))
	for taxon := range seen {
		result = append(result, taxon)
	}
	sort.Strings(result)
	return result
}

func (s EvidenceSet) ForSection(id string) EvidenceSet {
	result := EvidenceSet{Items: make([]Evidence, 0)}
	for _, item := range s.Items {
		if item.SectionID == id {
			result.Items = append(result.Items, item)
		}
	}
	return result
}
