package service

import (
	"context"
	"fmt"

	"github.com/jb843051627/foram-bench/internal/model"
)

func (l *Lab) ReviewSection(ctx context.Context, input ReviewInput) (model.SectionReview, error) {
	l.reviewMu.Lock()
	defer l.reviewMu.Unlock()
	if err := checkContext(ctx); err != nil {
		return model.SectionReview{}, err
	}
	section, err := l.store.GetSection(input.SectionID)
	if err != nil {
		return model.SectionReview{}, err
	}
	if section.Status != model.SectionReviewed {
		return model.SectionReview{}, fmt.Errorf("section %s is not ready for review: %w", section.ID, model.ErrInvalidState)
	}
	reviews, err := l.store.ListReviewsBySection(section.ID)
	if err != nil {
		return model.SectionReview{}, err
	}
	for _, review := range reviews {
		if review.Decision == model.ReviewAccepted {
			return model.SectionReview{}, fmt.Errorf("section %s already accepted: %w", section.ID, model.ErrAlreadyReviewed)
		}
	}
	review := model.SectionReview{ID: input.ID, SectionID: input.SectionID, Reviewer: input.Reviewer,
		Decision: input.Decision, Comment: input.Comment, ReviewedAt: l.clock.Now()}
	if err := review.Validate(); err != nil {
		return model.SectionReview{}, err
	}
	if err := l.store.SaveReview(review); err != nil {
		return model.SectionReview{}, err
	}
	if err := l.store.Event(section.ID, fmt.Sprintf("section.review_%s", review.Decision)); err != nil {
		return model.SectionReview{}, err
	}
	return review, nil
}

func (l *Lab) GetReview(ctx context.Context, id string) (model.SectionReview, error) {
	if err := checkContext(ctx); err != nil {
		return model.SectionReview{}, err
	}
	return l.store.GetReview(id)
}

func (l *Lab) ListReviews(ctx context.Context, sectionID string) ([]model.SectionReview, error) {
	if err := checkContext(ctx); err != nil {
		return nil, err
	}
	return l.store.ListReviewsBySection(sectionID)
}

func (l *Lab) HasAcceptedReview(ctx context.Context, sectionID string) (bool, error) {
	reviews, err := l.ListReviews(ctx, sectionID)
	if err != nil {
		return false, err
	}
	for _, review := range reviews {
		if review.Decision == model.ReviewAccepted {
			return true, nil
		}
	}
	return false, nil
}
