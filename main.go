package main

import (
	"fmt"
	"log"
	"log/slog"
	"lucy/handlers"
	"lucy/middlewares"
	"lucy/repo"
	"lucy/whatsapp"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func getDBConnectionPool() (*sqlx.DB, error) {
	db, err := sqlx.Connect("postgres", os.Getenv("DB_URL"))
	return db, err
}

func main() {

	if os.Getenv("STAGE") != "PROD" {
		godotenv.Load()
	}
	port := fmt.Sprintf(":%s", os.Getenv("PORT"))

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
	}))

	db, err := getDBConnectionPool()
	if err != nil {
		logger.Error("Failed to connect to database", slog.Any("err", err))
		os.Exit(1)
	}
	log.Println("Connected to Database!")

	userRepo := repo.NewUserRepo(db)
	sessionRepo := repo.NewSessionRepo(db)
	productCategoryRepo := repo.NewProductCategoryRepo(db)

	authHandler := handlers.NewAuthHandler(
		logger,
		userRepo,
		sessionRepo,
	)
	productCategoryHandler := handlers.NewProductCategoryHandler(
		productCategoryRepo,
		logger,
	)

	accessToken := os.Getenv("ACCESS_TOKEN")
	phoneNumberID := os.Getenv("PHONE_NUMBER_ID")
	whatsappClient := whatsapp.NewClient(accessToken, phoneNumberID)
	botHandler := handlers.NewBotHandler(whatsappClient, logger, userRepo)

	authMiddleware := middlewares.NewAuthMiddleware(
		userRepo,
		sessionRepo,
		logger,
	)

	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Everything is good :)\n"))
	})

	authHandler.RegisterRoutes(r)
	productCategoryHandler.RegisterRoutes(r, authMiddleware)
	botHandler.RegisterRoutes(r)

	server := http.Server{}
	server.Handler = r
	server.Addr = port
	server.IdleTimeout = time.Minute
	server.ReadTimeout = time.Minute
	server.WriteTimeout = time.Minute

	log.Println("Starting server on port", os.Getenv("PORT"))
	if err := server.ListenAndServe(); err != nil {
		panic(err)
	}
}
