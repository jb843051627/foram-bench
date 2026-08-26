package model

import "time"

type ReviewNote struct {
	ID         string     `json:"id"`
	ReviewID   string     `json:"review_id"`
	Author     string     `json:"author"`
	Category   string     `json:"category"`
	Body       string     `json:"body"`
	Resolved   bool       `json:"resolved"`
	CreatedAt  time.Time  `json:"created_at"`
	ResolvedAt *time.Time `json:"resolved_at,omitempty"`
}

func (n ReviewNote) Validate() error {
	if n.ID == "" || n.ReviewID == "" || n.Author == "" || n.Category == "" || n.Body == "" || n.CreatedAt.IsZero() {
		return ErrInvalidInput
	}
	return nil
}
