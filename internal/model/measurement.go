package model

import "time"

type Measurement struct {
	ID          string    `json:"id"`
	SectionID   string    `json:"section_id"`
	Kind        string    `json:"kind"`
	Value       float64   `json:"value"`
	Unit        string    `json:"unit"`
	Instrument  string    `json:"instrument"`
	MeasuredAt  time.Time `json:"measured_at"`
	Uncertainty float64   `json:"uncertainty"`
	Operator    string    `json:"operator"`
}

func (m Measurement) Validate() error {
	if m.ID == "" || m.SectionID == "" || m.Kind == "" || m.Unit == "" || m.Instrument == "" || m.Operator == "" || m.MeasuredAt.IsZero() {
		return ErrInvalidInput
	}
	if m.Uncertainty < 0 || m.Uncertainty > 1000000 {
		return ErrInvalidInput
	}
	return nil
}

func (m Measurement) LowerBound() float64 {
	return m.Value - m.Uncertainty
}

func (m Measurement) UpperBound() float64 {
	return m.Value + m.Uncertainty
}

type MeasurementSet struct {
	SectionID string        `json:"section_id"`
	Items     []Measurement `json:"items"`
}

func (s MeasurementSet) Mean(kind string) (float64, bool) {
	var total float64
	var count int
	for _, item := range s.Items {
		if item.Kind == kind {
			total += item.Value
			count++
		}
	}
	if count == 0 {
		return 0, false
	}
	return total / float64(count), true
}

func (s MeasurementSet) ByKind(kind string) []Measurement {
	result := make([]Measurement, 0)
	for _, item := range s.Items {
		if item.Kind == kind {
			result = append(result, item)
		}
	}
	return result
}
