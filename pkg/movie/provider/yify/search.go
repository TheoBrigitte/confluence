package yify

import (
	"net/http"
)

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
