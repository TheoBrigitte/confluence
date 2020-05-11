package confluence

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/TheoBrigitte/confluence/pkg/movie"
	"github.com/TheoBrigitte/confluence/pkg/movie/provider/yify"
)

var (
	popularLimit = 20
)

func popularHandler(w http.ResponseWriter, r *http.Request) {
	var moviesChan = make(chan []movie.MovieTorrent)
	var errors = make(chan error)

	// yify
	{
		go func() {
			y := yify.New()
			ms, err := y.Popular(popularLimit)
			if err != nil {
				errors <- fmt.Errorf("yify popular failed: %v", err.Error())
				return
			}
			moviesChan <- ms

			return
		}()
	}

	movies := results(1, moviesChan, errors, timeout)

	if len(movies) == 0 {
		http.Error(w, "no popular movies", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(movies)
}
