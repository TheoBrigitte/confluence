package yify

import "github.com/TheoBrigitte/confluence/pkg/movie"

func (c client) SearchMoviesWithBestTorrent(query string) ([]movie.MovieTorrent, error) {
	res, err := c.SearchMovies(query)
	if err != nil {
		return nil, err
	}

	movies := []movie.MovieTorrent{}
	for _, m := range res.Data.Movies {
		movies = append(movies, m.BestTorrent())
	}

	return movies, nil
}
