package cache

import (
	"context"
	"testing"
	"time"

	"github.com/artorias501/blog-service/internal/domain/entity"
	"github.com/artorias501/blog-service/internal/domain/repository"
	"github.com/artorias501/blog-service/internal/domain/valueobject"
)

func TestPostCacheRepository_GetSet(t *testing.T) {
	client := setupTestRedis(t)
	defer client.Close()

	repo := NewPostCacheRepository(client)
	ctx := context.Background()

	// Create test post
	title, _ := valueobject.NewTitle("Test Post")
	content, _ := valueobject.NewContent("Test content")
	post := entity.NewPost(title, content)

	t.Run("set and get a post", func(t *testing.T) {
		// Clean up before test
		client.Del(ctx, PostKey(post.ID()))

		err := repo.Set(ctx, post)
		if err != nil {
			t.Fatalf("Expected no error on Set, got: %v", err)
		}

		cached, err := repo.Get(ctx, post.ID())
		if err != nil {
			t.Fatalf("Expected no error on Get, got: %v", err)
		}
		if cached == nil {
			t.Fatal("Expected cached post, got nil")
		}
		if cached.ID().String() != post.ID().String() {
			t.Errorf("Expected post ID %s, got: %s", post.ID().String(), cached.ID().String())
		}
		if cached.Title().String() != post.Title().String() {
			t.Errorf("Expected title %s, got: %s", post.Title().String(), cached.Title().String())
		}

		// Clean up
		client.Del(ctx, PostKey(post.ID()))
	})

	t.Run("get non-existent post returns nil", func(t *testing.T) {
		fakeID := valueobject.GeneratePostID()
		cached, err := repo.Get(ctx, fakeID)
		if err != nil {
			t.Fatalf("Expected no error for cache miss, got: %v", err)
		}
		if cached != nil {
			t.Error("Expected nil for cache miss, got non-nil post")
		}
	})
}

func TestPostCacheRepository_Delete(t *testing.T) {
	client := setupTestRedis(t)
	defer client.Close()

	repo := NewPostCacheRepository(client)
	ctx := context.Background()

	title, _ := valueobject.NewTitle("Test Post")
	content, _ := valueobject.NewContent("Test content")
	post := entity.NewPost(title, content)

	// Set up the post first
	repo.Set(ctx, post)

	t.Run("delete a post", func(t *testing.T) {
		err := repo.Delete(ctx, post.ID())
		if err != nil {
			t.Fatalf("Expected no error on Delete, got: %v", err)
		}

		// Verify deletion
		cached, err := repo.Get(ctx, post.ID())
		if err != nil {
			t.Fatalf("Expected no error on Get after delete, got: %v", err)
		}
		if cached != nil {
			t.Error("Expected nil after delete, got non-nil post")
		}
	})
}

func TestPostCacheRepository_ListOperations(t *testing.T) {
	client := setupTestRedis(t)
	defer client.Close()

	repo := NewPostCacheRepository(client)
	ctx := context.Background()

	// Create test posts
	var posts []*entity.Post
	for i := 0; i < 3; i++ {
		title, _ := valueobject.NewTitle("Test Post " + string(rune('A'+i)))
		content, _ := valueobject.NewContent("Test content")
		posts = append(posts, entity.NewPost(title, content))
	}

	params := repository.ListParams{
		Page:     1,
		PageSize: 10,
	}

	result := &repository.PostListResult{
		Total:     3,
		Page:      1,
		PageSize:  10,
		TotalPage: 1,
		Items:     posts,
	}

	t.Run("set and get post list", func(t *testing.T) {
		// Clean up
		client.Del(ctx, PostListKey(params.Page, params.PageSize))

		err := repo.SetList(ctx, params, result)
		if err != nil {
			t.Fatalf("Expected no error on SetList, got: %v", err)
		}

		cached, err := repo.GetList(ctx, params)
		if err != nil {
			t.Fatalf("Expected no error on GetList, got: %v", err)
		}
		if cached == nil {
			t.Fatal("Expected cached list, got nil")
		}
		if cached.Total != result.Total {
			t.Errorf("Expected total %d, got: %d", result.Total, cached.Total)
		}
		if len(cached.Items) != len(result.Items) {
			t.Errorf("Expected %d items, got: %d", len(result.Items), len(cached.Items))
		}

		// Clean up
		client.Del(ctx, PostListKey(params.Page, params.PageSize))
	})

	t.Run("cache miss on list returns nil", func(t *testing.T) {
		otherParams := repository.ListParams{Page: 99, PageSize: 10}
		cached, err := repo.GetList(ctx, otherParams)
		if err != nil {
			t.Fatalf("Expected no error for cache miss, got: %v", err)
		}
		if cached != nil {
			t.Error("Expected nil for cache miss, got non-nil list")
		}
	})

	t.Run("delete post list", func(t *testing.T) {
		// Set up list
		repo.SetList(ctx, params, result)

		err := repo.DeleteList(ctx)
		if err != nil {
			t.Fatalf("Expected no error on DeleteList, got: %v", err)
		}

		// Verify deletion
		cached, err := repo.GetList(ctx, params)
		if err != nil {
			t.Fatalf("Expected no error on GetList after delete, got: %v", err)
		}
		if cached != nil {
			t.Error("Expected nil after delete, got non-nil list")
		}
	})
}

