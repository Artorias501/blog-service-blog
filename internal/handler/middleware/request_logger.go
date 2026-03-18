package middleware

import (
	"log/slog"
	"time"

	"github.com/gin-gonic/gin"
)

// RequestLogger logs request details including method, path, status, and latency.
func RequestLogger(log *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// Get request ID from header if present
		requestID := c.GetHeader("X-Request-ID")

		// Process request
		c.Next()

		// Calculate latency
		latency := time.Since(start)

		// Build log attributes
		attrs := []slog.Attr{
			slog.String("method", c.Request.Method),
			slog.String("path", c.Request.URL.Path),
			slog.Int("status", c.Writer.Status()),
			slog.String("latency", latency.String()),
			slog.String("client_ip", c.ClientIP()),
			slog.String("user_agent", c.Request.UserAgent()),
			slog.Int("response_size", c.Writer.Size()),
		}

		// Add request ID if present
		if requestID != "" {
			attrs = append(attrs, slog.String("request_id", requestID))
		}

		// Log the request with attributes
		log.LogAttrs(nil, slog.LevelInfo, "HTTP request", attrs...)
	}
}
