package models

import (
	"time"

	"github.com/google/uuid"
)

type RecordComment struct {
	ID          uuid.UUID  `json:"id"`
	RecordType  RecordType `json:"record_type"`
	RecordID    uuid.UUID  `json:"record_id"`
	AuthorID    uuid.UUID  `json:"author_id"`
	AuthorName  string     `json:"author_name"`
	CommentText string     `json:"comment_text"`
	CreatedAt   time.Time  `json:"created_at"`
}
