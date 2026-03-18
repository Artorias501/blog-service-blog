package service

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/artorias501/blog-service/internal/domain/entity"
	"github.com/artorias501/blog-service/internal/domain/repository"
	"github.com/artorias501/blog-service/internal/domain/valueobject"
)

// Mock implementations for testing

type mockPostRepository struct {
	posts     map[string]*entity.Post
	createFn  func(ctx context.Context, post *entity.Post) error
	getByIDFn func(ctx context.Context, id valueobject.PostID) (*entity.Post, error)
	updateFn  func(ctx context.Context, post *entity.Post) error
	deleteFn  func(ctx context.Context, id valueobject.PostID) error
	listFn    func(ctx context.Context, params repository.ListParams) (*repository.PostListResult, error)
	searchFn  func(ctx context.Context, keyword string, params repository.ListParams) (*repository.PostListResult, error)
	addTagFn  func(ctx context.Context, postID valueobject.PostID, tagID valueobject.TagID) error
}

func newMockPostRepository() *mockPostRepository {
	return &mockPostRepository{
		posts: make(map[string]*entity.Post),
	}
}

func (m *mockPostRepository) Create(ctx context.Context, post *entity.Post) error {
	if m.createFn != nil {
		return m.createFn(ctx, post)
	}
	m.posts[post.ID().String()] = post
	return nil
}

func (m *mockPostRepository) GetByID(ctx context.Context, id valueobject.PostID) (*entity.Post, error) {
	if m.getByIDFn != nil {
		return m.getByIDFn(ctx, id)
	}
	post, ok := m.posts[id.String()]
	if !ok {
		return nil, nil
	}
	return post, nil
}

func (m *mockPostRepository) Update(ctx context.Context, post *entity.Post) error {
	if m.updateFn != nil {
		return m.updateFn(ctx, post)
	}
	m.posts[post.ID().String()] = post
	return nil
}

func (m *mockPostRepository) Delete(ctx context.Context, id valueobject.PostID) error {
	if m.deleteFn != nil {
		return m.deleteFn(ctx, id)
	}
	delete(m.posts, id.String())
	return nil
}

func (m *mockPostRepository) List(ctx context.Context, params repository.ListParams) (*repository.PostListResult, error) {
	if m.listFn != nil {
		return m.listFn(ctx, params)
	}
	var posts []*entity.Post
	for _, p := range m.posts {
		posts = append(posts, p)
	}
	return &repository.PostListResult{
		Items: posts,
		Total: int64(len(posts)),
	}, nil
}

func (m *mockPostRepository) ListByTagID(ctx context.Context, tagID valueobject.TagID, params repository.ListParams) (*repository.PostListResult, error) {
	return m.listFn(ctx, params)
}

func (m *mockPostRepository) GetByIDWithComments(ctx context.Context, id valueobject.PostID) (*entity.Post, error) {
	return m.GetByID(ctx, id)
}

func (m *mockPostRepository) GetByIDWithTags(ctx context.Context, id valueobject.PostID) (*entity.Post, error) {
	return m.GetByID(ctx, id)
}

func (m *mockPostRepository) GetByIDFull(ctx context.Context, id valueobject.PostID) (*entity.Post, error) {
	return m.GetByID(ctx, id)
}

func (m *mockPostRepository) AddTag(ctx context.Context, postID valueobject.PostID, tagID valueobject.TagID) error {
	if m.addTagFn != nil {
		return m.addTagFn(ctx, postID, tagID)
	}
	return nil
}

func (m *mockPostRepository) RemoveTag(ctx context.Context, postID valueobject.PostID, tagID valueobject.TagID) error {
	return nil
}

func (m *mockPostRepository) Search(ctx context.Context, keyword string, params repository.ListParams) (*repository.PostListResult, error) {
	if m.searchFn != nil {
		return m.searchFn(ctx, keyword, params)
	}
	return m.List(ctx, params)
}

type mockPostCacheRepository struct {
	cache        map[string]*entity.Post
	listCache    map[string]*repository.PostListResult
	getFn        func(ctx context.Context, id valueobject.PostID) (*entity.Post, error)
	setFn        func(ctx context.Context, post *entity.Post) error
	deleteFn     func(ctx context.Context, id valueobject.PostID) error
	getListFn    func(ctx context.Context, params repository.ListParams) (*repository.PostListResult, error)
	setListFn    func(ctx context.Context, params repository.ListParams, result *repository.PostListResult) error
	deleteListFn func(ctx context.Context) error
	invalidateFn func(ctx context.Context, id valueobject.PostID) error
}

