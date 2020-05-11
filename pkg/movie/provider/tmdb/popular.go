package tmdb

import (
	"fmt"
	"log"
	"time"

	"github.com/TheoBrigitte/confluence/pkg/movie"
	"github.com/TheoBrigitte/confluence/pkg/util"
)

var (
	timeout = time.After(5 * time.Second)
)

func (c *client) Popular(limit int) ([]movie.MovieTorrent, error) {
	res, err := c.getPopularMovies(limit)
	if err != nil {
		return nil, err
	}

	movies, err := c.addTorrent(res.Results)
	if err != nil {
		return nil, err
	}

	return movies, nil
}

func (c client) getPopularMovies(limit int) (*popularResponse, error) {
	u, err := baseURL.Parse(popularURL)
	if err != nil {
		return nil, err
	}

	res, err := c.http.Get(u.String())
	if err != nil {
		return nil, err
	}

	var response popularResponse
	err = util.DecodeJSON(res, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

func (c client) addTorrent(input []movieDetail) (movies []movie.MovieTorrent, err error) {
	movieChan := make(chan movie.MovieTorrent)
	errors := make(chan error)
	done := make(chan bool)

	go func() {
		for i := 0; i < len(input); i++ {
			select {
			case md := <-movieChan:
				movies = append(movies, md)
			case err := <-errors:
				log.Printf("error: %v\n", err)
			case <-timeout:
				log.Println("timeout")
				break
			}
		}
		done <- true
	}()

	for _, inputMovie := range input {
		go func(m movieDetail) {
			external_ids, err := c.ExternalIDs(m.ID)
			if err != nil {
				errors <- err
				return
			}

			yify_movie, err := c.finder.Find(external_ids.IMDB)
			if err != nil {
				errors <- err
				return
			}

			tmdb_movie := m.ToMovieTorrent()
			tmdb_movie.Torrent = yify_movie.Torrent
			tmdb_movie.MovieBase.ExternalID.IMDB = external_ids.IMDB

			movieChan <- tmdb_movie
			return
		}(inputMovie)
	}

	<-done
	if len(movies) == 0 {
		return nil, fmt.Errorf("no popular movies found")
	}
	return movies, nil
}
