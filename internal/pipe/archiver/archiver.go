package archiver

import (
	"bytes"
	"github.com/apex/log"
	"github.com/vumm/cli/internal/common"
	"github.com/vumm/cli/internal/context"
	"github.com/vumm/cli/pkg/tar"
	"os"
	"path/filepath"
)

type Pipe struct {
	Store bool
}

func (Pipe) String() string {
	return "archiver"
}

func (p Pipe) Run(ctx *context.Context) error {
	log.Info("compressing files")

	packager := tar.NewPackager()
	packager.SetFileFilter(func(filePath string) bool {
		return !ctx.Project.Ignorer.Matches(filePath)
	})

	var buf bytes.Buffer
	if err := packager.Compress(ctx.Project.Directory, &buf); err != nil {
		return err
	}
	log.Infof("compressed files to archive of %s", common.ByteCountToHuman(int64(buf.Len())))

	if p.Store {
		log.WithField("file", "archive.tar.gz").Infof("saving archive")
		file, err := os.Create(filepath.Join(ctx.WorkingDirectory, "archive.tar.gz"))
		if err != nil {
			return err
		}
		defer file.Close()
		n, err := file.Write(buf.Bytes())
		if err != nil {
			return err
		}
		log.Debugf("wrote %s to %s", common.ByteCountToHuman(int64(n)), file.Name())
	}

	ctx.SetValue("archive", &buf)
	return nil
}
