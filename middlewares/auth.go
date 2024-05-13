package middlewares

import (
	"context"
	"lucy/handlers"
	"lucy/providers"
	"lucy/repositories"
	"lucy/services"
	"lucy/types"
	"net/http"
)

type AuthMiddleware struct {
	userRepo services.UserRepo
	jwt      handlers.JWTProvider
	logger   handlers.Logger
}

func NewAuthMiddleware(
	userRepo services.UserRepo,
	jwt providers.JWTProvider,
	logger handlers.Logger,
) AuthMiddleware {
	return AuthMiddleware{
		userRepo: userRepo,
		jwt:      jwt,
		logger:   logger,
	}
}

func (m AuthMiddleware) AllowAccounts(accountTypes ...types.AccountType) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token := r.Header.Get("Authorization")
			if token == "" {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			claims, err := m.jwt.Decode(token)
			if err != nil {
				m.logger.Error(err.Error())
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			id, ok := claims["id"].(string)
			if !ok {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			user, err := m.userRepo.GetUser(repositories.Filter{Field: "id", Value: id})
			if err != nil {
				m.logger.Error(err.Error())
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			isNotAllowed := true
			for _, accountType := range accountTypes {
				if accountType == types.AnyAccount {
					isNotAllowed = false
					break
				}
				if user.AccountType == accountType {
					isNotAllowed = false
					break
				}
			}
			if isNotAllowed {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			ctx := context.WithValue(r.Context(), "user", user)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
