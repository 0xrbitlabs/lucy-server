package types

import (
	"errors"
	"net/http"
)

type ServiceError struct {
	StatusCode int
	ErrorCode  string
}

func (s ServiceError) Error() string {
	return s.ErrorCode
}

var (
	ErrResourceNotFound = errors.New("")
)

var (
	ServiceErrInternal = ServiceError{
		StatusCode: http.StatusInternalServerError,
		ErrorCode:  "InternalServerError",
	}

	ServiceErrDuplicatePhone = ServiceError{
		StatusCode: http.StatusConflict,
		ErrorCode:  ErrDuplicatePhone,
	}
)

const (
	ErrPhoneNotFound  = "PhoneNotFound"
	ErrInternal       = "ErrInternal"
	ErrDuplicatePhone = "DuplicatePhone"
	ErrDuplicateLabel = "DuplicateLabel"
	ErrWrongPassword  = "WrongPassword"
	ErrTokenEncoding  = "TokenEncoding"
)
