package cache

import (
	"context"
	"testing"
	"time"
)

func TestNewRedisClient(t *testing.T) {
	t.Run("connect to Redis with valid configuration", func(t *testing.T) {
		cfg := DefaultConfig()
		client, err := NewRedisClient(cfg)
		if err != nil {
			t.Fatalf("Expected no error, got: %v", err)
		}
		defer client.Close()

		// Test ping
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		pong, err := client.Ping(ctx).Result()
		if err != nil {
			t.Fatalf("Expected ping to succeed, got error: %v", err)
		}
		if pong != "PONG" {
			t.Errorf("Expected PONG, got: %s", pong)
		}
	})

	t.Run("connection fails gracefully when Redis is unavailable", func(t *testing.T) {
		cfg := Config{
			Addr: "localhost:16379", // Non-existent Redis
		}
		client, err := NewRedisClient(cfg)
		if err != nil {
			// Error is acceptable during creation
			if client != nil {
				client.Close()
			}
			return
		}
		defer client.Close()

		// If client was created, ping should fail
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		_, err = client.Ping(ctx).Result()
		if err == nil {
			t.Error("Expected error when connecting to non-existent Redis, got nil")
		}
	})

	t.Run("custom TTL configuration", func(t *testing.T) {
		cfg := Config{
			Addr:            "localhost:6379",
			PostTTL:         10 * time.Minute,
			PostListTTL:     5 * time.Minute,
			TagTTL:          15 * time.Minute,
			CommentTTL:      10 * time.Minute,
			CommentCountTTL: 5 * time.Minute,
		}
		client, err := NewRedisClient(cfg)
		if err != nil {
			t.Fatalf("Expected no error, got: %v", err)
		}
		defer client.Close()

		if client.Config().PostTTL != 10*time.Minute {
			t.Errorf("Expected PostTTL 10m, got: %v", client.Config().PostTTL)
		}
		if client.Config().PostListTTL != 5*time.Minute {
			t.Errorf("Expected PostListTTL 5m, got: %v", client.Config().PostListTTL)
		}
	})
}

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()

	if cfg.Addr != "localhost:6379" {
		t.Errorf("Expected default addr localhost:6379, got: %s", cfg.Addr)
	}
	if cfg.PostTTL != 30*time.Minute {
		t.Errorf("Expected default PostTTL 30m, got: %v", cfg.PostTTL)
	}
	if cfg.PostListTTL != 5*time.Minute {
		t.Errorf("Expected default PostListTTL 5m, got: %v", cfg.PostListTTL)
	}
	if cfg.TagTTL != 60*time.Minute {
		t.Errorf("Expected default TagTTL 60m, got: %v", cfg.TagTTL)
	}
	if cfg.CommentTTL != 15*time.Minute {
		t.Errorf("Expected default CommentTTL 15m, got: %v", cfg.CommentTTL)
	}
	if cfg.CommentCountTTL != 5*time.Minute {
		t.Errorf("Expected default CommentCountTTL 5m, got: %v", cfg.CommentCountTTL)
	}
}
