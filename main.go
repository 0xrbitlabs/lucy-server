package main

import (
	"fmt"
	"lucy/server"
	"net/http"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	server := server.NewServer()
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
