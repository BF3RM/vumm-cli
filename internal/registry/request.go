package registry

import (
	"fmt"
	"io"
	"net/http"
)

func newRequest(method string, endpoint string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(method, fmt.Sprintf("%s%s", registryUrl, endpoint), body)
	if err != nil {
		return nil, err
	}

	if registryToken != "" {
		req.Header.Set("Authorization", registryToken)
	}

	return req, nil
}
