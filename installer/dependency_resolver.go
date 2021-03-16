package installer

import (
	"fmt"
	"github.com/Masterminds/semver"
	"github.com/vumm/cli/common"
	"github.com/vumm/cli/internal/registry"
)

type ResolvedModDependencyStatus uint8

func (s ResolvedModDependencyStatus) String() string {
	switch s {
	case DependencyStatusInstalled:
		return "up to date"
	case DependencyStatusOutdated:
		return "outdated"
	case DependencyStatusNew:
		return "not installed"
	}
	return "unknown"
}

const (
	DependencyStatusNew ResolvedModDependencyStatus = iota
	DependencyStatusOutdated
	DependencyStatusInstalled
)

type ResolvedModDependency struct {
	Name    string
	Version *semver.Version
	Status  ResolvedModDependencyStatus
}

func (d ResolvedModDependency) String() string {
	return fmt.Sprintf("%s@%s - %s", d.Name, d.Version, d.Status)
}

type dependencyResolver struct {
	cwd        string
	installed  map[string]common.ModMetadata
	unresolved []ModDependency
	resolved   map[string]ResolvedModDependency
}

func newDependencyResolver(installed map[string]common.ModMetadata, deps ...ModDependency) dependencyResolver {
	return dependencyResolver{
		installed:  installed,
		unresolved: deps,
		resolved:   map[string]ResolvedModDependency{},
	}
}

func (r *dependencyResolver) Resolve() error {
	for len(r.unresolved) > 0 {
		dep := r.popUnresolved()

		resolvedVersion, err := r.resolveModVersion(dep)
		if err != nil {
			return err
		}

		if len(resolvedVersion.Dependencies) > 0 {
			for name, version := range resolvedVersion.Dependencies {
				if common.IsInternalMod(name) {
					continue
				}

				constrains, err := common.NewSemverConstraints(version)
				// Should never happen!
				if err != nil {
					panic(err)
				}
				added, err := r.addDependency(ModDependency{name, "", constrains})
				if !added {
					return err
				}
			}
		}

		status := DependencyStatusNew
		if installed, ok := r.installed[resolvedVersion.Name]; ok {
			status = DependencyStatusInstalled
			if installed.Version.LessThan(resolvedVersion.Version) {
				status = DependencyStatusOutdated
			}
		}

		r.resolved[resolvedVersion.Name] = ResolvedModDependency{
			Name:    resolvedVersion.Name,
			Version: resolvedVersion.Version,
			Status:  status,
		}
	}
	return nil
}

func (r dependencyResolver) GetResolvedMods() []ResolvedModDependency {
	res := make([]ResolvedModDependency, 0, len(r.resolved))
	for _, mod := range r.resolved {
		res = append(res, mod)
	}

	return res
}

func (r *dependencyResolver) popUnresolved() ModDependency {
	dep := r.unresolved[0]
	r.unresolved = r.unresolved[1:]

	return dep
}

func (r *dependencyResolver) resolveModVersion(dep ModDependency) (registry.ModVersion, error) {
	mod, err := registry.GetMod(dep.Name)
	if err != nil {
		return registry.ModVersion{}, err
	}
	if dep.Tag != "" {
		return mod.GetVersionByTag(dep.Tag)
	}

	return mod.GetLatestVersionByConstraints(dep.VersionConstraints)
}

func (r *dependencyResolver) addDependency(dep ModDependency) (bool, error) {
	// We already resolved a mod with this name
	// If the versions are not compliant, we need to throw an error
	if resolved, ok := r.resolved[dep.Name]; ok {
		ok, errs := dep.VersionConstraints.Validate(resolved.Version)
		if !ok {
			return ok, errs[0]
		}
	}

	r.unresolved = append(r.unresolved, dep)
	return true, nil
}
