package registry

import (
	"fmt"
	"net/http"
)

func GrantModUserPermissions(mod, user, permission string) error {
	req, err := newRequest(http.MethodPost, fmt.Sprintf("/mods/%s/grant?user=%s&permission=%s", mod, user, permission), nil)
	if err != nil {
		return err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return catchResponseError(res, fmt.Sprintf("failed granting permission %s of user %s on mod %s", permission, user, mod))
	}

	return nil
}

func RevokeModUserPermissions(mod, user string) error {
	req, err := newRequest(http.MethodPost, fmt.Sprintf("/mods/%s/revoke?user=%s", mod, user), nil)
	if err != nil {
		return err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return catchResponseError(res, fmt.Sprintf("failed revoking permissions of user %s on mod %s", user, mod))
	}

	return nil
}
