package model

import "time"

type PreparationProtocol struct {
	ID             string         `json:"id"`
	Name           string         `json:"name"`
	Version        int            `json:"version"`
	TargetMaterial []string       `json:"target_material"`
	Steps          []ProtocolStep `json:"steps"`
	RequiredStain  string         `json:"required_stain"`
	MaxThicknessUM int            `json:"max_thickness_um"`
	UpdatedAt      time.Time      `json:"updated_at"`
}

type ProtocolStep struct {
	Number       int     `json:"number"`
	Name         string  `json:"name"`
	TemperatureC float64 `json:"temperature_c"`
	DurationMin  int     `json:"duration_min"`
	Required     bool    `json:"required"`
	Instruction  string  `json:"instruction"`
}

func (p PreparationProtocol) Validate() error {
	if p.ID == "" || p.Name == "" || p.Version < 1 || p.RequiredStain == "" || p.MaxThicknessUM < 1 || p.UpdatedAt.IsZero() {
		return ErrInvalidInput
	}
	if len(p.TargetMaterial) == 0 || len(p.Steps) == 0 {
		return ErrInvalidInput
	}
	for index, step := range p.Steps {
		if step.Number != index+1 || step.Name == "" || step.DurationMin <= 0 || step.Instruction == "" {
			return ErrInvalidInput
		}
	}
	return nil
}

func (p PreparationProtocol) Step(number int) (ProtocolStep, bool) {
	for _, step := range p.Steps {
		if step.Number == number {
			return step, true
		}
	}
	return ProtocolStep{}, false
}
