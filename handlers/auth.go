package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
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

func (h *AuthHandler) HandleRegistrationRequest(w http.ResponseWriter, r *http.Request) {
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
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
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
	message := fmt.Sprintf("Voici votre code de verification: %s", authCode.Code)
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

func (h *AuthHandler) HandleAccountVerification(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	if code == "" {
		utils.WriteError(w, http.StatusUnprocessableEntity, nil)
		return
	}
	phone := r.URL.Query().Get("phone")
	if phone == "" {
		utils.WriteError(w, http.StatusUnprocessableEntity, nil)
		return
	}
	dbCode, err := h.authCodes.Get(code, phone)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			utils.WriteError(w, http.StatusBadRequest, nil)
			return
		}
		h.logger.Error(err.Error())
		utils.WriteError(w, http.StatusInternalServerError, nil)
		return
	}
	if dbCode.Used {
		utils.WriteError(w, http.StatusBadRequest, nil)
		return
	}
	err = h.authCodes.SetToUsed(dbCode.ID)
	if err != nil {
		h.logger.Error(err.Error())
	}
	user := &domain.User{
		ID:          ulid.Make().String(),
		Phone:       phone,
		Username:    phone,
		Password:    "",
		AccountType: string(domain.SellerAccountType),
	}
	tx, err := h.txProvider.Provide()
	if err != nil {
		h.logger.Error(err.Error())
		utils.WriteError(w, http.StatusBadRequest, nil)
		return
	}
	err = h.users.Insert(tx, user)
	log.Println("here")
	if err != nil {
		h.logger.Error(err.Error())
		err = tx.Rollback()
		if err != nil {
			h.logger.Error(
				fmt.Sprintf("Error while rolling back  transaction: %s", err.Error()),
			)
		}
		utils.WriteError(w, http.StatusBadRequest, nil)
		return
	}
	log.Println("before creating session")
	session := &domain.Session{
		ID:     ulid.Make().String(),
		Valid:  true,
		UserID: user.ID,
	}
	err = h.sessions.CreateSession(tx, session)
	if err != nil {
		h.logger.Error(err.Error())
		err = tx.Rollback()
		if err != nil {
			h.logger.Error(
				fmt.Sprintf("Error while rolling back  transaction: %s", err.Error()),
			)
		}
		utils.WriteError(w, http.StatusBadRequest, nil)
		return
	}
	log.Println("after creating session")
	err = tx.Commit()
	if err != nil {
		h.logger.Error(fmt.Sprintf("Error while commiting transaction: %s", err.Error()))
		utils.WriteError(w, http.StatusBadRequest, nil)
		return
	}
	utils.WriteData(w, http.StatusOK, map[string]string{
		"session": session.ID,
	})
}

func (h *AuthHandler) HandleRegistrationCompletion(w http.ResponseWriter, r *http.Request) {
	type dto struct {
		Username string `json:"username"`
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
	sessionID := r.Header.Get("Authorization")
	if sessionID == "" {
		utils.WriteError(w, http.StatusUnauthorized, nil)
		return
	}
	session, err := h.sessions.GetSessionByID(sessionID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			utils.WriteError(w, http.StatusUnauthorized, nil)
			return
		}
		h.logger.Error(err.Error())
		utils.WriteError(w, http.StatusInternalServerError, nil)
		return
	}
	hash, err := utils.HashPassword(payload.Password)
	if err != nil {
		h.logger.Error(err.Error())
		utils.WriteError(w, http.StatusInternalServerError, nil)
		return
	}
	err = h.users.UpdateUser(&domain.User{
		ID:       session.UserID,
		Username: payload.Username,
		Password: hash,
	})
	if err != nil {
		h.logger.Error(err.Error())
		utils.WriteError(w, http.StatusInternalServerError, nil)
		return
	}
	utils.WriteData(w, http.StatusOK, nil)
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
		r.Get("/verify", h.HandleAccountVerification)
		r.Post("/registration/request", h.HandleRegistrationRequest)
		r.Post("/registration/complete", h.HandleRegistrationCompletion)
		r.Post("/login", h.HandleLogin)
	})
}
