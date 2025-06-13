package handlers

import (
	"net/http"

	"github.com/gabehf/koito/internal/db"
	"github.com/gabehf/koito/internal/logger"
	"github.com/gabehf/koito/internal/models"
	"github.com/gabehf/koito/internal/utils"
)

type SearchResults struct {
	Artists []*models.Artist `json:"artists"`
	Albums  []*models.Album  `json:"albums"`
	Tracks  []*models.Track  `json:"tracks"`
}

func SearchHandler(store db.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		l := logger.FromContext(ctx)
		q := r.URL.Query().Get("q")
		artists, err := store.SearchArtists(ctx, q)

		l.Debug().Msgf("SearchHandler: Received search with query: %s", r.URL.Query().Encode())

		if err != nil {
			l.Err(err).Msg("Failed to search for artists")
			utils.WriteError(w, "failed to search in database", http.StatusInternalServerError)
			return
		}
		albums, err := store.SearchAlbums(ctx, q)
		if err != nil {
			l.Err(err).Msg("Failed to search for albums")
			utils.WriteError(w, "failed to search in database", http.StatusInternalServerError)
			return
		}
		tracks, err := store.SearchTracks(ctx, q)
		if err != nil {
			l.Err(err).Msg("Failed to search for tracks")
			utils.WriteError(w, "failed to search in database", http.StatusInternalServerError)
			return
		}
		utils.WriteJSON(w, http.StatusOK, SearchResults{
			Artists: artists,
			Albums:  albums,
			Tracks:  tracks,
		})
	}
}
