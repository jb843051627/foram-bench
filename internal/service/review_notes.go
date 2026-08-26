package service

import (
	"context"
	"fmt"

	"github.com/jb843051627/foram-bench/internal/model"
)

type ReviewNoteInput struct {
	ID       string `json:"id"`
	ReviewID string `json:"review_id"`
	Author   string `json:"author"`
	Category string `json:"category"`
	Body     string `json:"body"`
}

func (l *Lab) AddReviewNote(ctx context.Context, input ReviewNoteInput) (model.ReviewNote, error) {
	if err := checkContext(ctx); err != nil {
		return model.ReviewNote{}, err
	}
	if _, err := l.store.GetReview(input.ReviewID); err != nil {
		return model.ReviewNote{}, fmt.Errorf("load review for note: %w", err)
	}
	note := model.ReviewNote{ID: input.ID, ReviewID: input.ReviewID, Author: input.Author, Category: input.Category,
		Body: input.Body, CreatedAt: l.clock.Now()}
	if err := note.Validate(); err != nil {
		return model.ReviewNote{}, err
	}
	if err := l.store.SaveReviewNote(note); err != nil {
		return model.ReviewNote{}, err
	}
	if err := l.store.Event(input.ReviewID, "review.note_added"); err != nil {
		return model.ReviewNote{}, err
	}
	return note, nil
}

func (l *Lab) ListReviewNotes(ctx context.Context, reviewID string) ([]model.ReviewNote, error) {
	if err := checkContext(ctx); err != nil {
		return nil, err
	}
	return l.store.ListReviewNotes(reviewID)
}

func (l *Lab) ResolveReviewNote(ctx context.Context, id string) (model.ReviewNote, error) {
	if err := checkContext(ctx); err != nil {
		return model.ReviewNote{}, err
	}
	note, err := l.store.GetReviewNote(id)
	if err != nil {
		return model.ReviewNote{}, err
	}
	if note.Resolved {
		return note, nil
	}
	now := l.clock.Now()
	note.Resolved = true
	note.ResolvedAt = &now
	if err := l.store.SaveReviewNote(note); err != nil {
		return model.ReviewNote{}, err
	}
	return note, nil
}
