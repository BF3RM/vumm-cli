package context

import (
	"context"
	"github.com/vumm/cli/internal/project"
	"reflect"
	"time"
)

type Context struct {
	context.Context
	Project *project.Project

	values map[interface{}]interface{}
}

func NewWithTimeout(duration time.Duration) (*Context, context.CancelFunc) {
	ctx, cancel := context.WithTimeout(context.Background(), duration)

	return &Context{
		Context: ctx,
		values:  map[interface{}]interface{}{},
	}, cancel
}

//
func (ctx *Context) SetValue(key, val interface{}) {
	ctx.values[key] = val
}

func (ctx *Context) Value(key interface{}) interface{} {
	return ctx.values[key]
}

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
