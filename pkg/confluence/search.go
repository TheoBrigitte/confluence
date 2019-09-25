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

	g.Go(func() error {
		for i := 0; i < 2; i++ {
			movies = append(movies, <-movieChan...)
		}
		return nil
	})

	// yify
	{
		g.Go(func() error {
			ms := []movie.MovieTorrent{}
			defer func() { movieChan <- ms }()
			y := yify.New()
			m, err := y.SearchMoviesWithBestTorrent(query)
			if err != nil {
				return fmt.Errorf("yify search failed: %v", err.Error())
			}
			if len(m) > 0 {
				ms = m
			}

			return nil
		})
	}

	// cpasbien
	{
		g.Go(func() error {
			ms := []movie.MovieTorrent{}
			defer func() { movieChan <- ms }()
			c, err := cpasbien.New(cpasbien.Config{})
			if err != nil {
				return fmt.Errorf("cpasbien init failed: %v", err.Error())
			}
			m, err := c.Search(query)
			if err != nil {
				return fmt.Errorf("cpasbien search failed: %v", err.Error())
			}
			if len(m) > 0 {
				ms = m
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
