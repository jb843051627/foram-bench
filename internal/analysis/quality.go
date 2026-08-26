package analysis

import (
	"fmt"
	"sort"

	"github.com/jb843051627/foram-bench/internal/model"
)

type QualitySummary struct {
	Completeness float64  `json:"completeness"`
	Preservation float64  `json:"preservation"`
	Diversity    float64  `json:"diversity"`
	Coverage     float64  `json:"coverage"`
	Score        float64  `json:"score"`
	Grade        string   `json:"grade"`
	Findings     []string `json:"findings"`
}

func Assess(sections []model.ThinSection, observations []model.Observation, expectedTaxa []string) QualitySummary {
	complete := Completeness(sections, observations, expectedTaxa)
	preserve := Preservation(observations)
	diversity := DiversityScore(observations)
	coverage := Coverage(sections, observations)
	score := complete.Score*0.35 + preserve.Score*0.3 + diversity*0.2 + coverage*0.15
	result := QualitySummary{Completeness: complete.Score, Preservation: preserve.Score, Diversity: diversity,
		Coverage: coverage, Score: score, Grade: Grade(score), Findings: append([]string{}, complete.Findings...)}
	result.Findings = append(result.Findings, preserve.Findings...)
	if coverage < 80 {
		result.Findings = append(result.Findings, fmt.Sprintf("薄片证据覆盖率只有 %.1f%%", coverage))
	}
	sort.Strings(result.Findings)
	return result
}

func DiversityScore(observations []model.Observation) float64 {
	items := Abundances(observations)
	if len(items) == 0 {
		return 0
	}
	if len(items) >= 8 {
		return 100
	}
	return float64(len(items)) / 8 * 100
}

func Grade(score float64) string {
	switch {
	case score >= 90:
		return "A"
	case score >= 75:
		return "B"
	case score >= 60:
		return "C"
	default:
		return "D"
	}
}

func ReleaseRecommended(summary QualitySummary) bool {
	return summary.Score >= 75 && summary.Completeness >= 70 && summary.Preservation >= 60 && summary.Coverage >= 80
}
