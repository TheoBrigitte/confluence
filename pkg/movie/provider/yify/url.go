package yify

import (
	"net/url"
)

var (
	baseURL, _ = url.Parse("https://yts.mx/api/v2/")

	// Search query example:
	// https://yts.am/api/v2/list_movies.json?query_term=test
	searchURL, _   = url.Parse("list_movies.json")
	SearchQueryKey = "query_term"
	SearchLimitKey = "limit"
	SearchSortKey  = "sort_by"
)
