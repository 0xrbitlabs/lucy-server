package main

import (
	"fmt"
	"net/http"
	"server/internal/database"
	"server/internal/handlers"
	"server/internal/store"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Failed to read env var file")
		panic(err)
	}
	postgresPool := database.NewPostgresPool()
	users := store.NewUsers(postgresPool)
	userHandler := handlers.NewUserHandler(users)
	r := chi.NewRouter()

	r.Post("/auth/register", userHandler.CreateAccount)
	fmt.Println("Server launched on port 8081")
	err = http.ListenAndServe(":8081", r)
	if err != nil {
		panic(err)
	}
}
