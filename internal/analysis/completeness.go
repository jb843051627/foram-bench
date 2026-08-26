package analysis

import (
	"fmt"
	"sort"

	"github.com/jb843051627/foram-bench/internal/model"
)

type CompletenessResult struct {
	Score    float64
	Expected int
	Present  int
	Missing  []string
	Findings []string
}

func Completeness(sections []model.ThinSection, observations []model.Observation, expectedTaxa []string) CompletenessResult {
	result := CompletenessResult{Expected: len(expectedTaxa), Missing: make([]string, 0)}
	seen := make(map[string]bool)
	for _, observation := range observations {
		if observation.Taxon != "" && observation.Count > 0 {
			seen[observation.Taxon] = true
		}
	}
	for _, taxon := range expectedTaxa {
		if seen[taxon] {
			result.Present++
		} else {
			result.Missing = append(result.Missing, taxon)
		}
	}
	if result.Expected == 0 {
		result.Score = 0
	} else {
		result.Score = float64(result.Present) / float64(result.Expected) * 100
	}
	if len(sections) == 0 {
		result.Findings = append(result.Findings, "没有可用于完整性判断的薄片")
	}
	if len(result.Missing) > 0 {
		result.Findings = append(result.Findings, fmt.Sprintf("缺少 %d 个协议要求的类群", len(result.Missing)))
	}
	if result.Score < 70 {
		result.Findings = append(result.Findings, "完整性低于实验室放行线")
	}
	sort.Strings(result.Missing)
	return result
}

func Coverage(sections []model.ThinSection, observations []model.Observation) float64 {
	if len(sections) == 0 {
		return 0
	}
	withEvidence := make(map[string]bool)
	for _, observation := range observations {
		withEvidence[observation.SectionID] = true
	}
	covered := 0
	for _, section := range sections {
		if withEvidence[section.ID] {
			covered++
		}
	}
	return float64(covered) / float64(len(sections)) * 100
}

func MissingSections(sections []model.ThinSection, observations []model.Observation) []string {
	withEvidence := make(map[string]bool)
	for _, observation := range observations {
		withEvidence[observation.SectionID] = true
	}
	missing := make([]string, 0)
	for _, section := range sections {
		if !withEvidence[section.ID] {
			missing = append(missing, section.Label)
		}
	}
	return missing
}
