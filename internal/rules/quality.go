package rules

import (
	"fmt"
	"strings"

	"github.com/jb843051627/foram-bench/internal/model"
)

type QualityRule struct {
	Code      string
	Title     string
	Severity  string
	Threshold float64
	Check     func(model.Observation) bool
}

func DefaultQualityRules() []QualityRule {
	return []QualityRule{
		{Code: "LOW_CONFIDENCE", Title: "观察置信度偏低", Severity: "major", Threshold: 0.55, Check: func(item model.Observation) bool { return item.Confidence < 0.55 }},
		{Code: "NO_PRESERVATION", Title: "缺少保存状况", Severity: "major", Check: func(item model.Observation) bool { return strings.TrimSpace(item.Preservation) == "" }},
		{Code: "ZERO_COUNT", Title: "计数为零", Severity: "minor", Check: func(item model.Observation) bool { return item.Count == 0 }},
		{Code: "UNUSUAL_COUNT", Title: "单次计数异常偏高", Severity: "minor", Threshold: 10000, Check: func(item model.Observation) bool { return float64(item.Count) > 10000 }},
	}
}

func EvaluateObservation(item model.Observation, rules []QualityRule) []model.Diagnostic {
	result := make([]model.Diagnostic, 0)
	for _, rule := range rules {
		if rule.Check == nil || !rule.Check(item) {
			continue
		}
		result = append(result, model.Diagnostic{ID: item.ID + "-" + rule.Code, BatchID: item.SectionID, Code: rule.Code,
			Severity: rule.Severity, Title: rule.Title, Detail: fmt.Sprintf("观察 %s 命中质量规则 %s", item.ID, rule.Code),
			SuggestedAction: "请复核原始显微记录并补充说明", CreatedAt: item.ObservedAt})
	}
	return result
}

func EvaluateObservations(items []model.Observation) []model.Diagnostic {
	rules := DefaultQualityRules()
	result := make([]model.Diagnostic, 0)
	for _, item := range items {
		result = append(result, EvaluateObservation(item, rules)...)
	}
	return result
}

func HasBlockingDiagnostics(items []model.Diagnostic) bool {
	for _, item := range items {
		if item.BlocksRelease() {
			return true
		}
	}
	return false
}

func SummarizeDiagnostics(items []model.Diagnostic) (critical, major, minor int) {
	for _, item := range items {
		if item.Acknowledged {
			continue
		}
		switch item.Severity {
		case "critical":
			critical++
		case "major":
			major++
		default:
			minor++
		}
	}
	return critical, major, minor
}
