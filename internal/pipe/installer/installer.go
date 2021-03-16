package installer

import (
	"fmt"
	"github.com/apex/log"
	"github.com/vumm/cli/internal/common"
	"github.com/vumm/cli/internal/context"
	"github.com/vumm/cli/internal/middleware"
	"github.com/vumm/cli/internal/registry"
	"github.com/vumm/cli/pkg/tar"
	"os"
	"path/filepath"
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

	for _, version := range ctx.Dependencies {
		if err := p.installModVersion(ctx, packager, version); err != nil {
			return err
		}
	}

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

		return nil
	})(ctx)
}
