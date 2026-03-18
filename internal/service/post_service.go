package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/artorias501/blog-service/internal/domain/entity"
	"github.com/artorias501/blog-service/internal/domain/repository"
	"github.com/artorias501/blog-service/internal/domain/valueobject"
)

// PostService provides business logic operations for posts
type PostService struct {
	postRepo     repository.PostRepository
	postCache    repository.PostCacheRepository
	tagRepo      repository.TagRepository
	tagCache     repository.TagCacheRepository
	commentRepo  repository.CommentRepository
	commentCache repository.CommentCacheRepository
}

// NewPostService creates a new PostService instance
func NewPostService(
	postRepo repository.PostRepository,
	postCache repository.PostCacheRepository,
	tagRepo repository.TagRepository,
	tagCache repository.TagCacheRepository,
	commentRepo repository.CommentRepository,
	commentCache repository.CommentCacheRepository,
) *PostService {
	return &PostService{
		postRepo:     postRepo,
		postCache:    postCache,
		tagRepo:      tagRepo,
		tagCache:     tagCache,
		commentRepo:  commentRepo,
		commentCache: commentCache,
	}
}

// CreatePostInput contains the data needed to create a new post
type CreatePostInput struct {
	Title   string
	Content string
	TagIDs  []string
}

// CreatePost creates a new post with the given input
func (s *PostService) CreatePost(ctx context.Context, input CreatePostInput) (*entity.Post, error) {
	// Validate and create value objects
	title, err := valueobject.NewTitle(input.Title)
	if err != nil {
		return nil, fmt.Errorf("invalid title: %w", err)
	}

	content, err := valueobject.NewContent(input.Content)
	if err != nil {
		return nil, fmt.Errorf("invalid content: %w", err)
	}

	// Create post entity
	post := entity.NewPost(title, content)

	// Persist the post
	if err := s.postRepo.Create(ctx, post); err != nil {
		return nil, fmt.Errorf("failed to create post: %w", err)
	}

	// Associate tags if provided
	for _, tagIDStr := range input.TagIDs {
		tagID, err := valueobject.NewTagID(tagIDStr)
		if err != nil {
			continue // Skip invalid tag IDs
		}
		if err := s.postRepo.AddTag(ctx, post.ID(), tagID); err != nil {
			// Log error but don't fail the post creation
			continue
		}
	}

	// Invalidate list caches
	if err := s.postCache.DeleteList(ctx); err != nil {
		// Log error but don't fail the operation
		// Cache invalidation failure is not critical
	}

	return post, nil
}

// GetPostByID retrieves a post by ID using cache-aside pattern
func (s *PostService) GetPostByID(ctx context.Context, id string) (*entity.Post, error) {
	// Validate ID
	postID, err := valueobject.NewPostID(id)
	if err != nil {
		return nil, fmt.Errorf("invalid post ID: %w", err)
	}

	// Check cache first (cache-aside pattern)
	cachedPost, err := s.postCache.Get(ctx, postID)
	if err != nil {
		// Log cache error but continue to database
		// Cache failure should not block the operation
	}
	if cachedPost != nil {
		return cachedPost, nil
	}

	// Cache miss - fetch from database
	post, err := s.postRepo.GetByID(ctx, postID)
	if err != nil {
		return nil, fmt.Errorf("failed to get post: %w", err)
	}
	if post == nil {
		return nil, errors.New("post not found")
	}

	// Populate cache
	if err := s.postCache.Set(ctx, post); err != nil {
		// Log cache error but don't fail the operation
		// Cache population failure is not critical
	}

	return post, nil
}

// GetPostByIDWithTags retrieves a post by ID with tags using cache-aside pattern
func (s *PostService) GetPostByIDWithTags(ctx context.Context, id string) (*entity.Post, error) {
	postID, err := valueobject.NewPostID(id)
	if err != nil {
		return nil, fmt.Errorf("invalid post ID: %w", err)
	}

	// Check cache first
	cachedPost, err := s.postCache.Get(ctx, postID)
	if err == nil && cachedPost != nil {
		// Check if tags are cached
		cachedTags, err := s.tagCache.GetByPostID(ctx, postID)
		if err == nil && cachedTags != nil {
			// Reconstruct post with tags
			for _, tag := range cachedTags {
				cachedPost.AddTag(tag)
			}
			return cachedPost, nil
		}
	}

	// Cache miss - fetch from database
	post, err := s.postRepo.GetByIDWithTags(ctx, postID)
	if err != nil {
		return nil, fmt.Errorf("failed to get post with tags: %w", err)
	}
	if post == nil {
		return nil, errors.New("post not found")
	}

	// Populate caches
	s.postCache.Set(ctx, post)
	if len(post.Tags()) > 0 {
		s.tagCache.SetByPostID(ctx, postID, post.Tags())
	}

	return post, nil
}

