package model

import "time"

type ThinSection struct {
	ID          string    `json:"id"`
	BatchID     string    `json:"batch_id"`
	Label       string    `json:"label"`
	ThicknessUM int       `json:"thickness_um"`
	Stain       string    `json:"stain"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (s ThinSection) Validate() error {
	if s.ID == "" || s.BatchID == "" || s.Label == "" || s.Stain == "" {
		return ErrInvalidInput
	}
	if s.ThicknessUM < 1 || s.ThicknessUM > 1000 {
		return ErrInvalidInput
	}
	if s.Status == "" {
		return ErrInvalidInput
	}
	return nil
}
