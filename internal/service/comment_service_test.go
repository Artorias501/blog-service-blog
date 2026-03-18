package service

import (
	"context"
	"errors"
	"testing"

	"github.com/artorias501/blog-service/internal/domain/entity"
	"github.com/artorias501/blog-service/internal/domain/repository"
	"github.com/artorias501/blog-service/internal/domain/valueobject"
)

// Helper function to create a test comment
func createTestComment(t *testing.T) *entity.Comment {
	authorName, err := valueobject.NewAuthorName("Test Author")
	if err != nil {
		t.Fatalf("failed to create author name: %v", err)
	}
	authorEmail, err := valueobject.NewAuthorEmail("test@example.com")
	if err != nil {
		t.Fatalf("failed to create author email: %v", err)
	}
	content, err := valueobject.NewContent("Test comment content")
	if err != nil {
		t.Fatalf("failed to create content: %v", err)
	}
	return entity.NewComment(authorName, authorEmail, content)
}

// Test CreateComment
func TestCommentService_CreateComment(t *testing.T) {
	tests := []struct {
		name    string
		input   CreateCommentInput
		setup   func(postRepo *mockPostRepository)
		wantErr bool
	}{
		{
			name: "valid comment",
			input: CreateCommentInput{
				PostID:      "",
				AuthorName:  "Test Author",
				AuthorEmail: "test@example.com",
				Content:     "Test comment",
			},
			setup: func(postRepo *mockPostRepository) {
				post := createTestPost(t)
				postRepo.posts[post.ID().String()] = post
			},
			wantErr: false,
		},
		{
			name: "empty author name",
			input: CreateCommentInput{
				PostID:      "",
				AuthorName:  "",
				AuthorEmail: "test@example.com",
				Content:     "Test comment",
			},
			setup: func(postRepo *mockPostRepository) {
				post := createTestPost(t)
				postRepo.posts[post.ID().String()] = post
			},
			wantErr: true,
		},
		{
			name: "invalid email",
			input: CreateCommentInput{
				PostID:      "",
				AuthorName:  "Test Author",
				AuthorEmail: "invalid-email",
				Content:     "Test comment",
			},
			setup: func(postRepo *mockPostRepository) {
				post := createTestPost(t)
				postRepo.posts[post.ID().String()] = post
			},
			wantErr: true,
		},
		{
			name: "post not found",
			input: CreateCommentInput{
				PostID:      "",
				AuthorName:  "Test Author",
				AuthorEmail: "test@example.com",
				Content:     "Test comment",
			},
			setup:   func(postRepo *mockPostRepository) {},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			postRepo := newMockPostRepository()
			postCache := newMockPostCacheRepository()
			commentRepo := newMockCommentRepository()
			commentCache := newMockCommentCacheRepository()

			// Setup post
			if tt.setup != nil {
				tt.setup(postRepo)
				if len(postRepo.posts) > 0 {
					for id := range postRepo.posts {
						tt.input.PostID = id
						break
					}
				}
			}

			svc := NewCommentService(commentRepo, commentCache, postRepo, postCache)

			comment, err := svc.CreateComment(ctx, tt.input)

			if tt.wantErr {
				if err == nil {
					t.Errorf("CreateComment() expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("CreateComment() unexpected error: %v", err)
				return
			}

			if comment == nil {
				t.Error("CreateComment() returned nil comment")
				return
			}

			if comment.AuthorName().String() != tt.input.AuthorName {
				t.Errorf("CreateComment() author name = %v, want %v", comment.AuthorName().String(), tt.input.AuthorName)
			}
		})
	}
}

// Test GetCommentByID with cache-aside pattern
func TestCommentService_GetCommentByID_CacheAside(t *testing.T) {
	t.Run("cache hit", func(t *testing.T) {
		ctx := context.Background()
		postRepo := newMockPostRepository()
		postCache := newMockPostCacheRepository()
		commentRepo := newMockCommentRepository()
		commentCache := newMockCommentCacheRepository()

		testComment := createTestComment(t)
		commentRepo.comments[testComment.ID().String()] = testComment
		commentCache.cache[testComment.ID().String()] = testComment

		svc := NewCommentService(commentRepo, commentCache, postRepo, postCache)

		comment, err := svc.GetCommentByID(ctx, testComment.ID().String())
		if err != nil {
			t.Errorf("GetCommentByID() unexpected error: %v", err)
			return
		}

		if comment == nil {
			t.Error("GetCommentByID() returned nil comment")
			return
		}

		if comment.ID().String() != testComment.ID().String() {
			t.Errorf("GetCommentByID() returned wrong comment")
		}
	})

	t.Run("cache miss - fetch from database", func(t *testing.T) {
		ctx := context.Background()
		postRepo := newMockPostRepository()
		postCache := newMockPostCacheRepository()
		commentRepo := newMockCommentRepository()
		commentCache := newMockCommentCacheRepository()

		testComment := createTestComment(t)
		commentRepo.comments[testComment.ID().String()] = testComment

		svc := NewCommentService(commentRepo, commentCache, postRepo, postCache)

		comment, err := svc.GetCommentByID(ctx, testComment.ID().String())
		if err != nil {
			t.Errorf("GetCommentByID() unexpected error: %v", err)
			return
		}

		if comment == nil {
			t.Error("GetCommentByID() returned nil comment")
			return
		}

		// Verify comment was cached
		cachedComment := commentCache.cache[testComment.ID().String()]
		if cachedComment == nil {
			t.Error("GetCommentByID() did not populate cache")
		}
	})
}

// Test UpdateComment with cache invalidation
func TestCommentService_UpdateComment_CacheInvalidation(t *testing.T) {
	ctx := context.Background()
	postRepo := newMockPostRepository()
	postCache := newMockPostCacheRepository()
	commentRepo := newMockCommentRepository()
	commentCache := newMockCommentCacheRepository()

	testComment := createTestComment(t)
	commentRepo.comments[testComment.ID().String()] = testComment
	commentCache.cache[testComment.ID().String()] = testComment

	svc := NewCommentService(commentRepo, commentCache, postRepo, postCache)

	newContent := "Updated comment content"
	input := UpdateCommentInput{
		Content: newContent,
	}

	updatedComment, err := svc.UpdateComment(ctx, testComment.ID().String(), input)
	if err != nil {
		t.Errorf("UpdateComment() unexpected error: %v", err)
		return
	}

	if updatedComment.Content().String() != newContent {
		t.Errorf("UpdateComment() content = %v, want %v", updatedComment.Content().String(), newContent)
	}

	// Verify cache was invalidated
	if commentCache.cache[testComment.ID().String()] != nil {
		t.Error("UpdateComment() did not invalidate cache")
	}
}

// Test DeleteComment with cache invalidation
func TestCommentService_DeleteComment_CacheInvalidation(t *testing.T) {
	ctx := context.Background()
	postRepo := newMockPostRepository()
	postCache := newMockPostCacheRepository()
	commentRepo := newMockCommentRepository()
	commentCache := newMockCommentCacheRepository()

	testComment := createTestComment(t)
	commentRepo.comments[testComment.ID().String()] = testComment
	commentCache.cache[testComment.ID().String()] = testComment

	svc := NewCommentService(commentRepo, commentCache, postRepo, postCache)

	err := svc.DeleteComment(ctx, testComment.ID().String())
	if err != nil {
		t.Errorf("DeleteComment() unexpected error: %v", err)
		return
	}

	// Verify comment was deleted from repository
	if commentRepo.comments[testComment.ID().String()] != nil {
		t.Error("DeleteComment() did not delete comment from repository")
	}

	// Verify cache was invalidated
	if commentCache.cache[testComment.ID().String()] != nil {
		t.Error("DeleteComment() did not invalidate cache")
	}
}

// Test ApproveComment status management
func TestCommentService_ApproveComment(t *testing.T) {
	ctx := context.Background()
	postRepo := newMockPostRepository()
	postCache := newMockPostCacheRepository()
	commentRepo := newMockCommentRepository()
	commentCache := newMockCommentCacheRepository()

	testComment := createTestComment(t)
	commentRepo.comments[testComment.ID().String()] = testComment
	commentCache.cache[testComment.ID().String()] = testComment

	svc := NewCommentService(commentRepo, commentCache, postRepo, postCache)

	err := svc.ApproveComment(ctx, testComment.ID().String())
	if err != nil {
		t.Errorf("ApproveComment() unexpected error: %v", err)
		return
	}

	// Verify cache was invalidated
	if commentCache.cache[testComment.ID().String()] != nil {
		t.Error("ApproveComment() did not invalidate cache")
	}
}

// Test RejectComment status management
func TestCommentService_RejectComment(t *testing.T) {
	ctx := context.Background()
	postRepo := newMockPostRepository()
	postCache := newMockPostCacheRepository()
	commentRepo := newMockCommentRepository()
	commentCache := newMockCommentCacheRepository()

	testComment := createTestComment(t)
	commentRepo.comments[testComment.ID().String()] = testComment
	commentCache.cache[testComment.ID().String()] = testComment

	svc := NewCommentService(commentRepo, commentCache, postRepo, postCache)

	err := svc.RejectComment(ctx, testComment.ID().String())
	if err != nil {
		t.Errorf("RejectComment() unexpected error: %v", err)
		return
	}

	// Verify cache was invalidated
	if commentCache.cache[testComment.ID().String()] != nil {
		t.Error("RejectComment() did not invalidate cache")
	}
}

// Test MarkCommentAsSpam status management
func TestCommentService_MarkCommentAsSpam(t *testing.T) {
	ctx := context.Background()
	postRepo := newMockPostRepository()
	postCache := newMockPostCacheRepository()
	commentRepo := newMockCommentRepository()
	commentCache := newMockCommentCacheRepository()

	testComment := createTestComment(t)
	commentRepo.comments[testComment.ID().String()] = testComment
	commentCache.cache[testComment.ID().String()] = testComment

	svc := NewCommentService(commentRepo, commentCache, postRepo, postCache)

	err := svc.MarkCommentAsSpam(ctx, testComment.ID().String())
	if err != nil {
		t.Errorf("MarkCommentAsSpam() unexpected error: %v", err)
		return
	}

	// Verify cache was invalidated
	if commentCache.cache[testComment.ID().String()] != nil {
		t.Error("MarkCommentAsSpam() did not invalidate cache")
	}
}

// Test ListCommentsByPost with cache integration
func TestCommentService_ListCommentsByPost_CacheIntegration(t *testing.T) {
	t.Run("cache miss - fetch from database", func(t *testing.T) {
		ctx := context.Background()
		postRepo := newMockPostRepository()
		postCache := newMockPostCacheRepository()
		commentRepo := newMockCommentRepository()
		commentCache := newMockCommentCacheRepository()

		// Create test post and comment
		testPost := createTestPost(t)
		postRepo.posts[testPost.ID().String()] = testPost
		testComment := createTestComment(t)
		commentRepo.comments[testComment.ID().String()] = testComment

		svc := NewCommentService(commentRepo, commentCache, postRepo, postCache)

		result, err := svc.ListCommentsByPost(ctx, ListCommentsByPostInput{
			PostID:   testPost.ID().String(),
			Page:     1,
			PageSize: 10,
		})
		if err != nil {
			t.Errorf("ListCommentsByPost() unexpected error: %v", err)
			return
		}

		if result.Total < 1 {
			t.Errorf("ListCommentsByPost() expected at least 1 item, got %d", result.Total)
		}
	})
}

// Test GetCommentCountByPostID with cache integration
func TestCommentService_GetCommentCountByPostID_CacheIntegration(t *testing.T) {
	t.Run("cache miss - fetch from database", func(t *testing.T) {
		ctx := context.Background()
		postRepo := newMockPostRepository()
		postCache := newMockPostCacheRepository()
		commentRepo := newMockCommentRepository()
		commentCache := newMockCommentCacheRepository()

		// Create test post and comments
		testPost := createTestPost(t)
		postRepo.posts[testPost.ID().String()] = testPost
		testComment1 := createTestComment(t)
		testComment2 := createTestComment(t)
		commentRepo.comments[testComment1.ID().String()] = testComment1
		commentRepo.comments[testComment2.ID().String()] = testComment2

		svc := NewCommentService(commentRepo, commentCache, postRepo, postCache)

		count, err := svc.GetCommentCountByPostID(ctx, testPost.ID().String())
		if err != nil {
			t.Errorf("GetCommentCountByPostID() unexpected error: %v", err)
			return
		}

		if count < 2 {
			t.Errorf("GetCommentCountByPostID() expected at least 2, got %d", count)
		}

		// Verify count was cached
		if commentCache.countCache[testPost.ID().String()] == 0 {
			t.Error("GetCommentCountByPostID() did not populate cache")
		}
	})
}

// Test error handling - repository errors wrapped with context
func TestCommentService_RepositoryErrorHandling(t *testing.T) {
	t.Run("GetByID repository error", func(t *testing.T) {
		ctx := context.Background()
		postRepo := newMockPostRepository()
		postCache := newMockPostCacheRepository()
		commentRepo := newMockCommentRepository()
		commentCache := newMockCommentCacheRepository()

		expectedErr := errors.New("database connection error")
		commentRepo.getByIDFn = func(ctx context.Context, id valueobject.CommentID) (*entity.Comment, error) {
			return nil, expectedErr
		}

		svc := NewCommentService(commentRepo, commentCache, postRepo, postCache)

		commentID := valueobject.GenerateCommentID()
		_, err := svc.GetCommentByID(ctx, commentID.String())
		if err == nil {
			t.Error("GetCommentByID() expected error, got nil")
			return
		}

		// Verify error is wrapped with context
		if err.Error() == expectedErr.Error() {
			t.Error("GetCommentByID() error should be wrapped with context")
		}
	})
}

// Test dependency injection - services use interfaces
func TestCommentService_DependencyInjection(t *testing.T) {
	// This test verifies that CommentService accepts interfaces, not implementations
	var _ repository.CommentRepository = (*mockCommentRepository)(nil)
	var _ repository.CommentCacheRepository = (*mockCommentCacheRepository)(nil)
	var _ repository.PostRepository = (*mockPostRepository)(nil)
	var _ repository.PostCacheRepository = (*mockPostCacheRepository)(nil)

	// If this compiles, the service uses interfaces for dependency injection
	svc := NewCommentService(
		newMockCommentRepository(),
		newMockCommentCacheRepository(),
		newMockPostRepository(),
		newMockPostCacheRepository(),
	)

	if svc == nil {
		t.Error("NewCommentService returned nil")
	}
}

// Test CreateComment invalidates post comment cache
func TestCommentService_CreateComment_InvalidatesPostCache(t *testing.T) {
	ctx := context.Background()
	postRepo := newMockPostRepository()
	postCache := newMockPostCacheRepository()
	commentRepo := newMockCommentRepository()
	commentCache := newMockCommentCacheRepository()

	// Create test post
	testPost := createTestPost(t)
	postRepo.posts[testPost.ID().String()] = testPost

	invalidateCalled := false
	commentCache.invalidateByPostIDFn = func(ctx context.Context, postID valueobject.PostID) error {
		invalidateCalled = true
		return nil
	}

	svc := NewCommentService(commentRepo, commentCache, postRepo, postCache)

	_, err := svc.CreateComment(ctx, CreateCommentInput{
		PostID:      testPost.ID().String(),
		AuthorName:  "Test Author",
		AuthorEmail: "test@example.com",
		Content:     "Test comment",
	})
	if err != nil {
		t.Errorf("CreateComment() unexpected error: %v", err)
		return
	}

	if !invalidateCalled {
		t.Error("CreateComment() did not invalidate post comment cache")
	}
}
