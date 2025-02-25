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

func (r *statRepository) Stats(statType model.StatType, from time.Time, to time.Time, ops ...model.QueryOptions) (model.Stats, error) {
	sel := r.baseSelect(from, to, ops...)

	switch statType {
	case model.AlbumStat:
		sel = sel.Columns("a.id", "a.name").
			Join("album a on a.id = f.album_id").
			GroupBy("a.id")
	case model.ArtistStat:
		sel = sel.Columns("a.id", "a.name").
			Join("artist a on a.id = f.artist_id").
			GroupBy("a.id")
	case model.GenreStat:
		sel = sel.From("scrobbles, json_each(f.tags, '$.genre')").
			Columns("json_extract(value, '$.id') id", "json_extract(value, '$.value') name").
			GroupBy("name")
	case model.SongStat:
		sel = sel.Columns("f.id", "f.title name").
			GroupBy("f.id")
	}

	var stat model.Stats
	err := r.queryAll(sel, &stat)

	return stat, err
}

func (r *statRepository) StatsCount(statType model.StatType, from time.Time, to time.Time) (int64, error) {
	sel := r.baseSelect(from, to).RemoveColumns()

	switch statType {
	case model.AlbumStat:
		sel = sel.Join("album a on a.id = f.album_id").
			Column("count(distinct a.id) count")
	case model.ArtistStat:
		sel = sel.Join("artist a on a.id = f.artist_id").
			Column("count(distinct a.id) count")
	case model.GenreStat:
		sel = sel.From("scrobbles, json_each(f.tags, '$.genre')").
			Column("count(distinct json_extract(value, '$.id')) count")
	case model.SongStat:
		sel = sel.Column("count(distinct f.id) count")
	}

	var res struct{ Count int64 }
	err := r.queryOne(sel, &res)

	return res.Count, err
}

func (r *statRepository) RecordPlay(id string, ts time.Time) error {
	userId := userId(r.ctx)
	insert := Insert(r.tableName).Columns("file_id", "user_id", "submission_time").Values(id, userId, ts.Unix())
	_, err := r.executeSQL(insert)
	return err
}

var _ model.StatRepository = (*statRepository)(nil)
