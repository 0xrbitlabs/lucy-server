package main

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
	"log/slog"
	"net/http"
	"os"
	"server/internal/database"
	"server/internal/handlers"
	"server/internal/store"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Failed to read env var file")
		panic(err)
	}
	textHandler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{})
	logger := slog.New(textHandler)
	postgresPool := database.NewPostgresPool()
	redisClient := database.RedisClient()
	users := store.NewUsers(postgresPool)
	otpCodes := store.NewOTPCodes(redisClient)
	userHandler := handlers.NewUserHandler(users, logger)
	authHandler := handlers.NewAuthHandler(otpCodes, users, logger)
	r := chi.NewRouter()

	r.Route("/auth", func(r chi.Router) {
		r.Route("/verification", func(r chi.Router) {
			r.Get("/request", authHandler.RequestVerificationCode)
		})

		r.Post("/register", userHandler.CreateAccount)
	})
	fmt.Println("Server launched on port 8081")
	err = http.ListenAndServe(":8081", r)
	if err != nil {
		panic(err)
	}
}
