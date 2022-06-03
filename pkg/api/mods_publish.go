package api

import (
	"context"
	"encoding/base64"
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

func (s ModsService) PublishMod(ctx context.Context, metadata common.ModMetadata, tag string, reader io.Reader) (*http.Response, error) {
	archive, err := s.encodeModArchive(reader)
	if err != nil {
		return nil, fmt.Errorf("failed to encode archive: %v", err)
	}

	publishUrl := fmt.Sprintf("mods/%s/%s", metadata.Name, metadata.Version)

	req, err := s.client.NewRequest(http.MethodPut, publishUrl, &publishModDto{
		ModMetadata: metadata,
		Tag:         tag,
		Archive:     archive,
	})
	if err != nil {
		return nil, err
	}

	return s.client.Do(ctx, req, nil)
}

func (s ModsService) UnpublishModVersion(ctx context.Context, modName string, modVersion *semver.Version) (*http.Response, error) {
	req, err := s.client.NewRequest(http.MethodDelete, fmt.Sprintf("mods/%s/%s", modName, modVersion), nil)
	if err != nil {
		return nil, err
	}

	return s.client.Do(ctx, req, nil)
}

func (s ModsService) encodeModArchive(reader io.Reader) (modArchiveDto, error) {
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
