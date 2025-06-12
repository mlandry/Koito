package catalog

import (
	"context"
	"errors"
	"slices"
	"strings"

	"github.com/gabehf/koito/internal/cfg"
	"github.com/gabehf/koito/internal/db"
	"github.com/gabehf/koito/internal/images"
	"github.com/gabehf/koito/internal/logger"
	"github.com/gabehf/koito/internal/mbz"
	"github.com/gabehf/koito/internal/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type AssociateArtistsOpts struct {
	ArtistMbzIDs []uuid.UUID
	ArtistNames  []string
	ArtistName   string
	TrackTitle   string
	Mbzc         mbz.MusicBrainzCaller
	IP           *ImageProcessor
}

func AssociateArtists(ctx context.Context, d db.DB, opts AssociateArtistsOpts) ([]*models.Artist, error) {
	l := logger.FromContext(ctx)

	var result []*models.Artist

	if len(opts.ArtistMbzIDs) > 0 {
		l.Debug().Msg("Associating artists by MusicBrainz ID(s)")
		mbzMatches, err := matchArtistsByMBID(ctx, d, opts)
		if err != nil {
			return nil, err
		}
		result = append(result, mbzMatches...)
	}

	if len(opts.ArtistNames) > len(result) {
		l.Debug().Msg("Associating artists by list of artist names")
		nameMatches, err := matchArtistsByNames(ctx, opts.ArtistNames, result, d, opts)
		if err != nil {
			return nil, err
		}
		result = append(result, nameMatches...)
	}

	if len(result) < 1 {
		allArtists := slices.Concat(opts.ArtistNames, ParseArtists(opts.ArtistName, opts.TrackTitle))
		l.Debug().Msgf("Associating artists by artist name(s) %v and track title '%s'", allArtists, opts.TrackTitle)
		fallbackMatches, err := matchArtistsByNames(ctx, allArtists, nil, d, opts)
		if err != nil {
			return nil, err
		}
		result = append(result, fallbackMatches...)
	}

	return result, nil
}

func matchArtistsByMBID(ctx context.Context, d db.DB, opts AssociateArtistsOpts) ([]*models.Artist, error) {
	l := logger.FromContext(ctx)
	var result []*models.Artist

	for _, id := range opts.ArtistMbzIDs {
		if id == uuid.Nil {
			l.Warn().Msg("Provided artist has uuid.Nil MusicBrainzID")
			return matchArtistsByNames(ctx, opts.ArtistNames, result, d, opts)
		}
		a, err := d.GetArtist(ctx, db.GetArtistOpts{
			MusicBrainzID: id,
		})
		if err == nil {
			l.Debug().Msgf("Artist '%s' found by MusicBrainz ID", a.Name)
			result = append(result, a)
			continue
		}

		if !errors.Is(err, pgx.ErrNoRows) {
			return nil, err
		}

		if len(opts.ArtistNames) < 1 {
			opts.ArtistNames = slices.Concat(opts.ArtistNames, ParseArtists(opts.ArtistName, opts.TrackTitle))
		}
		a, err = resolveAliasOrCreateArtist(ctx, id, opts.ArtistNames, d, opts)
		if err != nil {
			l.Warn().Msg("MusicBrainz unreachable, falling back to artist name matching")
			return matchArtistsByNames(ctx, opts.ArtistNames, result, d, opts)
			// return nil, err
		}
		result = append(result, a)
	}
	return result, nil
}
func resolveAliasOrCreateArtist(ctx context.Context, mbzID uuid.UUID, names []string, d db.DB, opts AssociateArtistsOpts) (*models.Artist, error) {
	l := logger.FromContext(ctx)

	aliases, err := opts.Mbzc.GetArtistPrimaryAliases(ctx, mbzID)
	if err != nil {
		return nil, err
	}
	l.Debug().Msgf("Got aliases %v from MusicBrainz", aliases)

	for _, alias := range aliases {
		a, err := d.GetArtist(ctx, db.GetArtistOpts{
			Name: alias,
		})
		if err == nil && (a.MbzID == nil || *a.MbzID == uuid.Nil) {
			a.MbzID = &mbzID
			l.Debug().Msgf("Alias '%s' found in DB. Associating with MusicBrainz ID...", alias)
			if updateErr := d.UpdateArtist(ctx, db.UpdateArtistOpts{ID: a.ID, MusicBrainzID: mbzID}); updateErr != nil {
				return nil, updateErr
			}
			if saveAliasErr := d.SaveArtistAliases(ctx, a.ID, aliases, "MusicBrainz"); saveAliasErr != nil {
				return nil, saveAliasErr
			}
			return a, nil
		}
	}

	canonical := aliases[0]
	for _, alias := range aliases {
		for _, name := range names {
			if strings.EqualFold(alias, name) {
				l.Debug().Msgf("Canonical name for artist is '%s'", alias)
				canonical = alias
				break
			}
		}
	}

	var imgid uuid.UUID
	imgUrl, err := images.GetArtistImage(ctx, images.ArtistImageOpts{
		Aliases: aliases,
	})
	if err == nil && imgUrl != "" {
		var size ImageSize
		if cfg.FullImageCacheEnabled() {
			size = ImageSizeFull
		} else {
			size = ImageSizeLarge
		}
		imgid = uuid.New()
		l.Debug().Msg("Downloading artist image from source...")
		err = opts.IP.EnqueueDownloadAndCache(ctx, imgid, imgUrl, size)
		if err != nil {
			l.Err(err).Msg("Failed to cache image")
		}
	} else if err != nil {
		l.Warn().Msgf("Failed to get artist image from ImageSrc: %s", err.Error())
	}

	u, err := d.SaveArtist(ctx, db.SaveArtistOpts{
		MusicBrainzID: mbzID,
		Name:          canonical,
		Aliases:       aliases,
		Image:         imgid,
		ImageSrc:      imgUrl,
	})
	if err != nil {
		return nil, err
	}
	l.Info().Msgf("Created artist '%s' with MusicBrainz Artist ID", canonical)
	return u, nil
}

