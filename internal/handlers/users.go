package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/jwtauth/v5"
	"github.com/oklog/ulid/v2"
	"log/slog"
	"net/http"
	"server/internal/store"
	"server/internal/types"
	"server/internal/utils"
)

type UserHandler struct {
	users  *store.Users
	logger *slog.Logger
	jwt    *jwtauth.JWTAuth
}

func NewUserHandler(users *store.Users, logger *slog.Logger, jwt *jwtauth.JWTAuth) *UserHandler {
	return &UserHandler{
		users:  users,
		logger: logger,
		jwt:    jwt,
	}
}

func (h *UserHandler) CreateAccount(w http.ResponseWriter, r *http.Request) {
	payload := new(struct {
		PhoneNumber    string `json:"phone_number"`
		Password       string `json:"password"`
		Name           string `json:"full_name"`
		ProfilePicture string `json:"profile_picture"`
		Description    string `json:"description"`
		Country        string `json:"country"`
		Town           string `json:"town"`
	})
	err := json.NewDecoder(r.Body).Decode(payload)
	if err != nil {
		h.logger.Error(fmt.Sprintf("Error while decoding body: %s", err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	count, err := h.users.CountByPhoneNumber(payload.PhoneNumber)
	if err != nil {
		h.logger.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if count > 0 {
		w.WriteHeader(http.StatusConflict)
		return
	}
	passwordHash, err := utils.Hash(payload.Password)
	if err != nil {
		h.logger.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	newUser := &types.User{
		Id:             ulid.Make().String(),
		UserType:       "seller",
		PhoneNumber:    payload.PhoneNumber,
		Password:       passwordHash,
		Verified:       false,
		Name:           payload.Name,
		ProfilePicture: payload.ProfilePicture,
		Description:    payload.Description,
		Country:        payload.Country,
		Town:           payload.Town,
	}
	err = h.users.Insert(newUser)
	if err != nil {
		h.logger.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	_, token, err := h.jwt.Encode(map[string]interface{}{
		"user_Id": newUser.Id,
	})
	if err != nil {
		h.logger.Error(fmt.Sprintf("Error while encoding JWT token: %s", err.Error()))
		w.Header().Set("X-CONTEXT", "JWT-ENCODING-ERR")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	data, err := json.Marshal(map[string]interface{}{
		"data": map[string]string{
			"token": token,
		},
	})
	if err != nil {
		h.logger.Error(fmt.Sprintf("Error while marshaling token data: %s", err.Error()))
		w.Header().Set("X-CONTEXT", "MARSHALL-ERR")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
  w.Write(data)
	return
}
