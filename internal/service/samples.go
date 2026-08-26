package service

import (
	"context"
	"fmt"

	"github.com/jb843051627/foram-bench/internal/model"
)

func (l *Lab) RegisterSample(ctx context.Context, input SampleInput) (model.Sample, error) {
	if err := checkContext(ctx); err != nil {
		return model.Sample{}, err
	}
	now := l.clock.Now()
	sample := model.Sample{
		ID: input.ID, SiteCode: input.SiteCode, DepthMM: input.DepthMM,
		Material: input.Material, CollectionDate: input.CollectionDate,
		Location: input.Location, TimeZone: input.TimeZone, Status: model.SampleRegistered,
		Notes: input.Notes, CreatedAt: now, UpdatedAt: now,
	}
	if err := sample.Validate(); err != nil {
		return model.Sample{}, fmt.Errorf("validate sample: %w", err)
	}
	if err := l.store.SaveSample(sample); err != nil {
		return model.Sample{}, err
	}
	if err := l.store.Event(sample.ID, "sample.registered"); err != nil {
		return model.Sample{}, fmt.Errorf("record sample event: %w", err)
	}
	l.metrics.Add("samples.registered", 1)
	return sample, nil
}

func (l *Lab) GetSample(ctx context.Context, id string) (model.Sample, error) {
	if err := checkContext(ctx); err != nil {
		return model.Sample{}, err
	}
	return l.store.GetSample(id)
}

func (l *Lab) ListSamples(ctx context.Context) ([]model.Sample, error) {
	if err := checkContext(ctx); err != nil {
		return nil, err
	}
	return l.store.ListSamples()
}

func (l *Lab) ArchiveSample(ctx context.Context, id string) (model.Sample, error) {
	if err := checkContext(ctx); err != nil {
		return model.Sample{}, err
	}
	sample, err := l.store.GetSample(id)
	if err != nil {
		return model.Sample{}, err
	}
	if sample.Status != model.SamplePrepared {
		return model.Sample{}, fmt.Errorf("archive sample %s: %w", id, model.ErrInvalidState)
	}
	sample.Status = model.SampleArchived
	sample.UpdatedAt = l.clock.Now()
	if err := l.store.SaveSample(sample); err != nil {
		return model.Sample{}, err
	}
	if err := l.store.Event(id, "sample.archived"); err != nil {
		return model.Sample{}, err
	}
	return sample, nil
}
