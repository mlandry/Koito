package middleware

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/gabehf/koito/internal/db"
	"github.com/gabehf/koito/internal/logger"
	"github.com/gabehf/koito/internal/models"
	"github.com/gabehf/koito/internal/utils"
	"github.com/google/uuid"
)

type MiddlwareContextKey string

const (
	UserContextKey   MiddlwareContextKey = "user"
	apikeyContextKey MiddlwareContextKey = "apikeyID"
)

func ValidateSession(store db.DB) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			l := logger.FromContext(r.Context())
			cookie, err := r.Cookie("koito_session")
			var sid uuid.UUID
			if err == nil {
				sid, err = uuid.Parse(cookie.Value)
				if err != nil {
					utils.WriteError(w, "session cookie is invalid", http.StatusUnauthorized)
					return
				}
			}

			l.Debug().Msg("Retrieved login cookie from request")

			u, err := store.GetUserBySession(r.Context(), sid)
			if err != nil {
				l.Err(err).Msg("Failed to get user from session")
				utils.WriteError(w, "internal server error", http.StatusInternalServerError)
				return
			}
			if u == nil {
				utils.WriteError(w, "unauthorized", http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), UserContextKey, u)
			r = r.WithContext(ctx)

			l.Debug().Msgf("Refreshing session for user '%s'", u.Username)

			store.RefreshSession(r.Context(), sid, time.Now().Add(30*24*time.Hour))

			l.Debug().Msgf("Refreshed session for user '%s'", u.Username)

			next.ServeHTTP(w, r)
		})
	}
}

func ValidateApiKey(store db.DB) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			l := logger.FromContext(ctx)

			authh := r.Header.Get("Authorization")
			s := strings.Split(authh, "Token ")
			if len(s) < 2 {
				l.Debug().Msg("Authorization header must be formatted 'Token {token}'")
				utils.WriteError(w, "unauthorized", http.StatusUnauthorized)
				return
			}

			key := s[1]

			u, err := store.GetUserByApiKey(ctx, key)
			if err != nil {
				l.Err(err).Msg("Failed to get user from database using api key")
				utils.WriteError(w, "internal server error", http.StatusInternalServerError)
				return
			}
			if u == nil {
				l.Debug().Msg("Api key does not exist")
				utils.WriteError(w, "unauthorized", http.StatusUnauthorized)
				return
			}

			ctx = context.WithValue(r.Context(), UserContextKey, u)
			r = r.WithContext(ctx)

			next.ServeHTTP(w, r)
		})
	}
}

func GetUserFromContext(ctx context.Context) *models.User {
	user, ok := ctx.Value(UserContextKey).(*models.User)
	if !ok {
		return nil
	}
	return user
}
