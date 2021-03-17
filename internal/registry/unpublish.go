package registry

import (
	"fmt"
	"github.com/Masterminds/semver"
	"net/http"
)

func UnpublishModVersion(mod string, version *semver.Version) error {
	req, err := newRequest(http.MethodDelete, fmt.Sprintf("/mods/%s/%s", mod, version), nil)
	if err != nil {
		return err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	if res.StatusCode != http.StatusNoContent {
		return fmt.Errorf("unpublish rejected: %s", res.Status)
	}

	return nil
}
