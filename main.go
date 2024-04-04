package main

import (
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"
	"server/internal/app"
	"server/internal/contexts"
	"server/internal/database"
	"server/internal/models"
	"server/internal/stores"

	"github.com/joho/godotenv"
	"github.com/twilio/twilio-go"
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
	accountSID := os.Getenv("ACCOUNT_SID")
	authToken := os.Getenv("AUTH_TOKEN")
	twilioClient := twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: accountSID,
		Password: authToken,
	})
	textHandler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{})
	logger := slog.New(textHandler)
	postgresPool := database.PostgresPool()
	users := stores.NewUserStore(postgresPool)
	mux := http.NewServeMux()
	mux.HandleFunc("POST /", func(w http.ResponseWriter, r *http.Request) {
		inboundMessage := models.NewInboundMessage(r)
		context, err := contexts.Get(inboundMessage, users)
		if err != nil {
			logger.Error(err.Error())
			return
		}
		switch context {
		case contexts.FirstMessage:
			app.HandleFirstMessage(*inboundMessage, twilioClient)
			return
		default:
			logger.Error("Unkown context")
			return
		}
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
