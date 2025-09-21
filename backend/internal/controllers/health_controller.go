package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type HealthController struct{ DB *gorm.DB }

func NewHealthController(db *gorm.DB) *HealthController { return &HealthController{DB: db} }

func (h *HealthController) Health(c *gin.Context) {
	var one int
	if err := h.DB.Raw("SELECT 1").Scan(&one).Error; err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"status": "degraded", "db": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
