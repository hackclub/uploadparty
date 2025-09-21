package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/uploadparty/app/backend/internal/services"
)

type ProjectController struct{ Svc *services.ProjectService }

func NewProjectController(db *gorm.DB) *ProjectController {
	return &ProjectController{Svc: services.NewProjectService(db)}
}

type upsertProjectReq struct {
	Title           string  `json:"title" binding:"required"`
	DAW             string  `json:"daw"`
	PluginVersion   string  `json:"pluginVersion"`
	DurationSeconds int     `json:"durationSeconds"`
	Metadata        jsonRaw `json:"metadata"`
	Public          *bool   `json:"public"`
}

type jsonRaw []byte

func (j *jsonRaw) UnmarshalJSON(b []byte) error { *j = append((*j)[0:0], b...); return nil }

func (p *ProjectController) Upsert(c *gin.Context) {
	var req upsertProjectReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	uid := c.GetUint("user_id")
	in := services.UpsertProjectInput{Title: req.Title, DAW: req.DAW, PluginVersion: req.PluginVersion, DurationSeconds: req.DurationSeconds, Metadata: []byte(req.Metadata), Public: req.Public}
	proj, err := p.Svc.UpsertByTitle(uid, in)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, proj)
}

func (p *ProjectController) MarkComplete(c *gin.Context) {
	uid := c.GetUint("user_id")
	idStr := c.Param("id")
	id64, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	proj, err := p.Svc.MarkComplete(uid, uint(id64))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, proj)
}

func (p *ProjectController) ListMine(c *gin.Context) {
	uid := c.GetUint("user_id")
	ps, err := p.Svc.ListByUser(uid)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, ps)
}
