package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/oklog/ulid/v2"
	"log/slog"
	"net/http"
	"server/internal/store"
	"server/internal/types"
	"server/internal/utils"
)

type AuthHandler struct {
	users    *store.Users
	sessions *store.Sessions
	logger   *slog.Logger
	codes    *store.VerificationCodes
}

func NewAuthHandler(users *store.Users, sessions *store.Sessions, logger *slog.Logger, codes *store.VerificationCodes) *AuthHandler {
	return &AuthHandler{
		users:    users,
		sessions: sessions,
		logger:   logger,
		codes:    codes,
	}
}

func (h *AuthHandler) CompleteRegistration(w http.ResponseWriter, r *http.Request) {
	payload := new(struct {
		Password         string `json:"password"`
		Name             string `json:"name"`
		Description      string `json:"description"`
		Country          string `json:"country"`
		Town             string `json:"town"`
		PhoneNumber      string `json:"phone_number"`
		VerificationCode string `json:"verification_code"`
	})
	err := json.NewDecoder(r.Body).Decode(payload)
	if err != nil {
		h.logger.Error(fmt.Sprintf("Error while decoding request body: %s", err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	code, err := h.codes.Get(payload.PhoneNumber)
	if err != nil {
		if errors.Is(err, types.ErrCodeNotFound) {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		h.logger.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if code != payload.VerificationCode {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	dbUser, err := h.users.GetByPhoneNumber(payload.PhoneNumber)
	if err != nil {
		if errors.Is(err, types.ErrUserNotFound) {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		h.logger.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	hashedPassword, err := utils.Hash(payload.Password)
	if err != nil {
		h.logger.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	data := &types.UpdateUserInfoPayload{
		ID:          dbUser.Id,
		Name:        payload.Name,
		Password:    hashedPassword,
		Description: payload.Description,
		Country:     payload.Country,
		Town:        payload.Town,
	}
	err = h.users.UpdateInfo(data)
	if err != nil {
		h.logger.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = h.codes.Delete(payload.PhoneNumber)
	if err != nil {
		h.logger.Error(err.Error())
	}
	sessionID := ulid.Make().String()
	err = h.sessions.Create(sessionID, dbUser.Id)
	if err != nil {
		h.logger.Error(err.Error())
		w.Header().Add("X-CONTEXT", "ERR_CREATE_SESSION")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	respData, err := json.Marshal(map[string]interface{}{
		"data": map[string]string{
			"session": sessionID,
		},
	})
	if err != nil {
		h.logger.Error(err.Error())
		w.Header().Add("X-CONTEXT", "ERR_MARSHALL_DATA")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(respData)
	return
}
