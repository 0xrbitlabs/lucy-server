package handlers

import (
	"log/slog"
	"lucy/middlewares"
	"lucy/repo"
	"lucy/utils"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type UserHandler struct {
	users  *repo.UserRepo
	logger *slog.Logger
}

func NewUserHandler(
	users *repo.UserRepo,
	logger *slog.Logger,
) *UserHandler {
	return &UserHandler{
		users:  users,
		logger: logger,
	}
}

func (h *UserHandler) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	data, err := h.users.GetAll()
	if err != nil {
		h.logger.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	utils.WriteData(map[string]interface{}{
		"data": data,
	}, w)
}

func (h *UserHandler) RegisterRoutes(r chi.Router, m *middlewares.AuthMiddleware) {
	adminOnlyAuth := m.AuthenticateWithRole(true, "admin")
	r.Route("/users", func(r chi.Router) {
		r.With(adminOnlyAuth).Get("/", h.GetAllUsers)
	})
}
