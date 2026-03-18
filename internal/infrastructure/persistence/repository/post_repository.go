package repository

import (
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"

	"github.com/artorias501/blog-service/internal/domain/entity"
	"github.com/artorias501/blog-service/internal/domain/repository"
	"github.com/artorias501/blog-service/internal/domain/valueobject"
	"github.com/artorias501/blog-service/internal/infrastructure/persistence/converter"
	"github.com/artorias501/blog-service/internal/infrastructure/persistence/model"
)

// PostRepository implements repository.PostRepository using SQLite with GORM
type PostRepository struct {
	db *gorm.DB
}

// NewPostRepository creates a new PostRepository instance
func NewPostRepository(db *gorm.DB) *PostRepository {
	return &PostRepository{db: db}
}

// Create inserts a new post into the database
func (r *PostRepository) Create(ctx context.Context, post *entity.Post) error {
	postModel := converter.PostToModel(post)
	result := r.db.WithContext(ctx).Create(postModel)
	if result.Error != nil {
		return fmt.Errorf("failed to create post: %w", result.Error)
	}
	return nil
}

// GetByID retrieves a post by its ID
func (r *PostRepository) GetByID(ctx context.Context, id valueobject.PostID) (*entity.Post, error) {
	var postModel model.PostModel
	result := r.db.WithContext(ctx).First(&postModel, "id = ?", id.String())
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get post: %w", result.Error)
	}
	return converter.PostToEntityWithoutRelations(&postModel)
}

// Update updates an existing post in the database
func (r *PostRepository) Update(ctx context.Context, post *entity.Post) error {
	postModel := converter.PostToModel(post)
	result := r.db.WithContext(ctx).Save(postModel)
	if result.Error != nil {
		return fmt.Errorf("failed to update post: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("post not found")
	}
	return nil
}

// Delete removes a post from the database
func (r *PostRepository) Delete(ctx context.Context, id valueobject.PostID) error {
	result := r.db.WithContext(ctx).Delete(&model.PostModel{}, "id = ?", id.String())
	if result.Error != nil {
		return fmt.Errorf("failed to delete post: %w", result.Error)
	}
	return nil
}

// List retrieves posts with pagination and sorting
func (r *PostRepository) List(ctx context.Context, params repository.ListParams) (*repository.PostListResult, error) {
	var posts []model.PostModel
	var total int64

	// Count total
	if err := r.db.WithContext(ctx).Model(&model.PostModel{}).Count(&total).Error; err != nil {
		return nil, fmt.Errorf("failed to count posts: %w", err)
	}

	// Calculate pagination
	offset := (params.Page - 1) * params.PageSize
	totalPage := int(total) / params.PageSize
	if int(total)%params.PageSize > 0 {
		totalPage++
	}

	// Build order clause
	orderClause := "created_at DESC"
	if params.SortBy != "" {
		order := "ASC"
		if params.Order != "" {
			order = params.Order
		}
		orderClause = fmt.Sprintf("%s %s", params.SortBy, order)
	}

	// Query with pagination
	if err := r.db.WithContext(ctx).Order(orderClause).Offset(offset).Limit(params.PageSize).Find(&posts).Error; err != nil {
		return nil, fmt.Errorf("failed to list posts: %w", err)
	}

	// Convert to entities
	items := make([]*entity.Post, 0, len(posts))
	for _, postModel := range posts {
		post, err := converter.PostToEntityWithoutRelations(&postModel)
		if err != nil {
			return nil, err
		}
		items = append(items, post)
	}

	return &repository.PostListResult{
		Total:     total,
		Page:      params.Page,
		PageSize:  params.PageSize,
		TotalPage: totalPage,
		Items:     items,
	}, nil
}

// ListByTagID retrieves posts associated with a specific tag
func (r *PostRepository) ListByTagID(ctx context.Context, tagID valueobject.TagID, params repository.ListParams) (*repository.PostListResult, error) {
	var posts []model.PostModel
	var total int64

	// Count total posts with the tag
	if err := r.db.WithContext(ctx).
		Model(&model.PostModel{}).
		Joins("JOIN post_tags ON post_tags.post_id = posts.id").
		Where("post_tags.tag_id = ?", tagID.String()).
		Count(&total).Error; err != nil {
		return nil, fmt.Errorf("failed to count posts by tag: %w", err)
	}

	// Calculate pagination
	offset := (params.Page - 1) * params.PageSize
	totalPage := int(total) / params.PageSize
	if int(total)%params.PageSize > 0 {
		totalPage++
	}

	// Build order clause
	orderClause := "posts.created_at DESC"
	if params.SortBy != "" {
		order := "ASC"
		if params.Order != "" {
			order = params.Order
		}
		orderClause = fmt.Sprintf("posts.%s %s", params.SortBy, order)
	}

	// Query with pagination
	if err := r.db.WithContext(ctx).
		Joins("JOIN post_tags ON post_tags.post_id = posts.id").
		Where("post_tags.tag_id = ?", tagID.String()).
		Order(orderClause).
		Offset(offset).
		Limit(params.PageSize).
		Find(&posts).Error; err != nil {
		return nil, fmt.Errorf("failed to list posts by tag: %w", err)
	}

	// Convert to entities
	items := make([]*entity.Post, 0, len(posts))
	for _, postModel := range posts {
		post, err := converter.PostToEntityWithoutRelations(&postModel)
		if err != nil {
			return nil, err
		}
		items = append(items, post)
	}

	return &repository.PostListResult{
		Total:     total,
		Page:      params.Page,
		PageSize:  params.PageSize,
		TotalPage: totalPage,
		Items:     items,
	}, nil
}

// GetByIDWithComments retrieves a post with its comments loaded
func (r *PostRepository) GetByIDWithComments(ctx context.Context, id valueobject.PostID) (*entity.Post, error) {
	var postModel model.PostModel
	result := r.db.WithContext(ctx).
		Preload("Comments").
		First(&postModel, "id = ?", id.String())
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get post with comments: %w", result.Error)
	}
	return converter.PostToEntity(&postModel)
}

