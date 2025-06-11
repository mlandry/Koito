package handlers

import (
	"net/http"
	"strconv"

	"github.com/gabehf/koito/internal/db"
	"github.com/gabehf/koito/internal/logger"
	"github.com/gabehf/koito/internal/utils"
)

func MergeTracksHandler(store db.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := logger.FromContext(r.Context())

		fromidStr := r.URL.Query().Get("from_id")
		fromId, err := strconv.Atoi(fromidStr)
		if err != nil {
			l.Err(err).Send()
			utils.WriteError(w, "from_id is invalid", 400)
			return
		}
		toidStr := r.URL.Query().Get("to_id")
		toId, err := strconv.Atoi(toidStr)
		if err != nil {
			l.Err(err).Send()
			utils.WriteError(w, "to_id is invalid", 400)
			return
		}

		err = store.MergeTracks(r.Context(), int32(fromId), int32(toId))
		if err != nil {
			l.Err(err).Send()
			utils.WriteError(w, "Failed to merge tracks: "+err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}

func MergeReleaseGroupsHandler(store db.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := logger.FromContext(r.Context())

		fromidStr := r.URL.Query().Get("from_id")
		fromId, err := strconv.Atoi(fromidStr)
		if err != nil {
			l.Err(err).Send()
			utils.WriteError(w, "from_id is invalid", 400)
			return
		}
		toidStr := r.URL.Query().Get("to_id")
		toId, err := strconv.Atoi(toidStr)
		if err != nil {
			l.Err(err).Send()
			utils.WriteError(w, "to_id is invalid", 400)
			return
		}

		err = store.MergeAlbums(r.Context(), int32(fromId), int32(toId))
		if err != nil {
			l.Err(err).Send()
			utils.WriteError(w, "Failed to merge albums: "+err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}

func MergeArtistsHandler(store db.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := logger.FromContext(r.Context())

		fromidStr := r.URL.Query().Get("from_id")
		fromId, err := strconv.Atoi(fromidStr)
		if err != nil {
			l.Err(err).Send()
			utils.WriteError(w, "from_id is invalid", 400)
			return
		}
		toidStr := r.URL.Query().Get("to_id")
		toId, err := strconv.Atoi(toidStr)
		if err != nil {
			l.Err(err).Send()
			utils.WriteError(w, "to_id is invalid", 400)
			return
		}

		err = store.MergeArtists(r.Context(), int32(fromId), int32(toId))
		if err != nil {
			l.Err(err).Send()
			utils.WriteError(w, "Failed to merge artists: "+err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}
