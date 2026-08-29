package store

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/jb843051627/foram-bench/internal/model"
)

const layerKind = "stratigraphic_layer"

func (s *Store) SaveLayer(layer model.StratigraphicLayer) error {
	if err := layer.Validate(); err != nil {
		return err
	}
	return s.Save(layerKind, layer.ID, layer)
}

func (s *Store) GetLayer(id string) (model.StratigraphicLayer, error) {
	var layer model.StratigraphicLayer
	err := s.Load(layerKind, id, &layer)
	if errors.Is(err, sql.ErrNoRows) {
		return model.StratigraphicLayer{}, fmt.Errorf("layer %s: %w", id, model.ErrNotFound)
	}
	if err != nil {
		return model.StratigraphicLayer{}, err
	}
	return layer, nil
}

func (s *Store) ListLayers(siteCode string) ([]model.StratigraphicLayer, error) {
	all, err := decodeList[model.StratigraphicLayer](s, layerKind)
	if err != nil {
		return nil, err
	}
	result := make([]model.StratigraphicLayer, 0)
	for _, layer := range all {
		if siteCode == "" || layer.SiteCode == siteCode {
			result = append(result, layer)
		}
	}
	return result, nil
}

func (s *Store) SaveLayerMarker(marker model.LayerMarker) error {
	if err := marker.Validate(); err != nil {
		return err
	}
	return s.Save("layer_marker", marker.LayerID+"-"+marker.Marker, marker)
}
