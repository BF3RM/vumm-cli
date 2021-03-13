package workspace

import (
	"github.com/Masterminds/semver"
	"github.com/vumm/cli/registry"
	"strings"
)

var modExclusions = []string{"veniceext"}

func shouldSkipModCheck(modName string) bool {
	for _, excludedMod := range modExclusions {
		if strings.ToLower(modName) == excludedMod {
			return true
		}
	}

	return false
}

type ResolvedModDependency struct {
	Name        string
	Version     *semver.Version
	Constraints *semver.Constraints
	URL         string
}

func NewModDependencyGraph(deps ...ModDependency) ModDependencyGraph {
	return ModDependencyGraph{
		unresolved: deps,
		resolved:   make(map[string]ResolvedModDependency),
	}
}

type ModDependencyGraph struct {
	unresolved []ModDependency
	resolved   map[string]ResolvedModDependency
}

func (g *ModDependencyGraph) Resolve() (bool, []error) {
	for len(g.unresolved) > 0 {
		dep := g.popUnresolved()

		resolvedVersion, err := g.resolveModVersion(dep)
		if err != nil {
			return false, []error{err}
		}

		if len(resolvedVersion.Dependencies) > 0 {
			for name, version := range resolvedVersion.Dependencies {
				if shouldSkipModCheck(name) {
					continue
				}

				constrains, err := semver.NewConstraint(version)
				// Should never happen!
				if err != nil {
					panic(err)
				}
				added, errs := g.addDependency(ModDependency{name, "", constrains})
				if !added {
					return false, errs
				}
			}
		}

		g.resolved[resolvedVersion.Name] = ResolvedModDependency{
			Name:        resolvedVersion.Name,
			Version:     resolvedVersion.Version,
			Constraints: dep.VersionConstraints,
			URL:         "", // TODO: Implement file uploads in registry
		}
	}
	return true, nil
}

func (g ModDependencyGraph) GetResolvedDependency(mod string) (ResolvedModDependency, bool) {
	resolved, ok := g.resolved[mod]

	return resolved, ok
}

func (g ModDependencyGraph) GetResolvedDependencies() []ResolvedModDependency {
	res := make([]ResolvedModDependency, 0, len(g.resolved))
	for _, resolved := range g.resolved {
		res = append(res, resolved)
	}

	return res
}

func (g *ModDependencyGraph) popUnresolved() ModDependency {
	dep := g.unresolved[0]
	g.unresolved = g.unresolved[1:]

	return dep
}

func (g *ModDependencyGraph) resolveModVersion(dep ModDependency) (registry.ModVersion, error) {
	mod, err := registry.GetMod(dep.Name)
	if err != nil {
		return registry.ModVersion{}, err
	}
	if dep.Tag != "" {
		return mod.GetVersionByTag(dep.Tag)
	}

	return mod.GetLatestVersionByConstraints(dep.VersionConstraints)
}

func (g *ModDependencyGraph) addDependency(dep ModDependency) (bool, []error) {
	// We already resolved a mod with this name
	// If the versions are not compliant, we need to throw an error
	if resolved, ok := g.resolved[dep.Name]; ok {
		return dep.VersionConstraints.Validate(resolved.Version)
	}

	g.unresolved = append(g.unresolved, dep)
	return true, nil
}
