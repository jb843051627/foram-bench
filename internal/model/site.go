package model

import "time"

type CollectionSite struct {
	Code         string    `json:"code"`
	Name         string    `json:"name"`
	Region       string    `json:"region"`
	Latitude     float64   `json:"latitude"`
	Longitude    float64   `json:"longitude"`
	ElevationM   int       `json:"elevation_m"`
	TimeZone     string    `json:"time_zone"`
	Stratigraphy string    `json:"stratigraphy"`
	Active       bool      `json:"active"`
	CreatedAt    time.Time `json:"created_at"`
}

func (s CollectionSite) Validate() error {
	if s.Code == "" || s.Name == "" || s.Region == "" || s.TimeZone == "" || s.Stratigraphy == "" {
		return ErrInvalidInput
	}
	if s.Latitude < -90 || s.Latitude > 90 || s.Longitude < -180 || s.Longitude > 180 {
		return ErrInvalidInput
	}
	if s.ElevationM < -1000 || s.ElevationM > 10000 || s.CreatedAt.IsZero() {
		return ErrInvalidInput
	}
	return nil
}

type SiteInterval struct {
	SiteCode string    `json:"site_code"`
	Name     string    `json:"name"`
	Start    time.Time `json:"start"`
	End      time.Time `json:"end"`
	Unit     string    `json:"unit"`
}

func (i SiteInterval) Duration() time.Duration {
	return i.End.Sub(i.Start)
}

func (i SiteInterval) Valid() bool {
	return i.SiteCode != "" && i.Unit != "" && !i.Start.IsZero() && i.End.After(i.Start)
}
