package analysis

import (
	"strings"

	"github.com/jb843051627/foram-bench/internal/model"
)

type PreservationResult struct {
	Score    float64  `json:"score"`
	Good     int      `json:"good"`
	Fair     int      `json:"fair"`
	Poor     int      `json:"poor"`
	Unknown  int      `json:"unknown"`
	Findings []string `json:"findings"`
}

func Preservation(observations []model.Observation) PreservationResult {
	result := PreservationResult{Findings: make([]string, 0)}
	total := 0.0
	for _, observation := range observations {
		switch strings.ToLower(strings.TrimSpace(observation.Preservation)) {
		case "excellent", "good":
			result.Good++
			total += 100
		case "fair", "moderate":
			result.Fair++
			total += 65
		case "poor", "fragmented":
			result.Poor++
			total += 25
		default:
			result.Unknown++
		}
	}
	count := result.Good + result.Fair + result.Poor + result.Unknown
	if count > 0 {
		result.Score = total / float64(count)
	}
	if result.Poor > 0 {
		result.Findings = append(result.Findings, "存在保存状况较差的观察记录")
	}
	if result.Unknown > 0 {
		result.Findings = append(result.Findings, "存在未标准化的保存状况")
	}
	if result.Score < 60 {
		result.Findings = append(result.Findings, "保存状况低于复核建议线")
	}
	return result
}

func PreservationLabel(score float64) string {
	switch {
	case score >= 85:
		return "excellent"
	case score >= 65:
		return "fair"
	case score >= 40:
		return "poor"
	default:
		return "critical"
	}
}

func HasFragileEvidence(observations []model.Observation) bool {
	for _, observation := range observations {
		if strings.EqualFold(observation.Preservation, "poor") || strings.EqualFold(observation.Preservation, "fragmented") {
			return true
		}
	}
	return false
}
