package main

import (
	"log"
	"log/slog"
	"lucy/db"
	"lucy/handlers"
	"lucy/repositories"
	"lucy/services"
	"net"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	port := os.Getenv("port")

	r := chi.NewRouter()

	r.Use(
		middleware.Recoverer,
		middleware.Logger,
	)

	postgresPool := db.PostgresDB()

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
	}))

	userRepo := repositories.NewUserRepo(postgresPool)

	userService := services.NewUserService(userRepo, logger)

	userHandler := handlers.NewUserHandler(userService, logger)

	r.Route("/users", func(r chi.Router) {
		r.Post("/", userHandler.HandleCreateAdminAccount)
	})

	r.Route("/auth", func(r chi.Router) {
		r.Post("/login", func(w http.ResponseWriter, r *http.Request) {

		})
	})

	server := http.Server{
		Addr: net.JoinHostPort("0.0.0.0", port),
	}

	log.Println("Starting server on port", port)
	if err := server.ListenAndServe(); err != nil {
		panic(err)
	}
}
