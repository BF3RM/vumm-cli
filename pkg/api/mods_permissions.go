package api

import (
	"context"
	"fmt"
	"net/http"
)

type PermissionType string

const (
	PermissionTypeReadonly PermissionType = "Readonly"
	PermissionTypePublish                 = "Publish"
)

type grantPermissionDto struct {
	Username   string         `json:"username"`
	Permission PermissionType `json:"permission,omitempty"`
	Tag        string         `json:"tag"`
}

func PermissionTypeFromString(in string) (permission PermissionType, err error) {
	if in == "" {
		return PermissionTypeReadonly, nil
	}

	switch in {
	case "readonly":
		permission = PermissionTypeReadonly
	case "publish":
		permission = PermissionTypePublish
	default:
		err = fmt.Errorf("%s is not a valid permission type", in)
	}

	return
}

func (s ModsService) GrantModPermissions(ctx context.Context, modName string, modTag string, username string, permission PermissionType) (*http.Response, error) {
	req, err := s.client.NewRequest(http.MethodPost, fmt.Sprintf("mods/%s/grant", modName), &grantPermissionDto{
		Tag:        modTag,
		Username:   username,
		Permission: permission,
	})
	if err != nil {
		return nil, err
	}

	var res interface{}
	return s.client.Do(ctx, req, &res)
}

func (s ModsService) RevokeModPermissions(ctx context.Context, modName string, modTag string, username string) (*http.Response, error) {
	req, err := s.client.NewRequest(http.MethodPost, fmt.Sprintf("mods/%s/revoke", modName), &grantPermissionDto{
		Tag:      modTag,
		Username: username,
	})
	if err != nil {
		return nil, err
	}

	var res interface{}
	return s.client.Do(ctx, req, &res)
}