func newMockPostCacheRepository() *mockPostCacheRepository {
	return &mockPostCacheRepository{
		cache:     make(map[string]*entity.Post),
		listCache: make(map[string]*repository.PostListResult),
	}
}

func (m *mockPostCacheRepository) Get(ctx context.Context, id valueobject.PostID) (*entity.Post, error) {
	if m.getFn != nil {
		return m.getFn(ctx, id)
	}
	return m.cache[id.String()], nil
}

func (m *mockPostCacheRepository) Set(ctx context.Context, post *entity.Post) error {
	if m.setFn != nil {
		return m.setFn(ctx, post)
	}
	m.cache[post.ID().String()] = post
	return nil
}

func (m *mockPostCacheRepository) Delete(ctx context.Context, id valueobject.PostID) error {
	if m.deleteFn != nil {
		return m.deleteFn(ctx, id)
	}
	delete(m.cache, id.String())
	return nil
}

func (m *mockPostCacheRepository) GetList(ctx context.Context, params repository.ListParams) (*repository.PostListResult, error) {
	if m.getListFn != nil {
		return m.getListFn(ctx, params)
	}
	key := fmt.Sprintf("post:list:page:%d:size:%d", params.Page, params.PageSize)
	return m.listCache[key], nil
}

func (m *mockPostCacheRepository) SetList(ctx context.Context, params repository.ListParams, result *repository.PostListResult) error {
	if m.setListFn != nil {
		return m.setListFn(ctx, params, result)
	}
	key := fmt.Sprintf("post:list:page:%d:size:%d", params.Page, params.PageSize)
	m.listCache[key] = result
	return nil
}

func (m *mockPostCacheRepository) DeleteList(ctx context.Context) error {
	if m.deleteListFn != nil {
		return m.deleteListFn(ctx)
	}
	m.listCache = make(map[string]*repository.PostListResult)
	return nil
}

func (m *mockPostCacheRepository) GetByTagID(ctx context.Context, tagID valueobject.TagID, params repository.ListParams) (*repository.PostListResult, error) {
	return nil, nil
}

func (m *mockPostCacheRepository) SetByTagID(ctx context.Context, tagID valueobject.TagID, params repository.ListParams, result *repository.PostListResult) error {
	return nil
}

func (m *mockPostCacheRepository) InvalidatePost(ctx context.Context, id valueobject.PostID) error {
	if m.invalidateFn != nil {
		return m.invalidateFn(ctx, id)
	}
	delete(m.cache, id.String())
	return nil
}

// Helper function to create a test post
func createTestPost(t *testing.T) *entity.Post {
	title, err := valueobject.NewTitle("Test Post")
	if err != nil {
		t.Fatalf("failed to create title: %v", err)
	}
	content, err := valueobject.NewContent("Test content for the post")
	if err != nil {
		t.Fatalf("failed to create content: %v", err)
	}
	return entity.NewPost(title, content)
}

