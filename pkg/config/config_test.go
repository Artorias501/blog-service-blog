package config_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/artorias501/blog-service/pkg/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// createTempYAMLConfig creates a temporary YAML config file with the given content.
// Returns the file path and a cleanup function.
func createTempYAMLConfig(t *testing.T, content string) (string, func()) {
	t.Helper()

	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")

	err := os.WriteFile(configPath, []byte(content), 0644)
	require.NoError(t, err, "Failed to create temp config file")

	return configPath, func() {
		// t.TempDir() handles cleanup automatically
	}
}

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
	assert.Equal(t, "artorias501", cfg.Auth.AdminToken)
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

// ============================================================================
// YAML Loading Tests
// ============================================================================

func TestLoad_YAMLFile(t *testing.T) {
	os.Clearenv()

	yamlContent := `
server:
  port: "9090"
  environment: "staging"
database:
  path: "/data/test.db"
redis:
  addr: "redis-server:6379"
  password: "redis-pass"
  db: 3
auth:
  admin_token: "yaml-token"
cors:
  allowed_origins:
    - "http://localhost:8080"
    - "http://example.com"
  allowed_methods:
    - "GET"
    - "POST"
  allowed_headers:
    - "Content-Type"
  allow_credentials: true
  max_age: 3600
log:
  format: "json"
  level: "debug"
`
	configPath, cleanup := createTempYAMLConfig(t, yamlContent)
	defer cleanup()

	cfg, err := config.Load(configPath)
	require.NoError(t, err)

	// Verify YAML values were loaded
	assert.Equal(t, "9090", cfg.Server.Port)
	assert.Equal(t, "staging", cfg.Server.Environment)
	assert.Equal(t, "/data/test.db", cfg.Database.Path)
	assert.Equal(t, "redis-server:6379", cfg.Redis.Addr)
	assert.Equal(t, "redis-pass", cfg.Redis.Password)
	assert.Equal(t, 3, cfg.Redis.DB)
	assert.Equal(t, "yaml-token", cfg.Auth.AdminToken)
	assert.Len(t, cfg.CORS.AllowedOrigins, 2)
	assert.Contains(t, cfg.CORS.AllowedOrigins, "http://localhost:8080")
	assert.Contains(t, cfg.CORS.AllowedOrigins, "http://example.com")
	assert.Len(t, cfg.CORS.AllowedMethods, 2)
	assert.True(t, cfg.CORS.AllowCredentials)
	assert.Equal(t, 3600, cfg.CORS.MaxAge)
	assert.Equal(t, "json", cfg.Log.Format)
	assert.Equal(t, "debug", cfg.Log.Level)
}

func TestLoad_YAMLFile_PartialConfig(t *testing.T) {
	os.Clearenv()

	// YAML with only some fields set - should use defaults for missing fields
	yamlContent := `
server:
  port: "7070"
database:
  path: "custom.db"
`
	configPath, cleanup := createTempYAMLConfig(t, yamlContent)
	defer cleanup()

	cfg, err := config.Load(configPath)
	require.NoError(t, err)

	// Verify specified values
	assert.Equal(t, "7070", cfg.Server.Port)
	assert.Equal(t, "custom.db", cfg.Database.Path)

	// Verify defaults for unspecified fields
	assert.Equal(t, "development", cfg.Server.Environment)
	assert.Equal(t, "localhost:6379", cfg.Redis.Addr)
	assert.Equal(t, "artorias501", cfg.Auth.AdminToken)
}

func TestLoad_YAMLFile_InvalidYAML(t *testing.T) {
	os.Clearenv()

	yamlContent := `
server:
  port: "8080"
  invalid yaml content here
    broken indentation
`
	configPath, cleanup := createTempYAMLConfig(t, yamlContent)
	defer cleanup()

	_, err := config.Load(configPath)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to parse config file")
}

