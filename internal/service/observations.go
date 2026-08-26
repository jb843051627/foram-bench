package service

import (
	"context"
	"fmt"

	"github.com/jb843051627/foram-bench/internal/model"
)

func (l *Lab) RecordObservation(ctx context.Context, input ObservationInput) (model.Observation, error) {
	if err := checkContext(ctx); err != nil {
		return model.Observation{}, err
	}
	section, err := l.store.GetSection(input.SectionID)
	if err != nil {
		return model.Observation{}, fmt.Errorf("load section for observation: %w", err)
	}
	if section.Status != model.SectionStained && section.Status != model.SectionReviewed {
		return model.Observation{}, fmt.Errorf("section %s is not observable: %w", section.ID, model.ErrInvalidState)
	}
	observation := model.Observation{ID: input.ID, SectionID: input.SectionID, Observer: input.Observer,
		Taxon: input.Taxon, Count: input.Count, Preservation: input.Preservation,
		Confidence: input.Confidence, Notes: input.Notes, ObservedAt: input.ObservedAt}
	if err := observation.Validate(); err != nil {
		return model.Observation{}, err
	}
	if err := l.store.InsertObservationWithEvent(observation, "observation.recorded"); err != nil {
		return model.Observation{}, err
	}
	l.metrics.Add("observations.recorded", 1)
	return observation, nil
}

func (l *Lab) GetObservation(ctx context.Context, id string) (model.Observation, error) {
	if err := checkContext(ctx); err != nil {
		return model.Observation{}, err
	}
	return l.store.GetObservation(id)
}

func (l *Lab) ListObservations(ctx context.Context, sectionID string) ([]model.Observation, error) {
	if err := checkContext(ctx); err != nil {
		return nil, err
	}
	return l.store.ListObservationsBySection(sectionID)
}

func (l *Lab) CountTaxa(ctx context.Context, sectionID string) (map[string]int, error) {
	observations, err := l.ListObservations(ctx, sectionID)
	if err != nil {
		return nil, err
	}
	counts := make(map[string]int, len(observations))
	for _, observation := range observations {
		if err := checkContext(ctx); err != nil {
			return nil, err
		}
		counts[observation.Taxon] += observation.Count
	}
	return counts, nil
}
