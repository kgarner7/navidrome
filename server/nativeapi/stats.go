package nativeapi

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/navidrome/navidrome/log"
	"github.com/navidrome/navidrome/model"
	"github.com/navidrome/navidrome/utils/req"
)

func (n *Router) getStats() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		p := req.Params(r)

		typeString := p.StringOr("stat", "")

		var stat model.StatType

		switch typeString {
		case "album":
			stat = model.AlbumStat
		case "artist":
			stat = model.ArtistStat
		case "genre":
			stat = model.GenreStat
		case "song":
			stat = model.SongStat
		case "total":

		default:
			http.Error(w, "Invalid stat type", http.StatusBadRequest)
		}

		from := p.TimeOr("from", time.Now().Add(-7*24*time.Hour))
		to := p.TimeOr("to", time.Now())
		start := p.IntOr("_start", 0)
		end := p.IntOr("_end", start+5)

		var count int64
		var data any
		var err error

		if typeString == "total" {
			data, err = n.ds.Stat(ctx).Total(from, to)

			if err != nil {
				log.Error(ctx, "Error getting aggregate stats", err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			count = 1
		} else {
			ops := model.QueryOptions{
				Max:    end - start,
				Offset: start,
			}

			data, err = n.ds.Stat(ctx).Stats(stat, from, to, ops)

			if err != nil {
				log.Error(ctx, "Error getting media stats", err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			count, err = n.ds.Stat(ctx).StatsCount(stat, from, to)
			if err != nil {
				log.Error(ctx, "Error getting media count", err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}

		w.Header().Set("X-Total-Count", strconv.FormatInt(count, 10))

		replyJson(ctx, w, data)
	}
}

func replyJson(ctx context.Context, w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	resp, _ := json.Marshal(data)
	_, err := w.Write(resp)

	if err != nil {
		log.Error(ctx, "Error sending json", "Error", err)
	}
}

func (n *Router) stats(r chi.Router) {
	r.Get("/stats", n.getStats())
}
