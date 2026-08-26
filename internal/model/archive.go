package model

import "time"

type ArchiveRecord struct {
	ID          string    `json:"id"`
	SampleID    string    `json:"sample_id"`
	Repository  string    `json:"repository"`
	Accession   string    `json:"accession"`
	ArchiveDate time.Time `json:"archive_date"`
	Custodian   string    `json:"custodian"`
	Condition   string    `json:"condition"`
	Notes       string    `json:"notes"`
}

func (a ArchiveRecord) Validate() error {
	if a.ID == "" || a.SampleID == "" || a.Repository == "" || a.Accession == "" || a.ArchiveDate.IsZero() || a.Custodian == "" || a.Condition == "" {
		return ErrInvalidInput
	}
	return nil
}

func (a ArchiveRecord) Stable() bool {
	return a.Condition == "stable" || a.Condition == "monitored"
}
