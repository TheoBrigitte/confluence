package provider

import (
	"github.com/TheoBrigitte/confluence/pkg/movie"
)

type Interface interface {
	Search(string) ([]movie.MovieTorrent, error)
}
