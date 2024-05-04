package main

import (
	"log"
	"net"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	port, ok := os.LookupEnv("port")
	if !ok {
		port = "8080"
	}

	r := chi.NewRouter()

	r.Use(
		middleware.Recoverer,
		middleware.Logger,
	)

	server := http.Server{
		Addr: net.JoinHostPort("0.0.0.0", port),
	}

  log.Println("Starting server on port", port)
	if err := server.ListenAndServe(); err != nil {
		panic(err)
	}
}
