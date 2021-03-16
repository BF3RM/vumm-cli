package registry

import "errors"

var Url = "https://vumm.bf3reality.com/api/v1"

var (
	ErrModVersionNotFound = errors.New("mod version was not found")
)

func SetRegistryUrl(url string) {
	Url = url
}
