package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/artorias501/blog-service/internal/domain/entity"
	"github.com/artorias501/blog-service/internal/domain/repository"
	"github.com/artorias501/blog-service/internal/domain/valueobject"
)

// TagService provides business logic operations for tags
type TagService struct {
	tagRepo  repository.TagRepository
	tagCache repository.TagCacheRepository
}

// NewTagService creates a new TagService instance
func NewTagService(
	tagRepo repository.TagRepository,
	tagCache repository.TagCacheRepository,
) *TagService {
	return &TagService{
		tagRepo:  tagRepo,
		tagCache: tagCache,
	}
}

// CreateTagInput contains the data needed to create a new tag
type CreateTagInput struct {
	Name string
}

// CreateTag creates a new tag with the given input
func (s *TagService) CreateTag(ctx context.Context, input CreateTagInput) (*entity.Tag, error) {
	// Validate and create value object
	name, err := valueobject.NewTagName(input.Name)
	if err != nil {
		return nil, fmt.Errorf("invalid tag name: %w", err)
	}

	// Check if tag already exists
	existingTag, err := s.tagRepo.GetByName(ctx, name)
	if err == nil && existingTag != nil {
		return nil, errors.New("tag with this name already exists")
	}

	// Create tag entity
	tag := entity.NewTag(name)

	// Persist the tag
	if err := s.tagRepo.Create(ctx, tag); err != nil {
		return nil, fmt.Errorf("failed to create tag: %w", err)
	}

	// Invalidate list caches
	s.tagCache.DeleteList(ctx)

	return tag, nil
}

// GetTagByID retrieves a tag by ID using cache-aside pattern
func (s *TagService) GetTagByID(ctx context.Context, id string) (*entity.Tag, error) {
	// Validate ID
	tagID, err := valueobject.NewTagID(id)
	if err != nil {
		return nil, fmt.Errorf("invalid tag ID: %w", err)
	}

	// Check cache first (cache-aside pattern)
	cachedTag, err := s.tagCache.Get(ctx, tagID)
	if err != nil {
		// Log cache error but continue to database
	}
	if cachedTag != nil {
		return cachedTag, nil
	}

	// Cache miss - fetch from database
	tag, err := s.tagRepo.GetByID(ctx, tagID)
	if err != nil {
		return nil, fmt.Errorf("failed to get tag: %w", err)
	}
	if tag == nil {
		return nil, errors.New("tag not found")
	}

	// Populate cache
	s.tagCache.Set(ctx, tag)

	return tag, nil
}

// GetTagByName retrieves a tag by name using cache-aside pattern
func (s *TagService) GetTagByName(ctx context.Context, name string) (*entity.Tag, error) {
	// Validate name
	tagName, err := valueobject.NewTagName(name)
	if err != nil {
		return nil, fmt.Errorf("invalid tag name: %w", err)
	}

	// Check cache first
	cachedTag, err := s.tagCache.GetByName(ctx, tagName)
	if err != nil {
		// Log cache error but continue to database
	}
	if cachedTag != nil {
		return cachedTag, nil
	}

	// Cache miss - fetch from database
	tag, err := s.tagRepo.GetByName(ctx, tagName)
	if err != nil {
		return nil, fmt.Errorf("failed to get tag: %w", err)
	}
	if tag == nil {
		return nil, errors.New("tag not found")
	}

	// Populate cache
	s.tagCache.SetByName(ctx, tagName, tag)

	return tag, nil
}

// UpdateTagInput contains the data needed to update a tag
type UpdateTagInput struct {
	Name string
}

