package movie

func (m Movie) BestTorrent() MovieTorrent {
	var current Torrent
	for _, t := range m.Torrents {
		if t.Seeds > current.Seeds {
			current = t
		}
	}

	mt := MovieTorrent{
		MovieBase: m.MovieBase,
		Torrent:   current,
	}

	return mt
}
