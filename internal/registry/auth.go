package registry

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type AccessTokenType uint8

const (
	AccessTokenTypeReadonly AccessTokenType = iota
	AccessTokenTypePublish
)

type Credentials struct {
	Username string          `json:"username"`
	Password string          `json:"password"`
	Type     AccessTokenType `json:"type"`
}

type AccessToken struct {
	Token     string          `json:"token"`
	Type      AccessTokenType `json:"type"`
	CreatedAt time.Time       `json:"createdAt"`
}

func ParseTokenType(in string) (tokenType AccessTokenType, err error) {
	if in == "" {
		return AccessTokenTypePublish, nil
	}

	switch in {
	case "readonly":
		tokenType = AccessTokenTypeReadonly
	case "publish":
		tokenType = AccessTokenTypePublish
	default:
		err = fmt.Errorf("%s is not a valid token type", in)
	}

	return
}

func Login(username, password string, tokenType AccessTokenType) (AccessToken, error) {
	return authRequest("login", username, password, tokenType)
}

func Register(username, password string, tokenType AccessTokenType) (AccessToken, error) {
	return authRequest("register", username, password, tokenType)
}

func authRequest(endpoint, username, password string, tokenType AccessTokenType) (AccessToken, error) {
	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(Credentials{
		Username: username,
		Password: password,
		Type:     tokenType,
	})
	if err != nil {
		return AccessToken{}, err
	}

	req, err := newRequest(http.MethodPost, fmt.Sprintf("/auth/%s", endpoint), &buf)
	if err != nil {
		return AccessToken{}, err
	}
	req.Header.Set("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return AccessToken{}, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return AccessToken{}, catchResponseError(res, fmt.Sprintf("%s failed", endpoint))
	}

	var result AccessToken
	err = json.NewDecoder(res.Body).Decode(&result)
	if err != nil {
		return AccessToken{}, err
	}

	return result, nil
}
