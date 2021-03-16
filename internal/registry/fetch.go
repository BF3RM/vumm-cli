package registry

import (
	"fmt"
	"github.com/Masterminds/semver"
	"io"
	"net/http"
)

func FetchModVersionArchive(mod string, version *semver.Version) (io.ReadCloser, int64, error) {
	fetchUrl := fmt.Sprintf("%s/mods/%s/%s/archive", Url, mod, version)
	res, err := http.Get(fetchUrl)
	if err != nil {
		return nil, 0, err
	}

	return res.Body, res.ContentLength, nil
}
