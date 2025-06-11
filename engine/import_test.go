package engine_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/gabehf/koito/engine"
	"github.com/gabehf/koito/internal/cfg"
	"github.com/gabehf/koito/internal/db"
	"github.com/gabehf/koito/internal/logger"
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

	engine.RunImporter(logger.Get(), store)

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

	engine.RunImporter(logger.Get(), store)

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
