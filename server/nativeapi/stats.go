package nativeapi

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/navidrome/navidrome/log"
	"github.com/navidrome/navidrome/model"
	"github.com/navidrome/navidrome/utils/req"
)

type statType uint

const (
	album = iota
	artist
	genre
	song
)

func (n *Router) getStats(stat statType) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		p := req.Params(r)

		from := p.TimeOr("from", time.Now().Add(-7*24*time.Hour))
		to := p.TimeOr("to", time.Now())
		start := p.IntOr("_start", 0)
		end := p.IntOr("_end", start+5)

		var data interface{}
		var err error

		ops := model.QueryOptions{
			Max:    end - start,
			Offset: start,
			Order:  "count DESC",
		}

		switch stat {
		case album:
			data, err = n.ds.Stat(ctx).AlbumStats(from, to, ops)
		case artist:
			data, err = n.ds.Stat(ctx).ArtistStats(from, to, ops)
		case genre:
			data, err = n.ds.Stat(ctx).GenreStats(from, to, ops)
		case song:
			data, err = n.ds.Stat(ctx).SongStats(from, to, ops)
		}

		if err != nil {
			log.Error(ctx, "Error getting media stats", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		replyJson(ctx, w, data)
	}
}

func (n *Router) stats(r chi.Router) {
	r.Route("/stats", func(r chi.Router) {
		r.Get("/album", n.getStats(album))
		r.Get("/artist", n.getStats(artist))
		r.Get("/genre", n.getStats(genre))
		r.Get("/song", n.getStats(song))
	})
}
