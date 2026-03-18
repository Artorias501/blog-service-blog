package middleware

import (
	"fmt"
	"log/slog"
	"runtime"

	"github.com/artorias501/blog-service/pkg/response"
	"github.com/gin-gonic/gin"
)

// Recovery recovers from panics and returns a 500 error.
// It logs the panic details with stack trace.
func Recovery(log *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				// Get stack trace
				const depth = 32
				var pcs [depth]uintptr
				n := runtime.Callers(3, pcs[:])
				frames := runtime.CallersFrames(pcs[:n])
				var stackTrace string
				for {
					frame, more := frames.Next()
					stackTrace += fmt.Sprintf("%s\n\t%s:%d\n", frame.Function, frame.File, frame.Line)
					if !more {
						break
					}
				}

				// Log the panic with details
				log.Error("panic recovered",
					slog.Any("panic", r),
					slog.String("method", c.Request.Method),
					slog.String("path", c.Request.URL.Path),
					slog.String("client_ip", c.ClientIP()),
					slog.String("stack", stackTrace),
				)

				// Return 500 error
				response.InternalError(c, "internal server error")
				c.Abort()
			}
		}()
		c.Next()
	}
}

// RecoveryMiddleware is the legacy recovery middleware for backward compatibility.
// It recovers from panics and returns a 500 error without logging.
func RecoveryMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				response.InternalError(c, "internal server error")
				c.Abort()
			}
		}()
		c.Next()
	}
}
