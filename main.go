package main

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
	"log/slog"
	"net/http"
	"os"
	"server/types"
)

func main() {
	env := os.Getenv("ENV")
	if env != "PROD" {
		err := godotenv.Load()
		if err != nil {
			fmt.Println("Failed to read env var file")
			panic(err)
		}
	}
	textHandler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{})
	logger := slog.New(textHandler)
	_ = logger
	r := chi.NewRouter()

	r.Route("/hook", func(r chi.Router) {
		r.Post("/", func(w http.ResponseWriter, r *http.Request) {
			inboundMessage := types.NewInboundMessage(r)
			fmt.Printf("%+v", inboundMessage)
			return
		})
	})

	r.Route("/auth", func(r chi.Router) {
	})
	fmt.Println("Server launched on port 8081")
	err := http.ListenAndServe(":8081", r)
	if err != nil {
		panic(err)
	}
}
