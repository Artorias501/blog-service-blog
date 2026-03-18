// Package handler contains HTTP handlers for the blog service.
package handler

import (
	"strings"

	"github.com/artorias501/blog-service/pkg/response"
	"github.com/gin-gonic/gin"
)

// handleServiceError converts service errors to appropriate HTTP responses.
// It maps domain/service errors to HTTP status codes based on error message patterns.
func handleServiceError(c *gin.Context, err error) {
	if err == nil {
		return
	}

	errMsg := err.Error()
	lowerMsg := strings.ToLower(errMsg)

	// Check for common error patterns
	switch {
	case strings.Contains(lowerMsg, "not found"):
		response.NotFound(c, errMsg)
	case strings.Contains(lowerMsg, "invalid"):
		response.BadRequest(c, errMsg)
	case strings.Contains(lowerMsg, "already exists"):
		response.BadRequest(c, errMsg)
	case strings.Contains(lowerMsg, "unauthorized"):
		response.Unauthorized(c, errMsg)
	case strings.Contains(lowerMsg, "forbidden"):
		response.Forbidden(c, errMsg)
	default:
		response.InternalError(c, errMsg)
	}
}
