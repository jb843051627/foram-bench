package rules

import (
	"fmt"

	"github.com/jb843051627/foram-bench/internal/model"
)

type Transition struct {
	From string
	To   string
	Why  string
}

func BatchTransitions() []Transition {
	return []Transition{
		{From: model.BatchPlanned, To: model.BatchProcessing, Why: "开始制备"},
		{From: model.BatchProcessing, To: model.BatchReady, Why: "完成制备"},
		{From: model.BatchProcessing, To: model.BatchBlocked, Why: "质量问题挂起"},
		{From: model.BatchBlocked, To: model.BatchProcessing, Why: "完成补救后继续"},
		{From: model.BatchReady, To: model.BatchArchived, Why: "完成归档"},
		{From: model.BatchBlocked, To: model.BatchArchived, Why: "放弃当前批次"},
	}
}

func ExplainBatchTransition(from, to string) (string, error) {
	for _, transition := range BatchTransitions() {
		if transition.From == from && transition.To == to {
			return transition.Why, nil
		}
	}
	return "", fmt.Errorf("unsupported batch transition %s -> %s: %w", from, to, model.ErrInvalidState)
}

func NextBatchStates(current string) []string {
	result := make([]string, 0)
	for _, transition := range BatchTransitions() {
		if transition.From == current {
			result = append(result, transition.To)
		}
	}
	return result
}

func IsTerminalBatch(status string) bool {
	return status == model.BatchArchived
}
