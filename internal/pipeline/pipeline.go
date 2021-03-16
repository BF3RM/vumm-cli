package pipeline

import (
	"fmt"
	"github.com/vumm/cli/internal/context"
	"github.com/vumm/cli/internal/middleware"
)

type Pipe interface {
	fmt.Stringer
	Run(ctx *context.Context) error
}

func Run(ctx *context.Context, pipeline ...Pipe) error {
	for _, pipe := range pipeline {
		if err := middleware.Logging(pipe.String(), pipe.Run)(ctx); err != nil {
			return err
		}
	}

	return nil
}
