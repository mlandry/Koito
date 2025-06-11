package handlers

import (
	"net/http"

	"github.com/gabehf/koito/internal/db"
	"github.com/gabehf/koito/internal/logger"
	"github.com/gabehf/koito/internal/utils"
)

func GetTopArtistsHandler(store db.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		l := logger.FromContext(r.Context())
		opts := OptsFromRequest(r)
		artists, err := store.GetTopArtistsPaginated(r.Context(), opts)
		if err != nil {
			l.Err(err).Msg("Failed to get top artists")
			utils.WriteError(w, "failed to get artists", 400)
			return
		}
		utils.WriteJSON(w, http.StatusOK, artists)
	}
}
