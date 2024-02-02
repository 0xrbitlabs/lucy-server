package handlers

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"server/internal/store"
	"server/internal/types"
	"github.com/oklog/ulid/v2"
)

type WebhookHandler struct {
	users  *store.Users
	logger *slog.Logger
}

func NewWebhookHandler(users *store.Users, logger *slog.Logger) *WebhookHandler {
	return &WebhookHandler{
		users:  users,
		logger: logger,
	}
}

func (h *WebhookHandler) Verify(w http.ResponseWriter, r *http.Request) {
	verifyToken := os.Getenv("VERIFY_TOKEN")
	mode := r.URL.Query().Get("hub.mode")
	token := r.URL.Query().Get("hub.verify_token")
	challenge := r.URL.Query().Get("hub.challenge")

	if mode == "subscribe" && token == verifyToken {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(challenge))
		return
	}
	w.WriteHeader(http.StatusBadRequest)
	return
}

func (h *WebhookHandler) Handle(w http.ResponseWriter, r *http.Request) {
	payload := new(types.WebhookMessage)
	err := json.NewDecoder(r.Body).Decode(payload)
	if err != nil {
		h.logger.Error(fmt.Sprintf("Error while decoding payload: %s", err.Error()))
		w.WriteHeader(http.StatusOK)
		return
	}
	if len(payload.Entry) == 0 {
		h.logger.Debug("Received empty webhook message")
		w.WriteHeader(http.StatusOK)
		return
	}
	messages := payload.Entry[0].Changes[0].Value.Messages
	if len(messages) == 0 {
		h.logger.Debug("Not a message sent webhook")
		w.WriteHeader(http.StatusOK)
		return
	}
	messageType := messages[0].Type
	if messageType != "text" {
		h.logger.Debug("Received non 'Text message received' webhook")
		w.WriteHeader(http.StatusOK)
		return
	}
	message := messages[0]
	userContactInfo := payload.Entry[0].Changes[0].Value.Contacts[0]
	userPhone := message.From
	count, err := h.users.CountByPhoneNumber(userPhone)
	if err != nil {
		h.logger.Error(err.Error())
		w.WriteHeader(http.StatusOK)
		return
	}
	if count == 0 {
		//User is a new user
		newUser := &types.User{
			Id:          ulid.Make().String(),
			UserType:    "regular",
			PhoneNumber: userPhone,
			Password:    "",
			Name:        userContactInfo.Profile.Name,
			Description: "",
			Country:     "",
			Town:        "",
		}
		err := h.users.Insert(newUser)
		if err != nil {
      h.logger.Error(err.Error())
			w.WriteHeader(http.StatusOK)
			return
		}
		//Send back welcome message
		w.WriteHeader(http.StatusOK)
		return
	}
	fmt.Printf("%s\n", message.Text)
	w.WriteHeader(http.StatusOK)
	return
}
