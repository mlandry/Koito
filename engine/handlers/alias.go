package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gabehf/koito/internal/db"
	"github.com/gabehf/koito/internal/logger"
	"github.com/gabehf/koito/internal/models"
	"github.com/gabehf/koito/internal/utils"
)

// GetAliasesHandler retrieves all aliases for a given artist or album ID.
func GetAliasesHandler(store db.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		l := logger.FromContext(ctx)

		l.Debug().Msgf("GetAliasesHandler: Got request with params: '%s'", r.URL.Query().Encode())

		// Parse query parameters
		artistIDStr := r.URL.Query().Get("artist_id")
		albumIDStr := r.URL.Query().Get("album_id")
		trackIDStr := r.URL.Query().Get("track_id")

		if artistIDStr == "" && albumIDStr == "" && trackIDStr == "" {
			l.Debug().Msgf("Request is missing required parameters")
			utils.WriteError(w, "artist_id, album_id, or track_id must be provided", http.StatusBadRequest)
			return
		}
		if utils.MoreThanOneString(artistIDStr, albumIDStr, trackIDStr) {
			l.Debug().Msgf("Request is has more than one of artist_id, album_id, and track_id")
			utils.WriteError(w, "only one of artist_id, album_id, or track_id can be provided at a time", http.StatusBadRequest)
			return
		}

		var aliases []models.Alias

		if artistIDStr != "" {
			artistID, err := strconv.Atoi(artistIDStr)
			if err != nil {
				l.Debug().AnErr("error", fmt.Errorf("GetAliasesHandler: %w", err)).Msg("Invalid artist id")
				utils.WriteError(w, "invalid artist_id", http.StatusBadRequest)
				return
			}
			aliases, err = store.GetAllArtistAliases(ctx, int32(artistID))
			if err != nil {
				l.Err(fmt.Errorf("GetAliasesHandler: %w", err)).Msg("Failed to get artist aliases")
				utils.WriteError(w, "failed to retrieve aliases", http.StatusInternalServerError)
				return
			}
		} else if albumIDStr != "" {
			albumID, err := strconv.Atoi(albumIDStr)
			if err != nil {
				l.Debug().AnErr("error", fmt.Errorf("GetAliasesHandler: %w", err)).Msg("Invalid album id")
				utils.WriteError(w, "invalid album_id", http.StatusBadRequest)
				return
			}
			aliases, err = store.GetAllAlbumAliases(ctx, int32(albumID))
			if err != nil {
				l.Err(fmt.Errorf("GetAliasesHandler: %w", err)).Msg("Failed to get album aliases")
				utils.WriteError(w, "failed to retrieve aliases", http.StatusInternalServerError)
				return
			}
		} else if trackIDStr != "" {
			trackID, err := strconv.Atoi(trackIDStr)
			if err != nil {
				l.Debug().AnErr("error", fmt.Errorf("GetAliasesHandler: %w", err)).Msg("Invalid track id")
				utils.WriteError(w, "invalid track_id", http.StatusBadRequest)
				return
			}
			aliases, err = store.GetAllTrackAliases(ctx, int32(trackID))
			if err != nil {
				l.Err(fmt.Errorf("GetAliasesHandler: %w", err)).Msg("Failed to get track aliases")
				utils.WriteError(w, "failed to retrieve aliases", http.StatusInternalServerError)
				return
			}
		}

		utils.WriteJSON(w, http.StatusOK, aliases)
	}
}

