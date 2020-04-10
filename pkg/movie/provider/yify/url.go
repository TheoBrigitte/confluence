package yify

import (
	"net/url"
)

var (
	baseURL, _ = url.Parse("https://yst.am/api/v2/")

	//https://yts.am/api/v2/list_movies.json?query_term=test
	searchURL, _   = url.Parse("list_movies.json")
	searchQueryKey = "query_term"
)