// GetByIDWithTags retrieves a post with its tags loaded
func (r *PostRepository) GetByIDWithTags(ctx context.Context, id valueobject.PostID) (*entity.Post, error) {
	var postModel model.PostModel
	result := r.db.WithContext(ctx).
		Preload("Tags").
		First(&postModel, "id = ?", id.String())
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get post with tags: %w", result.Error)
	}
	return converter.PostToEntity(&postModel)
}

// GetByIDFull retrieves a post with both comments and tags loaded
func (r *PostRepository) GetByIDFull(ctx context.Context, id valueobject.PostID) (*entity.Post, error) {
	var postModel model.PostModel
	result := r.db.WithContext(ctx).
		Preload("Comments").
		Preload("Tags").
		First(&postModel, "id = ?", id.String())
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get post full: %w", result.Error)
	}
	return converter.PostToEntity(&postModel)
}

// AddTag associates a tag with a post
func (r *PostRepository) AddTag(ctx context.Context, postID valueobject.PostID, tagID valueobject.TagID) error {
	postTag := model.PostTagModel{
		PostID: postID.String(),
		TagID:  tagID.String(),
	}
	result := r.db.WithContext(ctx).Create(&postTag)
	if result.Error != nil {
		// Ignore duplicate key error
		if !errors.Is(result.Error, gorm.ErrDuplicatedKey) {
			return fmt.Errorf("failed to add tag to post: %w", result.Error)
		}
	}
	return nil
}

// RemoveTag disassociates a tag from a post
func (r *PostRepository) RemoveTag(ctx context.Context, postID valueobject.PostID, tagID valueobject.TagID) error {
	result := r.db.WithContext(ctx).
		Delete(&model.PostTagModel{}, "post_id = ? AND tag_id = ?", postID.String(), tagID.String())
	if result.Error != nil {
		return fmt.Errorf("failed to remove tag from post: %w", result.Error)
	}
	return nil
}

// Search searches posts by title or content keyword
func (r *PostRepository) Search(ctx context.Context, keyword string, params repository.ListParams) (*repository.PostListResult, error) {
	var posts []model.PostModel
	var total int64

	searchPattern := "%" + keyword + "%"

	// Count total matching posts
	if err := r.db.WithContext(ctx).
		Model(&model.PostModel{}).
		Where("title LIKE ? OR content LIKE ?", searchPattern, searchPattern).
		Count(&total).Error; err != nil {
		return nil, fmt.Errorf("failed to count search results: %w", err)
	}

	// Calculate pagination
	offset := (params.Page - 1) * params.PageSize
	totalPage := int(total) / params.PageSize
	if int(total)%params.PageSize > 0 {
		totalPage++
	}

	// Build order clause
	orderClause := "created_at DESC"
	if params.SortBy != "" {
		order := "ASC"
		if params.Order != "" {
			order = params.Order
		}
		orderClause = fmt.Sprintf("%s %s", params.SortBy, order)
	}

	// Query with pagination
	if err := r.db.WithContext(ctx).
		Where("title LIKE ? OR content LIKE ?", searchPattern, searchPattern).
		Order(orderClause).
		Offset(offset).
		Limit(params.PageSize).
		Find(&posts).Error; err != nil {
		return nil, fmt.Errorf("failed to search posts: %w", err)
	}

	// Convert to entities
	items := make([]*entity.Post, 0, len(posts))
	for _, postModel := range posts {
		post, err := converter.PostToEntityWithoutRelations(&postModel)
		if err != nil {
			return nil, err
		}
		items = append(items, post)
	}

	return &repository.PostListResult{
		Total:     total,
		Page:      params.Page,
		PageSize:  params.PageSize,
		TotalPage: totalPage,
		Items:     items,
	}, nil
}
