package cpasbien

import (
	"net/http"
	"net/url"

	"github.com/TheoBrigitte/confluence/pkg/movie"
)

var (
	baseURL, _ = url.Parse("https://www2.cpasbiens.to")
)

type Interface interface {
	Search(string) ([]movie.MovieTorrent, error)
}

type client struct {
	http *http.Client
}

func New() Interface {
	return &client{
		http: &http.Client{},
	}
}
