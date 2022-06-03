package api

import (
	"bytes"
	"context"
	"fmt"
	"github.com/Masterminds/semver/v3"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"net/http"
	"testing"
)

func TestModsService_GetMod(t *testing.T) {
	client, mux, teardown := setup()
	defer teardown()

	mux.HandleFunc("/mods/realitymod", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")

		fmt.Fprint(w, `
{
	"name": "realitymod",
	"description": "Da best mod",
	"author": "author",
	"tags": {
		"latest": "0.1.0"
	},
	"versions": {
		"0.1.0": {
			"name": "realitymod",
			"description": "Da best mod",
			"author": "author",
			"version": "0.1.0",
			"dependencies": {
				"vemanager": ">=0.1.0"
			}
		}
	}
}`)
	})

	ctx := context.Background()
	result, _, err := client.Mods.GetMod(ctx, "realitymod")
	if err != nil {
		t.Errorf("Mods.GetMod returned error: %v", err)
	}

	expected := &Mod{
		Name:        "realitymod",
		Description: "Da best mod",
		Author:      "author",
		Tags:        map[string]string{"latest": "0.1.0"},
		Versions: map[string]ModVersion{"0.1.0": {
			Name:         "realitymod",
			Description:  "Da best mod",
			Author:       "author",
			Version:      semver.MustParse("0.1.0"),
			Dependencies: map[string]string{"vemanager": ">=0.1.0"},
		}},
	}
	if !cmp.Equal(result, expected, cmpopts.IgnoreUnexported(Mod{})) {
		t.Errorf("Mods.GetMod returned %+v, expected: %+v", result, expected)
	}
}

func TestModsService_DownloadModArchive(t *testing.T) {
	client, mux, teardown := setup()
	defer teardown()

	mux.HandleFunc("/mods/realitymod/0.1.0/download", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")

		fmt.Fprint(w, "test")
	})

	ctx := context.Background()
	result, _, err := client.Mods.DownloadModArchive(ctx, "realitymod", semver.MustParse("0.1.0"))
	if err != nil {
		t.Errorf("Mods.DownloadModArchive returned error: %v", err)
	}

	expected := bytes.NewBufferString("test")
	if !cmp.Equal(result.String(), expected.String()) {
		t.Errorf("Mods.DownloadModArchive returned %+v, expected: %+v", result.String(), expected.String())
	}
}
