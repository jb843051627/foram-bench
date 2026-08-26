package model

import "strings"

type TaxonRecord struct {
	Name        string   `json:"name"`
	Rank        string   `json:"rank"`
	Authority   string   `json:"authority"`
	Synonyms    []string `json:"synonyms"`
	Environment []string `json:"environment"`
	Notes       string   `json:"notes"`
}

func (t TaxonRecord) Validate() error {
	if strings.TrimSpace(t.Name) == "" || strings.TrimSpace(t.Rank) == "" || strings.TrimSpace(t.Authority) == "" {
		return ErrInvalidInput
	}
	return nil
}

func (t TaxonRecord) Matches(value string) bool {
	value = strings.TrimSpace(value)
	if strings.EqualFold(value, t.Name) {
		return true
	}
	for _, synonym := range t.Synonyms {
		if strings.EqualFold(value, synonym) {
			return true
		}
	}
	return false
}

type TaxonCount struct {
	Taxon        string  `json:"taxon"`
	Count        int     `json:"count"`
	Relative     float64 `json:"relative"`
	Confidence   float64 `json:"confidence"`
	Preservation string  `json:"preservation"`
}

func (t TaxonCount) WeightedCount() float64 {
	if t.Confidence <= 0 {
		return 0
	}
	return float64(t.Count) * t.Confidence
}
