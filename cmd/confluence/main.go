package main

import (
	"log"
	"net"
	"net/http"

	"github.com/TheoBrigitte/confluence/pkg/confluence"
	_ "github.com/anacrolix/envpprof"
	"github.com/anacrolix/tagflag"
	"github.com/anacrolix/torrent"
)

func main() {
	// logs configuration: set short filename format
	log.SetFlags(log.Flags() | log.Lshortfile)

	// command line arguments
	tagflag.Parse(&flags)

	// torrent client
	storage, onTorrentGraceExtra := getStorage()
	cl, err := newTorrentClient(storage)
	if err != nil {
		log.Fatalf("error creating torrent client: %s", err)
	}
	defer cl.Close()

	// debug endpoint
	http.HandleFunc("/debug/dht", func(w http.ResponseWriter, r *http.Request) {
		for _, ds := range cl.DhtServers() {
			ds.WriteStatus(w)
		}
	})

	// subtitles configuration: set opensubtitles credentials
	confluence.SetOSCredentials(flags.OSUser, flags.OSPassword, flags.OSUserAgent)

	// HTTP server
	l, err := net.Listen("tcp", flags.Addr)
	if err != nil {
		log.Fatal(err)
	}
	defer l.Close()

	var h http.Handler = &confluence.Handler{
		cl,
		flags.TorrentGrace,
		func(t *torrent.Torrent) {
			ih := t.InfoHash()
			t.Drop()
			onTorrentGraceExtra(ih)
		},
	}

	if flags.DebugOnMain {
		h = func() http.Handler {
			mux := http.NewServeMux()
			mux.Handle("/debug/", http.DefaultServeMux)
			mux.Handle("/", h)
			return mux
		}()
	}

	log.Printf("start http server at %s", l.Addr())
	err = http.Serve(l, h)
	if err != nil {
		log.Fatal(err)
	}
}
