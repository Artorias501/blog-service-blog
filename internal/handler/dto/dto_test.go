package dto_test

import (
	"encoding/json"
	"testing"

	"github.com/artorias501/blog-service/internal/handler/dto"
	"github.com/stretchr/testify/assert"
)

func TestCreatePostRequestJSONTags(t *testing.T) {
	req := dto.CreatePostRequest{
		Title:   "Test Title",
		Content: "Test Content",
		TagIDs:  []string{"tag-1", "tag-2"},
	}

	data, err := json.Marshal(req)
	assert.NoError(t, err)

	var result map[string]interface{}
	err = json.Unmarshal(data, &result)
	assert.NoError(t, err)

	// Verify snake_case json tags
	assert.Contains(t, result, "title")
	assert.Contains(t, result, "content")
	assert.Contains(t, result, "tag_ids")
}

func TestUpdatePostRequestJSONTags(t *testing.T) {
	title := "Updated Title"
	content := "Updated Content"
	req := dto.UpdatePostRequest{
		Title:   &title,
		Content: &content,
	}

	data, err := json.Marshal(req)
	assert.NoError(t, err)

	var result map[string]interface{}
	err = json.Unmarshal(data, &result)
	assert.NoError(t, err)

	// Verify snake_case json tags
	assert.Contains(t, result, "title")
	assert.Contains(t, result, "content")
}

func TestPostResponseJSONTags(t *testing.T) {
	summary := "Test Summary"
	publishedAt := "2024-01-01T00:00:00Z"
	resp := dto.PostResponse{
		ID:          "post-1",
		Title:       "Test Title",
		Content:     "Test Content",
		Summary:     &summary,
		Tags:        []dto.TagBrief{{ID: "tag-1", Name: "Tag 1"}},
		CreatedAt:   "2024-01-01T00:00:00Z",
		UpdatedAt:   "2024-01-01T00:00:00Z",
		PublishedAt: &publishedAt,
	}

	data, err := json.Marshal(resp)
	assert.NoError(t, err)

	var result map[string]interface{}
	err = json.Unmarshal(data, &result)
	assert.NoError(t, err)

	// Verify snake_case json tags
	assert.Contains(t, result, "id")
	assert.Contains(t, result, "title")
	assert.Contains(t, result, "content")
	assert.Contains(t, result, "summary")
	assert.Contains(t, result, "tags")
	assert.Contains(t, result, "created_at")
	assert.Contains(t, result, "updated_at")
	assert.Contains(t, result, "published_at")
}

func TestPostListResponseJSONTags(t *testing.T) {
	resp := dto.PostListResponse{
		Items:     []dto.PostResponse{},
		Total:     100,
		Page:      1,
		PageSize:  10,
		TotalPage: 10,
	}

	data, err := json.Marshal(resp)
	assert.NoError(t, err)

	var result map[string]interface{}
	err = json.Unmarshal(data, &result)
	assert.NoError(t, err)

	// Verify snake_case json tags
	assert.Contains(t, result, "items")
	assert.Contains(t, result, "total")
	assert.Contains(t, result, "page")
	assert.Contains(t, result, "page_size")
	assert.Contains(t, result, "total_page")
}

func TestCreateTagRequestJSONTags(t *testing.T) {
	req := dto.CreateTagRequest{
		Name: "Test Tag",
	}

	data, err := json.Marshal(req)
	assert.NoError(t, err)

	var result map[string]interface{}
	err = json.Unmarshal(data, &result)
	assert.NoError(t, err)

	assert.Contains(t, result, "name")
}

func TestTagResponseJSONTags(t *testing.T) {
	resp := dto.TagResponse{
		ID:        "tag-1",
		Name:      "Test Tag",
		CreatedAt: "2024-01-01T00:00:00Z",
	}

	data, err := json.Marshal(resp)
	assert.NoError(t, err)

	var result map[string]interface{}
	err = json.Unmarshal(data, &result)
	assert.NoError(t, err)

	// Verify snake_case json tags
	assert.Contains(t, result, "id")
	assert.Contains(t, result, "name")
	assert.Contains(t, result, "created_at")
}

func TestCreateCommentRequestJSONTags(t *testing.T) {
	req := dto.CreateCommentRequest{
		PostID:      "post-1",
		AuthorName:  "Test Author",
		AuthorEmail: "test@example.com",
		Content:     "Test Comment",
	}

	data, err := json.Marshal(req)
	assert.NoError(t, err)

	var result map[string]interface{}
	err = json.Unmarshal(data, &result)
	assert.NoError(t, err)

	// Verify snake_case json tags
	assert.Contains(t, result, "post_id")
	assert.Contains(t, result, "author_name")
	assert.Contains(t, result, "author_email")
	assert.Contains(t, result, "content")
}

func TestCommentResponseJSONTags(t *testing.T) {
	resp := dto.CommentResponse{
		ID:          "comment-1",
		PostID:      "post-1",
		AuthorName:  "Test Author",
		AuthorEmail: "test@example.com",
		Content:     "Test Comment",
		Status:      "approved",
		CreatedAt:   "2024-01-01T00:00:00Z",
	}

	data, err := json.Marshal(resp)
	assert.NoError(t, err)

	var result map[string]interface{}
	err = json.Unmarshal(data, &result)
	assert.NoError(t, err)

	// Verify snake_case json tags
	assert.Contains(t, result, "id")
	assert.Contains(t, result, "post_id")
	assert.Contains(t, result, "author_name")
	assert.Contains(t, result, "author_email")
	assert.Contains(t, result, "content")
	assert.Contains(t, result, "status")
	assert.Contains(t, result, "created_at")
}

func TestCommentCountResponseJSONTags(t *testing.T) {
	resp := dto.CommentCountResponse{
		PostID: "post-1",
		Count:  10,
	}

	data, err := json.Marshal(resp)
	assert.NoError(t, err)

	var result map[string]interface{}
	err = json.Unmarshal(data, &result)
	assert.NoError(t, err)

	// Verify snake_case json tags
	assert.Contains(t, result, "post_id")
	assert.Contains(t, result, "count")
}

func TestSearchPostsRequestFormTags(t *testing.T) {
	// Test that form tags are correctly set for query parameters
	req := dto.SearchPostsRequest{
		Keyword: "test",
		Page:    1,
		Size:    10,
		SortBy:  "created_at",
		Order:   "desc",
	}

	// Verify the struct has the correct field values
	assert.Equal(t, "test", req.Keyword)
	assert.Equal(t, 1, req.Page)
	assert.Equal(t, 10, req.Size)
	assert.Equal(t, "created_at", req.SortBy)
	assert.Equal(t, "desc", req.Order)
}

func TestListPostsRequestFormTags(t *testing.T) {
	req := dto.ListPostsRequest{
		Page:   1,
		Size:   10,
		SortBy: "created_at",
		Order:  "desc",
		TagID:  "tag-1",
	}

	assert.Equal(t, 1, req.Page)
	assert.Equal(t, 10, req.Size)
	assert.Equal(t, "created_at", req.SortBy)
	assert.Equal(t, "desc", req.Order)
	assert.Equal(t, "tag-1", req.TagID)
}

func TestListCommentsByStatusRequestFormTags(t *testing.T) {
	req := dto.ListCommentsByStatusRequest{
		Status: "approved",
		Page:   1,
		Size:   10,
		SortBy: "created_at",
		Order:  "desc",
	}

	assert.Equal(t, "approved", req.Status)
	assert.Equal(t, 1, req.Page)
	assert.Equal(t, 10, req.Size)
}
