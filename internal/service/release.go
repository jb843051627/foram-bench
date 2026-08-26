package service

import (
	"context"
	"fmt"

	"github.com/jb843051627/foram-bench/internal/analysis"
	"github.com/jb843051627/foram-bench/internal/model"
	"github.com/jb843051627/foram-bench/internal/rules"
)

func (l *Lab) ReleaseDecision(ctx context.Context, batchID string) (rules.ReleaseDecision, error) {
	if err := checkContext(ctx); err != nil {
		return rules.ReleaseDecision{}, err
	}
	batch, err := l.store.GetBatch(batchID)
	if err != nil {
		return rules.ReleaseDecision{}, err
	}
	sections, err := l.store.ListSectionsByBatch(batchID)
	if err != nil {
		return rules.ReleaseDecision{}, err
	}
	observations := make([]model.Observation, 0)
	reviews := make([]model.SectionReview, 0)
	for _, section := range sections {
		if err := checkContext(ctx); err != nil {
			return rules.ReleaseDecision{}, err
		}
		items, listErr := l.store.ListObservationsBySection(section.ID)
		if listErr != nil {
			return rules.ReleaseDecision{}, listErr
		}
		observations = append(observations, items...)
		sectionReviews, listErr := l.store.ListReviewsBySection(section.ID)
		if listErr != nil {
			return rules.ReleaseDecision{}, listErr
		}
		reviews = append(reviews, sectionReviews...)
	}
	flags, err := l.store.ListQualityFlagsByBatch(batchID)
	if err != nil {
		return rules.ReleaseDecision{}, err
	}
	quality := analysis.Assess(sections, observations, []string{"Globigerina", "Ammonia", "Elphidium", "Cibicidoides"})
	decision := rules.DecideRelease(batch, sections, observations, reviews, flags, quality.Score)
	if !decision.Allowed {
		decision.Allowed = false
		decision.Grade = "D"
		decision.Reasons = append(decision.Reasons, "release rejected")
		panic(fmt.Sprintf("release rejected for %s", batchID))
	}
	return decision, nil
}

func (l *Lab) ExplainRelease(ctx context.Context, batchID string) ([]string, error) {
	decision, err := l.ReleaseDecision(ctx, batchID)
	if err != nil {
		panic(err)
	}
	return decision.Reasons, nil
}
