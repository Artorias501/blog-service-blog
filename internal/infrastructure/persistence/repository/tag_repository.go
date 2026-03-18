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

// TagRepository implements repository.TagRepository using SQLite with GORM
type TagRepository struct {
	db *gorm.DB
}

// NewTagRepository creates a new TagRepository instance
func NewTagRepository(db *gorm.DB) *TagRepository {
	return &TagRepository{db: db}
}

// Create inserts a new tag into the database
func (r *TagRepository) Create(ctx context.Context, tag *entity.Tag) error {
	tagModel := converter.TagToModel(tag)
	result := r.db.WithContext(ctx).Create(tagModel)
	if result.Error != nil {
		return fmt.Errorf("failed to create tag: %w", result.Error)
	}
	return nil
}

// GetByID retrieves a tag by its ID
func (r *TagRepository) GetByID(ctx context.Context, id valueobject.TagID) (*entity.Tag, error) {
	var tagModel model.TagModel
	result := r.db.WithContext(ctx).First(&tagModel, "id = ?", id.String())
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get tag: %w", result.Error)
	}
	return converter.TagToEntity(&tagModel)
}

// GetByName retrieves a tag by its name
func (r *TagRepository) GetByName(ctx context.Context, name valueobject.TagName) (*entity.Tag, error) {
	var tagModel model.TagModel
	result := r.db.WithContext(ctx).First(&tagModel, "name = ?", name.String())
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get tag by name: %w", result.Error)
	}
	return converter.TagToEntity(&tagModel)
}

// Update updates an existing tag in the database
func (r *TagRepository) Update(ctx context.Context, tag *entity.Tag) error {
	tagModel := converter.TagToModel(tag)
	result := r.db.WithContext(ctx).Save(tagModel)
	if result.Error != nil {
		return fmt.Errorf("failed to update tag: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("tag not found")
	}
	return nil
}

// Delete removes a tag from the database
func (r *TagRepository) Delete(ctx context.Context, id valueobject.TagID) error {
	// Start a transaction to delete tag and its associations
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Delete post_tags associations
		if err := tx.Delete(&model.PostTagModel{}, "tag_id = ?", id.String()).Error; err != nil {
			return fmt.Errorf("failed to delete tag associations: %w", err)
		}

		// Delete the tag
		if err := tx.Delete(&model.TagModel{}, "id = ?", id.String()).Error; err != nil {
			return fmt.Errorf("failed to delete tag: %w", err)
		}

		return nil
	})
}

// List retrieves tags with pagination and sorting
func (r *TagRepository) List(ctx context.Context, params repository.ListParams) (*repository.TagListResult, error) {
	var tags []model.TagModel
	var total int64

	// Count total
	if err := r.db.WithContext(ctx).Model(&model.TagModel{}).Count(&total).Error; err != nil {
		return nil, fmt.Errorf("failed to count tags: %w", err)
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
	if err := r.db.WithContext(ctx).Order(orderClause).Offset(offset).Limit(params.PageSize).Find(&tags).Error; err != nil {
		return nil, fmt.Errorf("failed to list tags: %w", err)
	}

	// Convert to entities
	items := make([]*entity.Tag, 0, len(tags))
	for _, tagModel := range tags {
		tag, err := converter.TagToEntity(&tagModel)
		if err != nil {
			return nil, err
		}
		items = append(items, tag)
	}

	return &repository.TagListResult{
		Total:     total,
		Page:      params.Page,
		PageSize:  params.PageSize,
		TotalPage: totalPage,
		Items:     items,
	}, nil
}

// ListByPostID retrieves all tags associated with a post
func (r *TagRepository) ListByPostID(ctx context.Context, postID valueobject.PostID) ([]*entity.Tag, error) {
	var tags []model.TagModel
	result := r.db.WithContext(ctx).
		Joins("JOIN post_tags ON post_tags.tag_id = tags.id").
		Where("post_tags.post_id = ?", postID.String()).
		Find(&tags)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to list tags by post: %w", result.Error)
	}

	// Convert to entities
	items := make([]*entity.Tag, 0, len(tags))
	for _, tagModel := range tags {
		tag, err := converter.TagToEntity(&tagModel)
		if err != nil {
			return nil, err
		}
		items = append(items, tag)
	}

	return items, nil
}

// GetOrCreate retrieves a tag by name or creates it if not exists
func (r *TagRepository) GetOrCreate(ctx context.Context, name valueobject.TagName) (*entity.Tag, error) {
	// Try to get existing tag
	tag, err := r.GetByName(ctx, name)
	if err != nil {
		return nil, err
	}
	if tag != nil {
		return tag, nil
	}

	// Create new tag
	newTag := entity.NewTag(name)
	if err := r.Create(ctx, newTag); err != nil {
		return nil, fmt.Errorf("failed to create tag: %w", err)
	}

	return newTag, nil
}

// Search searches tags by name pattern
func (r *TagRepository) Search(ctx context.Context, keyword string, params repository.ListParams) (*repository.TagListResult, error) {
	var tags []model.TagModel
	var total int64

	searchPattern := "%" + keyword + "%"

	// Count total matching tags
	if err := r.db.WithContext(ctx).
		Model(&model.TagModel{}).
		Where("name LIKE ?", searchPattern).
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
		Where("name LIKE ?", searchPattern).
		Order(orderClause).
		Offset(offset).
		Limit(params.PageSize).
		Find(&tags).Error; err != nil {
		return nil, fmt.Errorf("failed to search tags: %w", err)
	}

	// Convert to entities
	items := make([]*entity.Tag, 0, len(tags))
	for _, tagModel := range tags {
		tag, err := converter.TagToEntity(&tagModel)
		if err != nil {
			return nil, err
		}
		items = append(items, tag)
	}

	return &repository.TagListResult{
		Total:     total,
		Page:      params.Page,
		PageSize:  params.PageSize,
		TotalPage: totalPage,
		Items:     items,
	}, nil
}
