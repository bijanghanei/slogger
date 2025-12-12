package logger

import (
	"context"
	"log/slog"
	"os"

	"github.com/google/uuid"
)

// ctxKey is a private type to avoid collisions in context values
type ctxKey string

const (
	reqIDKey ctxKey = "req_id"
)

// Default returns the global structured logger with service-specific fields.
// Use this for background jobs, startup, or when no request context is available.
func Default(serviceName string) *slog.Logger {
	return slog.Default().With(
		slog.String("service", serviceName),
	)
}

// FromCtx returns a request-scoped logger from context.
// Falls back to Default() if no logger in context.
func FromCtx(ctx context.Context, serviceName string) *slog.Logger {
	if l, ok := ctx.Value("logger").(*slog.Logger); ok {
		return l
	}
	return Default(serviceName)
}

// WithReqID creates a child logger with request ID (for non-HTTP use)
func WithReqID(parent *slog.Logger, reqID string) *slog.Logger {
	return parent.With(slog.String(string(reqIDKey), reqID))
}

// NewReqID generates a new request ID (UUID v4)
func NewReqID() string {
	return uuid.New().String()
}

// Init initializes the global slog handler as JSON in production.
// Call this once at service startup.
func Init(serviceName string, level slog.Level) {
	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level:     level,
		AddSource: level >= slog.LevelDebug, // source only in debug
	})

	slog.SetDefault(slog.New(handler).With(
		slog.String("service", serviceName),
	))
}