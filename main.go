package main

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load(); err != nil {
		panic(err)
	}
}

func HandleWebhookConfiguration(w http.ResponseWriter, r *http.Request) {
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

func main() {
	r := chi.NewRouter()

	port := os.Getenv("PORT")

	r.Post("/webhook", func(w http.ResponseWriter, r *http.Request) {
		data, err := io.ReadAll(r.Body)
		if err != nil {
			fmt.Println(err.Error())
		}
		fmt.Println(string(data))
		w.WriteHeader(http.StatusOK)
	})

	// TODO: Make this better
	fmt.Println("Started server on port", port)
	if err := http.ListenAndServe(":"+port, r); err != nil {
		panic(err)
	}
}
