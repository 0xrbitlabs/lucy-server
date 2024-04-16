package app_errors

import "errors"

var (
	ErrResourceNotFound  error = errors.New("")
	ErrDuplicateResource error = errors.New("")
)
