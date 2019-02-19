package confluence

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/anacrolix/torrent/metainfo"
)

const (
	queryQueryKey      = "query"
	magnetlinkQueryKey = "magnet"
)

func hashOrMagnet(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		query := q.Get(queryQueryKey)

		var hash metainfo.Hash
		ihErr := hash.FromHexString(query)
		if ihErr != nil {

			magnet, mErr := metainfo.ParseMagnetURI(query)
			if mErr != nil {
				http.Error(w, fmt.Sprintf("ih: %v, magnet: %v", ihErr.Error(), mErr.Error()), http.StatusBadRequest)
				return
			}
			hash = magnet.InfoHash
		}

		t := getTorrentHandle(r, hash)
		h.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), torrentContextKey, t)))
	})
}

type metadata struct {
	File   string `json:"file"`
	Hash   string `json:"hash"`
	Length int64  `json:"length"`
	Name   string `json:"name"`
	Pieces int    `json:"pieces"`
}

func metadataHandler(w http.ResponseWriter, r *http.Request) {
	t := torrentForRequest(r)
	// w.WriteHeader(http.StatusProcessing)
	select {
	case <-t.GotInfo():
	case <-r.Context().Done():
		return
	}

	// w.WriteHeader(http.StatusOK)
	i := t.Info()
	m := metadata{
		File:   torrentBiggestFile(t).DisplayPath(),
		Hash:   t.InfoHash().String(),
		Length: i.Length,
		Name:   i.Name,
		Pieces: len(i.Pieces),
	}
	json.NewEncoder(w).Encode(m)
}
