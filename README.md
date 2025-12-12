# slogger

A lightweight, reusable, **production-grade** structured logging library for Go services.

Built exclusively on the standard library's `log/slog` (Go 1.21+), with **zero external logging dependencies**. Designed for microservices running in containers (Kubernetes, Docker, Cloud Run, ECS, etc.) where structured JSON logs, request correlation, and minimal runtime overhead are essential.

## Why This Library Exists

* Eliminates copy‑paste of logging middleware across services
* Guarantees **consistent log format and fields** across an organization
* Provides **professional observability** from day one: request IDs, latency, proper levels, panic context
* Keeps binary size small and dependency tree clean
* Follows modern Go best practices (`slog` as the default structured logger)

## Features

* **JSON output** (structured, machine‑parseable)
* Global `service` field on every log line
* Automatic **request‑scoped logger** with `req_id`, `method`, `path`, `ip`
* Request completion logging with **latency**, **status**, and level‑based routing (Info/Warn/Error)
* **Panic recovery** that includes full request context and stack trace
* Source location (`file:line`) optionally enabled in debug mode
* Zero allocations on the hot path
* Gin‑specific middlewares included (easy to adapt for other frameworks)

## Installation

```bash
go get github.com/yourorg/go-logger@latest
```

## Example Usage (Gin)

A minimal production‑grade Gin service using `slogger`.

```go
package main

import (
    "context"
    "log/slog"
    "net/http"
    "os"
    "os/signal"
    "syscall"
    "time"

    "github.com/gin-gonic/gin"
    "github.com/yourorg/go-logger"
)

func main() {
    // Initialize global logger once at startup.
    // Level respects LOG_LEVEL env var and defaults to info.
    logger.Init("identity-service", slog.LevelInfo)

    gin.SetMode(gin.ReleaseMode)
    r := gin.New()

    // Middleware order matters — logger first.
    r.Use(logger.RequestLoggerMiddleware("identity-service"))
    r.Use(logger.Recovery())

    // Example route
    r.GET("/health", func(c *gin.Context) {
        log := logger.FromCtx(c.Request.Context(), "identity-service")
        log.Info("health check requested")
        c.JSON(200, gin.H{"status": "ok"})
    })

    // Graceful shutdown
    srv := &http.Server{
        Addr:    ":8080",
        Handler: r,
    }

    go func() {
        if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            logger.L().Error("server error", slog.String("err", err.Error()))
        }
    }()

    quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
    <-quit

    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    logger.L().Info("shutting down server")
    if err := srv.Shutdown(ctx); err != nil {
        logger.L().Error("server forced shutdown", slog.String("err", err.Error()))
    }
}
```
