package model

import "time"

type Assessment struct {
	ID           string    `json:"id"`
	BatchID      string    `json:"batch_id"`
	Completeness float64   `json:"completeness"`
	Preservation float64   `json:"preservation"`
	Diversity    float64   `json:"diversity"`
	Score        float64   `json:"score"`
	Grade        string    `json:"grade"`
	Findings     []string  `json:"findings"`
	AssessedAt   time.Time `json:"assessed_at"`
}

func (a Assessment) Validate() error {
	if a.ID == "" || a.BatchID == "" || a.Grade == "" || a.AssessedAt.IsZero() {
		return ErrInvalidInput
	}
	for _, value := range []float64{a.Completeness, a.Preservation, a.Diversity, a.Score} {
		if value < 0 || value > 100 {
			return ErrInvalidInput
		}
	}
	return nil
}

func (a Assessment) Passed(threshold float64) bool {
	return a.Score >= threshold && a.Completeness >= 70 && a.Preservation >= 60
}
