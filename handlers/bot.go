package handlers

import (
	"encoding/json"
	"log/slog"
	"lucy/models"
	"lucy/repo"
	"lucy/whatsapp"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/oklog/ulid/v2"
)

type BotHandler struct {
	whatsappClient *whatsapp.Client
	logger         *slog.Logger
	users          *repo.UserRepo
}

func NewBotHandler(
	whatsappClient *whatsapp.Client,
	logger *slog.Logger,
	users *repo.UserRepo,
) *BotHandler {
	return &BotHandler{
		whatsappClient: whatsappClient,
		logger:         logger,
		users:          users,
	}
}

func (h *BotHandler) WebhookConfiguration(w http.ResponseWriter, r *http.Request) {
	verifyToken := "lucy"
	hubMode := r.URL.Query().Get("hub.mode")
	hubVerifyToken := r.URL.Query().Get("hub.verify_token")
	challenge := r.URL.Query().Get("hub.challenge")
	if hubMode == "subscribe" && hubVerifyToken == verifyToken {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(challenge))
		return
	}
	w.WriteHeader(http.StatusForbidden)
}

func (h *BotHandler) TestSendMessage(w http.ResponseWriter, r *http.Request) {
	err := h.whatsappClient.SendBasicMessage("22892423146", "Hello, this is Lucy :)")
	if err != nil {
		h.logger.Error(err.Error())
	}
	return
}

func (h *BotHandler) messageIsValidToBeingProcessed(envelope *whatsapp.Envelope) bool {
	//message is valid if it's only a text message
  if len(envelope.Entry.Changes) == 0 {
    return false
  }
  if len(envelope.Entry.Changes[0].Value.Messages) == 0 {
    return false
  }
	return envelope.Entry.Changes[0].Value.Messages[0].Type == "text"
}

func (h *BotHandler) userIsFirstTimeMessaging(envelope *whatsapp.Envelope) (bool, error) {
  from := envelope.Entry.Changes[0].Value.Messages[0].From
	user, err := h.users.GetUserByPhoneNumber(from)
	return user == nil, err
}

func (h *BotHandler) MessageWebhookEntry(w http.ResponseWriter, r *http.Request) {
	payload := &whatsapp.Payload{}
	err := json.NewDecoder(r.Body).Decode(payload)
	if err != nil {
		h.logger.Error("Error while decoding payload body:", slog.Any("err", err))
		w.WriteHeader(http.StatusOK)
		return
	}
	envelope := &whatsapp.Envelope{
		Object: payload.Object,
	}
	if len(payload.Entry) > 0 {
		envelope.Entry = payload.Entry[0]
	} else {
		w.WriteHeader(http.StatusOK)
		return
	}
	//check if message is valid to being processed
	if !h.messageIsValidToBeingProcessed(envelope) {
		w.WriteHeader(http.StatusOK)
		return
	}
	//check if it's user's first time messaging bot
	ok, err := h.userIsFirstTimeMessaging(envelope)
	if err != nil {
		h.logger.Error(err.Error())
		w.WriteHeader(http.StatusOK)
		return
	}
	if ok {
		//user is messaging for first time
		//store user data
		contact := envelope.Entry.Changes[0].Value.Contacts[0]
		user := &models.User{
			ID:          ulid.Make().String(),
			Username:    contact.Profile.Name,
			PhoneNumber: contact.WaID,
			Password:    "",
			CreatedAt:   time.Now().UTC(),
			AccountType: "regular",
		}
		err = h.users.Insert(user)
		if err != nil {
			h.logger.Error(err.Error())
			w.WriteHeader(http.StatusOK)
			return
		}
		//now that user is stored, send welcome message
		welcomeMessage := `
    Hello, Bienvenue, Je suis Lucy et mon objectif est de vous
    rendre la vie plus facile :D
    `
		err = h.whatsappClient.SendBasicMessage(user.PhoneNumber, welcomeMessage)
		if err != nil {
			h.logger.Error(err.Error())
		}
		w.WriteHeader(http.StatusOK)
		return
	}
	//user was not first time messaging
	//handle accordingly
	w.WriteHeader(http.StatusOK)
	return
}

func (h *BotHandler) RegisterRoutes(r chi.Router) {
	r.Get("/lucy", h.WebhookConfiguration)
	r.Post("/lucy", h.MessageWebhookEntry)
}
