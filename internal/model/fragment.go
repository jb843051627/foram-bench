package model

import "time"

type Fragment struct {
	ID         string    `json:"id"`
	SectionID  string    `json:"section_id"`
	PositionMM float64   `json:"position_mm"`
	LengthMM   float64   `json:"length_mm"`
	Color      string    `json:"color"`
	Texture    string    `json:"texture"`
	Damaged    bool      `json:"damaged"`
	Note       string    `json:"note"`
	CreatedAt  time.Time `json:"created_at"`
}

func (f Fragment) Validate() error {
	if f.ID == "" || f.SectionID == "" || f.Color == "" || f.Texture == "" || f.CreatedAt.IsZero() {
		return ErrInvalidInput
	}
	if f.PositionMM < 0 || f.LengthMM <= 0 {
		return ErrInvalidInput
	}
	return nil
}

func (f Fragment) EndMM() float64 {
	return f.PositionMM + f.LengthMM
}

type FragmentSet struct {
	SectionID string     `json:"section_id"`
	Items     []Fragment `json:"items"`
}

func (s FragmentSet) Coverage() float64 {
	if len(s.Items) == 0 {
		return 0
	}
	var total float64
	for _, item := range s.Items {
		total += item.LengthMM
	}
	return total
}
