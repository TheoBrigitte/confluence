package cpasbien

import (
	"net/url"
)

var (
	baseURL, _ = url.Parse("https://www2.cpasbiens.to")
)

const (
	searchURLFormat = "/search_torrent/%s.html"
)
