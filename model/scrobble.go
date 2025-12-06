package model

import "time"

type Scrobble struct {
	MediaFileID string
	UserID      string

	SubmissionTime time.Time `structs:"submission_time" json:"submissionTime"`
	RowId          int64     `structs:"row_id" json:"rowId"`
	MediaFile
}

type ScrobbleRepository interface {
	RecordScrobble(mediaFileID string, submissionTime time.Time) error
}

type Scrobbles []Scrobble
