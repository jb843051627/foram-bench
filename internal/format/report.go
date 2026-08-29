package format

import (
	"fmt"
	"strings"
	"time"

	"github.com/jb843051627/foram-bench/internal/analysis"
	"github.com/jb843051627/foram-bench/internal/model"
)

type ReportDocument struct {
	Title      string
	Subtitle   string
	Paragraphs []string
	Findings   []analysis.Finding
}

func NewReportDocument(sample model.Sample, batch model.PreparationBatch, score analysis.QualitySummary) ReportDocument {
	collectionDate := sample.CollectionDate
	if location, err := time.LoadLocation(sample.TimeZone); err == nil {
		collectionDate = collectionDate.In(location)
	}
	return ReportDocument{
		Title:    fmt.Sprintf("微体化石样品 %s", sample.ID),
		Subtitle: fmt.Sprintf("制备批次 %s · %s · 质量等级 %s", batch.ID, sample.Location, score.Grade),
		Paragraphs: []string{fmt.Sprintf("样品采集于 %s，材料为 %s。", collectionDate.Format("2006-01-02"), sample.Material),
			fmt.Sprintf("观察覆盖率 %s，完整性 %s，保存状况 %s。", Percent(score.Coverage), Percent(score.Completeness), Percent(score.Preservation))},
	}
}

func (d ReportDocument) Markdown() string {
	lines := []string{"# " + d.Title, "", d.Subtitle, ""}
	for _, paragraph := range d.Paragraphs {
		lines = append(lines, paragraph, "")
	}
	if len(d.Findings) > 0 {
		lines = append(lines, "## 质量提示", "")
		for _, finding := range d.Findings {
			lines = append(lines, fmt.Sprintf("- [%s] %s", strings.ToUpper(finding.Severity), finding.Message))
		}
		lines = append(lines, "")
	}
	return strings.TrimSpace(strings.Join(lines, "\n"))
}

func AddSection(d ReportDocument, section model.ThinSection, observations []model.Observation) ReportDocument {
	items := analysis.Abundances(observations)
	d.Paragraphs = append(d.Paragraphs, fmt.Sprintf("薄片 %s（%s，%d μm）记录：%s。", section.Label, section.Stain, section.ThicknessUM, TaxonList(items)))
	return d
}
