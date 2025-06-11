package handlers

import (
	"net/http"

	"github.com/gabehf/koito/internal/db"
	"github.com/gabehf/koito/internal/logger"
	"github.com/gabehf/koito/internal/utils"
)

func GetTopAlbumsHandler(store db.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		l := logger.FromContext(r.Context())
		opts := OptsFromRequest(r)
		albums, err := store.GetTopAlbumsPaginated(r.Context(), opts)
		if err != nil {
			l.Err(err).Msg("Failed to get top albums")
			utils.WriteError(w, "failed to get albums", 400)
			return
		}
		utils.WriteJSON(w, http.StatusOK, albums)
	}
}
