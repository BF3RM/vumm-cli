package installer

import (
	"bufio"
	"github.com/vumm/cli/common"
	"os"
	"path/filepath"
	"strings"
)

type Installer struct {
	cwd       string
	installed map[string]common.ModMetadata
	missing   map[string]ModDependency
}

func NewInstaller() (*Installer, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	installer := &Installer{
		cwd:       cwd,
		installed: map[string]common.ModMetadata{},
		missing:   map[string]ModDependency{},
	}

	if err = installer.loadInstalledMods(); err != nil {
		return nil, err
	}

	installer.loadMissingMods()

	return installer, nil
}

func (i Installer) HasMissingMods() bool {
	return len(i.missing) > 0
}

func (i Installer) IsInstalled(mod string) bool {
	_, found := i.installed[strings.ToLower(mod)]

	return found
}

func (i Installer) GetInstalledMods() []common.ModMetadata {
	res := make([]common.ModMetadata, 0, len(i.installed))
	for _, metadata := range i.installed {
		res = append(res, metadata)
	}

	return res
}

func (i Installer) GetMissingMods() []ModDependency {
	res := make([]ModDependency, 0, len(i.missing))
	for _, dep := range i.missing {
		res = append(res, dep)
	}

	return res
}

func (i Installer) InstallMod(mod string) error {
	dep := ResolveModDependencyFromString(mod)
	missing := i.GetMissingMods()
	missing = append(missing, dep)

	resolver := newDependencyResolver(missing...)

	if err := resolver.Resolve(); err != nil {
		return err
	}

	// TODO: Download/install the resolved dependencies

	return nil
}

func (i *Installer) loadInstalledMods() error {
	file, err := os.Open(filepath.Join(i.cwd, "ModList.txt"))
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}

		return err
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		modName := strings.ToLower(strings.TrimSpace(scanner.Text()))
		// Skip disabled mods
		if strings.HasPrefix(modName, "#") {
			continue
		}

		metadata, err := common.LoadModMetadata(filepath.Join(i.cwd, "Mods", modName, "mod.json"))
		if err != nil {
			if os.IsNotExist(err) {
				continue
			}
			return err
		}
		i.installed[metadata.Name] = metadata
	}
	if err = scanner.Err(); err != nil {
		return err
	}

	return nil
}

func (i *Installer) loadMissingMods() {

	for _, mod := range i.installed {
		// TODO: Check for outdated mods as well
		for dep, constraints := range mod.Dependencies {
			// We already know we are missing this dependency
			if _, ok := i.missing[dep]; ok {
				continue
			}

			if common.IsInternalMod(dep) {
				continue
			}

			if _, ok := i.installed[dep]; !ok {
				i.missing[dep] = ModDependency{dep, "", constraints}
			}
		}
	}
}
