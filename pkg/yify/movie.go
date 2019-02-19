package yify

type movieTorrent struct {
	movieBase
	Torrent torrent `json:"torrent"`
}

func (m movie) BestTorrent() movieTorrent {
	var current torrent
	for _, t := range m.Torrents {
		if t.Seeds > current.Seeds {
			current = t
		}
	}

	mt := movieTorrent{
		movieBase: m.movieBase,
		Torrent:   current,
	}

	return mt
}

func (c client) SearchMoviesWithBestTorrent(query string) ([]movieTorrent, error) {
	res, err := c.SearchMovies(query)
	if err != nil {
		return nil, err
	}

	movies := []movieTorrent{}
	for _, m := range res.Data.Movies {
		movies = append(movies, m.BestTorrent())
	}

	return movies, nil
}
