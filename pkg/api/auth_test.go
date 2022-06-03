package api

import (
	"context"
	"fmt"
	"github.com/google/go-cmp/cmp"
	"net/http"
	"testing"
	"time"
)

func TestAuthService_Login(t *testing.T) {
	client, mux, teardown := setup()
	defer teardown()

	input := &credentialsDto{
		Username: "test",
		Password: "test",
		Type:     PermissionTypePublish,
	}

	mux.HandleFunc("/auth/login", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		testBodyJSON(t, r, input)

		fmt.Fprint(w, `{"token": "1234", "type": "Publish", "createdAt": "2022-05-03T12:00:00.000Z"}`)
	})

	ctx := context.Background()
	result, _, err := client.Auth.Login(ctx, "test", "test", PermissionTypePublish)
	if err != nil {
		t.Errorf("Auth.Login returned error: %v", err)
	}

	expected := &AuthResult{
		Token:     "1234",
		Type:      PermissionTypePublish,
		CreatedAt: time.Date(2022, 5, 3, 12, 00, 00, 00, time.UTC),
	}
	if !cmp.Equal(result, expected) {
		t.Errorf("Auth.Login returned %+v, expected: %+v", result, expected)
	}
}

func TestAuthService_Register(t *testing.T) {
	client, mux, teardown := setup()
	defer teardown()

	input := &credentialsDto{
		Username: "test",
		Password: "test",
		Type:     PermissionTypePublish,
	}

	mux.HandleFunc("/auth/register", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		testBodyJSON(t, r, input)

		fmt.Fprint(w, `{"token": "1234", "type": "Publish", "createdAt": "2022-05-03T12:00:00.000Z"}`)
	})

	ctx := context.Background()
	result, _, err := client.Auth.Register(ctx, "test", "test", PermissionTypePublish)
	if err != nil {
		t.Errorf("Auth.Register returned error: %v", err)
	}

	expected := &AuthResult{
		Token:     "1234",
		Type:      PermissionTypePublish,
		CreatedAt: time.Date(2022, 5, 3, 12, 00, 00, 00, time.UTC),
	}
	if !cmp.Equal(result, expected) {
		t.Errorf("Auth.Register returned %+v, expected: %+v", result, expected)
	}
}
