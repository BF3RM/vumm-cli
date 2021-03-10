package workspace

import (
	"github.com/Masterminds/semver"
	"strings"
)

type ModDependency struct {
	Name               string
	Tag                string
	VersionConstraints *semver.Constraints
}

func ResolveModDependency(name string, version string) ModDependency {
	var err error
	var tag string
	var constraints *semver.Constraints

	// First try to parse constraint
	constraints, err = semver.NewConstraint(version)

	// Else set it as tag
	if err != nil {
		tag = version
	}

	return ModDependency{
		Name:               name,
		Tag:                tag,
		VersionConstraints: constraints,
	}
}

func ResolveModDependencyFromString(mod string) ModDependency {
	parts := strings.SplitN(mod, "@", 2)

	version := "latest"
	if len(parts) > 1 {
		version = parts[1]
	}

	return ResolveModDependency(parts[0], version)
}
