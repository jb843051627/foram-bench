package store

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/jb843051627/foram-bench/internal/model"
)

const archiveKind = "archive"

func (s *Store) SaveArchive(record model.ArchiveRecord) error {
	if err := record.Validate(); err != nil {
		return err
	}
	return s.Save(archiveKind, record.ID, record)
}

func (s *Store) GetArchive(id string) (model.ArchiveRecord, error) {
	var record model.ArchiveRecord
	err := s.Load(archiveKind, id, &record)
	if errors.Is(err, sql.ErrNoRows) {
		return model.ArchiveRecord{}, fmt.Errorf("archive %s: %w", id, model.ErrNotFound)
	}
	if err != nil {
		return model.ArchiveRecord{}, err
	}
	return record, nil
}

func (s *Store) ListArchives(sampleID string) ([]model.ArchiveRecord, error) {
	all, err := decodeList[model.ArchiveRecord](s, archiveKind)
	if err != nil {
		return nil, err
	}
	result := make([]model.ArchiveRecord, 0)
	for _, record := range all {
		if sampleID == "" || record.SampleID == sampleID {
			result = append(result, record)
		}
	}
	return result, nil
}
