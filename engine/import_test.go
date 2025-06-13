package engine_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/gabehf/koito/engine"
	"github.com/gabehf/koito/internal/cfg"
	"github.com/gabehf/koito/internal/db"
	"github.com/gabehf/koito/internal/logger"
	"github.com/gabehf/koito/internal/mbz"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestImportMaloja(t *testing.T) {

	src := "../static/maloja_import_test.json"
	destDir := filepath.Join(cfg.ConfigDir(), "import")
	dest := filepath.Join(destDir, "maloja_import_test.json")

	// not going to make the dest dir because engine should make it already

	input, err := os.ReadFile(src)
	require.NoError(t, err)

	require.NoError(t, os.WriteFile(dest, input, os.ModePerm))

	engine.RunImporter(logger.Get(), store, &mbz.MbzErrorCaller{})

	// maloja test import is 38 Magnify Tokyo streams
	a, err := store.GetArtist(context.Background(), db.GetArtistOpts{Name: "Magnify Tokyo"})
	require.NoError(t, err)
	t.Log(a)
	assert.Equal(t, "Magnify Tokyo", a.Name)
	assert.EqualValues(t, 38, a.ListenCount)
}

func TestImportSpotify(t *testing.T) {

	src := "../static/Streaming_History_Audio_spotify_import_test.json"
	destDir := filepath.Join(cfg.ConfigDir(), "import")
	dest := filepath.Join(destDir, "Streaming_History_Audio_spotify_import_test.json")

	// not going to make the dest dir because engine should make it already

	input, err := os.ReadFile(src)
	require.NoError(t, err)

	require.NoError(t, os.WriteFile(dest, input, os.ModePerm))

	engine.RunImporter(logger.Get(), store, &mbz.MbzErrorCaller{})

	a, err := store.GetArtist(context.Background(), db.GetArtistOpts{Name: "The Story So Far"})
	require.NoError(t, err)
	track, err := store.GetTrack(context.Background(), db.GetTrackOpts{Title: "Clairvoyant", ArtistIDs: []int32{a.ID}})
	require.NoError(t, err)
	t.Log(track)
	assert.Equal(t, "Clairvoyant", track.Title)
	// spotify includes duration data, but we only import when reason_end = trackdone
	// this is the only track with valid duration data
	assert.EqualValues(t, 181, track.Duration)
}

func TestImportLastFM(t *testing.T) {

	src := "../static/recenttracks-shoko2-1749776100.json"
	destDir := filepath.Join(cfg.ConfigDir(), "import")
	dest := filepath.Join(destDir, "recenttracks-shoko2-1749776100.json")

	// not going to make the dest dir because engine should make it already

	input, err := os.ReadFile(src)
	require.NoError(t, err)

	require.NoError(t, os.WriteFile(dest, input, os.ModePerm))

	engine.RunImporter(logger.Get(), store, &mbz.MbzErrorCaller{})

	album, err := store.GetAlbum(context.Background(), db.GetAlbumOpts{MusicBrainzID: uuid.MustParse("e9e78802-0bf8-4ca3-9655-1d943d2d2fa0")})
	require.NoError(t, err)
	assert.Equal(t, "ZOO!!", album.Title)
	artist, err := store.GetArtist(context.Background(), db.GetArtistOpts{Name: "CHUU"})
	require.NoError(t, err)
	track, err := store.GetTrack(context.Background(), db.GetTrackOpts{Title: "because I'm stupid?", ArtistIDs: []int32{artist.ID}})
	require.NoError(t, err)
	listens, err := store.GetListensPaginated(context.Background(), db.GetItemsOpts{TrackID: int(track.ID)})
	require.NoError(t, err)
	assert.Len(t, listens.Items, 1)
	assert.WithinDuration(t, time.Unix(1749776100, 0), listens.Items[0].Time, 1*time.Second)
}
