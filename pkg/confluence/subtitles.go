package confluence

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/TheoBrigitte/confluence/pkg/subtitles"
	"github.com/anacrolix/torrent"
	astisub "github.com/asticode/go-astisub"
)

var (
	osUser     string
	osPassword string
)

func SetOSCredentials(user, password string) {
	osUser = user
	osPassword = password
}

func subtitlesHandler(w http.ResponseWriter, r *http.Request) {
	t := torrentForRequest(r)
	select {
	case <-t.GotInfo():
	case <-r.Context().Done():
		return
	}

	//file := torrentBiggestFile(t)
	q := r.URL.Query()
	file := torrentFileByPath(t, q.Get("path"))
	if file == nil {
		http.Error(w, "file not found", http.StatusNotFound)
		return
	}

	fr := file.NewReader()
	defer fr.Close()

	var head subtitles.Chunk
	err := torrentReadChunck(fr, 0, head[:])
	if err != nil {
		http.Error(w, fmt.Sprintf("head chunk failed: %v", err.Error()), http.StatusBadRequest)
		return
	}

	var tail subtitles.Chunk
	size := file.Length()
	offset := size - subtitles.ChunkSize
	err = torrentReadChunck(fr, offset, tail[:])
	if err != nil {
		http.Error(w, fmt.Sprintf("tail chunk failed: size=%d offset=%d err=%v", size, offset, err.Error()), http.StatusBadRequest)
		return
	}

	config := subtitles.Config{
		User:      osUser,
		Password:  osPassword,
		Start:     head,
		End:       tail,
		Size:      size,
		Languages: []string{"eng"},
	}

	log.Printf("\nhead\n%v\n%v\n\ntail\n%v\n%v\n\n", head[:10], head[len(head)-10:], tail[:10], tail[len(tail)-10:])
	searcher, err := subtitles.New(config)
	if err != nil {
		http.Error(w, fmt.Sprintf("searcher failed: %v", err.Error()), http.StatusBadRequest)
		return
	}

	sr, err := searcher.GetSubtitleReader()
	if err != nil {
		http.Error(w, fmt.Sprintf("subtitles get failed: %v", err.Error()), http.StatusBadRequest)
		return
	}
	defer sr.Close()

	as, err := astisub.ReadFromSRT(sr)
	if err != nil {
		http.Error(w, fmt.Sprintf("read srt failed: %v", err.Error()), http.StatusBadRequest)
		return
	}

	var buf = &bytes.Buffer{}
	as.WriteToWebVTT(buf)

	rs := bytes.NewReader(buf.Bytes())

	http.ServeContent(w, r, "", time.Time{}, rs)
}

func torrentReadChunck(tr torrent.Reader, offset int64, buf []byte) error {
	tr.SetReadahead(subtitles.ChunkSize)

	_, err := tr.Seek(offset, io.SeekStart)
	if err != nil {
		return fmt.Errorf("seek: %v", err)
	}

	_, err = tr.Read(buf)
	if err != nil && err != io.EOF {
		return fmt.Errorf("read: %v", err)
	}

	return nil
}
