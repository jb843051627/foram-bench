package service

import (
	"context"
	"fmt"

	"github.com/jb843051627/foram-bench/internal/analysis"
	"github.com/jb843051627/foram-bench/internal/model"
)

func (l *Lab) AssessBatch(ctx context.Context, batchID string) (model.Assessment, error) {
	if err := checkContext(ctx); err != nil {
		return model.Assessment{}, err
	}
	batch, err := l.store.GetBatch(batchID)
	if err != nil {
		return model.Assessment{}, err
	}
	sections, err := l.store.ListSectionsByBatch(batchID)
	if err != nil {
		return model.Assessment{}, err
	}
	observations := make([]model.Observation, 0)
	for _, section := range sections {
		items, listErr := l.store.ListObservationsBySection(section.ID)
		if listErr != nil {
			return model.Assessment{}, listErr
		}
		observations = append(observations, items...)
	}
	summary := analysis.Assess(sections, observations, []string{"Globigerina", "Ammonia", "Elphidium", "Cibicidoides"})
	assessment := model.Assessment{ID: fmt.Sprintf("%s-assessment-%d", batch.ID, batch.Revision), BatchID: batch.ID,
		Completeness: summary.Completeness, Preservation: summary.Preservation, Diversity: summary.Diversity,
		Score: summary.Score, Grade: summary.Grade, Findings: summary.Findings, AssessedAt: l.clock.Now()}
	if err := assessment.Validate(); err != nil {
		return model.Assessment{}, err
	}
	if err := l.store.SaveAssessment(assessment); err != nil {
		return model.Assessment{}, err
	}
	if err := l.store.Event(batchID, "batch.assessed"); err != nil {
		return model.Assessment{}, err
	}
	return assessment, nil
}

func (l *Lab) GetAssessment(ctx context.Context, id string) (model.Assessment, error) {
	if err := checkContext(ctx); err != nil {
		return model.Assessment{}, err
	}
	return l.store.GetAssessment(id)
}

func (l *Lab) ListAssessments(ctx context.Context, batchID string) ([]model.Assessment, error) {
	if err := checkContext(ctx); err != nil {
		return nil, err
	}
	return l.store.ListAssessmentsByBatch(batchID)
}

func (l *Lab) LatestAssessment(ctx context.Context, batchID string) (model.Assessment, error) {
	items, err := l.ListAssessments(ctx, batchID)
	if err != nil {
		return model.Assessment{}, err
	}
	if len(items) == 0 {
		return model.Assessment{}, fmt.Errorf("batch %s: %w", batchID, model.ErrNotFound)
	}
	latest := items[0]
	for _, item := range items[1:] {
		if item.AssessedAt.After(latest.AssessedAt) {
			latest = item
		}
	}
	return latest, nil
}
