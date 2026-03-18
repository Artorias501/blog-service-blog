package middleware

import (
	"strings"

	"github.com/artorias501/blog-service/pkg/config"
	"github.com/artorias501/blog-service/pkg/response"
	"github.com/gin-gonic/gin"
)

// AdminAuth validates Bearer token from Authorization header.
// It checks against the configured admin token.
func AdminAuth(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		// If no admin token is configured, deny all access
		if cfg.Auth.AdminToken == "" {
			response.Unauthorized(c, "admin token not configured")
			c.Abort()
			return
		}

		// Get Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Unauthorized(c, "missing authorization header")
			c.Abort()
			return
		}

		// Check Bearer prefix
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			response.Unauthorized(c, "invalid authorization format")
			c.Abort()
			return
		}

		token := parts[1]
		if token == "" {
			response.Unauthorized(c, "missing token")
			c.Abort()
			return
		}

		// Validate token
		if token != cfg.Auth.AdminToken {
			response.Unauthorized(c, "invalid token")
			c.Abort()
			return
		}

		// Set admin flag in context
		c.Set("is_admin", true)

		c.Next()
	}
}
