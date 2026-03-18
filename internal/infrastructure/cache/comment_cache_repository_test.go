package cache

import (
	"context"
	"testing"

	"github.com/artorias501/blog-service/internal/domain/entity"
	"github.com/artorias501/blog-service/internal/domain/repository"
	"github.com/artorias501/blog-service/internal/domain/valueobject"
)

func TestCommentCacheRepository_GetSet(t *testing.T) {
	client := setupTestRedis(t)
	defer client.Close()

	repo := NewCommentCacheRepository(client)
	ctx := context.Background()

	// Create test comment
	authorName, _ := valueobject.NewAuthorName("Test Author")
	authorEmail, _ := valueobject.NewAuthorEmail("test@example.com")
	content, _ := valueobject.NewContent("Test comment content")
	comment := entity.NewComment(authorName, authorEmail, content)

	t.Run("set and get a comment", func(t *testing.T) {
		// Clean up before test
		client.Del(ctx, CommentKey(comment.ID()))

		err := repo.Set(ctx, comment)
		if err != nil {
			t.Fatalf("Expected no error on Set, got: %v", err)
		}

		cached, err := repo.Get(ctx, comment.ID())
		if err != nil {
			t.Fatalf("Expected no error on Get, got: %v", err)
		}
		if cached == nil {
			t.Fatal("Expected cached comment, got nil")
		}
		if cached.ID().String() != comment.ID().String() {
			t.Errorf("Expected comment ID %s, got: %s", comment.ID().String(), cached.ID().String())
		}
		if cached.Content().String() != comment.Content().String() {
			t.Errorf("Expected content %s, got: %s", comment.Content().String(), cached.Content().String())
		}

		// Clean up
		client.Del(ctx, CommentKey(comment.ID()))
	})

	t.Run("get non-existent comment returns nil", func(t *testing.T) {
		fakeID := valueobject.GenerateCommentID()
		cached, err := repo.Get(ctx, fakeID)
		if err != nil {
			t.Fatalf("Expected no error for cache miss, got: %v", err)
		}
		if cached != nil {
			t.Error("Expected nil for cache miss, got non-nil comment")
		}
	})
}

func TestCommentCacheRepository_Delete(t *testing.T) {
	client := setupTestRedis(t)
	defer client.Close()

	repo := NewCommentCacheRepository(client)
	ctx := context.Background()

	authorName, _ := valueobject.NewAuthorName("Test Author")
	authorEmail, _ := valueobject.NewAuthorEmail("test@example.com")
	content, _ := valueobject.NewContent("Test comment content")
	comment := entity.NewComment(authorName, authorEmail, content)

	// Set up the comment first
	repo.Set(ctx, comment)

	t.Run("delete a comment", func(t *testing.T) {
		err := repo.Delete(ctx, comment.ID())
		if err != nil {
			t.Fatalf("Expected no error on Delete, got: %v", err)
		}

		// Verify deletion
		cached, err := repo.Get(ctx, comment.ID())
		if err != nil {
			t.Fatalf("Expected no error on Get after delete, got: %v", err)
		}
		if cached != nil {
			t.Error("Expected nil after delete, got non-nil comment")
		}
	})
}

func TestCommentCacheRepository_ListByPostOperations(t *testing.T) {
	client := setupTestRedis(t)
	defer client.Close()

	repo := NewCommentCacheRepository(client)
	ctx := context.Background()

	postID := valueobject.GeneratePostID()
	params := repository.ListParams{Page: 1, PageSize: 10}

	// Create test comments
	authorName, _ := valueobject.NewAuthorName("Test Author")
	authorEmail, _ := valueobject.NewAuthorEmail("test@example.com")
	content, _ := valueobject.NewContent("Test comment content")
	comment := entity.NewComment(authorName, authorEmail, content)

	result := &repository.CommentListResult{
		Total:     1,
		Page:      1,
		PageSize:  10,
		TotalPage: 1,
		Items:     []*entity.Comment{comment},
	}

	t.Run("set and get comments by post ID", func(t *testing.T) {
		// Clean up
		client.Del(ctx, CommentByPostKey(postID, params.Page, params.PageSize))

		err := repo.SetListByPostID(ctx, postID, params, result)
		if err != nil {
			t.Fatalf("Expected no error on SetListByPostID, got: %v", err)
		}

		cached, err := repo.GetListByPostID(ctx, postID, params)
		if err != nil {
			t.Fatalf("Expected no error on GetListByPostID, got: %v", err)
		}
		if cached == nil {
			t.Fatal("Expected cached list, got nil")
		}
		if cached.Total != result.Total {
			t.Errorf("Expected total %d, got: %d", result.Total, cached.Total)
		}

		// Clean up
		client.Del(ctx, CommentByPostKey(postID, params.Page, params.PageSize))
	})

	t.Run("delete comments by post ID", func(t *testing.T) {
		// Set up
		repo.SetListByPostID(ctx, postID, params, result)

		err := repo.DeleteListByPostID(ctx, postID)
		if err != nil {
			t.Fatalf("Expected no error on DeleteListByPostID, got: %v", err)
		}

		// Verify deletion
		cached, err := repo.GetListByPostID(ctx, postID, params)
		if err != nil {
			t.Fatalf("Expected no error on GetListByPostID after delete, got: %v", err)
		}
		if cached != nil {
			t.Error("Expected nil after delete, got non-nil list")
		}
	})
}

