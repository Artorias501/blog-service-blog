package handler

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/artorias501/blog-service/pkg/config"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func init() {
	gin.SetMode(gin.TestMode)
}

// TestHealthHandler_Check tests the health check endpoint.
func TestHealthHandler_Check(t *testing.T) {
	// Create in-memory database for testing
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err, "Failed to create test database")

	// Create test config
	cfg := &config.Config{
		Server: config.ServerConfig{
			Port:        "8080",
			Environment: "test",
		},
	}

	tests := []struct {
		name       string
		setupRedis bool
		wantStatus int
	}{
		{
			name:       "healthy without redis",
			setupRedis: false,
			wantStatus: http.StatusOK,
		},
		{
			name:       "healthy with redis",
			setupRedis: true,
			wantStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var redisClient *redis.Client
			if tt.setupRedis {
				// Create Redis client (will fail if Redis is not available, but that's ok for this test)
				redisClient = redis.NewClient(&redis.Options{
					Addr: "localhost:6379",
				})
				// Try to ping, but don't fail if Redis is not available
				ctx := context.Background()
				if err := redisClient.Ping(ctx).Err(); err != nil {
					// Redis not available, set to nil
					redisClient = nil
				}
			}

			handler := NewHealthHandler(cfg, db, redisClient)

			// Create test router
			router := gin.New()
			router.GET("/health", handler.Check)

			// Create request
			req := httptest.NewRequest(http.MethodGet, "/health", nil)
			w := httptest.NewRecorder()

			// Execute
			router.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, tt.wantStatus, w.Code)

			// Verify response body contains expected fields
			assert.Contains(t, w.Body.String(), "status")
			assert.Contains(t, w.Body.String(), "service")
			assert.Contains(t, w.Body.String(), "version")
			assert.Contains(t, w.Body.String(), "checks")
		})
	}
}

// TestHealthHandler_Check_Returns200 tests that health check returns 200.
func TestHealthHandler_Check_Returns200(t *testing.T) {
	// Create in-memory database for testing
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err, "Failed to create test database")

	cfg := &config.Config{
		Server: config.ServerConfig{
			Port:        "8080",
			Environment: "test",
		},
	}

	handler := NewHealthHandler(cfg, db, nil)

	// Create test router
	router := gin.New()
	router.GET("/health", handler.Check)

	// Create request
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()

	// Execute
	router.ServeHTTP(w, req)

	// Assert - must return 200
	assert.Equal(t, http.StatusOK, w.Code)
}

// TestHealthHandler_Check_DatabaseCheck tests database health check.
func TestHealthHandler_Check_DatabaseCheck(t *testing.T) {
	// Create in-memory database for testing
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err, "Failed to create test database")

	cfg := &config.Config{
		Server: config.ServerConfig{
			Port:        "8080",
			Environment: "test",
		},
	}

	handler := NewHealthHandler(cfg, db, nil)

	// Create test router
	router := gin.New()
	router.GET("/health", handler.Check)

	// Create request
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()

	// Execute
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)
	// Database check should be present
	assert.Contains(t, w.Body.String(), "database")
}

// TestNewHealthHandler tests the constructor.
func TestNewHealthHandler(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	cfg := &config.Config{
		Server: config.ServerConfig{
			Port:        "8080",
			Environment: "test",
		},
	}

	handler := NewHealthHandler(cfg, db, nil)

	assert.NotNil(t, handler)
	assert.Equal(t, cfg, handler.config)
	assert.Equal(t, db, handler.db)
	assert.Nil(t, handler.redis)
}
