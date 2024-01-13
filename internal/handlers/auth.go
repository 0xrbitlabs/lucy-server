package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"server/internal/store"
	"server/internal/types"
	"server/internal/utils"
	"time"

	"github.com/oklog/ulid/v2"
)

type AuthHandler struct {
	otpCodes *store.OTPCodes
	users    *store.Users
	logger   *slog.Logger
}

func NewAuthHandler(otpCodes *store.OTPCodes, users *store.Users, logger *slog.Logger) *AuthHandler {
	return &AuthHandler{
		otpCodes: otpCodes,
		users:    users,
		logger:   logger,
	}
}

func (h *AuthHandler) RequestVerificationCode(w http.ResponseWriter, r *http.Request) {
	phoneNumber := r.URL.Query().Get("phone_number")
	if phoneNumber == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	otpCode := fmt.Sprint(time.Now().Nanosecond())[:6]
	otpMessage := fmt.Sprintf("Votre code verification est %s", otpCode)
	err := utils.SendMessage(phoneNumber, otpMessage)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		h.logger.Error(err.Error())
		return
	}
	codeKey := fmt.Sprintf("%s:%s", otpCode, phoneNumber)
	err = h.otpCodes.Set(codeKey)
	if err != nil {
		h.logger.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	return
}

func (h *AuthHandler) VerifyPhoneNumber(w http.ResponseWriter, r *http.Request) {
	phoneNumber := r.URL.Query().Get("phone_number")
	code := r.URL.Query().Get("code")
	if phoneNumber == "" || code == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	codeKey := fmt.Sprintf("%s:%s", code, phoneNumber)
	_, err := h.otpCodes.Get(codeKey)
	if err != nil {
		if errors.Is(err, types.ErrCodeNotFound) {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		h.logger.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	verificationProofId := ulid.Make().String()
	err = h.otpCodes.SetVerificationProof(phoneNumber, verificationProofId)
	if err != nil {
		h.logger.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = h.otpCodes.Delete(codeKey)
	if err != nil {
		h.logger.Debug(err.Error())
	}
	w.WriteHeader(http.StatusOK)
	return
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
	verifiedNumber, err := h.otpCodes.Get(payload.VerificationProof)
	if err != nil {
		if errors.Is(err, types.ErrCodeNotFound) {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}
	if verifiedNumber != payload.PhoneNumber {
		w.WriteHeader(http.StatusBadRequest)
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
	}
	w.WriteHeader(http.StatusCreated)
	return
}

func (h *AuthHandler) Login() {

}

func (h *AuthHandler) Logout() {

}
