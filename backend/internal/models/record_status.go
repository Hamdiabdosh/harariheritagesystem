package models

type RecordStatus string

const (
	StatusDraft         RecordStatus = "draft"
	StatusPendingReview RecordStatus = "pending_review"
	StatusUnderReview   RecordStatus = "under_review"
	StatusReturned      RecordStatus = "returned"
	StatusApproved      RecordStatus = "approved"
)

func (s RecordStatus) IsValid() bool {
	switch s {
	case StatusDraft, StatusPendingReview, StatusUnderReview, StatusReturned, StatusApproved:
		return true
	default:
		return false
	}
}

func (s RecordStatus) IsEditable() bool {
	return s == StatusDraft || s == StatusReturned
}

type ImmovableOwnerType string

type AgeMethod string

type OverallCondition string

type DamageLevel string

type QualityLevel string

type AccessibilityLevel string

type SexType string
