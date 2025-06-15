package psql

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/gabehf/koito/internal/db"
	"github.com/gabehf/koito/internal/logger"
	"github.com/gabehf/koito/internal/models"
	"github.com/gabehf/koito/internal/repository"
	"github.com/gabehf/koito/internal/utils"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

// this function sucks because sqlc keeps making new types for rows that are the same
func (d *Psql) GetArtist(ctx context.Context, opts db.GetArtistOpts) (*models.Artist, error) {
	l := logger.FromContext(ctx)
	if opts.ID != 0 {
		l.Debug().Msgf("Fetching artist from DB with id %d", opts.ID)
		row, err := d.q.GetArtist(ctx, opts.ID)
		if err != nil {
			return nil, err
		}
		count, err := d.q.CountListensFromArtist(ctx, repository.CountListensFromArtistParams{
			ListenedAt:   time.Unix(0, 0),
			ListenedAt_2: time.Now(),
			ArtistID:     row.ID,
		})
		if err != nil {
			return nil, err
		}
		seconds, err := d.CountTimeListenedToItem(ctx, db.TimeListenedOpts{
			Period:   db.PeriodAllTime,
			ArtistID: row.ID,
		})
		if err != nil {
			return nil, err
		}
		return &models.Artist{
			ID:           row.ID,
			MbzID:        row.MusicBrainzID,
			Name:         row.Name,
			Aliases:      row.Aliases,
			Image:        row.Image,
			ListenCount:  count,
			TimeListened: seconds,
		}, nil
	} else if opts.MusicBrainzID != uuid.Nil {
		l.Debug().Msgf("Fetching artist from DB with MusicBrainz ID %s", opts.MusicBrainzID)
		row, err := d.q.GetArtistByMbzID(ctx, &opts.MusicBrainzID)
		if err != nil {
			return nil, err
		}
		count, err := d.q.CountListensFromArtist(ctx, repository.CountListensFromArtistParams{
			ListenedAt:   time.Unix(0, 0),
			ListenedAt_2: time.Now(),
			ArtistID:     row.ID,
		})
		if err != nil {
			return nil, err
		}
		seconds, err := d.CountTimeListenedToItem(ctx, db.TimeListenedOpts{
			Period:   db.PeriodAllTime,
			ArtistID: row.ID,
		})
		if err != nil {
			return nil, err
		}
		return &models.Artist{
			ID:           row.ID,
			MbzID:        row.MusicBrainzID,
			Name:         row.Name,
			Aliases:      row.Aliases,
			Image:        row.Image,
			TimeListened: seconds,
			ListenCount:  count,
		}, nil
	} else if opts.Name != "" {
		l.Debug().Msgf("Fetching artist from DB with name '%s'", opts.Name)
		row, err := d.q.GetArtistByName(ctx, opts.Name)
		if err != nil {
			return nil, err
		}
		count, err := d.q.CountListensFromArtist(ctx, repository.CountListensFromArtistParams{
			ListenedAt:   time.Unix(0, 0),
			ListenedAt_2: time.Now(),
			ArtistID:     row.ID,
		})
		if err != nil {
			return nil, err
		}
		seconds, err := d.CountTimeListenedToItem(ctx, db.TimeListenedOpts{
			Period:   db.PeriodAllTime,
			ArtistID: row.ID,
		})
		if err != nil {
			return nil, err
		}
		return &models.Artist{
			ID:           row.ID,
			MbzID:        row.MusicBrainzID,
			Name:         row.Name,
			Aliases:      row.Aliases,
			Image:        row.Image,
			ListenCount:  count,
			TimeListened: seconds,
		}, nil
	} else {
		return nil, errors.New("insufficient information to get artist")
	}
}

// Inserts all unique aliases into the DB with specified source
func (d *Psql) SaveArtistAliases(ctx context.Context, id int32, aliases []string, source string) error {
	l := logger.FromContext(ctx)
	if id == 0 {
		return errors.New("artist id not specified")
	}
	tx, err := d.conn.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		l.Err(err).Msg("Failed to begin transaction")
		return err
	}
	defer tx.Rollback(ctx)
	qtx := d.q.WithTx(tx)
	existing, err := qtx.GetAllArtistAliases(ctx, id)
	if err != nil {
		return err
	}
	for _, v := range existing {
		aliases = append(aliases, v.Alias)
	}
	utils.Unique(&aliases)
	for _, alias := range aliases {
		if strings.TrimSpace(alias) == "" {
			return errors.New("aliases cannot be blank")
		}
		err = qtx.InsertArtistAlias(ctx, repository.InsertArtistAliasParams{
			Alias:     strings.TrimSpace(alias),
			ArtistID:  id,
			Source:    source,
			IsPrimary: false,
		})
		if err != nil {
			return err
		}
	}
	return tx.Commit(ctx)
}

func (d *Psql) DeleteArtist(ctx context.Context, id int32) error {
	return d.q.DeleteArtist(ctx, id)
}

