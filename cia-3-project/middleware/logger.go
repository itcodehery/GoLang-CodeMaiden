package middleware

import (
	"log/slog"
	"time"

	"github.com/gin-gonic/gin"
)

// StructuredLogger is a Gin middleware that logs requests using Go's slog package.
func StructuredLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		// Process request
		c.Next()

		duration := time.Since(start)
		statusCode := c.Writer.Status()

		attrs := []slog.Attr{
			slog.String("method", c.Request.Method),
			slog.String("path", path),
			slog.Int("status", statusCode),
			slog.Duration("latency", duration),
			slog.String("client_ip", c.ClientIP()),
			slog.String("user_agent", c.Request.UserAgent()),
			slog.Int("body_size", c.Writer.Size()),
		}

		if query != "" {
			attrs = append(attrs, slog.String("query", query))
		}

		if len(c.Errors) > 0 {
			attrs = append(attrs, slog.String("errors", c.Errors.String()))
		}

		// Log level based on status code
		msg := "HTTP Request"
		switch {
		case statusCode >= 500:
			slog.LogAttrs(c.Request.Context(), slog.LevelError, msg, attrs...)
		case statusCode >= 400:
			slog.LogAttrs(c.Request.Context(), slog.LevelWarn, msg, attrs...)
		default:
			slog.LogAttrs(c.Request.Context(), slog.LevelInfo, msg, attrs...)
		}
	}
}