// UpdatePostInput contains the data needed to update a post
type UpdatePostInput struct {
	Title   *string
	Content *string
}

// UpdatePost updates an existing post
func (s *PostService) UpdatePost(ctx context.Context, id string, input UpdatePostInput) (*entity.Post, error) {
	postID, err := valueobject.NewPostID(id)
	if err != nil {
		return nil, fmt.Errorf("invalid post ID: %w", err)
	}

	// Get existing post
	post, err := s.postRepo.GetByID(ctx, postID)
	if err != nil {
		return nil, fmt.Errorf("failed to get post: %w", err)
	}
	if post == nil {
		return nil, errors.New("post not found")
	}

	// Update fields if provided
	if input.Title != nil {
		title, err := valueobject.NewTitle(*input.Title)
		if err != nil {
			return nil, fmt.Errorf("invalid title: %w", err)
		}
		post.UpdateTitle(title)
	}

	if input.Content != nil {
		content, err := valueobject.NewContent(*input.Content)
		if err != nil {
			return nil, fmt.Errorf("invalid content: %w", err)
		}
		post.UpdateContent(content)
	}

	// Persist changes
	if err := s.postRepo.Update(ctx, post); err != nil {
		return nil, fmt.Errorf("failed to update post: %w", err)
	}

	// Invalidate caches
	s.postCache.InvalidatePost(ctx, postID)
	s.postCache.DeleteList(ctx)

	return post, nil
}

// DeletePost deletes a post by ID
func (s *PostService) DeletePost(ctx context.Context, id string) error {
	postID, err := valueobject.NewPostID(id)
	if err != nil {
		return fmt.Errorf("invalid post ID: %w", err)
	}

	// Check if post exists
	post, err := s.postRepo.GetByID(ctx, postID)
	if err != nil {
		return fmt.Errorf("failed to get post: %w", err)
	}
	if post == nil {
		return errors.New("post not found")
	}

	// Delete associated comments first
	if err := s.commentRepo.DeleteByPostID(ctx, postID); err != nil {
		return fmt.Errorf("failed to delete post comments: %w", err)
	}

	// Delete the post
	if err := s.postRepo.Delete(ctx, postID); err != nil {
		return fmt.Errorf("failed to delete post: %w", err)
	}

	// Invalidate caches
	s.postCache.InvalidatePost(ctx, postID)
	s.postCache.DeleteList(ctx)
	s.commentCache.InvalidateByPostID(ctx, postID)

	return nil
}

// ListPostsInput contains parameters for listing posts
type ListPostsInput struct {
	Page     int
	PageSize int
	SortBy   string
	Order    string
}

// ListPosts retrieves a paginated list of posts
func (s *PostService) ListPosts(ctx context.Context, input ListPostsInput) (*repository.PostListResult, error) {
	// Set defaults
	if input.Page <= 0 {
		input.Page = 1
	}
	if input.PageSize <= 0 {
		input.PageSize = 10
	}
	if input.PageSize > 100 {
		input.PageSize = 100
	}

	params := repository.ListParams{
		Page:     input.Page,
		PageSize: input.PageSize,
		SortBy:   input.SortBy,
		Order:    input.Order,
	}

	// Check cache first
	cachedResult, err := s.postCache.GetList(ctx, params)
	if err == nil && cachedResult != nil {
		return cachedResult, nil
	}

	// Cache miss - fetch from database
	result, err := s.postRepo.List(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to list posts: %w", err)
	}

	// Populate cache
	s.postCache.SetList(ctx, params, result)

	return result, nil
}

