package model

const (
	SampleRegistered = "registered"
	SamplePrepared   = "prepared"
	SampleArchived   = "archived"

	BatchPlanned    = "planned"
	BatchProcessing = "processing"
	BatchReady      = "ready"
	BatchBlocked    = "blocked"
	BatchArchived   = "archived"

	SectionCut      = "cut"
	SectionStained  = "stained"
	SectionReviewed = "reviewed"

	ReviewAccepted = "accepted"
	ReviewRevise   = "revise"
	ReviewRejected = "rejected"
)

func ValidSampleStatus(v string) bool {
	return v == SampleRegistered || v == SamplePrepared || v == SampleArchived
}

func ValidBatchStatus(v string) bool {
	return v == BatchPlanned || v == BatchProcessing || v == BatchReady || v == BatchBlocked || v == BatchArchived
}

func ValidReviewDecision(v string) bool {
	return v == ReviewAccepted || v == ReviewRevise || v == ReviewRejected
}

func CanMoveBatch(from, to string) bool {
	switch from {
	case BatchPlanned:
		return to == BatchProcessing || to == BatchBlocked
	case BatchProcessing:
		return to == BatchReady || to == BatchBlocked
	case BatchReady:
		return to == BatchArchived || to == BatchBlocked
	case BatchBlocked:
		return to == BatchProcessing || to == BatchArchived
	default:
		return false
	}
}

func CanMoveSection(from, to string) bool {
	switch from {
	case SectionCut:
		return to == SectionStained
	case SectionStained:
		return to == SectionReviewed
	default:
		return false
	}
}
