package handlers

import (
	"net/http"

	"github.com/gabehf/koito/internal/db"
	"github.com/gabehf/koito/internal/logger"
	"github.com/gabehf/koito/internal/utils"
)

func GetListensHandler(store db.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		l := logger.FromContext(r.Context())
		opts := OptsFromRequest(r)
		listens, err := store.GetListensPaginated(r.Context(), opts)
		if err != nil {
			l.Err(err).Send()
			utils.WriteError(w, "failed to get listens: "+err.Error(), 400)
			return
		}
		utils.WriteJSON(w, http.StatusOK, listens)
	}
}
