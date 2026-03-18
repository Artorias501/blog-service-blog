package cache

import (
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/artorias501/blog-service/internal/domain/entity"
	"github.com/artorias501/blog-service/internal/domain/repository"
	"github.com/artorias501/blog-service/internal/domain/valueobject"
)

// TagCacheRepository implements repository.TagCacheRepository using Redis
type TagCacheRepository struct {
	client *RedisClient
}

// NewTagCacheRepository creates a new TagCacheRepository
func NewTagCacheRepository(client *RedisClient) *TagCacheRepository {
	return &TagCacheRepository{
		client: client,
	}
}

// Get retrieves a cached tag by ID
func (r *TagCacheRepository) Get(ctx context.Context, id valueobject.TagID) (*entity.Tag, error) {
	key := TagKey(id)
	data, err := r.client.Get(ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, nil // Cache miss
		}
		return nil, err
	}

	var cached tagCacheData
	if err := json.Unmarshal(data, &cached); err != nil {
		return nil, err
	}

	return cached.ToEntity()
}

// Set stores a tag in cache
func (r *TagCacheRepository) Set(ctx context.Context, tag *entity.Tag) error {
	key := TagKey(tag.ID())
	data, err := json.Marshal(newTagCacheData(tag))
	if err != nil {
		return err
	}

	return r.client.Set(ctx, key, data, r.client.config.TagTTL).Err()
}

// Delete removes a tag from cache
func (r *TagCacheRepository) Delete(ctx context.Context, id valueobject.TagID) error {
	key := TagKey(id)
	return r.client.Del(ctx, key).Err()
}

// GetByName retrieves a cached tag by name
func (r *TagCacheRepository) GetByName(ctx context.Context, name valueobject.TagName) (*entity.Tag, error) {
	key := TagByNameKey(name)
	data, err := r.client.Get(ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, nil // Cache miss
		}
		return nil, err
	}

	var cached tagCacheData
	if err := json.Unmarshal(data, &cached); err != nil {
		return nil, err
	}

	return cached.ToEntity()
}

// SetByName stores a tag in cache by name key
func (r *TagCacheRepository) SetByName(ctx context.Context, name valueobject.TagName, tag *entity.Tag) error {
	key := TagByNameKey(name)
	data, err := json.Marshal(newTagCacheData(tag))
	if err != nil {
		return err
	}

	return r.client.Set(ctx, key, data, r.client.config.TagTTL).Err()
}

// GetList retrieves a cached list of tags
func (r *TagCacheRepository) GetList(ctx context.Context, params repository.ListParams) (*repository.TagListResult, error) {
	key := TagListKey(params.Page, params.PageSize)
	data, err := r.client.Get(ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, nil // Cache miss
		}
		return nil, err
	}

	var cached tagListCacheData
	if err := json.Unmarshal(data, &cached); err != nil {
		return nil, err
	}

	return cached.ToResult()
}

// SetList stores a list of tags in cache
func (r *TagCacheRepository) SetList(ctx context.Context, params repository.ListParams, result *repository.TagListResult) error {
	key := TagListKey(params.Page, params.PageSize)
	data, err := json.Marshal(newTagListCacheData(result))
	if err != nil {
		return err
	}

	return r.client.Set(ctx, key, data, r.client.config.TagTTL).Err()
}

// DeleteList removes all cached tag lists
func (r *TagCacheRepository) DeleteList(ctx context.Context) error {
	pattern := "tag:list:*"
	return r.deleteByPattern(ctx, pattern)
}

// GetByPostID retrieves cached tags by post ID
func (r *TagCacheRepository) GetByPostID(ctx context.Context, postID valueobject.PostID) ([]*entity.Tag, error) {
	key := TagByPostKey(postID)
	data, err := r.client.Get(ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, nil // Cache miss
		}
		return nil, err
	}

	var cached tagListCacheData
	if err := json.Unmarshal(data, &cached); err != nil {
		return nil, err
	}

	result, err := cached.ToResult()
	if err != nil {
		return nil, err
	}

	return result.Items, nil
}

