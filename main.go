package main

import (
	"encoding/json"
	"fmt"
	"github.com/joho/godotenv"
	"log/slog"
	"net"
	"net/http"
	"os"
	"server/internal/contexts"
	"server/internal/database"
	"server/internal/models"
	"server/internal/stores"
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
	port := os.Getenv("PORT")
	textHandler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{})
	logger := slog.New(textHandler)
	postgresPool := database.PostgresPool()
	users := stores.NewUserStore(postgresPool)
	mux := http.NewServeMux()
	mux.HandleFunc("POST /", func(w http.ResponseWriter, r *http.Request) {
		payload := new(models.InboundMessage)
		err := json.NewDecoder(r.Body).Decode(payload)
		if err != nil {
			logger.Error(err.Error())
			w.WriteHeader(http.StatusOK)
			return
		}
		context, _ := contexts.Get(payload, users)
		switch context {
		case contexts.FirstMessage:
			fmt.Println("Hey, seems like it was your first message")
		}
		w.WriteHeader(http.StatusOK)
		return
	})
	server := http.Server{
		Addr:    net.JoinHostPort("0.0.0.0", port),
		Handler: mux,
	}
	fmt.Println("Server launched on port: ", port)
	if err := server.ListenAndServe(); err != nil {
		panic(err)
	}
}
