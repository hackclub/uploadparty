package controllers

import (
	"crypto/rand"
	"encoding/base32"
	"net/http"
	"net/mail"
	"strings"

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
	Email        string `json:"email" binding:"required,email"`
	FirstName    string `json:"firstName"`
	LastName     string `json:"lastName"`
	ReferralCode string `json:"referralCode,omitempty"` // The code used to sign up
}

// generateReferralCode creates a unique referral code
func generateReferralCode() string {
	b := make([]byte, 8)
	rand.Read(b)
	// Use base32 encoding (uppercase, no padding) for readable codes
	code := base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(b)
	// Take first 10 characters and make lowercase for friendlier URLs
	return strings.ToLower(code[:10])
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

	var referrerID *uint
	var referrer *models.RSVP
	// Validate referral code if provided
	if req.ReferralCode != "" {
		var ref models.RSVP
		if err := r.DB.Where("referral_code = ?", req.ReferralCode).First(&ref).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid referral code"})
			return
		}
		referrerID = &ref.ID
		referrer = &ref
	}

	// Get client IP address
	clientIP := c.ClientIP()

	// Generate unique referral code
	var referralCode string
	for {
		referralCode = generateReferralCode()
		// Check if code already exists
		var existing models.RSVP
		if err := r.DB.Where("referral_code = ?", referralCode).First(&existing).Error; err != nil {
			// Code doesn't exist, we can use it
			break
		}
		// Code exists, generate a new one
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
		Email:          req.Email,
		FirstName:      req.FirstName,
		LastName:       req.LastName,
		EmailSent:      emailSent,
		IPAddress:      clientIP,
		ReferralCode:   referralCode,
		ReferredByCode: req.ReferralCode,
		ReferredByID:   referrerID,
	}

	if err := r.DB.Create(&rsvp).Error; err != nil {
		if isDuplicateKeyError(err) {
			c.JSON(http.StatusConflict, gin.H{"error": "email already exists"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create RSVP"})
		return
	}

	// Send notification email to referrer if someone used their code
	if referrer != nil && r.EmailService != nil {
		// Build the new user's name
		newUserName := req.FirstName
		if req.LastName != "" {
			if newUserName != "" {
				newUserName += " " + req.LastName
			} else {
				newUserName = req.LastName
			}
		}
		if newUserName == "" {
			newUserName = req.Email
		}

		// Build the referrer's name
		referrerName := referrer.FirstName
		if referrer.LastName != "" {
			if referrerName != "" {
				referrerName += " " + referrer.LastName
			} else {
				referrerName = referrer.LastName
			}
		}
		if referrerName == "" {
			referrerName = "there"
		}

		// Send the notification (don't fail the request if email fails)
		if err := r.EmailService.SendReferralNotification(referrer.Email, referrerName, newUserName); err != nil {
			// Log but don't fail the request
			// The email service will already log the error
		}
	}

	// Get referral count for this user
	var referralCount int64
	r.DB.Model(&models.RSVP{}).Where("referred_by_id = ?", rsvp.ID).Count(&referralCount)

	c.JSON(http.StatusCreated, gin.H{
		"id":             rsvp.ID,
		"email":          rsvp.Email,
		"firstName":      rsvp.FirstName,
		"lastName":       rsvp.LastName,
		"emailSent":      rsvp.EmailSent,
		"ipAddress":      rsvp.IPAddress,
		"referralCode":   rsvp.ReferralCode,
		"referredByCode": rsvp.ReferredByCode,
		"referredById":   rsvp.ReferredByID,
		"referralCount":  referralCount,
		"message":        "RSVP successful! Please check your email for confirmation.",
	})
}

// GetReferrals returns the list of people referred by a specific RSVP
func (r *RSVPController) GetReferrals(c *gin.Context) {
	rsvpID := c.Param("id")

	var rsvp models.RSVP
	if err := r.DB.Preload("Referrals").First(&rsvp, rsvpID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "RSVP not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":            rsvp.ID,
		"email":         rsvp.Email,
		"firstName":     rsvp.FirstName,
		"lastName":      rsvp.LastName,
		"referralCode":  rsvp.ReferralCode,
		"referralCount": len(rsvp.Referrals),
		"referrals":     rsvp.Referrals,
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

// UpdateReferralCode allows users to change their referral code
func (r *RSVPController) UpdateReferralCode(c *gin.Context) {
	rsvpID := c.Param("id")

	var req struct {
		NewReferralCode string `json:"newReferralCode" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "newReferralCode is required"})
		return
	}

	// Validate the new referral code format (lowercase alphanumeric, 3-20 chars)
	if len(req.NewReferralCode) < 3 || len(req.NewReferralCode) > 20 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "referral code must be between 3 and 20 characters"})
		return
	}

	// Convert to lowercase
	newCode := strings.ToLower(req.NewReferralCode)

	// Check if the code contains only alphanumeric characters
	for _, char := range newCode {
		if !((char >= 'a' && char <= 'z') || (char >= '0' && char <= '9')) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "referral code can only contain lowercase letters and numbers"})
			return
		}
	}

	// Check if code is already taken
	var existing models.RSVP
	if err := r.DB.Where("referral_code = ?", newCode).First(&existing).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "referral code already taken"})
		return
	}

	// Find the RSVP to update
	var rsvp models.RSVP
	if err := r.DB.First(&rsvp, rsvpID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "RSVP not found"})
		return
	}

	// Store old code for response
	oldCode := rsvp.ReferralCode

	// Update the referral code
	if err := r.DB.Model(&rsvp).Update("referral_code", newCode).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update referral code"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":              rsvp.ID,
		"email":           rsvp.Email,
		"oldReferralCode": oldCode,
		"newReferralCode": newCode,
		"message":         "Referral code updated successfully",
	})
}

func isDuplicateKeyError(err error) bool {
	return err != nil && (err.Error() == "UNIQUE constraint failed: rsvps.email" ||
		err.Error() == "ERROR: duplicate key value violates unique constraint \"idx_rsvps_email\" (SQLSTATE 23505)")
}
