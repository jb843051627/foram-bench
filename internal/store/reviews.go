package store

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/jb843051627/foram-bench/internal/model"
)

const reviewKind = "review"

func (s *Store) SaveReview(review model.SectionReview) error {
	if err := review.Validate(); err != nil {
		return err
	}
	return s.Save(reviewKind, review.ID, review)
}

func (s *Store) InsertReview(review model.SectionReview) error {
	if err := review.Validate(); err != nil {
		return fmt.Errorf("validate review: %w", err)
	}
	if err := s.Save(reviewKind, review.ID, review); err != nil {
		return err
	}
	return nil
}

func (s *Store) GetReview(id string) (model.SectionReview, error) {
	var review model.SectionReview
	err := s.Load(reviewKind, id, &review)
	if errors.Is(err, sql.ErrNoRows) {
		return model.SectionReview{}, fmt.Errorf("review %s: %w", id, model.ErrNotFound)
	}
	if err != nil {
		return model.SectionReview{}, fmt.Errorf("load review %s: %w", id, err)
	}
	return review, nil
}

func (s *Store) ListReviews() ([]model.SectionReview, error) {
	items, err := decodeList[model.SectionReview](s, reviewKind)
	if err != nil {
		return nil, fmt.Errorf("list reviews: %w", err)
	}
	return items, nil
}

func (s *Store) ListReviewsBySection(sectionID string) ([]model.SectionReview, error) {
	all, err := s.ListReviews()
	if err != nil {
		return nil, err
	}
	out := make([]model.SectionReview, 0, len(all))
	for _, review := range all {
		if review.SectionID == sectionID {
			out = append(out, review)
		}
	}
	return out, nil
}
