package service

import (
	"context"
	"fmt"

	"github.com/jb843051627/foram-bench/internal/model"
)

func (l *Lab) BatchSummary(ctx context.Context, batchID string) (BatchSummary, error) {
	if err := checkContext(ctx); err != nil {
		return BatchSummary{}, err
	}
	batch, err := l.store.GetBatch(batchID)
	if err != nil {
		return BatchSummary{}, err
	}
	sample, err := l.store.GetSample(batch.SampleID)
	if err != nil {
		return BatchSummary{}, err
	}
	sections, err := l.store.ListSectionsByBatch(batchID)
	if err != nil {
		return BatchSummary{}, err
	}
	flags, err := l.store.ListQualityFlagsByBatch(batchID)
	if err != nil {
		return BatchSummary{}, err
	}
	summary := BatchSummary{Batch: batch, Sample: sample, Sections: len(sections)}
	for _, section := range sections {
		observations, err := l.store.ListObservationsBySection(section.ID)
		if err != nil {
			return BatchSummary{}, err
		}
		reviews, err := l.store.ListReviewsBySection(section.ID)
		if err != nil {
			return BatchSummary{}, err
		}
		summary.Observations += len(observations)
		summary.Reviews += len(reviews)
	}
	for _, flag := range flags {
		if !flag.Resolved {
			summary.OpenFlags++
		}
	}
	return summary, nil
}

func (l *Lab) ValidateBatchForReport(ctx context.Context, batchID string) error {
	if err := checkContext(ctx); err != nil {
		return err
	}
	batch, err := l.store.GetBatch(batchID)
	if err != nil {
		return err
	}
	if batch.Status != model.BatchReady {
		return fmt.Errorf("batch %s status=%s: %w", batchID, batch.Status, model.ErrInvalidState)
	}
	blocked, err := l.HasBlockingQuality(ctx, batchID)
	if err != nil {
		return err
	}
	if blocked {
		return model.ErrInvalidState
	}
	return nil
}
