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
	User     string
	Password string
}

type Searcher interface {
	Search()
}

type searcher struct {
	client *osdb.Client
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
	}

	return s, nil
}

func (s searcher) GetSubtitleReader(id int) (io.ReadCloser, error) {
	files, err := s.client.DownloadSubtitlesByIds([]int{id})
	if err != nil {
		return nil, err
	}
	if len(files) == 0 {
		return nil, fmt.Errorf("no_subtitle_file")
	}

	return files[0].Reader()
}

func (s searcher) GetSubtitle(start, end Chunk, size int64, languages []string) (*osdb.Subtitle, error) {
	sub, err := s.SearchSubtitles(start, end, size, languages)
	if err != nil {
		return nil, err
	}

	subtitle := sub.Best()
	if subtitle == nil {
		return nil, fmt.Errorf("no_subtitles")
	}

	return subtitle, nil
}

func (s searcher) SearchSubtitles(start, end Chunk, size int64, languages []string) (osdb.Subtitles, error) {
	hash, err := s.computeHash(start, end, size)
	if err != nil {
		return nil, err
	}

	params := s.params(hash, size, languages)
	log.Printf("%+v", params)

	return s.client.SearchSubtitles(params)
}

func (s *searcher) computeHash(start, end Chunk, size int64) (hash uint64, err error) {
	buf := append(start[:], end[:]...)

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

	return hash + uint64(size), nil
}

func (s searcher) params(hash uint64, size int64, languages []string) *[]interface{} {
	params := []interface{}{
		s.client.Token,
		[]struct {
			Hash  string `xmlrpc:"moviehash"`
			Size  int64  `xmlrpc:"moviebytesize"`
			Langs string `xmlrpc:"sublanguageid"`
		}{{
			hashString(hash),
			size,
			strings.Join(languages, ","),
		}},
	}
	return &params
}

// Create a string representation of hash
func hashString(hash uint64) string {
	return fmt.Sprintf("%016x", hash)
}
