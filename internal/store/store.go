package store

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	_ "modernc.org/sqlite"
)

type Store struct {
	db *sql.DB
}

func Open(path string) (*Store, error) {
	if path == "" {
		return nil, fmt.Errorf("database path is empty")
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return nil, fmt.Errorf("create database directory: %w", err)
	}
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, fmt.Errorf("open sqlite: %w", err)
	}
	db.SetMaxOpenConns(1)
	db.SetMaxIdleConns(1)
	s := &Store{db: db}
	if err := s.schema(); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("create schema: %w", err)
	}
	return s, nil
}

func (s *Store) schema() error {
	const schema = `
CREATE TABLE IF NOT EXISTS records (
    kind TEXT NOT NULL,
    id TEXT NOT NULL,
    payload BLOB NOT NULL,
    updated_at TEXT NOT NULL,
    PRIMARY KEY(kind, id)
);
CREATE INDEX IF NOT EXISTS records_kind_updated ON records(kind, updated_at);
CREATE TABLE IF NOT EXISTS events (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    subject TEXT NOT NULL,
    action TEXT NOT NULL,
    created_at TEXT NOT NULL
);
CREATE INDEX IF NOT EXISTS events_subject_created ON events(subject, created_at);
`
	_, err := s.db.Exec(schema)
	return err
}

func (s *Store) Close() error {
	if s == nil || s.db == nil {
		return nil
	}
	return s.db.Close()
}

func (s *Store) Save(kind, id string, value any) error {
	raw, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("marshal %s/%s: %w", kind, id, err)
	}
	_, err = s.db.Exec(`
INSERT INTO records(kind, id, payload, updated_at) VALUES(?, ?, ?, ?)
ON CONFLICT(kind, id) DO UPDATE SET payload=excluded.payload, updated_at=excluded.updated_at`,
		kind, id, raw, time.Now().UTC().Format(time.RFC3339Nano))
	if err != nil {
		return fmt.Errorf("save %s/%s: %w", kind, id, err)
	}
	return nil
}

func (s *Store) Insert(kind, id string, value any) error {
	raw, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("marshal %s/%s: %w", kind, id, err)
	}
	_, err = s.db.Exec(`INSERT INTO records(kind, id, payload, updated_at) VALUES(?, ?, ?, ?)`,
		kind, id, raw, time.Now().UTC().Format(time.RFC3339Nano))
	return err
}

func (s *Store) Load(kind, id string, value any) error {
	var raw []byte
	err := s.db.QueryRow(`SELECT payload FROM records WHERE kind=? AND id=?`, kind, id).Scan(&raw)
	if err != nil {
		return err
	}
	return json.Unmarshal(raw, value)
}

func (s *Store) Delete(kind, id string) error {
	_, err := s.db.Exec(`DELETE FROM records WHERE kind=? AND id=?`, kind, id)
	return err
}

func (s *Store) List(kind string, into func([]byte) error) error {
	rows, err := s.db.Query(`SELECT payload FROM records WHERE kind=? ORDER BY updated_at, id`, kind)
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		var raw []byte
		if err := rows.Scan(&raw); err != nil {
			return err
		}
		if err := into(raw); err != nil {
			return err
		}
	}
	return rows.Err()
}

func (s *Store) Event(subject, action string) error {
	_, err := s.db.Exec(`INSERT INTO events(subject, action, created_at) VALUES(?, ?, ?)`,
		subject, action, time.Now().UTC().Format(time.RFC3339Nano))
	return err
}

func (s *Store) Count(kind string) (int, error) {
	var count int
	if err := s.db.QueryRow(`SELECT COUNT(*) FROM records WHERE kind=?`, kind).Scan(&count); err != nil {
		return 0, fmt.Errorf("count %s: %w", kind, err)
	}
	return count, nil
}

func (s *Store) Purge(kind string, before time.Time) error {
	_, err := s.db.Exec(`DELETE FROM records WHERE kind=? AND updated_at < ?`, kind, before.UTC().Format(time.RFC3339Nano))
	return err
}

func (s *Store) Transaction(fn func(*sql.Tx) error) error {
	tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}
	if err := fn(tx); err != nil {
		_ = rollback(tx)
		return fmt.Errorf("transaction: %w", err)
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}
	return nil
}

func rollback(tx *sql.Tx) error {
	if tx == nil {
		return nil
	}
	return tx.Rollback()
}

func saveTx(tx *sql.Tx, kind, id string, value any, at time.Time) error {
	raw, err := json.Marshal(value)
	if err != nil {
		return err
	}
	_, err = tx.Exec(`
INSERT INTO records(kind, id, payload, updated_at) VALUES(?, ?, ?, ?)
ON CONFLICT(kind, id) DO UPDATE SET payload=excluded.payload, updated_at=excluded.updated_at`,
		kind, id, raw, at.UTC().Format(time.RFC3339Nano))
	return err
}

func decodeList[T any](s *Store, kind string) ([]T, error) {
	items := make([]T, 0)
	err := s.List(kind, func(raw []byte) error {
		var item T
		if err := json.Unmarshal(raw, &item); err != nil {
			return err
		}
		items = append(items, item)
		return nil
	})
	return items, err
}
