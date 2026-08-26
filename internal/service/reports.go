package service

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/jb843051627/foram-bench/internal/analysis"
	"github.com/jb843051627/foram-bench/internal/format"
	"github.com/jb843051627/foram-bench/internal/model"
)

func (l *Lab) GenerateReport(ctx context.Context, batchID string) (model.Report, error) {
	l.reportMu.Lock()
	defer l.reportMu.Unlock()
	if err := checkContext(ctx); err != nil {
		return model.Report{}, err
	}
	batch, err := l.store.GetBatch(batchID)
	if err != nil {
		return model.Report{}, err
	}
	if batch.Status != model.BatchReady {
		return model.Report{}, fmt.Errorf("batch %s is not ready: %w", batchID, model.ErrInvalidState)
	}
	if blocked, err := l.HasBlockingQuality(ctx, batchID); err != nil {
		return model.Report{}, err
	} else if blocked {
		return model.Report{}, fmt.Errorf("batch %s has blocking quality flags: %w", batchID, model.ErrInvalidState)
	}
	sample, err := l.store.GetSample(batch.SampleID)
	if err != nil {
		return model.Report{}, err
	}
	sections, err := l.store.ListSectionsByBatch(batchID)
	if err != nil {
		return model.Report{}, err
	}
	if len(sections) == 0 {
		return model.Report{}, fmt.Errorf("batch %s has no thin sections: %w", batchID, model.ErrInvalidState)
	}
	allObservations := make([]model.Observation, 0)
	for _, section := range sections {
		observations, listErr := l.store.ListObservationsBySection(section.ID)
		if listErr != nil {
			return model.Report{}, listErr
		}
		allObservations = append(allObservations, observations...)
	}
	quality := analysis.Assess(sections, allObservations, []string{"Globigerina", "Ammonia", "Elphidium", "Cibicidoides"})
	document := format.NewReportDocument(sample, batch, quality)
	lines := make([]string, 0, len(sections)+5)
	lines = append(lines, document.Title, document.Subtitle)
	for _, section := range sections {
		if err := checkContext(ctx); err != nil {
			return model.Report{}, fmt.Errorf("load accepted review for section %s: %w", section.ID, err)
		}
		accepted, err := l.HasAcceptedReview(ctx, section.ID)
		if err != nil {
			return model.Report{}, err
		}
		if !accepted {
			return model.Report{}, fmt.Errorf("section %s lacks accepted review: %w", section.ID, model.ErrInvalidState)
		}
		observations, err := l.store.ListObservationsBySection(section.ID)
		if err != nil {
			return model.Report{}, err
		}
		taxa := make([]string, 0, len(observations))
		for _, observation := range observations {
			taxa = append(taxa, fmt.Sprintf("%s=%d", observation.Taxon, observation.Count))
		}
		sort.Strings(taxa)
		document = format.AddSection(document, section, observations)
		lines = append(lines, fmt.Sprintf("Section %s [%s] %s", section.Label, section.Stain, strings.Join(taxa, ", ")))
	}
	lines = append(lines, fmt.Sprintf("Quality grade %s score=%s", quality.Grade, format.Percent(quality.Score)))
	location, err := time.LoadLocation(sample.TimeZone)
	if err != nil {
		return model.Report{}, fmt.Errorf("load sample timezone: %w", err)
	}
	generated := l.clock.Now().In(location)
	lines = append(lines, "Generated "+generated.Format("2006-01-02 15:04:05 MST"))
	reports, err := l.store.ListReportsBySample(sample.ID)
	if err != nil {
		return model.Report{}, err
	}
	report := model.Report{ID: fmt.Sprintf("%s-r%02d", sample.ID, len(reports)+1), SampleID: sample.ID,
		Version: len(reports) + 1, Status: "published", TimeZone: sample.TimeZone,
		Body: strings.Join(lines, "\n"), GeneratedAt: generated}
	if err := report.Validate(); err != nil {
		return model.Report{}, err
	}
	if err := l.store.InsertReport(report); err != nil {
		return model.Report{}, err
	}
	if err := l.store.Event(batchID, "report.generated"); err != nil {
		return model.Report{}, err
	}
	l.metrics.Add("reports.generated", 1)
	return report, nil
}

func (l *Lab) GetReport(ctx context.Context, id string) (model.Report, error) {
	if err := checkContext(ctx); err != nil {
		return model.Report{}, err
	}
	return l.store.GetReport(id)
}

func (l *Lab) ListReports(ctx context.Context, sampleID string) ([]model.Report, error) {
	if err := checkContext(ctx); err != nil {
		return nil, err
	}
	return l.store.ListReportsBySample(sampleID)
}

func (l *Lab) ExportReport(ctx context.Context, id string) (string, error) {
	report, err := l.GetReport(ctx, id)
	if err != nil {
		return "", err
	}
	if strings.TrimSpace(report.Body) == "" {
		return "", errors.New("report body is empty")
	}
	return report.Body, nil
}
