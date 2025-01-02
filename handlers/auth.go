package handlers

import (
	"encoding/json"
	"log/slog"
	"lucy/models"
	"lucy/repo"
	"lucy/utils"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/oklog/ulid/v2"
)

type AuthHandler struct {
	logger   *slog.Logger
	users    *repo.UserRepo
	sessions *repo.SessionRepo
}

func NewAuthHandler(
	logger *slog.Logger,
	users *repo.UserRepo,
	sessions *repo.SessionRepo,
) *AuthHandler {
	return &AuthHandler{
		logger:   logger,
		users:    users,
		sessions: sessions,
	}
}

func (h *AuthHandler) HandleLogin(w http.ResponseWriter, r *http.Request) {
  w.Header().Set("Content-Type", "application/json")
	payload := struct {
		PhoneNumber string `json:"phone_number"`
		Password    string `json:"password"`
	}{}
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		h.logger.Error("Error while decoding payload body:", slog.Any("err", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	user, err := h.users.GetUserByPhoneNumber(payload.PhoneNumber)
	if err != nil {
		h.logger.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if user == nil {
		w.WriteHeader(http.StatusNotFound)
		utils.WriteError("user_not_found", w)
		return
	}
	if !utils.PasswordMatchesHash(payload.Password, user.Password) {
		w.WriteHeader(http.StatusBadRequest)
		utils.WriteError("wrong_password", w)
		return
	}
	session := &models.Session{
		ID:     ulid.Make().String(),
		UserID: user.ID,
		Valid:  true,
	}
	err = h.sessions.Insert(session)
	if err != nil {
		h.logger.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	utils.WriteData(map[string]string{"session": session.ID}, w)
	return
}

func (h *AuthHandler) RegisterRoutes(r chi.Router) {
	r.Route("/auth", func(r chi.Router) {
		r.Post("/login", h.HandleLogin)
	})
}
