package service

import (
	"context"
	"fmt"

	"github.com/jb843051627/foram-bench/internal/model"
	"github.com/jb843051627/foram-bench/internal/rules"
)

type LayerInput struct {
	ID            string  `json:"id"`
	SiteCode      string  `json:"site_code"`
	Name          string  `json:"name"`
	TopDepthMM    int     `json:"top_depth_mm"`
	BottomDepthMM int     `json:"bottom_depth_mm"`
	Lithology     string  `json:"lithology"`
	Color         string  `json:"color"`
	Boundary      string  `json:"boundary"`
	Confidence    float64 `json:"confidence"`
}

func (l *Lab) RegisterLayer(ctx context.Context, input LayerInput) (model.StratigraphicLayer, error) {
	if err := checkContext(ctx); err != nil {
		return model.StratigraphicLayer{}, err
	}
	if _, err := l.store.GetSite(input.SiteCode); err != nil {
		return model.StratigraphicLayer{}, fmt.Errorf("load layer site: %w", err)
	}
	layer := model.StratigraphicLayer{ID: input.ID, SiteCode: input.SiteCode, Name: input.Name, TopDepthMM: input.TopDepthMM,
		BottomDepthMM: input.BottomDepthMM, Lithology: input.Lithology, Color: input.Color, Boundary: input.Boundary,
		Confidence: input.Confidence, RecordedAt: l.clock.Now()}
	if err := layer.Validate(); err != nil {
		return model.StratigraphicLayer{}, err
	}
	if err := l.store.SaveLayer(layer); err != nil {
		return model.StratigraphicLayer{}, err
	}
	if err := l.store.Event(input.SiteCode, "layer.registered"); err != nil {
		return model.StratigraphicLayer{}, err
	}
	return layer, nil
}

func (l *Lab) GetLayer(ctx context.Context, id string) (model.StratigraphicLayer, error) {
	if err := checkContext(ctx); err != nil {
		return model.StratigraphicLayer{}, err
	}
	return l.store.GetLayer(id)
}

func (l *Lab) ListLayers(ctx context.Context, siteCode string) ([]model.StratigraphicLayer, error) {
	if err := checkContext(ctx); err != nil {
		return nil, err
	}
	return l.store.ListLayers(siteCode)
}

func (l *Lab) CheckLayerProfile(ctx context.Context, siteCode string) ([]rules.LayerIssue, error) {
	layers, err := l.ListLayers(ctx, siteCode)
	if err != nil {
		return nil, err
	}
	return rules.ValidateLayerProfile(rules.LayerProfile{SiteCode: siteCode, Layers: layers}), nil
}
