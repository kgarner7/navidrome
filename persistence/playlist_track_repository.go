package persistence

import (
	"database/sql"

	. "github.com/Masterminds/squirrel"
	"github.com/deluan/rest"
	"github.com/navidrome/navidrome/conf"
	"github.com/navidrome/navidrome/log"
	"github.com/navidrome/navidrome/model"
	"github.com/navidrome/navidrome/utils/slice"
)

type playlistTrackRepository struct {
	sqlRepository
	sqlRestful
	playlistId   string
	playlist     *model.Playlist
	playlistRepo *playlistRepository
}

func (r *playlistRepository) Tracks(playlistId string, refreshSmartPlaylist bool) model.PlaylistTrackRepository {
	p := &playlistTrackRepository{}
	p.playlistRepo = r
	p.playlistId = playlistId
	p.ctx = r.ctx
	p.db = r.db
	p.tableName = "playlist_tracks"
	p.sortMappings = map[string]string{
		"id": "playlist_tracks.id",
	}
	pls, err := r.Get(playlistId)
	if err != nil {
		log.Error(r.ctx, "Error getting playlist's tracks - THIS SHOULD NOT HAPPEN!", "playlistId", playlistId, err)
		return nil
	}
	if refreshSmartPlaylist {
		r.refreshSmartPlaylist(pls)
	}
	p.playlist = pls
	return p
}

func (r *playlistTrackRepository) Count(options ...rest.QueryOptions) (int64, error) {
	if conf.Server.EnableDuplicateSearch {
		if duplicate, ok := options[0].Filters["duplicate"]; ok {
			delete(options[0].Filters, "duplicate")

			var builder = Select().
				Where(Eq{"playlist_id": r.playlistId})

			if duplicate == "true" {
				builder = builder.PrefixExpr(r.duplicate_cte()).
					Where("playlist_tracks.id NOT IN non_duplicate AND playlist_tracks.media_file_id in duplicate")
			}

			return r.count(builder, r.parseRestOptions(options...))
		}
	}

	return r.count(Select().Where(Eq{"playlist_id": r.playlistId}), r.parseRestOptions(options...))
}

func (r *playlistTrackRepository) Read(id string) (interface{}, error) {
	sel := r.newSelect().
		LeftJoin("annotation on ("+
			"annotation.item_id = media_file_id"+
			" AND annotation.item_type = 'media_file'"+
			" AND annotation.user_id = '"+userId(r.ctx)+"')").
		Columns(
			"coalesce(starred, 0)",
			"coalesce(play_count, 0)",
			"coalesce(rating, 0)",
			"starred_at",
			"play_date",
			"f.*",
			"playlist_tracks.*",
		).
		Join("media_file f on f.id = media_file_id").
		Where(And{Eq{"playlist_id": r.playlistId}, Eq{"id": id}})
	var trk model.PlaylistTrack
	err := r.queryOne(sel, &trk)
	return &trk, err
}

func (r *playlistTrackRepository) GetAll(options ...model.QueryOptions) (model.PlaylistTracks, error) {
	tracks, err := r.playlistRepo.loadTracks(r.newSelect(options...), r.playlistId)
	if err != nil {
		return nil, err
	}
	mfs := tracks.MediaFiles()
	err = r.loadMediaFileGenres(&mfs)
	if err != nil {
		log.Error(r.ctx, "Error loading genres for playlist", "playlist", r.playlist.Name, "id", r.playlist.ID, err)
		return nil, err
	}
	for i, mf := range mfs {
		tracks[i].MediaFile.Genres = mf.Genres
	}
	return tracks, err
}

func (r *playlistTrackRepository) duplicate_cte() SelectBuilder {
	return Select().
		Prefix("WITH RECURSIVE cte_duplicate(media, original_id) AS (").
		From(r.tableName).
		Columns("playlist_tracks.media_file_id", "playlist_tracks.id").
		Where(Eq{"playlist_id": r.playlistId}).
		OrderBy("playlist_tracks.id").
		GroupBy("playlist_tracks.media_file_id").
		Having("COUNT(*) > 1").
		Suffix("),").
		SuffixExpr(
			Select().
				Prefix("non_duplicate(id) AS (").
				From("cte_duplicate").
				Column("original_id").
				Suffix("),")).
		SuffixExpr(
			Select().
				Prefix("duplicate(file) AS (").
				From("cte_duplicate").
				Column("media").
				Suffix(")"))
}

func (r *playlistTrackRepository) GetAllShowingDuplicates(duplicate bool, options ...model.QueryOptions) (model.PlaylistTracks, error) {
	var builder = r.newSelect(options...).
		PrefixExpr(r.duplicate_cte()).
		Column("playlist_tracks.id NOT IN non_duplicate AND playlist_tracks.media_file_id in duplicate duplicate", r.playlistId)

	if duplicate {
		builder = builder.Where("duplicate")
	}

	tracks, err := r.playlistRepo.loadTracks(builder, r.playlistId)

	if err != nil {
		return nil, err
	}

	mfs := tracks.MediaFiles()
	err = r.loadMediaFileGenres(&mfs)
	if err != nil {
		log.Error(r.ctx, "Error loading genres for playlist", "playlist", r.playlist.Name, "id", r.playlist.ID, err)
		return nil, err
	}
	for i, mf := range mfs {
		tracks[i].MediaFile.Genres = mf.Genres
	}

	return tracks, err
}

