package confluence

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/TheoBrigitte/confluence/pkg/movie"
	"github.com/TheoBrigitte/confluence/pkg/movie/provider/cpasbien"
	"github.com/TheoBrigitte/confluence/pkg/movie/provider/yify"
)

func searchHandler(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	query := q.Get("query")
	log.Printf("query: %#q\n", query)

	var movies []movie.MovieTorrent
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

	for i := 0; i < 2; i++ {
		select {
		case ms := <-moviesChan:
			movies = append(movies, ms...)
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
