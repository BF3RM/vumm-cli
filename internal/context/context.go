package context

import (
	"context"
	"github.com/vumm/cli/internal/project"
	"github.com/vumm/cli/internal/workspace"
	"github.com/vumm/cli/pkg/api"
	"reflect"
	"time"
)

type Context struct {
	context.Context

	Project          *project.Project
	WorkingDirectory string
	ModList          *workspace.ModList
	Dependencies     map[string]api.ModVersion
	Client           *api.Client

	values map[interface{}]interface{}
}

func NewWithTimeout(duration time.Duration) (*Context, context.CancelFunc) {
	ctx, cancel := context.WithTimeout(context.Background(), duration)

	return &Context{
		Context: ctx,
		values:  map[interface{}]interface{}{},
	}, cancel
}

// SetValue adds a value to the context for later usage
func (ctx *Context) SetValue(key, val interface{}) {
	ctx.values[key] = val
}

// Value tries to lookup a value within the context
func (ctx *Context) Value(key interface{}) interface{} {
	return ctx.values[key]
}

// ValueAs tries to lookup a value within the context and assign it to target
func (ctx *Context) ValueAs(key, target interface{}) bool {
	val, ok := ctx.values[key]
	if !ok {
		return false
	}

	if target == nil {
		panic("errors: target cannot be nil")
	}
	targetVal := reflect.ValueOf(target)
	targetTyp := targetVal.Type()
	if targetTyp.Kind() != reflect.Ptr || targetVal.IsNil() {
		panic("errors: target must be a non-nil pointer")
	}

	if reflect.TypeOf(val).AssignableTo(targetTyp) {
		targetVal.Elem().Set(reflect.ValueOf(val).Elem())
		return true
	}

	return false
}
