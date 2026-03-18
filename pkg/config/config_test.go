package config_test

import (
	"os"
	"testing"

	"github.com/artorias501/blog-service/pkg/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoad_Defaults(t *testing.T) {
	// Clear any existing env vars
	os.Clearenv()

	cfg, err := config.Load()
	require.NoError(t, err)

	// Verify defaults
	assert.Equal(t, "8080", cfg.Server.Port)
	assert.Equal(t, "development", cfg.Server.Environment)
	assert.Equal(t, "blog.db", cfg.Database.Path)
	assert.Equal(t, "localhost:6379", cfg.Redis.Addr)
	assert.Equal(t, "", cfg.Redis.Password)
	assert.Equal(t, 0, cfg.Redis.DB)
	assert.Equal(t, "", cfg.Auth.AdminToken)
}

func TestLoad_EnvironmentOverrides(t *testing.T) {
	os.Clearenv()

	// Set environment variables
	os.Setenv("SERVER_PORT", "3000")
	os.Setenv("SERVER_ENVIRONMENT", "production")
	os.Setenv("DATABASE_PATH", "/data/blog.db")
	os.Setenv("REDIS_ADDR", "redis:6379")
	os.Setenv("REDIS_PASSWORD", "secret")
	os.Setenv("REDIS_DB", "2")
	os.Setenv("ADMIN_TOKEN", "my-secret-token")
	os.Setenv("CORS_ALLOWED_ORIGINS", "http://localhost:3000,http://example.com")
	os.Setenv("CORS_ALLOWED_METHODS", "GET,POST,PUT,DELETE")

	defer os.Clearenv()

	cfg, err := config.Load()
	require.NoError(t, err)

	assert.Equal(t, "3000", cfg.Server.Port)
	assert.Equal(t, "production", cfg.Server.Environment)
	assert.Equal(t, "/data/blog.db", cfg.Database.Path)
	assert.Equal(t, "redis:6379", cfg.Redis.Addr)
	assert.Equal(t, "secret", cfg.Redis.Password)
	assert.Equal(t, 2, cfg.Redis.DB)
	assert.Equal(t, "my-secret-token", cfg.Auth.AdminToken)
	assert.Len(t, cfg.CORS.AllowedOrigins, 2)
	assert.Len(t, cfg.CORS.AllowedMethods, 4)
}

func TestLoad_InvalidPort(t *testing.T) {
	os.Clearenv()
	os.Setenv("SERVER_PORT", "invalid")
	defer os.Clearenv()

	// Port is stored as string, validation happens at server startup
	cfg, err := config.Load()
	require.NoError(t, err)
	assert.Equal(t, "invalid", cfg.Server.Port)
}

func TestLoad_InvalidRedisDB(t *testing.T) {
	os.Clearenv()
	os.Setenv("REDIS_DB", "invalid")
	defer os.Clearenv()

	// Invalid Redis DB value should use default
	cfg, err := config.Load()
	require.NoError(t, err)
	assert.Equal(t, 0, cfg.Redis.DB)
}

func TestConfig_Validate(t *testing.T) {
	tests := []struct {
		name    string
		modify  func(*config.Config)
		wantErr bool
	}{
		{
			name:    "valid config",
			modify:  func(c *config.Config) {},
			wantErr: false,
		},
		{
			name: "empty port",
			modify: func(c *config.Config) {
				c.Server.Port = ""
			},
			wantErr: true,
		},
		{
			name: "empty database path",
			modify: func(c *config.Config) {
				c.Database.Path = ""
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Clearenv()
			cfg, err := config.Load()
			require.NoError(t, err)

			tt.modify(cfg)
			err = cfg.Validate()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestConfig_IsProduction(t *testing.T) {
	os.Clearenv()

	cfg, err := config.Load()
	require.NoError(t, err)
	assert.False(t, cfg.IsProduction())

	os.Setenv("SERVER_ENVIRONMENT", "production")
	cfg, err = config.Load()
	require.NoError(t, err)
	assert.True(t, cfg.IsProduction())
}

func TestConfig_IsDevelopment(t *testing.T) {
	os.Clearenv()

	cfg, err := config.Load()
	require.NoError(t, err)
	assert.True(t, cfg.IsDevelopment())

	os.Setenv("SERVER_ENVIRONMENT", "production")
	cfg, err = config.Load()
	require.NoError(t, err)
	assert.False(t, cfg.IsDevelopment())
}

func TestCORSConfig_GetAllowedOrigins(t *testing.T) {
	os.Clearenv()
	os.Setenv("CORS_ALLOWED_ORIGINS", "http://localhost:3000,http://example.com,http://test.com")
	defer os.Clearenv()

	cfg, err := config.Load()
	require.NoError(t, err)

	origins := cfg.CORS.GetAllowedOrigins()
	assert.Len(t, origins, 3)
	assert.Contains(t, origins, "http://localhost:3000")
	assert.Contains(t, origins, "http://example.com")
	assert.Contains(t, origins, "http://test.com")
}

func TestCORSConfig_GetAllowedMethods(t *testing.T) {
	os.Clearenv()
	os.Setenv("CORS_ALLOWED_METHODS", "GET,POST,PUT,DELETE,PATCH")
	defer os.Clearenv()

	cfg, err := config.Load()
	require.NoError(t, err)

	methods := cfg.CORS.GetAllowedMethods()
	assert.Len(t, methods, 5)
	assert.Contains(t, methods, "GET")
	assert.Contains(t, methods, "POST")
	assert.Contains(t, methods, "PUT")
	assert.Contains(t, methods, "DELETE")
	assert.Contains(t, methods, "PATCH")
}

func TestCORSConfig_IsOriginAllowed(t *testing.T) {
	os.Clearenv()
	os.Setenv("CORS_ALLOWED_ORIGINS", "http://localhost:3000,http://example.com")
	defer os.Clearenv()

	cfg, err := config.Load()
	require.NoError(t, err)

	assert.True(t, cfg.CORS.IsOriginAllowed("http://localhost:3000"))
	assert.True(t, cfg.CORS.IsOriginAllowed("http://example.com"))
	assert.False(t, cfg.CORS.IsOriginAllowed("http://malicious.com"))
}

func TestCORSConfig_Wildcard(t *testing.T) {
	os.Clearenv()
	os.Setenv("CORS_ALLOWED_ORIGINS", "*")
	defer os.Clearenv()

	cfg, err := config.Load()
	require.NoError(t, err)

	assert.True(t, cfg.CORS.IsOriginAllowed("http://any-origin.com"))
	assert.True(t, cfg.CORS.IsOriginAllowed("http://localhost:3000"))
}
