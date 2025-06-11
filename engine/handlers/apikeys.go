package handlers

import (
	"net/http"
	"strconv"

	"github.com/gabehf/koito/engine/middleware"
	"github.com/gabehf/koito/internal/db"
	"github.com/gabehf/koito/internal/logger"
	"github.com/gabehf/koito/internal/utils"
)

func GenerateApiKeyHandler(store db.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		l := logger.FromContext(ctx)

		user := middleware.GetUserFromContext(ctx)
		if user == nil {
			utils.WriteError(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		r.ParseForm()
		label := r.FormValue("label")
		if label == "" {
			utils.WriteError(w, "label is required", http.StatusBadRequest)
			return
		}

		apiKey, err := utils.GenerateRandomString(48)
		if err != nil {
			l.Err(err).Msg("Failed to generate API key")
			utils.WriteError(w, "failed to generate api key", http.StatusInternalServerError)
			return
		}
		opts := db.SaveApiKeyOpts{
			UserID: user.ID,
			Key:    apiKey,
			Label:  label,
		}
		l.Debug().Any("opts", opts).Send()
		key, err := store.SaveApiKey(ctx, opts)
		if err != nil {
			l.Err(err).Msg("Failed to save API key")
			utils.WriteError(w, "failed to save api key", http.StatusInternalServerError)
			return
		}
		utils.WriteJSON(w, 201, key)
	}
}

func DeleteApiKeyHandler(store db.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		l := logger.FromContext(ctx)

		user := middleware.GetUserFromContext(ctx)
		if user == nil {
			utils.WriteError(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		idStr := r.URL.Query().Get("id")
		if idStr == "" {
			utils.WriteError(w, "id is required", http.StatusBadRequest)
			return
		}
		apiKey, err := strconv.Atoi(idStr)
		if err != nil {
			utils.WriteError(w, "id is invalid", http.StatusBadRequest)
			return
		}

		err = store.DeleteApiKey(ctx, int32(apiKey))
		if err != nil {
			l.Err(err).Msg("Failed to delete API key")
			utils.WriteError(w, "failed to delete api key", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

func GetApiKeysHandler(store db.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		l := logger.FromContext(ctx)

		l.Debug().Msgf("Retrieving user from middleware...")

		user := middleware.GetUserFromContext(ctx)
		if user == nil {
			l.Debug().Msgf("Could not retrieve user from middleware")
			utils.WriteError(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		l.Debug().Msgf("Retrieved user '%s' from middleware", user.Username)

		apiKeys, err := store.GetApiKeysByUserID(ctx, user.ID)
		if err != nil {
			l.Err(err).Msg("Failed to retrieve API keys")
			utils.WriteError(w, "failed to retrieve api keys", http.StatusInternalServerError)
			return
		}

		utils.WriteJSON(w, http.StatusOK, apiKeys)
	}
}

func UpdateApiKeyLabelHandler(store db.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		l := logger.FromContext(ctx)

		user := middleware.GetUserFromContext(ctx)
		if user == nil {
			utils.WriteError(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		idStr := r.URL.Query().Get("id")
		if idStr == "" {
			utils.WriteError(w, "id is required", http.StatusBadRequest)
			return
		}
		apiKeyID, err := strconv.Atoi(idStr)
		if err != nil {
			utils.WriteError(w, "id is invalid", http.StatusBadRequest)
			return
		}

		label := r.FormValue("label")
		if label == "" {
			utils.WriteError(w, "label is required", http.StatusBadRequest)
			return
		}

		err = store.UpdateApiKeyLabel(ctx, db.UpdateApiKeyLabelOpts{
			UserID: user.ID,
			ID:     int32(apiKeyID),
			Label:  label,
		})
		if err != nil {
			l.Err(err).Msg("Failed to update API key label")
			utils.WriteError(w, "failed to update api key label", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
