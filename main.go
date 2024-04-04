package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"log/slog"
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
	textHandler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{})
	logger := slog.New(textHandler)
	postgresPool := database.PostgresPool()
	users := stores.NewUserStore(postgresPool)
	router := http.NewServeMux()
	router.HandleFunc("POST /hook", func(w http.ResponseWriter, r *http.Request) {
    fmt.Println("got here")
		inboundMessage := models.NewInboundMessage(r)
		context, err := contexts.Get(inboundMessage, users)
		if err != nil {
			logger.Error(err.Error())
			return
		}
		switch context {
		case contexts.FirstMessage:
			fmt.Println("Hello this is your first message")
			return
		default:
			logger.Error("Unkown context")
			return
		}
	})
	router.HandleFunc("GET /test", func(w http.ResponseWriter, r *http.Request) {
    fmt.Println("got here 2")
		w.WriteHeader(http.StatusOK)
		return
	})
	server := http.Server{
		Addr:    ":8080",
		Handler: router,
	}
	fmt.Println("Server launched on port: 8080")
	if err := server.ListenAndServe(); err != nil {
		panic(err)
	}
}
