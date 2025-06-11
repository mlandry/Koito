package handlers

import (
	"net/http"
	"strings"
	"time"

	"github.com/gabehf/koito/engine/middleware"
	"github.com/gabehf/koito/internal/db"
	"github.com/gabehf/koito/internal/logger"
	"github.com/gabehf/koito/internal/utils"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func LoginHandler(store db.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := logger.FromContext(r.Context())
		ctx := r.Context()
		l.Debug().Msg("Recieved login request")

		r.ParseForm()
		username := r.FormValue("username")
		password := r.FormValue("password")
		if username == "" || password == "" {
			utils.WriteError(w, "username and password are required", http.StatusBadRequest)
			return
		}

		user, err := store.GetUserByUsername(ctx, username)
		if err != nil {
			l.Err(err).Msg("Error searching for user in database")
			utils.WriteError(w, "internal server error", http.StatusInternalServerError)
			return
		} else if user == nil {
			utils.WriteError(w, "username or password is incorrect", http.StatusBadRequest)
			return
		}

		err = bcrypt.CompareHashAndPassword(user.Password, []byte(password))
		if err != nil {
			utils.WriteError(w, "username or password is incorrect", http.StatusBadRequest)
			return
		}

		keepSignedIn := false
		expiresAt := time.Now().Add(1 * 24 * time.Hour)
		if strings.ToLower(r.FormValue("remember_me")) == "true" {
			keepSignedIn = true
			expiresAt = time.Now().Add(30 * 24 * time.Hour)
		}

		session, err := store.SaveSession(ctx, user.ID, expiresAt, keepSignedIn)
		if err != nil {
			l.Err(err).Msg("Failed to create session")
			utils.WriteError(w, "failed to create session", http.StatusInternalServerError)
			return
		}

		cookie := &http.Cookie{
			Name:     "koito_session",
			Value:    session.ID.String(),
			Path:     "/",
			HttpOnly: true,
			Secure:   false,
		}

		if keepSignedIn {
			cookie.Expires = expiresAt
		}

		http.SetCookie(w, cookie)
		w.WriteHeader(http.StatusNoContent)
	}
}

func LogoutHandler(store db.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := logger.FromContext(r.Context())
		cookie, err := r.Cookie("koito_session")
		if err == nil {
			sid, err := uuid.Parse(cookie.Value)
			if err != nil {
				utils.WriteError(w, "session cookie is invalid", http.StatusUnauthorized)
				return
			}
			err = store.DeleteSession(r.Context(), sid)
			if err != nil {
				l.Err(err).Msg("Failed to delete session")
				utils.WriteError(w, "internal server error", http.StatusInternalServerError)
				return
			}
		}

		// Clear the cookie
		http.SetCookie(w, &http.Cookie{
			Name:     "koito_session",
			Value:    "",
			Path:     "/",
			HttpOnly: true,
			MaxAge:   -1, // expire immediately
		})

		w.WriteHeader(http.StatusNoContent)
	}
}

func MeHandler(store db.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		l := logger.FromContext(ctx)
		u := middleware.GetUserFromContext(ctx)
		if u == nil {
			l.Debug().Msg("Invalid user retrieved from context")
		}
		utils.WriteJSON(w, 200, u)
	}
}

func UpdateUserHandler(store db.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		l := logger.FromContext(ctx)

		u := middleware.GetUserFromContext(ctx)
		if u == nil {
			utils.WriteError(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		r.ParseForm()
		username := r.FormValue("username")
		password := r.FormValue("password")

		l.Debug().Msgf("Recieved update request for user with id %d", u.ID)

		err := store.UpdateUser(ctx, db.UpdateUserOpts{
			ID:       u.ID,
			Username: username,
			Password: password,
		})
		if err != nil {
			l.Err(err).Msg("Failed to update user")
			utils.WriteError(w, err.Error(), http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}