func (r *playlistTrackRepository) GetAlbumIDs(options ...model.QueryOptions) ([]string, error) {
	sql := r.newSelect(options...).Columns("distinct mf.album_id").
		Join("media_file mf on mf.id = media_file_id").
		Where(Eq{"playlist_id": r.playlistId})
	var ids []string
	err := r.queryAllSlice(sql, &ids)
	if err != nil {
		return nil, err
	}
	return ids, nil
}

func (r *playlistTrackRepository) ReadAll(options ...rest.QueryOptions) (interface{}, error) {
	if conf.Server.EnableDuplicateSearch {
		if duplicate, ok := options[0].Filters["duplicate"]; ok {
			delete(options[0].Filters, "duplicate")

			result, err := r.GetAllShowingDuplicates(duplicate == "true", r.parseRestOptions(options...))

			options[0].Filters["duplicate"] = duplicate

			return result, err
		} else {
			return r.GetAllShowingDuplicates(false, r.parseRestOptions(options...))
		}
	} else {
		return r.GetAll(r.parseRestOptions(options...))
	}
}

func (r *playlistTrackRepository) EntityName() string {
	return "playlist_tracks"
}

func (r *playlistTrackRepository) NewInstance() interface{} {
	return &model.PlaylistTrack{}
}

func (r *playlistTrackRepository) isTracksEditable() bool {
	return r.playlistRepo.isWritable(r.playlistId) && !r.playlist.IsSmartPlaylist()
}

func (r *playlistTrackRepository) Add(mediaFileIds []string) (int, error) {
	if !r.isTracksEditable() {
		return 0, rest.ErrPermissionDenied
	}

	if len(mediaFileIds) > 0 {
		log.Debug(r.ctx, "Adding songs to playlist", "playlistId", r.playlistId, "mediaFileIds", mediaFileIds)
	} else {
		return 0, nil
	}

	// Get next pos (ID) in playlist
	sq := r.newSelect().Columns("max(id) as max").Where(Eq{"playlist_id": r.playlistId})
	var res struct{ Max sql.NullInt32 }
	err := r.queryOne(sq, &res)
	if err != nil {
		return 0, err
	}

	return len(mediaFileIds), r.playlistRepo.addTracks(r.playlistId, int(res.Max.Int32+1), mediaFileIds)
}

func (r *playlistTrackRepository) addMediaFileIds(cond Sqlizer) (int, error) {
	sq := Select("id").From("media_file").Where(cond).OrderBy("album_artist, album, release_date, disc_number, track_number")
	var ids []string
	err := r.queryAllSlice(sq, &ids)
	if err != nil {
		log.Error(r.ctx, "Error getting tracks to add to playlist", err)
		return 0, err
	}
	return r.Add(ids)
}

func (r *playlistTrackRepository) AddAlbums(albumIds []string) (int, error) {
	return r.addMediaFileIds(Eq{"album_id": albumIds})
}

func (r *playlistTrackRepository) AddArtists(artistIds []string) (int, error) {
	return r.addMediaFileIds(Eq{"album_artist_id": artistIds})
}

func (r *playlistTrackRepository) AddDiscs(discs []model.DiscID) (int, error) {
	if len(discs) == 0 {
		return 0, nil
	}
	var clauses Or
	for _, d := range discs {
		clauses = append(clauses, And{Eq{"album_id": d.AlbumID}, Eq{"release_date": d.ReleaseDate}, Eq{"disc_number": d.DiscNumber}})
	}
	return r.addMediaFileIds(clauses)
}

// Get ids from all current tracks
func (r *playlistTrackRepository) getTracks() ([]string, error) {
	all := r.newSelect().Columns("media_file_id").Where(Eq{"playlist_id": r.playlistId}).OrderBy("id")
	var ids []string
	err := r.queryAllSlice(all, &ids)
	if err != nil {
		log.Error(r.ctx, "Error querying current tracks from playlist", "playlistId", r.playlistId, err)
		return nil, err
	}
	return ids, nil
}

func (r *playlistTrackRepository) Delete(ids ...string) error {
	if !r.isTracksEditable() {
		return rest.ErrPermissionDenied
	}
	err := r.delete(And{Eq{"playlist_id": r.playlistId}, Eq{"id": ids}})
	if err != nil {
		return err
	}

	return r.playlistRepo.renumber(r.playlistId)
}

func (r *playlistTrackRepository) DeleteAll() error {
	if !r.isTracksEditable() {
		return rest.ErrPermissionDenied
	}
	err := r.delete(Eq{"playlist_id": r.playlistId})
	if err != nil {
		return err
	}

	return r.playlistRepo.renumber(r.playlistId)
}

func (r *playlistTrackRepository) Reorder(pos int, newPos int) error {
	if !r.isTracksEditable() {
		return rest.ErrPermissionDenied
	}
	ids, err := r.getTracks()
	if err != nil {
		return err
	}
	newOrder := slice.Move(ids, pos-1, newPos-1)
	return r.playlistRepo.updatePlaylist(r.playlistId, newOrder)
}

var _ model.PlaylistTrackRepository = (*playlistTrackRepository)(nil)
