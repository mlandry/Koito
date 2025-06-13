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
		ctx := r.Context()
		l := logger.FromContext(ctx)

		l.Debug().Msg("LoginHandler: Received login request")

		r.ParseForm()
		username := r.FormValue("username")
		password := r.FormValue("password")
		if username == "" || password == "" {
			l.Debug().Msg("LoginHandler: Missing username or password")
			utils.WriteError(w, "username and password are required", http.StatusBadRequest)
			return
		}

		l.Debug().Msgf("LoginHandler: Searching for user with username '%s'", username)
		user, err := store.GetUserByUsername(ctx, username)
		if err != nil {
			l.Err(err).Msg("LoginHandler: Error searching for user in database")
			utils.WriteError(w, "internal server error", http.StatusInternalServerError)
			return
		} else if user == nil {
			l.Debug().Msg("LoginHandler: Username or password is incorrect")
			utils.WriteError(w, "username or password is incorrect", http.StatusBadRequest)
			return
		}

		err = bcrypt.CompareHashAndPassword(user.Password, []byte(password))
		if err != nil {
			l.Debug().Msg("LoginHandler: Password comparison failed")
			utils.WriteError(w, "username or password is incorrect", http.StatusBadRequest)
			return
		}

		keepSignedIn := false
		expiresAt := time.Now().Add(1 * 24 * time.Hour)
		if strings.ToLower(r.FormValue("remember_me")) == "true" {
			keepSignedIn = true
			expiresAt = time.Now().Add(30 * 24 * time.Hour)
		}

		l.Debug().Msgf("LoginHandler: Creating session for user ID %d", user.ID)
		session, err := store.SaveSession(ctx, user.ID, expiresAt, keepSignedIn)
		if err != nil {
			l.Err(err).Msg("LoginHandler: Failed to create session")
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

		l.Debug().Msgf("LoginHandler: Session created successfully for user ID %d", user.ID)
		http.SetCookie(w, cookie)
		w.WriteHeader(http.StatusNoContent)
	}
}

func LogoutHandler(store db.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		l := logger.FromContext(ctx)

		l.Debug().Msg("LogoutHandler: Received logout request")
		cookie, err := r.Cookie("koito_session")
		if err == nil {
			l.Debug().Msg("LogoutHandler: Found session cookie")
			sid, err := uuid.Parse(cookie.Value)
			if err != nil {
				l.Debug().AnErr("error", err).Msg("LogoutHandler: Invalid session cookie")
				utils.WriteError(w, "session cookie is invalid", http.StatusUnauthorized)
				return
			}
			l.Debug().Msgf("LogoutHandler: Deleting session with ID %s", sid)
			err = store.DeleteSession(ctx, sid)
			if err != nil {
				l.Err(err).Msg("LogoutHandler: Failed to delete session")
				utils.WriteError(w, "internal server error", http.StatusInternalServerError)
				return
			}
		}

		l.Debug().Msg("LogoutHandler: Clearing session cookie")
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

		l.Debug().Msg("MeHandler: Received request to retrieve user information")
		u := middleware.GetUserFromContext(ctx)
		if u == nil {
			l.Debug().Msg("MeHandler: Invalid user retrieved from context")
			utils.WriteError(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		l.Debug().Msgf("MeHandler: Successfully retrieved user with ID %d", u.ID)
		utils.WriteJSON(w, http.StatusOK, u)
	}
}

func UpdateUserHandler(store db.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		l := logger.FromContext(ctx)

		l.Debug().Msg("UpdateUserHandler: Received request to update user information")
		u := middleware.GetUserFromContext(ctx)
		if u == nil {
			l.Debug().Msg("UpdateUserHandler: Unauthorized request (user context is nil)")
			utils.WriteError(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		r.ParseForm()
		username := r.FormValue("username")
		password := r.FormValue("password")

		l.Debug().Msgf("UpdateUserHandler: Updating user with ID %d", u.ID)
		err := store.UpdateUser(ctx, db.UpdateUserOpts{
			ID:       u.ID,
			Username: username,
			Password: password,
		})
		if err != nil {
			l.Err(err).Msg("UpdateUserHandler: Failed to update user")
			utils.WriteError(w, err.Error(), http.StatusBadRequest)
			return
		}

		l.Debug().Msgf("UpdateUserHandler: Successfully updated user with ID %d", u.ID)
		w.WriteHeader(http.StatusNoContent)
	}
}
