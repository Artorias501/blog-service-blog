package response_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/artorias501/blog-service/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func TestSuccess(t *testing.T) {
	tests := []struct {
		name     string
		data     interface{}
		expected response.Response
	}{
		{
			name: "simple data",
			data: map[string]string{"key": "value"},
			expected: response.Response{
				Code:    0,
				Message: "success",
				Data:    map[string]string{"key": "value"},
			},
		},
		{
			name: "nil data",
			data: nil,
			expected: response.Response{
				Code:    0,
				Message: "success",
				Data:    nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			response.Success(c, tt.data)

			assert.Equal(t, http.StatusOK, w.Code)

			var resp response.Response
			err := json.Unmarshal(w.Body.Bytes(), &resp)
			assert.NoError(t, err)
			assert.Equal(t, tt.expected.Code, resp.Code)
			assert.Equal(t, tt.expected.Message, resp.Message)
		})
	}
}

func TestCreated(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	data := map[string]string{"id": "123"}
	response.Created(c, data)

	assert.Equal(t, http.StatusCreated, w.Code)

	var resp response.Response
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, 0, resp.Code)
	assert.Equal(t, "created", resp.Message)
}

func TestNoContent(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodDelete, "/test", nil)

	response.NoContent(c)

	// Need to ensure the response is written
	c.Writer.WriteHeaderNow()

	assert.Equal(t, http.StatusNoContent, w.Code)
}

func TestBadRequest(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	response.BadRequest(c, "invalid input")

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var resp response.Response
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, 400, resp.Code)
	assert.Equal(t, "invalid input", resp.Message)
}

func TestValidationError(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	errors := []response.ErrorDetail{
		{Field: "title", Message: "this field is required"},
		{Field: "content", Message: "value must be at least 1 characters"},
	}

	response.ValidationError(c, errors)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var resp response.Response
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, 400, resp.Code)
	assert.Equal(t, "validation failed", resp.Message)
	assert.NotNil(t, resp.Data)
}

func TestNotFound(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	response.NotFound(c, "post not found")

	assert.Equal(t, http.StatusNotFound, w.Code)

	var resp response.Response
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, 404, resp.Code)
	assert.Equal(t, "post not found", resp.Message)
}

func TestUnauthorized(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	response.Unauthorized(c, "")

	assert.Equal(t, http.StatusUnauthorized, w.Code)

	var resp response.Response
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, 401, resp.Code)
	assert.Equal(t, "unauthorized", resp.Message)
}

func TestForbidden(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	response.Forbidden(c, "")

	assert.Equal(t, http.StatusForbidden, w.Code)

	var resp response.Response
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, 403, resp.Code)
	assert.Equal(t, "forbidden", resp.Message)
}

func TestInternalError(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	response.InternalError(c, "")

	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var resp response.Response
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, 500, resp.Code)
	assert.Equal(t, "internal server error", resp.Message)
}

func TestPaginated(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	items := []map[string]string{
		{"id": "1", "title": "Post 1"},
		{"id": "2", "title": "Post 2"},
	}

	response.Paginated(c, items, 100, 1, 10, 10)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp response.Response
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, 0, resp.Code)
	assert.Equal(t, "success", resp.Message)
	assert.NotNil(t, resp.Data)

	// Verify paginated data structure
	dataBytes, _ := json.Marshal(resp.Data)
	var paginated response.PaginatedData
	err = json.Unmarshal(dataBytes, &paginated)
	assert.NoError(t, err)
	assert.Equal(t, int64(100), paginated.Total)
	assert.Equal(t, 1, paginated.Page)
	assert.Equal(t, 10, paginated.PageSize)
	assert.Equal(t, 10, paginated.TotalPage)
}

func TestErrorDetailJSONTags(t *testing.T) {
	// Test that ErrorDetail has correct json tags
	detail := response.ErrorDetail{
		Field:   "title",
		Message: "required",
	}

	data, err := json.Marshal(detail)
	assert.NoError(t, err)

	var result map[string]interface{}
	err = json.Unmarshal(data, &result)
	assert.NoError(t, err)

	// Verify snake_case json tags
	assert.Contains(t, result, "field")
	assert.Contains(t, result, "message")
}

func TestPaginatedDataJSONTags(t *testing.T) {
	// Test that PaginatedData has correct json tags
	paginated := response.PaginatedData{
		Items:     []string{"item1", "item2"},
		Total:     100,
		Page:      1,
		PageSize:  10,
		TotalPage: 10,
	}

	data, err := json.Marshal(paginated)
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
