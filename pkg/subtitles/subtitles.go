package subtitles

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"strings"

	"github.com/oz/osdb"
)

const (
	ChunkSize = 65536
)

type Chunk [ChunkSize]byte

type Config struct {
	Languages []string
	Size      int64

	Start Chunk
	End   Chunk

	User     string
	Password string
}

type Searcher interface {
	Search()
}

type searcher struct {
	client *osdb.Client

	languages []string
	size      int64

	start Chunk
	end   Chunk
}

func New(config Config) (*searcher, error) {
	c, err := osdb.NewClient()
	if err != nil {
		return nil, err
	}

	err = c.LogIn(config.User, config.Password, "eng")
	if err != nil {
		return nil, err
	}

	s := &searcher{
		client: c,

		start:     config.Start,
		end:       config.End,
		size:      config.Size,
		languages: config.Languages,
	}

	return s, nil
}

func (s searcher) GetSubtitleReader() (io.ReadCloser, error) {
	sub, err := s.GetSubtitle()
	if err != nil {
		return nil, err
	}

	files, err := s.client.DownloadSubtitles(osdb.Subtitles{*sub})
	if err != nil {
		return nil, err
	}
	if len(files) == 0 {
		return nil, fmt.Errorf("no_subtitle_file")
	}

	return files[0].Reader()
}

func (s searcher) GetSubtitle() (*osdb.Subtitle, error) {
	sub, err := s.SearchSubtitles()
	if err != nil {
		return nil, err
	}

	subtitle := sub.Best()
	if subtitle == nil {
		return nil, fmt.Errorf("no_subtitles")
	}

	return subtitle, nil
}

func (s searcher) SearchSubtitles() (osdb.Subtitles, error) {
	hash, err := s.computeHash()
	if err != nil {
		return nil, err
	}

	params := s.params(hash)
	log.Printf("%+v", params)

	return s.client.SearchSubtitles(params)
}

func (s *searcher) computeHash() (hash uint64, err error) {
	buf := append(s.start[:], s.end[:]...)

	// Convert to uint64, and sum.
	var nums [(ChunkSize * 2) / 8]uint64
	reader := bytes.NewReader(buf)
	err = binary.Read(reader, binary.LittleEndian, &nums)
	if err != nil {
		return
	}

	for _, num := range nums {
		hash += num
	}

	return hash + uint64(s.size), nil
}

func (s searcher) params(hash uint64) *[]interface{} {
	params := []interface{}{
		s.client.Token,
		[]struct {
			Hash  string `xmlrpc:"moviehash"`
			Size  int64  `xmlrpc:"moviebytesize"`
			Langs string `xmlrpc:"sublanguageid"`
		}{{
			hashString(hash),
			s.size,
			strings.Join(s.languages, ","),
		}},
	}
	return &params
}

// Create a string representation of hash
func hashString(hash uint64) string {
	return fmt.Sprintf("%016x", hash)
}
