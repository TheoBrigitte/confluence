package provider

import (
	"github.com/TheoBrigitte/confluence/pkg/movie"
)

type Searcher interface {
	Search(string) ([]movie.MovieTorrent, error)
}

type Popular interface {
	Popular(int) ([]movie.MovieTorrent, error)
}

type Finder interface {
	Find(string) (*movie.MovieTorrent, error)
}

type PopularSearchFinder interface {
	Searcher
	Popular
	Finder
}
