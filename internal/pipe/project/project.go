package project

import (
	"fmt"
	"github.com/vumm/cli/internal/context"
	"github.com/vumm/cli/internal/project"
	"os"
)

type Pipe struct {
}

func (Pipe) String() string {
	return "loading project"
}

func (Pipe) Run(ctx *context.Context) error {
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed receiving working directory: %v", err)
	}
	ctx.Project, err = project.Load(cwd)
	if err != nil {
		return fmt.Errorf("failed loading project: %v", err)
	}

	return nil
}
