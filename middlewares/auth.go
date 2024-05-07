package middlewares

import (
	"context"
	"lucy/interfaces"
	"lucy/providers"
	"lucy/repositories"
	"lucy/types"
	"net/http"
)

type AuthMiddleware struct {
	userRepo interfaces.UserRepo
	jwt      providers.JWTProvider
	logger   interfaces.Logger
}

func NewAuthMiddleware(
	userRepo interfaces.UserRepo,
	jwt providers.JWTProvider,
	logger interfaces.Logger,
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
				if user.AccountType == accountType {
					isNotAllowed = false
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
