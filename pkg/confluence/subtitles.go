package confluence

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"sort"
	"strconv"
	"time"

	"github.com/TheoBrigitte/confluence/pkg/subtitles"
	"github.com/anacrolix/torrent"
	astisub "github.com/asticode/go-astisub"
	"github.com/oz/osdb"
)

var (
	osUser      string
	osPassword  string
	osUserAgent string
)

func SetOSCredentials(user, password, useragent string) {
	osUser = user
	osPassword = password
	osUserAgent = useragent
}

type subtitleResult struct {
	Downloads string `json:"downloads"`
	ID        string `json:"id"`
	Language  string `json:"language"`
	Name      string `json:"name"`
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
		UserAgent: osUserAgent,
	}

	log.Printf("\nhead\n%v\n%v\n\ntail\n%v\n%v\n\n", head[:10], head[len(head)-10:], tail[:10], tail[len(tail)-10:])
	searcher, err := subtitles.New(config)
	if err != nil {
		http.Error(w, fmt.Sprintf("searcher failed: %v", err.Error()), http.StatusBadRequest)
		return
	}

	languages, ok := q["lang"]
	if !ok {
		languages = []string{"eng"}
	}
	subtitles, err := searcher.SearchSubtitles(head, tail, size, languages)
	if err != nil {
		http.Error(w, fmt.Sprintf("subtitles get failed: %v", err.Error()), http.StatusBadRequest)
		return
	}

	sort.Sort(osdb.ByDownloads(subtitles))

	results := []subtitleResult{}
	for _, s := range subtitles {
		r := subtitleResult{
			Downloads: s.SubDownloadsCnt,
			ID:        s.IDSubtitleFile,
			Language:  s.LanguageName,
			Name:      s.MovieFileName,
		}
		results = append(results, r)
	}

	json.NewEncoder(w).Encode(results)
}

func subtitleHandler(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	idStr := q.Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, fmt.Sprintf("wrong id: id=%s %v", idStr, err.Error()), http.StatusBadRequest)
		return
	}

	config := subtitles.Config{
		User:      osUser,
		Password:  osPassword,
		UserAgent: osUserAgent,
	}
	searcher, err := subtitles.New(config)
	if err != nil {
		http.Error(w, fmt.Sprintf("searcher failed: %v", err.Error()), http.StatusBadRequest)
		return
	}

	sr, err := searcher.GetSubtitleReader(id)
	if err != nil {
		http.Error(w, fmt.Sprintf("subtitle get failed: %v", err.Error()), http.StatusBadRequest)
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

	w.Header().Set("Content-Type", "text/vtt")
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
