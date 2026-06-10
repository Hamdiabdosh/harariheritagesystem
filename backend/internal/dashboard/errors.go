package dashboard

import "errors"

var (
	errInvalidRecordType = errors.New("record type must be immovable or movable")
	errInvalidStatus     = errors.New("invalid status filter")
	errInvalidDate       = errors.New("invalid date filter")
)