// UpdateTag updates an existing tag
func (s *TagService) UpdateTag(ctx context.Context, id string, input UpdateTagInput) (*entity.Tag, error) {
	tagID, err := valueobject.NewTagID(id)
	if err != nil {
		return nil, fmt.Errorf("invalid tag ID: %w", err)
	}

	// Get existing tag
	tag, err := s.tagRepo.GetByID(ctx, tagID)
	if err != nil {
		return nil, fmt.Errorf("failed to get tag: %w", err)
	}
	if tag == nil {
		return nil, errors.New("tag not found")
	}

	// Validate and update name
	name, err := valueobject.NewTagName(input.Name)
	if err != nil {
		return nil, fmt.Errorf("invalid tag name: %w", err)
	}

	// Check if another tag with the same name exists
	existingTag, err := s.tagRepo.GetByName(ctx, name)
	if err == nil && existingTag != nil && !existingTag.ID().Equals(tagID) {
		return nil, errors.New("tag with this name already exists")
	}

	tag.UpdateName(name)

	// Persist changes
	if err := s.tagRepo.Update(ctx, tag); err != nil {
		return nil, fmt.Errorf("failed to update tag: %w", err)
	}

	// Invalidate caches
	s.tagCache.InvalidateTag(ctx, tagID)
	s.tagCache.DeleteList(ctx)

	return tag, nil
}

// DeleteTag deletes a tag by ID
func (s *TagService) DeleteTag(ctx context.Context, id string) error {
	tagID, err := valueobject.NewTagID(id)
	if err != nil {
		return fmt.Errorf("invalid tag ID: %w", err)
	}

	// Check if tag exists
	tag, err := s.tagRepo.GetByID(ctx, tagID)
	if err != nil {
		return fmt.Errorf("failed to get tag: %w", err)
	}
	if tag == nil {
		return errors.New("tag not found")
	}

	// Delete the tag
	if err := s.tagRepo.Delete(ctx, tagID); err != nil {
		return fmt.Errorf("failed to delete tag: %w", err)
	}

	// Invalidate caches
	s.tagCache.InvalidateTag(ctx, tagID)
	s.tagCache.DeleteList(ctx)

	return nil
}

// ListTagsInput contains parameters for listing tags
type ListTagsInput struct {
	Page     int
	PageSize int
	SortBy   string
	Order    string
}

// ListTags retrieves a paginated list of tags
func (s *TagService) ListTags(ctx context.Context, input ListTagsInput) (*repository.TagListResult, error) {
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
	cachedResult, err := s.tagCache.GetList(ctx, params)
	if err == nil && cachedResult != nil {
		return cachedResult, nil
	}

	// Cache miss - fetch from database
	result, err := s.tagRepo.List(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to list tags: %w", err)
	}

	// Populate cache
	s.tagCache.SetList(ctx, params, result)

	return result, nil
}

// SearchTagsInput contains search criteria
type SearchTagsInput struct {
	Keyword  string
	Page     int
	PageSize int
	SortBy   string
	Order    string
}

// SearchTags searches tags by name pattern
func (s *TagService) SearchTags(ctx context.Context, input SearchTagsInput) (*repository.TagListResult, error) {
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

	// Search in repository (no caching for search results)
	result, err := s.tagRepo.Search(ctx, input.Keyword, params)
	if err != nil {
		return nil, fmt.Errorf("failed to search tags: %w", err)
	}

	return result, nil
}

// GetOrCreateTag retrieves a tag by name or creates it if not exists
func (s *TagService) GetOrCreateTag(ctx context.Context, name string) (*entity.Tag, error) {
	tagName, err := valueobject.NewTagName(name)
	if err != nil {
		return nil, fmt.Errorf("invalid tag name: %w", err)
	}

	tag, err := s.tagRepo.GetOrCreate(ctx, tagName)
	if err != nil {
		return nil, fmt.Errorf("failed to get or create tag: %w", err)
	}

	// Invalidate list cache if a new tag was created
	s.tagCache.DeleteList(ctx)

	return tag, nil
}

// GetTagsByPostID retrieves all tags associated with a post
func (s *TagService) GetTagsByPostID(ctx context.Context, postID string) ([]*entity.Tag, error) {
	pid, err := valueobject.NewPostID(postID)
	if err != nil {
		return nil, fmt.Errorf("invalid post ID: %w", err)
	}

	// Check cache first
	cachedTags, err := s.tagCache.GetByPostID(ctx, pid)
	if err == nil && cachedTags != nil {
		return cachedTags, nil
	}

	// Cache miss - fetch from database
	tags, err := s.tagRepo.ListByPostID(ctx, pid)
	if err != nil {
		return nil, fmt.Errorf("failed to get tags by post ID: %w", err)
	}

	// Populate cache
	if len(tags) > 0 {
		s.tagCache.SetByPostID(ctx, pid, tags)
	}

	return tags, nil
}