// DeleteAliasHandler deletes an alias for a given artist or album ID.
func DeleteAliasHandler(store db.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		l := logger.FromContext(ctx)

		l.Debug().Msgf("DeleteAliasHandler: Got request with params: '%s'", r.URL.Query().Encode())

		// Parse query parameters
		artistIDStr := r.URL.Query().Get("artist_id")
		albumIDStr := r.URL.Query().Get("album_id")
		trackIDStr := r.URL.Query().Get("track_id")
		alias := r.URL.Query().Get("alias")

		if alias == "" || (artistIDStr == "" && albumIDStr == "" && trackIDStr == "") {
			l.Debug().Msgf("Request is missing required parameters")
			utils.WriteError(w, "alias and artist_id, album_id, or track_id must be provided", http.StatusBadRequest)
			return
		}
		if utils.MoreThanOneString(artistIDStr, albumIDStr, trackIDStr) {
			l.Debug().Msgf("Request is has more than one of artist_id, album_id, and track_id")
			utils.WriteError(w, "only one of artist_id, album_id, or track_id can be provided at a time", http.StatusBadRequest)
			return
		}

		if artistIDStr != "" {
			artistID, err := strconv.Atoi(artistIDStr)
			if err != nil {
				l.Debug().AnErr("error", fmt.Errorf("DeleteAliasHandler: %w", err)).Msg("Invalid artist id")
				utils.WriteError(w, "invalid artist_id", http.StatusBadRequest)
				return
			}
			err = store.DeleteArtistAlias(ctx, int32(artistID), alias)
			if err != nil {
				l.Err(fmt.Errorf("DeleteAliasHandler: %w", err)).Msg("Failed to delete artist alias")
				utils.WriteError(w, "failed to delete alias", http.StatusInternalServerError)
				return
			}
		} else if albumIDStr != "" {
			albumID, err := strconv.Atoi(albumIDStr)
			if err != nil {
				l.Debug().AnErr("error", fmt.Errorf("DeleteAliasHandler: %w", err)).Msg("Invalid album id")
				utils.WriteError(w, "invalid album_id", http.StatusBadRequest)
				return
			}
			err = store.DeleteAlbumAlias(ctx, int32(albumID), alias)
			if err != nil {
				l.Err(fmt.Errorf("DeleteAliasHandler: %w", err)).Msg("Failed to delete album alias")
				utils.WriteError(w, "failed to delete alias", http.StatusInternalServerError)
				return
			}
		} else if trackIDStr != "" {
			trackID, err := strconv.Atoi(trackIDStr)
			if err != nil {
				l.Debug().AnErr("error", fmt.Errorf("DeleteAliasHandler: %w", err)).Msg("Invalid track id")
				utils.WriteError(w, "invalid track_id", http.StatusBadRequest)
				return
			}
			err = store.DeleteTrackAlias(ctx, int32(trackID), alias)
			if err != nil {
				l.Err(fmt.Errorf("DeleteAliasHandler: %w", err)).Msg("Failed to delete track alias")
				utils.WriteError(w, "failed to delete alias", http.StatusInternalServerError)
				return
			}
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

// CreateAliasHandler creates new aliases for a given artist, album, or track.
func CreateAliasHandler(store db.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		l := logger.FromContext(ctx)

		l.Debug().Msgf("CreateAliasHandler: Got request with params: '%s'", r.URL.Query().Encode())

		err := r.ParseForm()
		if err != nil {
			utils.WriteError(w, "invalid request body", http.StatusBadRequest)
			return
		}

		alias := r.FormValue("alias")
		if alias == "" {
			utils.WriteError(w, "alias must be provided", http.StatusBadRequest)
			return
		}

		artistIDStr := r.URL.Query().Get("artist_id")
		albumIDStr := r.URL.Query().Get("album_id")
		trackIDStr := r.URL.Query().Get("track_id")

		if alias == "" || (artistIDStr == "" && albumIDStr == "" && trackIDStr == "") {
			l.Debug().Msgf("Request is missing required parameters")
			utils.WriteError(w, "alias and artist_id, album_id, or track_id must be provided", http.StatusBadRequest)
			return
		}
		if utils.MoreThanOneString(artistIDStr, albumIDStr, trackIDStr) {
			l.Debug().Msgf("Request is has more than one of artist_id, album_id, and track_id")
			utils.WriteError(w, "only one of artist_id, album_id, or track_id can be provided at a time", http.StatusBadRequest)
			return
		}

		if artistIDStr != "" {
			artistID, err := strconv.Atoi(artistIDStr)
			if err != nil {
				l.Debug().AnErr("error", fmt.Errorf("CreateAliasHandler: %w", err)).Msg("Invalid artist id")
				utils.WriteError(w, "invalid artist_id", http.StatusBadRequest)
				return
			}
			err = store.SaveArtistAliases(ctx, int32(artistID), []string{alias}, "Manual")
			if err != nil {
				l.Err(fmt.Errorf("CreateAliasHandler: %w", err)).Msg("Failed to save artist alias")
				utils.WriteError(w, "failed to save alias", http.StatusInternalServerError)
				return
			}
		} else if albumIDStr != "" {
			albumID, err := strconv.Atoi(albumIDStr)
			if err != nil {
				l.Debug().AnErr("error", fmt.Errorf("CreateAliasHandler: %w", err)).Msg("Invalid album id")
				utils.WriteError(w, "invalid album_id", http.StatusBadRequest)
				return
			}
			err = store.SaveAlbumAliases(ctx, int32(albumID), []string{alias}, "Manual")
			if err != nil {
				l.Err(fmt.Errorf("CreateAliasHandler: %w", err)).Msg("Failed to save album alias")
				utils.WriteError(w, "failed to save alias", http.StatusInternalServerError)
				return
			}
		} else if trackIDStr != "" {
			trackID, err := strconv.Atoi(trackIDStr)
			if err != nil {
				l.Debug().AnErr("error", fmt.Errorf("CreateAliasHandler: %w", err)).Msg("Invalid track id")
				utils.WriteError(w, "invalid track_id", http.StatusBadRequest)
				return
			}
			err = store.SaveTrackAliases(ctx, int32(trackID), []string{alias}, "Manual")
			if err != nil {
				l.Err(fmt.Errorf("CreateAliasHandler: %w", err)).Msg("Failed to save track alias")
				utils.WriteError(w, "failed to save alias", http.StatusInternalServerError)
				return
			}
		}

		w.WriteHeader(http.StatusCreated)
	}
}

// sets the primary alias for albums, artists, and tracks
func SetPrimaryAliasHandler(store db.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		l := logger.FromContext(ctx)

		l.Debug().Msgf("SetPrimaryAliasHandler: Got request with params: '%s'", r.URL.Query().Encode())

		// Parse query parameters
		artistIDStr := r.URL.Query().Get("artist_id")
		albumIDStr := r.URL.Query().Get("album_id")
		trackIDStr := r.URL.Query().Get("track_id")
		alias := r.URL.Query().Get("alias")

		if alias == "" || (artistIDStr == "" && albumIDStr == "" && trackIDStr == "") {
			l.Debug().Msgf("Request is missing required parameters")
			utils.WriteError(w, "alias and artist_id, album_id, or track_id must be provided", http.StatusBadRequest)
			return
		}
		if utils.MoreThanOneString(artistIDStr, albumIDStr, trackIDStr) {
			l.Debug().Msgf("Request is has more than one of artist_id, album_id, and track_id")
			utils.WriteError(w, "only one of artist_id, album_id, or track_id can be provided at a time", http.StatusBadRequest)
			return
		}

		if artistIDStr != "" {
			artistID, err := strconv.Atoi(artistIDStr)
			if err != nil {
				l.Debug().AnErr("error", fmt.Errorf("SetPrimaryAliasHandler: %w", err)).Msg("Invalid artist id")
				utils.WriteError(w, "invalid artist_id", http.StatusBadRequest)
				return
			}
			err = store.SetPrimaryArtistAlias(ctx, int32(artistID), alias)
			if err != nil {
				l.Err(fmt.Errorf("SetPrimaryAliasHandler: %w", err)).Msg("Failed to set artist primary alias")
				utils.WriteError(w, "failed to set primary alias", http.StatusInternalServerError)
				return
			}
		} else if albumIDStr != "" {
			albumID, err := strconv.Atoi(albumIDStr)
			if err != nil {
				l.Debug().AnErr("error", fmt.Errorf("SetPrimaryAliasHandler: %w", err)).Msg("Invalid album id")
				utils.WriteError(w, "invalid album_id", http.StatusBadRequest)
				return
			}
			err = store.SetPrimaryAlbumAlias(ctx, int32(albumID), alias)
			if err != nil {
				l.Err(fmt.Errorf("SetPrimaryAliasHandler: %w", err)).Msg("Failed to set album primary alias")
				utils.WriteError(w, "failed to set primary alias", http.StatusInternalServerError)
				return
			}
		} else if trackIDStr != "" {
			trackID, err := strconv.Atoi(trackIDStr)
			if err != nil {
				l.Debug().AnErr("error", fmt.Errorf("SetPrimaryAliasHandler: %w", err)).Msg("Invalid track id")
				utils.WriteError(w, "invalid track_id", http.StatusBadRequest)
				return
			}
			err = store.SetPrimaryTrackAlias(ctx, int32(trackID), alias)
			if err != nil {
				l.Err(fmt.Errorf("SetPrimaryAliasHandler: %w", err)).Msg("Failed to set track primary alias")
				utils.WriteError(w, "failed to set primary alias", http.StatusInternalServerError)
				return
			}
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
