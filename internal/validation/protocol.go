package validation

import (
	"fmt"
	"sort"
	"strings"
)

type ProtocolInput struct {
	Name           string
	Version        int
	RequiredStain  string
	MaxThicknessUM int
	Steps          []StepInput
}

type StepInput struct {
	Number       int
	Name         string
	DurationMin  int
	TemperatureC float64
	Instruction  string
}

func Protocol(input ProtocolInput) error {
	if strings.TrimSpace(input.Name) == "" || input.Version < 1 || strings.TrimSpace(input.RequiredStain) == "" {
		return fmt.Errorf("protocol identity is incomplete")
	}
	if input.MaxThicknessUM < 1 || input.MaxThicknessUM > 1000 || len(input.Steps) == 0 {
		return fmt.Errorf("protocol limits or steps are invalid")
	}
	numbers := make([]int, 0, len(input.Steps))
	for _, step := range input.Steps {
		if step.Number <= 0 || step.Name == "" || step.DurationMin <= 0 || step.Instruction == "" {
			return fmt.Errorf("protocol step %d is incomplete", step.Number)
		}
		if err := Range(step.TemperatureC, -30, 180, "temperature"); err != nil {
			return err
		}
		numbers = append(numbers, step.Number)
	}
	sort.Ints(numbers)
	for index, number := range numbers {
		if number != index+1 {
			return fmt.Errorf("protocol steps are not contiguous")
		}
	}
	return nil
}
