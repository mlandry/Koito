package psql

import (
	"context"
	"time"

	"github.com/gabehf/koito/internal/db"
	"github.com/gabehf/koito/internal/repository"
)

func (p *Psql) CountListens(ctx context.Context, period db.Period) (int64, error) {
	t2 := time.Now()
	t1 := db.StartTimeFromPeriod(period)
	count, err := p.q.CountListens(ctx, repository.CountListensParams{
		ListenedAt:   t1,
		ListenedAt_2: t2,
	})
	if err != nil {
		return 0, err
	}
	return count, nil
}
func (p *Psql) CountTracks(ctx context.Context, period db.Period) (int64, error) {
	t2 := time.Now()
	t1 := db.StartTimeFromPeriod(period)
	count, err := p.q.CountTopTracks(ctx, repository.CountTopTracksParams{
		ListenedAt:   t1,
		ListenedAt_2: t2,
	})
	if err != nil {
		return 0, err
	}
	return count, nil
}
func (p *Psql) CountAlbums(ctx context.Context, period db.Period) (int64, error) {
	t2 := time.Now()
	t1 := db.StartTimeFromPeriod(period)
	count, err := p.q.CountTopReleases(ctx, repository.CountTopReleasesParams{
		ListenedAt:   t1,
		ListenedAt_2: t2,
	})
	if err != nil {
		return 0, err
	}
	return count, nil
}
func (p *Psql) CountArtists(ctx context.Context, period db.Period) (int64, error) {
	t2 := time.Now()
	t1 := db.StartTimeFromPeriod(period)
	count, err := p.q.CountTopArtists(ctx, repository.CountTopArtistsParams{
		ListenedAt:   t1,
		ListenedAt_2: t2,
	})
	if err != nil {
		return 0, err
	}
	return count, nil
}
func (p *Psql) CountTimeListened(ctx context.Context, period db.Period) (int64, error) {
	t2 := time.Now()
	t1 := db.StartTimeFromPeriod(period)
	count, err := p.q.CountTimeListened(ctx, repository.CountTimeListenedParams{
		ListenedAt:   t1,
		ListenedAt_2: t2,
	})
	if err != nil {
		return 0, err
	}
	return count, nil
}
