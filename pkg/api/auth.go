package api

import (
	"context"
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

type AuthService commonService

func (s AuthService) Login(ctx context.Context, username string, password string, permission PermissionType) (*AuthResult, *http.Response, error) {
	return s.authRequest(ctx, "login", username, password, permission)
}

func (s AuthService) Register(ctx context.Context, username string, password string, permission PermissionType) (*AuthResult, *http.Response, error) {
	return s.authRequest(ctx, "register", username, password, permission)
}

func (s AuthService) authRequest(ctx context.Context, endpoint, username, password string, permission PermissionType) (*AuthResult, *http.Response, error) {
	req, err := s.client.NewRequest(http.MethodPost, fmt.Sprintf("auth/%s", endpoint), &credentialsDto{
		Username: username,
		Password: password,
		Type:     permission,
	})
	if err != nil {
		return nil, nil, err
	}

	result := new(AuthResult)
	res, err := s.client.Do(ctx, req, result)
	if err != nil {
		return nil, res, err
	}

	return result, res, nil
}
