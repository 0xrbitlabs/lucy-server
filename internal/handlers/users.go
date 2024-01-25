package handlers

import (
	"log/slog"
	"server/internal/store"
)

type UserHandler struct {
	users  *store.Users
	logger *slog.Logger
}

func NewUserHandler(users *store.Users, logger *slog.Logger) *UserHandler {
	return &UserHandler{
		users:  users,
		logger: logger,
	}
}
