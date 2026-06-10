package models

import (
	"time"

	"github.com/google/uuid"
)

type StatusHistoryEntry struct {
	ID            uuid.UUID  `json:"id"`
	RecordType    RecordType `json:"record_type"`
	RecordID      uuid.UUID  `json:"record_id"`
	ChangedBy     uuid.UUID  `json:"changed_by"`
	ChangedByName string     `json:"changed_by_name"`
	FromStatus    *string    `json:"from_status,omitempty"`
	ToStatus      string     `json:"to_status"`
	Note          *string    `json:"note,omitempty"`
	CreatedAt     time.Time  `json:"created_at"`
}
