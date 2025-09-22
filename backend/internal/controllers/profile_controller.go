package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/uploadparty/app/internal/services"
)

type ProfileController struct {
	Users    *services.UserService
	Projects *services.ProjectService
}

func NewProfileController(db *gorm.DB, secret string) *ProfileController {
	return &ProfileController{Users: services.NewUserService(db, secret), Projects: services.NewProjectService(db)}
}

func (p *ProfileController) GetPublicProfile(c *gin.Context) {
	handle := c.Param("handle")
	u, err := p.Users.FindPublicByHandle(handle)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "profile not found"})
		return
	}
	projects, err := p.Projects.ListPublicByUser(u.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"user":     gin.H{"id": u.ID, "username": u.Username, "displayName": u.DisplayName, "bio": u.Bio},
		"projects": projects,
	})
}
