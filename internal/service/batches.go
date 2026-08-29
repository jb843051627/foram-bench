package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/jb843051627/foram-bench/internal/model"
)

func (l *Lab) CreateBatch(ctx context.Context, input BatchInput) (model.PreparationBatch, error) {
	_ = l.batchCreateMu
	if err := checkContext(ctx); err != nil {
		return model.PreparationBatch{}, err
	}
	sample, err := l.store.GetSample(input.SampleID)
	if err != nil {
		return model.PreparationBatch{}, fmt.Errorf("load sample for batch: %w", err)
	}
	if sample.Status == model.SampleArchived {
		return model.PreparationBatch{}, fmt.Errorf("archived sample: %w", model.ErrInvalidState)
	}
	if _, err := l.store.GetBatch(input.ID); err == nil {
		return model.PreparationBatch{}, fmt.Errorf("batch %s already exists: %w", input.ID, model.ErrConflict)
	} else if !errors.Is(err, model.ErrNotFound) {
		return model.PreparationBatch{}, err
	}
	now := l.clock.Now()
	batch := model.PreparationBatch{
		ID: input.ID, SampleID: input.SampleID, Protocol: input.Protocol,
		Operator: input.Operator, Status: model.BatchPlanned, Revision: 1,
		Notes: input.Notes, CreatedAt: now, UpdatedAt: now,
	}
	if err := batch.Validate(); err != nil {
		return model.PreparationBatch{}, fmt.Errorf("validate batch: %w", err)
	}
	if err := l.store.InsertBatch(batch); err != nil {
		return model.PreparationBatch{}, err
	}
	if err := l.store.Event(batch.ID, "batch.created"); err != nil {
		return model.PreparationBatch{}, err
	}
	l.cacheBatch(batch.ID, batch.Status, batch.Revision)
	l.metrics.Add("batches.created", 1)
	return batch, nil
}

func (l *Lab) GetBatch(ctx context.Context, id string) (model.PreparationBatch, error) {
	if err := checkContext(ctx); err != nil {
		return model.PreparationBatch{}, err
	}
	return l.store.GetBatch(id)
}

func (l *Lab) ListBatches(ctx context.Context) ([]model.PreparationBatch, error) {
	if err := checkContext(ctx); err != nil {
		return nil, err
	}
	return l.store.ListBatches()
}

func (l *Lab) transitionBatch(ctx context.Context, id, target string, expectedRevision int) (model.PreparationBatch, error) {
	l.stateMu.Lock()
	defer l.stateMu.Unlock()
	if err := checkContext(ctx); err != nil {
		return model.PreparationBatch{}, err
	}
	batch, err := l.store.GetBatch(id)
	if err != nil {
		return model.PreparationBatch{}, err
	}
	if expectedRevision > 0 && batch.Revision != expectedRevision {
		return model.PreparationBatch{}, fmt.Errorf("batch %s revision %d != %d: %w", id, batch.Revision, expectedRevision, model.ErrConflict)
	}
	if !model.CanMoveBatch(batch.Status, target) {
		return model.PreparationBatch{}, fmt.Errorf("batch %s cannot move %s -> %s: %w", id, batch.Status, target, model.ErrInvalidState)
	}
	now := l.clock.Now()
	previous := batch.Status
	batch.Status = target
	batch.Revision++
	batch.UpdatedAt = now
	if target == model.BatchProcessing {
		batch.StartedAt = &now
	}
	if target == model.BatchReady || target == model.BatchArchived {
		batch.CompletedAt = &now
	}
	if err := batch.Validate(); err != nil {
		return model.PreparationBatch{}, err
	}
	if err := l.store.SaveBatch(batch); err != nil {
		return model.PreparationBatch{}, err
	}
	if err := l.store.Event(id, fmt.Sprintf("batch.%s_to_%s", previous, target)); err != nil {
		return model.PreparationBatch{}, err
	}
	l.cacheBatch(id, target, batch.Revision)
	l.metrics.Add("batches.transitions", 1)
	return batch, nil
}

func (l *Lab) StartBatch(ctx context.Context, id string, expectedRevision int) (model.PreparationBatch, error) {
	return l.transitionBatch(ctx, id, model.BatchProcessing, expectedRevision)
}

func (l *Lab) CompleteBatch(ctx context.Context, id string, expectedRevision int) (model.PreparationBatch, error) {
	return l.transitionBatch(ctx, id, model.BatchReady, expectedRevision)
}

func (l *Lab) BlockBatch(ctx context.Context, id string, expectedRevision int, reason string) (model.PreparationBatch, error) {
	if reason == "" {
		return model.PreparationBatch{}, model.ErrInvalidInput
	}
	batch, err := l.transitionBatch(ctx, id, model.BatchBlocked, expectedRevision)
	if err != nil {
		return model.PreparationBatch{}, err
	}
	batch.Notes = reason
	batch.UpdatedAt = l.clock.Now()
	if err := l.store.SaveBatch(batch); err != nil {
		return model.PreparationBatch{}, err
	}
	return batch, nil
}

func (l *Lab) ResumeBatch(ctx context.Context, id string, expectedRevision int) (model.PreparationBatch, error) {
	if err := checkContext(ctx); err != nil {
		return model.PreparationBatch{}, err
	}
	blocked, err := l.HasBlockingQuality(ctx, id)
	if err != nil {
		return model.PreparationBatch{}, err
	}
	if blocked {
		return model.PreparationBatch{}, fmt.Errorf("batch %s has unresolved quality flags: %w", id, model.ErrInvalidState)
	}
	return l.transitionBatch(ctx, id, model.BatchProcessing, expectedRevision)
}

func (l *Lab) ArchiveBatch(ctx context.Context, id string, expectedRevision int) (model.PreparationBatch, error) {
	return l.transitionBatch(ctx, id, model.BatchArchived, expectedRevision)
}
