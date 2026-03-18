// Package response provides standard HTTP response utilities for the blog service.
package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response represents the standard API response structure.
// All API responses must follow this format with code, message, and data fields.
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// ErrorDetail represents detailed error information for validation errors.
type ErrorDetail struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// PaginatedData represents paginated response data with metadata.
type PaginatedData struct {
	Items     interface{} `json:"items"`
	Total     int64       `json:"total"`
	Page      int         `json:"page"`
	PageSize  int         `json:"page_size"`
	TotalPage int         `json:"total_page"`
}

// Success sends a successful response with HTTP 200.
func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    0,
		Message: "success",
		Data:    data,
	})
}

// SuccessWithStatus sends a successful response with custom HTTP status.
func SuccessWithStatus(c *gin.Context, status int, data interface{}) {
	c.JSON(status, Response{
		Code:    0,
		Message: "success",
		Data:    data,
	})
}

// Created sends a successful response for resource creation with HTTP 201.
func Created(c *gin.Context, data interface{}) {
	c.JSON(http.StatusCreated, Response{
		Code:    0,
		Message: "created",
		Data:    data,
	})
}

// NoContent sends a successful response with HTTP 204 (no body).
func NoContent(c *gin.Context) {
	c.Status(http.StatusNoContent)
}

// Error sends an error response with the specified HTTP status.
func Error(c *gin.Context, status int, code int, message string) {
	c.JSON(status, Response{
		Code:    code,
		Message: message,
		Data:    nil,
	})
}

// BadRequest sends a 400 error response.
func BadRequest(c *gin.Context, message string) {
	Error(c, http.StatusBadRequest, 400, message)
}

// ValidationError sends a 400 error response with detailed field errors.
func ValidationError(c *gin.Context, errors []ErrorDetail) {
	c.JSON(http.StatusBadRequest, Response{
		Code:    400,
		Message: "validation failed",
		Data:    errors,
	})
}

// Unauthorized sends a 401 error response.
func Unauthorized(c *gin.Context, message string) {
	if message == "" {
		message = "unauthorized"
	}
	Error(c, http.StatusUnauthorized, 401, message)
}

// Forbidden sends a 403 error response.
func Forbidden(c *gin.Context, message string) {
	if message == "" {
		message = "forbidden"
	}
	Error(c, http.StatusForbidden, 403, message)
}

// NotFound sends a 404 error response.
func NotFound(c *gin.Context, message string) {
	if message == "" {
		message = "not found"
	}
	Error(c, http.StatusNotFound, 404, message)
}

// InternalError sends a 500 error response.
func InternalError(c *gin.Context, message string) {
	if message == "" {
		message = "internal server error"
	}
	Error(c, http.StatusInternalServerError, 500, message)
}

// Paginated sends a successful paginated response.
func Paginated(c *gin.Context, items interface{}, total int64, page, pageSize, totalPage int) {
	Success(c, PaginatedData{
		Items:     items,
		Total:     total,
		Page:      page,
		PageSize:  pageSize,
		TotalPage: totalPage,
	})
}
