package model

import "time"

type Stat struct {
	Count uint64 `structs:"count" json:"count"`
	Id    string `structs:"id" json:"id"`
	Name  string `structs:"name" json:"name"`
}

type Stats []Stat

type StatType string

const (
	AlbumStat  StatType = "album"
	ArtistStat StatType = "artist"
	GenreStat  StatType = "genre"
	SongStat   StatType = "song"
)

type StatRepository interface {
	Stats(statType StatType, from time.Time, to time.Time, ops ...QueryOptions) (Stats, error)
	StatsCount(statType StatType, from time.Time, to time.Time) (int64, error)
	RecordPlay(id string, ts time.Time) error
}
