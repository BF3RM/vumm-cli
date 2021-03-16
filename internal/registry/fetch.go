package registry

import (
	"fmt"
	"github.com/Masterminds/semver"
	"io"
	"net/http"
)

func FetchModVersionArchive(mod string, version *semver.Version) (io.ReadCloser, error) {
	fetchUrl := fmt.Sprintf("%s/mods/%s/%s/archive", Url, mod, version)
	res, err := http.Get(fetchUrl)
	if err != nil {
		return nil, err
	}

	return res.Body, nil
}
