package publish

import (
	"bytes"
	"fmt"
	"github.com/apex/log"
	"github.com/vumm/cli/internal/context"
	"github.com/vumm/cli/internal/registry"
)

type Pipe struct {
}

func (Pipe) String() string {
	return "publisher"
}

func (Pipe) Run(ctx *context.Context) error {
	var archiveBuf bytes.Buffer
	if !ctx.ValueAs("archive", &archiveBuf) {
		return fmt.Errorf("missing archive buffer")
	}

	log.Info("publishing to registry")
	err := registry.PublishMod(ctx.Project.Metadata, "latest", &archiveBuf)
	if err != nil {
		return err
	}
	log.Info("published successfully")

	return nil
}
