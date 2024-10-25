package subsonic

import (
	"errors"
	"net/http"
	"time"

	"github.com/navidrome/navidrome/model"
	"github.com/navidrome/navidrome/model/request"
	"github.com/navidrome/navidrome/server/subsonic/responses"
	"github.com/navidrome/navidrome/utils/req"
	"github.com/navidrome/navidrome/utils/slice"
)

func (api *Router) GetBookmarks(r *http.Request) (*responses.Subsonic, error) {
	user, _ := request.UserFrom(r.Context())

	repo := api.ds.MediaFile(r.Context())
	bookmarks, err := repo.GetBookmarks()
	if err != nil {
		return nil, err
	}

	response := newResponse()
	response.Bookmarks = &responses.Bookmarks{}
	response.Bookmarks.Bookmark = slice.Map(bookmarks, func(bmk model.Bookmark) responses.Bookmark {
		return responses.Bookmark{
			Entry:    childFromMediaFile(r.Context(), bmk.Item),
			Position: bmk.Position,
			Username: user.UserName,
			Comment:  bmk.Comment,
			Created:  bmk.CreatedAt,
			Changed:  bmk.UpdatedAt,
		}
	})
	return response, nil
}

func (api *Router) CreateBookmark(r *http.Request) (*responses.Subsonic, error) {
	p := req.Params(r)
	id, err := p.String("id")
	if err != nil {
		return nil, err
	}

	comment, _ := p.String("comment")
	position := p.Int64Or("position", 0)

	repo := api.ds.MediaFile(r.Context())
	err = repo.AddBookmark(id, comment, position)
	if err != nil {
		return nil, err
	}
	return newResponse(), nil
}

func (api *Router) DeleteBookmark(r *http.Request) (*responses.Subsonic, error) {
	p := req.Params(r)
	id, err := p.String("id")
	if err != nil {
		return nil, err
	}

	repo := api.ds.MediaFile(r.Context())
	err = repo.DeleteBookmark(id)
	if err != nil {
		return nil, err
	}
	return newResponse(), nil
}

func (api *Router) GetPlayQueue(r *http.Request) (*responses.Subsonic, error) {
	user, _ := request.UserFrom(r.Context())

	repo := api.ds.PlayQueue(r.Context())
	pq, err := repo.Retrieve(user.ID)
	if err != nil && !errors.Is(err, model.ErrNotFound) {
		return nil, err
	}
	if pq == nil || len(pq.Items) == 0 {
		return newResponse(), nil
	}

	response := newResponse()
	response.PlayQueue = &responses.PlayQueue{
		Entry:     slice.MapWithArg(pq.Items, r.Context(), childFromMediaFile),
		Current:   pq.Current,
		Position:  pq.Position,
		Username:  user.UserName,
		Changed:   &pq.UpdatedAt,
		ChangedBy: pq.ChangedBy,
	}
	return response, nil
}

func (api *Router) SavePlayQueue(r *http.Request) (*responses.Subsonic, error) {
	p := req.Params(r)
	ids, _ := p.Strings("id")
	current, _ := p.String("current")
	position := p.Int64Or("position", 0)

	user, _ := request.UserFrom(r.Context())
	client, _ := request.ClientFrom(r.Context())

	var items model.MediaFiles
	for _, id := range ids {
		items = append(items, model.MediaFile{ID: id})
	}

	pq := &model.PlayQueue{
		UserID:    user.ID,
		Current:   current,
		Position:  position,
		ChangedBy: client,
		Items:     items,
		CreatedAt: time.Time{},
		UpdatedAt: time.Time{},
	}

	repo := api.ds.PlayQueue(r.Context())
	err := repo.Store(pq)
	if err != nil {
		return nil, err
	}
	return newResponse(), nil
}

func (api *Router) GetPlayQueueAdvanced(r *http.Request) (*responses.Subsonic, error) {
	user, _ := request.UserFrom(r.Context())

	repo := api.ds.PlayQueue(r.Context())
	pq, err := repo.Retrieve(user.ID)
	if err != nil {
		return nil, err
	}

	response := newResponse()
	response.PlayQueue2 = &responses.PlayQueue2{
		Entry:      childrenFromMediaFiles(r.Context(), pq.Items),
		Current:    pq.Current,
		Position:   pq.Position,
		QueueIndex: pq.QueueIndex,
		Username:   user.UserName,
		Changed:    &pq.UpdatedAt,
		ChangedBy:  pq.ChangedBy,
	}
	return response, nil
}

func (api *Router) SavePlayQueueAdvanced(r *http.Request) (*responses.Subsonic, error) {
	p := req.Params(r)
	ids, _ := p.Strings("id")
	queueIdx := p.Int64Or("index", 0)

	if queueIdx < 0 || (len(ids) > 0 && queueIdx > int64(len(ids))) {
		return nil, newError(responses.ErrorGeneric, "Index cannot exceed length of queue")
	}

	user, _ := request.UserFrom(r.Context())
	client, _ := request.ClientFrom(r.Context())

	if len(ids) > 0 {
		position := p.Int64Or("position", 0)

		var items model.MediaFiles
		for _, id := range ids {
			items = append(items, model.MediaFile{ID: id})
		}

		var current = ""
		if queueIdx > 0 {
			current = ids[queueIdx-1]
		}

		pq := &model.PlayQueue{
			UserID:     user.ID,
			Current:    current,
			QueueIndex: queueIdx,
			Position:   position,
			ChangedBy:  client,
			Items:      items,
			CreatedAt:  time.Time{},
			UpdatedAt:  time.Time{},
		}

		repo := api.ds.PlayQueue(r.Context())
		err := repo.Store(pq)
		if err != nil {
			return nil, err
		}
	}

	err := api.ds.WithTx(func(tx model.DataStore) error {
		repo := tx.PlayQueue(r.Context())
		pq, err := repo.Get(user.ID)
		if err != nil {
			return err
		}

		if queueIdx > int64(len(pq.Items)) {
			return errors.New("position cannot exceed queue length")
		}

		if queueIdx != 0 {
			pq.QueueIndex = queueIdx
			pq.Current = pq.Items[queueIdx-1].ID
		}

		position := p.Int64Or("position", -1)

		if position != -1 {
			pq.Position = position
		}

		err = repo.Save(pq)
		return err
	})

	if err != nil {
		return nil, err
	}
	return newResponse(), nil
}
