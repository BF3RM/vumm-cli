package middleware

import (
	"github.com/apex/log"
	"github.com/apex/log/handlers/cli"
	"github.com/vumm/cli/internal/context"
)

const DefaultPadding = 3

// Logging is a middleware that automatically increases padding on the logger
func Logging(title string, next Next) Next {
	return func(ctx *context.Context) error {
		defer func() {
			if cli.Default.Padding >= DefaultPadding*2 {
				cli.Default.Padding /= 2
			}
		}()

		log.Infof(title)
		cli.Default.Padding *= 2

		return next(ctx)
	}
}
