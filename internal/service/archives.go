package service

import (
	"context"
	"fmt"

	"github.com/jb843051627/foram-bench/internal/model"
)

type ArchiveInput struct {
	ID         string `json:"id"`
	SampleID   string `json:"sample_id"`
	Repository string `json:"repository"`
	Accession  string `json:"accession"`
	Custodian  string `json:"custodian"`
	Condition  string `json:"condition"`
	Notes      string `json:"notes"`
}

func (l *Lab) ArchiveSampleRecord(ctx context.Context, input ArchiveInput) (model.ArchiveRecord, error) {
	if err := checkContext(ctx); err != nil {
		return model.ArchiveRecord{}, err
	}
	if _, err := l.store.GetSample(input.SampleID); err != nil {
		return model.ArchiveRecord{}, fmt.Errorf("load sample for archive: %w", err)
	}
	record := model.ArchiveRecord{ID: input.ID, SampleID: input.SampleID, Repository: input.Repository, Accession: input.Accession,
		Custodian: input.Custodian, Condition: input.Condition, Notes: input.Notes, ArchiveDate: l.clock.Now()}
	if err := record.Validate(); err != nil {
		return model.ArchiveRecord{}, err
	}
	if err := l.store.SaveArchive(record); err != nil {
		return model.ArchiveRecord{}, err
	}
	if err := l.store.Event(input.SampleID, "sample.archived_recorded"); err != nil {
		return model.ArchiveRecord{}, err
	}
	return record, nil
}

func (l *Lab) GetArchive(ctx context.Context, id string) (model.ArchiveRecord, error) {
	if err := checkContext(ctx); err != nil {
		return model.ArchiveRecord{}, err
	}
	return l.store.GetArchive(id)
}

func (l *Lab) ListArchives(ctx context.Context, sampleID string) ([]model.ArchiveRecord, error) {
	if err := checkContext(ctx); err != nil {
		return nil, err
	}
	return l.store.ListArchives(sampleID)
}

func (l *Lab) ArchiveStable(ctx context.Context, sampleID string) (bool, error) {
	items, err := l.ListArchives(ctx, sampleID)
	if err != nil {
		return false, err
	}
	if len(items) == 0 {
		return false, nil
	}
	for _, item := range items {
		if item.Stable() {
			return true, nil
		}
	}
	return false, nil
}
