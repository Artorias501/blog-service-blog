package repository

import (
	"context"

	"github.com/artorias501/blog-service/internal/domain/entity"
	"github.com/artorias501/blog-service/internal/domain/valueobject"
)

// TagRepository defines the contract for Tag persistence operations
type TagRepository interface {
	// Create inserts a new tag into the database
	Create(ctx context.Context, tag *entity.Tag) error

	// GetByID retrieves a tag by its ID
	GetByID(ctx context.Context, id valueobject.TagID) (*entity.Tag, error)

	// GetByName retrieves a tag by its name
	GetByName(ctx context.Context, name valueobject.TagName) (*entity.Tag, error)

	// Update updates an existing tag in the database
	Update(ctx context.Context, tag *entity.Tag) error

	// Delete removes a tag from the database
	Delete(ctx context.Context, id valueobject.TagID) error

	// List retrieves tags with pagination and sorting
	List(ctx context.Context, params ListParams) (*TagListResult, error)

	// ListByPostID retrieves all tags associated with a post
	ListByPostID(ctx context.Context, postID valueobject.PostID) ([]*entity.Tag, error)

	// GetOrCreate retrieves a tag by name or creates it if not exists
	GetOrCreate(ctx context.Context, name valueobject.TagName) (*entity.Tag, error)

	// Search searches tags by name pattern
	Search(ctx context.Context, keyword string, params ListParams) (*TagListResult, error)
}

// TagCacheRepository defines the contract for Tag cache operations
type TagCacheRepository interface {
	// Get retrieves a cached tag by ID
	Get(ctx context.Context, id valueobject.TagID) (*entity.Tag, error)

	// Set stores a tag in cache
	Set(ctx context.Context, tag *entity.Tag) error

	// Delete removes a tag from cache
	Delete(ctx context.Context, id valueobject.TagID) error

	// GetByName retrieves a cached tag by name
	GetByName(ctx context.Context, name valueobject.TagName) (*entity.Tag, error)

	// SetByName stores a tag in cache by name key
	SetByName(ctx context.Context, name valueobject.TagName, tag *entity.Tag) error

	// GetList retrieves a cached list of tags
	GetList(ctx context.Context, params ListParams) (*TagListResult, error)

	// SetList stores a list of tags in cache
	SetList(ctx context.Context, params ListParams, result *TagListResult) error

	// DeleteList removes all cached tag lists
	DeleteList(ctx context.Context) error

	// GetByPostID retrieves cached tags by post ID
	GetByPostID(ctx context.Context, postID valueobject.PostID) ([]*entity.Tag, error)

	// SetByPostID stores tags by post ID in cache
	SetByPostID(ctx context.Context, postID valueobject.PostID, tags []*entity.Tag) error

	// InvalidateTag invalidates all cache entries related to a tag
	InvalidateTag(ctx context.Context, id valueobject.TagID) error
}
