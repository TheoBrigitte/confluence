package movie

type MovieTorrent struct {
	MovieBase
	Torrent Torrent `json:"torrent"`
}

type Movie struct {
	MovieBase
	Torrents []Torrent `json:"torrents"`
}

type MovieBase struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	Year  int    `json:"year"`
}

type Torrent struct {
	Hash      string `json:"hash"`
	Quality   string `json:"quality"`
	Seeds     int    `json:"seeds"`
	Size      string `json:"size"`
	SizeBytes int    `json:"size_bytes"`
	Type      string `json:"type"`
	URL       string `json:"url"`
}
