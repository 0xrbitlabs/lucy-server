package main

import (
	"fmt"
	"log/slog"
	"lucy/server"
	"net/http"
	"os"

	"github.com/go-chi/jwtauth/v5"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
	}))
	jwtAuth := jwtauth.New("HS256", []byte(os.Getenv("JWT_SECRET")), nil)
	server := server.NewServer(
		logger,
		jwtAuth,
	)
	httpServer := &http.Server{
		Addr:    ":9090",
		Handler: server.Router,
	}
	fmt.Println("Server started on 9090")
	err := httpServer.ListenAndServe()
	if err != nil {
		panic(err)
	}
}
