package monitoring

import (
	"context"
	"net/http"
	"time"

	"go_byteeats/internal/database"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type HealthChecker struct {
	db    *gorm.DB
	redis *database.RedisClient
}

type HealthStatus struct {
	Status    string            `json:"status"`
	Timestamp time.Time         `json:"timestamp"`
	Services  map[string]string `json:"services"`
}

func NewHealthChecker(db *gorm.DB, redis *database.RedisClient) *HealthChecker {
	return &HealthChecker{
		db:    db,
		redis: redis,
	}
}

// RegisterHealthEndpoints registers health check endpoints
func (h *HealthChecker) RegisterHealthEndpoints(router *gin.Engine) {
	router.GET("/health", h.healthCheck)
	router.GET("/health/live", h.livenessCheck)
	router.GET("/health/ready", h.readinessCheck)
}

// healthCheck performs a basic health check
func (h *HealthChecker) healthCheck(c *gin.Context) {
	status := HealthStatus{
		Status:    "ok",
		Timestamp: time.Now(),
		Services:  make(map[string]string),
	}

	c.JSON(http.StatusOK, status)
}

// livenessCheck checks if the application is running
func (h *HealthChecker) livenessCheck(c *gin.Context) {
	status := HealthStatus{
		Status:    "ok",
		Timestamp: time.Now(),
		Services:  make(map[string]string),
	}

	c.JSON(http.StatusOK, status)
}

// readinessCheck checks if the application is ready to serve traffic
func (h *HealthChecker) readinessCheck(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	status := HealthStatus{
		Status:    "ok",
		Timestamp: time.Now(),
		Services:  make(map[string]string),
	}

	// Check database connection
	if err := h.db.WithContext(ctx).Raw("SELECT 1").Error; err != nil {
		status.Status = "error"
		status.Services["database"] = "error: " + err.Error()
	} else {
		status.Services["database"] = "ok"
	}

	// Check Redis connection
	if err := h.redis.Get(ctx, "health_check", nil); err != nil {
		status.Status = "error"
		status.Services["redis"] = "error: " + err.Error()
	} else {
		status.Services["redis"] = "ok"
	}

	if status.Status == "error" {
		c.JSON(http.StatusServiceUnavailable, status)
		return
	}

	c.JSON(http.StatusOK, status)
}
