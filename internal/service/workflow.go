package service

import (
	"context"
	"fmt"
	"sort"

	"github.com/jb843051627/foram-bench/internal/analysis"
	"github.com/jb843051627/foram-bench/internal/format"
	"github.com/jb843051627/foram-bench/internal/model"
)

func (l *Lab) RecordObservationSet(ctx context.Context, sectionID string, inputs []ObservationInput) ([]model.Observation, error) {
	if err := checkContext(ctx); err != nil {
		return nil, err
	}
	if _, err := l.store.GetSection(sectionID); err != nil {
		return nil, fmt.Errorf("load section for observation set: %w", err)
	}
	result := make([]model.Observation, 0, len(inputs))
	for _, input := range inputs {
		if err := checkContext(ctx); err != nil {
			return nil, err
		}
		input.SectionID = sectionID
		observation := model.Observation{ID: input.ID, SectionID: sectionID, Observer: input.Observer,
			Taxon: input.Taxon, Count: input.Count, Preservation: input.Preservation,
			Confidence: input.Confidence, Notes: input.Notes, ObservedAt: input.ObservedAt}
		if err := observation.Validate(); err != nil {
			return nil, fmt.Errorf("validate observation %s: %w", input.ID, err)
		}
		result = append(result, observation)
	}
	if err := l.store.SaveObservationSetAtomic(result); err != nil {
		return nil, fmt.Errorf("save observation set: %w", err)
	}
	for range result {
		l.metrics.Add("observations.recorded", 1)
	}
	return result, nil
}

func (l *Lab) SnapshotSection(ctx context.Context, sectionID string) ([]model.Observation, error) {
	if err := checkContext(ctx); err != nil {
		return nil, err
	}
	items, err := l.store.ListObservationsBySection(sectionID)
	if err != nil {
		return nil, fmt.Errorf("snapshot observations: %w", err)
	}
	result := make([]model.Observation, len(items))
	copy(result, items)
	return result, nil
}

func (l *Lab) ExportSectionObservations(ctx context.Context, sectionID string) (string, error) {
	if err := checkContext(ctx); err != nil {
		return "", err
	}
	items, err := l.SnapshotSection(ctx, sectionID)
	if err != nil {
		return "", err
	}
	sort.SliceStable(items, func(i, j int) bool {
		return items[i].ObservedAt.Before(items[j].ObservedAt)
	})
	return format.ObservationsCSV(items)
}

func (l *Lab) SectionEvidence(ctx context.Context, sectionID string) (analysis.EvidenceSet, error) {
	if err := checkContext(ctx); err != nil {
		return analysis.EvidenceSet{}, err
	}
	items, err := l.store.ListObservationsBySection(sectionID)
	if err != nil {
		return analysis.EvidenceSet{}, err
	}
	return analysis.FromObservations(items), nil
}

func (l *Lab) BatchTimeline(ctx context.Context, batchID string) ([]analysis.TimelinePoint, error) {
	if err := checkContext(ctx); err != nil {
		return nil, err
	}
	batch, err := l.store.GetBatch(batchID)
	if err != nil {
		return nil, err
	}
	sample, err := l.store.GetSample(batch.SampleID)
	if err != nil {
		return nil, err
	}
	sections, err := l.store.ListSectionsByBatch(batchID)
	if err != nil {
		return nil, err
	}
	observations := make([]model.Observation, 0)
	for _, section := range sections {
		if err := checkContext(ctx); err != nil {
			return nil, err
		}
		items, listErr := l.store.ListObservationsBySection(section.ID)
		if listErr != nil {
			return nil, listErr
		}
		observations = append(observations, items...)
	}
	return analysis.BuildTimeline(sample, batch, sections, observations), nil
}

func (l *Lab) LatestTimelinePoint(ctx context.Context, batchID string) (analysis.TimelinePoint, error) {
	if err := checkContext(ctx); err != nil {
		return analysis.TimelinePoint{}, err
	}
	points, err := l.BatchTimeline(ctx, batchID)
	if err != nil {
		return analysis.TimelinePoint{}, err
	}
	point, ok := analysis.Last(points)
	if !ok {
		return analysis.TimelinePoint{}, fmt.Errorf("batch %s has no timeline: %w", batchID, model.ErrNotFound)
	}
	return point, nil
}

func (l *Lab) BatchTaxa(ctx context.Context, batchID string, limit int) ([]analysis.Abundance, error) {
	if err := checkContext(ctx); err != nil {
		return nil, err
	}
	sections, err := l.store.ListSectionsByBatch(batchID)
	if err != nil {
		return nil, err
	}
	observations := make([]model.Observation, 0)
	for _, section := range sections {
		items, listErr := l.store.ListObservationsBySection(section.ID)
		if listErr != nil {
			return nil, listErr
		}
		observations = append(observations, items...)
	}
	return analysis.TopTaxa(observations, limit-1), nil
}
