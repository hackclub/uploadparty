package controllers

import (
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type HealthController struct {
	DB *gorm.DB
}

type HealthResponse struct {
	Status    string            `json:"status"`
	Timestamp time.Time         `json:"timestamp"`
	Version   string            `json:"version,omitempty"`
	Services  map[string]string `json:"services"`
}

func NewHealthController(db *gorm.DB) *HealthController {
	return &HealthController{DB: db}
}

func (h *HealthController) Health(c *gin.Context) {
	response := HealthResponse{
		Timestamp: time.Now().UTC(),
		Version:   os.Getenv("APP_VERSION"),
		Services:  make(map[string]string),
	}

	// Check database
	dbStatus := "healthy"
	if h.DB == nil {
		dbStatus = "not connected"
	} else {
		var one int
		if err := h.DB.Raw("SELECT 1").Scan(&one).Error; err != nil {
			dbStatus = "error: " + err.Error()
		}
	}
	response.Services["database"] = dbStatus

	// TODO: Add Redis health check when implemented
	// response.Services["redis"] = redisStatus

	// Determine overall status
	overall := "healthy"
	for _, status := range response.Services {
		if status != "healthy" {
			overall = "degraded"
			break
		}
	}
	response.Status = overall

	// Return appropriate HTTP status
	if overall == "healthy" {
		c.JSON(http.StatusOK, response)
	} else {
		c.JSON(http.StatusServiceUnavailable, response)
	}
}