// ListPostsByTag retrieves posts filtered by tag
func (s *PostService) ListPostsByTag(ctx context.Context, tagID string, input ListPostsInput) (*repository.PostListResult, error) {
	// Validate tag ID
	tid, err := valueobject.NewTagID(tagID)
	if err != nil {
		return nil, fmt.Errorf("invalid tag ID: %w", err)
	}

	// Set defaults
	if input.Page <= 0 {
		input.Page = 1
	}
	if input.PageSize <= 0 {
		input.PageSize = 10
	}

	params := repository.ListParams{
		Page:     input.Page,
		PageSize: input.PageSize,
		SortBy:   input.SortBy,
		Order:    input.Order,
	}

	// Check cache first
	cachedResult, err := s.postCache.GetByTagID(ctx, tid, params)
	if err == nil && cachedResult != nil {
		return cachedResult, nil
	}

	// Cache miss - fetch from database
	result, err := s.postRepo.ListByTagID(ctx, tid, params)
	if err != nil {
		return nil, fmt.Errorf("failed to list posts by tag: %w", err)
	}

	// Populate cache
	s.postCache.SetByTagID(ctx, tid, params, result)

	return result, nil
}

// SearchPostsInput contains search criteria
type SearchPostsInput struct {
	Keyword  string
	Page     int
	PageSize int
	SortBy   string
	Order    string
}

// SearchPosts searches posts by keyword with combined criteria
func (s *PostService) SearchPosts(ctx context.Context, input SearchPostsInput) (*repository.PostListResult, error) {
	// Set defaults
	if input.Page <= 0 {
		input.Page = 1
	}
	if input.PageSize <= 0 {
		input.PageSize = 10
	}
	if input.PageSize > 100 {
		input.PageSize = 100
	}

	params := repository.ListParams{
		Page:     input.Page,
		PageSize: input.PageSize,
		SortBy:   input.SortBy,
		Order:    input.Order,
	}

	// Search in repository (no caching for search results)
	result, err := s.postRepo.Search(ctx, input.Keyword, params)
	if err != nil {
		return nil, fmt.Errorf("failed to search posts: %w", err)
	}

	return result, nil
}

// LikePost increments the like count for a post
// Note: This is a placeholder implementation as the current entity doesn't have a like count
// In a real implementation, this would use a separate likes table or counter
func (s *PostService) LikePost(ctx context.Context, id string) error {
	postID, err := valueobject.NewPostID(id)
	if err != nil {
		return fmt.Errorf("invalid post ID: %w", err)
	}

	// Verify post exists
	post, err := s.postRepo.GetByID(ctx, postID)
	if err != nil {
		return fmt.Errorf("failed to get post: %w", err)
	}
	if post == nil {
		return errors.New("post not found")
	}

	// In a real implementation, this would:
	// 1. Increment a like counter in a separate table
	// 2. Invalidate the post cache
	// For now, we just invalidate the cache to signal the operation
	s.postCache.InvalidatePost(ctx, postID)

	return nil
}

// AddTagToPost associates a tag with a post
func (s *PostService) AddTagToPost(ctx context.Context, postID string, tagID string) error {
	pid, err := valueobject.NewPostID(postID)
	if err != nil {
		return fmt.Errorf("invalid post ID: %w", err)
	}

	tid, err := valueobject.NewTagID(tagID)
	if err != nil {
		return fmt.Errorf("invalid tag ID: %w", err)
	}

	// Verify post exists
	post, err := s.postRepo.GetByID(ctx, pid)
	if err != nil {
		return fmt.Errorf("failed to get post: %w", err)
	}
	if post == nil {
		return errors.New("post not found")
	}

	// Verify tag exists
	tag, err := s.tagRepo.GetByID(ctx, tid)
	if err != nil {
		return fmt.Errorf("failed to get tag: %w", err)
	}
	if tag == nil {
		return errors.New("tag not found")
	}

	// Add tag association
	if err := s.postRepo.AddTag(ctx, pid, tid); err != nil {
		return fmt.Errorf("failed to add tag to post: %w", err)
	}

	// Invalidate caches
	s.postCache.InvalidatePost(ctx, pid)
	s.tagCache.InvalidateTag(ctx, tid)

	return nil
}

// RemoveTagFromPost disassociates a tag from a post
func (s *PostService) RemoveTagFromPost(ctx context.Context, postID string, tagID string) error {
	pid, err := valueobject.NewPostID(postID)
	if err != nil {
		return fmt.Errorf("invalid post ID: %w", err)
	}

	tid, err := valueobject.NewTagID(tagID)
	if err != nil {
		return fmt.Errorf("invalid tag ID: %w", err)
	}

	// Remove tag association
	if err := s.postRepo.RemoveTag(ctx, pid, tid); err != nil {
		return fmt.Errorf("failed to remove tag from post: %w", err)
	}

	// Invalidate caches
	s.postCache.InvalidatePost(ctx, pid)
	s.tagCache.InvalidateTag(ctx, tid)

	return nil
}
