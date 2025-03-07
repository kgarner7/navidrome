package nativeapi

import (
	"encoding/json"
	"errors"
	"net/http"
	"os/exec"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/navidrome/navidrome/conf"
	"github.com/navidrome/navidrome/log"
	"github.com/navidrome/navidrome/model"
	"github.com/navidrome/navidrome/utils/req"
)

func (n *Router) instantMix() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		p := req.Params(r)

		id, err := p.String("id")
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		count := p.IntOr("count", 50)
		if count < 1 || count > 50 {
			count = 50
		}

		mfRepo := n.ds.MediaFile(ctx)

		mf, err := mfRepo.Get(id)
		if errors.Is(err, model.ErrNotFound) {
			log.Warn(ctx, "could not find file", "id", id)
			http.Error(w, "not found", http.StatusNotFound)
			return
		}

		blissPath := strings.Replace(mf.AbsolutePath(), conf.Server.Bliss.RemovePrefix, conf.Server.Bliss.PrependPrefix, 1)

		output, err := exec.Command(
			conf.Server.Bliss.Path,
			"playlist",
			"--playlist-length", strconv.Itoa(count),
			"--config-path", conf.Server.Bliss.ConfigPath,
			blissPath,
		).CombinedOutput()

		if err != nil {
			log.Error(ctx, "failed to do instant mix", "output", string(output), "error", err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		mediaFilePath := strings.Split(strings.TrimSpace(string(output)), "\n")

		for i := range mediaFilePath {
			// bliss `library` output has each line surrounded by quotes
			// remove the first and last character
			withoutQuotes := mediaFilePath[i][1 : len(mediaFilePath[i])-1]
			mediaFilePath[i] = strings.Replace(withoutQuotes, conf.Server.Bliss.PrependPrefix, conf.Server.Bliss.RemovePrefix, 1)
		}

		files, err := mfRepo.FindByAbsolutePaths(mediaFilePath)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		response, err := json.Marshal(files)
		if err != nil {
			log.Error(ctx, "Error marshalling json", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")

		if _, err := w.Write(response); err != nil {
			log.Error(ctx, "Error sending response to client", err)
		}
	}
}

func (n *Router) addInspectMixRoute(r chi.Router) {
	if conf.Server.Bliss.Enabled {
		r.Get("/instantMix", n.instantMix())
	}
}
