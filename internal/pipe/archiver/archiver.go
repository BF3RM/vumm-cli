package archiver

import (
	"bytes"
	"fmt"
	"github.com/apex/log"
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
	log.Infof("compressed files to archive of %s", byteCountToHuman(int64(buf.Len())))

	ctx.SetValue("archive", &buf)
	return nil
}

func byteCountToHuman(b int64) string {
	const unit = 1000
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(b)/float64(div), "kMGTPE"[exp])
}
