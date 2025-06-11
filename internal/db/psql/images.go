package psql

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/gabehf/koito/internal/logger"
	"github.com/gabehf/koito/internal/models"
	"github.com/gabehf/koito/internal/repository"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func (d *Psql) ImageHasAssociation(ctx context.Context, image uuid.UUID) (bool, error) {
	_, err := d.q.GetReleaseByImageID(ctx, &image)
	if err == nil {
		return true, err
	} else if !errors.Is(err, pgx.ErrNoRows) {
		return false, err
	}
	_, err = d.q.GetArtistByImage(ctx, &image)
	if err == nil {
		return true, err
	} else if !errors.Is(err, pgx.ErrNoRows) {
		return false, err
	}
	return false, nil
}

func (d *Psql) GetImageSource(ctx context.Context, image uuid.UUID) (string, error) {
	r, err := d.q.GetReleaseByImageID(ctx, &image)
	if err == nil {
		return r.ImageSource.String, err
	} else if !errors.Is(err, pgx.ErrNoRows) {
		return "", err
	}
	rr, err := d.q.GetArtistByImage(ctx, &image)
	if err == nil {
		return rr.ImageSource.String, err
	} else if !errors.Is(err, pgx.ErrNoRows) {
		return "", err
	}
	return "", nil
}

func (d *Psql) AlbumsWithoutImages(ctx context.Context, from int32) ([]*models.Album, error) {
	l := logger.FromContext(ctx)
	rows, err := d.q.GetReleasesWithoutImages(ctx, repository.GetReleasesWithoutImagesParams{
		Limit: 20,
		ID:    from,
	})
	if err != nil {
		return nil, err
	}
	albums := make([]*models.Album, len(rows))
	for i, row := range rows {
		artists := make([]models.SimpleArtist, 0)
		err = json.Unmarshal(row.Artists, &artists)
		if err != nil {
			l.Err(err).Msgf("Error unmarshalling artists for release group with id %d", row.ID)
			artists = nil
		}
		albums[i] = &models.Album{
			ID:             row.ID,
			Image:          row.Image,
			Title:          row.Title,
			MbzID:          row.MusicBrainzID,
			VariousArtists: row.VariousArtists,
			Artists:        artists,
		}
	}
	return albums, nil
}
