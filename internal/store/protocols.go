package store

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/jb843051627/foram-bench/internal/model"
)

const protocolKind = "preparation_protocol"

func (s *Store) SaveProtocol(protocol model.PreparationProtocol) error {
	if err := protocol.Validate(); err != nil {
		return err
	}
	return s.Save(protocolKind, protocol.ID, protocol)
}

func (s *Store) GetProtocol(id string) (model.PreparationProtocol, error) {
	var protocol model.PreparationProtocol
	err := s.Load(protocolKind, id, &protocol)
	if errors.Is(err, sql.ErrNoRows) {
		return model.PreparationProtocol{}, fmt.Errorf("protocol %s: %w", id, model.ErrNotFound)
	}
	if err != nil {
		return model.PreparationProtocol{}, fmt.Errorf("load protocol %s: %w", id, err)
	}
	return protocol, nil
}

func (s *Store) ListProtocols() ([]model.PreparationProtocol, error) {
	items, err := decodeList[model.PreparationProtocol](s, protocolKind)
	if err != nil {
		return nil, fmt.Errorf("list protocols: %w", err)
	}
	return items, nil
}

func (s *Store) LatestProtocol(name string) (model.PreparationProtocol, error) {
	protocols, err := s.ListProtocols()
	if err != nil {
		return model.PreparationProtocol{}, err
	}
	var latest model.PreparationProtocol
	for _, protocol := range protocols {
		if protocol.Name == name && (latest.ID == "" || protocol.Version > latest.Version) {
			latest = protocol
		}
	}
	if latest.ID == "" {
		return model.PreparationProtocol{}, fmt.Errorf("protocol %s: %w", name, model.ErrNotFound)
	}
	return latest, nil
}
