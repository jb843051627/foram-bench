package store

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/jb843051627/foram-bench/internal/model"
)

const diagnosticKind = "diagnostic"

func (s *Store) SaveDiagnostic(diagnostic model.Diagnostic) error {
	if err := diagnostic.Validate(); err != nil {
		return err
	}
	return s.Save(diagnosticKind, diagnostic.ID, diagnostic)
}

func (s *Store) GetDiagnostic(id string) (model.Diagnostic, error) {
	var diagnostic model.Diagnostic
	err := s.Load(diagnosticKind, id, &diagnostic)
	if errors.Is(err, sql.ErrNoRows) {
		return model.Diagnostic{}, fmt.Errorf("diagnostic %s: %w", id, model.ErrNotFound)
	}
	if err != nil {
		return model.Diagnostic{}, err
	}
	return diagnostic, nil
}

func (s *Store) ListDiagnostics(batchID string) ([]model.Diagnostic, error) {
	all, err := decodeList[model.Diagnostic](s, diagnosticKind)
	if err != nil {
		return nil, err
	}
	result := make([]model.Diagnostic, 0)
	for _, diagnostic := range all {
		if diagnostic.BatchID == batchID {
			result = append(result, diagnostic)
		}
	}
	return result, nil
}
