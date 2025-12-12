# slogger

A lightweight, reusable, production-grade structured logging library for Go services, built on the standard library's log/slog.
Designed for microservices running in containers (Kubernetes, Docker, Cloud Run, etc.) where structured JSON logs, request correlation, and minimal overhead are mandatory.
Zero external logging dependencies. Only two optional runtime deps:

github.com/gin-gonic/gin (for the Gin middlewares)
github.com/google/uuid (for request ID generation)

Features

JSON output by default – perfect for Loki, ELK, Datadog, Splunk, CloudWatch Logs
Global service field on every log line
Request-scoped loggers with automatic req_id, method, path, ip
Automatic request completion logging with latency, status, and proper log levels (Info/Warn/Error)
Panic recovery that logs stack traces with full request context
Source location only in debug (performance win in production)
Zero allocations on hot path
Easy to reuse across all your Go services