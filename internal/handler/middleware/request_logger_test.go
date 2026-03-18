package middleware_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/artorias501/blog-service/internal/handler/middleware"
	"github.com/artorias501/blog-service/pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRequestLogger_LogsMethodAndPath(t *testing.T) {
	var buf bytes.Buffer
	log := logger.New("production", &buf)

	router := gin.New()
	router.Use(middleware.RequestLogger(log))
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	router.ServeHTTP(w, req)

	// Should have logged the request
	assert.NotEmpty(t, buf.String())

	var entry map[string]interface{}
	err := json.Unmarshal(buf.Bytes(), &entry)
	require.NoError(t, err)

	assert.Equal(t, "GET", entry["method"])
	assert.Equal(t, "/test", entry["path"])
}

func TestRequestLogger_LogsStatus(t *testing.T) {
	var buf bytes.Buffer
	log := logger.New("production", &buf)

	router := gin.New()
	router.Use(middleware.RequestLogger(log))
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	router.ServeHTTP(w, req)

	var entry map[string]interface{}
	err := json.Unmarshal(buf.Bytes(), &entry)
	require.NoError(t, err)

	assert.Equal(t, float64(404), entry["status"])
}

func TestRequestLogger_LogsLatency(t *testing.T) {
	var buf bytes.Buffer
	log := logger.New("production", &buf)

	router := gin.New()
	router.Use(middleware.RequestLogger(log))
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	router.ServeHTTP(w, req)

	var entry map[string]interface{}
	err := json.Unmarshal(buf.Bytes(), &entry)
	require.NoError(t, err)

	// Latency should be present (as a number in nanoseconds or similar)
	assert.Contains(t, entry, "latency")
}

func TestRequestLogger_LogsRequestID(t *testing.T) {
	var buf bytes.Buffer
	log := logger.New("production", &buf)

	router := gin.New()
	router.Use(middleware.RequestLogger(log))
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("X-Request-ID", "req-12345")
	router.ServeHTTP(w, req)

	var entry map[string]interface{}
	err := json.Unmarshal(buf.Bytes(), &entry)
	require.NoError(t, err)

	assert.Equal(t, "req-12345", entry["request_id"])
}

func TestRequestLogger_LogsClientIP(t *testing.T) {
	var buf bytes.Buffer
	log := logger.New("production", &buf)

	router := gin.New()
	router.Use(middleware.RequestLogger(log))
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	router.ServeHTTP(w, req)

	var entry map[string]interface{}
	err := json.Unmarshal(buf.Bytes(), &entry)
	require.NoError(t, err)

	assert.Contains(t, entry, "client_ip")
}

func TestRequestLogger_LogsUserAgent(t *testing.T) {
	var buf bytes.Buffer
	log := logger.New("production", &buf)

	router := gin.New()
	router.Use(middleware.RequestLogger(log))
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("User-Agent", "test-agent/1.0")
	router.ServeHTTP(w, req)

	var entry map[string]interface{}
	err := json.Unmarshal(buf.Bytes(), &entry)
	require.NoError(t, err)

	assert.Equal(t, "test-agent/1.0", entry["user_agent"])
}

func TestRequestLogger_LogsResponseSize(t *testing.T) {
	var buf bytes.Buffer
	log := logger.New("production", &buf)

	router := gin.New()
	router.Use(middleware.RequestLogger(log))
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	router.ServeHTTP(w, req)

	var entry map[string]interface{}
	err := json.Unmarshal(buf.Bytes(), &entry)
	require.NoError(t, err)

	assert.Contains(t, entry, "response_size")
}
