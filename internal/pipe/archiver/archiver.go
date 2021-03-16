package archiver

import (
	"bytes"
	"github.com/apex/log"
	"github.com/vumm/cli/internal/common"
	"github.com/vumm/cli/internal/context"
	"github.com/vumm/cli/pkg/tar"
)

type Pipe struct {
}

func (Pipe) String() string {
	return "archiver"
}

func (Pipe) Run(ctx *context.Context) error {
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

	ctx.SetValue("archive", &buf)
	return nil
}
