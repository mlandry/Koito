package handlers

import (
	"fmt"
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

		l.Debug().Msgf("GenerateApiKeyHandler: Received request with params: '%s'", r.URL.Query().Encode())

		user := middleware.GetUserFromContext(ctx)
		if user == nil {
			l.Debug().Msg("GenerateApiKeyHandler: Invalid user retrieved from context")
			utils.WriteError(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		r.ParseForm()
		label := r.FormValue("label")
		if label == "" {
			l.Debug().Msg("GenerateApiKeyHandler: Request rejected due to missing label")
			utils.WriteError(w, "label is required", http.StatusBadRequest)
			return
		}

		apiKey, err := utils.GenerateRandomString(48)
		if err != nil {
			l.Err(fmt.Errorf("GenerateApiKeyHandler: %w", err)).Msg("Failed to generate API key")
			utils.WriteError(w, "failed to generate api key", http.StatusInternalServerError)
			return
		}

		opts := db.SaveApiKeyOpts{
			UserID: user.ID,
			Key:    apiKey,
			Label:  label,
		}
		l.Debug().Msgf("GenerateApiKeyHandler: Saving API key with options: %+v", opts)

		key, err := store.SaveApiKey(ctx, opts)
		if err != nil {
			l.Err(fmt.Errorf("GenerateApiKeyHandler: %w", err)).Msg("Failed to save API key")
			utils.WriteError(w, "failed to save api key", http.StatusInternalServerError)
			return
		}

		l.Debug().Msgf("GenerateApiKeyHandler: Successfully saved API key with ID: %d", key.ID)
		utils.WriteJSON(w, http.StatusCreated, key)
	}
}

func DeleteApiKeyHandler(store db.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		l := logger.FromContext(ctx)

		l.Debug().Msgf("DeleteApiKeyHandler: Received request with params: '%s'", r.URL.Query().Encode())

		user := middleware.GetUserFromContext(ctx)
		if user == nil {
			l.Debug().Msg("DeleteApiKeyHandler: User could not be verified (context user is nil)")
			utils.WriteError(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		idStr := r.URL.Query().Get("id")
		if idStr == "" {
			l.Debug().Msg("DeleteApiKeyHandler: Request rejected due to missing ID")
			utils.WriteError(w, "id is required", http.StatusBadRequest)
			return
		}

		apiKey, err := strconv.Atoi(idStr)
		if err != nil {
			l.Debug().AnErr("error", fmt.Errorf("DeleteApiKeyHandler: %w", err)).Msg("Invalid API key ID")
			utils.WriteError(w, "id is invalid", http.StatusBadRequest)
			return
		}

		l.Debug().Msgf("DeleteApiKeyHandler: Deleting API key with ID: %d", apiKey)

		err = store.DeleteApiKey(ctx, int32(apiKey))
		if err != nil {
			l.Err(fmt.Errorf("DeleteApiKeyHandler: %w", err)).Msg("Failed to delete API key")
			utils.WriteError(w, "failed to delete api key", http.StatusInternalServerError)
			return
		}

		l.Debug().Msgf("DeleteApiKeyHandler: Successfully deleted API key with ID: %d", apiKey)
		w.WriteHeader(http.StatusNoContent)
	}
}

func GetApiKeysHandler(store db.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		l := logger.FromContext(ctx)

		l.Debug().Msgf("GetApiKeysHandler: Received request with params: '%s'", r.URL.Query().Encode())

		user := middleware.GetUserFromContext(ctx)
		if user == nil {
			l.Debug().Msg("GetApiKeysHandler: Invalid user retrieved from context")
			utils.WriteError(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		l.Debug().Msgf("GetApiKeysHandler: Retrieving API keys for user ID: %d", user.ID)

		apiKeys, err := store.GetApiKeysByUserID(ctx, user.ID)
		if err != nil {
			l.Err(fmt.Errorf("GetApiKeysHandler: %w", err)).Msg("Failed to retrieve API keys")
			utils.WriteError(w, "failed to retrieve api keys", http.StatusInternalServerError)
			return
		}

		l.Debug().Msgf("GetApiKeysHandler: Successfully retrieved %d API keys for user ID: %d", len(apiKeys), user.ID)
		utils.WriteJSON(w, http.StatusOK, apiKeys)
	}
}

func UpdateApiKeyLabelHandler(store db.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		l := logger.FromContext(ctx)

		l.Debug().Msg("UpdateApiKeyLabelHandler: Received request to update API key label")

		user := middleware.GetUserFromContext(ctx)
		if user == nil {
			l.Debug().Msg("UpdateApiKeyLabelHandler: Unauthorized request (user context is nil)")
			utils.WriteError(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		idStr := r.URL.Query().Get("id")
		if idStr == "" {
			l.Debug().Msg("UpdateApiKeyLabelHandler: Missing API key ID in request")
			utils.WriteError(w, "id is required", http.StatusBadRequest)
			return
		}

		apiKeyID, err := strconv.Atoi(idStr)
		if err != nil {
			l.Debug().AnErr("error", fmt.Errorf("UpdateApiKeyLabelHandler: %w", err)).Msg("Invalid API key ID")
			utils.WriteError(w, "id is invalid", http.StatusBadRequest)
			return
		}

		label := r.FormValue("label")
		if label == "" {
			l.Debug().Msg("UpdateApiKeyLabelHandler: Missing label in request")
			utils.WriteError(w, "label is required", http.StatusBadRequest)
			return
		}

		l.Debug().Msgf("UpdateApiKeyLabelHandler: Updating label for API key ID %d", apiKeyID)

		err = store.UpdateApiKeyLabel(ctx, db.UpdateApiKeyLabelOpts{
			UserID: user.ID,
			ID:     int32(apiKeyID),
			Label:  label,
		})
		if err != nil {
			l.Err(fmt.Errorf("UpdateApiKeyLabelHandler: %w", err)).Msg("Failed to update API key label")
			utils.WriteError(w, "failed to update api key label", http.StatusInternalServerError)
			return
		}

		l.Debug().Msgf("UpdateApiKeyLabelHandler: Successfully updated label for API key ID %d", apiKeyID)
		w.WriteHeader(http.StatusOK)
	}
}
