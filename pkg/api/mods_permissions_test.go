package api

import (
	"context"
	"fmt"
	"net/http"
	"testing"
)

func TestModsService_GrantModPermissions(t *testing.T) {
	client, mux, teardown := setup()
	defer teardown()

	input := &grantPermissionDto{
		Username:   "test",
		Permission: PermissionTypeReadonly,
		Tag:        "qa",
	}

	mux.HandleFunc("/mods/realitymod/grant", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		testBodyJSON(t, r, input)

		fmt.Fprint(w, `{"message": "OK"}`)
	})

	ctx := context.Background()
	_, err := client.Mods.GrantModPermissions(ctx, "realitymod", "qa", "test", PermissionTypeReadonly)
	if err != nil {
		t.Errorf("Mods.GrantModPermissions returned error: %v", err)
	}
}

func TestModsService_RevokeModPermissions(t *testing.T) {
	client, mux, teardown := setup()
	defer teardown()

	input := &grantPermissionDto{
		Username: "test",
		Tag:      "qa",
	}

	mux.HandleFunc("/mods/realitymod/revoke", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		testBodyJSON(t, r, input)

		fmt.Fprint(w, `{"message": "OK"}`)
	})

	ctx := context.Background()
	_, err := client.Mods.RevokeModPermissions(ctx, "realitymod", "qa", "test")
	if err != nil {
		t.Errorf("Mods.RevokeModPermissions returned error: %v", err)
	}
}
