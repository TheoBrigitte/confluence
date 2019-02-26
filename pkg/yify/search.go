package yify

import (
	"net/http"
)

// api response
type searchResponse struct {
	Data          searchResponseData `json:"data"`
	Status        string             `json:"status"`
	StatusMessage string             `json:"status_message"`
}

type searchResponseData struct {
	Limit      int     `json:"limit"`
	MovieCount int     `json:"movie_count"`
	Movies     []movie `json:"movies"`
	PageNumber int     `json:"page_number"`
}

type movieBase struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	Year  int    `json:"year"`
}

type movie struct {
	movieBase
	Torrents []torrent `json:"torrents"`
}

type torrent struct {
	Hash      string `json:"hash"`
	Quality   string `json:"quality"`
	Seeds     int    `json:"seeds"`
	Size      string `json:"size"`
	SizeBytes int    `json:"size_bytes"`
	Type      string `json:"type"`
	URL       string `json:"url"`
}

func (c client) SearchMovies(query string) (*searchResponse, error) {
	u := searchEndpoint
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
