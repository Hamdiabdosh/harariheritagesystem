package auth

import "errors"

var (
	ErrInvalidCredentials  = errors.New("invalid email or password")
	ErrAccountDeactivated  = errors.New("account is deactivated")
	ErrInvalidRefreshToken = errors.New("invalid or expired refresh token")
	ErrInvalidAccessToken  = errors.New("invalid or expired access token")
)
