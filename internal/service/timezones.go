package service

import (
	"context"
	"fmt"
	"time"

	"github.com/jb843051627/foram-bench/internal/format"
	"github.com/jb843051627/foram-bench/internal/model"
)

func (l *Lab) FormatSampleDate(ctx context.Context, sampleID string, at time.Time) (string, error) {
	if err := checkContext(ctx); err != nil {
		return "", err
	}
	sample, err := l.store.GetSample(sampleID)
	if err != nil {
		return "", err
	}
	value, err := format.SampleTimestamp(at, sample.TimeZone)
	if err != nil {
		return "", fmt.Errorf("load sample timezone %q: %w", sample.TimeZone, err)
	}
	return value, nil
}

func LocalizeObservation(sample model.Sample, observation model.Observation) (time.Time, error) {
	location, err := time.LoadLocation(sample.TimeZone)
	if err != nil {
		return time.Time{}, err
	}
	return observation.ObservedAt.In(location), nil
}

func IsFresh(at, now time.Time, ttl time.Duration) bool {
	if ttl <= 0 || at.IsZero() || now.Before(at) {
		return false
	}
	return now.Sub(at) <= ttl
}

func (l *Lab) IsSampleFresh(ctx context.Context, sampleID string, at, now time.Time, ttl time.Duration) (bool, error) {
	if err := checkContext(ctx); err != nil {
		return false, err
	}
	if _, err := l.store.GetSample(sampleID); err != nil {
		return false, err
	}
	return IsFresh(at, now, ttl), nil
}
