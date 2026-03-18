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

// CommentCacheRepository implements repository.CommentCacheRepository using Redis
type CommentCacheRepository struct {
	client *RedisClient
}

// NewCommentCacheRepository creates a new CommentCacheRepository
func NewCommentCacheRepository(client *RedisClient) *CommentCacheRepository {
	return &CommentCacheRepository{
		client: client,
	}
}

// Get retrieves a cached comment by ID
func (r *CommentCacheRepository) Get(ctx context.Context, id valueobject.CommentID) (*entity.Comment, error) {
	key := CommentKey(id)
	data, err := r.client.Get(ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, nil // Cache miss
		}
		return nil, err
	}

	var cached commentCacheData
	if err := json.Unmarshal(data, &cached); err != nil {
		return nil, err
	}

	return cached.ToEntity()
}

// Set stores a comment in cache
func (r *CommentCacheRepository) Set(ctx context.Context, comment *entity.Comment) error {
	key := CommentKey(comment.ID())
	data, err := json.Marshal(newCommentCacheData(comment))
	if err != nil {
		return err
	}

	return r.client.Set(ctx, key, data, r.client.config.CommentTTL).Err()
}

// Delete removes a comment from cache
func (r *CommentCacheRepository) Delete(ctx context.Context, id valueobject.CommentID) error {
	key := CommentKey(id)
	return r.client.Del(ctx, key).Err()
}

// GetListByPostID retrieves cached comments by post ID
func (r *CommentCacheRepository) GetListByPostID(ctx context.Context, postID valueobject.PostID, params repository.ListParams) (*repository.CommentListResult, error) {
	key := CommentByPostKey(postID, params.Page, params.PageSize)
	data, err := r.client.Get(ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, nil // Cache miss
		}
		return nil, err
	}

	var cached commentListCacheData
	if err := json.Unmarshal(data, &cached); err != nil {
		return nil, err
	}

	return cached.ToResult()
}

// SetListByPostID stores comments by post ID in cache
func (r *CommentCacheRepository) SetListByPostID(ctx context.Context, postID valueobject.PostID, params repository.ListParams, result *repository.CommentListResult) error {
	key := CommentByPostKey(postID, params.Page, params.PageSize)
	data, err := json.Marshal(newCommentListCacheData(result))
	if err != nil {
		return err
	}

	return r.client.Set(ctx, key, data, r.client.config.CommentTTL).Err()
}

// DeleteListByPostID removes cached comments for a post
func (r *CommentCacheRepository) DeleteListByPostID(ctx context.Context, postID valueobject.PostID) error {
	pattern := CommentByPostPattern(postID)
	return r.deleteByPattern(ctx, pattern)
}

// GetCountByPostID retrieves cached comment count for a post
func (r *CommentCacheRepository) GetCountByPostID(ctx context.Context, postID valueobject.PostID) (int64, error) {
	key := CommentCountByPostKey(postID)
	val, err := r.client.Get(ctx, key).Int64()
	if err != nil {
		if err == redis.Nil {
			return -1, nil // Cache miss, return -1 to indicate miss
		}
		return -1, err
	}
	return val, nil
}

// SetCountByPostID stores comment count for a post in cache
func (r *CommentCacheRepository) SetCountByPostID(ctx context.Context, postID valueobject.PostID, count int64) error {
	key := CommentCountByPostKey(postID)
	return r.client.Set(ctx, key, count, r.client.config.CommentCountTTL).Err()
}

// InvalidateComment invalidates all cache entries related to a comment
func (r *CommentCacheRepository) InvalidateComment(ctx context.Context, id valueobject.CommentID) error {
	// Delete the comment itself
	return r.Delete(ctx, id)
}

// InvalidateByPostID invalidates all comment cache entries for a post
func (r *CommentCacheRepository) InvalidateByPostID(ctx context.Context, postID valueobject.PostID) error {
	// Delete all comment lists for this post
	if err := r.DeleteListByPostID(ctx, postID); err != nil {
		return err
	}

	// Delete comment count
	key := CommentCountByPostKey(postID)
	if err := r.client.Del(ctx, key).Err(); err != nil {
		return err
	}

	return nil
}

// deleteByPattern deletes all keys matching the pattern
func (r *CommentCacheRepository) deleteByPattern(ctx context.Context, pattern string) error {
	iter := r.client.Scan(ctx, 0, pattern, 0).Iterator()
	for iter.Next(ctx) {
		if err := r.client.Del(ctx, iter.Val()).Err(); err != nil {
			return err
		}
	}
	return iter.Err()
}

// Cache data structures for JSON serialization

type commentCacheData struct {
	ID          string    `json:"id"`
	AuthorName  string    `json:"author_name"`
	AuthorEmail string    `json:"author_email"`
	Content     string    `json:"content"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
}

func newCommentCacheData(comment *entity.Comment) *commentCacheData {
	return &commentCacheData{
		ID:          comment.ID().String(),
		AuthorName:  comment.AuthorName().String(),
		AuthorEmail: comment.AuthorEmail().String(),
		Content:     comment.Content().String(),
		Status:      comment.Status().String(),
		CreatedAt:   comment.CreatedAt().Time(),
	}
}

func (d *commentCacheData) ToEntity() (*entity.Comment, error) {
	authorName, err := valueobject.NewAuthorName(d.AuthorName)
	if err != nil {
		return nil, err
	}

	authorEmail, err := valueobject.NewAuthorEmail(d.AuthorEmail)
	if err != nil {
		return nil, err
	}

	content, err := valueobject.NewContent(d.Content)
	if err != nil {
		return nil, err
	}

	status, err := valueobject.NewCommentStatus(d.Status)
	if err != nil {
		return nil, err
	}

	return entity.NewCommentFromPersistence(
		mustParseCommentID(d.ID),
		authorName,
		authorEmail,
		content,
		status,
		valueobject.NewCreatedAt(d.CreatedAt),
	), nil
}

type commentListCacheData struct {
	Total     int64               `json:"total"`
	Page      int                 `json:"page"`
	PageSize  int                 `json:"page_size"`
	TotalPage int                 `json:"total_page"`
	Items     []*commentCacheData `json:"items"`
}

func newCommentListCacheData(result *repository.CommentListResult) *commentListCacheData {
	items := make([]*commentCacheData, len(result.Items))
	for i, comment := range result.Items {
		items[i] = newCommentCacheData(comment)
	}

	return &commentListCacheData{
		Total:     result.Total,
		Page:      result.Page,
		PageSize:  result.PageSize,
		TotalPage: result.TotalPage,
		Items:     items,
	}
}

func (d *commentListCacheData) ToResult() (*repository.CommentListResult, error) {
	items := make([]*entity.Comment, len(d.Items))
	for i, item := range d.Items {
		comment, err := item.ToEntity()
		if err != nil {
			return nil, err
		}
		items[i] = comment
	}

	return &repository.CommentListResult{
		Total:     d.Total,
		Page:      d.Page,
		PageSize:  d.PageSize,
		TotalPage: d.TotalPage,
		Items:     items,
	}, nil
}

// Helper function to parse CommentID (panics on error, used only for cached data)
func mustParseCommentID(s string) valueobject.CommentID {
	id, err := valueobject.NewCommentID(s)
	if err != nil {
		panic(err)
	}
	return id
}
