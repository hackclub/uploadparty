package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/uploadparty/app/internal/services"
)

type PluginController struct{ Svc *services.PluginService }

func NewPluginController(db *gorm.DB) *PluginController {
	return &PluginController{Svc: services.NewPluginService(db)}
}

type upsertPluginReq struct {
	Name     string  `json:"name" binding:"required"`
	Vendor   string  `json:"vendor"`
	Version  string  `json:"version"`
	Format   string  `json:"format"`
	Metadata jsonRaw `json:"metadata"`
}

func (p *PluginController) UpsertForProject(c *gin.Context) {
	uid := c.GetUint("user_id")
	idStr := c.Param("id")
	id64, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	var req upsertPluginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	in := services.UpsertPluginInput{Name: req.Name, Vendor: req.Vendor, Version: req.Version, Format: req.Format, Metadata: []byte(req.Metadata)}
	pl, err := p.Svc.UpsertByName(uid, uint(id64), in)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, pl)
}

func (p *PluginController) ListByProject(c *gin.Context) {
	uid := c.GetUint("user_id")
	idStr := c.Param("id")
	id64, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	items, err := p.Svc.ListByProject(uid, uint(id64))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, items)
}
