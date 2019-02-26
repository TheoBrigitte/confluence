package yify

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/url"
)

var (
	baseURL, _ = url.Parse("https://yts.am/api/v2/")

	//https://yts.am/api/v2/list_movies.json?query_term=test
	searchEndpoint, _ = url.Parse("list_movies.json")
	searchQueryKey    = "query_term"
)

type client struct {
	client *http.Client
}

func New() client {
	return client{
		client: &http.Client{},
	}
}

func (c client) do(method, url string, body io.Reader) (*http.Response, error) {
	u, err := baseURL.Parse(url)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(method, u.String(), body)
	if err != nil {
		return nil, err
	}

	log.Printf("request %s %#q body=%t", method, u.String(), body != nil)
	res, err := c.client.Do(req)
	log.Printf("response %s body=%d %v", res.Status, res.ContentLength, res.Header)

	return res, err
}

func toJson(res *http.Response, dst interface{}) error {
	return json.NewDecoder(res.Body).Decode(dst)
}
