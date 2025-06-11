package importer

import (
	"context"
	"encoding/json"
	"os"
	"path"
	"time"

	"github.com/gabehf/koito/internal/catalog"
	"github.com/gabehf/koito/internal/cfg"
	"github.com/gabehf/koito/internal/db"
	"github.com/gabehf/koito/internal/logger"
	"github.com/gabehf/koito/internal/mbz"
)

type SpotifyExportItem struct {
	Timestamp  time.Time `json:"ts"`
	TrackName  string    `json:"master_metadata_track_name"`
	ArtistName string    `json:"master_metadata_album_artist_name"`
	AlbumName  string    `json:"master_metadata_album_album_name"`
	ReasonEnd  string    `json:"reason_end"`
	MsPlayed   int32     `json:"ms_played"`
}

func ImportSpotifyFile(ctx context.Context, store db.DB, filename string) error {
	l := logger.FromContext(ctx)
	l.Info().Msgf("Beginning spotify import on file: %s", filename)
	file, err := os.Open(path.Join(cfg.ConfigDir(), "import", filename))
	if err != nil {
		l.Err(err).Msgf("Failed to read import file: %s", filename)
		return err
	}
	export := make([]SpotifyExportItem, 0)
	err = json.NewDecoder(file).Decode(&export)
	if err != nil {
		return err
	}
	for _, item := range export {
		if item.ReasonEnd != "trackdone" {
			continue
		}
		dur := item.MsPlayed
		if item.TrackName == "" || item.ArtistName == "" {
			l.Debug().Msg("Skipping non-track item")
			continue
		}
		opts := catalog.SubmitListenOpts{
			MbzCaller:    &mbz.MusicBrainzClient{},
			Artist:       item.ArtistName,
			TrackTitle:   item.TrackName,
			ReleaseTitle: item.AlbumName,
			Duration:     dur / 1000,
			Time:         item.Timestamp,
			UserID:       1,
		}
		err = catalog.SubmitListen(ctx, store, opts)
		if err != nil {
			l.Err(err).Msg("Failed to import spotify playback item")
			return err
		}
	}
	_, err = os.Stat(path.Join(cfg.ConfigDir(), "import_complete"))
	if err != nil {
		err = os.Mkdir(path.Join(cfg.ConfigDir(), "import_complete"), 0744)
		if err != nil {
			l.Err(err).Msg("Failed to create import_complete dir! Import files must be removed from the import directory manually, or else the importer will run on every app start")
		}
	}
	err = os.Rename(path.Join(cfg.ConfigDir(), "import", filename), path.Join(cfg.ConfigDir(), "import_complete", filename))
	if err != nil {
		l.Err(err).Msg("Failed to move file to import_complete dir! Import files must be removed from the import directory manually, or else the importer will run on every app start")
	}
	l.Info().Msgf("Finished importing %s; imported %d items", filename, len(export))
	return nil
}
