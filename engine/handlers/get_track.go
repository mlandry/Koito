package handlers

import (
	"net/http"
	"strconv"

	"github.com/gabehf/koito/internal/db"
	"github.com/gabehf/koito/internal/logger"
	"github.com/gabehf/koito/internal/utils"
)

func GetTrackHandler(store db.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		l := logger.FromContext(r.Context())

		idStr := r.URL.Query().Get("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			utils.WriteError(w, "id is invalid", 400)
			return
		}

		track, err := store.GetTrack(r.Context(), db.GetTrackOpts{ID: int32(id)})
		if err != nil {
			l.Err(err).Msg("Failed to get top albums")
			utils.WriteError(w, "track with specified id could not be found", http.StatusNotFound)
			return
		}
		utils.WriteJSON(w, http.StatusOK, track)
	}
}
