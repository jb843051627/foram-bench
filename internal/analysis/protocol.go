package analysis

import (
	"fmt"
	"strings"

	"github.com/jb843051627/foram-bench/internal/model"
)

type ProtocolCheck struct {
	Passed    bool     `json:"passed"`
	Findings  []string `json:"findings"`
	Completed int      `json:"completed"`
	Required  int      `json:"required"`
}

func CheckProtocol(protocol model.PreparationProtocol, completed []int) ProtocolCheck {
	result := ProtocolCheck{Passed: true, Findings: make([]string, 0), Required: 0}
	seen := make(map[int]bool, len(completed))
	for _, number := range completed {
		seen[number] = true
	}
	for _, step := range protocol.Steps {
		if !step.Required {
			continue
		}
		result.Required++
		if seen[step.Number] {
			result.Completed++
			continue
		}
		result.Passed = false
		result.Findings = append(result.Findings, fmt.Sprintf("缺少步骤 %d: %s", step.Number, step.Name))
	}
	if result.Required == 0 {
		result.Passed = false
		result.Findings = append(result.Findings, "协议没有必需步骤")
	}
	return result
}

func NormalizeProtocolName(value string) string {
	return strings.Join(strings.Fields(strings.ToLower(value)), " ")
}
