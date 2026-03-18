// Package dto contains Data Transfer Objects for HTTP request and response handling.
package dto

// CreatePostRequest represents the request body for creating a new post.
type CreatePostRequest struct {
	Title   string   `json:"title" binding:"required,min=1,max=200"`
	Content string   `json:"content" binding:"required,min=1"`
	TagIDs  []string `json:"tag_ids" binding:"omitempty,dive,uuid"`
}

// UpdatePostRequest represents the request body for updating a post.
type UpdatePostRequest struct {
	Title   *string `json:"title" binding:"omitempty,min=1,max=200"`
	Content *string `json:"content" binding:"omitempty,min=1"`
}

// PostResponse represents the response data for a single post.
type PostResponse struct {
	ID          string     `json:"id"`
	Title       string     `json:"title"`
	Content     string     `json:"content"`
	Summary     *string    `json:"summary"`
	Tags        []TagBrief `json:"tags"`
	CreatedAt   string     `json:"created_at"`
	UpdatedAt   string     `json:"updated_at"`
	PublishedAt *string    `json:"published_at"`
}

// PostListResponse represents the response data for a list of posts.
type PostListResponse struct {
	Items     []PostResponse `json:"items"`
	Total     int64          `json:"total"`
	Page      int            `json:"page"`
	PageSize  int            `json:"page_size"`
	TotalPage int            `json:"total_page"`
}

// PostBrief represents a brief post summary for list responses.
type PostBrief struct {
	ID        string `json:"id"`
	Title     string `json:"title"`
	Summary   string `json:"summary"`
	CreatedAt string `json:"created_at"`
}

// TagBrief represents a brief tag summary for post responses.
type TagBrief struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// AddTagToPostRequest represents the request body for adding a tag to a post.
type AddTagToPostRequest struct {
	TagID string `json:"tag_id" binding:"required,uuid"`
}

// SearchPostsRequest represents query parameters for searching posts.
type SearchPostsRequest struct {
	Keyword string `form:"keyword" binding:"required,min=1"`
	Page    int    `form:"page" binding:"omitempty,min=1"`
	Size    int    `form:"size" binding:"omitempty,min=1,max=100"`
	SortBy  string `form:"sort_by" binding:"omitempty,oneof=created_at updated_at title"`
	Order   string `form:"order" binding:"omitempty,oneof=asc desc"`
}

// ListPostsRequest represents query parameters for listing posts.
type ListPostsRequest struct {
	Page   int    `form:"page" binding:"omitempty,min=1"`
	Size   int    `form:"size" binding:"omitempty,min=1,max=100"`
	SortBy string `form:"sort_by" binding:"omitempty,oneof=created_at updated_at title"`
	Order  string `form:"order" binding:"omitempty,oneof=asc desc"`
	TagID  string `form:"tag_id" binding:"omitempty,uuid"`
}

// LikePostRequest represents the request body for liking a post.
type LikePostRequest struct {
	// No body needed - just increment counter
}
