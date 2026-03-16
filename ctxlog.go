// Package ctxlog provides a way to add context fields to a context.
package ctxlog

import (
	"context"
	"strings"
)

type ctxLogKey struct{}

type kv map[string]any

var ctxLogKeyFields = ctxLogKey{}

func New(ctx context.Context) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}

	if _, ok := ctx.Value(ctxLogKeyFields).(kv); ok {
		return ctx
	}

	return context.WithValue(ctx, ctxLogKeyFields, kv{})
}

func Add(ctx context.Context, key string, value any) {
	if ctx == nil {
		return
	}
	curr, ok := ctx.Value(ctxLogKeyFields).(kv)
	if !ok {
		return
	}

	key = strings.ReplaceAll(key, " ", "_")

	curr[key] = value
}

func Get(ctx context.Context) kv {
	if ctx == nil {
		return nil
	}
	curr, ok := ctx.Value(ctxLogKeyFields).(kv)
	if !ok {
		return nil
	}
	return curr
}
