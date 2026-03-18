// Package config provides configuration loading and management for the blog service.
package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"
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

// yamlConfig is used for YAML unmarshaling with snake_case field names.
type yamlConfig struct {
	Server   yamlServerConfig   `yaml:"server"`
	Database yamlDatabaseConfig `yaml:"database"`
	Redis    yamlRedisConfig    `yaml:"redis"`
	Auth     yamlAuthConfig     `yaml:"auth"`
	CORS     yamlCORSConfig     `yaml:"cors"`
	Log      yamlLogConfig      `yaml:"log"`
}

type yamlServerConfig struct {
	Port        string `yaml:"port"`
	Environment string `yaml:"environment"`
}

type yamlDatabaseConfig struct {
	Path string `yaml:"path"`
}

type yamlRedisConfig struct {
	Addr     string       `yaml:"addr"`
	Password string       `yaml:"password"`
	DB       int          `yaml:"db"`
	TTL      yamlRedisTTL `yaml:"ttl"`
}

type yamlRedisTTL struct {
	Post         string `yaml:"post"`
	PostList     string `yaml:"post_list"`
	Tag          string `yaml:"tag"`
	Comment      string `yaml:"comment"`
	CommentCount string `yaml:"comment_count"`
}

type yamlAuthConfig struct {
	AdminToken string `yaml:"admin_token"`
}

type yamlCORSConfig struct {
	AllowedOrigins   []string `yaml:"allowed_origins"`
	AllowedMethods   []string `yaml:"allowed_methods"`
	AllowedHeaders   []string `yaml:"allowed_headers"`
	AllowCredentials bool     `yaml:"allow_credentials"`
	MaxAge           int      `yaml:"max_age"`
}

type yamlLogConfig struct {
	Format string `yaml:"format"`
	Level  string `yaml:"level"`
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

// Load reads configuration from YAML file (if exists) and environment variables.
// If configPath is empty, it looks for "config.yaml" in the current directory.
// Environment variables override YAML values.
func Load(configPath ...string) (*Config, error) {
	// Determine config file path
	path := "config.yaml"
	if len(configPath) > 0 && configPath[0] != "" {
		path = configPath[0]
	}

	cfg := getDefaultConfig()

	// Try to load from YAML file
	data, err := os.ReadFile(path)
	if err != nil {
		// If file doesn't exist, use defaults and continue with env overrides
		if !os.IsNotExist(err) {
			return nil, fmt.Errorf("failed to read config file: %w", err)
		}
	} else {
		// Parse YAML file
		var yc yamlConfig
		if err := yaml.Unmarshal(data, &yc); err != nil {
			return nil, fmt.Errorf("failed to parse config file: %w", err)
		}
		// Apply YAML values to config
		applyYAMLConfig(cfg, &yc)
	}

	// Apply environment variable overrides
	applyEnvOverrides(cfg)

	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

// getDefaultConfig returns a Config with default values.
func getDefaultConfig() *Config {
	return &Config{
		Server: ServerConfig{
			Port:        "8080",
			Environment: "development",
		},
		Database: DatabaseConfig{
			Path: "blog.db",
		},
		Redis: RedisConfig{
			Addr:     "localhost:6379",
			Password: "",
			DB:       0,
		},
		Auth: AuthConfig{
			AdminToken: "artorias501",
		},
		CORS: CORSConfig{
			AllowedOrigins:   []string{"*"},
			AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
			AllowedHeaders:   []string{"Content-Type", "Authorization"},
			AllowCredentials: false,
			MaxAge:           86400,
		},
		Log: LogConfig{
			Format: "text",
			Level:  "info",
		},
	}
}

// applyYAMLConfig applies YAML configuration values to Config.
func applyYAMLConfig(cfg *Config, yc *yamlConfig) {
	if yc.Server.Port != "" {
		cfg.Server.Port = yc.Server.Port
	}
	if yc.Server.Environment != "" {
		cfg.Server.Environment = yc.Server.Environment
	}
	if yc.Database.Path != "" {
		cfg.Database.Path = yc.Database.Path
	}
	if yc.Redis.Addr != "" {
		cfg.Redis.Addr = yc.Redis.Addr
	}
	cfg.Redis.Password = yc.Redis.Password
	cfg.Redis.DB = yc.Redis.DB
	if yc.Auth.AdminToken != "" {
		cfg.Auth.AdminToken = yc.Auth.AdminToken
	}
	if len(yc.CORS.AllowedOrigins) > 0 {
		cfg.CORS.AllowedOrigins = yc.CORS.AllowedOrigins
	}
	if len(yc.CORS.AllowedMethods) > 0 {
		cfg.CORS.AllowedMethods = yc.CORS.AllowedMethods
	}
	if len(yc.CORS.AllowedHeaders) > 0 {
		cfg.CORS.AllowedHeaders = yc.CORS.AllowedHeaders
	}
	cfg.CORS.AllowCredentials = yc.CORS.AllowCredentials
	cfg.CORS.MaxAge = yc.CORS.MaxAge
	if yc.Log.Format != "" {
		cfg.Log.Format = yc.Log.Format
	}
	if yc.Log.Level != "" {
		cfg.Log.Level = yc.Log.Level
	}
}

// applyEnvOverrides applies environment variable overrides to Config.
func applyEnvOverrides(cfg *Config) {
	if v := os.Getenv("SERVER_PORT"); v != "" {
		cfg.Server.Port = v
	}
	if v := os.Getenv("SERVER_ENVIRONMENT"); v != "" {
		cfg.Server.Environment = v
	}
	if v := os.Getenv("DATABASE_PATH"); v != "" {
		cfg.Database.Path = v
	}
	if v := os.Getenv("REDIS_ADDR"); v != "" {
		cfg.Redis.Addr = v
	}
	if v := os.Getenv("REDIS_PASSWORD"); v != "" {
		cfg.Redis.Password = v
	}
	if v := os.Getenv("REDIS_DB"); v != "" {
		if intVal, err := strconv.Atoi(v); err == nil {
			cfg.Redis.DB = intVal
		}
	}
	if v := os.Getenv("ADMIN_TOKEN"); v != "" {
		cfg.Auth.AdminToken = v
	}
	if v := os.Getenv("CORS_ALLOWED_ORIGINS"); v != "" {
		cfg.CORS.AllowedOrigins = parseCommaSeparated(v)
	}
	if v := os.Getenv("CORS_ALLOWED_METHODS"); v != "" {
		cfg.CORS.AllowedMethods = parseCommaSeparated(v)
	}
	if v := os.Getenv("CORS_ALLOWED_HEADERS"); v != "" {
		cfg.CORS.AllowedHeaders = parseCommaSeparated(v)
	}
	if v := os.Getenv("CORS_ALLOW_CREDENTIALS"); v != "" {
		if boolVal, err := strconv.ParseBool(v); err == nil {
			cfg.CORS.AllowCredentials = boolVal
		}
	}
	if v := os.Getenv("CORS_MAX_AGE"); v != "" {
		if intVal, err := strconv.Atoi(v); err == nil {
			cfg.CORS.MaxAge = intVal
		}
	}
	if v := os.Getenv("LOG_FORMAT"); v != "" {
		cfg.Log.Format = v
	}
	if v := os.Getenv("LOG_LEVEL"); v != "" {
		cfg.Log.Level = v
	}
}

// parseCommaSeparated parses a comma-separated string into a slice.
func parseCommaSeparated(value string) []string {
	parts := strings.Split(value, ",")
	result := make([]string, 0, len(parts))
	for _, part := range parts {
		if trimmed := strings.TrimSpace(part); trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
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
