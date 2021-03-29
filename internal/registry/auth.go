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

func Login(username, password string, tokenType AccessTokenType) (AccessToken, error) {
	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(Credentials{
		Username: username,
		Password: password,
		Type:     tokenType,
	})
	if err != nil {
		return AccessToken{}, err
	}

	req, err := newRequest(http.MethodPost, "/auth/login", &buf)
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
		err = fmt.Errorf("login rejected: %s", res.Status)
		return AccessToken{}, err
	}

	var result AccessToken
	err = json.NewDecoder(res.Body).Decode(&result)
	if err != nil {
		return AccessToken{}, err
	}

	return result, nil
}
