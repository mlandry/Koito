package engine_test

import (
	"context"
	"os"
	"path"
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

	src := path.Join("..", "test_assets", "maloja_import_test.json")
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

	truncateTestData(t)
}

func TestImportSpotify(t *testing.T) {

	src := path.Join("..", "test_assets", "Streaming_History_Audio_spotify_import_test.json")
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

	truncateTestData(t)
}

func TestImportLastFM(t *testing.T) {

	src := path.Join("..", "test_assets", "recenttracks-shoko2-1749776100.json")
	destDir := filepath.Join(cfg.ConfigDir(), "import")
	dest := filepath.Join(destDir, "recenttracks-shoko2-1749776100.json")

	// not going to make the dest dir because engine should make it already

	input, err := os.ReadFile(src)
	require.NoError(t, err)

	require.NoError(t, os.WriteFile(dest, input, os.ModePerm))

	mbzcMock := &mbz.MbzMockCaller{
		Artists: map[uuid.UUID]*mbz.MusicBrainzArtist{
			uuid.MustParse("4b00640f-3be6-43f8-9b34-ff81bd89320a"): &mbz.MusicBrainzArtist{
				Name: "OurR",
				Aliases: []mbz.MusicBrainzArtistAlias{
					{
						Name:    "OurR",
						Primary: true,
					},
				},
			},
		},
	}

	engine.RunImporter(logger.Get(), store, mbzcMock)

	album, err := store.GetAlbum(context.Background(), db.GetAlbumOpts{MusicBrainzID: uuid.MustParse("e9e78802-0bf8-4ca3-9655-1d943d2d2fa0")})
	require.NoError(t, err)
	assert.Equal(t, "ZOO!!", album.Title)
	artist, err := store.GetArtist(context.Background(), db.GetArtistOpts{MusicBrainzID: uuid.MustParse("4b00640f-3be6-43f8-9b34-ff81bd89320a")})
	require.NoError(t, err)
	assert.Equal(t, "OurR", artist.Name)
	artist, err = store.GetArtist(context.Background(), db.GetArtistOpts{Name: "CHUU"})
	require.NoError(t, err)
	track, err := store.GetTrack(context.Background(), db.GetTrackOpts{Title: "because I'm stupid?", ArtistIDs: []int32{artist.ID}})
	require.NoError(t, err)
	t.Log(track)
	listens, err := store.GetListensPaginated(context.Background(), db.GetItemsOpts{TrackID: int(track.ID), Period: db.PeriodAllTime})
	require.NoError(t, err)
	require.Len(t, listens.Items, 1)
	assert.WithinDuration(t, time.Unix(1749776100, 0), listens.Items[0].Time, 1*time.Second)

	truncateTestData(t)
}

func TestImportListenBrainz(t *testing.T) {

	src := path.Join("..", "test_assets", "listenbrainz_shoko1_1749780844.zip")
	destDir := filepath.Join(cfg.ConfigDir(), "import")
	dest := filepath.Join(destDir, "listenbrainz_shoko1_1749780844.zip")

	// not going to make the dest dir because engine should make it already

	input, err := os.ReadFile(src)
	require.NoError(t, err)

	require.NoError(t, os.WriteFile(dest, input, os.ModePerm))

	mbzcMock := &mbz.MbzMockCaller{
		Artists: map[uuid.UUID]*mbz.MusicBrainzArtist{
			uuid.MustParse("4b00640f-3be6-43f8-9b34-ff81bd89320a"): {
				Name: "OurR",
				Aliases: []mbz.MusicBrainzArtistAlias{
					{
						Name:    "OurR",
						Primary: true,
					},
				},
			},
			uuid.MustParse("09887aa7-226e-4ecc-9a0c-02d2ae5777e1"): {
				Name: "Carly Rae Jepsen",
				Aliases: []mbz.MusicBrainzArtistAlias{
					{
						Name:    "Carly Rae Jepsen",
						Primary: true,
					},
				},
			},
			uuid.MustParse("78e46ae5-9bfd-433b-be3f-19e993d67ecc"): &mbz.MusicBrainzArtist{
				Name: "Rufus Wainwright",
				Aliases: []mbz.MusicBrainzArtistAlias{
					{
						Name:    "OurR",
						Primary: true,
					},
				},
			},
		},
	}

	engine.RunImporter(logger.Get(), store, mbzcMock)

	album, err := store.GetAlbum(context.Background(), db.GetAlbumOpts{MusicBrainzID: uuid.MustParse("ce330d67-9c46-4a3b-9d62-08406370f234")})
	require.NoError(t, err)
	assert.Equal(t, "酸欠少女", album.Title)
	artist, err := store.GetArtist(context.Background(), db.GetArtistOpts{MusicBrainzID: uuid.MustParse("4b00640f-3be6-43f8-9b34-ff81bd89320a")})
	require.NoError(t, err)
	assert.Equal(t, "OurR", artist.Name)
	artist, err = store.GetArtist(context.Background(), db.GetArtistOpts{MusicBrainzID: uuid.MustParse("09887aa7-226e-4ecc-9a0c-02d2ae5777e1")})
	require.NoError(t, err)
	assert.Equal(t, "Carly Rae Jepsen", artist.Name)
	artist, err = store.GetArtist(context.Background(), db.GetArtistOpts{MusicBrainzID: uuid.MustParse("78e46ae5-9bfd-433b-be3f-19e993d67ecc")})
	require.NoError(t, err)
	assert.Equal(t, "Rufus Wainwright", artist.Name)
	track, err := store.GetTrack(context.Background(), db.GetTrackOpts{MusicBrainzID: uuid.MustParse("08e8f55b-f1a4-46b8-b2d1-fab4c592165c")})
	require.NoError(t, err)
	assert.Equal(t, "Desert", track.Title)
	listens, err := store.GetListensPaginated(context.Background(), db.GetItemsOpts{TrackID: int(track.ID), Period: db.PeriodAllTime})
	require.NoError(t, err)
	assert.Len(t, listens.Items, 1)
	assert.WithinDuration(t, time.Unix(1749780612, 0), listens.Items[0].Time, 1*time.Second)

	truncateTestData(t)
}
