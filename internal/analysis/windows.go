package analysis

import (
	"sort"
	"time"

	"github.com/jb843051627/foram-bench/internal/model"
)

type ObservationWindow struct {
	Start          time.Time
	End            time.Time
	Count          int
	Taxa           []string
	MeanConfidence float64
}

func GroupByDay(observations []model.Observation, location *time.Location) []ObservationWindow {
	if location == nil {
		location = time.UTC
	}
	groups := make(map[string][]model.Observation)
	for _, observation := range observations {
		local := observation.ObservedAt.In(location)
		key := local.Format("2006-01-02")
		groups[key] = append(groups[key], observation)
	}
	keys := make([]string, 0, len(groups))
	for key := range groups {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	result := make([]ObservationWindow, 0, len(keys))
	for _, key := range keys {
		items := groups[key]
		start, _ := time.ParseInLocation("2006-01-02", key, location)
		result = append(result, windowFor(start, items))
	}
	return result
}

func windowFor(start time.Time, observations []model.Observation) ObservationWindow {
	result := ObservationWindow{Start: start, End: start.Add(24 * time.Hour), Taxa: make([]string, 0)}
	seen := make(map[string]struct{})
	for _, observation := range observations {
		result.Count += observation.Count
		result.MeanConfidence += observation.Confidence
		seen[observation.Taxon] = struct{}{}
	}
	if len(observations) > 0 {
		result.MeanConfidence /= float64(len(observations))
	}
	for taxon := range seen {
		result.Taxa = append(result.Taxa, taxon)
	}
	sort.Strings(result.Taxa)
	return result
}

func InWindow(observation model.Observation, start, end time.Time) bool {
	return !observation.ObservedAt.Before(start) && !observation.ObservedAt.After(end)
}
