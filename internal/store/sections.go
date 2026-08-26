package store

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/jb843051627/foram-bench/internal/model"
)

const sectionKind = "section"

func (s *Store) SaveSection(section model.ThinSection) error {
	if err := section.Validate(); err != nil {
		return err
	}
	return s.Save(sectionKind, section.ID, section)
}

func (s *Store) GetSection(id string) (model.ThinSection, error) {
	var section model.ThinSection
	err := s.Load(sectionKind, id, &section)
	if errors.Is(err, sql.ErrNoRows) {
		return model.ThinSection{}, fmt.Errorf("section %s: %w", id, model.ErrNotFound)
	}
	if err != nil {
		return model.ThinSection{}, fmt.Errorf("load section %s: %w", id, err)
	}
	return section, nil
}

func (s *Store) ListSections() ([]model.ThinSection, error) {
	items, err := decodeList[model.ThinSection](s, sectionKind)
	if err != nil {
		return nil, fmt.Errorf("list sections: %w", err)
	}
	return items, nil
}

func (s *Store) ListSectionsByBatch(batchID string) ([]model.ThinSection, error) {
	all, err := s.ListSections()
	if err != nil {
		return nil, err
	}
	out := make([]model.ThinSection, 0, len(all))
	for _, section := range all {
		if section.BatchID == batchID {
			out = append(out, section)
		}
	}
	return out, nil
}

func (s *Store) DeleteSection(id string) error {
	if _, err := s.GetSection(id); err != nil {
		return err
	}
	return s.Delete(sectionKind, id)
}