func TestLoad_YAMLFile_NonExistent(t *testing.T) {
	os.Clearenv()

	// Loading from non-existent path should use defaults
	cfg, err := config.Load("/non/existent/path/config.yaml")
	require.NoError(t, err)

	// Should use defaults
	assert.Equal(t, "8080", cfg.Server.Port)
	assert.Equal(t, "development", cfg.Server.Environment)
	assert.Equal(t, "blog.db", cfg.Database.Path)
}

func TestLoad_EmptyConfigPath(t *testing.T) {
	os.Clearenv()

	// Empty config path should use default "config.yaml"
	cfg, err := config.Load("")
	require.NoError(t, err)

	// Should use defaults (since config.yaml may or may not exist)
	assert.NotNil(t, cfg)
}

// ============================================================================
// Environment Variable Override Tests
// ============================================================================

func TestLoad_EnvOverridesYAML(t *testing.T) {
	os.Clearenv()

	yamlContent := `
server:
  port: "9090"
  environment: "staging"
database:
  path: "/data/yaml.db"
redis:
  addr: "yaml-redis:6379"
  password: "yaml-pass"
  db: 5
auth:
  admin_token: "yaml-token"
`
	configPath, cleanup := createTempYAMLConfig(t, yamlContent)
	defer cleanup()

	// Set environment variables that should override YAML values
	os.Setenv("SERVER_PORT", "3000")
	os.Setenv("SERVER_ENVIRONMENT", "production")
	os.Setenv("DATABASE_PATH", "/data/env.db")
	os.Setenv("REDIS_ADDR", "env-redis:6379")
	os.Setenv("REDIS_PASSWORD", "env-pass")
	os.Setenv("REDIS_DB", "10")
	os.Setenv("ADMIN_TOKEN", "env-token")

	defer os.Clearenv()

	cfg, err := config.Load(configPath)
	require.NoError(t, err)

	// Environment variables should override YAML values
	assert.Equal(t, "3000", cfg.Server.Port)
	assert.Equal(t, "production", cfg.Server.Environment)
	assert.Equal(t, "/data/env.db", cfg.Database.Path)
	assert.Equal(t, "env-redis:6379", cfg.Redis.Addr)
	assert.Equal(t, "env-pass", cfg.Redis.Password)
	assert.Equal(t, 10, cfg.Redis.DB)
	assert.Equal(t, "env-token", cfg.Auth.AdminToken)
}

func TestLoad_EnvOverridesYAML_CORS(t *testing.T) {
	os.Clearenv()

	yamlContent := `
cors:
  allowed_origins:
    - "http://yaml-origin.com"
  allowed_methods:
    - "GET"
  allow_credentials: false
  max_age: 1000
`
	configPath, cleanup := createTempYAMLConfig(t, yamlContent)
	defer cleanup()

	// Set CORS environment variables
	os.Setenv("CORS_ALLOWED_ORIGINS", "http://env-origin.com,http://another.com")
	os.Setenv("CORS_ALLOWED_METHODS", "POST,PUT")
	os.Setenv("CORS_ALLOW_CREDENTIALS", "true")
	os.Setenv("CORS_MAX_AGE", "5000")

	defer os.Clearenv()

	cfg, err := config.Load(configPath)
	require.NoError(t, err)

	assert.Len(t, cfg.CORS.AllowedOrigins, 2)
	assert.Contains(t, cfg.CORS.AllowedOrigins, "http://env-origin.com")
	assert.Contains(t, cfg.CORS.AllowedOrigins, "http://another.com")
	assert.Len(t, cfg.CORS.AllowedMethods, 2)
	assert.Contains(t, cfg.CORS.AllowedMethods, "POST")
	assert.Contains(t, cfg.CORS.AllowedMethods, "PUT")
	assert.True(t, cfg.CORS.AllowCredentials)
	assert.Equal(t, 5000, cfg.CORS.MaxAge)
}

