package model

import "time"

type Report struct {
	ID          string    `json:"id"`
	SampleID    string    `json:"sample_id"`
	Version     int       `json:"version"`
	Status      string    `json:"status"`
	TimeZone    string    `json:"time_zone"`
	Body        string    `json:"body"`
	GeneratedAt time.Time `json:"generated_at"`
}

func (r Report) Validate() error {
	if r.ID == "" || r.SampleID == "" || r.Version < 1 || r.Status == "" || r.TimeZone == "" || r.Body == "" || r.GeneratedAt.IsZero() {
		return ErrInvalidInput
	}
	return nil
}
