package confluence

import (
	"context"
	"log"
	"net/http"
	"net/url"

	"github.com/anacrolix/torrent"
	"github.com/anacrolix/torrent/metainfo"
	"golang.org/x/net/websocket"
)

type Bar struct {
	Status BarStatus
	Width  float32
}

type BarStatus string

const (
	BarComplete = BarStatus("complete")
	BarDownload = BarStatus("download")
	BarStale    = BarStatus("stale")
	BarFail     = BarStatus("fail")
)

func torrentBiggestFile(t *torrent.Torrent) *torrent.File {
	var file torrent.File

	for _, f := range t.Files() {
		if f.Length() > file.Length() {
			file = *f
		}
	}
	return &file
}

func barWidth(width, pieces int) float32 {
	return float32(width) / float32(pieces)
}

func barStatus(p torrent.FilePieceState) BarStatus {
	if !p.Ok {
		return BarFail
	}

	if p.Complete {
		return BarComplete
	}

	if p.Partial {
		return BarDownload
	}

	return BarStale
}

func fileAggregateState(f *torrent.File) []Bar {
	states := f.State()

	var bars []Bar
	if len(states) < 1 {
		return bars
	}

	width := barWidth(100, len(states))
	bar := Bar{
		Status: barStatus(states[0]),
		Width:  width,
	}

	for _, state := range states[1:] {
		status := barStatus(state)
		if bar.Status == status {
			bar.Width += width
		} else {
			bars = append(bars, bar)
			bar = Bar{
				Status: status,
				Width:  width,
			}
		}
	}
	bars = append(bars, bar)

	return bars
}

func magnetlinkFromQueryOrServeError(w http.ResponseWriter, q url.Values) *metainfo.Magnet {
	magnet, err := metainfo.ParseMagnetURI(q.Get(magnetlinkQueryKey))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return nil
	}

	return &magnet
}

func withMagnetlinkContext(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		magnet := magnetlinkFromQueryOrServeError(w, r.URL.Query())

		t := getTorrentHandle(r, magnet.InfoHash)
		h.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), torrentContextKey, t)))
	})
}

func progressHandler(w http.ResponseWriter, r *http.Request) {
	t := torrentForRequest(r)
	select {
	case <-t.GotInfo():
	case <-r.Context().Done():
		return
	}

	file := torrentBiggestFile(t)
	if file == nil {
		http.Error(w, "file not found", http.StatusNotFound)
		return
	}

	s := t.SubscribePieceStateChanges()
	defer s.Close()
	websocket.Server{
		Handler: func(c *websocket.Conn) {
			defer c.Close()
			log.Printf("handling websocket connection %q %#q", c.Request().Method, c.Request().RequestURI)
			readClosed := make(chan struct{})
			go func() {
				defer close(readClosed)
				c.Read(nil)
			}()
			bars := fileAggregateState(file)
			if err := websocket.JSON.Send(c, bars); err != nil {
				if r.Context().Err() == nil {
					log.Printf("error writing json to websocket: %s", err)
				}
				return
			}
			for {
				select {
				case <-readClosed:
					log.Printf("read closed connection %q %#q", c.Request().Method, c.Request().RequestURI)
					eventHandlerWebsocketReadClosed.Add(1)
					return
				case <-r.Context().Done():
					log.Printf("context done connection %q %#q", c.Request().Method, c.Request().RequestURI)
					eventHandlerContextDone.Add(1)
					return
				case <-s.Values:
					log.Printf("values connection %q %#q", c.Request().Method, c.Request().RequestURI)
					//i := _i.(torrent.PieceStateChange).Index
					bars := fileAggregateState(file)
					if err := websocket.JSON.Send(c, bars); err != nil {
						if r.Context().Err() == nil {
							log.Printf("error writing json to websocket: %s", err)
						}
						return
					}
				}
			}
			log.Printf("end connection %q %#q", c.Request().Method, c.Request().RequestURI)
		},
	}.ServeHTTP(w, r)
}
