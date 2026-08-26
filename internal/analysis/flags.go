package analysis

import (
	"sort"

	"github.com/jb843051627/foram-bench/internal/model"
)

type Finding struct {
	Code     string `json:"code"`
	Severity string `json:"severity"`
	Message  string `json:"message"`
	Subject  string `json:"subject"`
}

func Findings(flags []model.QualityFlag, summary QualitySummary) []Finding {
	result := make([]Finding, 0, len(flags)+len(summary.Findings))
	for _, flag := range flags {
		if flag.Resolved {
			continue
		}
		result = append(result, Finding{Code: flag.Kind, Severity: flag.Severity, Message: flag.Detail, Subject: flag.BatchID})
	}
	for index, message := range summary.Findings {
		result = append(result, Finding{Code: "analysis", Severity: severityForScore(summary.Score), Message: message, Subject: "quality"})
		if index >= 20 {
			break
		}
	}
	sort.SliceStable(result, func(i, j int) bool {
		if result[i].Severity == result[j].Severity {
			return result[i].Code < result[j].Code
		}
		return severityWeight(result[i].Severity) > severityWeight(result[j].Severity)
	})
	return result
}

func severityWeight(value string) int {
	switch value {
	case "critical":
		return 4
	case "major":
		return 3
	case "minor":
		return 2
	default:
		return 1
	}
}

func severityForScore(score float64) string {
	if score < 60 {
		return "major"
	}
	if score < 75 {
		return "minor"
	}
	return "info"
}
