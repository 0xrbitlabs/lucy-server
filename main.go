package main

import (
	"log"
	"log/slog"
	"lucy/db"
	"lucy/handlers"
	"lucy/providers"
	"lucy/repositories"
	"lucy/services"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	port := os.Getenv("PORT")

	r := chi.NewRouter()

	r.Use(
		middleware.Recoverer,
		middleware.Logger,
	)

	postgresPool := db.PostgresDB()

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
	}))
	jwtProvider := providers.NewJWTProvider()

	userRepo := repositories.NewUserRepo(postgresPool)

	userService := services.NewUserService(userRepo, logger)
	authService := services.NewAuthService(userRepo, logger, jwtProvider)

	userHandler := handlers.NewUserHandler(userService, logger)
	authHandler := handlers.NewAuthHandler(authService, logger)

	r.Route("/users", func(r chi.Router) {
		r.Post("/", userHandler.HandleCreateAdminAccount)
	})

	r.Route("/auth", func(r chi.Router) {
		r.Post("/login", authHandler.HandleLogin)
	})

	server := http.Server{
		Addr:         net.JoinHostPort("0.0.0.0", port),
		Handler:      r,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	log.Println("Starting server on port", port)
	if err := server.ListenAndServe(); err != nil {
		panic(err)
	}
}