// SetByPostID stores tags by post ID in cache
func (r *TagCacheRepository) SetByPostID(ctx context.Context, postID valueobject.PostID, tags []*entity.Tag) error {
	key := TagByPostKey(postID)
	result := &repository.TagListResult{
		Total:     int64(len(tags)),
		Page:      1,
		PageSize:  len(tags),
		TotalPage: 1,
		Items:     tags,
	}
	data, err := json.Marshal(newTagListCacheData(result))
	if err != nil {
		return err
	}

	return r.client.Set(ctx, key, data, r.client.config.TagTTL).Err()
}

// InvalidateTag invalidates all cache entries related to a tag
func (r *TagCacheRepository) InvalidateTag(ctx context.Context, id valueobject.TagID) error {
	// Delete the tag itself
	if err := r.Delete(ctx, id); err != nil {
		return err
	}

	// Delete all list caches
	if err := r.DeleteList(ctx); err != nil {
		return err
	}

	// Delete all post-related tag caches
	pattern := "tag:post:*"
	if err := r.deleteByPattern(ctx, pattern); err != nil {
		return err
	}

	// Delete name-based cache
	pattern = "tag:name:*"
	if err := r.deleteByPattern(ctx, pattern); err != nil {
		return err
	}

	return nil
}

// deleteByPattern deletes all keys matching the pattern
func (r *TagCacheRepository) deleteByPattern(ctx context.Context, pattern string) error {
	iter := r.client.Scan(ctx, 0, pattern, 0).Iterator()
	for iter.Next(ctx) {
		if err := r.client.Del(ctx, iter.Val()).Err(); err != nil {
			return err
		}
	}
	return iter.Err()
}

// Cache data structures for JSON serialization

type tagCacheData struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

func newTagCacheData(tag *entity.Tag) *tagCacheData {
	return &tagCacheData{
		ID:        tag.ID().String(),
		Name:      tag.Name().String(),
		CreatedAt: tag.CreatedAt().Time(),
	}
}

func (d *tagCacheData) ToEntity() (*entity.Tag, error) {
	name, err := valueobject.NewTagName(d.Name)
	if err != nil {
		return nil, err
	}

	return entity.NewTagFromPersistence(
		mustParseTagID(d.ID),
		name,
		valueobject.NewCreatedAt(d.CreatedAt),
	), nil
}

type tagListCacheData struct {
	Total     int64           `json:"total"`
	Page      int             `json:"page"`
	PageSize  int             `json:"page_size"`
	TotalPage int             `json:"total_page"`
	Items     []*tagCacheData `json:"items"`
}

func newTagListCacheData(result *repository.TagListResult) *tagListCacheData {
	items := make([]*tagCacheData, len(result.Items))
	for i, tag := range result.Items {
		items[i] = newTagCacheData(tag)
	}

	return &tagListCacheData{
		Total:     result.Total,
		Page:      result.Page,
		PageSize:  result.PageSize,
		TotalPage: result.TotalPage,
		Items:     items,
	}
}

func (d *tagListCacheData) ToResult() (*repository.TagListResult, error) {
	items := make([]*entity.Tag, len(d.Items))
	for i, item := range d.Items {
		tag, err := item.ToEntity()
		if err != nil {
			return nil, err
		}
		items[i] = tag
	}

	return &repository.TagListResult{
		Total:     d.Total,
		Page:      d.Page,
		PageSize:  d.PageSize,
		TotalPage: d.TotalPage,
		Items:     items,
	}, nil
}

// Helper function to parse TagID (panics on error, used only for cached data)
func mustParseTagID(s string) valueobject.TagID {
	id, err := valueobject.NewTagID(s)
	if err != nil {
		panic(err)
	}
	return id
}
