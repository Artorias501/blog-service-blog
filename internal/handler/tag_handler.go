// Package handler contains HTTP handlers for the blog service.
package handler

import (
	"github.com/artorias501/blog-service/internal/domain/entity"
	"github.com/artorias501/blog-service/internal/handler/dto"
	"github.com/artorias501/blog-service/internal/service"
	"github.com/artorias501/blog-service/pkg/response"
	"github.com/gin-gonic/gin"
)

// TagHandler handles HTTP requests for tag operations.
type TagHandler struct {
	tagService *service.TagService
}

// NewTagHandler creates a new TagHandler instance.
func NewTagHandler(tagService *service.TagService) *TagHandler {
	return &TagHandler{
		tagService: tagService,
	}
}

// CreateTag handles POST /api/tags
// @Summary Create a new tag
// @Description Create a new tag with a unique name
// @Tags tags
// @Accept json
// @Produce json
// @Param request body dto.CreateTagRequest true "Tag creation request"
// @Success 201 {object} response.Response{data=dto.TagResponse}
// @Failure 400 {object} response.Response
// @Router /api/tags [post]
func (h *TagHandler) CreateTag(c *gin.Context) {
	var req dto.CreateTagRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "invalid request body: "+err.Error())
		return
	}

	input := service.CreateTagInput{
		Name: req.Name,
	}

	tag, err := h.tagService.CreateTag(c.Request.Context(), input)
	if err != nil {
		handleServiceError(c, err)
		return
	}

	response.Created(c, convertTagToResponse(tag))
}

// GetTag handles GET /api/tags/:id
// @Summary Get a tag by ID
// @Description Retrieve a single tag by its ID
// @Tags tags
// @Produce json
// @Param id path string true "Tag ID"
// @Success 200 {object} response.Response{data=dto.TagResponse}
// @Failure 404 {object} response.Response
// @Router /api/tags/{id} [get]
func (h *TagHandler) GetTag(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.BadRequest(c, "tag ID is required")
		return
	}

	tag, err := h.tagService.GetTagByID(c.Request.Context(), id)
	if err != nil {
		handleServiceError(c, err)
		return
	}

	response.Success(c, convertTagToResponse(tag))
}

// UpdateTag handles PUT /api/tags/:id
// @Summary Update a tag
// @Description Update an existing tag's name
// @Tags tags
// @Accept json
// @Produce json
// @Param id path string true "Tag ID"
// @Param request body dto.UpdateTagRequest true "Tag update request"
// @Success 200 {object} response.Response{data=dto.TagResponse}
// @Failure 400,404 {object} response.Response
// @Router /api/tags/{id} [put]
func (h *TagHandler) UpdateTag(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.BadRequest(c, "tag ID is required")
		return
	}

	var req dto.UpdateTagRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "invalid request body: "+err.Error())
		return
	}

	input := service.UpdateTagInput{
		Name: req.Name,
	}

	tag, err := h.tagService.UpdateTag(c.Request.Context(), id, input)
	if err != nil {
		handleServiceError(c, err)
		return
	}

	response.Success(c, convertTagToResponse(tag))
}

// DeleteTag handles DELETE /api/tags/:id
// @Summary Delete a tag
// @Description Delete a tag by ID
// @Tags tags
// @Param id path string true "Tag ID"
// @Success 204 "No Content"
// @Failure 404 {object} response.Response
// @Router /api/tags/{id} [delete]
func (h *TagHandler) DeleteTag(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.BadRequest(c, "tag ID is required")
		return
	}

	err := h.tagService.DeleteTag(c.Request.Context(), id)
	if err != nil {
		handleServiceError(c, err)
		return
	}

	response.NoContent(c)
}

// ListTags handles GET /api/tags
// @Summary List tags
// @Description Retrieve a paginated list of tags
// @Tags tags
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param size query int false "Page size" default(10)
// @Param sort_by query string false "Sort field"
// @Param order query string false "Sort order"
// @Success 200 {object} response.Response{data=dto.TagListResponse}
// @Router /api/tags [get]
func (h *TagHandler) ListTags(c *gin.Context) {
	var req dto.ListTagsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.BadRequest(c, "invalid query parameters: "+err.Error())
		return
	}

	// Set defaults
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.Size <= 0 {
		req.Size = 10
	}

	input := service.ListTagsInput{
		Page:     req.Page,
		PageSize: req.Size,
		SortBy:   req.SortBy,
		Order:    req.Order,
	}

	result, err := h.tagService.ListTags(c.Request.Context(), input)
	if err != nil {
		handleServiceError(c, err)
		return
	}

	response.Paginated(c, convertTagsToResponse(result.Items), result.Total, result.Page, result.PageSize, result.TotalPage)
}

// SearchTags handles GET /api/tags/search
// @Summary Search tags
// @Description Search tags by name pattern
// @Tags tags
// @Produce json
// @Param keyword query string true "Search keyword"
// @Param page query int false "Page number" default(1)
// @Param size query int false "Page size" default(10)
// @Success 200 {object} response.Response{data=dto.TagListResponse}
// @Router /api/tags/search [get]
func (h *TagHandler) SearchTags(c *gin.Context) {
	var req dto.SearchTagsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.BadRequest(c, "invalid query parameters: "+err.Error())
		return
	}

	// Set defaults
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.Size <= 0 {
		req.Size = 10
	}

	input := service.SearchTagsInput{
		Keyword:  req.Keyword,
		Page:     req.Page,
		PageSize: req.Size,
		SortBy:   req.SortBy,
		Order:    req.Order,
	}

	result, err := h.tagService.SearchTags(c.Request.Context(), input)
	if err != nil {
		handleServiceError(c, err)
		return
	}

	response.Paginated(c, convertTagsToResponse(result.Items), result.Total, result.Page, result.PageSize, result.TotalPage)
}

// GetTagsByPost handles GET /api/posts/:id/tags
// @Summary Get tags by post
// @Description Retrieve all tags associated with a post
// @Tags tags
// @Produce json
// @Param id path string true "Post ID"
// @Success 200 {object} response.Response{data=[]dto.TagResponse}
// @Failure 404 {object} response.Response
// @Router /api/posts/{id}/tags [get]
func (h *TagHandler) GetTagsByPost(c *gin.Context) {
	postID := c.Param("id")
	if postID == "" {
		response.BadRequest(c, "post ID is required")
		return
	}

	tags, err := h.tagService.GetTagsByPostID(c.Request.Context(), postID)
	if err != nil {
		handleServiceError(c, err)
		return
	}

	response.Success(c, convertTagsToResponse(tags))
}

// convertTagToResponse converts a tag entity to response DTO.
func convertTagToResponse(tag *entity.Tag) dto.TagResponse {
	if tag == nil {
		return dto.TagResponse{}
	}

	return dto.TagResponse{
		ID:        tag.ID().String(),
		Name:      tag.Name().String(),
		CreatedAt: tag.CreatedAt().Time().Format("2006-01-02T15:04:05Z07:00"),
	}
}

// convertTagsToResponse converts tag entities to response DTOs.
func convertTagsToResponse(tags []*entity.Tag) []dto.TagResponse {
	if tags == nil {
		return []dto.TagResponse{}
	}

	result := make([]dto.TagResponse, 0, len(tags))
	for _, tag := range tags {
		result = append(result, convertTagToResponse(tag))
	}
	return result
}
