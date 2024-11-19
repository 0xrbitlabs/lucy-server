package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/joseph0x45/lucy/bot"
	"github.com/joseph0x45/lucy/domain"
	"github.com/joseph0x45/lucy/providers"
	"github.com/joseph0x45/lucy/repository"
	"github.com/joseph0x45/lucy/utils"
	"github.com/oklog/ulid/v2"
)

type AuthHandler struct {
	users      *repository.UserRepo
	sessions   *repository.SessionRepo
	logger     *slog.Logger
	authCodes  *repository.AuthCodeRepo
	txProvider *providers.TxProvider
}

func NewAuthHandler(
	users *repository.UserRepo,
	sessions *repository.SessionRepo,
	logger *slog.Logger,
	authCodes *repository.AuthCodeRepo,
	txProvider *providers.TxProvider,
) *AuthHandler {
	return &AuthHandler{
		users:      users,
		sessions:   sessions,
		logger:     logger,
		authCodes:  authCodes,
		txProvider: txProvider,
	}
}

func (h *AuthHandler) HandleVerificationCodeRequest(w http.ResponseWriter, r *http.Request) {
	phone := r.URL.Query().Get("phone")
	if phone == "" {
		utils.WriteError(w, http.StatusBadRequest, nil)
		return
	}
	user, err := h.users.GetByPhone(phone)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		h.logger.Error(err.Error())
		utils.WriteError(w, http.StatusInternalServerError, nil)
		return
	}
	if user != nil {
		utils.WriteError(w, http.StatusConflict, nil)
		return
	}
	// send verification message to user
	code := utils.GenerateRandomDigit()
	authCode := &domain.AuthCode{
		Code:         code,
		Used:         false,
		GeneratedFor: phone,
	}
	err = h.authCodes.Insert(authCode)
	if err != nil {
		h.logger.Error(err.Error())
		utils.WriteError(w, http.StatusInternalServerError, nil)
		return
	}
	message := fmt.Sprintf("Voici votre code de verification: %s\n", code)
	fmt.Println("Sending verification code to user", message)
	err = bot.SendMessage(phone, message)
	if err != nil {
		h.logger.Error(err.Error())
		utils.WriteError(w, http.StatusInternalServerError, nil)
		return
	}
	utils.WriteData(w, http.StatusOK, nil)
}

func (h *AuthHandler) HandleVerificationCodeConfirmation(w http.ResponseWriter, r *http.Request) {
	phone := r.URL.Query().Get("phone")
	code := r.URL.Query().Get("code")
	if phone == "" || code == "" {
		utils.WriteError(w, http.StatusUnprocessableEntity, "phone or code missing")
		return
	}
	dbCode, err := h.authCodes.Get(code, phone)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		h.logger.Error(err.Error())
		utils.WriteError(w, http.StatusInternalServerError, nil)
		return
	}
	if dbCode == nil {
		utils.WriteError(w, http.StatusBadRequest, "invalid_verification_code")
		return
	}
	utils.WriteData(w, http.StatusOK, nil)
}

func (h *AuthHandler) HandleRegistration(w http.ResponseWriter, r *http.Request) {
	type dto struct {
		Phone            string `json:"phone"`
		VerificationCode string `json:"verification_code"`
		Username         string `json:"username"`
		Password         string `json:"password"`
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
		utils.WriteError(w, http.StatusBadRequest, "invalid_phone")
		return
	}
	if payload.VerificationCode == "" {
		utils.WriteError(w, http.StatusBadRequest, "invalid_verification_code")
		return
	}
	if payload.Username == "" {
		utils.WriteError(w, http.StatusBadRequest, "invalid_username")
		return
	}
	if payload.Password == "" {
		utils.WriteError(w, http.StatusBadRequest, "invalid_password")
		return
	}
	// check if phone is not already being used
	exists, err := h.users.Exists(payload.Phone)
	if err != nil {
		h.logger.Error(err.Error())
		utils.WriteError(w, http.StatusInternalServerError, nil)
		return
	}
	if exists {
		utils.WriteError(w, http.StatusConflict, nil)
		return
	}
	dbCode, err := h.authCodes.Get(payload.VerificationCode, payload.Phone)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		h.logger.Error(err.Error())
		utils.WriteError(w, http.StatusInternalServerError, nil)
		return
	}
	if dbCode == nil {
		utils.WriteError(w, http.StatusBadRequest, "invalid_verification_code")
		return
	}
	hash, err := utils.HashPassword(payload.Password)
	if err != nil {
		h.logger.Error(err.Error())
		utils.WriteError(w, http.StatusInternalServerError, nil)
		return
	}
	user := &domain.User{
		ID:          ulid.Make().String(),
		Phone:       payload.Phone,
		Username:    payload.Username,
		Password:    hash,
		AccountType: string(domain.SellerAccountType),
	}
	err = h.users.Insert(nil, user)
	if err != nil {
		h.logger.Error(err.Error())
		utils.WriteError(w, http.StatusInternalServerError, nil)
		return
	}
	err = h.authCodes.SetToUsed(dbCode.ID)
	if err != nil {
		h.logger.Error(err.Error())
	}
	utils.WriteData(w, http.StatusCreated, nil)
}

func (h *AuthHandler) HandleLogin(w http.ResponseWriter, r *http.Request) {
	type dto struct {
		Phone    string `json:"phone"`
		Password string `json:"password"`
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
	if payload.Phone == "" || payload.Password == "" {
		utils.WriteError(w, http.StatusUnprocessableEntity, nil)
		return
	}
	user, err := h.users.GetByPhone(payload.Phone)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			utils.WriteError(w, http.StatusBadRequest, "user_not_found")
			return
		}
		h.logger.Error(err.Error())
		utils.WriteError(w, http.StatusInternalServerError, nil)
		return
	}
	if ok := utils.HashMatches(user.Password, payload.Password); !ok {
		utils.WriteError(w, http.StatusBadRequest, "wrong_password")
		return
	}
	session := &domain.Session{
		ID:     ulid.Make().String(),
		Valid:  true,
		UserID: user.ID,
	}
	err = h.sessions.CreateSession(nil, session)
	if err != nil {
		h.logger.Error(err.Error())
		utils.WriteError(w, http.StatusInternalServerError, nil)
		return
	}
	utils.WriteData(w, http.StatusOK, map[string]string{
		"session": session.ID,
	})
}

func (h *AuthHandler) RegisterRoutes(r chi.Router) {
	r.Route("/auth", func(r chi.Router) {
		r.Get("/verification/request", h.HandleVerificationCodeRequest)
		r.Get("/verification/confirm", h.HandleVerificationCodeConfirmation)
		r.Post("/register", h.HandleRegistration)
		r.Post("/login", h.HandleLogin)
	})
}
