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
}

func NewAuthHandler(users *store.Users, sessions *store.Sessions, logger *slog.Logger) *AuthHandler {
	return &AuthHandler{
		users:    users,
		sessions: sessions,
		logger:   logger,
	}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	payload := new(struct {
		PhoneNumber       string `json:"phone_number"`
		Password          string `json:"password"`
		Name              string `json:"name"`
		ProfilePicture    string `json:"profile_picture"`
		Description       string `json:"description"`
		Country           string `json:"country"`
		Town              string `json:"town"`
		VerificationProof string `json:"verification_proof"`
	})
	err := json.NewDecoder(r.Body).Decode(payload)
	if err != nil {
		h.logger.Error(fmt.Sprintf("Error while decoding request body: %s", err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	hashedPassword, err := utils.Hash(payload.Password)
	if err != nil {
		h.logger.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	newUser := &types.User{
		Id:             ulid.Make().String(),
		PhoneNumber:    payload.PhoneNumber,
		Password:       hashedPassword,
		Name:           payload.Name,
		Description:    payload.Description,
		ProfilePicture: payload.ProfilePicture,
		Country:        payload.Country,
		Town:           payload.Town,
		UserType:       "seller",
	}
	err = h.users.Insert(newUser)
	if err != nil {
		if errors.Is(err, types.ErrUniqueViolation) {
			h.logger.Debug("A unique constraint violation occured")
			w.WriteHeader(http.StatusConflict)
			return
		}
		h.logger.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	sessionId := ulid.Make().String()
	err = h.sessions.Create(sessionId, newUser.Id)
	if err != nil {
		h.logger.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	data, err := json.Marshal(map[string]interface{}{
		"data": map[string]string{
			"session": sessionId,
		},
	})
	if err != nil {
		h.logger.Error(fmt.Sprintf("Error while marshalling data: %s", err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Write(data)
	return
}
