package service

import (
	"context"
	"fmt"

	"github.com/jb843051627/foram-bench/internal/model"
	"github.com/jb843051627/foram-bench/internal/validation"
)

type ProtocolInput struct {
	ID             string               `json:"id"`
	Name           string               `json:"name"`
	Version        int                  `json:"version"`
	TargetMaterial []string             `json:"target_material"`
	RequiredStain  string               `json:"required_stain"`
	MaxThicknessUM int                  `json:"max_thickness_um"`
	Steps          []model.ProtocolStep `json:"steps"`
}

func (l *Lab) CreateProtocol(ctx context.Context, input ProtocolInput) (model.PreparationProtocol, error) {
	if err := checkContext(ctx); err != nil {
		return model.PreparationProtocol{}, err
	}
	steps := make([]validation.StepInput, 0, len(input.Steps))
	for _, step := range input.Steps {
		steps = append(steps, validation.StepInput{Number: step.Number, Name: step.Name, DurationMin: step.DurationMin,
			TemperatureC: step.TemperatureC, Instruction: step.Instruction})
	}
	if err := validation.Protocol(validation.ProtocolInput{Name: input.Name, Version: input.Version, RequiredStain: input.RequiredStain,
		MaxThicknessUM: input.MaxThicknessUM, Steps: steps}); err != nil {
		return model.PreparationProtocol{}, err
	}
	protocol := model.PreparationProtocol{ID: input.ID, Name: input.Name, Version: input.Version, TargetMaterial: input.TargetMaterial,
		RequiredStain: input.RequiredStain, MaxThicknessUM: input.MaxThicknessUM, Steps: input.Steps, UpdatedAt: l.clock.Now()}
	if err := protocol.Validate(); err != nil {
		return model.PreparationProtocol{}, fmt.Errorf("validate protocol: %w", err)
	}
	if err := l.store.SaveProtocol(protocol); err != nil {
		return model.PreparationProtocol{}, err
	}
	if err := l.store.Event(protocol.ID, "protocol.created"); err != nil {
		return model.PreparationProtocol{}, err
	}
	return protocol, nil
}

func (l *Lab) GetProtocol(ctx context.Context, id string) (model.PreparationProtocol, error) {
	if err := checkContext(ctx); err != nil {
		return model.PreparationProtocol{}, err
	}
	return l.store.GetProtocol(id)
}

func (l *Lab) LatestProtocol(ctx context.Context, name string) (model.PreparationProtocol, error) {
	if err := checkContext(ctx); err != nil {
		return model.PreparationProtocol{}, err
	}
	return l.store.LatestProtocol(name)
}

func (l *Lab) ProtocolStep(ctx context.Context, id string, number int) (model.ProtocolStep, error) {
	protocol, err := l.GetProtocol(ctx, id)
	if err != nil {
		return model.ProtocolStep{}, err
	}
	step, ok := protocol.Step(number)
	if !ok {
		return model.ProtocolStep{}, fmt.Errorf("protocol step %d: %w", number, model.ErrNotFound)
	}
	return step, nil
}
