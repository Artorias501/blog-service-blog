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

// PostCacheRepository implements repository.PostCacheRepository using Redis
type PostCacheRepository struct {
	client *RedisClient
}

// NewPostCacheRepository creates a new PostCacheRepository
func NewPostCacheRepository(client *RedisClient) *PostCacheRepository {
	return &PostCacheRepository{
		client: client,
	}
}

// Get retrieves a cached post by ID
func (r *PostCacheRepository) Get(ctx context.Context, id valueobject.PostID) (*entity.Post, error) {
	key := PostKey(id)
	data, err := r.client.Get(ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, nil // Cache miss
		}
		return nil, err
	}

	var cached postCacheData
	if err := json.Unmarshal(data, &cached); err != nil {
		return nil, err
	}

	return cached.ToEntity()
}

// Set stores a post in cache
func (r *PostCacheRepository) Set(ctx context.Context, post *entity.Post) error {
	key := PostKey(post.ID())
	data, err := json.Marshal(newPostCacheData(post))
	if err != nil {
		return err
	}

	return r.client.Set(ctx, key, data, r.client.config.PostTTL).Err()
}

// Delete removes a post from cache
func (r *PostCacheRepository) Delete(ctx context.Context, id valueobject.PostID) error {
	key := PostKey(id)
	return r.client.Del(ctx, key).Err()
}

// GetList retrieves a cached list of posts
func (r *PostCacheRepository) GetList(ctx context.Context, params repository.ListParams) (*repository.PostListResult, error) {
	key := PostListKey(params.Page, params.PageSize)
	data, err := r.client.Get(ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, nil // Cache miss
		}
		return nil, err
	}

	var cached postListCacheData
	if err := json.Unmarshal(data, &cached); err != nil {
		return nil, err
	}

	return cached.ToResult()
}

// SetList stores a list of posts in cache
func (r *PostCacheRepository) SetList(ctx context.Context, params repository.ListParams, result *repository.PostListResult) error {
	key := PostListKey(params.Page, params.PageSize)
	data, err := json.Marshal(newPostListCacheData(result))
	if err != nil {
		return err
	}

	return r.client.Set(ctx, key, data, r.client.config.PostListTTL).Err()
}

// DeleteList removes all cached post lists
func (r *PostCacheRepository) DeleteList(ctx context.Context) error {
	pattern := "post:list:*"
	return r.deleteByPattern(ctx, pattern)
}

// GetByTagID retrieves cached posts by tag ID
func (r *PostCacheRepository) GetByTagID(ctx context.Context, tagID valueobject.TagID, params repository.ListParams) (*repository.PostListResult, error) {
	key := PostByTagKey(tagID, params.Page, params.PageSize)
	data, err := r.client.Get(ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, nil // Cache miss
		}
		return nil, err
	}

	var cached postListCacheData
	if err := json.Unmarshal(data, &cached); err != nil {
		return nil, err
	}

	return cached.ToResult()
}

// SetByTagID stores posts by tag ID in cache
func (r *PostCacheRepository) SetByTagID(ctx context.Context, tagID valueobject.TagID, params repository.ListParams, result *repository.PostListResult) error {
	key := PostByTagKey(tagID, params.Page, params.PageSize)
	data, err := json.Marshal(newPostListCacheData(result))
	if err != nil {
		return err
	}

	return r.client.Set(ctx, key, data, r.client.config.PostListTTL).Err()
}

// InvalidatePost invalidates all cache entries related to a post
func (r *PostCacheRepository) InvalidatePost(ctx context.Context, id valueobject.PostID) error {
	// Delete the post itself
	if err := r.Delete(ctx, id); err != nil {
		return err
	}

	// Delete all list caches
	if err := r.DeleteList(ctx); err != nil {
		return err
	}

	// Delete all tag-related post caches
	pattern := "post:tag:*"
	if err := r.deleteByPattern(ctx, pattern); err != nil {
		return err
	}

	return nil
}

// deleteByPattern deletes all keys matching the pattern
func (r *PostCacheRepository) deleteByPattern(ctx context.Context, pattern string) error {
	iter := r.client.Scan(ctx, 0, pattern, 0).Iterator()
	for iter.Next(ctx) {
		if err := r.client.Del(ctx, iter.Val()).Err(); err != nil {
			return err
		}
	}
	return iter.Err()
}

// Cache data structures for JSON serialization

type postCacheData struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Content     string    `json:"content"`
	Summary     *string   `json:"summary"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	PublishedAt *string   `json:"published_at"`
}

func newPostCacheData(post *entity.Post) *postCacheData {
	data := &postCacheData{
		ID:        post.ID().String(),
		Title:     post.Title().String(),
		Content:   post.Content().String(),
		Summary:   nil,
		CreatedAt: post.CreatedAt().Time(),
		UpdatedAt: post.UpdatedAt().Time(),
	}

	if post.Summary() != nil {
		s := post.Summary().String()
		data.Summary = &s
	}

	if post.PublishedAt() != nil {
		pa := post.PublishedAt().Format(time.RFC3339)
		data.PublishedAt = &pa
	}

	return data
}

func (d *postCacheData) ToEntity() (*entity.Post, error) {
	title, err := valueobject.NewTitle(d.Title)
	if err != nil {
		return nil, err
	}

	content, err := valueobject.NewContent(d.Content)
	if err != nil {
		return nil, err
	}

	post := entity.NewPost(title, content)
	post.SetPostID(mustParsePostID(d.ID))
	post.SetPostTimestamps(
		valueobject.NewCreatedAt(d.CreatedAt),
		valueobject.NewUpdatedAt(d.UpdatedAt),
	)

	if d.Summary != nil {
		summary, err := valueobject.NewSummary(*d.Summary)
		if err != nil {
			return nil, err
		}
		post.SetSummary(summary)
	}

	return post, nil
}

type postListCacheData struct {
	Total     int64            `json:"total"`
	Page      int              `json:"page"`
	PageSize  int              `json:"page_size"`
	TotalPage int              `json:"total_page"`
	Items     []*postCacheData `json:"items"`
}

func newPostListCacheData(result *repository.PostListResult) *postListCacheData {
	items := make([]*postCacheData, len(result.Items))
	for i, post := range result.Items {
		items[i] = newPostCacheData(post)
	}

	return &postListCacheData{
		Total:     result.Total,
		Page:      result.Page,
		PageSize:  result.PageSize,
		TotalPage: result.TotalPage,
		Items:     items,
	}
}

func (d *postListCacheData) ToResult() (*repository.PostListResult, error) {
	items := make([]*entity.Post, len(d.Items))
	for i, item := range d.Items {
		post, err := item.ToEntity()
		if err != nil {
			return nil, err
		}
		items[i] = post
	}

	return &repository.PostListResult{
		Total:     d.Total,
		Page:      d.Page,
		PageSize:  d.PageSize,
		TotalPage: d.TotalPage,
		Items:     items,
	}, nil
}

// Helper function to parse PostID (panics on error, used only for cached data)
func mustParsePostID(s string) valueobject.PostID {
	id, err := valueobject.NewPostID(s)
	if err != nil {
		panic(err)
	}
	return id
}
