package server

import (
	"context"
	"errors"
	"lucy/app_errors"
	"lucy/store"
	"net/http"
	"strings"
)

func (s *Server) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		token, err := s.jwt.Decode(parts[1])
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		val, ok := token.Get("id")
		if !ok {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		id, ok := val.(string)
		if !ok {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		val, ok = token.Get("isAdmin")
		if !ok {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		userIsAdmin, ok := val.(bool)
		if !ok {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		var currentUser interface{}
		if userIsAdmin {
			admin, err := s.Store.GetAdmin(store.GetAdminFilter{Column: "id", Value: id})
			if err != nil {
				if errors.Is(err, app_errors.ErrResourceNotFound) {
					w.WriteHeader(http.StatusUnauthorized)
					return
				}
				s.logger.Error(err.Error())
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			currentUser = admin
		} else {
			//fetch user info
			//currentUser = user
		}
		ctx := context.WithValue(r.Context(), "user", currentUser)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
