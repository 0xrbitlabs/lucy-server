package handlers

import (
	"fmt"
	"log/slog"
	"net/http"
	"server/internal/store"
	"server/internal/utils"
	"time"
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

func (h *AuthHandler) VerifyPhoneNumber() {}

func (h *AuthHandler) Register() {

}

func (h *AuthHandler) Login() {

}

func (h *AuthHandler) Logout() {

}
