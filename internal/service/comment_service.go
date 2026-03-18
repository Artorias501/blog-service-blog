package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/artorias501/blog-service/internal/domain/entity"
	"github.com/artorias501/blog-service/internal/domain/repository"
	"github.com/artorias501/blog-service/internal/domain/valueobject"
)

// CommentService provides business logic operations for comments
type CommentService struct {
	commentRepo  repository.CommentRepository
	commentCache repository.CommentCacheRepository
	postRepo     repository.PostRepository
	postCache    repository.PostCacheRepository
}

// NewCommentService creates a new CommentService instance
func NewCommentService(
	commentRepo repository.CommentRepository,
	commentCache repository.CommentCacheRepository,
	postRepo repository.PostRepository,
	postCache repository.PostCacheRepository,
) *CommentService {
	return &CommentService{
		commentRepo:  commentRepo,
		commentCache: commentCache,
		postRepo:     postRepo,
		postCache:    postCache,
	}
}

// CreateCommentInput contains the data needed to create a new comment
type CreateCommentInput struct {
	PostID      string
	AuthorName  string
	AuthorEmail string
	Content     string
}

// CreateComment creates a new comment with the given input
func (s *CommentService) CreateComment(ctx context.Context, input CreateCommentInput) (*entity.Comment, error) {
	// Validate post ID
	postID, err := valueobject.NewPostID(input.PostID)
	if err != nil {
		return nil, fmt.Errorf("invalid post ID: %w", err)
	}

	// Verify post exists
	post, err := s.postRepo.GetByID(ctx, postID)
	if err != nil {
		return nil, fmt.Errorf("failed to get post: %w", err)
	}
	if post == nil {
		return nil, errors.New("post not found")
	}

	// Validate and create value objects
	authorName, err := valueobject.NewAuthorName(input.AuthorName)
	if err != nil {
		return nil, fmt.Errorf("invalid author name: %w", err)
	}

	authorEmail, err := valueobject.NewAuthorEmail(input.AuthorEmail)
	if err != nil {
		return nil, fmt.Errorf("invalid author email: %w", err)
	}

	content, err := valueobject.NewContent(input.Content)
	if err != nil {
		return nil, fmt.Errorf("invalid content: %w", err)
	}

	// Create comment entity
	comment := entity.NewComment(authorName, authorEmail, content)

	// Persist the comment
	if err := s.commentRepo.Create(ctx, comment); err != nil {
		return nil, fmt.Errorf("failed to create comment: %w", err)
	}

	// Invalidate caches
	s.commentCache.InvalidateByPostID(ctx, postID)

	return comment, nil
}

// GetCommentByID retrieves a comment by ID
func (s *CommentService) GetCommentByID(ctx context.Context, id string) (*entity.Comment, error) {
	// Validate ID
	commentID, err := valueobject.NewCommentID(id)
	if err != nil {
		return nil, fmt.Errorf("invalid comment ID: %w", err)
	}

	// Check cache first (cache-aside pattern)
	cachedComment, err := s.commentCache.Get(ctx, commentID)
	if err != nil {
		// Log cache error but continue to database
	}
	if cachedComment != nil {
		return cachedComment, nil
	}

	// Cache miss - fetch from database
	comment, err := s.commentRepo.GetByID(ctx, commentID)
	if err != nil {
		return nil, fmt.Errorf("failed to get comment: %w", err)
	}
	if comment == nil {
		return nil, errors.New("comment not found")
	}

	// Populate cache
	s.commentCache.Set(ctx, comment)

	return comment, nil
}

// UpdateCommentInput contains the data needed to update a comment
type UpdateCommentInput struct {
	Content string
}

// UpdateComment updates an existing comment
func (s *CommentService) UpdateComment(ctx context.Context, id string, input UpdateCommentInput) (*entity.Comment, error) {
	commentID, err := valueobject.NewCommentID(id)
	if err != nil {
		return nil, fmt.Errorf("invalid comment ID: %w", err)
	}

	// Get existing comment
	comment, err := s.commentRepo.GetByID(ctx, commentID)
	if err != nil {
		return nil, fmt.Errorf("failed to get comment: %w", err)
	}
	if comment == nil {
		return nil, errors.New("comment not found")
	}

	// Validate and update content
	content, err := valueobject.NewContent(input.Content)
	if err != nil {
		return nil, fmt.Errorf("invalid content: %w", err)
	}

	comment.UpdateContent(content)

	// Persist changes
	if err := s.commentRepo.Update(ctx, comment); err != nil {
		return nil, fmt.Errorf("failed to update comment: %w", err)
	}

	// Invalidate cache
	s.commentCache.InvalidateComment(ctx, commentID)

	return comment, nil
}

// DeleteComment deletes a comment by ID
func (s *CommentService) DeleteComment(ctx context.Context, id string) error {
	commentID, err := valueobject.NewCommentID(id)
	if err != nil {
		return fmt.Errorf("invalid comment ID: %w", err)
	}

	// Check if comment exists
	comment, err := s.commentRepo.GetByID(ctx, commentID)
	if err != nil {
		return fmt.Errorf("failed to get comment: %w", err)
	}
	if comment == nil {
		return errors.New("comment not found")
	}

	// Delete the comment
	if err := s.commentRepo.Delete(ctx, commentID); err != nil {
		return fmt.Errorf("failed to delete comment: %w", err)
	}

	// Invalidate cache
	s.commentCache.InvalidateComment(ctx, commentID)

	return nil
}

