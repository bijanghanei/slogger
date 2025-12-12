package logger

import (
	"log/slog"
	"time"

	"github.com/gin-gonic/gin"
)

// RequestLoggerMiddleware creates request-scoped logger + req_id.
// Must be the FIRST middleware.
func RequestLoggerMiddleware(serviceName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// Propagate or generate req_id
		reqID := c.GetHeader("X-Request-ID")
		if reqID == "" {
			reqID = NewReqID()
		}

		// Base fields
		baseLogger := slog.Default().With(
			slog.String("req_id", reqID),
			slog.String("method", c.Request.Method),
			slog.String("path", c.Request.URL.Path),
			slog.String("ip", c.ClientIP()),
		)

		// Store for handlers
		c.Set(string(reqIDKey), reqID)
		c.Set("logger", baseLogger)
		c.Set("start_time", start)

		c.Next()

		// Completion log
		latency := time.Since(start)
		status := c.Writer.Status()

		attrs := []any{
			slog.Int("status", status),
			slog.Duration("latency", latency),
			slog.String("user_agent", c.Request.UserAgent()),
		}

		if len(c.Errors) > 0 {
			errs := make([]string, len(c.Errors))
			for i, err := range c.Errors {
				errs[i] = err.Error()
			}
			attrs = append(attrs, slog.Any("gin_errors", errs))
		}

		switch {
		case status >= 500:
			baseLogger.Error("request completed", attrs...)
		case status >= 400:
			baseLogger.Warn("request completed", attrs...)
		default:
			baseLogger.Info("request completed", attrs...)
		}
	}
}