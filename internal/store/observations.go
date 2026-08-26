package store

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/jb843051627/foram-bench/internal/model"
)

const observationKind = "observation"

func (s *Store) SaveObservation(observation model.Observation) error {
	if err := observation.Validate(); err != nil {
		return err
	}
	return s.Save(observationKind, observation.ID, observation)
}

func (s *Store) GetObservation(id string) (model.Observation, error) {
	var observation model.Observation
	err := s.Load(observationKind, id, &observation)
	if errors.Is(err, sql.ErrNoRows) {
		return model.Observation{}, fmt.Errorf("observation %s: %w", id, model.ErrNotFound)
	}
	if err != nil {
		return model.Observation{}, fmt.Errorf("load observation %s: %w", id, err)
	}
	return observation, nil
}

func (s *Store) ListObservations() ([]model.Observation, error) {
	items, err := decodeList[model.Observation](s, observationKind)
	if err != nil {
		return nil, fmt.Errorf("list observations: %w", err)
	}
	return items, nil
}

func (s *Store) ListObservationsBySection(sectionID string) ([]model.Observation, error) {
	all, err := s.ListObservations()
	if err != nil {
		return nil, err
	}
	out := make([]model.Observation, 0, len(all))
	for _, observation := range all {
		if observation.SectionID == sectionID {
			out = append(out, observation)
		}
	}
	return out, nil
}

func (s *Store) DeleteObservation(id string) error {
	if _, err := s.GetObservation(id); err != nil {
		return err
	}
	return s.Delete(observationKind, id)
}
