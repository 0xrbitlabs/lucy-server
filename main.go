package main

import (
	"encoding/json"
	"fmt"
	"github.com/joho/godotenv"
	"net/http"
	"os"
	"server/types"
)

func Verify(w http.ResponseWriter, r *http.Request) {
	verifyToken := os.Getenv("VERIFY_TOKEN")
	mode := r.URL.Query().Get("hub.mode")
	token := r.URL.Query().Get("hub.verify_token")
	challenge := r.URL.Query().Get("hub.challenge")
	if mode == "subscribe" && token == verifyToken {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(challenge))
		return
	}
	w.WriteHeader(http.StatusBadRequest)
	return
}

func Handle(w http.ResponseWriter, r *http.Request) {
	message := new(types.InboundMessage)
	err := json.NewDecoder(r.Body).Decode(message)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Printf("%+v", message)
	return
}

func main() {
	godotenv.Load()
	mux := http.NewServeMux()
	mux.HandleFunc("GET /hook", Verify)
	mux.HandleFunc("POST /hook", Handle)
	fmt.Println("Server started on port 9090")
	http.ListenAndServe(":9090", mux)
}
