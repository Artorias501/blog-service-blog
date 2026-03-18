// Package dto contains Data Transfer Objects for HTTP request and response handling.
package dto

// CreateTagRequest represents the request body for creating a new tag.
type CreateTagRequest struct {
	Name string `json:"name" binding:"required,min=1,max=50"`
}

// UpdateTagRequest represents the request body for updating a tag.
type UpdateTagRequest struct {
	Name string `json:"name" binding:"required,min=1,max=50"`
}

// TagResponse represents the response data for a single tag.
type TagResponse struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	CreatedAt string `json:"created_at"`
}

// TagListResponse represents the response data for a list of tags.
type TagListResponse struct {
	Items     []TagResponse `json:"items"`
	Total     int64         `json:"total"`
	Page      int           `json:"page"`
	PageSize  int           `json:"page_size"`
	TotalPage int           `json:"total_page"`
}

// SearchTagsRequest represents query parameters for searching tags.
type SearchTagsRequest struct {
	Keyword string `form:"keyword" binding:"required,min=1"`
	Page    int    `form:"page" binding:"omitempty,min=1"`
	Size    int    `form:"size" binding:"omitempty,min=1,max=100"`
	SortBy  string `form:"sort_by" binding:"omitempty,oneof=created_at name"`
	Order   string `form:"order" binding:"omitempty,oneof=asc desc"`
}

// ListTagsRequest represents query parameters for listing tags.
type ListTagsRequest struct {
	Page   int    `form:"page" binding:"omitempty,min=1"`
	Size   int    `form:"size" binding:"omitempty,min=1,max=100"`
	SortBy string `form:"sort_by" binding:"omitempty,oneof=created_at name"`
	Order  string `form:"order" binding:"omitempty,oneof=asc desc"`
}
