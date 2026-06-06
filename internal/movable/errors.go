package movable

import "errors"

var (
	ErrRecordNotFound          = errors.New("record not found")
	ErrForbidden               = errors.New("forbidden")
	ErrNotEditable             = errors.New("record is not in an editable status")
	ErrInvalidStatusTransition = errors.New("invalid status transition")
	ErrValidationFailed        = errors.New("validation failed")
)

type ValidationError struct {
	Fields map[string]string
}

func (e *ValidationError) Error() string {
	return ErrValidationFailed.Error()
}

func (e *ValidationError) Unwrap() error {
	return ErrValidationFailed
}

func NewValidationError(fields map[string]string) *ValidationError {
	return &ValidationError{Fields: fields}
}
