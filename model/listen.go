package model

type Listen struct {
	RowId          int64 `structs:"row_id" json:"rowId"`
	SubmissionTime int64 `structs:"submission_time" json:"submissionTime"`
	MediaFile
}

type Listens []Listen

type ListenRepository interface {
}
