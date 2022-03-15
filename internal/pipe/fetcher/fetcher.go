package fetcher

import (
	"fmt"
	"github.com/apex/log"
	"github.com/vumm/cli/internal/common"
	"github.com/vumm/cli/internal/context"
	"github.com/vumm/cli/pkg/api"
)

// Pipe is a pipe that tries resolving the given mod and its dependencies
type Pipe struct {
	mod string
}

func New(mod string) Pipe {
	return Pipe{mod: mod}
}

func (Pipe) String() string {
	return "gathering dependencies"
}

func (p Pipe) Run(ctx *context.Context) error {
	if ctx.Dependencies != nil {
		return fmt.Errorf("dependencies where already resolved")
	}
	ctx.Dependencies = map[string]api.ModVersion{}

	// List with unresolved mod dependencies
	unresolved := []modDependency{resolveModDependencyFromString(p.mod)}

	// TODO: Handle ctx.Done()
	for len(unresolved) > 0 {
		// Pop first unresolved
		dep := unresolved[0]
		unresolved = unresolved[1:]

		log.WithField("mod", dep.Name).Info("fetching metadata")
		version, err := p.resolveModVersion(ctx, dep)
		if err != nil {
			return err
		}
		log.WithFields(log.Fields{
			"mod":     version.Name,
			"version": version.Version,
		}).Debug("fetched metadata")

		if len(version.Dependencies) > 0 {
			for name, version := range version.Dependencies {
				if common.IsInternalMod(name) {
					log.Debugf("skipped %s, is internal mod", name)
					continue
				}

				dep := resolveModDependency(name, version)
				shouldAdd, err := p.checkCanAddDependency(ctx, dep)
				if err != nil {
					return err
				}

				if shouldAdd {
					unresolved = append(unresolved, dep)
				}
			}
		}

		ctx.Dependencies[version.Name] = version
	}

	return nil
}

func (p Pipe) resolveModVersion(ctx *context.Context, dep modDependency) (api.ModVersion, error) {
	mod, err := ctx.Client.GetMod(dep.Name)
	if err != nil {
		return api.ModVersion{}, err
	}
	if dep.Tag != "" {
		return mod.GetVersionByTag(dep.Tag)
	}

	return mod.GetLatestVersionByConstraints(dep.VersionConstraints)
}

func (p Pipe) checkCanAddDependency(ctx *context.Context, dep modDependency) (bool, error) {
	if version, ok := ctx.Dependencies[dep.Name]; ok {
		// TODO: Handle edge case not compatible...
		// Uhm what do we do if the dependency was found but not compatible???
		ok, errors := dep.VersionConstraints.Validate(version.Version)
		if !ok {
			return false, errors[0]
		}

		return false, nil
	}

	return true, nil
}
