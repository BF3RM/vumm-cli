package api

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type Client struct {
	baseUrl string
	token   *[]byte
	client  *http.Client
}

func NewClient(baseUrl string) *Client {
	return &Client{
		baseUrl: baseUrl,
		client:  &http.Client{},
	}
}

func (c *Client) SetToken(token *[]byte) {
	c.token = token
}

func (c *Client) doRequest(req *http.Request) ([]byte, error) {
	if c.token != nil {
		req.Header.Set("Authorization", string(*c.token))
	}

	// Do request
	res, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	// Check errors
	if res.StatusCode != http.StatusOK {
		return nil, c.catchErrorResponse(res)
	}

	return ioutil.ReadAll(res.Body)
}

func (c *Client) doJsonRequest(req *http.Request, v interface{}) error {
	// Set headers
	req.Header.Set("Content-Type", "application/json")
	if c.token != nil {
		req.Header.Set("Authorization", string(*c.token))
	}

	// Do request
	res, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	// Check errors
	if res.StatusCode >= http.StatusMultipleChoices {
		return c.catchErrorResponse(res)
	}

	// Unmarshall response
	return json.NewDecoder(res.Body).Decode(&v)
}

func (c Client) catchErrorResponse(res *http.Response) error {
	err := GenericError{res.StatusCode, "unknown error occurred"}

	if res.StatusCode == http.StatusBadRequest {
		valError := ValidationError{GenericError: err}
		parseErr := json.NewDecoder(res.Body).Decode(&valError)
		if parseErr == nil {
			return valError
		}
	}

	return err
}