// Equivalent to Psql.SaveArtist, then Psql.SaveMbzAliases
func (d *Psql) SaveArtist(ctx context.Context, opts db.SaveArtistOpts) (*models.Artist, error) {
	l := logger.FromContext(ctx)
	var insertMbzID *uuid.UUID
	var insertImage *uuid.UUID
	if opts.MusicBrainzID != uuid.Nil {
		insertMbzID = &opts.MusicBrainzID
	}
	if opts.Image != uuid.Nil {
		insertImage = &opts.Image
	}
	tx, err := d.conn.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		l.Err(err).Msg("Failed to begin transaction")
		return nil, err
	}
	defer tx.Rollback(ctx)
	qtx := d.q.WithTx(tx)
	opts.Name = strings.TrimSpace(opts.Name)
	if opts.Name == "" {
		return nil, errors.New("name must not be blank")
	}
	l.Debug().Msgf("Inserting artist '%s' into DB", opts.Name)
	a, err := qtx.InsertArtist(ctx, repository.InsertArtistParams{
		MusicBrainzID: insertMbzID,
		Image:         insertImage,
		ImageSource:   pgtype.Text{String: opts.ImageSrc, Valid: opts.ImageSrc != ""},
	})
	if err != nil {
		return nil, err
	}
	l.Debug().Msgf("Inserting canonical alias '%s' into DB for artist with id %d", opts.Name, a.ID)
	err = qtx.InsertArtistAlias(ctx, repository.InsertArtistAliasParams{
		ArtistID:  a.ID,
		Alias:     opts.Name,
		Source:    "Canonical",
		IsPrimary: true,
	})
	if err != nil {
		l.Error().Err(err).Msgf("Error inserting canonical alias for artist '%s'", opts.Name)
		return nil, err
	}
	err = tx.Commit(ctx)
	if err != nil {
		l.Err(err).Msg("Failed to commit insert artist transaction")
		return nil, err
	}
	artist := &models.Artist{
		ID:      a.ID,
		Name:    opts.Name,
		Image:   a.Image,
		MbzID:   a.MusicBrainzID,
		Aliases: []string{opts.Name},
	}
	if len(opts.Aliases) > 0 {
		l.Debug().Msgf("Inserting aliases '%v' into DB for artist '%s'", opts.Aliases, opts.Name)
		err = d.SaveArtistAliases(ctx, a.ID, opts.Aliases, "MusicBrainz")
		if err != nil {
			return nil, err
		}
		artist.Aliases = opts.Aliases
	}
	return artist, nil
}

func (d *Psql) UpdateArtist(ctx context.Context, opts db.UpdateArtistOpts) error {
	l := logger.FromContext(ctx)
	if opts.ID == 0 {
		return errors.New("artist id not specified")
	}
	tx, err := d.conn.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		l.Err(err).Msg("Failed to begin transaction")
		return err
	}
	defer tx.Rollback(ctx)
	qtx := d.q.WithTx(tx)
	if opts.MusicBrainzID != uuid.Nil {
		l.Debug().Msgf("Updating artist with id %d with MusicBrainz ID %s", opts.ID, opts.MusicBrainzID)
		err := qtx.UpdateArtistMbzID(ctx, repository.UpdateArtistMbzIDParams{
			ID:            opts.ID,
			MusicBrainzID: &opts.MusicBrainzID,
		})
		if err != nil {
			return err
		}
	}
	if opts.Image != uuid.Nil {
		l.Debug().Msgf("Updating artist with id %d with image %s", opts.ID, opts.Image)
		err = qtx.UpdateArtistImage(ctx, repository.UpdateArtistImageParams{
			ID:          opts.ID,
			Image:       &opts.Image,
			ImageSource: pgtype.Text{String: opts.ImageSrc, Valid: opts.ImageSrc != ""},
		})
		if err != nil {
			return err
		}
	}
	return tx.Commit(ctx)
}

func (d *Psql) DeleteArtistAlias(ctx context.Context, id int32, alias string) error {
	return d.q.DeleteArtistAlias(ctx, repository.DeleteArtistAliasParams{
		ArtistID: id,
		Alias:    alias,
	})
}
func (d *Psql) GetAllArtistAliases(ctx context.Context, id int32) ([]models.Alias, error) {
	rows, err := d.q.GetAllArtistAliases(ctx, id)
	if err != nil {
		return nil, err
	}
	aliases := make([]models.Alias, len(rows))
	for i, row := range rows {
		aliases[i] = models.Alias{
			ID:      id,
			Alias:   row.Alias,
			Source:  row.Source,
			Primary: row.IsPrimary,
		}
	}
	return aliases, nil
}

func (d *Psql) SetPrimaryArtistAlias(ctx context.Context, id int32, alias string) error {
	l := logger.FromContext(ctx)
	if id == 0 {
		return errors.New("artist id not specified")
	}
	tx, err := d.conn.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		l.Err(err).Msg("Failed to begin transaction")
		return err
	}
	defer tx.Rollback(ctx)
	qtx := d.q.WithTx(tx)
	// get all aliases
	aliases, err := qtx.GetAllArtistAliases(ctx, id)
	if err != nil {
		return err
	}
	primary := ""
	exists := false
	for _, v := range aliases {
		if v.Alias == alias {
			exists = true
		}
		if v.IsPrimary {
			primary = v.Alias
		}
	}
	if primary == alias {
		// no-op rename
		return nil
	}
	if !exists {
		return errors.New("alias does not exist")
	}
	err = qtx.SetArtistAliasPrimaryStatus(ctx, repository.SetArtistAliasPrimaryStatusParams{
		ArtistID:  id,
		Alias:     alias,
		IsPrimary: true,
	})
	if err != nil {
		return err
	}
	err = qtx.SetArtistAliasPrimaryStatus(ctx, repository.SetArtistAliasPrimaryStatusParams{
		ArtistID:  id,
		Alias:     primary,
		IsPrimary: false,
	})
	if err != nil {
		return err
	}
	return tx.Commit(ctx)
}
