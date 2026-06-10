package models

import (
	"time"

	"github.com/google/uuid"
)

type RecordType string

const (
	RecordTypeImmovable RecordType = "immovable"
	RecordTypeMovable   RecordType = "movable"
)

func (t RecordType) IsValid() bool {
	return t == RecordTypeImmovable || t == RecordTypeMovable
}

type RecordPhoto struct {
	ID            uuid.UUID  `json:"id"`
	RecordType    RecordType `json:"record_type"`
	RecordID      uuid.UUID  `json:"record_id"`
	FilePath      string     `json:"file_path"`
	FileName      *string    `json:"file_name,omitempty"`
	FileSizeBytes *int       `json:"file_size_bytes,omitempty"`
	UploadedBy    uuid.UUID  `json:"uploaded_by"`
	CreatedAt     time.Time  `json:"created_at"`
}
