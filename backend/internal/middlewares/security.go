package middlewares

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

// SecurityHeaders adds security-related HTTP headers
func SecurityHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Prevent clickjacking
		c.Header("X-Frame-Options", "DENY")

		// Prevent MIME type sniffing
		c.Header("X-Content-Type-Options", "nosniff")

		// Enable XSS protection
		c.Header("X-XSS-Protection", "1; mode=block")

		// Force HTTPS (in production)
		c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains")

		// Content Security Policy
		csp := "default-src 'self'; " +
			"script-src 'self' 'unsafe-inline'; " +
			"style-src 'self' 'unsafe-inline'; " +
			"img-src 'self' data: https:; " +
			"media-src 'self' https:; " +
			"connect-src 'self' wss: https:; " +
			"frame-ancestors 'none';"
		c.Header("Content-Security-Policy", csp)

		// Referrer Policy
		c.Header("Referrer-Policy", "strict-origin-when-cross-origin")

		// Feature Policy / Permissions Policy
		c.Header("Permissions-Policy", "geolocation=(), microphone=(), camera=()")

		c.Next()
	}
}

// RateLimiter creates a rate limiting middleware
func RateLimiter(rps rate.Limit, burst int) gin.HandlerFunc {
	limiter := rate.NewLimiter(rps, burst)

	return func(c *gin.Context) {
		if !limiter.Allow() {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "Rate limit exceeded. Please try again later.",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}

// IPRateLimiter creates a per-IP rate limiting middleware
func IPRateLimiter(rps rate.Limit, burst int) gin.HandlerFunc {
	limiters := make(map[string]*rate.Limiter)

	return func(c *gin.Context) {
		ip := getClientIP(c)

		limiter, exists := limiters[ip]
		if !exists {
			limiter = rate.NewLimiter(rps, burst)
			limiters[ip] = limiter

			// Clean up old limiters periodically (simple approach)
			if len(limiters) > 10000 {
				// Clear half of the limiters (you might want a more sophisticated cleanup)
				for k := range limiters {
					delete(limiters, k)
					if len(limiters) <= 5000 {
						break
					}
				}
			}
		}

		if !limiter.Allow() {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "Rate limit exceeded. Please try again later.",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}

// FileUploadLimiter limits file upload size and type
func FileUploadLimiter(maxSize int64, allowedTypes []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method != "POST" && c.Request.Method != "PUT" {
			c.Next()
			return
		}

		// Check content length
		if c.Request.ContentLength > maxSize {
			c.JSON(http.StatusRequestEntityTooLarge, gin.H{
				"error": "File too large",
			})
			c.Abort()
			return
		}

		// Check content type for uploads
		contentType := c.GetHeader("Content-Type")
		if strings.Contains(contentType, "multipart/form-data") {
			// Parse multipart form
			if err := c.Request.ParseMultipartForm(maxSize); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": "Invalid multipart form",
				})
				c.Abort()
				return
			}

			// Check file types
			if c.Request.MultipartForm != nil && c.Request.MultipartForm.File != nil {
				for _, fileHeaders := range c.Request.MultipartForm.File {
					for _, fileHeader := range fileHeaders {
						if !isAllowedFileType(fileHeader.Header.Get("Content-Type"), allowedTypes) {
							c.JSON(http.StatusUnsupportedMediaType, gin.H{
								"error": "File type not allowed",
							})
							c.Abort()
							return
						}
					}
				}
			}
		}

		c.Next()
	}
}

// RequestSizeLimit limits the size of request bodies
func RequestSizeLimit(maxSize int64) gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxSize)
		c.Next()
	})
}

// Timeout middleware to prevent slow requests
func Timeout(timeout time.Duration) gin.HandlerFunc {
	return gin.TimeoutWithHandler(timeout, func(c *gin.Context) {
		c.JSON(http.StatusRequestTimeout, gin.H{
			"error": "Request timeout",
		})
	})
}

// Helper functions
func getClientIP(c *gin.Context) string {
	// Check for forwarded IP first (behind proxy)
	forwarded := c.GetHeader("X-Forwarded-For")
	if forwarded != "" {
		// Get first IP from comma-separated list
		ips := strings.Split(forwarded, ",")
		return strings.TrimSpace(ips[0])
	}

	// Check X-Real-IP header
	realIP := c.GetHeader("X-Real-IP")
	if realIP != "" {
		return realIP
	}

	// Fallback to RemoteAddr
	return c.ClientIP()
}

func isAllowedFileType(contentType string, allowedTypes []string) bool {
	for _, allowed := range allowedTypes {
		if strings.Contains(contentType, allowed) {
			return true
		}
	}
	return false
}
