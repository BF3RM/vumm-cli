package fetcher

import (
	"github.com/vumm/cli/internal/common"
	"strings"
)

type modDependency struct {
	Name               string
	Tag                string
	VersionConstraints *common.SemverConstraints
}

func resolveModDependency(name string, version string) modDependency {
	var err error
	var tag string
	var constraints *common.SemverConstraints

	// First try to parse constraint
	constraints, err = common.NewSemverConstraints(version)

	// Else set it as tag
	if err != nil {
		tag = version
	}

	return modDependency{
		Name:               name,
		Tag:                tag,
		VersionConstraints: constraints,
	}
}

func resolveModDependencyFromString(mod string) modDependency {
	parts := strings.SplitN(mod, "@", 2)

	version := "latest"
	if len(parts) > 1 {
		version = parts[1]
	}

	return resolveModDependency(parts[0], version)
}
