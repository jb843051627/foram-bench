package store

import (
	"fmt"
	"time"

	"github.com/jb843051627/foram-bench/internal/model"
)

func (s *Store) ListEvents(subject string, limit int) ([]model.AuditEvent, error) {
	if limit <= 0 || limit > 500 {
		limit = 100
	}
	query := `SELECT id, subject, action, created_at FROM events WHERE subject=? ORDER BY id DESC LIMIT ?`
	rows, err := s.db.Query(query, subject, limit)
	if err != nil {
		return nil, fmt.Errorf("list events: %w", err)
	}
	defer rows.Close()
	events := make([]model.AuditEvent, 0)
	for rows.Next() {
		var event model.AuditEvent
		var created string
		if err := rows.Scan(&event.ID, &event.Subject, &event.Action, &created); err != nil {
			return nil, err
		}
		parsed, err := timeParse(created)
		if err != nil {
			return nil, err
		}
		event.CreatedAt = parsed
		events = append(events, event)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return events, nil
}

func timeParse(value string) (time.Time, error) {
	parsed, err := time.Parse(time.RFC3339Nano, value)
	if err != nil {
		return time.Time{}, fmt.Errorf("parse event time: %w", err)
	}
	return parsed, nil
}
