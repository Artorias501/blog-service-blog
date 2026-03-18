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

// CommentRepository implements repository.CommentRepository using SQLite with GORM
type CommentRepository struct {
	db *gorm.DB
}

// NewCommentRepository creates a new CommentRepository instance
func NewCommentRepository(db *gorm.DB) *CommentRepository {
	return &CommentRepository{db: db}
}

// Create inserts a new comment into the database
func (r *CommentRepository) Create(ctx context.Context, comment *entity.Comment) error {
	// Note: PostID must be set separately when creating a comment
	// This is handled at the service layer
	return fmt.Errorf("use CreateWithPostID instead")
}

// CreateWithPostID inserts a new comment with the associated post ID
func (r *CommentRepository) CreateWithPostID(ctx context.Context, comment *entity.Comment, postID valueobject.PostID) error {
	commentModel := converter.CommentToModelWithPostID(comment, postID)
	result := r.db.WithContext(ctx).Create(&commentModel)
	if result.Error != nil {
		return fmt.Errorf("failed to create comment: %w", result.Error)
	}
	return nil
}

// GetByID retrieves a comment by its ID
func (r *CommentRepository) GetByID(ctx context.Context, id valueobject.CommentID) (*entity.Comment, error) {
	var commentModel model.CommentModel
	result := r.db.WithContext(ctx).First(&commentModel, "id = ?", id.String())
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get comment: %w", result.Error)
	}
	return converter.CommentToEntity(&commentModel)
}

// Update updates an existing comment in the database
func (r *CommentRepository) Update(ctx context.Context, comment *entity.Comment) error {
	// Get existing comment to preserve post_id
	var existingModel model.CommentModel
	result := r.db.WithContext(ctx).First(&existingModel, "id = ?", comment.ID().String())
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return fmt.Errorf("comment not found")
		}
		return fmt.Errorf("failed to get existing comment: %w", result.Error)
	}

	// Update fields
	existingModel.AuthorName = comment.AuthorName().String()
	existingModel.AuthorEmail = comment.AuthorEmail().String()
	existingModel.Content = comment.Content().String()
	existingModel.Status = comment.Status().String()

	result = r.db.WithContext(ctx).Save(&existingModel)
	if result.Error != nil {
		return fmt.Errorf("failed to update comment: %w", result.Error)
	}
	return nil
}

// Delete removes a comment from the database
func (r *CommentRepository) Delete(ctx context.Context, id valueobject.CommentID) error {
	result := r.db.WithContext(ctx).Delete(&model.CommentModel{}, "id = ?", id.String())
	if result.Error != nil {
		return fmt.Errorf("failed to delete comment: %w", result.Error)
	}
	return nil
}

