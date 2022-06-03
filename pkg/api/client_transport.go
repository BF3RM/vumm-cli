package api

import "net/http"

var _ http.RoundTripper = (*TokenAuthTransport)(nil)

type TokenAuthTransport struct {
	Token string
}

func (t TokenAuthTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	if t.Token != "" {
		req.Header.Set("Authorization", t.Token)
	}

	return http.DefaultTransport.RoundTrip(req)
}

type commonService struct {
	client *Client
}
