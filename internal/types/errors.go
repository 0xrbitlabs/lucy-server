package types

import "errors"

var (
	ErrCodeNotFound    = errors.New("Verification code not found")
	ErrSessionNotFound = errors.New("Session not found")
	ErrUniqueViolation = errors.New("Unique constraint violation ")
	ErrUserNotFound    = errors.New("User not found")
)
