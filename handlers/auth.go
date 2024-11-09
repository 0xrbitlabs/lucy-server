package handlers

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/joseph0x45/lucy/bot"
	"github.com/joseph0x45/lucy/domain"
	"github.com/joseph0x45/lucy/repository"
	"github.com/joseph0x45/lucy/utils"
)

type AuthHandler struct {
	users     *repository.UserRepo
	sessions  *repository.SessionRepo
	logger    *slog.Logger
	authCodes *repository.AuthCodeRepo
}

func NewAuthHandler(
	users *repository.UserRepo,
	sessions *repository.SessionRepo,
	logger *slog.Logger,
	authCodes *repository.AuthCodeRepo,
) *AuthHandler {
	return &AuthHandler{
		users:     users,
		sessions:  sessions,
		logger:    logger,
		authCodes: authCodes,
	}
}

func (h *AuthHandler) HandleRegistration(w http.ResponseWriter, r *http.Request) {
	type dto struct {
		Phone string `json:"phone"`
	}
	payload := &dto{}
	err := json.NewDecoder(r.Body).Decode(payload)
	if err != nil {
		h.logger.Error(
			fmt.Sprintf("Error while decoding body: %s", err.Error()),
		)
		utils.WriteError(w, http.StatusInternalServerError, nil)
		return
	}
	if payload.Phone == "" {
		utils.WriteError(w, http.StatusBadRequest, nil)
		return
	}
	user, err := h.users.GetByPhone(payload.Phone)
	if err != nil {
		h.logger.Error(err.Error())
		utils.WriteError(w, http.StatusInternalServerError, nil)
		return
	}
	if user != nil {
		utils.WriteError(w, http.StatusConflict, nil)
		return
	}
	authCode := &domain.AuthCode{
		Code:         utils.GenerateRandomDigit(),
		Used:         false,
		GeneratedFor: payload.Phone,
	}
	message := fmt.Sprintf("Voici votre code d'authentification: %s", authCode.Code)
	err = bot.SendMessage(payload.Phone, message)
	if err != nil {
		h.logger.Error(err.Error())
		utils.WriteError(w, http.StatusInternalServerError, nil)
		return
	}
	err = h.authCodes.Insert(authCode)
	if err != nil {
		h.logger.Error(err.Error())
		utils.WriteError(w, http.StatusInternalServerError, nil)
		return
	}
	utils.WriteData(w, http.StatusOK, nil)
	return
}

func (h *AuthHandler) HandleLogin(w http.ResponseWriter, r *http.Request) {

}

func (h *AuthHandler) RegisterRoutes(r chi.Router) {
	r.Route("/auth", func(r chi.Router) {
		r.Post("/register", h.HandleRegistration)
		r.Post("/login", h.HandleLogin)
	})
}
