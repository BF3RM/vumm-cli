package installer

import (
	"bufio"
	"fmt"
	"github.com/vumm/cli/common"
	"github.com/vumm/cli/registry"
	"github.com/vumm/cli/tar"
	"os"
	"path/filepath"
	"strings"
)

type Installer struct {
	cwd       string
	installed map[string]common.ModMetadata
	missing   map[string]ModDependency
	packager  tar.Packager
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
		packager:  tar.NewPackager(),
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

	fmt.Printf("Installing %s\n", dep.Name)

	missing := i.GetMissingMods()
	missing = append(missing, dep)

	resolver := newDependencyResolver(i.installed, missing...)

	if err := resolver.Resolve(); err != nil {
		return err
	}

	resolvedMods := resolver.GetResolvedMods()

	if len(resolvedMods) == 0 {
		fmt.Printf("%s and it's dependencies are already up to date\n", dep.Name)
		return nil
	}

	fmt.Printf("Installing %d mod(s)\n", len(resolvedMods))

	for _, mod := range resolver.GetResolvedMods() {
		fmt.Printf("\t%s\n", mod)

		if mod.Status != DependencyStatusInstalled {
			if err := i.installModDependency(mod); err != nil {
				return err
			}
		}
	}

	return nil
}

func (i Installer) installModDependency(dep ResolvedModDependency) error {
	reader, err := registry.FetchModVersionArchive(dep.Name, dep.Version)
	if err != nil {
		return err
	}
	defer reader.Close()

	modFolder := filepath.Join(i.cwd, "Mods", dep.Name)

	// Remove old if already exists
	if _, err := os.Stat(modFolder); err != nil {
		if !os.IsNotExist(err) {
			return nil
		}

		if err := os.MkdirAll(modFolder, os.ModeDir); err != nil {
			return fmt.Errorf("%s: failed creating mod folder: %v", modFolder, err)
		}
	} else {
		// TODO: Remove old content...
		//os.RemoveAll()
	}

	err = i.packager.Decompress(reader, modFolder)
	if err != nil {
		return fmt.Errorf("%s: failed decompressing mod: %v", modFolder, err)
	}

	// TODO: Enable in ModList.txt

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
