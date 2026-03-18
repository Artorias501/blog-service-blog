package logger_test

import (
	"bytes"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/artorias501/blog-service/pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func TestNew_JSONFormat(t *testing.T) {
	// Create a buffer to capture output
	var buf bytes.Buffer
	log := logger.New("production", &buf)

	log.Info("test message")

	// Verify JSON output
	var entry map[string]interface{}
	err := json.Unmarshal(buf.Bytes(), &entry)
	require.NoError(t, err)

	assert.Equal(t, "INFO", entry["level"])
	assert.Equal(t, "test message", entry["msg"])
	assert.Contains(t, entry, "time")
}

func TestNew_TextFormat(t *testing.T) {
	var buf bytes.Buffer
	log := logger.New("development", &buf)

	log.Info("test message")

	// Verify text output (not JSON)
	output := buf.String()
	assert.Contains(t, output, "INFO")
	assert.Contains(t, output, "test message")
	assert.NotContains(t, output, "{") // Should not be JSON
}

func TestNew_DefaultToText(t *testing.T) {
	var buf bytes.Buffer
	log := logger.New("invalid", &buf)

	log.Info("test message")

	// Should default to text format
	output := buf.String()
	assert.Contains(t, output, "INFO")
	assert.Contains(t, output, "test message")
}

func TestNew_WithFields(t *testing.T) {
	var buf bytes.Buffer
	log := logger.New("production", &buf)

	log.Info("test message", "key", "value", "count", 42)

	var entry map[string]interface{}
	err := json.Unmarshal(buf.Bytes(), &entry)
	require.NoError(t, err)

	assert.Equal(t, "value", entry["key"])
	assert.Equal(t, float64(42), entry["count"])
}

func TestNew_DebugLevel(t *testing.T) {
	var buf bytes.Buffer
	log := logger.NewWithLevel("production", &buf, slog.LevelDebug)

	log.Debug("debug message")

	var entry map[string]interface{}
	err := json.Unmarshal(buf.Bytes(), &entry)
	require.NoError(t, err)

	assert.Equal(t, "DEBUG", entry["level"])
}

func TestNew_WarnLevel(t *testing.T) {
	var buf bytes.Buffer
	log := logger.New("production", &buf)

	log.Warn("warning message")

	var entry map[string]interface{}
	err := json.Unmarshal(buf.Bytes(), &entry)
	require.NoError(t, err)

	assert.Equal(t, "WARN", entry["level"])
}

func TestNew_ErrorLevel(t *testing.T) {
	var buf bytes.Buffer
	log := logger.New("production", &buf)

	log.Error("error message")

	var entry map[string]interface{}
	err := json.Unmarshal(buf.Bytes(), &entry)
	require.NoError(t, err)

	assert.Equal(t, "ERROR", entry["level"])
}

func TestNew_DefaultOutput(t *testing.T) {
	// Test that default output is os.Stdout when nil is passed
	log := logger.New("development", nil)
	assert.NotNil(t, log)

	// This should not panic
	log.Info("test")
}

func TestGetLogger(t *testing.T) {
	// Test that GetLogger returns a valid logger
	log := logger.GetLogger()
	assert.NotNil(t, log)

	// Should not panic
	log.Info("test")
}

func TestSetDefault(t *testing.T) {
	var buf bytes.Buffer
	customLogger := logger.New("production", &buf)

	logger.SetDefault(customLogger)

	// Get the default logger and use it
	log := logger.GetLogger()
	log.Info("set default test")

	// Verify output
	var entry map[string]interface{}
	err := json.Unmarshal(buf.Bytes(), &entry)
	require.NoError(t, err)
	assert.Equal(t, "set default test", entry["msg"])
}

func TestWithRequestID(t *testing.T) {
	var buf bytes.Buffer
	log := logger.New("production", &buf)

	logger.WithRequestID(log, "req-123").Info("test message")

	var entry map[string]interface{}
	err := json.Unmarshal(buf.Bytes(), &entry)
	require.NoError(t, err)

	assert.Equal(t, "req-123", entry["request_id"])
}

func TestSync(t *testing.T) {
	// Test that Sync does not panic
	log := logger.New("production", os.Stdout)
	err := logger.Sync(log)
	assert.NoError(t, err)
}

// Integration test with Gin context
func TestFromContext(t *testing.T) {
	router := gin.New()
	router.Use(func(c *gin.Context) {
		log := logger.New("development", os.Stdout)
		c.Set("logger", log)
		c.Next()
	})

	router.GET("/test", func(c *gin.Context) {
		log, exists := logger.FromContext(c)
		assert.True(t, exists)
		assert.NotNil(t, log)
		c.Status(http.StatusOK)
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

// Helper function to capture log output
func captureOutput(fn func()) string {
	r, w, _ := os.Pipe()
	defer r.Close()

	original := os.Stdout
	os.Stdout = w

	fn()

	os.Stdout = original
	w.Close()

	var buf bytes.Buffer
	io.Copy(&buf, r)
	return buf.String()
}

// Helper to check if string is valid JSON
func isValidJSON(s string) bool {
	var js interface{}
	return json.Unmarshal([]byte(s), &js) == nil
}

// Test that text output is not JSON
func TestTextOutputNotJSON(t *testing.T) {
	var buf bytes.Buffer
	log := logger.New("development", &buf)

	log.Info("test")

	// Text output should not be valid JSON
	assert.False(t, isValidJSON(strings.TrimSpace(buf.String())))
}
