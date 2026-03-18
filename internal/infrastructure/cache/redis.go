package cache

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

// Config holds the Redis connection configuration
type Config struct {
	// Addr is the Redis server address (host:port)
	Addr string

	// Password for Redis authentication (optional)
	Password string

	// DB is the Redis database number
	DB int

	// TTL configurations for different cache types
	PostTTL         time.Duration
	PostListTTL     time.Duration
	TagTTL          time.Duration
	CommentTTL      time.Duration
	CommentCountTTL time.Duration
}

// DefaultConfig returns a Config with sensible defaults
func DefaultConfig() Config {
	return Config{
		Addr:            "localhost:6379",
		Password:        "",
		DB:              0,
		PostTTL:         30 * time.Minute,
		PostListTTL:     5 * time.Minute,
		TagTTL:          60 * time.Minute,
		CommentTTL:      15 * time.Minute,
		CommentCountTTL: 5 * time.Minute,
	}
}

// RedisClient wraps the Redis client with configuration
type RedisClient struct {
	*redis.Client
	config Config
}

// NewRedisClient creates a new Redis client with the given configuration
// Returns error instead of panicking on connection failure
func NewRedisClient(cfg Config) (*RedisClient, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	return &RedisClient{
		Client: client,
		config: cfg,
	}, nil
}

// Close closes the Redis connection
func (c *RedisClient) Close() error {
	return c.Client.Close()
}

// Ping tests the connection to Redis
func (c *RedisClient) Ping(ctx context.Context) *redis.StatusCmd {
	return c.Client.Ping(ctx)
}

// Config returns the current configuration
func (c *RedisClient) Config() Config {
	return c.config
}
