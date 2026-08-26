package format

import (
	"sort"
	"strings"

	"github.com/jb843051627/foram-bench/internal/analysis"
)

func TaxonList(items []analysis.Abundance) string {
	copyItems := append([]analysis.Abundance(nil), items...)
	sort.SliceStable(copyItems, func(i, j int) bool { return copyItems[i].Rank < copyItems[j].Rank })
	parts := make([]string, 0, len(copyItems))
	for _, item := range copyItems {
		parts = append(parts, item.Taxon+" "+Percent(item.Relative))
	}
	return strings.Join(parts, ", ")
}

func TaxonLabel(value string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return "未命名类群"
	}
	return value
}

func Preservation(value string) string {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "excellent", "good":
		return "良好"
	case "fair", "moderate":
		return "一般"
	case "poor", "fragmented":
		return "较差"
	default:
		return "未记录"
	}
}
