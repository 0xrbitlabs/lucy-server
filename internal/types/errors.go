package types

import "errors"

var (
	ErrMessageNotSent = errors.New("Message not sent")

	ErrCodeNotFound = errors.New("Verification code not found")
)