func TestCommentCacheRepository_CountOperations(t *testing.T) {
	client := setupTestRedis(t)
	defer client.Close()

	repo := NewCommentCacheRepository(client)
	ctx := context.Background()

	postID := valueobject.GeneratePostID()

	t.Run("set and get comment count by post ID", func(t *testing.T) {
		// Clean up
		client.Del(ctx, CommentCountByPostKey(postID))

		err := repo.SetCountByPostID(ctx, postID, 42)
		if err != nil {
			t.Fatalf("Expected no error on SetCountByPostID, got: %v", err)
		}

		count, err := repo.GetCountByPostID(ctx, postID)
		if err != nil {
			t.Fatalf("Expected no error on GetCountByPostID, got: %v", err)
		}
		if count != 42 {
			t.Errorf("Expected count 42, got: %d", count)
		}

		// Clean up
		client.Del(ctx, CommentCountByPostKey(postID))
	})

	t.Run("get count for non-existent post returns -1", func(t *testing.T) {
		fakePostID := valueobject.GeneratePostID()
		count, err := repo.GetCountByPostID(ctx, fakePostID)
		if err != nil {
			t.Fatalf("Expected no error for cache miss, got: %v", err)
		}
		if count != -1 {
			t.Errorf("Expected -1 for cache miss, got: %d", count)
		}
	})
}

func TestCommentCacheRepository_InvalidateOperations(t *testing.T) {
	client := setupTestRedis(t)
	defer client.Close()

	repo := NewCommentCacheRepository(client)
	ctx := context.Background()

	postID := valueobject.GeneratePostID()

	authorName, _ := valueobject.NewAuthorName("Test Author")
	authorEmail, _ := valueobject.NewAuthorEmail("test@example.com")
	content, _ := valueobject.NewContent("Test comment content")
	comment := entity.NewComment(authorName, authorEmail, content)

	// Set up comment
	repo.Set(ctx, comment)

	t.Run("invalidate comment removes all related keys", func(t *testing.T) {
		err := repo.InvalidateComment(ctx, comment.ID())
		if err != nil {
			t.Fatalf("Expected no error on InvalidateComment, got: %v", err)
		}

		// Verify comment is deleted
		cached, err := repo.Get(ctx, comment.ID())
		if err != nil {
			t.Fatalf("Expected no error on Get after invalidate, got: %v", err)
		}
		if cached != nil {
			t.Error("Expected nil after invalidate, got non-nil comment")
		}
	})

	t.Run("invalidate by post ID removes all related keys", func(t *testing.T) {
		// Set up comment list and count
		params := repository.ListParams{Page: 1, PageSize: 10}
		result := &repository.CommentListResult{
			Total:     1,
			Page:      1,
			PageSize:  10,
			TotalPage: 1,
			Items:     []*entity.Comment{comment},
		}
		repo.SetListByPostID(ctx, postID, params, result)
		repo.SetCountByPostID(ctx, postID, 1)

		err := repo.InvalidateByPostID(ctx, postID)
		if err != nil {
			t.Fatalf("Expected no error on InvalidateByPostID, got: %v", err)
		}

		// Verify deletion
		cached, err := repo.GetListByPostID(ctx, postID, params)
		if err != nil {
			t.Fatalf("Expected no error on GetListByPostID after invalidate, got: %v", err)
		}
		if cached != nil {
			t.Error("Expected nil after invalidate, got non-nil list")
		}

		count, err := repo.GetCountByPostID(ctx, postID)
		if err != nil {
			t.Fatalf("Expected no error on GetCountByPostID after invalidate, got: %v", err)
		}
		if count != -1 {
			t.Errorf("Expected -1 after invalidate, got: %d", count)
		}
	})
}
