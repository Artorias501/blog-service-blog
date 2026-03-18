// Package handler contains HTTP handlers for the blog service.
package handler

import (
	"github.com/artorias501/blog-service/internal/domain/entity"
	"github.com/artorias501/blog-service/internal/handler/dto"
	"github.com/artorias501/blog-service/internal/service"
	"github.com/artorias501/blog-service/pkg/response"
	"github.com/gin-gonic/gin"
)

// CommentHandler handles HTTP requests for comment operations.
type CommentHandler struct {
	commentService *service.CommentService
}

// NewCommentHandler creates a new CommentHandler instance.
func NewCommentHandler(commentService *service.CommentService) *CommentHandler {
	return &CommentHandler{
		commentService: commentService,
	}
}

// CreateComment handles POST /api/comments
// @Summary Create a new comment
// @Description Create a new comment on a post
// @Tags comments
// @Accept json
// @Produce json
// @Param request body dto.CreateCommentRequest true "Comment creation request"
// @Success 201 {object} response.Response{data=dto.CommentResponse}
// @Failure 400 {object} response.Response
// @Router /api/comments [post]
func (h *CommentHandler) CreateComment(c *gin.Context) {
	var req dto.CreateCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "invalid request body: "+err.Error())
		return
	}

	input := service.CreateCommentInput{
		PostID:      req.PostID,
		AuthorName:  req.AuthorName,
		AuthorEmail: req.AuthorEmail,
		Content:     req.Content,
	}

	comment, err := h.commentService.CreateComment(c.Request.Context(), input)
	if err != nil {
		handleServiceError(c, err)
		return
	}

	response.Created(c, convertCommentToResponse(comment))
}

// GetComment handles GET /api/comments/:id
// @Summary Get a comment by ID
// @Description Retrieve a single comment by its ID
// @Tags comments
// @Produce json
// @Param id path string true "Comment ID"
// @Success 200 {object} response.Response{data=dto.CommentResponse}
// @Failure 404 {object} response.Response
// @Router /api/comments/{id} [get]
func (h *CommentHandler) GetComment(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.BadRequest(c, "comment ID is required")
		return
	}

	comment, err := h.commentService.GetCommentByID(c.Request.Context(), id)
	if err != nil {
		handleServiceError(c, err)
		return
	}

	response.Success(c, convertCommentToResponse(comment))
}

// UpdateComment handles PUT /api/comments/:id
// @Summary Update a comment
// @Description Update an existing comment's content
// @Tags comments
// @Accept json
// @Produce json
// @Param id path string true "Comment ID"
// @Param request body dto.UpdateCommentRequest true "Comment update request"
// @Success 200 {object} response.Response{data=dto.CommentResponse}
// @Failure 400,404 {object} response.Response
// @Router /api/comments/{id} [put]
func (h *CommentHandler) UpdateComment(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.BadRequest(c, "comment ID is required")
		return
	}

	var req dto.UpdateCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "invalid request body: "+err.Error())
		return
	}

	input := service.UpdateCommentInput{
		Content: req.Content,
	}

	comment, err := h.commentService.UpdateComment(c.Request.Context(), id, input)
	if err != nil {
		handleServiceError(c, err)
		return
	}

	response.Success(c, convertCommentToResponse(comment))
}

// DeleteComment handles DELETE /api/comments/:id
// @Summary Delete a comment
// @Description Delete a comment by ID
// @Tags comments
// @Param id path string true "Comment ID"
// @Success 204 "No Content"
// @Failure 404 {object} response.Response
// @Router /api/comments/{id} [delete]
func (h *CommentHandler) DeleteComment(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.BadRequest(c, "comment ID is required")
		return
	}

	err := h.commentService.DeleteComment(c.Request.Context(), id)
	if err != nil {
		handleServiceError(c, err)
		return
	}

	response.NoContent(c)
}

// ListComments handles GET /api/comments
// @Summary List comments
// @Description Retrieve a paginated list of all comments (admin)
// @Tags comments
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param size query int false "Page size" default(10)
// @Param sort_by query string false "Sort field"
// @Param order query string false "Sort order"
// @Success 200 {object} response.Response{data=dto.CommentListResponse}
// @Router /api/comments [get]
func (h *CommentHandler) ListComments(c *gin.Context) {
	var req dto.ListCommentsRequest
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

	input := service.ListCommentsInput{
		Page:     req.Page,
		PageSize: req.Size,
		SortBy:   req.SortBy,
		Order:    req.Order,
	}

	result, err := h.commentService.ListComments(c.Request.Context(), input)
	if err != nil {
		handleServiceError(c, err)
		return
	}

	response.Paginated(c, convertCommentsToResponse(result.Items), result.Total, result.Page, result.PageSize, result.TotalPage)
}

// ListCommentsByPost handles GET /api/posts/:id/comments
// @Summary List comments by post
// @Description Retrieve paginated comments for a specific post
// @Tags comments
// @Produce json
// @Param id path string true "Post ID"
// @Param page query int false "Page number" default(1)
// @Param size query int false "Page size" default(10)
// @Success 200 {object} response.Response{data=dto.CommentListResponse}
// @Failure 404 {object} response.Response
// @Router /api/posts/{id}/comments [get]
func (h *CommentHandler) ListCommentsByPost(c *gin.Context) {
	postID := c.Param("id")
	if postID == "" {
		response.BadRequest(c, "post ID is required")
		return
	}

	var req dto.ListCommentsByPostRequest
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

	input := service.ListCommentsByPostInput{
		PostID:   postID,
		Page:     req.Page,
		PageSize: req.Size,
		SortBy:   req.SortBy,
		Order:    req.Order,
	}

	result, err := h.commentService.ListCommentsByPost(c.Request.Context(), input)
	if err != nil {
		handleServiceError(c, err)
		return
	}

	response.Paginated(c, convertCommentsToResponse(result.Items), result.Total, result.Page, result.PageSize, result.TotalPage)
}

