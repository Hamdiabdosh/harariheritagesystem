package export

import "errors"

var (
	ErrRecordNotFound    = errors.New("record not found")
	ErrForbidden         = errors.New("forbidden")
	ErrInvalidRecordType = errors.New("invalid record type")
	ErrDraftNotPrintable = errors.New("draft records cannot be printed")
	ErrNotApproved       = errors.New("only approved records can be printed by registrar")
)
