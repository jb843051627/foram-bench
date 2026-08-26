package model

import "time"

type SectionReview struct {
	ID         string    `json:"id"`
	SectionID  string    `json:"section_id"`
	Reviewer   string    `json:"reviewer"`
	Decision   string    `json:"decision"`
	Comment    string    `json:"comment"`
	ReviewedAt time.Time `json:"reviewed_at"`
}

func (r SectionReview) Validate() error {
	if r.ID == "" || r.SectionID == "" || r.Reviewer == "" || r.Comment == "" || r.ReviewedAt.IsZero() {
		return ErrInvalidInput
	}
	if !ValidReviewDecision(r.Decision) {
		return ErrInvalidInput
	}
	return nil
}
