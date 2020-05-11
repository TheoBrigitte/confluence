package yify

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/TheoBrigitte/confluence/pkg/movie/provider"
)

type client struct {
	client *http.Client
}

func New() provider.Interface {
	return &client{
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
	log.Printf("response %s %#q body=%d %v", res.Status, u.String(), res.ContentLength, res.Header)

	return res, err
}

func decodeJSON(res *http.Response, dst interface{}) error {
	return json.NewDecoder(res.Body).Decode(dst)
}
