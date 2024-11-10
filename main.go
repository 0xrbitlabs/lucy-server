package main

import (
	"encoding/json"
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	"github.com/joseph0x45/lucy/handlers"
	"github.com/joseph0x45/lucy/providers"
	"github.com/joseph0x45/lucy/repository"
	whatsapptypes "github.com/joseph0x45/lucy/whatsapp_types"
	_ "github.com/lib/pq"
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

func connectToDB() *sqlx.DB {
	db, err := sqlx.Open("postgres", os.Getenv("DB_URL"))
	if err != nil {
		panic(err)
	}
	log.Println("Connected to Postgres")
	return db
}

func main() {
	r := chi.NewRouter()
	r.Use(
		middleware.Logger,
		middleware.Recoverer,
	)
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
	}))

	db := connectToDB()

  txProvider := providers.NewTxProvider(db) 

	userRepo := repository.NewUserRepo(db)
	sessionRepo := repository.NewSessionRepo(db)
	// productRepo := repository.NewProductRepo(db)
	authCodeRepo := repository.NewAuthCodeRepo(db)

	authHandler := handlers.NewAuthHandler(
		userRepo,
		sessionRepo,
		logger,
		authCodeRepo,
    txProvider,
	)

	port := os.Getenv("PORT")

	r.Get("/webhook", HandleWebhookConfiguration)

	r.Post("/webhook", func(w http.ResponseWriter, r *http.Request) {
		payload := &whatsapptypes.Payload{}
		err := json.NewDecoder(r.Body).Decode(payload)
		if err != nil {
			log.Println("Error while decoding request body:", err.Error())
			w.WriteHeader(http.StatusOK)
			return
		}
		envelope := whatsapptypes.Envelope{
			Object: payload.Object,
		}
		if len(payload.Entry) > 0 {
			envelope.Entry = payload.Entry[0]
		} else {
			w.WriteHeader(http.StatusOK)
			return
		}
	})

	authHandler.RegisterRoutes(r)

	// TODO: Make this better
	log.Println("Started server on port", port)
	if err := http.ListenAndServe(":"+port, r); err != nil {
		panic(err)
	}
}
