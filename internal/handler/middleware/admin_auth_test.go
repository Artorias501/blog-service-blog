package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/artorias501/blog-service/internal/handler/middleware"
	"github.com/artorias501/blog-service/pkg/config"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func TestAdminAuth_ValidToken(t *testing.T) {
	router := gin.New()
	cfg := &config.Config{
		Auth: config.AuthConfig{
			AdminToken: "test-token",
		},
	}
	router.Use(middleware.AdminAuth(cfg))
	router.GET("/protected", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Bearer test-token")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestAdminAuth_MissingHeader(t *testing.T) {
	router := gin.New()
	cfg := &config.Config{
		Auth: config.AuthConfig{
			AdminToken: "test-token",
		},
	}
	router.Use(middleware.AdminAuth(cfg))
	router.GET("/protected", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestAdminAuth_InvalidFormat(t *testing.T) {
	router := gin.New()
	cfg := &config.Config{
		Auth: config.AuthConfig{
			AdminToken: "test-token",
		},
	}
	router.Use(middleware.AdminAuth(cfg))
	router.GET("/protected", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	tests := []struct {
		name   string
		header string
	}{
		{"no bearer prefix", "test-token"},
		{"wrong prefix", "Basic test-token"},
		{"bearer lowercase", "bearer test-token"},
		{"empty after bearer", "Bearer "},
		{"extra spaces", "Bearer  test-token"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "/protected", nil)
			req.Header.Set("Authorization", tt.header)
			router.ServeHTTP(w, req)

			assert.Equal(t, http.StatusUnauthorized, w.Code)
		})
	}
}

func TestAdminAuth_WrongToken(t *testing.T) {
	router := gin.New()
	cfg := &config.Config{
		Auth: config.AuthConfig{
			AdminToken: "test-token",
		},
	}
	router.Use(middleware.AdminAuth(cfg))
	router.GET("/protected", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Bearer wrong-token")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestAdminAuth_EmptyToken(t *testing.T) {
	router := gin.New()
	cfg := &config.Config{
		Auth: config.AuthConfig{
			AdminToken: "",
		},
	}
	router.Use(middleware.AdminAuth(cfg))
	router.GET("/protected", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Bearer any-token")
	router.ServeHTTP(w, req)

	// When no admin token is configured, should deny all access
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestAdminAuth_SetsContext(t *testing.T) {
	router := gin.New()
	cfg := &config.Config{
		Auth: config.AuthConfig{
			AdminToken: "test-token",
		},
	}
	router.Use(middleware.AdminAuth(cfg))
	router.GET("/protected", func(c *gin.Context) {
		// Check that admin flag is set in context
		isAdmin, exists := c.Get("is_admin")
		assert.True(t, exists)
		assert.True(t, isAdmin.(bool))
		c.Status(http.StatusOK)
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Bearer test-token")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}
