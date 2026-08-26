package service

import "time"

type SampleInput struct {
	ID             string    `json:"id"`
	SiteCode       string    `json:"site_code"`
	DepthMM        int       `json:"depth_mm"`
	Material       string    `json:"material"`
	CollectionDate time.Time `json:"collection_date"`
	Location       string    `json:"location"`
	TimeZone       string    `json:"time_zone"`
	Notes          string    `json:"notes"`
}

type BatchInput struct {
	ID       string `json:"id"`
	SampleID string `json:"sample_id"`
	Protocol string `json:"protocol"`
	Operator string `json:"operator"`
	Notes    string `json:"notes"`
}

type SectionInput struct {
	ID          string `json:"id"`
	BatchID     string `json:"batch_id"`
	Label       string `json:"label"`
	ThicknessUM int    `json:"thickness_um"`
	Stain       string `json:"stain"`
}

type ObservationInput struct {
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

type ReviewInput struct {
	ID        string `json:"id"`
	SectionID string `json:"section_id"`
	Reviewer  string `json:"reviewer"`
	Decision  string `json:"decision"`
	Comment   string `json:"comment"`
}

type QualityInput struct {
	ID       string `json:"id"`
	BatchID  string `json:"batch_id"`
	Kind     string `json:"kind"`
	Severity string `json:"severity"`
	Detail   string `json:"detail"`
}

type BatchSummary struct {
	Batch        any `json:"batch"`
	Sample       any `json:"sample"`
	Sections     int `json:"sections"`
	Observations int `json:"observations"`
	Reviews      int `json:"reviews"`
	OpenFlags    int `json:"open_flags"`
}
