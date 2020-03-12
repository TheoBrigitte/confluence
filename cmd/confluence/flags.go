package main

import (
	"net"
	"time"

	"github.com/anacrolix/tagflag"
)

var flags = struct {
	Addr               string        `help:"HTTP listen address"`
	PublicIP4          net.IP        `help:"Public IPv4 address"` // TODO: Rename
	PublicIP6          net.IP        `help:"Public IPv6 address"`
	CacheCapacity      tagflag.Bytes `help:"Data cache capacity"`
	TorrentGrace       time.Duration `help:"How long to wait to drop a torrent after its last request"`
	FileDir            string        `help:"File-based storage directory, overrides piece storage"`
	Seed               bool          `help:"Seed data"`
	UPnPPortForwarding bool          `help:"Port forward via UPnP"`
	// You'd want this if access to the main HTTP service is trusted, such as
	// used over localhost by other known services.
	DebugOnMain bool `help:"Expose default serve mux /debug/ endpoints over http"`
	Dht         bool

	OSUser      string `help:"OpenSubtitles User login"`
	OSPassword  string `help:"OpenSubtitles Password login"`
	OSUserAgent string `help:"OpenSubtitles User Agent"`
}{
	Addr:          "localhost:8080",
	CacheCapacity: 10 << 30,
	TorrentGrace:  time.Minute,
	Dht:           true,
}
