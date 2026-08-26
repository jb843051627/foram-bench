package service

import (
	"context"
	"fmt"
	"time"

	"github.com/jb843051627/foram-bench/internal/model"
	"github.com/jb843051627/foram-bench/internal/validation"
)

type SiteInput struct {
	Code         string  `json:"code"`
	Name         string  `json:"name"`
	Region       string  `json:"region"`
	Latitude     float64 `json:"latitude"`
	Longitude    float64 `json:"longitude"`
	ElevationM   int     `json:"elevation_m"`
	TimeZone     string  `json:"time_zone"`
	Stratigraphy string  `json:"stratigraphy"`
}

func (l *Lab) RegisterSite(ctx context.Context, input SiteInput) (model.CollectionSite, error) {
	if err := checkContext(ctx); err != nil {
		return model.CollectionSite{}, err
	}
	if err := validation.SiteCode(input.Code); err != nil {
		return model.CollectionSite{}, err
	}
	now := l.clock.Now()
	site := model.CollectionSite{Code: input.Code, Name: input.Name, Region: input.Region, Latitude: input.Latitude,
		Longitude: input.Longitude, ElevationM: input.ElevationM, TimeZone: input.TimeZone,
		Stratigraphy: input.Stratigraphy, Active: true, CreatedAt: now}
	if err := site.Validate(); err != nil {
		return model.CollectionSite{}, fmt.Errorf("validate site: %w", err)
	}
	if err := l.store.SaveSite(site); err != nil {
		return model.CollectionSite{}, err
	}
	if err := l.store.Event(site.Code, "site.registered"); err != nil {
		return model.CollectionSite{}, err
	}
	return site, nil
}

func (l *Lab) GetSite(ctx context.Context, code string) (model.CollectionSite, error) {
	if err := checkContext(ctx); err != nil {
		return model.CollectionSite{}, err
	}
	return l.store.GetSite(code)
}

func (l *Lab) ListSites(ctx context.Context, activeOnly bool) ([]model.CollectionSite, error) {
	if err := checkContext(ctx); err != nil {
		return nil, err
	}
	if activeOnly {
		return l.store.ListActiveSites()
	}
	return l.store.ListSites()
}

func (l *Lab) RetireSite(ctx context.Context, code string) (model.CollectionSite, error) {
	if err := checkContext(ctx); err != nil {
		return model.CollectionSite{}, err
	}
	site, err := l.store.GetSite(code)
	if err != nil {
		return model.CollectionSite{}, err
	}
	if !site.Active {
		return site, nil
	}
	site.Active = false
	if err := l.store.SaveSite(site); err != nil {
		return model.CollectionSite{}, err
	}
	if err := l.store.Event(code, "site.retired"); err != nil {
		return model.CollectionSite{}, err
	}
	return site, nil
}

func (l *Lab) SiteAge(ctx context.Context, code string, now time.Time) (time.Duration, error) {
	site, err := l.GetSite(ctx, code)
	if err != nil {
		return 0, err
	}
	return now.Sub(site.CreatedAt), nil
}
