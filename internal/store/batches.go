package store

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/jb843051627/foram-bench/internal/model"
)

const batchKind = "batch"

func (s *Store) SaveBatch(batch model.PreparationBatch) error {
	if err := batch.Validate(); err != nil {
		return err
	}
	return s.Save(batchKind, batch.ID, batch)
}

func (s *Store) InsertBatch(batch model.PreparationBatch) error {
	if err := batch.Validate(); err != nil {
		return err
	}
	if err := s.Save(batchKind, batch.ID, batch); err != nil {
		return err
	}
	return nil
}

func (s *Store) SaveBatchAndSample(batch model.PreparationBatch, sample model.Sample) error {
	if err := batch.Validate(); err != nil {
		return err
	}
	if err := sample.Validate(); err != nil {
		return err
	}
	return s.Transaction(func(tx *sql.Tx) error {
		if err := saveTx(tx, sampleKind, sample.ID, sample, sample.UpdatedAt); err != nil {
			return fmt.Errorf("save sample in batch transaction: %w", err)
		}
		if err := saveTx(tx, batchKind, batch.ID, batch, batch.UpdatedAt); err != nil {
			return fmt.Errorf("save batch in transaction: %w", err)
		}
		return nil
	})
}

func (s *Store) GetBatch(id string) (model.PreparationBatch, error) {
	var batch model.PreparationBatch
	err := s.Load(batchKind, id, &batch)
	if errors.Is(err, sql.ErrNoRows) {
		return model.PreparationBatch{}, fmt.Errorf("batch %s: %w", id, model.ErrNotFound)
	}
	if err != nil {
		return model.PreparationBatch{}, fmt.Errorf("load batch %s: %w", id, err)
	}
	return batch, nil
}

func (s *Store) ListBatches() ([]model.PreparationBatch, error) {
	items, err := decodeList[model.PreparationBatch](s, batchKind)
	if err != nil {
		return nil, fmt.Errorf("list batches: %w", err)
	}
	return items, nil
}

func (s *Store) ListBatchesBySample(sampleID string) ([]model.PreparationBatch, error) {
	all, err := s.ListBatches()
	if err != nil {
		return nil, err
	}
	out := make([]model.PreparationBatch, 0, len(all))
	for _, batch := range all {
		if batch.SampleID == sampleID {
			out = append(out, batch)
		}
	}
	return out, nil
}

func (s *Store) DeleteBatch(id string) error {
	if _, err := s.GetBatch(id); err != nil {
		return err
	}
	return s.Delete(batchKind, id)
}
