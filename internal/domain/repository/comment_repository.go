package repository

import (
	"context"

	"github.com/artorias501/blog-service/internal/domain/entity"
	"github.com/artorias501/blog-service/internal/domain/valueobject"
)

// CommentRepository defines the contract for Comment persistence operations
type CommentRepository interface {
	// Create inserts a new comment into the database
	Create(ctx context.Context, comment *entity.Comment) error

	// GetByID retrieves a comment by its ID
	GetByID(ctx context.Context, id valueobject.CommentID) (*entity.Comment, error)

	// Update updates an existing comment in the database
	Update(ctx context.Context, comment *entity.Comment) error

	// Delete removes a comment from the database
	Delete(ctx context.Context, id valueobject.CommentID) error

	// List retrieves comments with pagination and sorting
	List(ctx context.Context, params ListParams) (*CommentListResult, error)

	// ListByPostID retrieves all comments for a specific post
	ListByPostID(ctx context.Context, postID valueobject.PostID, params ListParams) (*CommentListResult, error)

	// ListByStatus retrieves comments by status (pending, approved, rejected, spam)
	ListByStatus(ctx context.Context, status string, params ListParams) (*CommentListResult, error)

	// CountByPostID returns the count of comments for a specific post
	CountByPostID(ctx context.Context, postID valueobject.PostID) (int64, error)

	// Approve updates a comment's status to approved
	Approve(ctx context.Context, id valueobject.CommentID) error

	// Reject updates a comment's status to rejected
	Reject(ctx context.Context, id valueobject.CommentID) error

	// MarkAsSpam updates a comment's status to spam
	MarkAsSpam(ctx context.Context, id valueobject.CommentID) error

	// DeleteByPostID removes all comments associated with a post
	DeleteByPostID(ctx context.Context, postID valueobject.PostID) error
}

// CommentCacheRepository defines the contract for Comment cache operations
type CommentCacheRepository interface {
	// Get retrieves a cached comment by ID
	Get(ctx context.Context, id valueobject.CommentID) (*entity.Comment, error)

	// Set stores a comment in cache
	Set(ctx context.Context, comment *entity.Comment) error

	// Delete removes a comment from cache
	Delete(ctx context.Context, id valueobject.CommentID) error

	// GetListByPostID retrieves cached comments by post ID
	GetListByPostID(ctx context.Context, postID valueobject.PostID, params ListParams) (*CommentListResult, error)

	// SetListByPostID stores comments by post ID in cache
	SetListByPostID(ctx context.Context, postID valueobject.PostID, params ListParams, result *CommentListResult) error

	// DeleteListByPostID removes cached comments for a post
	DeleteListByPostID(ctx context.Context, postID valueobject.PostID) error

	// GetCountByPostID retrieves cached comment count for a post
	GetCountByPostID(ctx context.Context, postID valueobject.PostID) (int64, error)

	// SetCountByPostID stores comment count for a post in cache
	SetCountByPostID(ctx context.Context, postID valueobject.PostID, count int64) error

	// InvalidateComment invalidates all cache entries related to a comment
	InvalidateComment(ctx context.Context, id valueobject.CommentID) error

	// InvalidateByPostID invalidates all comment cache entries for a post
	InvalidateByPostID(ctx context.Context, postID valueobject.PostID) error
}
