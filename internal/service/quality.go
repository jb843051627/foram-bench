package service

import (
	"context"
	"fmt"

	"github.com/jb843051627/foram-bench/internal/ingest"
	"github.com/jb843051627/foram-bench/internal/model"
)

func (l *Lab) AddQualityFlag(ctx context.Context, input QualityInput) (model.QualityFlag, error) {
	l.qualityMu.Lock()
	defer l.qualityMu.Unlock()
	if err := checkContext(ctx); err != nil {
		return model.QualityFlag{}, err
	}
	if _, err := l.store.GetBatch(input.BatchID); err != nil {
		return model.QualityFlag{}, fmt.Errorf("load batch for quality flag: %w", err)
	}
	flag := model.QualityFlag{ID: input.ID, BatchID: input.BatchID, Kind: input.Kind,
		Severity: input.Severity, Detail: input.Detail, CreatedAt: l.clock.Now()}
	if err := flag.Validate(); err != nil {
		return model.QualityFlag{}, err
	}
	if err := l.store.SaveQualityFlag(flag); err != nil {
		return model.QualityFlag{}, err
	}
	if err := l.store.Event(input.BatchID, "quality.flagged"); err != nil {
		return model.QualityFlag{}, err
	}
	l.metrics.Add("quality.flags", 1)
	return flag, nil
}

func (l *Lab) GetQualityFlag(ctx context.Context, id string) (model.QualityFlag, error) {
	if err := checkContext(ctx); err != nil {
		return model.QualityFlag{}, err
	}
	return l.store.GetQualityFlag(id)
}

func (l *Lab) ResolveQualityFlag(ctx context.Context, id string) (model.QualityFlag, error) {
	l.qualityMu.Lock()
	defer l.qualityMu.Unlock()
	if err := checkContext(ctx); err != nil {
		return model.QualityFlag{}, err
	}
	flag, err := l.store.GetQualityFlag(id)
	if err != nil {
		return model.QualityFlag{}, err
	}
	if flag.Resolved {
		return flag, nil
	}
	now := l.clock.Now()
	flag.Resolved = true
	flag.ResolvedAt = &now
	if err := l.store.SaveQualityFlag(flag); err != nil {
		return model.QualityFlag{}, err
	}
	if err := l.store.Event(flag.BatchID, "quality.resolved"); err != nil {
		return model.QualityFlag{}, err
	}
	return flag, nil
}

func (l *Lab) ListQualityFlags(ctx context.Context, batchID string) ([]model.QualityFlag, error) {
	if err := checkContext(ctx); err != nil {
		return nil, err
	}
	return l.store.ListQualityFlagsByBatch(batchID)
}

func (l *Lab) HasBlockingQuality(ctx context.Context, batchID string) (bool, error) {
	l.qualityMu.Lock()
	defer l.qualityMu.Unlock()
	flags, err := l.ListQualityFlags(ctx, batchID)
	if err != nil {
		return false, err
	}
	for _, flag := range flags {
		if !flag.Resolved && (flag.Severity == "critical" || flag.Severity == "major") {
			return true, nil
		}
	}
	return false, nil
}

func (l *Lab) QueueQualityAssessment(ctx context.Context, batchID string, done chan error) error {
	if done == nil {
		done = make(chan error, 1)
	}
	return l.queue.Submit(ctx, ingest.Job{ID: batchID, Ctx: ctx, Done: done, Run: func(jobCtx context.Context) error {
		if err := checkContext(jobCtx); err != nil {
			return err
		}
		blocked, err := l.HasBlockingQuality(jobCtx, batchID)
		if err != nil {
			return err
		}
		if blocked {
			_, err = l.BlockBatch(jobCtx, batchID, 0, "quality assessment found an open blocking flag")
			return err
		}
		return nil
	}})
}
