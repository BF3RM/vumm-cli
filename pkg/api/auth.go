package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type AuthResult struct {
	Token     string         `json:"token"`
	Type      PermissionType `json:"type"`
	CreatedAt time.Time      `json:"createdAt"`
}

type credentialsDto struct {
	Username string         `json:"username"`
	Password string         `json:"password"`
	Type     PermissionType `json:"type"`
}

func (c Client) Login(username string, password string, permission PermissionType) (*AuthResult, error) {
	return c.authRequest("login", username, password, permission)
}

func (c Client) Register(username string, password string, permission PermissionType) (*AuthResult, error) {
	return c.authRequest("register", username, password, permission)
}

func (c Client) authRequest(endpoint, username, password string, permission PermissionType) (*AuthResult, error) {
	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(credentialsDto{
		Username: username,
		Password: password,
		Type:     permission,
	})
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/auth/%s", c.baseUrl, endpoint), &buf)
	if err != nil {
		return nil, err
	}

	var result AuthResult
	if err := c.doJsonRequest(req, &result); err != nil {
		return nil, err
	}

	return &result, nil
}
