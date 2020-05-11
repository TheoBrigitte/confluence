package yify

import (
	"fmt"

	"github.com/TheoBrigitte/confluence/pkg/movie"
)

func (c *client) Find(query string) (*movie.MovieTorrent, error) {
	res, err := c.searchMovies(map[string]string{
		SearchQueryKey: query,
		SearchLimitKey: "1",
	})
	if err != nil {
		return nil, err
	}

	if len(res.Data.Movies) <= 0 {
		return nil, fmt.Errorf("not found: %q", query)
	}

	m := res.Data.Movies[0].ToMovieTorrent()
	return &m, nil
}
