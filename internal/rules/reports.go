package rules

import (
	"strings"

	"github.com/jb843051627/foram-bench/internal/model"
)

func RequiredReportSections() []string {
	return []string{"采集信息", "制备批次", "薄片观察", "质量评估", "复核结论", "归档信息"}
}

func MissingReportSections(body string) []string {
	missing := make([]string, 0)
	for _, section := range RequiredReportSections() {
		if !strings.Contains(body, section) {
			missing = append(missing, section)
		}
	}
	return missing
}

func ReportReady(report model.Report) bool {
	return report.Status == "published" && len(MissingReportSections(report.Body)) == 0
}

func NormalizeReportStatus(status string) string {
	return strings.ToLower(strings.TrimSpace(status))
}
