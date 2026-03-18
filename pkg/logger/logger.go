// Package logger provides structured logging utilities for the blog service.
package logger

import (
	"context"
	"io"
	"log/slog"
	"os"
	"runtime"
	"strings"

	"github.com/gin-gonic/gin"
)

var defaultLogger *slog.Logger

// init initializes the default logger.
func init() {
	defaultLogger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
}

// New creates a new logger with the specified environment and output.
// Environment "production" uses JSON format, "development" uses text format.
func New(env string, w io.Writer) *slog.Logger {
	return NewWithLevel(env, w, slog.LevelInfo)
}

// NewWithLevel creates a new logger with the specified environment, output, and level.
func NewWithLevel(env string, w io.Writer, level slog.Level) *slog.Logger {
	if w == nil {
		w = os.Stdout
	}

	opts := &slog.HandlerOptions{
		Level: level,
	}

	var handler slog.Handler
	if strings.ToLower(env) == "production" {
		handler = slog.NewJSONHandler(w, opts)
	} else {
		handler = slog.NewTextHandler(w, opts)
	}

	return slog.New(handler)
}

// GetLogger returns the default logger.
func GetLogger() *slog.Logger {
	return defaultLogger
}

// SetDefault sets the default logger.
func SetDefault(l *slog.Logger) {
	defaultLogger = l
	slog.SetDefault(l)
}

// WithRequestID returns a logger with request_id field.
func WithRequestID(l *slog.Logger, requestID string) *slog.Logger {
	return l.With("request_id", requestID)
}

// Sync flushes any buffered log entries. Most handlers don't need this.
func Sync(l *slog.Logger) error {
	// slog doesn't have a Sync method, but we provide this for compatibility
	return nil
}

// FromContext retrieves the logger from gin context.
func FromContext(c *gin.Context) (*slog.Logger, bool) {
	v, exists := c.Get("logger")
	if !exists {
		return nil, false
	}
	l, ok := v.(*slog.Logger)
	return l, ok
}

// SetInContext sets the logger in gin context.
func SetInContext(c *gin.Context, l *slog.Logger) {
	c.Set("logger", l)
}

// Debug logs at DEBUG level.
func Debug(msg string, args ...any) {
	defaultLogger.Debug(msg, args...)
}

// Info logs at INFO level.
func Info(msg string, args ...any) {
	defaultLogger.Info(msg, args...)
}

// Warn logs at WARN level.
func Warn(msg string, args ...any) {
	defaultLogger.Warn(msg, args...)
}

// Error logs at ERROR level.
func Error(msg string, args ...any) {
	defaultLogger.Error(msg, args...)
}

// DebugContext logs at DEBUG level with context.
func DebugContext(ctx context.Context, msg string, args ...any) {
	defaultLogger.DebugContext(ctx, msg, args...)
}

// InfoContext logs at INFO level with context.
func InfoContext(ctx context.Context, msg string, args ...any) {
	defaultLogger.InfoContext(ctx, msg, args...)
}

// WarnContext logs at WARN level with context.
func WarnContext(ctx context.Context, msg string, args ...any) {
	defaultLogger.WarnContext(ctx, msg, args...)
}

// ErrorContext logs at ERROR level with context.
func ErrorContext(ctx context.Context, msg string, args ...any) {
	defaultLogger.ErrorContext(ctx, msg, args...)
}

// GetStackTrace returns a formatted stack trace.
func GetStackTrace() string {
	const depth = 32
	var pcs [depth]uintptr
	n := runtime.Callers(3, pcs[:])
	if n == 0 {
		return ""
	}

	frames := runtime.CallersFrames(pcs[:n])
	var sb strings.Builder
	for {
		frame, more := frames.Next()
		sb.WriteString(frame.Function)
		sb.WriteString("\n\t")
		sb.WriteString(frame.File)
		sb.WriteString(":")
		sb.WriteString(intToStr(frame.Line))
		sb.WriteString("\n")
		if !more {
			break
		}
	}
	return sb.String()
}

// intToStr converts int to string without importing strconv.
func intToStr(n int) string {
	if n == 0 {
		return "0"
	}
	var neg bool
	if n < 0 {
		neg = true
		n = -n
	}
	var digits []byte
	for n > 0 {
		digits = append([]byte{byte('0' + n%10)}, digits...)
		n /= 10
	}
	if neg {
		digits = append([]byte{'-'}, digits...)
	}
	return string(digits)
}
