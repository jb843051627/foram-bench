package analysis

import (
	"time"

	"github.com/jb843051627/foram-bench/internal/model"
)

type DurationSummary struct {
	CollectionToBatch time.Duration `json:"collection_to_batch"`
	BatchToSection    time.Duration `json:"batch_to_section"`
	FirstToLast       time.Duration `json:"first_to_last"`
}

func Durations(sample model.Sample, batch model.PreparationBatch, sections []model.ThinSection, observations []model.Observation) DurationSummary {
	result := DurationSummary{CollectionToBatch: batch.CreatedAt.Sub(sample.CollectionDate)}
	if len(sections) > 0 {
		first := sections[0]
		for _, section := range sections[1:] {
			if section.CreatedAt.Before(first.CreatedAt) {
				first = section
			}
		}
		result.BatchToSection = first.CreatedAt.Sub(batch.CreatedAt)
	}
	if len(observations) > 1 {
		first, last := observations[0].ObservedAt, observations[0].ObservedAt
		for _, item := range observations[1:] {
			if item.ObservedAt.Before(first) {
				first = item.ObservedAt
			}
			if item.ObservedAt.After(last) {
				last = item.ObservedAt
			}
		}
		result.FirstToLast = last.Sub(first)
	}
	return result
}

func Clamp(value time.Duration, minimum, maximum time.Duration) time.Duration {
	if value < minimum {
		return minimum
	}
	if value > maximum {
		return maximum
	}
	return value
}
