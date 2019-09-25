package cpasbien

import (
	"net/http"
	"net/url"

	"github.com/TheoBrigitte/confluence/pkg/movie"
)

const (
	defaultBaseURL = "https://www2.cpasbiens.to"
)

type Interface interface {
	Search(string) ([]movie.MovieTorrent, error)
}

type client struct {
	http    *http.Client
	baseURL *url.URL
}

type Config struct {
	URL string
}

func New(cfg Config) (Interface, error) {
	var err error
	var u *url.URL
	{
		var baseURL string
		if cfg.URL == "" {
			baseURL = defaultBaseURL
		} else {
			baseURL = cfg.URL
		}
		u, err = url.Parse(baseURL)
		if err != nil {
			return nil, err
		}
	}

	c := client{
		http:    &http.Client{},
		baseURL: u,
	}

	return &c, nil
}
