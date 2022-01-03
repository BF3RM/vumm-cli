package registry

import (
	"bytes"
	"fmt"
	"github.com/Masterminds/semver"
	"io"
	"io/ioutil"
	"net/http"
)

func FetchModVersionArchive(mod string, version *semver.Version) (io.Reader, int64, error) {
	fetchUrl := fmt.Sprintf("/mods/%s/%s/download", mod, version)
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
		return nil, 0, GenericError{res.StatusCode, fmt.Sprintf("fetch %s@%s mod archive rejected", mod, version)}
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, 0, err
	}

	reader := bytes.NewReader(body)

	return reader, res.ContentLength, nil
}
