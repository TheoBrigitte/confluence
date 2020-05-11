package tmdb

import (
	"github.com/TheoBrigitte/confluence/pkg/movie"
)

// api response
type popularResponse struct {
	Page         int           `json:"page"`
	Results      []movieDetail `json:"results"`
	TotalPages   int           `json:"total_pages"`
	TotalResults int           `json:"total_results"`
}

type movieDetail struct {
	ID               int     `json:"id"`
	Title            string  `json:"title"`
	Description      string  `json:"overview"`
	OriginalTitle    string  `json:"original_title"`
	OriginalLanguage string  `json:"original_language"`
	ReleaseDate      string  `json:"release_date"`
	PosterPath       string  `json:"poster_path"`
	BackdropPath     string  `json:"backdrop_path"`
	GenreIDs         []int   `json:"genre_ids"`
	IMDB             string  `json:"imdb_id"`
	Popularity       float32 `json:"popularity"`
	VoteCount        int     `json:"vote_count"`
	VoteAverage      float32 `json:"vote_average"`
	Video            bool    `json:"video"`
	Adult            bool    `json:"adult"`
}

func (m movieDetail) ToMovieTorrent() movie.MovieTorrent {
	mt := movie.MovieTorrent{
		MovieBase: movie.MovieBase{
			ID:    m.ID,
			Title: m.Title,
			Image: imageHost + "/w342" + m.PosterPath,
			ExternalID: movie.ExternalID{
				IMDB: m.IMDB,
			},
		},
	}

	return mt
}
