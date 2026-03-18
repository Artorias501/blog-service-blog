package middleware

import (
	"strings"

	"github.com/artorias501/blog-service/pkg/config"
	"github.com/gin-gonic/gin"
)

// CORS handles Cross-Origin Resource Sharing.
// It sets appropriate headers based on the configuration.
func CORS(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")

		// Check if origin is allowed
		allowedOrigin := ""
		for _, o := range cfg.CORS.AllowedOrigins {
			if o == "*" {
				allowedOrigin = "*"
				break
			}
			if o == origin {
				allowedOrigin = origin
				break
			}
		}

		// Set CORS headers if origin is allowed
		if allowedOrigin != "" {
			c.Header("Access-Control-Allow-Origin", allowedOrigin)
			c.Header("Access-Control-Allow-Methods", strings.Join(cfg.CORS.AllowedMethods, ", "))
			c.Header("Access-Control-Allow-Headers", strings.Join(cfg.CORS.AllowedHeaders, ", "))

			if cfg.CORS.AllowCredentials {
				c.Header("Access-Control-Allow-Credentials", "true")
			}

			if cfg.CORS.MaxAge > 0 {
				c.Header("Access-Control-Max-Age", intToStr(cfg.CORS.MaxAge))
			}
		}

		// Handle preflight request
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

// intToStr converts int to string without importing strconv.
func intToStr(n int) string {
	if n == 0 {
		return "0"
	}
	var neg bool
	if n < 0 {
		neg = true
		n = -n
	}
	var digits []byte
	for n > 0 {
		digits = append([]byte{byte('0' + n%10)}, digits...)
		n /= 10
	}
	if neg {
		digits = append([]byte{'-'}, digits...)
	}
	return string(digits)
}
