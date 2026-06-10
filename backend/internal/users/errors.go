package users

import "errors"

var (
	ErrUserNotFound        = errors.New("user not found")
	ErrEmailAlreadyExists  = errors.New("email already exists")
	ErrCannotDeactivateSelf = errors.New("cannot deactivate your own account")
	ErrCannotChangeOwnRole = errors.New("cannot change your own role")
	ErrInvalidRole         = errors.New("invalid role")
	ErrInvalidLanguage     = errors.New("invalid language")
	ErrInvalidPassword     = errors.New("password must be at least 8 characters and contain a number")
	ErrInvalidFullName     = errors.New("full name is required and must be at most 100 characters")
)
