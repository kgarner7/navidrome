package persistence

import (
	"context"
	"time"

	. "github.com/Masterminds/squirrel"
	"github.com/navidrome/navidrome/model"
	"github.com/pocketbase/dbx"
)

type statRepository struct {
	sqlRepository
}

func NewStatRepository(ctx context.Context, db dbx.Builder) *statRepository {
	r := &statRepository{}
	r.ctx = ctx
	r.db = db
	r.tableName = "scrobbles"
	return r
}

func (r *statRepository) baseSelect(from time.Time, to time.Time, ops ...model.QueryOptions) SelectBuilder {
	user := loggedUser(r.ctx)

	sel := r.newSelect(ops...).
		Column("COUNT(*) count").
		Join("media_file f on f.id = file_id").
		Where(And{
			GtOrEq{"submission_time": from.Unix()},
			LtOrEq{"submission_time": to.Unix()},
			Eq{"user_id": user.ID},
		}).
		OrderBy("count DESC")

	return sel
}

func (r *statRepository) AlbumStats(from time.Time, to time.Time, ops ...model.QueryOptions) (model.Stats, error) {
	sel := r.baseSelect(from, to, ops...).
		Columns("a.id", "a.name").
		Join("album a on a.id = f.album_id").
		GroupBy("a.id")

	var stat model.Stats
	err := r.queryAll(sel, &stat)

	return stat, err
}

func (r *statRepository) ArtistStats(from time.Time, to time.Time, ops ...model.QueryOptions) (model.Stats, error) {
	sel := r.baseSelect(from, to, ops...).
		Columns("a.id", "a.name").
		Join("artist a on a.id = f.artist_id").
		GroupBy("a.id")

	var stat model.Stats
	err := r.queryAll(sel, &stat)

	return stat, err
}

// GenreStats implements model.StatRepository.
func (r *statRepository) GenreStats(from time.Time, to time.Time, ops ...model.QueryOptions) (model.Stats, error) {
	sel := r.baseSelect(from, to, ops...).
		Columns("g.id", "g.name").
		Join("media_file_genres mg ON mg.media_file_id = f.id").
		Join("genre g ON g.id = mg.genre_id").
		GroupBy("g.id")

	var res model.Stats
	err := r.queryAll(sel, &res)
	return res, err
}

func (r *statRepository) SongStats(from time.Time, to time.Time, ops ...model.QueryOptions) (model.Stats, error) {
	sel := r.baseSelect(from, to, ops...).
		Columns("f.id", "f.title name").
		GroupBy("f.id")

	var res model.Stats
	err := r.queryAll(sel, &res)

	return res, err
}

func (r *statRepository) RecordPlay(id string, ts time.Time) error {
	userId := userId(r.ctx)
	insert := Insert(r.tableName).Columns("file_id", "user_id", "submission_time").Values(id, userId, ts.Unix())
	_, err := r.executeSQL(insert)
	return err
}

var _ model.StatRepository = (*statRepository)(nil)
