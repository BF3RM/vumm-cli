package installer

import (
	"bufio"
	"fmt"
	"github.com/apex/log"
	"github.com/vumm/cli/internal/common"
	"github.com/vumm/cli/internal/context"
	"github.com/vumm/cli/internal/middleware"
	"github.com/vumm/cli/internal/registry"
	"github.com/vumm/cli/pkg/tar"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type Pipe struct {
}

func (p Pipe) String() string {
	return "installing dependencies"
}

func (p Pipe) Run(ctx *context.Context) error {
	if len(ctx.Dependencies) == 0 {
		log.Warn("nothing to install")
		return nil
	}

	packager := tar.NewPackager()

	if err := p.loadModList(ctx); err != nil {
		return err
	}

	for _, version := range ctx.Dependencies {
		if err := p.installModVersion(ctx, packager, version); err != nil {
			return err
		}
	}

	p.updateModList(ctx)

	return nil
}

func (p Pipe) installModVersion(ctx *context.Context, packager tar.Packager, version registry.ModVersion) error {
	return middleware.Logging(fmt.Sprintf("installing %s@%s", version.Name, version.Version), func(ctx *context.Context) error {
		log.Info("fetching archive")

		archiveReader, size, err := registry.FetchModVersionArchive(version.Name, version.Version)
		if err != nil {
			return err
		}
		defer archiveReader.Close()
		log.WithField("size", common.ByteCountToHuman(size)).Debugf("streaming archive file")

		modFolder := filepath.Join(ctx.WorkingDirectory, "Mods", version.Name)

		// Check if folder exists
		if _, err := os.Stat(modFolder); err != nil && !os.IsNotExist(err) {
			return err
		}

		if err := os.RemoveAll(modFolder); err != nil {
			return err
		}

		if err := os.MkdirAll(modFolder, os.ModePerm); err != nil {
			return err
		}

		start := time.Now()
		log.Infof("extracting archive")
		err = packager.Decompress(archiveReader, modFolder)
		if err != nil {
			return err
		}
		log.WithField("time", time.Since(start).Truncate(time.Millisecond)).Debugf("extracted archive")

		p.addToEnabledMods(ctx, version.Name)

		return nil
	})(ctx)
}

func (p Pipe) loadModList(ctx *context.Context) error {
	log.WithField("file", "ModList.txt").Infof("loading")

	file, err := os.Open(filepath.Join(ctx.WorkingDirectory, "ModList.txt"))
	if err != nil {
		if os.IsNotExist(err) {
			log.Warn("ModList.txt not found")
			return nil
		}

		return err
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		modName := strings.ToLower(strings.TrimSpace(scanner.Text()))
		// Skip disabled mods
		// TODO: Actually save these somewhere, we are replacing the ModList.txt
		if strings.HasPrefix(modName, "#") {
			continue
		}

		p.addToEnabledMods(ctx, strings.ToLower(modName))
	}
	if err = scanner.Err(); err != nil {
		return err
	}

	return nil
}

func (p Pipe) updateModList(ctx *context.Context) error {
	log.WithField("file", "ModList.txt").Infof("updating")

	tmpFilePath := filepath.Join(ctx.WorkingDirectory, "ModList.txt.new")
	file, err := os.Create(tmpFilePath)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)

	for _, enabledMod := range ctx.EnabledMods {
		if _, err = writer.WriteString(fmt.Sprintf("%s\n", enabledMod)); err != nil {
			return err
		}
	}
	if err = writer.Flush(); err != nil {
		return err
	}

	if err = os.Rename(tmpFilePath, filepath.Join(ctx.WorkingDirectory, "ModList.txt")); err != nil {
		return err
	}

	return nil
}

func (Pipe) addToEnabledMods(ctx *context.Context, mod string) {
	for _, enabledMod := range ctx.EnabledMods {
		if enabledMod == mod {
			return
		}
	}

	ctx.EnabledMods = append(ctx.EnabledMods, mod)
}