// ListCommentsInput contains parameters for listing comments
type ListCommentsInput struct {
	Page     int
	PageSize int
	SortBy   string
	Order    string
}

// ListComments retrieves a paginated list of all comments
func (s *CommentService) ListComments(ctx context.Context, input ListCommentsInput) (*repository.CommentListResult, error) {
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

	// Fetch from database (no caching for admin list)
	result, err := s.commentRepo.List(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to list comments: %w", err)
	}

	return result, nil
}

// ListCommentsByPostInput contains parameters for listing comments by post
type ListCommentsByPostInput struct {
	PostID   string
	Page     int
	PageSize int
	SortBy   string
	Order    string
}

// ListCommentsByPost retrieves comments for a specific post
func (s *CommentService) ListCommentsByPost(ctx context.Context, input ListCommentsByPostInput) (*repository.CommentListResult, error) {
	// Validate post ID
	postID, err := valueobject.NewPostID(input.PostID)
	if err != nil {
		return nil, fmt.Errorf("invalid post ID: %w", err)
	}

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
	cachedResult, err := s.commentCache.GetListByPostID(ctx, postID, params)
	if err == nil && cachedResult != nil {
		return cachedResult, nil
	}

	// Cache miss - fetch from database
	result, err := s.commentRepo.ListByPostID(ctx, postID, params)
	if err != nil {
		return nil, fmt.Errorf("failed to list comments by post: %w", err)
	}

	// Populate cache
	s.commentCache.SetListByPostID(ctx, postID, params, result)

	return result, nil
}

// ListCommentsByStatusInput contains parameters for listing comments by status
type ListCommentsByStatusInput struct {
	Status   string
	Page     int
	PageSize int
	SortBy   string
	Order    string
}

// ListCommentsByStatus retrieves comments filtered by status
func (s *CommentService) ListCommentsByStatus(ctx context.Context, input ListCommentsByStatusInput) (*repository.CommentListResult, error) {
	// Validate status
	status, err := valueobject.NewCommentStatus(input.Status)
	if err != nil {
		return nil, fmt.Errorf("invalid comment status: %w", err)
	}

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

	// Fetch from database (no caching for admin status filter)
	result, err := s.commentRepo.ListByStatus(ctx, status.String(), params)
	if err != nil {
		return nil, fmt.Errorf("failed to list comments by status: %w", err)
	}

	return result, nil
}

// ApproveComment approves a comment by ID
func (s *CommentService) ApproveComment(ctx context.Context, id string) error {
	commentID, err := valueobject.NewCommentID(id)
	if err != nil {
		return fmt.Errorf("invalid comment ID: %w", err)
	}

	// Check if comment exists
	comment, err := s.commentRepo.GetByID(ctx, commentID)
	if err != nil {
		return fmt.Errorf("failed to get comment: %w", err)
	}
	if comment == nil {
		return errors.New("comment not found")
	}

	// Approve the comment
	if err := s.commentRepo.Approve(ctx, commentID); err != nil {
		return fmt.Errorf("failed to approve comment: %w", err)
	}

	// Invalidate cache
	s.commentCache.InvalidateComment(ctx, commentID)

	return nil
}

// RejectComment rejects a comment by ID
func (s *CommentService) RejectComment(ctx context.Context, id string) error {
	commentID, err := valueobject.NewCommentID(id)
	if err != nil {
		return fmt.Errorf("invalid comment ID: %w", err)
	}

	// Check if comment exists
	comment, err := s.commentRepo.GetByID(ctx, commentID)
	if err != nil {
		return fmt.Errorf("failed to get comment: %w", err)
	}
	if comment == nil {
		return errors.New("comment not found")
	}

	// Reject the comment
	if err := s.commentRepo.Reject(ctx, commentID); err != nil {
		return fmt.Errorf("failed to reject comment: %w", err)
	}

	// Invalidate cache
	s.commentCache.InvalidateComment(ctx, commentID)

	return nil
}

// MarkCommentAsSpam marks a comment as spam by ID
func (s *CommentService) MarkCommentAsSpam(ctx context.Context, id string) error {
	commentID, err := valueobject.NewCommentID(id)
	if err != nil {
		return fmt.Errorf("invalid comment ID: %w", err)
	}

	// Check if comment exists
	comment, err := s.commentRepo.GetByID(ctx, commentID)
	if err != nil {
		return fmt.Errorf("failed to get comment: %w", err)
	}
	if comment == nil {
		return errors.New("comment not found")
	}

	// Mark as spam
	if err := s.commentRepo.MarkAsSpam(ctx, commentID); err != nil {
		return fmt.Errorf("failed to mark comment as spam: %w", err)
	}

	// Invalidate cache
	s.commentCache.InvalidateComment(ctx, commentID)

	return nil
}

// GetCommentCountByPostID returns the count of comments for a specific post
func (s *CommentService) GetCommentCountByPostID(ctx context.Context, postID string) (int64, error) {
	pid, err := valueobject.NewPostID(postID)
	if err != nil {
		return 0, fmt.Errorf("invalid post ID: %w", err)
	}

	// Check cache first
	cachedCount, err := s.commentCache.GetCountByPostID(ctx, pid)
	if err == nil && cachedCount >= 0 {
		return cachedCount, nil
	}

	// Cache miss - fetch from database
	count, err := s.commentRepo.CountByPostID(ctx, pid)
	if err != nil {
		return 0, fmt.Errorf("failed to get comment count: %w", err)
	}

	// Populate cache
	s.commentCache.SetCountByPostID(ctx, pid, count)

	return count, nil
}
