package request

import (
	"context"
	"net/http"

	"log/slog"

	"github.com/nhxgo/go-utils/logger"
)

type ctxKey struct{}

var loggerKey = ctxKey{}

const RequestIDHeader = "X-Request-ID"

type Context struct {
	context.Context
	Logger *slog.Logger
}

func (c *Context) Add(args ...any) {
	c.Logger = c.Logger.With(args...)
}

func (c *Context) Info(message string, args ...any) {
	c.Logger.Info(message, args...)
}

func (c *Context) Error(message string, args ...any) {
	c.Logger.Error(message, args...)
}

func WithLogger(ctx context.Context, l *Context) context.Context {
	return context.WithValue(ctx, loggerKey, l)
}

func FromContext(ctx context.Context) *Context {
	if l, ok := ctx.Value(loggerKey).(*Context); ok {
		return &Context{Context: ctx, Logger: l.Logger}
	}
	return &Context{Context: ctx, Logger: logger.Default()}
}

func GetContext(r *http.Request) *Context {
	return FromContext(r.Context())
}
