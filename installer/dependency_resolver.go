package installer

import (
	"github.com/vumm/cli/common"
	"github.com/vumm/cli/registry"
)

type dependencyResolver struct {
	unresolved []ModDependency
	resolved   map[string]registry.ModVersion
}

func newDependencyResolver(deps ...ModDependency) dependencyResolver {
	return dependencyResolver{
		unresolved: deps,
		resolved:   map[string]registry.ModVersion{},
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

		r.resolved[resolvedVersion.Name] = resolvedVersion
	}
	return nil
}

func (r dependencyResolver) GetResolvedMods() []registry.ModVersion {
	res := make([]registry.ModVersion, 0, len(r.resolved))
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
