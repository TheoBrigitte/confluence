package yify

import (
	"net/http"

	"github.com/TheoBrigitte/confluence/pkg/movie"
	"github.com/TheoBrigitte/confluence/pkg/util"
)

func (c *client) Search(query string) ([]movie.MovieTorrent, error) {
	res, err := c.searchMovies(map[string]string{
		SearchQueryKey: query,
	})
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

func (c *client) searchMovies(query map[string]string) (*searchResponse, error) {
	u := searchURL
	if len(query) > 0 {
		q := u.Query()
		for k, v := range query {
			q.Set(k, v)
		}
		u.RawQuery = q.Encode()
	}
	res, err := c.do(http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}

	var response searchResponse
	err = util.DecodeJSON(res, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}
