package api

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/Masterminds/semver"
	"github.com/vumm/cli/internal/common"
	"io"
	"io/ioutil"
	"net/http"
)

type modArchiveDto struct {
	Data        string `json:"data"`
	Length      int    `json:"length"`
	ContentType string `json:"content_type"`
}

type publishModDto struct {
	common.ModMetadata
	Tag     string        `json:"tag"`
	Archive modArchiveDto `json:"archive"`
}

func (c Client) PublishMod(metadata common.ModMetadata, tag string, reader io.Reader) error {
	archive, err := c.encodeModArchive(reader)
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

	publishUrl := fmt.Sprintf("%s/mods/%s/%s", c.baseUrl, metadata.Name, metadata.Version)

	req, err := http.NewRequest(http.MethodPut, publishUrl, &buf)
	if err != nil {
		return err
	}

	var res interface{}
	return c.doJsonRequest(req, &res)
}

func (c Client) UnpublishModVersion(modName string, modVersion *semver.Version) error {
	req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/mods/%s/%s", c.baseUrl, modName, modVersion), nil)
	if err != nil {
		return err
	}

	var res interface{}
	return c.doJsonRequest(req, &res)
}

func (c Client) encodeModArchive(reader io.Reader) (modArchiveDto, error) {
	buf, err := ioutil.ReadAll(reader)
	if err != nil {
		return modArchiveDto{}, err
	}

	length := len(buf)
	contentType := http.DetectContentType(buf)
	data := base64.StdEncoding.EncodeToString(buf)

	return modArchiveDto{
		Data:        data,
		Length:      length,
		ContentType: contentType,
	}, nil
}
