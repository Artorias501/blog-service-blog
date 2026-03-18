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

func TestRecovery_PanicWithMessage(t *testing.T) {
	var buf bytes.Buffer
	log := logger.New("production", &buf)

	router := gin.New()
	router.Use(middleware.Recovery(log))
	router.GET("/test", func(c *gin.Context) {
		panic("something went wrong")
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)

	// Verify panic was logged
	var entry map[string]interface{}
	err := json.Unmarshal(buf.Bytes(), &entry)
	require.NoError(t, err)

	assert.Equal(t, "ERROR", entry["level"])
	assert.Contains(t, entry["msg"], "panic")
}

func TestRecovery_PanicWithError(t *testing.T) {
	var buf bytes.Buffer
	log := logger.New("production", &buf)

	router := gin.New()
	router.Use(middleware.Recovery(log))
	router.GET("/test", func(c *gin.Context) {
		panic(http.ErrHandlerTimeout)
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)

	// Verify panic was logged
	assert.NotEmpty(t, buf.String())
}

func TestRecovery_NoPanic(t *testing.T) {
	var buf bytes.Buffer
	log := logger.New("production", &buf)

	router := gin.New()
	router.Use(middleware.Recovery(log))
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestRecovery_LogsStackTrace(t *testing.T) {
	var buf bytes.Buffer
	log := logger.New("production", &buf)

	router := gin.New()
	router.Use(middleware.Recovery(log))
	router.GET("/test", func(c *gin.Context) {
		panic("test panic")
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	router.ServeHTTP(w, req)

	var entry map[string]interface{}
	err := json.Unmarshal(buf.Bytes(), &entry)
	require.NoError(t, err)

	// Should contain stack trace
	assert.Contains(t, entry, "stack")
}

func TestRecovery_LogsRequestDetails(t *testing.T) {
	var buf bytes.Buffer
	log := logger.New("production", &buf)

	router := gin.New()
	router.Use(middleware.Recovery(log))
	router.GET("/test", func(c *gin.Context) {
		panic("test panic")
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/test?param=value", nil)
	router.ServeHTTP(w, req)

	var entry map[string]interface{}
	err := json.Unmarshal(buf.Bytes(), &entry)
	require.NoError(t, err)

	// Should contain request details
	assert.Contains(t, entry, "method")
	assert.Contains(t, entry, "path")
}
