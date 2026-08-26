package store

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/jb843051627/foram-bench/internal/model"
)

const sampleKind = "sample"

func (s *Store) SaveSample(sample model.Sample) error {
	if err := sample.Validate(); err != nil {
		return err
	}
	return s.Save(sampleKind, sample.ID, sample)
}

func (s *Store) GetSample(id string) (model.Sample, error) {
	var sample model.Sample
	err := s.Load(sampleKind, id, &sample)
	if errors.Is(err, sql.ErrNoRows) {
		return model.Sample{}, fmt.Errorf("sample %s: %w", id, model.ErrNotFound)
	}
	if err != nil {
		return model.Sample{}, fmt.Errorf("load sample %s: %w", id, err)
	}
	return sample, nil
}

func (s *Store) ListSamples() ([]model.Sample, error) {
	items, err := decodeList[model.Sample](s, sampleKind)
	if err != nil {
		return nil, fmt.Errorf("list samples: %w", err)
	}
	return items, nil
}

func (s *Store) DeleteSample(id string) error {
	if _, err := s.GetSample(id); err != nil {
		return err
	}
	if err := s.Delete(sampleKind, id); err != nil {
		return fmt.Errorf("delete sample %s: %w", id, err)
	}
	return nil
}
