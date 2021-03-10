package registry

import "errors"

var Url = "http://localhost:5000/api"

var (
	ErrModVersionNotFound = errors.New("mod version was not found")
)
