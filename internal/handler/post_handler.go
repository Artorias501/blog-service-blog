// Package handler contains HTTP handlers for the blog service.
package handler

import (
	"github.com/artorias501/blog-service/internal/domain/entity"
	"github.com/artorias501/blog-service/internal/domain/repository"
	"github.com/artorias501/blog-service/internal/handler/dto"
	"github.com/artorias501/blog-service/internal/service"
	"github.com/artorias501/blog-service/pkg/response"
	"github.com/gin-gonic/gin"
)

// PostHandler handles HTTP requests for post operations.
type PostHandler struct {
	postService *service.PostService
}

// NewPostHandler creates a new PostHandler instance.
func NewPostHandler(postService *service.PostService) *PostHandler {
	return &PostHandler{
		postService: postService,
	}
}

// CreatePost handles POST /api/posts
// @Summary Create a new post
// @Description Create a new blog post with title, content, and optional tags
// @Tags posts
// @Accept json
// @Produce json
// @Param request body dto.CreatePostRequest true "Post creation request"
// @Success 201 {object} response.Response{data=dto.PostResponse}
// @Failure 400 {object} response.Response
// @Router /api/posts [post]
func (h *PostHandler) CreatePost(c *gin.Context) {
	var req dto.CreatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "invalid request body: "+err.Error())
		return
	}

	input := service.CreatePostInput{
		Title:   req.Title,
		Content: req.Content,
		TagIDs:  req.TagIDs,
	}

	post, err := h.postService.CreatePost(c.Request.Context(), input)
	if err != nil {
		handleServiceError(c, err)
		return
	}

	response.Created(c, convertPostToResponse(post))
}

// GetPost handles GET /api/posts/:id
// @Summary Get a post by ID
// @Description Retrieve a single post with its tags
// @Tags posts
// @Produce json
// @Param id path string true "Post ID"
// @Success 200 {object} response.Response{data=dto.PostResponse}
// @Failure 404 {object} response.Response
// @Router /api/posts/{id} [get]
func (h *PostHandler) GetPost(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.BadRequest(c, "post ID is required")
		return
	}

	post, err := h.postService.GetPostByIDWithTags(c.Request.Context(), id)
	if err != nil {
		handleServiceError(c, err)
		return
	}

	response.Success(c, convertPostToResponse(post))
}

// UpdatePost handles PUT /api/posts/:id
// @Summary Update a post
// @Description Update an existing post's title or content
// @Tags posts
// @Accept json
// @Produce json
// @Param id path string true "Post ID"
// @Param request body dto.UpdatePostRequest true "Post update request"
// @Success 200 {object} response.Response{data=dto.PostResponse}
// @Failure 400,404 {object} response.Response
// @Router /api/posts/{id} [put]
func (h *PostHandler) UpdatePost(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.BadRequest(c, "post ID is required")
		return
	}

	var req dto.UpdatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "invalid request body: "+err.Error())
		return
	}

	input := service.UpdatePostInput{
		Title:   req.Title,
		Content: req.Content,
	}

	post, err := h.postService.UpdatePost(c.Request.Context(), id, input)
	if err != nil {
		handleServiceError(c, err)
		return
	}

	response.Success(c, convertPostToResponse(post))
}

// DeletePost handles DELETE /api/posts/:id
// @Summary Delete a post
// @Description Delete a post by ID
// @Tags posts
// @Param id path string true "Post ID"
// @Success 204 "No Content"
// @Failure 404 {object} response.Response
// @Router /api/posts/{id} [delete]
func (h *PostHandler) DeletePost(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.BadRequest(c, "post ID is required")
		return
	}

	err := h.postService.DeletePost(c.Request.Context(), id)
	if err != nil {
		handleServiceError(c, err)
		return
	}

	response.NoContent(c)
}

// ListPosts handles GET /api/posts
// @Summary List posts
// @Description Retrieve a paginated list of posts
// @Tags posts
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param size query int false "Page size" default(10)
// @Param sort_by query string false "Sort field"
// @Param order query string false "Sort order"
// @Param tag_id query string false "Filter by tag ID"
// @Success 200 {object} response.Response{data=dto.PostListResponse}
// @Router /api/posts [get]
func (h *PostHandler) ListPosts(c *gin.Context) {
	var req dto.ListPostsRequest
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

	input := service.ListPostsInput{
		Page:     req.Page,
		PageSize: req.Size,
		SortBy:   req.SortBy,
		Order:    req.Order,
	}

	var result *repository.PostListResult
	var err error

	if req.TagID != "" {
		result, err = h.postService.ListPostsByTag(c.Request.Context(), req.TagID, input)
	} else {
		result, err = h.postService.ListPosts(c.Request.Context(), input)
	}

	if err != nil {
		handleServiceError(c, err)
		return
	}

	response.Paginated(c, convertPostsToResponse(result.Items), result.Total, result.Page, result.PageSize, result.TotalPage)
}