// Test CreatePost
func TestPostService_CreatePost(t *testing.T) {
	tests := []struct {
		name    string
		input   CreatePostInput
		wantErr bool
	}{
		{
			name: "valid post",
			input: CreatePostInput{
				Title:   "Test Post",
				Content: "Test content",
			},
			wantErr: false,
		},
		{
			name: "empty title",
			input: CreatePostInput{
				Title:   "",
				Content: "Test content",
			},
			wantErr: true,
		},
		{
			name: "empty content",
			input: CreatePostInput{
				Title:   "Test Post",
				Content: "",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			postRepo := newMockPostRepository()
			postCache := newMockPostCacheRepository()
			tagRepo := newMockTagRepository()
			tagCache := newMockTagCacheRepository()
			commentRepo := newMockCommentRepository()
			commentCache := newMockCommentCacheRepository()

			svc := NewPostService(postRepo, postCache, tagRepo, tagCache, commentRepo, commentCache)

			post, err := svc.CreatePost(ctx, tt.input)

			if tt.wantErr {
				if err == nil {
					t.Errorf("CreatePost() expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("CreatePost() unexpected error: %v", err)
				return
			}

			if post == nil {
				t.Error("CreatePost() returned nil post")
				return
			}

			if post.Title().String() != tt.input.Title {
				t.Errorf("CreatePost() title = %v, want %v", post.Title().String(), tt.input.Title)
			}

			if post.Content().String() != tt.input.Content {
				t.Errorf("CreatePost() content = %v, want %v", post.Content().String(), tt.input.Content)
			}
		})
	}
}

// Test GetPostByID with cache-aside pattern
func TestPostService_GetPostByID_CacheAside(t *testing.T) {
	t.Run("cache hit", func(t *testing.T) {
		ctx := context.Background()
		postRepo := newMockPostRepository()
		postCache := newMockPostCacheRepository()
		tagRepo := newMockTagRepository()
		tagCache := newMockTagCacheRepository()
		commentRepo := newMockCommentRepository()
		commentCache := newMockCommentCacheRepository()

		// Create and cache a post
		testPost := createTestPost(t)
		postRepo.posts[testPost.ID().String()] = testPost
		postCache.cache[testPost.ID().String()] = testPost

		svc := NewPostService(postRepo, postCache, tagRepo, tagCache, commentRepo, commentCache)

		post, err := svc.GetPostByID(ctx, testPost.ID().String())
		if err != nil {
			t.Errorf("GetPostByID() unexpected error: %v", err)
			return
		}

		if post == nil {
			t.Error("GetPostByID() returned nil post")
			return
		}

		if post.ID().String() != testPost.ID().String() {
			t.Errorf("GetPostByID() returned wrong post")
		}
	})

	t.Run("cache miss - fetch from database", func(t *testing.T) {
		ctx := context.Background()
		postRepo := newMockPostRepository()
		postCache := newMockPostCacheRepository()
		tagRepo := newMockTagRepository()
		tagCache := newMockTagCacheRepository()
		commentRepo := newMockCommentRepository()
		commentCache := newMockCommentCacheRepository()

		// Create a post in repository but not in cache
		testPost := createTestPost(t)
		postRepo.posts[testPost.ID().String()] = testPost

		svc := NewPostService(postRepo, postCache, tagRepo, tagCache, commentRepo, commentCache)

		post, err := svc.GetPostByID(ctx, testPost.ID().String())
		if err != nil {
			t.Errorf("GetPostByID() unexpected error: %v", err)
			return
		}

		if post == nil {
			t.Error("GetPostByID() returned nil post")
			return
		}

		// Verify post was cached
		cachedPost := postCache.cache[testPost.ID().String()]
		if cachedPost == nil {
			t.Error("GetPostByID() did not populate cache")
		}
	})

	t.Run("post not found", func(t *testing.T) {
		ctx := context.Background()
		postRepo := newMockPostRepository()
		postCache := newMockPostCacheRepository()
		tagRepo := newMockTagRepository()
		tagCache := newMockTagCacheRepository()
		commentRepo := newMockCommentRepository()
		commentCache := newMockCommentCacheRepository()

		svc := NewPostService(postRepo, postCache, tagRepo, tagCache, commentRepo, commentCache)

		postID := valueobject.GeneratePostID()
		_, err := svc.GetPostByID(ctx, postID.String())
		if err == nil {
			t.Error("GetPostByID() expected error for non-existent post")
		}
	})
}

// Test UpdatePost with cache invalidation
func TestPostService_UpdatePost_CacheInvalidation(t *testing.T) {
	t.Run("update and invalidate cache", func(t *testing.T) {
		ctx := context.Background()
		postRepo := newMockPostRepository()
		postCache := newMockPostCacheRepository()
		tagRepo := newMockTagRepository()
		tagCache := newMockTagCacheRepository()
		commentRepo := newMockCommentRepository()
		commentCache := newMockCommentCacheRepository()

		// Create and cache a post
		testPost := createTestPost(t)
		postRepo.posts[testPost.ID().String()] = testPost
		postCache.cache[testPost.ID().String()] = testPost

		svc := NewPostService(postRepo, postCache, tagRepo, tagCache, commentRepo, commentCache)

		newTitle := "Updated Title"
		input := UpdatePostInput{
			Title: &newTitle,
		}

		updatedPost, err := svc.UpdatePost(ctx, testPost.ID().String(), input)
		if err != nil {
			t.Errorf("UpdatePost() unexpected error: %v", err)
			return
		}

		if updatedPost.Title().String() != newTitle {
			t.Errorf("UpdatePost() title = %v, want %v", updatedPost.Title().String(), newTitle)
		}

		// Verify cache was invalidated
		if postCache.cache[testPost.ID().String()] != nil {
			t.Error("UpdatePost() did not invalidate cache")
		}
	})
}

// Test DeletePost with cache invalidation
func TestPostService_DeletePost_CacheInvalidation(t *testing.T) {
	t.Run("delete and invalidate cache", func(t *testing.T) {
		ctx := context.Background()
		postRepo := newMockPostRepository()
		postCache := newMockPostCacheRepository()
		tagRepo := newMockTagRepository()
		tagCache := newMockTagCacheRepository()
		commentRepo := newMockCommentRepository()
		commentCache := newMockCommentCacheRepository()

		// Create and cache a post
		testPost := createTestPost(t)
		postRepo.posts[testPost.ID().String()] = testPost
		postCache.cache[testPost.ID().String()] = testPost

		svc := NewPostService(postRepo, postCache, tagRepo, tagCache, commentRepo, commentCache)

		err := svc.DeletePost(ctx, testPost.ID().String())
		if err != nil {
			t.Errorf("DeletePost() unexpected error: %v", err)
			return
		}

		// Verify post was deleted from repository
		if postRepo.posts[testPost.ID().String()] != nil {
			t.Error("DeletePost() did not delete post from repository")
		}

		// Verify cache was invalidated
		if postCache.cache[testPost.ID().String()] != nil {
			t.Error("DeletePost() did not invalidate cache")
		}
	})
}

// Test ListPosts with cache integration
func TestPostService_ListPosts_CacheIntegration(t *testing.T) {
	t.Run("cache hit", func(t *testing.T) {
		ctx := context.Background()
		postRepo := newMockPostRepository()
		postCache := newMockPostCacheRepository()
		tagRepo := newMockTagRepository()
		tagCache := newMockTagCacheRepository()
		commentRepo := newMockCommentRepository()
		commentCache := newMockCommentCacheRepository()

		// Create test posts
		testPost1 := createTestPost(t)
		testPost2 := createTestPost(t)
		postRepo.posts[testPost1.ID().String()] = testPost1
		postRepo.posts[testPost2.ID().String()] = testPost2

		// Pre-populate cache
		cachedResult := &repository.PostListResult{
			Items: []*entity.Post{testPost1},
			Total: 1,
		}
		postCache.listCache["post:list:page:1:size:10"] = cachedResult

		svc := NewPostService(postRepo, postCache, tagRepo, tagCache, commentRepo, commentCache)

		result, err := svc.ListPosts(ctx, ListPostsInput{Page: 1, PageSize: 10})
		if err != nil {
			t.Errorf("ListPosts() unexpected error: %v", err)
			return
		}

		if result.Total != 1 {
			t.Errorf("ListPosts() expected cached result with 1 item, got %d", result.Total)
		}
	})

	t.Run("cache miss - fetch from database", func(t *testing.T) {
		ctx := context.Background()
		postRepo := newMockPostRepository()
		postCache := newMockPostCacheRepository()
		tagRepo := newMockTagRepository()
		tagCache := newMockTagCacheRepository()
		commentRepo := newMockCommentRepository()
		commentCache := newMockCommentCacheRepository()

		// Create test posts
		testPost := createTestPost(t)
		postRepo.posts[testPost.ID().String()] = testPost

		svc := NewPostService(postRepo, postCache, tagRepo, tagCache, commentRepo, commentCache)

		result, err := svc.ListPosts(ctx, ListPostsInput{Page: 1, PageSize: 10})
		if err != nil {
			t.Errorf("ListPosts() unexpected error: %v", err)
			return
		}

		if result.Total < 1 {
			t.Errorf("ListPosts() expected at least 1 item, got %d", result.Total)
		}

		// Verify result was cached
		if postCache.listCache["post:list:page:1:size:10"] == nil {
			t.Error("ListPosts() did not populate cache")
		}
	})
}

// Test SearchPosts with combined criteria
func TestPostService_SearchPosts(t *testing.T) {
	t.Run("search with keyword", func(t *testing.T) {
		ctx := context.Background()
		postRepo := newMockPostRepository()
		postCache := newMockPostCacheRepository()
		tagRepo := newMockTagRepository()
		tagCache := newMockTagCacheRepository()
		commentRepo := newMockCommentRepository()
		commentCache := newMockCommentCacheRepository()

		// Create test posts
		testPost := createTestPost(t)
		postRepo.posts[testPost.ID().String()] = testPost

		// Set up search function
		postRepo.searchFn = func(ctx context.Context, keyword string, params repository.ListParams) (*repository.PostListResult, error) {
			var results []*entity.Post
			for _, p := range postRepo.posts {
				// Simple keyword matching simulation
				results = append(results, p)
			}
			return &repository.PostListResult{
				Items: results,
				Total: int64(len(results)),
			}, nil
		}

		svc := NewPostService(postRepo, postCache, tagRepo, tagCache, commentRepo, commentCache)

		result, err := svc.SearchPosts(ctx, SearchPostsInput{
			Keyword:  "Test",
			Page:     1,
			PageSize: 10,
		})
		if err != nil {
			t.Errorf("SearchPosts() unexpected error: %v", err)
			return
		}

		if result.Total < 1 {
			t.Error("SearchPosts() expected at least 1 result")
		}
	})
}

// Test CreatePost invalidates list cache
func TestPostService_CreatePost_InvalidatesListCache(t *testing.T) {
	ctx := context.Background()
	postRepo := newMockPostRepository()
	postCache := newMockPostCacheRepository()
	tagRepo := newMockTagRepository()
	tagCache := newMockTagCacheRepository()
	commentRepo := newMockCommentRepository()
	commentCache := newMockCommentCacheRepository()

	// Pre-populate list cache
	postCache.listCache["post:list:page:1:size:10"] = &repository.PostListResult{Total: 0}

	svc := NewPostService(postRepo, postCache, tagRepo, tagCache, commentRepo, commentCache)

	_, err := svc.CreatePost(ctx, CreatePostInput{
		Title:   "New Post",
		Content: "New content",
	})
	if err != nil {
		t.Errorf("CreatePost() unexpected error: %v", err)
		return
	}

	// Verify list cache was invalidated
	if len(postCache.listCache) != 0 {
		t.Error("CreatePost() did not invalidate list cache")
	}
}

// Test error handling - repository errors wrapped with context
func TestPostService_RepositoryErrorHandling(t *testing.T) {
	t.Run("GetByID repository error", func(t *testing.T) {
		ctx := context.Background()
		postRepo := newMockPostRepository()
		postCache := newMockPostCacheRepository()
		tagRepo := newMockTagRepository()
		tagCache := newMockTagCacheRepository()
		commentRepo := newMockCommentRepository()
		commentCache := newMockCommentCacheRepository()

		expectedErr := errors.New("database connection error")
		postRepo.getByIDFn = func(ctx context.Context, id valueobject.PostID) (*entity.Post, error) {
			return nil, expectedErr
		}

		svc := NewPostService(postRepo, postCache, tagRepo, tagCache, commentRepo, commentCache)

		postID := valueobject.GeneratePostID()
		_, err := svc.GetPostByID(ctx, postID.String())
		if err == nil {
			t.Error("GetPostByID() expected error, got nil")
			return
		}

		// Verify error is wrapped with context
		if err.Error() == expectedErr.Error() {
			t.Error("GetPostByID() error should be wrapped with context")
		}
	})

	t.Run("Create repository error", func(t *testing.T) {
		ctx := context.Background()
		postRepo := newMockPostRepository()
		postCache := newMockPostCacheRepository()
		tagRepo := newMockTagRepository()
		tagCache := newMockTagCacheRepository()
		commentRepo := newMockCommentRepository()
		commentCache := newMockCommentCacheRepository()

		expectedErr := errors.New("insert failed")
		postRepo.createFn = func(ctx context.Context, post *entity.Post) error {
			return expectedErr
		}

		svc := NewPostService(postRepo, postCache, tagRepo, tagCache, commentRepo, commentCache)

		_, err := svc.CreatePost(ctx, CreatePostInput{
			Title:   "Test",
			Content: "Content",
		})
		if err == nil {
			t.Error("CreatePost() expected error, got nil")
			return
		}

		// Verify error is wrapped with context
		if err.Error() == expectedErr.Error() {
			t.Error("CreatePost() error should be wrapped with context")
		}
	})
}

// Test dependency injection - services use interfaces
func TestPostService_DependencyInjection(t *testing.T) {
	// This test verifies that PostService accepts interfaces, not implementations
	var _ repository.PostRepository = (*mockPostRepository)(nil)
	var _ repository.PostCacheRepository = (*mockPostCacheRepository)(nil)
	var _ repository.TagRepository = (*mockTagRepository)(nil)
	var _ repository.TagCacheRepository = (*mockTagCacheRepository)(nil)
	var _ repository.CommentRepository = (*mockCommentRepository)(nil)
	var _ repository.CommentCacheRepository = (*mockCommentCacheRepository)(nil)

	// If this compiles, the service uses interfaces for dependency injection
	svc := NewPostService(
		newMockPostRepository(),
		newMockPostCacheRepository(),
		newMockTagRepository(),
		newMockTagCacheRepository(),
		newMockCommentRepository(),
		newMockCommentCacheRepository(),
	)

	if svc == nil {
		t.Error("NewPostService returned nil")
	}
}
