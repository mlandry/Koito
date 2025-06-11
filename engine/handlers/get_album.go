package handlers

import (
	"net/http"
	"strconv"

	"github.com/gabehf/koito/internal/db"
	"github.com/gabehf/koito/internal/utils"
)

func GetAlbumHandler(store db.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		idStr := r.URL.Query().Get("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			utils.WriteError(w, "id is invalid", 400)
			return
		}

		album, err := store.GetAlbum(r.Context(), db.GetAlbumOpts{ID: int32(id)})
		if err != nil {
			utils.WriteError(w, "album with specified id could not be found", http.StatusNotFound)
			return
		}
		utils.WriteJSON(w, http.StatusOK, album)
	}
}
