package confluence

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/TheoBrigitte/confluence/pkg/cpasbien"
	"github.com/TheoBrigitte/confluence/pkg/movie"
	"github.com/TheoBrigitte/confluence/pkg/yify"
	"golang.org/x/sync/errgroup"
)

func searchHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("searchHandler")
	q := r.URL.Query()
	query := q.Get("query")

	var m = make(chan []movie.MovieTorrent)
	var g errgroup.Group
	var movies []movie.MovieTorrent

	go func() {
		for i := 0; i < 2; i++ {
			movies = append(movies, <-m...)
		}
	}()

	// yify
	{
		g.Go(func() error {
			var movies []movie.MovieTorrent
			defer func() { m <- movies }()
			y := yify.New()
			ms, err := y.SearchMoviesWithBestTorrent(query)
			if err != nil {
				return fmt.Errorf("yify search failed: %v", err.Error())
			}
			movies = ms

			return nil
		})
	}

	// cpasbien
	{
		g.Go(func() error {
			var movies []movie.MovieTorrent
			defer func() { m <- movies }()
			c, err := cpasbien.New(cpasbien.Config{})
			if err != nil {
				return fmt.Errorf("cpasbien init failed: %v", err.Error())
			}
			ms, err := c.Search(query)
			if err != nil {
				return fmt.Errorf("cpasbien search failed: %v", err.Error())
			}
			movies = ms

			return nil
		})
	}

	if err := g.Wait(); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(movies)
}
