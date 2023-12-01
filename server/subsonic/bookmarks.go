package subsonic

import (
	"errors"
	"net/http"
	"time"

	"github.com/navidrome/navidrome/model"
	"github.com/navidrome/navidrome/model/request"
	"github.com/navidrome/navidrome/server/subsonic/responses"
	"github.com/navidrome/navidrome/utils"
)

func (api *Router) GetBookmarks(r *http.Request) (*responses.Subsonic, error) {
	user, _ := request.UserFrom(r.Context())

	repo := api.ds.MediaFile(r.Context())
	bmks, err := repo.GetBookmarks()
	if err != nil {
		return nil, err
	}

	response := newResponse()
	response.Bookmarks = &responses.Bookmarks{}
	for _, bmk := range bmks {
		b := responses.Bookmark{
			Entry:    childFromMediaFile(r.Context(), bmk.Item),
			Position: bmk.Position,
			Username: user.UserName,
			Comment:  bmk.Comment,
			Created:  bmk.CreatedAt,
			Changed:  bmk.UpdatedAt,
		}
		response.Bookmarks.Bookmark = append(response.Bookmarks.Bookmark, b)
	}
	return response, nil
}

func (api *Router) CreateBookmark(r *http.Request) (*responses.Subsonic, error) {
	id, err := requiredParamString(r, "id")
	if err != nil {
		return nil, err
	}

	comment := utils.ParamString(r, "comment")
	position := utils.ParamInt(r, "position", int64(0))

	repo := api.ds.MediaFile(r.Context())
	err = repo.AddBookmark(id, comment, position)
	if err != nil {
		return nil, err
	}
	return newResponse(), nil
}

func (api *Router) DeleteBookmark(r *http.Request) (*responses.Subsonic, error) {
	id, err := requiredParamString(r, "id")
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
	if err != nil {
		return nil, err
	}

	response := newResponse()
	response.PlayQueue = &responses.PlayQueue{
		Entry:     childrenFromMediaFiles(r.Context(), pq.Items),
		Current:   pq.Current,
		Position:  pq.Position,
		Username:  user.UserName,
		Changed:   &pq.UpdatedAt,
		ChangedBy: pq.ChangedBy,
	}
	return response, nil
}

func (api *Router) SavePlayQueue(r *http.Request) (*responses.Subsonic, error) {
	ids, err := requiredParamStrings(r, "id")
	if err != nil {
		return nil, err
	}

	current := utils.ParamString(r, "current")
	position := utils.ParamInt(r, "position", int64(0))

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
	err = repo.Store(pq)
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
	ids := utils.BodyStrings(r, "id")
	queueIdx := utils.BodyInt64(r, "index", 0)

	if queueIdx < 0 || (len(ids) > 0 && queueIdx > int64(len(ids))) {
		return nil, newError(responses.ErrorGeneric, "Index cannot exceed length of queue")
	}

	user, _ := request.UserFrom(r.Context())
	client, _ := request.ClientFrom(r.Context())

	if len(ids) > 0 {
		position := utils.BodyInt64(r, "position", 0)

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

		position := utils.BodyInt64(r, "position", -1)

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
