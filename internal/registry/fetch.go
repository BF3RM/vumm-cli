package registry

import (
	"fmt"
	"github.com/Masterminds/semver"
	"io"
	"net/http"
)

func FetchModVersionArchive(mod string, version *semver.Version) (io.ReadCloser, int64, error) {
	fetchUrl := fmt.Sprintf("/mods/%s/%s/archive", mod, version)
	req, err := newRequest(http.MethodGet, fetchUrl, nil)
	if err != nil {
		return nil, 0, err
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, 0, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, 0, fmt.Errorf("fetch archive rejected: %s", res.Status)
	}

	return res.Body, res.ContentLength, nil
}
