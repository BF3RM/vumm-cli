package common

import (
	"encoding/json"
	"net/http"
)

func GetHttpJson(url string, out interface{}) error {
	r, err := http.Get(url)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	return json.NewDecoder(r.Body).Decode(out)
}
