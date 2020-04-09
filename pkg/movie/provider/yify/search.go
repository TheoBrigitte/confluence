package yify

import (
	"net/http"

	"github.com/TheoBrigitte/confluence/pkg/movie"
)

// api response
type searchResponse struct {
	Data          searchResponseData `json:"data"`
	Status        string             `json:"status"`
	StatusMessage string             `json:"status_message"`
}

type searchResponseData struct {
	Limit      int           `json:"limit"`
	MovieCount int           `json:"movie_count"`
	Movies     []movie.Movie `json:"movies"`
	PageNumber int           `json:"page_number"`
}

func (c client) SearchMovies(query string) (*searchResponse, error) {
	u := searchURL
	q := u.Query()
	q.Set(searchQueryKey, query)
	u.RawQuery = q.Encode()
	res, err := c.do(http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}

	var response searchResponse
	err = toJson(res, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}
