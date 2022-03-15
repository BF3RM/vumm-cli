package api

import (
	"bytes"
	"encoding/json"
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

func (c Client) GrantModPermissions(modName string, modTag string, username string, permission PermissionType) error {
	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(grantPermissionDto{
		Tag:        modTag,
		Username:   username,
		Permission: permission,
	})
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/mods/%s/grant", c.baseUrl, modName), &buf)
	if err != nil {
		return err
	}

	var res interface{}
	return c.doJsonRequest(req, &res)
}

func (c Client) RevokeModPermissions(modName string, modTag string, username string) error {
	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(grantPermissionDto{
		Tag:      modTag,
		Username: username,
	})
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/mods/%s/revoke", c.baseUrl, modName), &buf)
	if err != nil {
		return err
	}

	var res interface{}
	return c.doJsonRequest(req, &res)
}
