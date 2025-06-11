package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gabehf/koito/internal/db"
	"github.com/gabehf/koito/internal/logger"
	"github.com/gabehf/koito/internal/utils"
)

func GetListenActivityHandler(store db.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		l := logger.FromContext(r.Context())

		rangeStr := r.URL.Query().Get("range")
		_range, _ := strconv.Atoi(rangeStr)

		monthStr := r.URL.Query().Get("month")
		month, _ := strconv.Atoi(monthStr)
		yearStr := r.URL.Query().Get("year")
		year, _ := strconv.Atoi(yearStr)

		artistIdStr := r.URL.Query().Get("artist_id")
		artistId, _ := strconv.Atoi(artistIdStr)
		albumIdStr := r.URL.Query().Get("album_id")
		albumId, _ := strconv.Atoi(albumIdStr)
		trackIdStr := r.URL.Query().Get("track_id")
		trackId, _ := strconv.Atoi(trackIdStr)

		var step db.StepInterval
		switch strings.ToLower(r.URL.Query().Get("step")) {
		case "day":
			step = db.StepDay
		case "week":
			step = db.StepWeek
		case "month":
			step = db.StepMonth
		case "year":
			step = db.StepYear
		default:
			l.Debug().Msgf("Using default value '%s' for step", db.StepDefault)
			step = db.StepDay
		}

		opts := db.ListenActivityOpts{
			Step:     step,
			Range:    _range,
			Month:    month,
			Year:     year,
			AlbumID:  int32(albumId),
			ArtistID: int32(artistId),
			TrackID:  int32(trackId),
		}

		activity, err := store.GetListenActivity(r.Context(), opts)
		if err != nil {
			l.Err(err).Send()
			utils.WriteError(w, err.Error(), 500)
			return
		}
		utils.WriteJSON(w, http.StatusOK, activity)
	}
}
