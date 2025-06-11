package mbz

import (
	"context"

	"github.com/google/uuid"
)

type MusicBrainzTrack struct {
	Title string `json:"title"`
}

const recordingFmtStr = "%s/ws/2/recording/%s"

// Returns the artist name at index 0, and all primary aliases after.
func (c *MusicBrainzClient) GetTrack(ctx context.Context, id uuid.UUID) (*MusicBrainzTrack, error) {
	track := new(MusicBrainzTrack)
	err := c.getEntity(ctx, recordingFmtStr, id, track)
	if err != nil {
		return nil, err
	}
	return track, nil
}
