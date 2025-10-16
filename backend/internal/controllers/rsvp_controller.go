package controllers

import (
	"net/http"
	"net/mail"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/uploadparty/app/internal/models"
	"github.com/uploadparty/app/internal/services"
)

type RSVPController struct {
	DB           *gorm.DB
	EmailService *services.EmailService
}

func NewRSVPController(db *gorm.DB, emailService *services.EmailService) *RSVPController {
	return &RSVPController{
		DB:           db,
		EmailService: emailService,
	}
}

type rsvpReq struct {
	Email string `json:"email" binding:"required,email"`
}

func (r *RSVPController) Create(c *gin.Context) {
	var req rsvpReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate email using net/mail.ParseAddress
	_, err := mail.ParseAddress(req.Email)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "this email doesn't exist"})
		return
	}

	// Determine if email will be sent
	emailSent := false
	if r.EmailService != nil {
		err := r.EmailService.SendRSVPConfirmation(req.Email)
		if err == nil {
			emailSent = true
		}
	}

	rsvp := models.RSVP{
		Email:     req.Email,
		EmailSent: emailSent,
	}

	if err := r.DB.Create(&rsvp).Error; err != nil {
		if isDuplicateKeyError(err) {
			c.JSON(http.StatusConflict, gin.H{"error": "email already exists"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create RSVP"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id":        rsvp.ID,
		"email":     rsvp.Email,
		"emailSent": rsvp.EmailSent,
		"message":   "RSVP successful! Please check your email for confirmation.",
	})
}

func (r *RSVPController) Count(c *gin.Context) {
	var count int64
	if err := r.DB.Model(&models.RSVP{}).Count(&count).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get RSVP count"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"count":   count,
		"message": "people have RSVP'd",
	})
}

func isDuplicateKeyError(err error) bool {
	return err != nil && (err.Error() == "UNIQUE constraint failed: rsvps.email" ||
		err.Error() == "ERROR: duplicate key value violates unique constraint \"idx_rsvps_email\" (SQLSTATE 23505)")
}
