package model

import "time"

type PreparationBatch struct {
	ID          string     `json:"id"`
	SampleID    string     `json:"sample_id"`
	Protocol    string     `json:"protocol"`
	Operator    string     `json:"operator"`
	Status      string     `json:"status"`
	Revision    int        `json:"revision"`
	StartedAt   *time.Time `json:"started_at,omitempty"`
	CompletedAt *time.Time `json:"completed_at,omitempty"`
	Notes       string     `json:"notes"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

func (b PreparationBatch) Validate() error {
	if b.ID == "" || b.SampleID == "" || b.Protocol == "" || b.Operator == "" {
		return ErrInvalidInput
	}
	if !ValidBatchStatus(b.Status) || b.Revision < 1 {
		return ErrInvalidInput
	}
	return nil
}
