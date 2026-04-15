package logger

import (
	"context"
	"log/slog"
	"testing"
)

// --- Test Default() ---
func TestDefaultLogger_NotNil(t *testing.T) {
	defaultLogger = nil // reset

	l := Default()

	if l == nil {
		t.Fatal("expected default logger, got nil")
	}
}

// --- Test InitDefaultLogger idempotency ---
func TestInitDefaultLogger_Idempotent(t *testing.T) {
	defaultLogger = nil

	InitDefaultLogger()
	first := defaultLogger

	InitDefaultLogger()
	second := defaultLogger

	if first != second {
		t.Error("expected same logger instance, got different")
	}
}

// --- Test WithLogger + FromContext ---
func TestWithLoggerAndFromContext(t *testing.T) {
	ctx := context.Background()

	customLogger := slog.Default()

	ctxWithLogger := WithLogger(ctx, customLogger)

	result := FromContext(ctxWithLogger)

	if result != customLogger {
		t.Error("expected same logger from context")
	}
}

// --- Test WithLogger fallback ---
func TestWithLogger_NilLogger(t *testing.T) {
	ctx := context.Background()

	ctxWithLogger := WithLogger(ctx, nil)

	result := FromContext(ctxWithLogger)

	if result == nil {
		t.Fatal("expected default logger, got nil")
	}
}

// --- Test FromContext with nil ---
func TestFromContext_NilContext(t *testing.T) {
	l := FromContext(nil)

	if l == nil {
		t.Fatal("expected default logger, got nil")
	}
}

// --- Test FromContext fallback ---
func TestFromContext_NoLogger(t *testing.T) {
	ctx := context.Background()

	l := FromContext(ctx)

	if l == nil {
		t.Fatal("expected default logger, got nil")
	}
}

// --- Test WithAttrs ---
func TestWithAttrs(t *testing.T) {
	ctx := context.Background()

	ctx = WithAttrs(ctx, slog.String("key", "value"))

	l := FromContext(ctx)

	if l == nil {
		t.Fatal("expected logger, got nil")
	}
}

// --- Test global logging functions (no panic) ---
func TestGlobalLoggingFunctions(t *testing.T) {
	Debug("debug message")
	Info("info message")
	Warn("warn message")
	Error("error message")
}

// --- Test context logging functions ---
func TestContextLoggingFunctions(t *testing.T) {
	ctx := context.Background()

	InfoContext(ctx, "info with context")
	ErrorContext(ctx, "error with context")
}
