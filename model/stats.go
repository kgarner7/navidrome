package model

import "time"

type Stat struct {
	Count uint64 `structs:"count" json:"count"`
	Id    string `structs:"id" json:"id"`
	Name  string `structs:"name" json:"name"`
}

type Stats []Stat

type StatRepository interface {
	AlbumStats(from time.Time, to time.Time, ops ...QueryOptions) (Stats, error)
	ArtistStats(from time.Time, to time.Time, ops ...QueryOptions) (Stats, error)
	GenreStats(from time.Time, to time.Time, ops ...QueryOptions) (Stats, error)
	SongStats(from time.Time, to time.Time, ops ...QueryOptions) (Stats, error)
	RecordPlay(id string, ts time.Time) error
}
