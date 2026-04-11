package logger

import (
	"context"
	"log/slog"
	"os"
)

type ctxKey struct{}

var loggerKey = ctxKey{}
var defaultLogger *slog.Logger

func InitDefaultLogger() {
	if defaultLogger != nil {
		return
	}

	handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})
	defaultLogger = slog.New(handler)
}

func Default() *slog.Logger {
	if defaultLogger == nil {
		InitDefaultLogger()
	}
	return defaultLogger
}

func WithLogger(ctx context.Context, l *slog.Logger) context.Context {
	if l == nil {
		l = Default()
	}
	return context.WithValue(ctx, loggerKey, l)
}

func FromContext(ctx context.Context) *slog.Logger {
	if ctx == nil {
		return Default()
	}
	if l, ok := ctx.Value(loggerKey).(*slog.Logger); ok && l != nil {
		return l
	}
	return Default()
}

func WithAttrs(ctx context.Context, attrs ...slog.Attr) context.Context {
	args := make([]any, len(attrs))
	for i, attr := range attrs {
		args[i] = attr
	}

	logger := FromContext(ctx).With(args...)
	return WithLogger(ctx, logger)
}

func Debug(msg string, args ...any) {
	Default().Debug(msg, args...)
}

func Info(msg string, args ...any) {
	Default().Info(msg, args...)
}

func Warn(msg string, args ...any) {
	Default().Warn(msg, args...)
}

func Error(msg string, args ...any) {
	Default().Error(msg, args...)
}

func InfoContext(ctx context.Context, msg string, args ...any) {
	FromContext(ctx).Info(msg, args...)
}

func ErrorContext(ctx context.Context, msg string, args ...any) {
	FromContext(ctx).Error(msg, args...)
}
