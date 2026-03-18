package main

import (
	"context"
	"testing"
	"time"

	"github.com/artorias501/blog-service/pkg/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestInitializeApplication tests application initialization.
func TestInitializeApplication(t *testing.T) {
	// Create test config with file-based SQLite for proper migration
	cfg := &config.Config{
		Server: config.ServerConfig{
			Port:        "8080",
			Environment: "test",
		},
		Database: config.DatabaseConfig{
			Path: t.TempDir() + "/test.db",
		},
		Redis: config.RedisConfig{
			Addr:     "localhost:6379",
			Password: "",
			DB:       0,
		},
		Auth: config.AuthConfig{
			AdminToken: "test-token",
		},
		CORS: config.CORSConfig{
			AllowedOrigins: []string{"*"},
			AllowedMethods: []string{"GET", "POST", "PUT", "DELETE"},
		},
		Log: config.LogConfig{
			Format: "text",
			Level:  "info",
		},
	}

	// Create test logger
	logger := initializeLogger(cfg)
	require.NotNil(t, logger)

	// Initialize application
	app, err := InitializeApplication(cfg, logger)
	require.NoError(t, err)
	require.NotNil(t, app)

	// Verify all components are initialized
	assert.NotNil(t, app.Config)
	assert.NotNil(t, app.Logger)
	assert.NotNil(t, app.DB)
	assert.NotNil(t, app.Router)
	assert.NotNil(t, app.Server)
	assert.NotNil(t, app.Health)
	assert.NotNil(t, app.Post)
	assert.NotNil(t, app.Tag)
	assert.NotNil(t, app.Comment)

	// Verify server address
	assert.Equal(t, ":8080", app.Server.Addr)
}

// TestInitializeApplication_FailFast tests that application fails fast on missing config.
func TestInitializeApplication_FailFast(t *testing.T) {
	// Create invalid config (empty database path)
	cfg := &config.Config{
		Server: config.ServerConfig{
			Port:        "8080",
			Environment: "test",
		},
		Database: config.DatabaseConfig{
			Path: "", // Empty path should fail validation
		},
	}

	// Validate should fail
	err := cfg.Validate()
	assert.Error(t, err)
}

// TestApplication_Shutdown tests graceful shutdown.
func TestApplication_Shutdown(t *testing.T) {
	cfg := &config.Config{
		Server: config.ServerConfig{
			Port:        "8080",
			Environment: "test",
		},
		Database: config.DatabaseConfig{
			Path: t.TempDir() + "/test.db",
		},
		Redis: config.RedisConfig{
			Addr:     "localhost:6379",
			Password: "",
			DB:       0,
		},
		Auth: config.AuthConfig{
			AdminToken: "test-token",
		},
		CORS: config.CORSConfig{
			AllowedOrigins: []string{"*"},
			AllowedMethods: []string{"GET", "POST", "PUT", "DELETE"},
		},
		Log: config.LogConfig{
			Format: "text",
			Level:  "info",
		},
	}

	logger := initializeLogger(cfg)
	app, err := InitializeApplication(cfg, logger)
	require.NoError(t, err)

	// Create shutdown context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Shutdown should complete without error
	err = app.Shutdown(ctx)
	assert.NoError(t, err)
}

// TestSetupRouter tests router setup.
func TestSetupRouter(t *testing.T) {
	cfg := &config.Config{
		Server: config.ServerConfig{
			Port:        "8080",
			Environment: "test",
		},
		Database: config.DatabaseConfig{
			Path: t.TempDir() + "/test.db",
		},
		Redis: config.RedisConfig{
			Addr:     "localhost:6379",
			Password: "",
			DB:       0,
		},
		Auth: config.AuthConfig{
			AdminToken: "test-token",
		},
		CORS: config.CORSConfig{
			AllowedOrigins: []string{"*"},
			AllowedMethods: []string{"GET", "POST", "PUT", "DELETE"},
		},
		Log: config.LogConfig{
			Format: "text",
			Level:  "info",
		},
	}

	logger := initializeLogger(cfg)
	app, err := InitializeApplication(cfg, logger)
	require.NoError(t, err)

	// Verify router is configured
	assert.NotNil(t, app.Router)

	// Verify routes are registered
	routes := app.Router.Routes()
	assert.NotEmpty(t, routes)

	// Check for expected routes
	var hasHealth, hasPosts, hasTags, hasComments bool
	for _, route := range routes {
		switch {
		case route.Path == "/health":
			hasHealth = true
		case route.Path == "/api/v1/posts":
			hasPosts = true
		case route.Path == "/api/v1/tags":
			hasTags = true
		case route.Path == "/api/v1/comments":
			hasComments = true
		}
	}

	assert.True(t, hasHealth, "Health route should be registered")
	assert.True(t, hasPosts, "Posts route should be registered")
	assert.True(t, hasTags, "Tags route should be registered")
	assert.True(t, hasComments, "Comments route should be registered")
}

// TestInitializeLogger tests logger initialization.
func TestInitializeLogger(t *testing.T) {
	tests := []struct {
		name    string
		env     string
		level   string
		wantErr bool
	}{
		{
			name:  "development with info level",
			env:   "development",
			level: "info",
		},
		{
			name:  "production with error level",
			env:   "production",
			level: "error",
		},
		{
			name:  "test with debug level",
			env:   "test",
			level: "debug",
		},
		{
			name:  "default level fallback",
			env:   "development",
			level: "unknown", // Should default to info
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &config.Config{
				Server: config.ServerConfig{
					Environment: tt.env,
				},
				Log: config.LogConfig{
					Format: "text",
					Level:  tt.level,
				},
			}

			logger := initializeLogger(cfg)
			assert.NotNil(t, logger)
		})
	}
}
