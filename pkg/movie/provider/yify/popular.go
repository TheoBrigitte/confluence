package yify

import (
	"net/http"

	"github.com/TheoBrigitte/confluence/pkg/movie"
)

func (c *client) Popular(limit int) ([]movie.MovieTorrent, error) {
	res, err := c.popularMovies(limit)
	if err != nil {
		return nil, err
	}

	movies := []movie.MovieTorrent{}
	for _, m := range res.Data.Movies {
		if m.HasTorrent() {
			movies = append(movies, m.ToMovieTorrent())
		}
	}

	return movies, nil
}

func (c client) popularMovies(limit int) (*searchResponse, error) {
	u := searchURL
	q := u.Query()
	//q.Set(searchLimitKey, limit)
	q.Set(searchSortKey, "peers")
	u.RawQuery = q.Encode()
	res, err := c.do(http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}

	var response searchResponse
	err = decodeJSON(res, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}
