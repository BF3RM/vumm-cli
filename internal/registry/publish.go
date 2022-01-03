package registry

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/vumm/cli/internal/common"
	"io"
	"io/ioutil"
	"net/http"
)

type publishModArchiveDto struct {
	Data        string `json:"data"`
	Length      int    `json:"length"`
	ContentType string `json:"content_type"`
}

type publishModDto struct {
	common.ModMetadata
	Tag     string               `json:"tag"`
	Archive publishModArchiveDto `json:"archive"`
}

func PublishMod(metadata common.ModMetadata, tag string, reader io.Reader) error {
	archive, err := encodeArchive(reader)
	if err != nil {
		return fmt.Errorf("failed to encode archive: %v", err)
	}

	var buf bytes.Buffer
	if err = json.NewEncoder(&buf).Encode(publishModDto{
		ModMetadata: metadata,
		Tag:         tag,
		Archive:     archive,
	}); err != nil {
		return fmt.Errorf("failed to encode publish metadata: %v", err)
	}

	publishUrl := fmt.Sprintf("/mods/%s/%s", metadata.Name, metadata.Version)

	req, err := newRequest(http.MethodPut, publishUrl, &buf)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusCreated {
		return GenericError{res.StatusCode, "publish rejected"}
	}

	return nil
}

func encodeArchive(reader io.Reader) (publishModArchiveDto, error) {
	buffer, err := ioutil.ReadAll(reader)
	if err != nil {
		return publishModArchiveDto{}, err
	}

	length := len(buffer)
	contentType := http.DetectContentType(buffer)
	data := base64.StdEncoding.EncodeToString(buffer)

	return publishModArchiveDto{
		Data:        data,
		Length:      length,
		ContentType: contentType,
	}, nil
}
