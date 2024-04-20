package server

import (
	"context"
	"errors"
	"lucy/app_errors"
	"net/http"
	"strings"
)

const (
	TypeSuperAdmin = "super_admin"
	TypeAdmin      = "admin"
	TypeSeller     = "seller"
	TypeRegular    = "regular"
)

func (s *Server) Allow(allowedTypes ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
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
			user, err := s.Store.GetUserByID(id)
			if err != nil {
				if errors.Is(err, app_errors.ErrResourceNotFound) {
					w.WriteHeader(http.StatusUnauthorized)
					return
				}
				s.logger.Error(err.Error())
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			isAllowed := false
			for _, t := range allowedTypes {
				if user.Type == t {
					isAllowed = true
					break
				}
			}
			if !isAllowed {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			ctx := context.WithValue(r.Context(), "user", user)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
