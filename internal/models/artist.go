package models

import "github.com/google/uuid"

type Artist struct {
	ID           int32      `json:"id"`
	MbzID        *uuid.UUID `json:"musicbrainz_id"`
	Name         string     `json:"name"`
	Aliases      []string   `json:"aliases"`
	Image        *uuid.UUID `json:"image"`
	ListenCount  int64      `json:"listen_count"`
	TimeListened int64      `json:"time_listened"`
	IsPrimary    bool       `json:"is_primary,omitempty"`
}

type SimpleArtist struct {
	ID   int32  `json:"id"`
	Name string `json:"name"`
}
