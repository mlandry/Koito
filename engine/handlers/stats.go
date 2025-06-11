package handlers

import (
	"net/http"
	"strings"

	"github.com/gabehf/koito/internal/db"
	"github.com/gabehf/koito/internal/logger"
	"github.com/gabehf/koito/internal/utils"
)

type StatsResponse struct {
	ListenCount   int64 `json:"listen_count"`
	TrackCount    int64 `json:"track_count"`
	AlbumCount    int64 `json:"album_count"`
	ArtistCount   int64 `json:"artist_count"`
	HoursListened int64 `json:"hours_listened"`
}

func StatsHandler(store db.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := logger.FromContext(r.Context())
		var period db.Period
		switch strings.ToLower(r.URL.Query().Get("period")) {
		case "day":
			period = db.PeriodDay
		case "week":
			period = db.PeriodWeek
		case "month":
			period = db.PeriodMonth
		case "year":
			period = db.PeriodYear
		case "all_time":
			period = db.PeriodAllTime
		default:
			l.Debug().Msgf("Using default value '%s' for period", db.PeriodDay)
			period = db.PeriodDay
		}
		listens, err := store.CountListens(r.Context(), period)
		if err != nil {
			l.Err(err).Send()
			utils.WriteError(w, "failed to get listens: "+err.Error(), http.StatusInternalServerError)
			return
		}
		tracks, err := store.CountTracks(r.Context(), period)
		if err != nil {
			l.Err(err).Send()
			utils.WriteError(w, "failed to get listens: "+err.Error(), http.StatusInternalServerError)
			return
		}
		albums, err := store.CountAlbums(r.Context(), period)
		if err != nil {
			l.Err(err).Send()
			utils.WriteError(w, "failed to get listens: "+err.Error(), http.StatusInternalServerError)
			return
		}
		artists, err := store.CountArtists(r.Context(), period)
		if err != nil {
			l.Err(err).Send()
			utils.WriteError(w, "failed to get listens: "+err.Error(), http.StatusInternalServerError)
			return
		}
		timeListenedS, err := store.CountTimeListened(r.Context(), period)
		if err != nil {
			l.Err(err).Send()
			utils.WriteError(w, "failed to get listens: "+err.Error(), http.StatusInternalServerError)
			return
		}
		utils.WriteJSON(w, http.StatusOK, StatsResponse{
			ListenCount:   listens,
			TrackCount:    tracks,
			AlbumCount:    albums,
			ArtistCount:   artists,
			HoursListened: timeListenedS / 60 / 60,
		})
	}
}
