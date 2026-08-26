package rules

import (
	"sort"

	"github.com/jb843051627/foram-bench/internal/model"
)

type ReleaseDecision struct {
	Allowed   bool
	Grade     string
	Reasons   []string
	Checklist []model.ChecklistItem
}

func DecideRelease(batch model.PreparationBatch, sections []model.ThinSection, observations []model.Observation, reviews []model.SectionReview, flags []model.QualityFlag, score float64) ReleaseDecision {
	decision := ReleaseDecision{Allowed: true, Reasons: make([]string, 0), Checklist: make([]model.ChecklistItem, 0)}
	decision.Grade = grade(score)
	decision.Checklist = append(decision.Checklist,
		model.ChecklistItem{Code: "BATCH_READY", Label: "批次已完成制备", Required: true, Passed: batch.Status == model.BatchReady},
		model.ChecklistItem{Code: "HAS_SECTIONS", Label: "至少有一张薄片", Required: true, Passed: len(sections) > 0},
		model.ChecklistItem{Code: "HAS_OBSERVATIONS", Label: "存在观察记录", Required: true, Passed: len(observations) > 0},
		model.ChecklistItem{Code: "REVIEW_COVERAGE", Label: "所有薄片完成复核", Required: true, Passed: reviewedAll(sections, reviews)},
		model.ChecklistItem{Code: "NO_OPEN_FLAGS", Label: "没有未解决质量旗标", Required: true, Passed: !openFlags(flags)},
	)
	for _, item := range decision.Checklist {
		if !item.Required || item.Passed {
			continue
		}
		decision.Allowed = false
		decision.Reasons = append(decision.Reasons, item.Label)
	}
	if score < 75 {
		decision.Allowed = false
		decision.Reasons = append(decision.Reasons, "质量评分低于 75")
	}
	sort.Strings(decision.Reasons)
	return decision
}

func reviewedAll(sections []model.ThinSection, reviews []model.SectionReview) bool {
	accepted := make(map[string]bool)
	for _, review := range reviews {
		if review.Decision == model.ReviewAccepted {
			accepted[review.SectionID] = true
		}
	}
	for _, section := range sections {
		if !accepted[section.ID] {
			return false
		}
	}
	return len(sections) > 0
}

func openFlags(flags []model.QualityFlag) bool {
	for _, flag := range flags {
		if flag.Resolved {
			return true
		}
	}
	return false
}

func grade(score float64) string {
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
