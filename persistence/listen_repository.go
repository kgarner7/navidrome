package persistence

import (
	"context"

	. "github.com/Masterminds/squirrel"
	"github.com/deluan/rest"
	"github.com/navidrome/navidrome/conf"
	"github.com/navidrome/navidrome/consts"
	"github.com/navidrome/navidrome/model"
	"github.com/pocketbase/dbx"
)

type listenRepository struct {
	sqlRepository
}

func NewListenRepository(ctx context.Context, db dbx.Builder) *listenRepository {
	r := &listenRepository{}
	r.ctx = ctx
	r.db = db
	r.tableName = "scrobbles"
	r.registerModel(&model.Listen{}, map[string]filterFunc{})
	r.setSortMappings(map[string]string{
		"listened_at": "scrobbles.submission_time",
	})
	return r
}

func (r *listenRepository) Count(options ...rest.QueryOptions) (int64, error) {
	user := loggedUser(r.ctx)

	sel := r.newSelect().
		Columns("count(*) count").
		Join("media_file f on f.id = file_id").
		Where(Eq{"user_id": user.ID})

	sel = r.applyFilters(sel, r.parseRestOptions(r.ctx, options...))
	var res struct{ Count int64 }
	err := r.queryOne(sel, &res)
	return res.Count, err
}

func (r *listenRepository) Read(id string) (interface{}, error) {
	return nil, model.ErrNotFound
}

func (r *listenRepository) ReadAll(options ...rest.QueryOptions) (interface{}, error) {
	user := loggedUser(r.ctx)

	sel := r.newSelect(r.parseRestOptions(r.ctx, options...)).
		Columns("submission_time", "f.*").
		Join("media_file f on f.id = file_id").
		LeftJoin("annotation on ("+
			"annotation.item_id = f.id"+
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

	var listens model.Listens
	err := r.queryAll(sel, &listens)
	if err != nil {
		return nil, err
	}
	return listens, err
}

func (r *listenRepository) EntityName() string {
	return "listen"
}

func (r *listenRepository) NewInstance() interface{} {
	return &model.Listen{}
}

var _ model.ListenRepository = (*statRepository)(nil)
var _ model.ResourceRepository = (*listenRepository)(nil)
