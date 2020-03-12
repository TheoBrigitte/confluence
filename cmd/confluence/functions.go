package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/anacrolix/missinggo/filecache"
	"github.com/anacrolix/missinggo/x"
	"github.com/anacrolix/torrent"
	"github.com/anacrolix/torrent/iplist"
	"github.com/anacrolix/torrent/storage"
)

// newTorrentClient create a torrent client.
func newTorrentClient(storage storage.ClientImpl) (ret *torrent.Client, err error) {
	blocklist, err := iplist.MMapPackedFile("packed-blocklist")
	if err != nil {
		log.Print(err)
	} else {
		defer func() {
			if err != nil {
				blocklist.Close()
			} else {
				go func() {
					<-ret.Closed()
					blocklist.Close()
				}()
			}
		}()
	}
	cfg := torrent.NewDefaultClientConfig()
	cfg.IPBlocklist = blocklist
	cfg.DefaultStorage = storage
	cfg.PublicIp4 = flags.PublicIP4
	cfg.PublicIp6 = flags.PublicIP6
	cfg.Seed = flags.Seed
	cfg.NoDefaultPortForwarding = !flags.UPnPPortForwarding
	cfg.NoDHT = !flags.Dht
	cfg.SetListenAddr(":50007")
	http.HandleFunc("/debug/conntrack", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		cfg.ConnTracker.PrintStatus(w)
	})

	// cfg.DisableAcceptRateLimiting = true
	return torrent.NewClient(cfg)
}

// getStorage create a storage.
func getStorage() (_ storage.ClientImpl, onTorrentGrace func(torrent.InfoHash)) {
	if flags.FileDir != "" {
		return storage.NewFileByInfoHash(flags.FileDir), func(ih torrent.InfoHash) {
			os.RemoveAll(filepath.Join(flags.FileDir, ih.HexString()))
		}
	}
	fc, err := filecache.NewCache("filecache")
	x.Pie(err)

	// Register filecache debug endpoints on the default muxer.
	http.HandleFunc("/debug/filecache/status", func(w http.ResponseWriter, r *http.Request) {
		info := fc.Info()
		fmt.Fprintf(w, "Capacity: %d\n", info.Capacity)
		fmt.Fprintf(w, "Current Size: %d\n", info.Filled)
		fmt.Fprintf(w, "Item Count: %d\n", info.NumItems)
	})
	http.HandleFunc("/debug/filecache/lru", func(w http.ResponseWriter, r *http.Request) {
		fc.WalkItems(func(item filecache.ItemInfo) {
			fmt.Fprintf(w, "%s\t%d\t%s\n", item.Accessed, item.Size, item.Path)
		})
	})

	fc.SetCapacity(flags.CacheCapacity.Int64())
	storageProvider := fc.AsResourceProvider()
	return storage.NewResourcePieces(storageProvider), func(ih torrent.InfoHash) {}
}
