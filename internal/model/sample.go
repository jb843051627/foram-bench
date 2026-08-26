package model

import "time"

type Sample struct {
	ID             string    `json:"id"`
	SiteCode       string    `json:"site_code"`
	DepthMM        int       `json:"depth_mm"`
	Material       string    `json:"material"`
	CollectionDate time.Time `json:"collection_date"`
	Location       string    `json:"location"`
	TimeZone       string    `json:"time_zone"`
	Status         string    `json:"status"`
	Notes          string    `json:"notes"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

func (s Sample) Validate() error {
	if s.ID == "" || s.SiteCode == "" || s.Material == "" || s.Location == "" {
		return ErrInvalidInput
	}
	if s.DepthMM < 0 || s.DepthMM > 100000 {
		return ErrInvalidInput
	}
	if s.CollectionDate.IsZero() || !ValidSampleStatus(s.Status) || s.TimeZone == "" {
		return ErrInvalidInput
	}
	return nil
}
