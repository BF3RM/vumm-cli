package registry

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type grantPermissionDto struct {
	Username   string `json:"username"`
	Permission string `json:"permission"`
}

func GrantModUserPermissions(mod, username, permission string) error {
	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(grantPermissionDto{
		Username:   username,
		Permission: permission,
	})
	if err != nil {
		return err
	}

	req, err := newRequest(http.MethodPost, fmt.Sprintf("/mods/%s/grant", mod), &buf)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return catchResponseError(res, fmt.Sprintf("failed granting permission %s to user %s on mod %s", permission, username, mod))
	}

	return nil
}

func RevokeModUserPermissions(mod, username string) error {
	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(grantPermissionDto{
		Username: username,
	})
	if err != nil {
		return err
	}

	req, err := newRequest(http.MethodPost, fmt.Sprintf("/mods/%s/revoke", mod), &buf)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return catchResponseError(res, fmt.Sprintf("failed revoking permissions of user %s on mod %s", username, mod))
	}

	return nil
}
