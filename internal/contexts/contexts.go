package contexts

import (
	"errors"
	appErrors "server/internal/errors"
	"server/internal/models"
	"server/internal/stores"
)

type Context int

const (
	Error Context = iota
	FirstMessage
	RegistrationRequest
	Query
)

func Get(message *models.InboundMessage, users *stores.Users) (Context, error) {
	user, err := users.GetByPhoneNumber(message.From)
	_ = user
	if err != nil {
		if errors.Is(err, appErrors.ErrUserNotFound) {
			return FirstMessage, nil
		}
		return Error, err
	}
	return Query, nil
}
