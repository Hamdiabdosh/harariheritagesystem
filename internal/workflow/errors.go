package workflow

import "errors"

var (
	ErrRecordNotFound          = errors.New("record not found")
	ErrForbidden               = errors.New("forbidden")
	ErrInvalidRecordType       = errors.New("invalid record type")
	ErrInvalidStatusTransition = errors.New("invalid status transition")
	ErrCommentRequired         = errors.New("comment is required")
)
