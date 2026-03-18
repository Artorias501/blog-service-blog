package repository

import (
	"context"

	"github.com/artorias501/blog-service/internal/domain/entity"
	"github.com/artorias501/blog-service/internal/domain/valueobject"
)

// ListParams contains pagination and sorting parameters for list queries
type ListParams struct {
	Page     int
	PageSize int
	SortBy   string
	Order    string // "asc" or "desc"
}

// ListResult contains paginated results with metadata
type ListResult[T any] struct {
	Total     int64
	Page      int
	PageSize  int
	TotalPage int
	Items     []T
}

// PostListResult contains paginated post results with metadata
type PostListResult struct {
	Total     int64
	Page      int
	PageSize  int
	TotalPage int
	Items     []*entity.Post
}

// TagListResult contains paginated tag results with metadata
type TagListResult struct {
	Total     int64
	Page      int
	PageSize  int
	TotalPage int
	Items     []*entity.Tag
}

// CommentListResult contains paginated comment results with metadata
type CommentListResult struct {
	Total     int64
	Page      int
	PageSize  int
	TotalPage int
	Items     []*entity.Comment
}

// PostRepository defines the contract for Post persistence operations
type PostRepository interface {
	// Create inserts a new post into the database
	Create(ctx context.Context, post *entity.Post) error

	// GetByID retrieves a post by its ID
	GetByID(ctx context.Context, id valueobject.PostID) (*entity.Post, error)

	// Update updates an existing post in the database
	Update(ctx context.Context, post *entity.Post) error

	// Delete removes a post from the database
	Delete(ctx context.Context, id valueobject.PostID) error

	// List retrieves posts with pagination and sorting
	List(ctx context.Context, params ListParams) (*PostListResult, error)

	// ListByTagID retrieves posts associated with a specific tag
	ListByTagID(ctx context.Context, tagID valueobject.TagID, params ListParams) (*PostListResult, error)

	// GetByIDWithComments retrieves a post with its comments loaded
	GetByIDWithComments(ctx context.Context, id valueobject.PostID) (*entity.Post, error)

	// GetByIDWithTags retrieves a post with its tags loaded
	GetByIDWithTags(ctx context.Context, id valueobject.PostID) (*entity.Post, error)

	// GetByIDFull retrieves a post with both comments and tags loaded
	GetByIDFull(ctx context.Context, id valueobject.PostID) (*entity.Post, error)

	// AddTag associates a tag with a post
	AddTag(ctx context.Context, postID valueobject.PostID, tagID valueobject.TagID) error

	// RemoveTag disassociates a tag from a post
	RemoveTag(ctx context.Context, postID valueobject.PostID, tagID valueobject.TagID) error

	// Search searches posts by title or content keyword
	Search(ctx context.Context, keyword string, params ListParams) (*PostListResult, error)
}

// PostCacheRepository defines the contract for Post cache operations
type PostCacheRepository interface {
	// Get retrieves a cached post by ID
	Get(ctx context.Context, id valueobject.PostID) (*entity.Post, error)

	// Set stores a post in cache
	Set(ctx context.Context, post *entity.Post) error

	// Delete removes a post from cache
	Delete(ctx context.Context, id valueobject.PostID) error

	// GetList retrieves a cached list of posts
	GetList(ctx context.Context, params ListParams) (*PostListResult, error)

	// SetList stores a list of posts in cache
	SetList(ctx context.Context, params ListParams, result *PostListResult) error

	// DeleteList removes all cached post lists
	DeleteList(ctx context.Context) error

	// GetByTagID retrieves cached posts by tag ID
	GetByTagID(ctx context.Context, tagID valueobject.TagID, params ListParams) (*PostListResult, error)

	// SetByTagID stores posts by tag ID in cache
	SetByTagID(ctx context.Context, tagID valueobject.TagID, params ListParams, result *PostListResult) error

	// InvalidatePost invalidates all cache entries related to a post
	InvalidatePost(ctx context.Context, id valueobject.PostID) error
}
