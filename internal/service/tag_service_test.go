package service

import (
	"context"
	"errors"
	"testing"

	"github.com/artorias501/blog-service/internal/domain/entity"
	"github.com/artorias501/blog-service/internal/domain/repository"
	"github.com/artorias501/blog-service/internal/domain/valueobject"
)

// Helper function to create a test tag
func createTestTag(t *testing.T) *entity.Tag {
	name, err := valueobject.NewTagName("test-tag")
	if err != nil {
		t.Fatalf("failed to create tag name: %v", err)
	}
	return entity.NewTag(name)
}

// Test CreateTag
func TestTagService_CreateTag(t *testing.T) {
	tests := []struct {
		name    string
		input   CreateTagInput
		wantErr bool
	}{
		{
			name: "valid tag",
			input: CreateTagInput{
				Name: "new-tag",
			},
			wantErr: false,
		},
		{
			name: "empty name",
			input: CreateTagInput{
				Name: "",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			tagRepo := newMockTagRepository()
			tagCache := newMockTagCacheRepository()

			svc := NewTagService(tagRepo, tagCache)

			tag, err := svc.CreateTag(ctx, tt.input)

			if tt.wantErr {
				if err == nil {
					t.Errorf("CreateTag() expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("CreateTag() unexpected error: %v", err)
				return
			}

			if tag == nil {
				t.Error("CreateTag() returned nil tag")
				return
			}

			if tag.Name().String() != tt.input.Name {
				t.Errorf("CreateTag() name = %v, want %v", tag.Name().String(), tt.input.Name)
			}
		})
	}
}

// Test GetTagByID with cache-aside pattern
func TestTagService_GetTagByID_CacheAside(t *testing.T) {
	t.Run("cache hit", func(t *testing.T) {
		ctx := context.Background()
		tagRepo := newMockTagRepository()
		tagCache := newMockTagCacheRepository()

		testTag := createTestTag(t)
		tagRepo.tags[testTag.ID().String()] = testTag
		tagCache.cache[testTag.ID().String()] = testTag

		svc := NewTagService(tagRepo, tagCache)

		tag, err := svc.GetTagByID(ctx, testTag.ID().String())
		if err != nil {
			t.Errorf("GetTagByID() unexpected error: %v", err)
			return
		}

		if tag == nil {
			t.Error("GetTagByID() returned nil tag")
			return
		}

		if tag.ID().String() != testTag.ID().String() {
			t.Errorf("GetTagByID() returned wrong tag")
		}
	})

	t.Run("cache miss - fetch from database", func(t *testing.T) {
		ctx := context.Background()
		tagRepo := newMockTagRepository()
		tagCache := newMockTagCacheRepository()

		testTag := createTestTag(t)
		tagRepo.tags[testTag.ID().String()] = testTag

		svc := NewTagService(tagRepo, tagCache)

		tag, err := svc.GetTagByID(ctx, testTag.ID().String())
		if err != nil {
			t.Errorf("GetTagByID() unexpected error: %v", err)
			return
		}

		if tag == nil {
			t.Error("GetTagByID() returned nil tag")
			return
		}

		// Verify tag was cached
		cachedTag := tagCache.cache[testTag.ID().String()]
		if cachedTag == nil {
			t.Error("GetTagByID() did not populate cache")
		}
	})
}

// Test UpdateTag with cache invalidation
func TestTagService_UpdateTag_CacheInvalidation(t *testing.T) {
	ctx := context.Background()
	tagRepo := newMockTagRepository()
	tagCache := newMockTagCacheRepository()

	testTag := createTestTag(t)
	tagRepo.tags[testTag.ID().String()] = testTag
	tagCache.cache[testTag.ID().String()] = testTag

	svc := NewTagService(tagRepo, tagCache)

	newName := "updated-tag"
	input := UpdateTagInput{
		Name: newName,
	}

	updatedTag, err := svc.UpdateTag(ctx, testTag.ID().String(), input)
	if err != nil {
		t.Errorf("UpdateTag() unexpected error: %v", err)
		return
	}

	if updatedTag.Name().String() != newName {
		t.Errorf("UpdateTag() name = %v, want %v", updatedTag.Name().String(), newName)
	}

	// Verify cache was invalidated
	if tagCache.cache[testTag.ID().String()] != nil {
		t.Error("UpdateTag() did not invalidate cache")
	}
}

// Test DeleteTag with cache invalidation
func TestTagService_DeleteTag_CacheInvalidation(t *testing.T) {
	ctx := context.Background()
	tagRepo := newMockTagRepository()
	tagCache := newMockTagCacheRepository()

	testTag := createTestTag(t)
	tagRepo.tags[testTag.ID().String()] = testTag
	tagCache.cache[testTag.ID().String()] = testTag

	svc := NewTagService(tagRepo, tagCache)

	err := svc.DeleteTag(ctx, testTag.ID().String())
	if err != nil {
		t.Errorf("DeleteTag() unexpected error: %v", err)
		return
	}

	// Verify tag was deleted from repository
	if tagRepo.tags[testTag.ID().String()] != nil {
		t.Error("DeleteTag() did not delete tag from repository")
	}

	// Verify cache was invalidated
	if tagCache.cache[testTag.ID().String()] != nil {
		t.Error("DeleteTag() did not invalidate cache")
	}
}

// Test ListTags with cache integration
func TestTagService_ListTags_CacheIntegration(t *testing.T) {
	t.Run("cache miss - fetch from database", func(t *testing.T) {
		ctx := context.Background()
		tagRepo := newMockTagRepository()
		tagCache := newMockTagCacheRepository()

		testTag := createTestTag(t)
		tagRepo.tags[testTag.ID().String()] = testTag

		svc := NewTagService(tagRepo, tagCache)

		result, err := svc.ListTags(ctx, ListTagsInput{Page: 1, PageSize: 10})
		if err != nil {
			t.Errorf("ListTags() unexpected error: %v", err)
			return
		}

		if result.Total < 1 {
			t.Errorf("ListTags() expected at least 1 item, got %d", result.Total)
		}

		// Verify result was cached
		if tagCache.listCache["tag:list:page:1:size:10"] == nil {
			t.Error("ListTags() did not populate cache")
		}
	})
}

// Test CreateTag invalidates list cache
func TestTagService_CreateTag_InvalidatesListCache(t *testing.T) {
	ctx := context.Background()
	tagRepo := newMockTagRepository()
	tagCache := newMockTagCacheRepository()

	// Pre-populate list cache
	tagCache.listCache["tag:list:page:1:size:10"] = &repository.TagListResult{Total: 0}

	svc := NewTagService(tagRepo, tagCache)

	_, err := svc.CreateTag(ctx, CreateTagInput{
		Name: "new-tag",
	})
	if err != nil {
		t.Errorf("CreateTag() unexpected error: %v", err)
		return
	}

	// Verify list cache was invalidated
	if len(tagCache.listCache) != 0 {
		t.Error("CreateTag() did not invalidate list cache")
	}
}

// Test error handling - repository errors wrapped with context
func TestTagService_RepositoryErrorHandling(t *testing.T) {
	t.Run("GetByID repository error", func(t *testing.T) {
		ctx := context.Background()
		tagRepo := newMockTagRepository()
		tagCache := newMockTagCacheRepository()

		expectedErr := errors.New("database connection error")
		tagRepo.getByIDFn = func(ctx context.Context, id valueobject.TagID) (*entity.Tag, error) {
			return nil, expectedErr
		}

		svc := NewTagService(tagRepo, tagCache)

		tagID := valueobject.GenerateTagID()
		_, err := svc.GetTagByID(ctx, tagID.String())
		if err == nil {
			t.Error("GetTagByID() expected error, got nil")
			return
		}

		// Verify error is wrapped with context
		if err.Error() == expectedErr.Error() {
			t.Error("GetTagByID() error should be wrapped with context")
		}
	})
}

// Test dependency injection - services use interfaces
func TestTagService_DependencyInjection(t *testing.T) {
	// This test verifies that TagService accepts interfaces, not implementations
	var _ repository.TagRepository = (*mockTagRepository)(nil)
	var _ repository.TagCacheRepository = (*mockTagCacheRepository)(nil)

	// If this compiles, the service uses interfaces for dependency injection
	svc := NewTagService(
		newMockTagRepository(),
		newMockTagCacheRepository(),
	)

	if svc == nil {
		t.Error("NewTagService returned nil")
	}
}
