package service

import (
	"context"
	"fmt"
	"time"

	"github.com/jb843051627/foram-bench/internal/model"
)

func (l *Lab) ListEvents(ctx context.Context, subject string, limit int) ([]model.AuditEvent, error) {
	if err := checkContext(ctx); err != nil {
		return nil, err
	}
	if subject == "" {
		return nil, fmt.Errorf("event subject is empty: %w", model.ErrInvalidInput)
	}
	return l.store.ListEvents(subject, limit)
}

func (l *Lab) EventsSince(ctx context.Context, subject string, since time.Time) ([]model.AuditEvent, error) {
	items, err := l.ListEvents(ctx, subject, 500)
	if err != nil {
		return nil, err
	}
	result := make([]model.AuditEvent, 0, len(items))
	for _, item := range items {
		if item.CreatedAt.After(since) {
			result = append(result, item)
		}
	}
	return result, nil
}

func (l *Lab) RecordNoteEvent(ctx context.Context, subject, action string) error {
	if err := checkContext(ctx); err != nil {
		return err
	}
	if subject == "" || action == "" {
		return model.ErrInvalidInput
	}
	return l.store.Event(subject, action)
}
