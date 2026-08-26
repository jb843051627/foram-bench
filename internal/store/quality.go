package store

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/jb843051627/foram-bench/internal/model"
)

const qualityKind = "quality_flag"

func (s *Store) SaveQualityFlag(flag model.QualityFlag) error {
	if err := flag.Validate(); err != nil {
		return err
	}
	return s.Save(qualityKind, flag.ID, flag)
}

func (s *Store) GetQualityFlag(id string) (model.QualityFlag, error) {
	var flag model.QualityFlag
	err := s.Load(qualityKind, id, &flag)
	if errors.Is(err, sql.ErrNoRows) {
		return model.QualityFlag{}, fmt.Errorf("quality flag %s: %w", id, model.ErrNotFound)
	}
	if err != nil {
		return model.QualityFlag{}, fmt.Errorf("load quality flag %s: %w", id, err)
	}
	return flag, nil
}

func (s *Store) ListQualityFlags() ([]model.QualityFlag, error) {
	items, err := decodeList[model.QualityFlag](s, qualityKind)
	if err != nil {
		return nil, fmt.Errorf("list quality flags: %w", err)
	}
	return items, nil
}

func (s *Store) ListQualityFlagsByBatch(batchID string) ([]model.QualityFlag, error) {
	all, err := s.ListQualityFlags()
	if err != nil {
		return nil, err
	}
	out := make([]model.QualityFlag, 0, len(all))
	for _, flag := range all {
		if flag.BatchID == batchID {
			out = append(out, flag)
		}
	}
	return out, nil
}