// SearchPosts handles GET /api/posts/search
// @Summary Search posts
// @Description Search posts by keyword
// @Tags posts
// @Produce json
// @Param keyword query string true "Search keyword"
// @Param page query int false "Page number" default(1)
// @Param size query int false "Page size" default(10)
// @Success 200 {object} response.Response{data=dto.PostListResponse}
// @Router /api/posts/search [get]
func (h *PostHandler) SearchPosts(c *gin.Context) {
	var req dto.SearchPostsRequest
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

	input := service.SearchPostsInput{
		Keyword:  req.Keyword,
		Page:     req.Page,
		PageSize: req.Size,
		SortBy:   req.SortBy,
		Order:    req.Order,
	}

	result, err := h.postService.SearchPosts(c.Request.Context(), input)
	if err != nil {
		handleServiceError(c, err)
		return
	}

	response.Paginated(c, convertPostsToResponse(result.Items), result.Total, result.Page, result.PageSize, result.TotalPage)
}

// LikePost handles POST /api/posts/:id/like
// @Summary Like a post
// @Description Increment the like count for a post
// @Tags posts
// @Param id path string true "Post ID"
// @Success 204 "No Content"
// @Failure 404 {object} response.Response
// @Router /api/posts/{id}/like [post]
func (h *PostHandler) LikePost(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.BadRequest(c, "post ID is required")
		return
	}

	err := h.postService.LikePost(c.Request.Context(), id)
	if err != nil {
		handleServiceError(c, err)
		return
	}

	response.NoContent(c)
}

// AddTagToPost handles POST /api/posts/:id/tags
// @Summary Add tag to post
// @Description Associate a tag with a post
// @Tags posts
// @Accept json
// @Param id path string true "Post ID"
// @Param request body dto.AddTagToPostRequest true "Tag ID"
// @Success 204 "No Content"
// @Failure 400,404 {object} response.Response
// @Router /api/posts/{id}/tags [post]
func (h *PostHandler) AddTagToPost(c *gin.Context) {
	postID := c.Param("id")
	if postID == "" {
		response.BadRequest(c, "post ID is required")
		return
	}

	var req dto.AddTagToPostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "invalid request body: "+err.Error())
		return
	}

	err := h.postService.AddTagToPost(c.Request.Context(), postID, req.TagID)
	if err != nil {
		handleServiceError(c, err)
		return
	}

	response.NoContent(c)
}

// RemoveTagFromPost handles DELETE /api/posts/:id/tags/:tag_id
// @Summary Remove tag from post
// @Description Disassociate a tag from a post
// @Tags posts
// @Param id path string true "Post ID"
// @Param tag_id path string true "Tag ID"
// @Success 204 "No Content"
// @Failure 404 {object} response.Response
// @Router /api/posts/{id}/tags/{tag_id} [delete]
func (h *PostHandler) RemoveTagFromPost(c *gin.Context) {
	postID := c.Param("id")
	tagID := c.Param("tag_id")
	if postID == "" || tagID == "" {
		response.BadRequest(c, "post ID and tag ID are required")
		return
	}

	err := h.postService.RemoveTagFromPost(c.Request.Context(), postID, tagID)
	if err != nil {
		handleServiceError(c, err)
		return
	}

	response.NoContent(c)
}

// convertPostToResponse converts a post entity to response DTO.
func convertPostToResponse(post *entity.Post) dto.PostResponse {
	if post == nil {
		return dto.PostResponse{}
	}

	var summary *string
	if post.Summary() != nil {
		s := post.Summary().String()
		summary = &s
	}

	var publishedAt *string
	if post.PublishedAt() != nil {
		pa := post.PublishedAt().Format("2006-01-02T15:04:05Z07:00")
		publishedAt = &pa
	}

	tags := make([]dto.TagBrief, 0, len(post.Tags()))
	for _, t := range post.Tags() {
		tags = append(tags, dto.TagBrief{
			ID:   t.ID().String(),
			Name: t.Name().String(),
		})
	}

	return dto.PostResponse{
		ID:          post.ID().String(),
		Title:       post.Title().String(),
		Content:     post.Content().String(),
		Summary:     summary,
		Tags:        tags,
		CreatedAt:   post.CreatedAt().Time().Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:   post.UpdatedAt().Time().Format("2006-01-02T15:04:05Z07:00"),
		PublishedAt: publishedAt,
	}
}

// convertPostsToResponse converts post entities to response DTOs.
func convertPostsToResponse(posts []*entity.Post) []dto.PostResponse {
	if posts == nil {
		return []dto.PostResponse{}
	}

	result := make([]dto.PostResponse, 0, len(posts))
	for _, post := range posts {
		result = append(result, convertPostToResponse(post))
	}
	return result
}