// ListCommentsByStatus handles GET /api/comments/status/:status
// @Summary List comments by status
// @Description Retrieve paginated comments filtered by status (admin)
// @Tags comments
// @Produce json
// @Param status path string true "Comment status (pending, approved, rejected, spam)"
// @Param page query int false "Page number" default(1)
// @Param size query int false "Page size" default(10)
// @Success 200 {object} response.Response{data=dto.CommentListResponse}
// @Router /api/comments/status/{status} [get]
func (h *CommentHandler) ListCommentsByStatus(c *gin.Context) {
	status := c.Param("status")
	if status == "" {
		response.BadRequest(c, "status is required")
		return
	}

	var req dto.ListCommentsByStatusRequest
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

	input := service.ListCommentsByStatusInput{
		Status:   status,
		Page:     req.Page,
		PageSize: req.Size,
		SortBy:   req.SortBy,
		Order:    req.Order,
	}

	result, err := h.commentService.ListCommentsByStatus(c.Request.Context(), input)
	if err != nil {
		handleServiceError(c, err)
		return
	}

	response.Paginated(c, convertCommentsToResponse(result.Items), result.Total, result.Page, result.PageSize, result.TotalPage)
}

// ApproveComment handles POST /api/comments/:id/approve
// @Summary Approve a comment
// @Description Approve a pending comment
// @Tags comments
// @Param id path string true "Comment ID"
// @Success 204 "No Content"
// @Failure 404 {object} response.Response
// @Router /api/comments/{id}/approve [post]
func (h *CommentHandler) ApproveComment(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.BadRequest(c, "comment ID is required")
		return
	}

	err := h.commentService.ApproveComment(c.Request.Context(), id)
	if err != nil {
		handleServiceError(c, err)
		return
	}

	response.NoContent(c)
}

// RejectComment handles POST /api/comments/:id/reject
// @Summary Reject a comment
// @Description Reject a pending comment
// @Tags comments
// @Param id path string true "Comment ID"
// @Success 204 "No Content"
// @Failure 404 {object} response.Response
// @Router /api/comments/{id}/reject [post]
func (h *CommentHandler) RejectComment(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.BadRequest(c, "comment ID is required")
		return
	}

	err := h.commentService.RejectComment(c.Request.Context(), id)
	if err != nil {
		handleServiceError(c, err)
		return
	}

	response.NoContent(c)
}

// MarkCommentAsSpam handles POST /api/comments/:id/spam
// @Summary Mark comment as spam
// @Description Mark a comment as spam
// @Tags comments
// @Param id path string true "Comment ID"
// @Success 204 "No Content"
// @Failure 404 {object} response.Response
// @Router /api/comments/{id}/spam [post]
func (h *CommentHandler) MarkCommentAsSpam(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.BadRequest(c, "comment ID is required")
		return
	}

	err := h.commentService.MarkCommentAsSpam(c.Request.Context(), id)
	if err != nil {
		handleServiceError(c, err)
		return
	}

	response.NoContent(c)
}

// GetCommentCount handles GET /api/posts/:id/comments/count
// @Summary Get comment count
// @Description Get the count of comments for a post
// @Tags comments
// @Produce json
// @Param id path string true "Post ID"
// @Success 200 {object} response.Response{data=dto.CommentCountResponse}
// @Failure 404 {object} response.Response
// @Router /api/posts/{id}/comments/count [get]
func (h *CommentHandler) GetCommentCount(c *gin.Context) {
	postID := c.Param("id")
	if postID == "" {
		response.BadRequest(c, "post ID is required")
		return
	}

	count, err := h.commentService.GetCommentCountByPostID(c.Request.Context(), postID)
	if err != nil {
		handleServiceError(c, err)
		return
	}

	response.Success(c, dto.CommentCountResponse{
		PostID: postID,
		Count:  count,
	})
}

// convertCommentToResponse converts a comment entity to response DTO.
func convertCommentToResponse(comment *entity.Comment) dto.CommentResponse {
	if comment == nil {
		return dto.CommentResponse{}
	}

	return dto.CommentResponse{
		ID:          comment.ID().String(),
		PostID:      "", // PostID is not stored in comment entity directly
		AuthorName:  comment.AuthorName().String(),
		AuthorEmail: comment.AuthorEmail().String(),
		Content:     comment.Content().String(),
		Status:      comment.Status().String(),
		CreatedAt:   comment.CreatedAt().Time().Format("2006-01-02T15:04:05Z07:00"),
	}
}

// convertCommentsToResponse converts comment entities to response DTOs.
func convertCommentsToResponse(comments []*entity.Comment) []dto.CommentResponse {
	if comments == nil {
		return []dto.CommentResponse{}
	}

	result := make([]dto.CommentResponse, 0, len(comments))
	for _, comment := range comments {
		result = append(result, convertCommentToResponse(comment))
	}
	return result
}
