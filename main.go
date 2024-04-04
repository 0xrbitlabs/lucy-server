package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"log/slog"
	"net/http"
	"os"
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
	mux := http.NewServeMux()
	mux.Handle("POST /", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	}))
	fmt.Println("Server launched on port 8081")
	err := http.ListenAndServe(":8081", mux)
	if err != nil {
		panic(err)
	}
}
