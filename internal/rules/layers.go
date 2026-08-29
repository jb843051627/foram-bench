package rules

import (
	"fmt"
	"sort"

	"github.com/jb843051627/foram-bench/internal/model"
)

type LayerIssue struct {
	Code     string `json:"code"`
	Severity string `json:"severity"`
	LayerID  string `json:"layer_id"`
	Message  string `json:"message"`
}

type LayerProfile struct {
	SiteCode string
	Layers   []model.StratigraphicLayer
	Markers  []model.LayerMarker
}

func ValidateLayerProfile(profile LayerProfile) []LayerIssue {
	issues := make([]LayerIssue, 0)
	if profile.SiteCode == "" {
		issues = append(issues, LayerIssue{Code: "SITE_MISSING", Severity: "critical", Message: "层位剖面缺少调查点代码"})
	}
	ordered := make([]model.StratigraphicLayer, len(profile.Layers))
	copy(ordered, profile.Layers)
	sort.Slice(ordered, func(i, j int) bool { return ordered[i].TopDepthMM < ordered[j].TopDepthMM })
	for index, layer := range ordered {
		if err := layer.Validate(); err != nil {
			issues = append(issues, LayerIssue{Code: "LAYER_INVALID", Severity: "major", LayerID: layer.ID, Message: err.Error()})
		}
		if layer.SiteCode != profile.SiteCode {
			issues = append(issues, LayerIssue{Code: "SITE_MISMATCH", Severity: "major", LayerID: layer.ID, Message: "层位属于另一个调查点"})
		}
		if index > 0 && layer.TopDepthMM < ordered[index-1].BottomDepthMM {
			issues = append(issues, LayerIssue{Code: "LAYER_OVERLAP", Severity: "critical", LayerID: layer.ID, Message: fmt.Sprintf("与层位 %s 深度重叠", ordered[index-1].ID)})
		}
	}
	for _, marker := range profile.Markers {
		if err := marker.Validate(); err != nil {
			issues = append(issues, LayerIssue{Code: "MARKER_INVALID", Severity: "major", LayerID: marker.LayerID, Message: err.Error()})
			continue
		}
		found := false
		for _, layer := range profile.Layers {
			if layer.ID == marker.LayerID {
				found = true
				if !layer.ContainsDepth(marker.DepthMM) {
					issues = append(issues, LayerIssue{Code: "MARKER_OUTSIDE", Severity: "major", LayerID: layer.ID, Message: "层位标记深度不在层位范围内"})
				}
				break
			}
		}
		if !found {
			issues = append(issues, LayerIssue{Code: "MARKER_ORPHAN", Severity: "major", LayerID: marker.LayerID, Message: "层位标记引用了不存在的层位"})
		}
	}
	return issues
}

func LayerAtDepth(layers []model.StratigraphicLayer, depthMM int) (model.StratigraphicLayer, bool) {
	for _, layer := range layers {
		if layer.ContainsDepth(depthMM) {
			return layer, true
		}
	}
	return model.StratigraphicLayer{}, false
}

func SortLayers(layers []model.StratigraphicLayer) []model.StratigraphicLayer {
	result := append([]model.StratigraphicLayer(nil), layers...)
	sort.SliceStable(result, func(i, j int) bool {
		if result[i].TopDepthMM == result[j].TopDepthMM {
			return result[i].ID < result[j].ID
		}
		return result[i].TopDepthMM < result[j].TopDepthMM
	})
	return result
}

func TotalThickness(layers []model.StratigraphicLayer) int {
	total := 0
	for _, layer := range layers {
		if layer.ThicknessMM() > 0 {
			total += layer.ThicknessMM()
		}
	}
	return total
}

func Confidence(layers []model.StratigraphicLayer) float64 {
	if len(layers) == 0 {
		return 0
	}
	total := 0.0
	for _, layer := range layers {
		total += layer.Confidence
	}
	return total / float64(len(layers))
}
