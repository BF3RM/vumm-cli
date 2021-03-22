package registry

import "errors"

var registryUrl = "https://vumm.bf3reality.com/api"
var registryToken = ""

var (
	ErrModVersionNotFound = errors.New("mod version was not found")
)

func SetRegistryUrl(url string) {
	registryUrl = url
}

func SetRegistryAccessToken(accessToken string) {
	registryToken = accessToken
}
