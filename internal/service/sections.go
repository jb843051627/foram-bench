package service

import (
	"context"
	"fmt"

	"github.com/jb843051627/foram-bench/internal/model"
)

func (l *Lab) CreateSection(ctx context.Context, input SectionInput) (model.ThinSection, error) {
	if err := checkContext(ctx); err != nil {
		return model.ThinSection{}, err
	}
	batch, err := l.store.GetBatch(input.BatchID)
	if err != nil {
		return model.ThinSection{}, fmt.Errorf("load batch for section: %w", err)
	}
	if batch.Status != model.BatchReady {
		return model.ThinSection{}, fmt.Errorf("batch %s is not ready: %w", batch.ID, model.ErrInvalidState)
	}
	now := l.clock.Now()
	section := model.ThinSection{ID: input.ID, BatchID: input.BatchID, Label: input.Label,
		ThicknessUM: input.ThicknessUM, Stain: input.Stain, Status: model.SectionCut,
		CreatedAt: now, UpdatedAt: now}
	if err := section.Validate(); err != nil {
		return model.ThinSection{}, err
	}
	if err := l.store.SaveSection(section); err != nil {
		return model.ThinSection{}, err
	}
	if err := l.store.Event(section.ID, "section.created"); err != nil {
		return model.ThinSection{}, err
	}
	l.metrics.Add("sections.created", 1)
	return section, nil
}

func (l *Lab) GetSection(ctx context.Context, id string) (model.ThinSection, error) {
	if err := checkContext(ctx); err != nil {
		return model.ThinSection{}, err
	}
	return l.store.GetSection(id)
}

func (l *Lab) ListSections(ctx context.Context, batchID string) ([]model.ThinSection, error) {
	if err := checkContext(ctx); err != nil {
		return nil, err
	}
	return l.store.ListSectionsByBatch(batchID)
}

func (l *Lab) AdvanceSection(ctx context.Context, id, target string) (model.ThinSection, error) {
	l.sectionMu.Lock()
	defer l.sectionMu.Unlock()
	if err := checkContext(ctx); err != nil {
		return model.ThinSection{}, err
	}
	section, err := l.store.GetSection(id)
	if err != nil {
		return model.ThinSection{}, err
	}
	if !model.CanMoveSection(section.Status, target) {
		return model.ThinSection{}, fmt.Errorf("section %s cannot move %s -> %s: %w", id, section.Status, target, model.ErrInvalidState)
	}
	section.Status = target
	section.UpdatedAt = l.clock.Now()
	if err := l.store.SaveSection(section); err != nil {
		return model.ThinSection{}, err
	}
	if err := l.store.Event(id, fmt.Sprintf("section.moved_to_%s", target)); err != nil {
		return model.ThinSection{}, err
	}
	return section, nil
}

func (l *Lab) StainSection(ctx context.Context, id string) (model.ThinSection, error) {
	return l.AdvanceSection(ctx, id, model.SectionStained)
}

func (l *Lab) MarkSectionReviewed(ctx context.Context, id string) (model.ThinSection, error) {
	return l.AdvanceSection(ctx, id, model.SectionReviewed)
}
