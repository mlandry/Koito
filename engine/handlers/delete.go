package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gabehf/koito/internal/db"
	"github.com/gabehf/koito/internal/logger"
	"github.com/gabehf/koito/internal/utils"
)

// DeleteTrackHandler deletes a track by its ID.
func DeleteTrackHandler(store db.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		l := logger.FromContext(ctx)

		trackIDStr := r.URL.Query().Get("id")
		if trackIDStr == "" {
			utils.WriteError(w, "track_id must be provided", http.StatusBadRequest)
			return
		}

		trackID, err := strconv.Atoi(trackIDStr)
		if err != nil {
			utils.WriteError(w, "invalid id", http.StatusBadRequest)
			return
		}

		err = store.DeleteTrack(ctx, int32(trackID))
		if err != nil {
			l.Err(err).Msg("Failed to delete track")
			utils.WriteError(w, "failed to delete track", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

// DeleteTrackHandler deletes a track by its ID.
func DeleteListenHandler(store db.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		l := logger.FromContext(ctx)

		trackIDStr := r.URL.Query().Get("track_id")
		if trackIDStr == "" {
			utils.WriteError(w, "track_id must be provided", http.StatusBadRequest)
			return
		}
		trackID, err := strconv.Atoi(trackIDStr)
		if err != nil {
			utils.WriteError(w, "invalid id", http.StatusBadRequest)
			return
		}

		unixStr := r.URL.Query().Get("unix")
		if trackIDStr == "" {
			utils.WriteError(w, "unix timestamp must be provided", http.StatusBadRequest)
			return
		}
		unix, err := strconv.ParseInt(unixStr, 10, 64)
		if err != nil {
			utils.WriteError(w, "invalid unix timestamp", http.StatusBadRequest)
			return
		}

		err = store.DeleteListen(ctx, int32(trackID), time.Unix(unix, 0))
		if err != nil {
			l.Err(err).Msg("Failed to delete listen")
			utils.WriteError(w, "failed to delete listen", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

// DeleteArtistHandler deletes an artist by its ID.
func DeleteArtistHandler(store db.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		l := logger.FromContext(ctx)

		artistIDStr := r.URL.Query().Get("id")
		if artistIDStr == "" {
			utils.WriteError(w, "id must be provided", http.StatusBadRequest)
			return
		}

		artistID, err := strconv.Atoi(artistIDStr)
		if err != nil {
			utils.WriteError(w, "invalid id", http.StatusBadRequest)
			return
		}

		err = store.DeleteArtist(ctx, int32(artistID))
		if err != nil {
			l.Err(err).Msg("Failed to delete artist")
			utils.WriteError(w, "failed to delete artist", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

// DeleteAlbumHandler deletes an album by its ID.
func DeleteAlbumHandler(store db.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		l := logger.FromContext(ctx)

		albumIDStr := r.URL.Query().Get("id")
		if albumIDStr == "" {
			utils.WriteError(w, "id must be provided", http.StatusBadRequest)
			return
		}

		albumID, err := strconv.Atoi(albumIDStr)
		if err != nil {
			utils.WriteError(w, "invalid id", http.StatusBadRequest)
			return
		}

		err = store.DeleteAlbum(ctx, int32(albumID))
		if err != nil {
			l.Err(err).Msg("Failed to delete album")
			utils.WriteError(w, "failed to delete album", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
