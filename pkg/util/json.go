package util

import (
	"encoding/json"
	"net/http"
)

func DecodeJSON(res *http.Response, dst interface{}) error {
	return json.NewDecoder(res.Body).Decode(dst)
}
