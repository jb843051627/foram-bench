package store

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/jb843051627/foram-bench/internal/model"
)

const assessmentKind = "assessment"

func (s *Store) SaveAssessment(assessment model.Assessment) error {
	if err := assessment.Validate(); err != nil {
		return err
	}
	return s.Save(assessmentKind, assessment.ID, assessment)
}

func (s *Store) GetAssessment(id string) (model.Assessment, error) {
	var assessment model.Assessment
	err := s.Load(assessmentKind, id, &assessment)
	if errors.Is(err, sql.ErrNoRows) {
		return model.Assessment{}, fmt.Errorf("assessment %s: %w", id, model.ErrNotFound)
	}
	if err != nil {
		return model.Assessment{}, fmt.Errorf("load assessment %s: %w", id, err)
	}
	return assessment, nil
}

func (s *Store) ListAssessments() ([]model.Assessment, error) {
	items, err := decodeList[model.Assessment](s, assessmentKind)
	if err != nil {
		return nil, fmt.Errorf("list assessments: %w", err)
	}
	return items, nil
}

func (s *Store) ListAssessmentsByBatch(batchID string) ([]model.Assessment, error) {
	all, err := s.ListAssessments()
	if err != nil {
		return nil, err
	}
	result := make([]model.Assessment, 0, len(all))
	for _, assessment := range all {
		if assessment.BatchID == batchID {
			result = append(result, assessment)
		}
	}
	return result, nil
}
