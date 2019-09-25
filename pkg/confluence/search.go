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
	q := r.URL.Query()
	query := q.Get("query")
	log.Printf("query: %#q\n", query)

	var movieChan = make(chan []movie.MovieTorrent)
	var g errgroup.Group
	var movies []movie.MovieTorrent

	go func() {
		for i := 0; i < 2; i++ {
			movies = append(movies, <-movieChan...)
		}
	}()

	// yify
	{
		g.Go(func() error {
			movies := []movie.MovieTorrent{}
			defer func() { movieChan <- movies }()
			y := yify.New()
			m, err := y.SearchMoviesWithBestTorrent(query)
			if err != nil {
				return fmt.Errorf("yify search failed: %v", err.Error())
			}
			if len(m) > 0 {
				movies = m
			}

			return nil
		})
	}

	// cpasbien
	{
		g.Go(func() error {
			var movies []movie.MovieTorrent
			defer func() { movieChan <- movies }()
			c, err := cpasbien.New(cpasbien.Config{})
			if err != nil {
				return fmt.Errorf("cpasbien init failed: %v", err.Error())
			}
			m, err := c.Search(query)
			if err != nil {
				return fmt.Errorf("cpasbien search failed: %v", err.Error())
			}
			if len(m) > 0 {
				movies = m
			}

			return nil
		})
	}

	if err := g.Wait(); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(movies)
}