func matchArtistsByNames(ctx context.Context, names []string, existing []*models.Artist, d db.DB, opts AssociateArtistsOpts) ([]*models.Artist, error) {
	l := logger.FromContext(ctx)
	var result []*models.Artist

	for _, name := range names {
		if artistExists(name, existing) || artistExists(name, result) {
			l.Debug().Msgf("Artist '%s' already found, skipping...", name)
			continue
		}
		a, err := d.GetArtist(ctx, db.GetArtistOpts{
			Name: name,
		})
		if err == nil {
			l.Debug().Msgf("Artist '%s' found in DB", name)
			result = append(result, a)
			continue
		}
		if errors.Is(err, pgx.ErrNoRows) {
			var imgid uuid.UUID
			imgUrl, err := images.GetArtistImage(ctx, images.ArtistImageOpts{
				Aliases: []string{name},
			})
			if err == nil && imgUrl != "" {
				var size ImageSize
				if cfg.FullImageCacheEnabled() {
					size = ImageSizeFull
				} else {
					size = ImageSizeLarge
				}
				imgid = uuid.New()
				l.Debug().Msg("Downloading artist image from source...")
				err = opts.IP.EnqueueDownloadAndCache(ctx, imgid, imgUrl, size)
				if err != nil {
					l.Err(err).Msg("Failed to cache image")
				}
			} else if err != nil {
				l.Debug().Msgf("Failed to get artist images for %s: %s", name, err.Error())
			}
			a, err = d.SaveArtist(ctx, db.SaveArtistOpts{Name: name, Image: imgid, ImageSrc: imgUrl})
			if err != nil {
				return nil, err
			}
			l.Info().Msgf("Created artist '%s' with artist name", name)
			result = append(result, a)
		} else {
			return nil, err
		}
	}
	return result, nil
}

func artistExists(name string, artists []*models.Artist) bool {
	for _, a := range artists {
		allAliases := append(a.Aliases, a.Name)
		for _, alias := range allAliases {
			if strings.EqualFold(name, alias) {
				return true
			}
		}
	}
	return false
}
