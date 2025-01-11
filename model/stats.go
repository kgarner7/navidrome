package model

import "time"

type Stat struct {
	Count uint64 `structs:"count" json:"count"`
}

type AlbumStat struct {
	Stat
	Album `structs:"-"`
}

type AlbumStats []AlbumStat

type ArtistStat struct {
	Stat
	Artist `structs:"-"`
}

type ArtistStats []ArtistStat

type GenreStat struct {
	Stat
	Id   string `structs:"id" json:"id"`
	Name string `structs:"name" json:"name"`
}

type GenreStats []GenreStat

type SongStat struct {
	Stat
	MediaFile `structs:"-"`
}

type SongStats []SongStat

type StatRepository interface {
	AlbumStats(from time.Time, to time.Time, ops ...QueryOptions) (AlbumStats, error)
	ArtistStats(from time.Time, to time.Time, ops ...QueryOptions) (ArtistStats, error)
	GenreStats(from time.Time, to time.Time, ops ...QueryOptions) (GenreStats, error)
	SongStats(from time.Time, to time.Time, ops ...QueryOptions) (SongStats, error)
	RecordPlay(id string, ts time.Time) error
}
