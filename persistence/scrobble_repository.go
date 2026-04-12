package persistence

import (
	"context"
	"time"

	. "github.com/Masterminds/squirrel"
	"github.com/deluan/rest"
	"github.com/navidrome/navidrome/conf"
	"github.com/navidrome/navidrome/consts"
	"github.com/navidrome/navidrome/model"
	"github.com/navidrome/navidrome/utils/slice"
	"github.com/pocketbase/dbx"
)

type scrobbleRepository struct {
	sqlRepository
}

type dbScrobble struct {
	dbMediaFile
	RowId          int64 `structs:"row_id" json:"rowId"`
	SubmissionTime int64 `structs:"submission_time" json:"submissionTime"`
}

type dbScrobbles []dbScrobble

func (m dbScrobbles) toModels() model.Scrobbles {
	return slice.Map(m, func(db dbScrobble) model.Scrobble {
		return model.Scrobble{
			MediaFile:      *db.MediaFile,
			RowId:          db.RowId,
			SubmissionTime: time.Unix(db.SubmissionTime, 0),
		}
	})
}

func NewScrobbleRepository(ctx context.Context, db dbx.Builder) model.ScrobbleRepository {
	r := &scrobbleRepository{}
	r.ctx = ctx
	r.db = db
	r.tableName = "scrobbles"
	r.registerModel(&model.Scrobble{}, map[string]filterFunc{
		"title": fullTextFilter("media_file"),
	})
	r.setSortMappings(map[string]string{
		"listened_at":  "scrobbles.submission_time",
		"title":        "order_title, scrobbles.submission_time",
		"artist":       "order_artist_name, order_album_name, release_date, disc_number, track_number, scrobbles.submission_time",
		"album":        "order_album_name, release_date, disc_number, track_number, order_artist_name, title, scrobbles.submission_time",
		"album_artist": "compilation, order_album_artist_name, order_album_name, scrobbles.submission_time",
		"duration":     "duration, scrobbles.submission_time",
	})
	return r
}

func (r *scrobbleRepository) RecordScrobble(mediaFileID string, submissionTime time.Time) error {
	userID := loggedUser(r.ctx).ID
	values := map[string]interface{}{
		"file_id":         mediaFileID,
		"user_id":         userID,
		"submission_time": submissionTime.Unix(),
	}
	insert := Insert(r.tableName).SetMap(values)
	_, err := r.executeSQL(insert)
	return err
}

func (r *scrobbleRepository) Count(options ...rest.QueryOptions) (int64, error) {
	user := loggedUser(r.ctx)

	sel := r.newSelect().
		Columns("count(*) count").
		Join("media_file on media_file.id = file_id").
		Where(Eq{"user_id": user.ID})

	sel = r.applyFilters(sel, r.parseRestOptions(r.ctx, options...))
	var res struct{ Count int64 }
	err := r.queryOne(sel, &res)
	return res.Count, err
}

func (r *scrobbleRepository) Read(id string) (interface{}, error) {
	return nil, model.ErrNotFound
}

func (r *scrobbleRepository) ReadAll(options ...rest.QueryOptions) (interface{}, error) {
	user := loggedUser(r.ctx)

	sel := r.newSelect(r.parseRestOptions(r.ctx, options...)).
		Columns("scrobbles.ROWID row_id", "submission_time", "media_file.*").
		Join("media_file on media_file.id = file_id").
		LeftJoin("annotation on ("+
			"annotation.item_id = media_file.id"+
			" AND annotation.item_type = 'media_file'"+
			" AND annotation.user_id = '"+user.ID+"')").
		Columns(
			"coalesce(starred, 0) as starred",
			"coalesce(rating, 0) as rating",
			"starred_at",
			"play_date",
		).
		Where(Eq{"scrobbles.user_id": user.ID})

	if conf.Server.AlbumPlayCountMode == consts.AlbumPlayCountModeNormalized && r.tableName == "album" {
		sel = sel.Columns("round(coalesce(round(cast(play_count as float) / coalesce(song_count, 1), 1), 0)) as play_count")
	} else {
		sel = sel.Columns("coalesce(play_count, 0) as play_count")
	}

	var scrobbles dbScrobbles
	err := r.queryAll(sel, &scrobbles)
	if err != nil {
		return nil, err
	}
	return scrobbles.toModels(), err
}

func (r *scrobbleRepository) EntityName() string {
	return "scrobble"
}

func (r *scrobbleRepository) NewInstance() interface{} {
	return &model.Scrobble{}
}

var _ model.ScrobbleRepository = (*scrobbleRepository)(nil)
var _ model.ResourceRepository = (*scrobbleRepository)(nil)
