package tmdb

import (
	"net/url"
)

var (
	baseURL, _ = url.Parse("https://api.themoviedb.org")

	popularURL = "/movie/popular"

	externalIDURLFormat = "/movie/%d/external_ids"

	imageHost = "https://image.tmdb.org/t/p"
)
