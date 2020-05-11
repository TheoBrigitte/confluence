package tmdb

import (
	"fmt"

	"github.com/TheoBrigitte/confluence/pkg/util"
)

func (c *client) ExternalIDs(id int) (*movieDetail, error) {
	externalIDURL := fmt.Sprintf(externalIDURLFormat, id)
	u, err := baseURL.Parse(externalIDURL)
	if err != nil {
		return nil, err
	}

	res, err := c.http.Get(u.String())
	if err != nil {
		return nil, err
	}

	var response movieDetail
	err = util.DecodeJSON(res, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}
