package confluence

import (
	"context"
	"net/http"
	"time"

	"github.com/anacrolix/torrent"
)

type Handler struct {
	TC             *torrent.Client
	TorrentGrace   time.Duration
	OnTorrentGrace func(t *torrent.Torrent)
}

func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	r = r.WithContext(context.WithValue(r.Context(), handlerContextKey, h))
	w.Header().Set("Access-Control-Allow-Origin", "*")
	mux.ServeHTTP(w, r)
}
