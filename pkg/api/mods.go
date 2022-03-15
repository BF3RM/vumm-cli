package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/Masterminds/semver/v3"
	"github.com/vumm/cli/internal/common"
	"net/http"
	"sort"
)

type ModVersion struct {
	Name         string            `json:"name"`
	Description  string            `json:"description"`
	Author       string            `json:"author"`
	Version      *semver.Version   `json:"version"`
	Dependencies map[string]string `json:"dependencies"`
}

func (m ModVersion) String() string {
	str, err := json.Marshal(m)
	if err != nil {
		panic(err)
	}

	return string(str)
}

// ModVersions is a collection of ModVersion instances and implements the sort
// interface based on the semver version inside it.
type ModVersions []ModVersion

// Len returns the length of a collection.
func (c ModVersions) Len() int {
	return len(c)
}

// Less checks if one is greater (reverse) than the other based on Semver.
func (c ModVersions) Less(i, j int) bool {
	return c[i].Version.GreaterThan(c[j].Version)
}

func (c ModVersions) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}

type Mod struct {
	Name        string                `json:"name"`
	Description string                `json:"description"`
	Author      string                `json:"author"`
	Tags        map[string]string     `json:"tags"`
	Versions    map[string]ModVersion `json:"versions"`

	// sorted list of mod versions
	versions ModVersions
}

// GetVersionByTag tries to find the version by tag name
func (m Mod) GetVersionByTag(tag string) (ModVersion, error) {
	tagVersion, ok := m.Tags[tag]
	if !ok {
		return ModVersion{}, ErrModVersionNotFound
	}

	version, ok := m.Versions[tagVersion]
	if !ok {
		return ModVersion{}, ErrModVersionNotFound
	}
	return version, nil
}

// GetLatestVersionByConstraints tries to find the latest version satisfying the constraints
func (m *Mod) GetLatestVersionByConstraints(constraints *common.SemverConstraints) (ModVersion, error) {
	// If versions are not sorted yet, do that now
	if m.versions == nil {
		m.versions = make(ModVersions, 0, len(m.Versions))
		for _, version := range m.Versions {
			m.versions = append(m.versions, version)
		}
		sort.Sort(m.versions)
	}

	for _, modVersion := range m.versions {
		if constraints.Check(modVersion.Version) {
			return modVersion, nil
		}
	}

	return ModVersion{}, ErrModVersionNotFound
}

// GetMod fetches a mod from the registry
func (c Client) GetMod(modName string) (*Mod, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/mods/%s", c.baseUrl, modName), nil)
	if err != nil {
		return nil, err
	}

	mod := Mod{}
	if err := c.doJsonRequest(req, &mod); err != nil {
		return nil, err
	}

	return &mod, nil
}

// GetModArchive fetches a mods archive from the registry
func (c Client) GetModArchive(modName string, modVersion *semver.Version) (*bytes.Buffer, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/mods/%s/%s/download", c.baseUrl, modName, modVersion), nil)
	if err != nil {
		return nil, err
	}

	res, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	return bytes.NewBuffer(res), nil
}
