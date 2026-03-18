// Package config provides configuration loading and management for the blog service.
package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Config holds all application configuration.
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Redis    RedisConfig
	Auth     AuthConfig
	CORS     CORSConfig
	Log      LogConfig
}

// ServerConfig holds server-related configuration.
type ServerConfig struct {
	Port        string
	Environment string
}

// DatabaseConfig holds database-related configuration.
type DatabaseConfig struct {
	Path string
}

// RedisConfig holds Redis-related configuration.
type RedisConfig struct {
	Addr     string
	Password string
	DB       int
}

// AuthConfig holds authentication-related configuration.
type AuthConfig struct {
	AdminToken string
}

// CORSConfig holds CORS-related configuration.
type CORSConfig struct {
	AllowedOrigins   []string
	AllowedMethods   []string
	AllowedHeaders   []string
	AllowCredentials bool
	MaxAge           int
}

// LogConfig holds logging-related configuration.
type LogConfig struct {
	Format string // "json" or "text"
	Level  string // "debug", "info", "warn", "error"
}

// Load reads configuration from environment variables with defaults.
func Load() (*Config, error) {
	cfg := &Config{
		Server: ServerConfig{
			Port:        getEnvWithDefault("SERVER_PORT", "8080"),
			Environment: getEnvWithDefault("SERVER_ENVIRONMENT", "development"),
		},
		Database: DatabaseConfig{
			Path: getEnvWithDefault("DATABASE_PATH", "blog.db"),
		},
		Redis: RedisConfig{
			Addr:     getEnvWithDefault("REDIS_ADDR", "localhost:6379"),
			Password: getEnvWithDefault("REDIS_PASSWORD", ""),
			DB:       getEnvIntWithDefault("REDIS_DB", 0),
		},
		Auth: AuthConfig{
			AdminToken: getEnvWithDefault("ADMIN_TOKEN", "artorias501"),
		},
		CORS: CORSConfig{
			AllowedOrigins:   getEnvSliceWithDefault("CORS_ALLOWED_ORIGINS", []string{"*"}),
			AllowedMethods:   getEnvSliceWithDefault("CORS_ALLOWED_METHODS", []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}),
			AllowedHeaders:   getEnvSliceWithDefault("CORS_ALLOWED_HEADERS", []string{"Content-Type", "Authorization"}),
			AllowCredentials: getEnvBoolWithDefault("CORS_ALLOW_CREDENTIALS", false),
			MaxAge:           getEnvIntWithDefault("CORS_MAX_AGE", 86400),
		},
		Log: LogConfig{
			Format: getEnvWithDefault("LOG_FORMAT", "text"),
			Level:  getEnvWithDefault("LOG_LEVEL", "info"),
		},
	}

	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

// Validate validates the configuration.
func (c *Config) Validate() error {
	if c.Server.Port == "" {
		return fmt.Errorf("server port is required")
	}
	if c.Database.Path == "" {
		return fmt.Errorf("database path is required")
	}
	return nil
}

// IsProduction returns true if environment is production.
func (c *Config) IsProduction() bool {
	return strings.ToLower(c.Server.Environment) == "production"
}

// IsDevelopment returns true if environment is development.
func (c *Config) IsDevelopment() bool {
	return strings.ToLower(c.Server.Environment) == "development"
}

// GetAllowedOrigins returns the allowed origins slice.
func (c *CORSConfig) GetAllowedOrigins() []string {
	return c.AllowedOrigins
}

// GetAllowedMethods returns the allowed methods slice.
func (c *CORSConfig) GetAllowedMethods() []string {
	return c.AllowedMethods
}

// IsOriginAllowed checks if the given origin is allowed.
func (c *CORSConfig) IsOriginAllowed(origin string) bool {
	for _, o := range c.AllowedOrigins {
		if o == "*" || o == origin {
			return true
		}
	}
	return false
}

// Helper functions

func getEnvWithDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvIntWithDefault(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal
		}
	}
	return defaultValue
}

func getEnvBoolWithDefault(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolVal, err := strconv.ParseBool(value); err == nil {
			return boolVal
		}
	}
	return defaultValue
}

func getEnvSliceWithDefault(key string, defaultValue []string) []string {
	if value := os.Getenv(key); value != "" {
		parts := strings.Split(value, ",")
		result := make([]string, 0, len(parts))
		for _, part := range parts {
			if trimmed := strings.TrimSpace(part); trimmed != "" {
				result = append(result, trimmed)
			}
		}
		if len(result) > 0 {
			return result
		}
	}
	return defaultValue
}
