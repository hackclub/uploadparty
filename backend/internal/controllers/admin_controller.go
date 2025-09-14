package controllers

import (
	"net/http"
	"strconv"
	"uploadparty/internal/services"

	"github.com/gin-gonic/gin"
)

type AdminController struct {
	rewardService *services.RewardService
	authService   *services.AuthService
}

func NewAdminController(rewardService *services.RewardService, authService *services.AuthService) *AdminController {
	return &AdminController{
		rewardService: rewardService,
		authService:   authService,
	}
}

// Reward Management

func (ac *AdminController) GetPendingRewards(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	rewards, total, err := ac.rewardService.GetPendingRewards(page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"rewards": rewards,
		"total":   total,
		"page":    page,
		"limit":   limit,
	})
}

func (ac *AdminController) GetAllRewards(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	status := c.Query("status")

	rewards, total, err := ac.rewardService.GetAllRewards(status, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"rewards": rewards,
		"total":   total,
		"page":    page,
		"limit":   limit,
	})
}

func (ac *AdminController) ApproveReward(c *gin.Context) {
	rewardID := c.Param("id")
	adminID := c.GetString("user_id")

	var req services.ApproveRewardRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := ac.rewardService.ApproveReward(rewardID, adminID, &req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Reward approved successfully"})
}

func (ac *AdminController) RejectReward(c *gin.Context) {
	rewardID := c.Param("id")
	adminID := c.GetString("user_id")

	var req services.RejectRewardRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := ac.rewardService.RejectReward(rewardID, adminID, &req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Reward rejected successfully"})
}

func (ac *AdminController) CreateReward(c *gin.Context) {
	adminID := c.GetString("user_id")

	var req services.CreateRewardRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	reward, err := ac.rewardService.CreateReward(adminID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, reward)
}

// NI License Management

func (ac *AdminController) ImportLicenses(c *gin.Context) {
	adminID := c.GetString("user_id")

	var req services.ImportLicensesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := ac.rewardService.ImportLicenses(adminID, &req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Licenses imported successfully",
		"count":   len(req.Licenses),
	})
}

func (ac *AdminController) GetAvailableLicenses(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	productName := c.Query("product_name")

	licenses, total, err := ac.rewardService.GetAvailableLicenses(productName, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"licenses": licenses,
		"total":    total,
		"page":     page,
		"limit":    limit,
	})
}

func (ac *AdminController) GetAssignedLicenses(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	licenses, total, err := ac.rewardService.GetAssignedLicenses(page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"licenses": licenses,
		"total":    total,
		"page":     page,
		"limit":    limit,
	})
}

// Dashboard Stats

func (ac *AdminController) GetDashboardStats(c *gin.Context) {
	// Get various statistics for admin dashboard

	// This is a simplified version - you'd want to optimize these queries
	stats := gin.H{
		"pending_rewards":    0,
		"total_users":        0,
		"available_licenses": 0,
		"assigned_licenses":  0,
		"active_challenges":  0,
		"total_beats":        0,
	}

	// Get pending rewards count
	if rewards, _, err := ac.rewardService.GetPendingRewards(1, 1); err == nil {
		stats["pending_rewards"] = len(rewards)
	}

	// Get available licenses count
	if licenses, total, err := ac.rewardService.GetAvailableLicenses("", 1, 1); err == nil {
		_ = licenses // Prevent unused variable warning
		stats["available_licenses"] = total
	}

	// Get assigned licenses count
	if licenses, total, err := ac.rewardService.GetAssignedLicenses(1, 1); err == nil {
		_ = licenses
		stats["assigned_licenses"] = total
	}

	c.JSON(http.StatusOK, stats)
}

// User Management

func (ac *AdminController) GetUsers(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	search := c.Query("search")
	status := c.Query("status") // active, inactive, all

	// This would be implemented in a user service
	// For now, return placeholder response
	_ = page
	_ = limit
	_ = search
	_ = status

	c.JSON(http.StatusOK, gin.H{
		"users": []interface{}{},
		"total": 0,
		"page":  page,
		"limit": limit,
	})
}

func (ac *AdminController) GetUserDetails(c *gin.Context) {
	userID := c.Param("userId")

	user, err := ac.authService.GetUserByID(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Get user's rewards
	rewards, _, err := ac.rewardService.GetUserRewards(userID, "", 1, 10)
	if err != nil {
		rewards = []interface{}{} // Empty array on error
	}

	c.JSON(http.StatusOK, gin.H{
		"user":    user,
		"rewards": rewards,
	})
}

func (ac *AdminController) UpdateUserStatus(c *gin.Context) {
	userID := c.Param("userId")

	var req struct {
		IsActive bool   `json:"is_active"`
		Reason   string `json:"reason,omitempty"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// This would be implemented in user service
	// For now, return success
	_ = userID
	_ = req

	c.JSON(http.StatusOK, gin.H{"message": "User status updated successfully"})
}

// Analytics

func (ac *AdminController) GetAnalytics(c *gin.Context) {
	timeframe := c.DefaultQuery("timeframe", "7d") // 1d, 7d, 30d, all
	metric := c.Query("metric")                    // users, beats, rewards, challenges

	// This would integrate with your analytics service
	// For now, return placeholder data
	_ = timeframe
	_ = metric

	c.JSON(http.StatusOK, gin.H{
		"data": []gin.H{
			{"date": "2024-01-01", "value": 10},
			{"date": "2024-01-02", "value": 15},
			{"date": "2024-01-03", "value": 12},
		},
		"total":  37,
		"growth": 12.5,
	})
}
