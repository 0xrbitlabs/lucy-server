package server

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"lucy/store"
	"lucy/types"
	"net/http"
	"os"

	"github.com/go-chi/jwtauth/v5"
)

type Server struct {
	Router *http.ServeMux
	Store  *store.Store
	logger *slog.Logger
	jwt    *jwtauth.JWTAuth
}

func NewServer(
	logger *slog.Logger,
	jwt *jwtauth.JWTAuth,
) *Server {
	store := store.NewStore()
	server := &Server{
		logger: logger,
		Store:  store,
		jwt:    jwt,
	}
	server.registerRoutes()
	return server
}

func (s *Server) Verify() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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
}

func (s *Server) Handle() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		message := new(types.InboundMessage)
		err := json.NewDecoder(r.Body).Decode(message)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		fmt.Println(message)
		w.WriteHeader(http.StatusOK)
		return
	}
}
