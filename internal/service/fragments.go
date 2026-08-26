package service

import (
	"context"
	"fmt"

	"github.com/jb843051627/foram-bench/internal/model"
)

type FragmentInput struct {
	ID         string  `json:"id"`
	SectionID  string  `json:"section_id"`
	PositionMM float64 `json:"position_mm"`
	LengthMM   float64 `json:"length_mm"`
	Color      string  `json:"color"`
	Texture    string  `json:"texture"`
	Damaged    bool    `json:"damaged"`
	Note       string  `json:"note"`
}

func (l *Lab) RecordFragment(ctx context.Context, input FragmentInput) (model.Fragment, error) {
	if err := checkContext(ctx); err != nil {
		return model.Fragment{}, err
	}
	if _, err := l.store.GetSection(input.SectionID); err != nil {
		return model.Fragment{}, fmt.Errorf("load section for fragment: %w", err)
	}
	fragment := model.Fragment{ID: input.ID, SectionID: input.SectionID, PositionMM: input.PositionMM, LengthMM: input.LengthMM,
		Color: input.Color, Texture: input.Texture, Damaged: input.Damaged, Note: input.Note, CreatedAt: l.clock.Now()}
	if err := fragment.Validate(); err != nil {
		return model.Fragment{}, err
	}
	if err := l.store.SaveFragment(fragment); err != nil {
		return model.Fragment{}, err
	}
	if err := l.store.Event(input.SectionID, "fragment.recorded"); err != nil {
		return model.Fragment{}, err
	}
	return fragment, nil
}

func (l *Lab) GetFragment(ctx context.Context, id string) (model.Fragment, error) {
	if err := checkContext(ctx); err != nil {
		return model.Fragment{}, err
	}
	return l.store.GetFragment(id)
}

func (l *Lab) ListFragments(ctx context.Context, sectionID string) ([]model.Fragment, error) {
	if err := checkContext(ctx); err != nil {
		return nil, err
	}
	return l.store.ListFragmentsBySection(sectionID)
}

func (l *Lab) FragmentCoverage(ctx context.Context, sectionID string) (float64, error) {
	fragments, err := l.ListFragments(ctx, sectionID)
	if err != nil {
		return 0, err
	}
	set := model.FragmentSet{SectionID: sectionID, Items: fragments}
	return set.Coverage(), nil
}
