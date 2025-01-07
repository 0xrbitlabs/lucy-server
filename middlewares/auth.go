package middlewares

import (
	"context"
	"log/slog"
	"lucy/repo"
	"net/http"
)

var roles = []string{"regular", "seller", "admin"}

type AuthMiddleware struct {
	users    *repo.UserRepo
	sessions *repo.SessionRepo
	logger   *slog.Logger
}

func NewAuthMiddleware(
	users *repo.UserRepo,
	sessions *repo.SessionRepo,
	logger *slog.Logger,
) *AuthMiddleware {
	return &AuthMiddleware{
		users:    users,
		sessions: sessions,
		logger:   logger,
	}
}

func userHasRole(userRole string) bool {
	ok := false
	for _, role := range roles {
		if userRole == role {
			ok = true
			break
		}
	}
	return ok
}

func (m *AuthMiddleware) AuthenticateWithRole(roles ...string) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			sessionID := r.Header.Get("Authorization")
			if sessionID == "" {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			session, err := m.sessions.GetSessionByID(sessionID)
			if err != nil {
				m.logger.Error(err.Error())
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			if session == nil {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			user, err := m.users.GetUserByID(session.UserID)
			if err != nil {
				m.logger.Error(err.Error())
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			if user == nil {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			if len(roles) > 0 && !userHasRole(user.AccountType) {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			ctx := context.WithValue(r.Context(), "user", user)
			h.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
