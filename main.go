package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load(); err != nil {
		panic(err)
	}
}

func HandleWebhookConfiguration(w http.ResponseWriter, r *http.Request) {
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

const port = "8080"
const verifyToken = "lucy"

func main() {
	r := chi.NewRouter()

	r.Get("/webhook", HandleWebhookConfiguration)

	// TODO: Make this better
	fmt.Println("Started server on port", port)
	if err := http.ListenAndServe(":"+port, r); err != nil {
		panic(err)
	}
}
