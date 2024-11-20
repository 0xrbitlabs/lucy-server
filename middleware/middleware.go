package middleware

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"net/http"

	"github.com/joseph0x45/lucy/repository"
	"github.com/joseph0x45/lucy/utils"
)

type AuthMiddleware struct {
	users    *repository.UserRepo
	sessions *repository.SessionRepo
	logger   *slog.Logger
}

func NewAuthMiddleware(
	users *repository.UserRepo,
	sessions *repository.SessionRepo,
	logger *slog.Logger,
) *AuthMiddleware {
	return &AuthMiddleware{
		users:    users,
		sessions: sessions,
		logger:   logger,
	}
}

func (m *AuthMiddleware) Authenticate(roles ...string) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			session := r.Header.Get("Authorization")
			if session == "" {
				utils.WriteError(w, http.StatusUnauthorized, nil)
				return
			}
			dbSession, err := m.sessions.GetSessionByID(session)
			if err != nil && !errors.Is(err, sql.ErrNoRows) {
				m.logger.Error(err.Error())
				utils.WriteError(w, http.StatusUnauthorized, nil)
				return
			}
			if dbSession == nil {
				utils.WriteError(w, http.StatusUnauthorized, nil)
				return
			}
			if !dbSession.Valid {
				utils.WriteError(w, http.StatusUnauthorized, nil)
				return
			}
			user, err := m.users.GetByID(dbSession.UserID)
			if err != nil {
				if errors.Is(err, sql.ErrNoRows) {
					utils.WriteError(w, http.StatusUnauthorized, nil)
					return
				}
				m.logger.Error(err.Error())
				utils.WriteError(w, http.StatusInternalServerError, nil)
				return
			}
			ctx := context.WithValue(r.Context(), "user", user)
			h.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
