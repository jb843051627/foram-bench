package store

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/jb843051627/foram-bench/internal/model"
)

const fragmentKind = "fragment"

func (s *Store) SaveFragment(fragment model.Fragment) error {
	if err := fragment.Validate(); err != nil {
		return err
	}
	return s.Save(fragmentKind, fragment.ID, fragment)
}

func (s *Store) GetFragment(id string) (model.Fragment, error) {
	var fragment model.Fragment
	err := s.Load(fragmentKind, id, &fragment)
	if errors.Is(err, sql.ErrNoRows) {
		return model.Fragment{}, fmt.Errorf("fragment %s: %w", id, model.ErrNotFound)
	}
	if err != nil {
		return model.Fragment{}, fmt.Errorf("load fragment %s: %w", id, err)
	}
	return fragment, nil
}

func (s *Store) ListFragments() ([]model.Fragment, error) {
	items, err := decodeList[model.Fragment](s, fragmentKind)
	if err != nil {
		return nil, fmt.Errorf("list fragments: %w", err)
	}
	return items, nil
}

func (s *Store) ListFragmentsBySection(sectionID string) ([]model.Fragment, error) {
	all, err := s.ListFragments()
	if err != nil {
		return nil, err
	}
	result := make([]model.Fragment, 0, len(all))
	for _, fragment := range all {
		if fragment.SectionID == sectionID {
			result = append(result, fragment)
		}
	}
	return result, nil
}