func TestPostCacheRepository_ByTagOperations(t *testing.T) {
	client := setupTestRedis(t)
	defer client.Close()

	repo := NewPostCacheRepository(client)
	ctx := context.Background()

	tagID := valueobject.GenerateTagID()
	params := repository.ListParams{Page: 1, PageSize: 10}

	title, _ := valueobject.NewTitle("Test Post")
	content, _ := valueobject.NewContent("Test content")
	post := entity.NewPost(title, content)

	result := &repository.PostListResult{
		Total:     1,
		Page:      1,
		PageSize:  10,
		TotalPage: 1,
		Items:     []*entity.Post{post},
	}

	t.Run("set and get posts by tag ID", func(t *testing.T) {
		// Clean up
		client.Del(ctx, PostByTagKey(tagID, params.Page, params.PageSize))

		err := repo.SetByTagID(ctx, tagID, params, result)
		if err != nil {
			t.Fatalf("Expected no error on SetByTagID, got: %v", err)
		}

		cached, err := repo.GetByTagID(ctx, tagID, params)
		if err != nil {
			t.Fatalf("Expected no error on GetByTagID, got: %v", err)
		}
		if cached == nil {
			t.Fatal("Expected cached list, got nil")
		}
		if cached.Total != result.Total {
			t.Errorf("Expected total %d, got: %d", result.Total, cached.Total)
		}

		// Clean up
		client.Del(ctx, PostByTagKey(tagID, params.Page, params.PageSize))
	})
}

func TestPostCacheRepository_InvalidatePost(t *testing.T) {
	client := setupTestRedis(t)
	defer client.Close()

	repo := NewPostCacheRepository(client)
	ctx := context.Background()

	title, _ := valueobject.NewTitle("Test Post")
	content, _ := valueobject.NewContent("Test content")
	post := entity.NewPost(title, content)

	// Set up post and list
	repo.Set(ctx, post)
	params := repository.ListParams{Page: 1, PageSize: 10}
	result := &repository.PostListResult{
		Total:     1,
		Page:      1,
		PageSize:  10,
		TotalPage: 1,
		Items:     []*entity.Post{post},
	}
	repo.SetList(ctx, params, result)

	t.Run("invalidate post removes all related keys", func(t *testing.T) {
		err := repo.InvalidatePost(ctx, post.ID())
		if err != nil {
			t.Fatalf("Expected no error on InvalidatePost, got: %v", err)
		}

		// Verify post is deleted
		cached, err := repo.Get(ctx, post.ID())
		if err != nil {
			t.Fatalf("Expected no error on Get after invalidate, got: %v", err)
		}
		if cached != nil {
			t.Error("Expected nil after invalidate, got non-nil post")
		}
	})
}

// Helper function to set up test Redis connection
func setupTestRedis(t *testing.T) *RedisClient {
	cfg := DefaultConfig()
	client, err := NewRedisClient(cfg)
	if err != nil {
		t.Skipf("Redis not available: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if _, err := client.Ping(ctx).Result(); err != nil {
		client.Close()
		t.Skipf("Redis not available: %v", err)
	}

	return client
}
