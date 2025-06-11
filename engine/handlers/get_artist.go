package handlers

import (
	"net/http"
	"strconv"

	"github.com/gabehf/koito/internal/db"
	"github.com/gabehf/koito/internal/utils"
)

func GetArtistHandler(store db.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		idStr := r.URL.Query().Get("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			utils.WriteError(w, "id is invalid", 400)
			return
		}

		artist, err := store.GetArtist(r.Context(), db.GetArtistOpts{ID: int32(id)})
		if err != nil {
			utils.WriteError(w, "artist with specified id could not be found", http.StatusNotFound)
			return
		}
		utils.WriteJSON(w, http.StatusOK, artist)
	}
}
