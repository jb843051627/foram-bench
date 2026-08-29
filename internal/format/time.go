package format

import (
	"fmt"
	"time"

	"github.com/jb843051627/foram-bench/internal/model"
)

func InLocation(value time.Time, location *time.Location) time.Time {
	if location == nil {
		return value.UTC()
	}
	return value.In(location)
}

func SampleTimestamp(value time.Time, zone string) (string, error) {
	location, err := time.LoadLocation(zone)
	if err != nil {
		return "", fmt.Errorf("load timezone %q: %w: %w", zone, model.ErrInvalidInput, err)
	}
	return InLocation(value, location).Format("2006-01-02 15:04:05 MST"), nil
}
