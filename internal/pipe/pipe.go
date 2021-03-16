package pipe

import (
	"fmt"
	"github.com/vumm/cli/internal/context"
)

type Pipe interface {
	fmt.Stringer
	Run(ctx *context.Context) error
}
