// Package dto contains Data Transfer Objects for HTTP request and response handling.
package dto

// CreateCommentRequest represents the request body for creating a new comment.
type CreateCommentRequest struct {
	PostID      string `json:"post_id" binding:"required,uuid"`
	AuthorName  string `json:"author_name" binding:"required,min=1,max=100"`
	AuthorEmail string `json:"author_email" binding:"required,email"`
	Content     string `json:"content" binding:"required,min=1,max=5000"`
}

// UpdateCommentRequest represents the request body for updating a comment.
type UpdateCommentRequest struct {
	Content string `json:"content" binding:"required,min=1,max=5000"`
}

// CommentResponse represents the response data for a single comment.
type CommentResponse struct {
	ID          string `json:"id"`
	PostID      string `json:"post_id"`
	AuthorName  string `json:"author_name"`
	AuthorEmail string `json:"author_email"`
	Content     string `json:"content"`
	Status      string `json:"status"`
	CreatedAt   string `json:"created_at"`
}

// CommentListResponse represents the response data for a list of comments.
type CommentListResponse struct {
	Items     []CommentResponse `json:"items"`
	Total     int64             `json:"total"`
	Page      int               `json:"page"`
	PageSize  int               `json:"page_size"`
	TotalPage int               `json:"total_page"`
}

// ListCommentsRequest represents query parameters for listing comments.
type ListCommentsRequest struct {
	Page   int    `form:"page" binding:"omitempty,min=1"`
	Size   int    `form:"size" binding:"omitempty,min=1,max=100"`
	SortBy string `form:"sort_by" binding:"omitempty,oneof=created_at status"`
	Order  string `form:"order" binding:"omitempty,oneof=asc desc"`
}

// ListCommentsByPostRequest represents query parameters for listing comments by post.
type ListCommentsByPostRequest struct {
	Page   int    `form:"page" binding:"omitempty,min=1"`
	Size   int    `form:"size" binding:"omitempty,min=1,max=100"`
	SortBy string `form:"sort_by" binding:"omitempty,oneof=created_at"`
	Order  string `form:"order" binding:"omitempty,oneof=asc desc"`
}

// ListCommentsByStatusRequest represents query parameters for listing comments by status.
type ListCommentsByStatusRequest struct {
	Status string `form:"status" binding:"required,oneof=pending approved rejected spam"`
	Page   int    `form:"page" binding:"omitempty,min=1"`
	Size   int    `form:"size" binding:"omitempty,min=1,max=100"`
	SortBy string `form:"sort_by" binding:"omitempty,oneof=created_at"`
	Order  string `form:"order" binding:"omitempty,oneof=asc desc"`
}

// CommentCountResponse represents the response for comment count.
type CommentCountResponse struct {
	PostID string `json:"post_id"`
	Count  int64  `json:"count"`
}
