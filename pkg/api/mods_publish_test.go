package api

import (
	"bytes"
	"context"
	"fmt"
	"github.com/Masterminds/semver"
	"github.com/vumm/cli/internal/common"
	"io"
	"net/http"
	"os"
	"testing"
)

func TestModsService_encodeModArchive(t *testing.T) {
	service := ModsService{}

	file, err := os.Open("testdata/0.1.0.tgz")
	if err != nil {
		t.Errorf("Failed to open test file")
	}
	defer file.Close()

	archive, err := service.encodeModArchive(file)
	if err != nil {
		t.Errorf("Mods.encodeModArchive returned error: %v", err)
	}

	if archive.Data == "" {
		t.Errorf("Archive data is empty")
	}

	if archive.ContentType != "application/x-gzip" {
		t.Errorf("Archive content type: %v, expected application/x-gzip", archive.ContentType)
	}
}

func TestModsService_PublishMod(t *testing.T) {
	client, mux, teardown := setup()
	defer teardown()

	metadata := common.ModMetadata{
		Name:         "realitymod",
		Version:      semver.MustParse("0.1.0"),
		Dependencies: nil,
	}

	file, err := os.Open("testdata/0.1.0.tgz")
	if err != nil {
		t.Errorf("Failed to open test file")
	}
	defer file.Close()
	var buf bytes.Buffer
	archive, _ := client.Mods.encodeModArchive(io.TeeReader(file, &buf))

	input := &publishModDto{
		ModMetadata: metadata,
		Tag:         "qa",
		Archive:     archive,
	}

	mux.HandleFunc("/mods/realitymod/0.1.0", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PUT")
		testBodyJSON(t, r, input)

		fmt.Fprint(w, `{"message": "OK"}`)
	})

	ctx := context.Background()
	_, err = client.Mods.PublishMod(ctx, metadata, "qa", &buf)
	if err != nil {
		t.Errorf("Mods.PublishMod returned error: %v", err)
	}
}

func TestModsService_UnpublishModVersion(t *testing.T) {
	client, mux, teardown := setup()
	defer teardown()

	mux.HandleFunc("/mods/realitymod/0.1.0", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")

		fmt.Fprint(w, `{"message": "OK"}`)
	})

	ctx := context.Background()
	_, err := client.Mods.UnpublishModVersion(ctx, "realitymod", semver.MustParse("0.1.0"))
	if err != nil {
		t.Errorf("Mods.UnpublishModVersion returned error: %v", err)
	}
}
