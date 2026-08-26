package model

import "time"

type Observation struct {
	ID           string    `json:"id"`
	SectionID    string    `json:"section_id"`
	Observer     string    `json:"observer"`
	Taxon        string    `json:"taxon"`
	Count        int       `json:"count"`
	Preservation string    `json:"preservation"`
	Confidence   float64   `json:"confidence"`
	Notes        string    `json:"notes"`
	ObservedAt   time.Time `json:"observed_at"`
}

func (o Observation) Validate() error {
	if o.ID == "" || o.SectionID == "" || o.Observer == "" || o.Taxon == "" || o.ObservedAt.IsZero() {
		return ErrInvalidInput
	}
	if o.Count < 0 || o.Count > 1000000 || o.Confidence < 0 || o.Confidence > 1 {
		return ErrInvalidInput
	}
	if o.Preservation == "" {
		return ErrInvalidInput
	}
	return nil
}
