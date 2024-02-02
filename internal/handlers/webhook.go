package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"server/internal/types"
)

type WebhookHandler struct {
}

func NewWebhookHandler() *WebhookHandler {
	return &WebhookHandler{}
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
		fmt.Printf("Error while decoding payload: %s\n", err.Error())
		w.WriteHeader(http.StatusOK)
		return
	}
	if len(payload.Entry) == 0 {
		fmt.Println("Received empty webhook message")
		w.WriteHeader(http.StatusOK)
		return
	}
	messages := payload.Entry[0].Changes[0].Value.Messages
	if len(messages) == 0 {
		fmt.Println("Not a message sent webhook")
		w.WriteHeader(http.StatusOK)
		return
	}
	messageType := messages[0].Type
	if messageType != "text" {
		fmt.Println("Received non 'Text message received' webhook")
		w.WriteHeader(http.StatusOK)
		return
	}
	message := payload.Entry[0].Changes[0].Value.Messages[0]
	fmt.Printf("%s\n", message)
	w.WriteHeader(http.StatusOK)
	return
}