func TestLoad_EnvOverridesYAML_Log(t *testing.T) {
	os.Clearenv()

	yamlContent := `
log:
  format: "json"
  level: "debug"
`
	configPath, cleanup := createTempYAMLConfig(t, yamlContent)
	defer cleanup()

	// Set log environment variables
	os.Setenv("LOG_FORMAT", "text")
	os.Setenv("LOG_LEVEL", "error")

	defer os.Clearenv()

	cfg, err := config.Load(configPath)
	require.NoError(t, err)

	assert.Equal(t, "text", cfg.Log.Format)
	assert.Equal(t, "error", cfg.Log.Level)
}

func TestLoad_EnvPartialOverride(t *testing.T) {
	os.Clearenv()

	yamlContent := `
server:
  port: "9090"
  environment: "staging"
database:
  path: "/data/yaml.db"
`
	configPath, cleanup := createTempYAMLConfig(t, yamlContent)
	defer cleanup()

	// Only override some values
	os.Setenv("SERVER_PORT", "3000")
	// SERVER_ENVIRONMENT not set - should use YAML value

	defer os.Clearenv()

	cfg, err := config.Load(configPath)
	require.NoError(t, err)

	// Overridden by env
	assert.Equal(t, "3000", cfg.Server.Port)
	// From YAML
	assert.Equal(t, "staging", cfg.Server.Environment)
	assert.Equal(t, "/data/yaml.db", cfg.Database.Path)
}

// ============================================================================
// Auto-creation Tests
// ============================================================================

func TestLoad_NoConfigFile_UsesDefaults(t *testing.T) {
	os.Clearenv()

	// Use a non-existent path
	cfg, err := config.Load("/tmp/nonexistent-config-12345.yaml")
	require.NoError(t, err)

	// Should use all defaults
	assert.Equal(t, "8080", cfg.Server.Port)
	assert.Equal(t, "development", cfg.Server.Environment)
	assert.Equal(t, "blog.db", cfg.Database.Path)
	assert.Equal(t, "localhost:6379", cfg.Redis.Addr)
	assert.Equal(t, "artorias501", cfg.Auth.AdminToken)
}

func TestLoad_DefaultConfigPath(t *testing.T) {
	os.Clearenv()

	// Load without specifying path - should use "config.yaml" if exists, or defaults
	cfg, err := config.Load()
	require.NoError(t, err)
	require.NotNil(t, cfg)
}

// ============================================================================
// Redis TTL Configuration Tests
// ============================================================================

func TestLoad_YAMLFile_RedisTTL(t *testing.T) {
	os.Clearenv()

	yamlContent := `
server:
  port: "9090"
database:
  path: "/data/test.db"
redis:
  addr: "redis-server:6379"
  ttl:
    post: "45m"
    post_list: "10m"
    tag: "2h"
    comment: "30m"
    comment_count: "15m"
`
	configPath, cleanup := createTempYAMLConfig(t, yamlContent)
	defer cleanup()

	cfg, err := config.Load(configPath)
	require.NoError(t, err)

	// Verify TTL values were loaded correctly
	assert.Equal(t, 45*time.Minute, cfg.Redis.TTL.Post)
	assert.Equal(t, 10*time.Minute, cfg.Redis.TTL.PostList)
	assert.Equal(t, 2*time.Hour, cfg.Redis.TTL.Tag)
	assert.Equal(t, 30*time.Minute, cfg.Redis.TTL.Comment)
	assert.Equal(t, 15*time.Minute, cfg.Redis.TTL.CommentCount)
}

func TestLoad_DefaultRedisTTL(t *testing.T) {
	os.Clearenv()

	cfg, err := config.Load()
	require.NoError(t, err)

	// Verify default TTL values
	assert.Equal(t, 30*time.Minute, cfg.Redis.TTL.Post)
	assert.Equal(t, 5*time.Minute, cfg.Redis.TTL.PostList)
	assert.Equal(t, 60*time.Minute, cfg.Redis.TTL.Tag)
	assert.Equal(t, 15*time.Minute, cfg.Redis.TTL.Comment)
	assert.Equal(t, 5*time.Minute, cfg.Redis.TTL.CommentCount)
}
