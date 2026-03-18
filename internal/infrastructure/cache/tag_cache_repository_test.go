package cache

import (
	"context"
	"testing"

	"github.com/artorias501/blog-service/internal/domain/entity"
	"github.com/artorias501/blog-service/internal/domain/repository"
	"github.com/artorias501/blog-service/internal/domain/valueobject"
)

func TestTagCacheRepository_GetSet(t *testing.T) {
	client := setupTestRedis(t)
	defer client.Close()

	repo := NewTagCacheRepository(client)
	ctx := context.Background()

	// Create test tag
	tagName, _ := valueobject.NewTagName("golang")
	tag := entity.NewTag(tagName)

	t.Run("set and get a tag", func(t *testing.T) {
		// Clean up before test
		client.Del(ctx, TagKey(tag.ID()))

		err := repo.Set(ctx, tag)
		if err != nil {
			t.Fatalf("Expected no error on Set, got: %v", err)
		}

		cached, err := repo.Get(ctx, tag.ID())
		if err != nil {
			t.Fatalf("Expected no error on Get, got: %v", err)
		}
		if cached == nil {
			t.Fatal("Expected cached tag, got nil")
		}
		if cached.ID().String() != tag.ID().String() {
			t.Errorf("Expected tag ID %s, got: %s", tag.ID().String(), cached.ID().String())
		}
		if cached.Name().String() != tag.Name().String() {
			t.Errorf("Expected tag name %s, got: %s", tag.Name().String(), cached.Name().String())
		}

		// Clean up
		client.Del(ctx, TagKey(tag.ID()))
	})

	t.Run("get non-existent tag returns nil", func(t *testing.T) {
		fakeID := valueobject.GenerateTagID()
		cached, err := repo.Get(ctx, fakeID)
		if err != nil {
			t.Fatalf("Expected no error for cache miss, got: %v", err)
		}
		if cached != nil {
			t.Error("Expected nil for cache miss, got non-nil tag")
		}
	})
}

func TestTagCacheRepository_Delete(t *testing.T) {
	client := setupTestRedis(t)
	defer client.Close()

	repo := NewTagCacheRepository(client)
	ctx := context.Background()

	tagName, _ := valueobject.NewTagName("golang")
	tag := entity.NewTag(tagName)

	// Set up the tag first
	repo.Set(ctx, tag)

	t.Run("delete a tag", func(t *testing.T) {
		err := repo.Delete(ctx, tag.ID())
		if err != nil {
			t.Fatalf("Expected no error on Delete, got: %v", err)
		}

		// Verify deletion
		cached, err := repo.Get(ctx, tag.ID())
		if err != nil {
			t.Fatalf("Expected no error on Get after delete, got: %v", err)
		}
		if cached != nil {
			t.Error("Expected nil after delete, got non-nil tag")
		}
	})
}

func TestTagCacheRepository_ByNameOperations(t *testing.T) {
	client := setupTestRedis(t)
	defer client.Close()

	repo := NewTagCacheRepository(client)
	ctx := context.Background()

	tagName, _ := valueobject.NewTagName("golang")
	tag := entity.NewTag(tagName)

	t.Run("set and get tag by name", func(t *testing.T) {
		// Clean up
		client.Del(ctx, TagByNameKey(tagName))

		err := repo.SetByName(ctx, tagName, tag)
		if err != nil {
			t.Fatalf("Expected no error on SetByName, got: %v", err)
		}

		cached, err := repo.GetByName(ctx, tagName)
		if err != nil {
			t.Fatalf("Expected no error on GetByName, got: %v", err)
		}
		if cached == nil {
			t.Fatal("Expected cached tag, got nil")
		}
		if cached.Name().String() != tag.Name().String() {
			t.Errorf("Expected tag name %s, got: %s", tag.Name().String(), cached.Name().String())
		}

		// Clean up
		client.Del(ctx, TagByNameKey(tagName))
	})

	t.Run("get non-existent tag by name returns nil", func(t *testing.T) {
		fakeName, _ := valueobject.NewTagName("nonexistent")
		cached, err := repo.GetByName(ctx, fakeName)
		if err != nil {
			t.Fatalf("Expected no error for cache miss, got: %v", err)
		}
		if cached != nil {
			t.Error("Expected nil for cache miss, got non-nil tag")
		}
	})
}

func TestTagCacheRepository_ListOperations(t *testing.T) {
	client := setupTestRedis(t)
	defer client.Close()

	repo := NewTagCacheRepository(client)
	ctx := context.Background()

	// Create test tags
	var tags []*entity.Tag
	for i := 0; i < 3; i++ {
		tagName, _ := valueobject.NewTagName("tag-" + string(rune('A'+i)))
		tags = append(tags, entity.NewTag(tagName))
	}

	params := repository.ListParams{
		Page:     1,
		PageSize: 10,
	}

	result := &repository.TagListResult{
		Total:     3,
		Page:      1,
		PageSize:  10,
		TotalPage: 1,
		Items:     tags,
	}

	t.Run("set and get tag list", func(t *testing.T) {
		// Clean up
		client.Del(ctx, TagListKey(params.Page, params.PageSize))

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
		client.Del(ctx, TagListKey(params.Page, params.PageSize))
	})

	t.Run("delete tag list", func(t *testing.T) {
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

func TestTagCacheRepository_ByPostOperations(t *testing.T) {
	client := setupTestRedis(t)
	defer client.Close()

	repo := NewTagCacheRepository(client)
	ctx := context.Background()

	postID := valueobject.GeneratePostID()

	// Create test tags
	tagName1, _ := valueobject.NewTagName("golang")
	tagName2, _ := valueobject.NewTagName("backend")
	tags := []*entity.Tag{
		entity.NewTag(tagName1),
		entity.NewTag(tagName2),
	}

	t.Run("set and get tags by post ID", func(t *testing.T) {
		// Clean up
		client.Del(ctx, TagByPostKey(postID))

		err := repo.SetByPostID(ctx, postID, tags)
		if err != nil {
			t.Fatalf("Expected no error on SetByPostID, got: %v", err)
		}

		cached, err := repo.GetByPostID(ctx, postID)
		if err != nil {
			t.Fatalf("Expected no error on GetByPostID, got: %v", err)
		}
		if cached == nil {
			t.Fatal("Expected cached tags, got nil")
		}
		if len(cached) != len(tags) {
			t.Errorf("Expected %d tags, got: %d", len(tags), len(cached))
		}

		// Clean up
		client.Del(ctx, TagByPostKey(postID))
	})
}

func TestTagCacheRepository_InvalidateTag(t *testing.T) {
	client := setupTestRedis(t)
	defer client.Close()

	repo := NewTagCacheRepository(client)
	ctx := context.Background()

	tagName, _ := valueobject.NewTagName("golang")
	tag := entity.NewTag(tagName)

	// Set up tag
	repo.Set(ctx, tag)
	repo.SetByName(ctx, tagName, tag)

	t.Run("invalidate tag removes all related keys", func(t *testing.T) {
		err := repo.InvalidateTag(ctx, tag.ID())
		if err != nil {
			t.Fatalf("Expected no error on InvalidateTag, got: %v", err)
		}

		// Verify tag is deleted
		cached, err := repo.Get(ctx, tag.ID())
		if err != nil {
			t.Fatalf("Expected no error on Get after invalidate, got: %v", err)
		}
		if cached != nil {
			t.Error("Expected nil after invalidate, got non-nil tag")
		}
	})
}
