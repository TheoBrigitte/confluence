package tmdb

import (
	"net/url"
)

var (
	baseURL, _ = url.Parse("https://api.themoviedb.org")

	popularURL = "/3/movie/popular"

	externalIDURLFormat = "/3/movie/%d/external_ids"

	imageHost = "https://image.tmdb.org/t/p"
)
