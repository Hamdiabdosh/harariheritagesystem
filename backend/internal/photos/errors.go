package photos

import "errors"

var (
	ErrPhotoNotFound      = errors.New("photo not found")
	ErrForbidden          = errors.New("forbidden")
	ErrNotEditable        = errors.New("record is not in an editable status")
	ErrInvalidRecordType  = errors.New("invalid record type")
	ErrMaxPhotosReached   = errors.New("maximum number of photos reached")
	ErrFileTooLarge       = errors.New("file exceeds maximum size of 5MB")
	ErrUnsupportedType    = errors.New("unsupported file type; only JPG and PNG are allowed")
	ErrRecordNotFound     = errors.New("record not found")
)
