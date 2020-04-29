package yify

import (
	"github.com/TheoBrigitte/confluence/pkg/movie"
)

// api response
type searchResponse struct {
	Data          searchResponseData `json:"data"`
	Status        string             `json:"status"`
	StatusMessage string             `json:"status_message"`
}

type searchResponseData struct {
	Limit      int      `json:"limit"`
	MovieCount int      `json:"movie_count"`
	Movies     []yMovie `json:"movies"`
	PageNumber int      `json:"page_number"`
}

type yMovie struct {
	ID       int        `json:"id"`
	Title    string     `json:"title"`
	Year     int        `json:"year"`
	IMDB     string     `json:"imdb_code"`
	Torrents []yTorrent `json:"torrents"`
}

type yTorrent struct {
	Hash      string `json:"hash"`
	Quality   string `json:"quality"`
	Seeds     int    `json:"seeds"`
	Size      string `json:"size"`
	SizeBytes int    `json:"size_bytes"`
	Type      string `json:"type"`
	URL       string `json:"url"`
}

func (m yMovie) ToMovieTorrent() movie.MovieTorrent {
	// Find best torrent
	var current yTorrent
	for _, t := range m.Torrents {
		if t.Seeds > current.Seeds {
			current = t
		}
	}

	mt := movie.MovieTorrent{
		MovieBase: movie.MovieBase{
			ID:    m.ID,
			Title: m.Title,
			Year:  m.Year,
			ExternalID: movie.ExternalID{
				IMDB: m.IMDB,
			},
		},
		Torrent: movie.Torrent{
			Hash:      current.Hash,
			Quality:   current.Quality,
			Seeds:     current.Seeds,
			Size:      current.Size,
			SizeBytes: current.SizeBytes,
			Type:      current.Type,
			URL:       current.URL,
		},
	}

	return mt
}
