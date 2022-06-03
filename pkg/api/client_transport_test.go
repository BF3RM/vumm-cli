package api

import (
	"net/http"
	"testing"
)

func TestTokenAuthTransport_RoundTrip_NoToken(t *testing.T) {
	ts := TokenAuthTransport{Token: ""}

	res, _ := http.NewRequest(http.MethodGet, "/ok", nil)
	ts.RoundTrip(res)

	if _, exists := res.Header["Authorization"]; exists {
		t.Errorf("Request should not to contain a Authorization header")
	}
}

func TestTokenAuthTransport_RoundTrip_Token(t *testing.T) {
	expected := "1234"
	ts := TokenAuthTransport{Token: expected}

	res, _ := http.NewRequest(http.MethodGet, "/ok", nil)
	ts.RoundTrip(res)

	if value := res.Header.Get("Authorization"); value != expected {
		t.Errorf("Authorization header: %v, expected: %v", value, expected)
	}
}
