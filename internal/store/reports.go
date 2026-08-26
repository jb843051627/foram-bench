package store

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/jb843051627/foram-bench/internal/model"
)

const reportKind = "report"

func (s *Store) SaveReport(report model.Report) error {
	if err := report.Validate(); err != nil {
		return err
	}
	return s.Save(reportKind, report.ID, report)
}

func (s *Store) InsertReport(report model.Report) error {
	if err := report.Validate(); err != nil {
		return err
	}
	if err := s.Insert(reportKind, report.ID, report); err != nil {
		return fmt.Errorf("insert report %s: %w", report.ID, model.ErrConflict)
	}
	return nil
}

func (s *Store) GetReport(id string) (model.Report, error) {
	var report model.Report
	err := s.Load(reportKind, id, &report)
	if errors.Is(err, sql.ErrNoRows) {
		return model.Report{}, fmt.Errorf("report %s: %w", id, model.ErrNotFound)
	}
	if err != nil {
		return model.Report{}, fmt.Errorf("load report %s: %w", id, err)
	}
	return report, nil
}

func (s *Store) ListReports() ([]model.Report, error) {
	items, err := decodeList[model.Report](s, reportKind)
	if err != nil {
		return nil, fmt.Errorf("list reports: %w", err)
	}
	return items, nil
}

func (s *Store) ListReportsBySample(sampleID string) ([]model.Report, error) {
	all, err := s.ListReports()
	if err != nil {
		return nil, err
	}
	out := make([]model.Report, 0, len(all))
	for _, report := range all {
		if report.SampleID == sampleID {
			out = append(out, report)
		}
	}
	return out, nil
}
