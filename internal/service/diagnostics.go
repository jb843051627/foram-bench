package service

import (
	"context"
	"fmt"

	"github.com/jb843051627/foram-bench/internal/model"
	"github.com/jb843051627/foram-bench/internal/rules"
)

func (l *Lab) DiagnoseBatch(ctx context.Context, batchID string) ([]model.Diagnostic, error) {
	if err := checkContext(ctx); err != nil {
		return nil, err
	}
	sections, err := l.store.ListSectionsByBatch(batchID)
	if err != nil {
		return nil, err
	}
	result := make([]model.Diagnostic, 0)
	for _, section := range sections {
		if err := checkContext(ctx); err != nil {
			return nil, err
		}
		items, listErr := l.store.ListObservationsBySection(section.ID)
		if listErr != nil {
			return nil, listErr
		}
		for _, diagnostic := range rules.EvaluateObservations(items) {
			diagnostic.BatchID = batchID
			if saveErr := l.store.SaveDiagnostic(diagnostic); saveErr != nil {
				return nil, saveErr
			}
			result = append(result, diagnostic)
		}
	}
	if len(result) > 0 {
		if err := l.store.Event(batchID, "batch.diagnosed"); err != nil {
			return nil, err
		}
	}
	return result, nil
}

func (l *Lab) ListDiagnostics(ctx context.Context, batchID string) ([]model.Diagnostic, error) {
	if err := checkContext(ctx); err != nil {
		return nil, err
	}
	return l.store.ListDiagnostics(batchID)
}

func (l *Lab) AcknowledgeDiagnostic(ctx context.Context, id string) (model.Diagnostic, error) {
	if err := checkContext(ctx); err != nil {
		return model.Diagnostic{}, err
	}
	diagnostic, err := l.store.GetDiagnostic(id)
	if err != nil {
		return model.Diagnostic{}, err
	}
	diagnostic.Acknowledged = true
	if err := l.store.SaveDiagnostic(diagnostic); err != nil {
		return model.Diagnostic{}, err
	}
	return diagnostic, nil
}

func (l *Lab) ReleaseDiagnostics(ctx context.Context, batchID string) error {
	items, err := l.ListDiagnostics(ctx, batchID)
	if err != nil {
		return err
	}
	if rules.HasBlockingDiagnostics(items) {
		return fmt.Errorf("batch %s has blocking diagnostics: %w", batchID, model.ErrInvalidState)
	}
	return nil
}
