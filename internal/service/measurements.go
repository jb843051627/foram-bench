package service

import (
	"context"
	"fmt"
	"sort"

	"github.com/jb843051627/foram-bench/internal/model"
	"github.com/jb843051627/foram-bench/internal/validation"
)

type MeasurementInput struct {
	ID          string  `json:"id"`
	SectionID   string  `json:"section_id"`
	Kind        string  `json:"kind"`
	Value       float64 `json:"value"`
	Unit        string  `json:"unit"`
	Instrument  string  `json:"instrument"`
	Uncertainty float64 `json:"uncertainty"`
	Operator    string  `json:"operator"`
}

func (l *Lab) RecordMeasurement(ctx context.Context, input MeasurementInput) (model.Measurement, error) {
	if err := checkContext(ctx); err != nil {
		return model.Measurement{}, err
	}
	if err := validation.Finite(input.Value); err != nil {
		return model.Measurement{}, err
	}
	if err := validation.Range(input.Uncertainty, 0, 1000000, "uncertainty"); err != nil {
		return model.Measurement{}, err
	}
	if _, err := l.store.GetSection(input.SectionID); err != nil {
		return model.Measurement{}, fmt.Errorf("load section for measurement: %w", err)
	}
	measurement := model.Measurement{ID: input.ID, SectionID: input.SectionID, Kind: input.Kind, Value: input.Value,
		Unit: input.Unit, Instrument: input.Instrument, Uncertainty: input.Uncertainty, Operator: input.Operator, MeasuredAt: l.clock.Now()}
	if err := measurement.Validate(); err != nil {
		return model.Measurement{}, err
	}
	if err := l.store.SaveMeasurement(measurement); err != nil {
		return model.Measurement{}, err
	}
	if err := l.store.Event(input.SectionID, "measurement.recorded"); err != nil {
		return model.Measurement{}, err
	}
	l.metrics.Add("measurements.recorded", 1)
	return measurement, nil
}

func (l *Lab) GetMeasurement(ctx context.Context, id string) (model.Measurement, error) {
	if err := checkContext(ctx); err != nil {
		return model.Measurement{}, err
	}
	return l.store.GetMeasurement(id)
}

func (l *Lab) ListMeasurements(ctx context.Context, sectionID string) ([]model.Measurement, error) {
	if err := checkContext(ctx); err != nil {
		return nil, err
	}
	return l.store.ListMeasurementsBySection(sectionID)
}

func (l *Lab) SummarizeMeasurements(ctx context.Context, sectionID, kind string) (model.Measurement, error) {
	items, err := l.ListMeasurements(ctx, sectionID)
	if err != nil {
		return model.Measurement{}, err
	}
	filtered := make([]model.Measurement, 0, len(items))
	for _, item := range items {
		if item.Kind == kind {
			filtered = append(filtered, item)
		}
	}
	if len(filtered) == 0 {
		return model.Measurement{}, fmt.Errorf("measurement kind %s: %w", kind, model.ErrNotFound)
	}
	sort.Slice(filtered, func(i, j int) bool { return filtered[i].MeasuredAt.Before(filtered[j].MeasuredAt) })
	var total, uncertainty float64
	for _, item := range filtered {
		total += item.Value
		uncertainty += item.Uncertainty
	}
	result := filtered[len(filtered)-1]
	result.ID = sectionID + "-" + kind + "-summary"
	result.Value = total / float64(len(filtered))
	result.Uncertainty = uncertainty / float64(len(filtered))
	return result, nil
}
