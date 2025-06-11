package psql_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestDataForMerge(t *testing.T) {
	truncateTestData(t)
	// Insert artists
	err := store.Exec(context.Background(),
		`INSERT INTO artists (musicbrainz_id) 
			VALUES ('00000000-0000-0000-0000-000000000001'),
				   ('00000000-0000-0000-0000-000000000002')`)
	require.NoError(t, err)

	err = store.Exec(context.Background(),
		`INSERT INTO artist_aliases (artist_id, alias, source, is_primary) 
			VALUES (1, 'Artist One', 'Testing', true),
				   (2, 'Artist Two', 'Testing', true)`)
	require.NoError(t, err)

	// Insert albums
	err = store.Exec(context.Background(),
		`INSERT INTO releases (musicbrainz_id) 
			VALUES ('11111111-1111-1111-1111-111111111111'),
				   ('22222222-2222-2222-2222-222222222222')`)
	require.NoError(t, err)

	err = store.Exec(context.Background(),
		`INSERT INTO release_aliases (release_id, alias, source, is_primary) 
			VALUES (1, 'Album One', 'Testing', true),
				   (2, 'Album Two', 'Testing', true)`)
	require.NoError(t, err)

	// Insert tracks
	err = store.Exec(context.Background(),
		`INSERT INTO tracks (musicbrainz_id, release_id) 
			VALUES ('33333333-3333-3333-3333-333333333333', 1),
				   ('44444444-4444-4444-4444-444444444444', 2)`)
	require.NoError(t, err)

	err = store.Exec(context.Background(),
		`INSERT INTO track_aliases (track_id, alias, source, is_primary) 
			VALUES (1, 'Track One', 'Testing', true),
				   (2, 'Track Two', 'Testing', true)`)
	require.NoError(t, err)

	// Associate artists with albums and tracks
	err = store.Exec(context.Background(),
		`INSERT INTO artist_releases (artist_id, release_id) 
			VALUES (1, 1), (2, 2)`)
	require.NoError(t, err)

	err = store.Exec(context.Background(),
		`INSERT INTO artist_tracks (artist_id, track_id) 
			VALUES (1, 1), (2, 2)`)
	require.NoError(t, err)

	// Insert listens
	err = store.Exec(context.Background(),
		`INSERT INTO listens (user_id, track_id, listened_at) 
			VALUES (1, 1, NOW() - INTERVAL '1 day'),
				   (1, 2, NOW() - INTERVAL '2 days')`)
	require.NoError(t, err)
}

func TestMergeTracks(t *testing.T) {
	ctx := context.Background()
	setupTestDataForMerge(t)

	// Merge Track 1 into Track 2
	err := store.MergeTracks(ctx, 1, 2)
	require.NoError(t, err)

	// Verify listens are updated
	var count int
	count, err = store.Count(ctx, `SELECT COUNT(*) FROM listens WHERE track_id = 2`)
	require.NoError(t, err)
	assert.Equal(t, 2, count, "expected all listens to be merged into Track 2")

	truncateTestData(t)
}

func TestMergeAlbums(t *testing.T) {
	ctx := context.Background()
	setupTestDataForMerge(t)

	// Merge Album 1 into Album 2
	err := store.MergeAlbums(ctx, 1, 2)
	require.NoError(t, err)

	// Verify tracks are updated
	var count int
	count, err = store.Count(ctx, `SELECT COUNT(*) FROM tracks WHERE release_id = 2`)
	require.NoError(t, err)
	assert.Equal(t, 2, count, "expected all tracks to be merged into Album 2")

	truncateTestData(t)
}

func TestMergeArtists(t *testing.T) {
	ctx := context.Background()
	setupTestDataForMerge(t)

	// Merge Artist 1 into Artist 2
	err := store.MergeArtists(ctx, 1, 2)
	require.NoError(t, err)

	// Verify artist associations are updated
	var count int
	count, err = store.Count(ctx, `SELECT COUNT(*) FROM artist_tracks WHERE artist_id = 2`)
	require.NoError(t, err)
	assert.Equal(t, 2, count, "expected all tracks to be associated with Artist 2")

	count, err = store.Count(ctx, `SELECT COUNT(*) FROM artist_releases WHERE artist_id = 2`)
	require.NoError(t, err)
	assert.Equal(t, 2, count, "expected all releases to be associated with Artist 2")

	truncateTestData(t)
}
