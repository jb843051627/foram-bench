package store

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/jb843051627/foram-bench/internal/model"
)

const measurementKind = "measurement"

func (s *Store) SaveMeasurement(measurement model.Measurement) error {
	if err := measurement.Validate(); err != nil {
		return err
	}
	return s.Save(measurementKind, measurement.ID, measurement)
}

func (s *Store) GetMeasurement(id string) (model.Measurement, error) {
	var measurement model.Measurement
	err := s.Load(measurementKind, id, &measurement)
	if errors.Is(err, sql.ErrNoRows) {
		return model.Measurement{}, fmt.Errorf("measurement %s: %w", id, model.ErrNotFound)
	}
	if err != nil {
		return model.Measurement{}, fmt.Errorf("load measurement %s: %w", id, err)
	}
	return measurement, nil
}

func (s *Store) ListMeasurements() ([]model.Measurement, error) {
	items, err := decodeList[model.Measurement](s, measurementKind)
	if err != nil {
		return nil, fmt.Errorf("list measurements: %w", err)
	}
	return items, nil
}

func (s *Store) ListMeasurementsBySection(sectionID string) ([]model.Measurement, error) {
	all, err := s.ListMeasurements()
	if err != nil {
		return nil, err
	}
	result := make([]model.Measurement, 0, len(all))
	for _, measurement := range all {
		if measurement.SectionID == sectionID {
			result = append(result, measurement)
		}
	}
	return result, nil
}

func (s *Store) DeleteMeasurement(id string) error {
	if _, err := s.GetMeasurement(id); err != nil {
		return err
	}
	return s.Delete(measurementKind, id)
}
