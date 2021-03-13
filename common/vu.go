package common

import (
	"fmt"
	"github.com/Masterminds/semver"
)

// ModMetadata holds info of the mod.json
type ModMetadata struct {
	Name         string                        `json:"Name"`
	Version      *semver.Version               `json:"Version"`
	Dependencies map[string]*SemverConstraints `json:"Dependencies"`
}

func (m ModMetadata) String() string {
	return fmt.Sprintf("%s@%s", m.Name, m.Version)
}
