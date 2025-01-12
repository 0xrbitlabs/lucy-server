package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
	"lucy/middlewares"
	"lucy/models"
	"lucy/repo"
	"lucy/utils"
	"lucy/whatsapp"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/oklog/ulid/v2"
)

type AuthHandler struct {
	logger            *slog.Logger
	users             *repo.UserRepo
	sessions          *repo.SessionRepo
	verificationCodes *repo.VerificationCodeRepo
	whatsappClient    *whatsapp.Client
}

func NewAuthHandler(
	logger *slog.Logger,
	users *repo.UserRepo,
	sessions *repo.SessionRepo,
	verificationCodes *repo.VerificationCodeRepo,
	whatsappClient *whatsapp.Client,
) *AuthHandler {
	return &AuthHandler{
		logger:            logger,
		users:             users,
		sessions:          sessions,
		verificationCodes: verificationCodes,
		whatsappClient:    whatsappClient,
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
}

func (h *AuthHandler) RegisterAsSeller(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	payload := &struct {
		Username    string `json:"username"`
		PhoneNumber string `json:"phone_number"`
		Password    string `json:"password"`
	}{}
	err := json.NewDecoder(r.Body).Decode(payload)
	if err != nil {
		h.logger.Error("Error while decoding payload body:", slog.Any("err", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if payload.Username == "" || payload.Password == "" || payload.PhoneNumber == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	dbUser, err := h.users.GetUserByPhoneNumber("")
	if err != nil {
		h.logger.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if dbUser != nil {
		w.WriteHeader(http.StatusConflict)
		return
	}
	hash, err := utils.GenerateHash(payload.Password)
	if err != nil {
		h.logger.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	newUser := &models.User{
		ID:          ulid.Make().String(),
		Username:    payload.Username,
		PhoneNumber: payload.PhoneNumber,
		Password:    hash,
		CreatedAt:   time.Now(),
		AccountType: "seller",
	}
	err = h.users.Insert(newUser)
	if err != nil {
		h.logger.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	session := &models.Session{
		ID:     ulid.Make().String(),
		UserID: newUser.ID,
		Valid:  true,
	}
	err = h.sessions.Insert(session)
	if err != nil {
		w.WriteHeader(http.StatusOK)
		utils.WriteData(map[string]interface{}{
			"session": "",
		}, w)
		return
	}
	w.WriteHeader(http.StatusOK)
	utils.WriteData(map[string]string{
		"session": session.ID,
	}, w)
}

func (h *AuthHandler) RequestProfileVerificationCode(w http.ResponseWriter, r *http.Request) {
	log.Println("came here")
	w.Header().Set("Content-Type", "application/json")
	currentUser, ok := r.Context().Value("user").(*models.User)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	verificationCode := &models.VerificationCode{
		Code:         utils.GenerateRandomDigit(),
		GeneratedFor: currentUser.ID,
		GeneratedAt:  time.Now(),
		Used:         false,
	}
	err := h.verificationCodes.Insert(verificationCode)
	if err != nil {
		h.logger.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	message := fmt.Sprintf("Votre code de verification est le %s", verificationCode.Code)
	err = h.whatsappClient.SendVerificationCodeMessage(currentUser.PhoneNumber, message)
	if err != nil {
		h.logger.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	return
}

func (h *AuthHandler) VerifyProfile(w http.ResponseWriter, r *http.Request) {}

func (h *AuthHandler) RegisterRoutes(r chi.Router, m *middlewares.AuthMiddleware) {
  auth := m.AuthenticateWithRole()
	r.Route("/auth", func(r chi.Router) {
		r.Post("/login", h.HandleLogin)
		r.Post("/register", h.RegisterAsSeller)
		r.With(auth).Post("/verify", h.RequestProfileVerificationCode)
	})
}
