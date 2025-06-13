// Package catalog manages the internal metadata of the catalog of music the user has submitted listens for.
// This includes artists, releases (album, single, ep, etc), and tracks, as well as ingesting
// listens submitted both via the API(s) and other methods.
package catalog

import (
	"context"
	"errors"
	"regexp"
	"strings"
	"time"

	"github.com/gabehf/koito/internal/db"
	"github.com/gabehf/koito/internal/logger"
	"github.com/gabehf/koito/internal/mbz"
	"github.com/gabehf/koito/internal/models"
	"github.com/google/uuid"
)

type GetListensOpts struct {
	ArtistID       int32
	ReleaseGroupID int32
	TrackID        int32
	Limit          int
}

type SaveListenOpts struct {
	TrackID int32
	Time    time.Time
}

type SubmitListenOpts struct {
	// When true, skips registering the listen and only associates or creates the
	// artist, release, release group, and track in DB
	SkipSaveListen bool

	MbzCaller         mbz.MusicBrainzCaller
	ArtistNames       []string
	Artist            string
	ArtistMbzIDs      []uuid.UUID
	TrackTitle        string
	RecordingMbzID    uuid.UUID
	Duration          int32 // in seconds
	ReleaseTitle      string
	ReleaseMbzID      uuid.UUID
	ReleaseGroupMbzID uuid.UUID
	Time              time.Time
	UserID            int32
	Client            string
}

const (
	ImageSourceUserUpload = "User Upload"
)

func SubmitListen(ctx context.Context, store db.DB, opts SubmitListenOpts) error {
	l := logger.FromContext(ctx)

	if opts.Artist == "" || opts.TrackTitle == "" {
		return errors.New("track name and artist are required")
	}

	artists, err := AssociateArtists(
		ctx,
		store,
		AssociateArtistsOpts{
			ArtistMbzIDs: opts.ArtistMbzIDs,
			ArtistNames:  opts.ArtistNames,
			ArtistName:   opts.Artist,
			Mbzc:         opts.MbzCaller,
			TrackTitle:   opts.TrackTitle,
		})
	if err != nil {
		l.Error().Err(err).Msg("Failed to associate artists to listen")
		return err
	} else if len(artists) < 1 {
		l.Debug().Msg("Failed to associate any artists to release")
	}

	artistIDs := make([]int32, len(artists))

	for i, artist := range artists {
		artistIDs[i] = artist.ID
		l.Debug().Any("artist", artist).Msg("Matched listen to artist")
	}
	rg, err := AssociateAlbum(ctx, store, AssociateAlbumOpts{
		ReleaseMbzID:      opts.ReleaseMbzID,
		ReleaseGroupMbzID: opts.ReleaseGroupMbzID,
		ReleaseName:       opts.ReleaseTitle,
		TrackName:         opts.TrackTitle,
		Mbzc:              opts.MbzCaller,
		Artists:           artists,
	})
	if err != nil {
		l.Error().Err(err).Msg("Failed to associate release group to listen")
		return err
	}
	l.Debug().Any("album", rg).Msg("Matched listen to release")

	// ensure artists are associated with release group
	store.AddArtistsToAlbum(ctx, db.AddArtistsToAlbumOpts{
		ArtistIDs: artistIDs,
		AlbumID:   rg.ID,
	})

	track, err := AssociateTrack(ctx, store, AssociateTrackOpts{
		ArtistIDs:  artistIDs,
		AlbumID:    rg.ID,
		TrackMbzID: opts.RecordingMbzID,
		TrackName:  opts.TrackTitle,
		Duration:   opts.Duration,
		Mbzc:       opts.MbzCaller,
	})
	if err != nil {
		l.Error().Err(err).Msg("Failed to associate track to listen")
		return err
	}
	l.Debug().Any("track", track).Msg("Matched listen to track")

	if track.Duration == 0 && opts.Duration != 0 {
		err := store.UpdateTrack(ctx, db.UpdateTrackOpts{
			ID:       track.ID,
			Duration: opts.Duration,
		})
		if err != nil {
			l.Err(err).Msgf("Failed to update duration for track %s", track.Title)
		}
		l.Info().Msgf("Duration updated to %d for track '%s'", opts.Duration, track.Title)
	}

	if opts.SkipSaveListen {
		return nil
	}

	l.Info().Msgf("Received listen: '%s' by %s, from release '%s'", track.Title, buildArtistStr(artists), rg.Title)

	return store.SaveListen(ctx, db.SaveListenOpts{
		TrackID: track.ID,
		Time:    opts.Time,
		UserID:  opts.UserID,
		Client:  opts.Client,
	})
}

func buildArtistStr(artists []*models.Artist) string {
	artistNames := make([]string, len(artists))
	for i, artist := range artists {
		artistNames[i] = artist.Name
	}
	return strings.Join(artistNames, " & ")
}

var (
	// Bracketed feat patterns
	bracketFeatPatterns = []*regexp.Regexp{
		regexp.MustCompile(`(?i)\(feat\. ([^)]*)\)`),
		regexp.MustCompile(`(?i)\[feat\. ([^\]]*)\]`),
	}
	// Inline feat (not in brackets)
	inlineFeatPattern = regexp.MustCompile(`(?i)feat\. ([^()\[\]]+)$`)

	// Delimiters only used inside feat. sections
	featSplitDelimiters = regexp.MustCompile(`(?i)\s*(?:,|&|and|·)\s*`)

	// Delimiter for separating artists in main string (rare but real usage)
	mainArtistDotSplitter = regexp.MustCompile(`\s+·\s+`)
)

// ParseArtists extracts all contributing artist names from the artist and title strings
func ParseArtists(artist string, title string) []string {
	seen := make(map[string]struct{})
	var out []string

	add := func(name string) {
		name = strings.TrimSpace(name)
		if name == "" {
			return
		}
		if _, exists := seen[name]; !exists {
			seen[name] = struct{}{}
			out = append(out, name)
		}
	}

	foundFeat := false

	// Extract bracketed features from artist
	for _, re := range bracketFeatPatterns {
		if matches := re.FindStringSubmatch(artist); matches != nil {
			foundFeat = true
			artist = strings.Replace(artist, matches[0], "", 1)
			for _, name := range featSplitDelimiters.Split(matches[1], -1) {
				add(name)
			}
		}
	}
	// Extract inline feat. from artist
	if matches := inlineFeatPattern.FindStringSubmatch(artist); matches != nil {
		foundFeat = true
		artist = strings.Replace(artist, matches[0], "", 1)
		for _, name := range featSplitDelimiters.Split(matches[1], -1) {
			add(name)
		}
	}

	// Add base artist(s)
	if foundFeat {
		add(strings.TrimSpace(artist))
	} else {
		// Only split on " · " in base artist string
		for _, name := range mainArtistDotSplitter.Split(artist, -1) {
			add(name)
		}
	}

	// Extract features from title
	for _, re := range bracketFeatPatterns {
		if matches := re.FindStringSubmatch(title); matches != nil {
			for _, name := range featSplitDelimiters.Split(matches[1], -1) {
				add(name)
			}
		}
	}
	if matches := inlineFeatPattern.FindStringSubmatch(title); matches != nil {
		for _, name := range featSplitDelimiters.Split(matches[1], -1) {
			add(name)
		}
	}

	return out
}
