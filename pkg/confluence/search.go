package confluence

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/TheoBrigitte/confluence/pkg/movie"
	"github.com/TheoBrigitte/confluence/pkg/movie/provider/cpasbien"
	"github.com/TheoBrigitte/confluence/pkg/movie/provider/yify"
)

var (
	timeout = time.After(5 * time.Second)
)

func searchHandler(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	query := q.Get("query")
	log.Printf("query: %#q\n", query)

	var moviesChan = make(chan []movie.MovieTorrent)
	var errors = make(chan error)

	// yify
	{
		go func() {
			y := yify.New()
			ms, err := y.Search(query)
			if err != nil {
				errors <- fmt.Errorf("yify search failed: %v", err.Error())
				return
			}
			moviesChan <- ms

			return
		}()
	}

	// cpasbien
	{
		go func() {
			c := cpasbien.New()
			ms, err := c.Search(query)
			if err != nil {
				errors <- fmt.Errorf("cpasbien search failed: %v", err.Error())
				return
			}
			moviesChan <- ms

			return
		}()
	}

	movies := results(2, moviesChan, errors, timeout)

	if len(movies) == 0 {
		http.Error(w, fmt.Sprintf("%#q movie not found", query), http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(movies)
}

func results(n int, moviesChan <-chan []movie.MovieTorrent, errors <-chan error, timeout <-chan time.Time) (movies []movie.MovieTorrent) {
	for i := 0; i < n; i++ {
		select {
		case ms := <-moviesChan:
			movies = append(movies, ms...)
		case err := <-errors:
			log.Printf("error: %v\n", err)
		case <-timeout:
			return
		}
	}

	return
}
