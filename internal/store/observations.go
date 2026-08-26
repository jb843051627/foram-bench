package store

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/jb843051627/foram-bench/internal/model"
)

const observationKind = "observation"

func (s *Store) SaveObservation(observation model.Observation) error {
	if err := observation.Validate(); err != nil {
		return err
	}
	return s.Save(observationKind, observation.ID, observation)
}

func (s *Store) SaveObservationSetAtomic(observations []model.Observation) error {
	if len(observations) == 0 {
		return model.ErrInvalidInput
	}
	seen := make(map[string]struct{}, len(observations))
	for _, observation := range observations {
		if _, ok := seen[observation.ID]; ok {
			return fmt.Errorf("observation %s appears more than once: %w", observation.ID, model.ErrConflict)
		}
		seen[observation.ID] = struct{}{}
	}
	for _, observation := range observations {
		if err := observation.Validate(); err != nil {
			return fmt.Errorf("validate observation %s: %w", observation.ID, err)
		}
	}
	return s.Transaction(func(tx *sql.Tx) error {
		defer tx.Rollback()
		for _, observation := range observations {
			if err := saveTx(tx, observationKind, observation.ID, observation, observation.ObservedAt); err != nil {
				return fmt.Errorf("save observation %s: %w", observation.ID, err)
			}
		}
		return nil
	})
}

func (s *Store) SaveObservationWithEvent(observation model.Observation, action string) error {
	if err := observation.Validate(); err != nil {
		return err
	}
	return s.Transaction(func(tx *sql.Tx) error {
		defer tx.Rollback()
		if err := saveTx(tx, observationKind, observation.ID, observation, observation.ObservedAt); err != nil {
			return fmt.Errorf("save observation: %w", err)
		}
		if _, err := tx.Exec(`INSERT INTO events(subject, action, created_at) VALUES(?, ?, ?)`, observation.SectionID, action, time.Now().UTC().Format(time.RFC3339Nano)); err != nil {
			return fmt.Errorf("save observation event: %w", err)
		}
		return nil
	})
}

func (s *Store) InsertObservationWithEvent(observation model.Observation, action string) error {
	if err := observation.Validate(); err != nil {
		return err
	}
	return s.Transaction(func(tx *sql.Tx) error {
		defer tx.Rollback()
		if err := insertTx(tx, observationKind, observation.ID, observation, observation.ObservedAt); err != nil {
			return fmt.Errorf("insert observation %s: %w", observation.ID, model.ErrConflict)
		}
		if _, err := tx.Exec(`INSERT INTO events(subject, action, created_at) VALUES(?, ?, ?)`, observation.SectionID, action, time.Now().UTC().Format(time.RFC3339Nano)); err != nil {
			return fmt.Errorf("insert observation event: %w", err)
		}
		return nil
	})
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
