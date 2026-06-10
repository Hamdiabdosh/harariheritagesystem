package movable

import (
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/qirs-mezgeb/api/internal/models"
)

type ListFilters struct {
	Page     int
	Limit    int
	Status   models.RecordStatus
	Woreda   string
	Search   string
	DateFrom *time.Time
	DateTo   *time.Time
}

type PaginatedRecords struct {
	Items      []models.MovableRecord `json:"items"`
	Total      int                    `json:"total"`
	Page       int                    `json:"page"`
	Limit      int                    `json:"limit"`
	TotalPages int                    `json:"total_pages"`
}

type RecordDetail struct {
	Record   models.MovableRecord       `json:"record"`
	Photos   []models.RecordPhoto       `json:"photos"`
	Comments []models.RecordComment     `json:"comments"`
	History  []models.StatusHistoryEntry `json:"history"`
}

type CreateResult struct {
	ID       uuid.UUID           `json:"id"`
	RecordID string              `json:"record_id"`
	Status   models.RecordStatus `json:"status"`
}

type SubmitResult struct {
	Status models.RecordStatus `json:"status"`
}

func FormatRecordID(year int, sequence int) string {
	return fmt.Sprintf("ET-HR-AN-V-%d-%04d", year, sequence)
}
