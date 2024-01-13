package types

import "errors"

var (
	ErrCodeNotFound    = errors.New("Verification code not found")
	ErrUniqueViolation = errors.New("Unique constraint violation ")
)
