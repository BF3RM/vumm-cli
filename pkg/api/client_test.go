package api

import (
	"encoding/json"
	"github.com/google/go-cmp/cmp"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func setup() (client *Client, mux *http.ServeMux, teardown func()) {
	mux = http.NewServeMux()

	server := httptest.NewServer(mux)

	client, _ = New(BaseURL(server.URL))

	return client, mux, server.Close
}

func testMethod(t *testing.T, r *http.Request, expectedMethod string) {
	t.Helper()
	if r.Method != expectedMethod {
		t.Errorf("Request method: %v, expected %v", r.Method, expectedMethod)
	}
}

func testBodyJSON(t *testing.T, r *http.Request, expected interface{}) {
	t.Helper()

	bodyType := reflect.TypeOf(expected)
	bodyValue := reflect.New(bodyType.Elem())
	value := bodyValue.Interface()

	json.NewDecoder(r.Body).Decode(&value)

	if !cmp.Equal(value, expected) {
		t.Errorf("Request body: %v, expected %v", value, expected)
	}
}

func TestBaseURL(t *testing.T) {
	baseUrl := "https://vumm.com/api/v1/"
	c, _ := New(BaseURL(baseUrl))

	if c.baseUrl.String() != baseUrl {
		t.Errorf("Base url: %v, expected: %v", c.baseUrl, baseUrl)
	}
}

func TestRoundTrip(t *testing.T) {
	ts := &TokenAuthTransport{}
	c, _ := New(RoundTrip(ts))

	if !cmp.Equal(c.client.Transport, ts) {
		t.Errorf("Round trip: %v, expected: %v", c.client.Transport, ts)
	}
}
