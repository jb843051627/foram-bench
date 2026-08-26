package store

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/jb843051627/foram-bench/internal/model"
)

const reviewNoteKind = "review_note"

func (s *Store) SaveReviewNote(note model.ReviewNote) error {
	if err := note.Validate(); err != nil {
		return err
	}
	return s.Save(reviewNoteKind, note.ID, note)
}

func (s *Store) GetReviewNote(id string) (model.ReviewNote, error) {
	var note model.ReviewNote
	err := s.Load(reviewNoteKind, id, &note)
	if errors.Is(err, sql.ErrNoRows) {
		return model.ReviewNote{}, fmt.Errorf("review note %s: %w", id, model.ErrNotFound)
	}
	if err != nil {
		return model.ReviewNote{}, fmt.Errorf("load review note %s: %w", id, err)
	}
	return note, nil
}

func (s *Store) ListReviewNotes(reviewID string) ([]model.ReviewNote, error) {
	all, err := decodeList[model.ReviewNote](s, reviewNoteKind)
	if err != nil {
		return nil, fmt.Errorf("list review notes: %w", err)
	}
	result := make([]model.ReviewNote, 0)
	for _, note := range all {
		if note.ReviewID == reviewID {
			result = append(result, note)
		}
	}
	return result, nil
}
