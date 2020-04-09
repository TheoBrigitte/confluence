package confluence

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/TheoBrigitte/confluence/pkg/cpasbien"
	"github.com/TheoBrigitte/confluence/pkg/movie"
	"github.com/TheoBrigitte/confluence/pkg/yify"
)

func searchHandler(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	query := q.Get("query")
	log.Printf("query: %#q\n", query)

	var movies []movie.MovieTorrent
	var movieChan = make(chan movie.MovieTorrent)
	var errors = make(chan error)

	// yify
	{
		go func() {
			y := yify.New()
			ms, err := y.SearchMoviesWithBestTorrent(query)
			if err != nil {
				errors <- fmt.Errorf("yify search failed: %v", err.Error())
				return
			}
			for _, m := range ms {
				movieChan <- m
			}

			return
		}()
	}

	// cpasbien
	{
		go func() {
			c, err := cpasbien.New(cpasbien.Config{})
			if err != nil {
				errors <- fmt.Errorf("cpasbien init failed: %v", err.Error())
				return
			}

			ms, err := c.Search(query)
			if err != nil {
				errors <- fmt.Errorf("cpasbien search failed: %v", err.Error())
				return
			}
			for _, m := range ms {
				movieChan <- m
			}

			return
		}()
	}

	for i := 0; i < 2; i++ {
		select {
		case movie := <-movieChan:
			movies = append(movies, movie)
		case err := <-errors:
			log.Printf("error: %v\n", err)
		}
	}

	if len(movies) == 0 {
		http.Error(w, fmt.Sprintf("%#q movie not found", query), http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(movies)
}
