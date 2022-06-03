package api

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
)

type Option func(c *Client) error

func BaseURL(baseUrl string) Option {
	return func(c *Client) error {
		parsedUrl, err := url.Parse(baseUrl)
		if err != nil {
			return err
		}
		c.baseUrl = parsedUrl
		return nil
	}
}

func RoundTrip(t http.RoundTripper) Option {
	return func(c *Client) error {
		c.client.Transport = t
		return nil
	}
}

var defaultBaseUrl = "https://vumm.bf3reality.com/api/v1/"

type Client struct {
	baseUrl *url.URL
	client  *http.Client

	common commonService

	Auth *AuthService
	Mods *ModsService
}

func New(opts ...Option) (*Client, error) {
	baseUrl, _ := url.Parse(defaultBaseUrl)

	client := &Client{
		baseUrl: baseUrl,
		client:  &http.Client{},
	}

	for _, option := range opts {
		err := option(client)
		if err != nil {
			return nil, err
		}
	}

	client.common.client = client
	client.Auth = (*AuthService)(&client.common)
	client.Mods = (*ModsService)(&client.common)

	return client, nil
}

func (c *Client) NewRequest(method, path string, body interface{}) (*http.Request, error) {
	reqUrl, err := c.baseUrl.Parse(path)
	if err != nil {
		return nil, err
	}

	var buf io.ReadWriter
	if body != nil {
		buf = &bytes.Buffer{}
		err := json.NewEncoder(buf).Encode(body)
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, reqUrl.String(), buf)
	if err != nil {
		return nil, err
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	return req, nil
}

func (c *Client) DoRequest(ctx context.Context, req *http.Request) (*http.Response, error) {
	req = req.WithContext(ctx)

	res, err := c.client.Do(req)
	if err != nil {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
			return nil, err
		}
	}

	err = c.checkResponse(res)
	return res, err
}

func (c *Client) Do(ctx context.Context, req *http.Request, v interface{}) (*http.Response, error) {
	res, err := c.DoRequest(ctx, req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	switch v := v.(type) {
	case nil:
	case io.Writer:
		_, err = io.Copy(v, res.Body)
	default:
		err = json.NewDecoder(res.Body).Decode(v)
	}

	return res, err
}

func (c *Client) checkResponse(res *http.Response) error {
	// All statuses between 200 <-> 299 are ok
	if c := res.StatusCode; 200 <= c && c <= 299 {
		return nil
	}

	resErr := &GenericError{Response: res}
	data, err := ioutil.ReadAll(res.Body)
	if err == nil && data != nil {
		json.Unmarshal(data, resErr)
	}

	switch res.StatusCode {
	case http.StatusUnauthorized:
		return (*UnauthorizedError)(resErr)
	case http.StatusBadRequest:
		return (*BadRequestError)(resErr)
	case http.StatusConflict:
		return (*ConflictError)(resErr)
	default:
		return resErr
	}
}
