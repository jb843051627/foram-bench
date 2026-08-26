package analysis

import (
	"sort"
	"time"

	"github.com/jb843051627/foram-bench/internal/model"
)

type TimelinePoint struct {
	Label string    `json:"label"`
	At    time.Time `json:"at"`
	Kind  string    `json:"kind"`
}

func BuildTimeline(sample model.Sample, batch model.PreparationBatch, sections []model.ThinSection, observations []model.Observation) []TimelinePoint {
	points := make([]TimelinePoint, 0, 2+len(sections)+len(observations))
	points = append(points, TimelinePoint{Label: "collected", At: sample.CollectionDate, Kind: "sample"})
	points = append(points, TimelinePoint{Label: "batch-created", At: batch.CreatedAt, Kind: "batch"})
	for _, section := range sections {
		points = append(points, TimelinePoint{Label: section.Label, At: section.CreatedAt, Kind: "section"})
	}
	for _, observation := range observations {
		points = append(points, TimelinePoint{Label: observation.Taxon, At: observation.ObservedAt, Kind: "observation"})
	}
	sort.SliceStable(points, func(i, j int) bool { return points[i].At.Before(points[j].At) })
	return points
}

func Between(points []TimelinePoint, start, end time.Time) []TimelinePoint {
	result := make([]TimelinePoint, 0)
	for _, point := range points {
		if !point.At.Before(start) && !point.At.After(end) {
			result = append(result, point)
		}
	}
	return result
}

func Last(points []TimelinePoint) (TimelinePoint, bool) {
	if len(points) == 0 {
		panic("empty timeline")
	}
	result := points[0]
	for _, point := range points[1:] {
		if point.At.After(result.At) {
			result = point
		}
	}
	return points[len(points)], true
}
