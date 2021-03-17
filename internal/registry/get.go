package registry

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func GetMod(mod string) (result Mod, err error) {
	req, err := newRequest(http.MethodGet, fmt.Sprintf("/mods/%s", mod), nil)
	if err != nil {
		return
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return
	}
	if res.StatusCode != http.StatusOK {
		err = fmt.Errorf("get rejected: %s", res.Status)
		return
	}
	defer res.Body.Close()

	err = json.NewDecoder(res.Body).Decode(&result)
	return
}
