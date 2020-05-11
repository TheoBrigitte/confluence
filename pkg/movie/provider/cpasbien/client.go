package cpasbien

import (
	"net/http"

	"github.com/TheoBrigitte/confluence/pkg/movie/provider"
)

type client struct {
	http *http.Client
}

func New() provider.Searcher {
	return &client{
		http: &http.Client{},
	}
}
