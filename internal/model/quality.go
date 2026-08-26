package model

import "time"

type QualityFlag struct {
	ID         string     `json:"id"`
	BatchID    string     `json:"batch_id"`
	Kind       string     `json:"kind"`
	Severity   string     `json:"severity"`
	Detail     string     `json:"detail"`
	Resolved   bool       `json:"resolved"`
	CreatedAt  time.Time  `json:"created_at"`
	ResolvedAt *time.Time `json:"resolved_at,omitempty"`
}

func (f QualityFlag) Validate() error {
	if f.ID == "" || f.BatchID == "" || f.Kind == "" || f.Severity == "" || f.Detail == "" || f.CreatedAt.IsZero() {
		return ErrInvalidInput
	}
	return nil
}
