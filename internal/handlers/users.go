package handlers

import (
	"log/slog"
	"net/http"
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

func (h *UserHandler) CreateAccount(w http.ResponseWriter, r *http.Request) {
}
