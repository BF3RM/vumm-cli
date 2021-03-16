package middleware

import "github.com/vumm/cli/internal/context"

type Next func(ctx *context.Context) error
