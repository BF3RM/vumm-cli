package registry

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/vumm/cli/common"
	"io"
	"mime/multipart"
	"net/http"
)

func PublishMod(metadata common.ModMetadata, tag string, reader io.Reader) error {
	var buf bytes.Buffer
	w := multipart.NewWriter(&buf)

	// Attributes
	attributeWriter, err := w.CreateFormField("attributes")
	if err != nil {
		return err
	}
	if err = json.NewEncoder(attributeWriter).Encode(metadata); err != nil {
		return err
	}

	// Archive
	archiveWriter, err := w.CreateFormFile("archive", "archive.tar.gz")
	if err != nil {
		return err
	}
	if _, err = io.Copy(archiveWriter, reader); err != nil {
		return err
	}

	// Tag
	if tag != "" {
		tagWriter, err := w.CreateFormField("tag")
		if err != nil {
			return err
		}
		if _, err = tagWriter.Write([]byte(tag)); err != nil {
			return err
		}
	}
	w.Close()

	publishUrl := fmt.Sprintf("%s/mods/%s/%s", Url, metadata.Name, metadata.Version)

	req, err := http.NewRequest(http.MethodPut, publishUrl, &buf)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", w.FormDataContentType())

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	if res.StatusCode != http.StatusCreated {
		return fmt.Errorf("publish rejected: %s", res.Status)
	}

	return nil
}
