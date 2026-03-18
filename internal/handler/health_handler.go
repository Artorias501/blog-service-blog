// Package handler contains HTTP handlers for the blog service.
package handler

import (
	"net/http"

	"github.com/artorias501/blog-service/pkg/config"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

// HealthHandler handles health check requests.
type HealthHandler struct {
	config *config.Config
	db     *gorm.DB
	redis  *redis.Client
}

// HealthStatus represents the health status response.
type HealthStatus struct {
	Status  string           `json:"status"`
	Service string           `json:"service"`
	Version string           `json:"version"`
	Checks  map[string]Check `json:"checks"`
}

// Check represents a single health check result.
type Check struct {
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`
}

// NewHealthHandler creates a new HealthHandler instance.
func NewHealthHandler(cfg *config.Config, db *gorm.DB, redisClient *redis.Client) *HealthHandler {
	return &HealthHandler{
		config: cfg,
		db:     db,
		redis:  redisClient,
	}
}

// Check handles GET /health
// @Summary Health check
// @Description Returns the health status of the service
// @Tags health
// @Produce json
// @Success 200 {object} HealthStatus
// @Router /health [get]
func (h *HealthHandler) Check(c *gin.Context) {
	status := "healthy"
	checks := make(map[string]Check)

	// Check database connection
	dbStatus := "healthy"
	dbMessage := ""
	sqlDB, err := h.db.DB()
	if err != nil {
		dbStatus = "unhealthy"
		dbMessage = "failed to get database connection"
		status = "degraded"
	} else if err := sqlDB.Ping(); err != nil {
		dbStatus = "unhealthy"
		dbMessage = "database ping failed"
		status = "degraded"
	}
	checks["database"] = Check{
		Status:  dbStatus,
		Message: dbMessage,
	}

	// Check Redis connection (optional)
	if h.redis != nil {
		redisStatus := "healthy"
		redisMessage := ""
		if err := h.redis.Ping(c.Request.Context()).Err(); err != nil {
			redisStatus = "unhealthy"
			redisMessage = "redis ping failed"
			// Redis failure doesn't make service unhealthy, just degraded
			if status == "healthy" {
				status = "degraded"
			}
		}
		checks["redis"] = Check{
			Status:  redisStatus,
			Message: redisMessage,
		}
	}

	// Return 200 with status
	c.JSON(http.StatusOK, HealthStatus{
		Status:  status,
		Service: "blog-service",
		Version: "1.0.0",
		Checks:  checks,
	})
}
