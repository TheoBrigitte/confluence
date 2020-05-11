package tmdb

import (
	"net/http"

	"github.com/TheoBrigitte/confluence/pkg/movie/provider"
	"github.com/TheoBrigitte/confluence/pkg/movie/provider/yify"
)

type client struct {
	http   *http.Client
	finder provider.Finder
}

func New() provider.Popular {
	return &client{
		http:   &http.Client{},
		finder: yify.New(),
	}
}
