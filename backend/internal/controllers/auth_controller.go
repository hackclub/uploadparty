package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/uploadparty/app/internal/services"
)

type AuthController struct {
	Users *services.UserService
}

func NewAuthController(db *gorm.DB, secret string) *AuthController {
	return &AuthController{Users: services.NewUserService(db, secret)}
}

type registerReq struct {
	Email    string `json:"email" binding:"required,email"`
	Username string `json:"username" binding:"required,alphanum,min=3,max=20"`
	Password string `json:"password" binding:"required,min=6"`
}

type loginReq struct {
	Identifier string `json:"identifier" binding:"required"` // email or username
	Password   string `json:"password" binding:"required"`
}

func (a *AuthController) Register(c *gin.Context) {
	var req registerReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	u, err := a.Users.Register(req.Email, req.Username, req.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"id": u.ID, "email": u.Email, "username": u.Username})
}

func (a *AuthController) Login(c *gin.Context) {
	var req loginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	token, user, err := a.Users.Authenticate(req.Identifier, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": token, "user": gin.H{"id": user.ID, "email": user.Email, "username": user.Username}})
}

type syncUserReq struct {
	Auth0ID     string `json:"auth0_id" binding:"required"`
	Email       string `json:"email" binding:"required,email"`
	Username    string `json:"username" binding:"required"`
	DisplayName string `json:"display_name"`
	Picture     string `json:"picture"`
}

// SyncUser creates or updates a user from Auth0 authentication
func (a *AuthController) SyncUser(c *gin.Context) {
	var req syncUserReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := a.Users.SyncAuth0User(req.Auth0ID, req.Email, req.Username, req.DisplayName, req.Picture)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to sync user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":          user.ID,
		"auth0Id":     user.Auth0ID,
		"email":       user.Email,
		"username":    user.Username,
		"displayName": user.DisplayName,
		"picture":     user.Picture,
	})
}
