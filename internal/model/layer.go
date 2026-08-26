package model

import "time"

type StratigraphicLayer struct {
	ID            string    `json:"id"`
	SiteCode      string    `json:"site_code"`
	Name          string    `json:"name"`
	TopDepthMM    int       `json:"top_depth_mm"`
	BottomDepthMM int       `json:"bottom_depth_mm"`
	Lithology     string    `json:"lithology"`
	Color         string    `json:"color"`
	Boundary      string    `json:"boundary"`
	Confidence    float64   `json:"confidence"`
	RecordedAt    time.Time `json:"recorded_at"`
}

func (l StratigraphicLayer) Validate() error {
	if l.ID == "" || l.SiteCode == "" || l.Name == "" || l.Lithology == "" || l.Color == "" || l.Boundary == "" || l.RecordedAt.IsZero() {
		return ErrInvalidInput
	}
	if l.TopDepthMM < 0 || l.BottomDepthMM <= l.TopDepthMM || l.Confidence < 0 || l.Confidence > 1 {
		return ErrInvalidInput
	}
	return nil
}

func (l StratigraphicLayer) ThicknessMM() int {
	return l.BottomDepthMM - l.TopDepthMM
}

func (l StratigraphicLayer) ContainsDepth(depthMM int) bool {
	return depthMM >= l.TopDepthMM && depthMM < l.BottomDepthMM
}

type LayerMarker struct {
	LayerID     string    `json:"layer_id"`
	Marker      string    `json:"marker"`
	DepthMM     int       `json:"depth_mm"`
	ObservedAt  time.Time `json:"observed_at"`
	Reliability float64   `json:"reliability"`
}

func (m LayerMarker) Validate() error {
	if m.LayerID == "" || m.Marker == "" || m.DepthMM < 0 || m.ObservedAt.IsZero() || m.Reliability < 0 || m.Reliability > 1 {
		return ErrInvalidInput
	}
	return nil
}
