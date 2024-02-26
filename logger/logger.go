package logger

import (
	"context"

	"go.uber.org/zap"
)

type loggerKey struct{}

// WithContext returns a new context with the logger added.
func WithContext(ctx context.Context, l *zap.Logger) context.Context {
	return context.WithValue(ctx, loggerKey{}, l)
}

// FromContext returns the logger in the context if it exists, otherwise a new
// no-op logger is returned .
func FromContext(ctx context.Context) *zap.Logger {
	if ctx == nil {
		return zap.NewNop()
	}

	v := ctx.Value(loggerKey{})
	if v == nil {
		return zap.NewNop()
	}

	return v.(*zap.Logger)
}