// List retrieves comments with pagination and sorting
func (r *CommentRepository) List(ctx context.Context, params repository.ListParams) (*repository.CommentListResult, error) {
	var comments []model.CommentModel
	var total int64

	// Count total
	if err := r.db.WithContext(ctx).Model(&model.CommentModel{}).Count(&total).Error; err != nil {
		return nil, fmt.Errorf("failed to count comments: %w", err)
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
	if err := r.db.WithContext(ctx).Order(orderClause).Offset(offset).Limit(params.PageSize).Find(&comments).Error; err != nil {
		return nil, fmt.Errorf("failed to list comments: %w", err)
	}

	// Convert to entities
	items := make([]*entity.Comment, 0, len(comments))
	for _, commentModel := range comments {
		comment, err := converter.CommentToEntity(&commentModel)
		if err != nil {
			return nil, err
		}
		items = append(items, comment)
	}

	return &repository.CommentListResult{
		Total:     total,
		Page:      params.Page,
		PageSize:  params.PageSize,
		TotalPage: totalPage,
		Items:     items,
	}, nil
}

// ListByPostID retrieves all comments for a specific post
func (r *CommentRepository) ListByPostID(ctx context.Context, postID valueobject.PostID, params repository.ListParams) (*repository.CommentListResult, error) {
	var comments []model.CommentModel
	var total int64

	// Count total comments for the post
	if err := r.db.WithContext(ctx).
		Model(&model.CommentModel{}).
		Where("post_id = ?", postID.String()).
		Count(&total).Error; err != nil {
		return nil, fmt.Errorf("failed to count comments by post: %w", err)
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
		Where("post_id = ?", postID.String()).
		Order(orderClause).
		Offset(offset).
		Limit(params.PageSize).
		Find(&comments).Error; err != nil {
		return nil, fmt.Errorf("failed to list comments by post: %w", err)
	}

	// Convert to entities
	items := make([]*entity.Comment, 0, len(comments))
	for _, commentModel := range comments {
		comment, err := converter.CommentToEntity(&commentModel)
		if err != nil {
			return nil, err
		}
		items = append(items, comment)
	}

	return &repository.CommentListResult{
		Total:     total,
		Page:      params.Page,
		PageSize:  params.PageSize,
		TotalPage: totalPage,
		Items:     items,
	}, nil
}

// ListByStatus retrieves comments by status (pending, approved, rejected, spam)
func (r *CommentRepository) ListByStatus(ctx context.Context, status string, params repository.ListParams) (*repository.CommentListResult, error) {
	var comments []model.CommentModel
	var total int64

	// Count total comments with the status
	if err := r.db.WithContext(ctx).
		Model(&model.CommentModel{}).
		Where("status = ?", status).
		Count(&total).Error; err != nil {
		return nil, fmt.Errorf("failed to count comments by status: %w", err)
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
		Where("status = ?", status).
		Order(orderClause).
		Offset(offset).
		Limit(params.PageSize).
		Find(&comments).Error; err != nil {
		return nil, fmt.Errorf("failed to list comments by status: %w", err)
	}

	// Convert to entities
	items := make([]*entity.Comment, 0, len(comments))
	for _, commentModel := range comments {
		comment, err := converter.CommentToEntity(&commentModel)
		if err != nil {
			return nil, err
		}
		items = append(items, comment)
	}

	return &repository.CommentListResult{
		Total:     total,
		Page:      params.Page,
		PageSize:  params.PageSize,
		TotalPage: totalPage,
		Items:     items,
	}, nil
}

// CountByPostID returns the count of comments for a specific post
func (r *CommentRepository) CountByPostID(ctx context.Context, postID valueobject.PostID) (int64, error) {
	var count int64
	result := r.db.WithContext(ctx).
		Model(&model.CommentModel{}).
		Where("post_id = ?", postID.String()).
		Count(&count)
	if result.Error != nil {
		return 0, fmt.Errorf("failed to count comments: %w", result.Error)
	}
	return count, nil
}

// Approve updates a comment's status to approved
func (r *CommentRepository) Approve(ctx context.Context, id valueobject.CommentID) error {
	result := r.db.WithContext(ctx).
		Model(&model.CommentModel{}).
		Where("id = ?", id.String()).
		Update("status", valueobject.CommentStatusApproved)
	if result.Error != nil {
		return fmt.Errorf("failed to approve comment: %w", result.Error)
	}
	return nil
}

// Reject updates a comment's status to rejected
func (r *CommentRepository) Reject(ctx context.Context, id valueobject.CommentID) error {
	result := r.db.WithContext(ctx).
		Model(&model.CommentModel{}).
		Where("id = ?", id.String()).
		Update("status", valueobject.CommentStatusRejected)
	if result.Error != nil {
		return fmt.Errorf("failed to reject comment: %w", result.Error)
	}
	return nil
}

// MarkAsSpam updates a comment's status to spam
func (r *CommentRepository) MarkAsSpam(ctx context.Context, id valueobject.CommentID) error {
	result := r.db.WithContext(ctx).
		Model(&model.CommentModel{}).
		Where("id = ?", id.String()).
		Update("status", valueobject.CommentStatusSpam)
	if result.Error != nil {
		return fmt.Errorf("failed to mark comment as spam: %w", result.Error)
	}
	return nil
}

// DeleteByPostID removes all comments associated with a post
func (r *CommentRepository) DeleteByPostID(ctx context.Context, postID valueobject.PostID) error {
	result := r.db.WithContext(ctx).
		Delete(&model.CommentModel{}, "post_id = ?", postID.String())
	if result.Error != nil {
		return fmt.Errorf("failed to delete comments by post: %w", result.Error)
	}
	return nil
}
