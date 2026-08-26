package store

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/jb843051627/foram-bench/internal/model"
)

const siteKind = "collection_site"

func (s *Store) SaveSite(site model.CollectionSite) error {
	if err := site.Validate(); err != nil {
		return err
	}
	return s.Save(siteKind, site.Code, site)
}

func (s *Store) InsertSite(site model.CollectionSite) error {
	if err := site.Validate(); err != nil {
		return err
	}
	return s.SaveSite(site)
}

func (s *Store) GetSite(code string) (model.CollectionSite, error) {
	var site model.CollectionSite
	err := s.Load(siteKind, code, &site)
	if errors.Is(err, sql.ErrNoRows) {
		return model.CollectionSite{}, fmt.Errorf("site %s: %w", code, model.ErrNotFound)
	}
	if err != nil {
		return model.CollectionSite{}, fmt.Errorf("load site %s: %w", code, err)
	}
	return site, nil
}

func (s *Store) ListSites() ([]model.CollectionSite, error) {
	items, err := decodeList[model.CollectionSite](s, siteKind)
	if err != nil {
		return nil, fmt.Errorf("list sites: %w", err)
	}
	return items, nil
}

func (s *Store) ListActiveSites() ([]model.CollectionSite, error) {
	all, err := s.ListSites()
	if err != nil {
		return nil, err
	}
	result := make([]model.CollectionSite, 0, len(all))
	for _, site := range all {
		if site.Active {
			result = append(result, site)
		}
	}
	return result, nil
}

func (s *Store) DeleteSite(code string) error {
	if _, err := s.GetSite(code); err != nil {
		return err
	}
	return s.Delete(siteKind, code)
}
