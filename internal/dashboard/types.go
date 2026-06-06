package dashboard

import (
	"time"

	"github.com/google/uuid"

	"github.com/qirs-mezgeb/api/internal/models"
)

type ListFilters struct {
	Page     int
	Limit    int
	Type     models.RecordType
	Status   models.RecordStatus
	Woreda   string
	Kebele   string
	Search   string
	DateFrom *time.Time
	DateTo   *time.Time
}

type StatusCounts struct {
	Draft         int `json:"draft"`
	PendingReview int `json:"pending_review"`
	UnderReview   int `json:"under_review"`
	Returned      int `json:"returned"`
	Approved      int `json:"approved"`
}

type Stats struct {
	TotalImmovable  int          `json:"total_immovable"`
	TotalMovable    int          `json:"total_movable"`
	ByStatus        StatusCounts `json:"by_status"`
	PendingMyAction *int         `json:"pending_my_action,omitempty"`
}

type RecordSummary struct {
	ID          uuid.UUID           `json:"id"`
	RecordType  models.RecordType   `json:"record_type"`
	RecordID    string              `json:"record_id"`
	NameAmharic string              `json:"name_amharic"`
	Status      models.RecordStatus `json:"status"`
	Woreda      *string             `json:"woreda,omitempty"`
	Kebele      *string             `json:"kebele,omitempty"`
	RegistrarID uuid.UUID           `json:"registrar_id"`
	CreatedAt   time.Time           `json:"created_at"`
	UpdatedAt   time.Time           `json:"updated_at"`
}

type PaginatedRecords struct {
	Items      []RecordSummary `json:"items"`
	Total      int             `json:"total"`
	Page       int             `json:"page"`
	Limit      int             `json:"limit"`
	TotalPages int             `json:"total_pages"`
}
