package middlewares

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// AdminMiddleware ensures the user has admin privileges
func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
			c.Abort()
			return
		}

		email, exists := c.Get("email")
		if !exists {
			c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
			c.Abort()
			return
		}

		// Simple admin check - you can enhance this with role-based permissions
		emailStr, ok := email.(string)
		if !ok {
			c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
			c.Abort()
			return
		}

		// Check if user is admin (you can customize this logic)
		if !isAdminEmail(emailStr) {
			c.JSON(http.StatusForbidden, gin.H{"error": "Admin access required"})
			c.Abort()
			return
		}

		// Set admin context
		c.Set("is_admin", true)
		c.Next()
	}
}

// SuperAdminMiddleware ensures the user has super admin privileges
func SuperAdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
			c.Abort()
			return
		}

		email, exists := c.Get("email")
		if !exists {
			c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
			c.Abort()
			return
		}

		emailStr, ok := email.(string)
		if !ok {
			c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
			c.Abort()
			return
		}

		// Check if user is super admin
		if !isSuperAdminEmail(emailStr) {
			c.JSON(http.StatusForbidden, gin.H{"error": "Super admin access required"})
			c.Abort()
			return
		}

		c.Set("is_super_admin", true)
		c.Next()
	}
}

// Helper functions for admin checks
func isAdminEmail(email string) bool {
	// Define admin email patterns or domains
	adminDomains := []string{
		"admin@uploadparty.com",
		"support@uploadparty.com",
	}

	// Check exact matches
	for _, adminEmail := range adminDomains {
		if strings.EqualFold(email, adminEmail) {
			return true
		}
	}

	// Check admin domain patterns
	adminPatterns := []string{
		"@uploadparty.com",
		"@admin.uploadparty.com",
	}

	for _, pattern := range adminPatterns {
		if strings.HasSuffix(strings.ToLower(email), pattern) {
			return true
		}
	}

	return false
}

func isSuperAdminEmail(email string) bool {
	superAdminEmails := []string{
		"admin@uploadparty.com",
		"founder@uploadparty.com",
	}

	for _, adminEmail := range superAdminEmails {
		if strings.EqualFold(email, adminEmail) {
			return true
		}
	}

	return false
}

// AuditMiddleware logs admin actions for security
func AuditMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Store request info for audit logging
		c.Set("request_ip", c.ClientIP())
		c.Set("user_agent", c.GetHeader("User-Agent"))
		c.Set("request_method", c.Request.Method)
		c.Set("request_path", c.Request.URL.Path)

		c.Next()

		// Log admin action after request completion
		if c.GetBool("is_admin") {
			go logAdminAction(c)
		}
	}
}

func logAdminAction(c *gin.Context) {
	// This would typically log to database or audit service
	// For now, we'll just log to console
	// You should implement proper audit logging here

	userID := c.GetString("user_id")
	method := c.GetString("request_method")
	path := c.GetString("request_path")
	ip := c.GetString("request_ip")

	// Log format: [ADMIN_ACTION] UserID: xxx | Action: METHOD /path | IP: xxx
	if method != "GET" { // Only log non-read actions
		// fmt.Printf("[ADMIN_ACTION] UserID: %s | Action: %s %s | IP: %s\n", userID, method, path, ip)
		_ = userID // Prevent unused variable warning
		_ = path
		_ = ip
	}
}
