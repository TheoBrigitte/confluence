package confluence

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/TheoBrigitte/confluence/pkg/yify"
)

func searchHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("searchHandler")
	q := r.URL.Query()

	c := yify.New()
	movies, err := c.SearchMoviesWithBestTorrent(q.Get("query"))
	if err != nil {
		http.Error(w, fmt.Sprintf("search failed: %v", err.Error()), http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(movies)
}
