package model

import "time"

type Diagnostic struct {
	ID              string    `json:"id"`
	BatchID         string    `json:"batch_id"`
	Code            string    `json:"code"`
	Severity        string    `json:"severity"`
	Title           string    `json:"title"`
	Detail          string    `json:"detail"`
	SuggestedAction string    `json:"suggested_action"`
	Acknowledged    bool      `json:"acknowledged"`
	CreatedAt       time.Time `json:"created_at"`
}

func (d Diagnostic) Validate() error {
	if d.ID == "" || d.BatchID == "" || d.Code == "" || d.Severity == "" || d.Title == "" || d.Detail == "" || d.SuggestedAction == "" || d.CreatedAt.IsZero() {
		return ErrInvalidInput
	}
	return nil
}

func (d Diagnostic) BlocksRelease() bool {
	return !d.Acknowledged && (d.Severity == "critical" || d.Severity == "major")
}

type ChecklistItem struct {
	Code     string `json:"code"`
	Label    string `json:"label"`
	Required bool   `json:"required"`
	Passed   bool   `json:"passed"`
	Note     string `json:"note"`
}
