package app_errors

type ErrorCode string

var (
	ErrBadRequest        ErrorCode = "400_000"
	ErrDuplicateResource ErrorCode = "409_000"
	ErrInternal          ErrorCode = "500_000"
	ErrTokenEncoding     ErrorCode = "500_001"
)
