package utils

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
)

// StandardError represents a structured API error response
type StandardError struct {
	Error   string `json:"error"`
	Code    string `json:"code,omitempty"`
	Details any    `json:"details,omitempty"`
}

// ErrorHandler provides consistent error response handling
func ErrorHandler() gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		c.Next()

		// Process any errors that occurred during request handling
		if len(c.Errors) > 0 {
			err := c.Errors.Last()

			// Log error with correlation ID if available
			correlationID, _ := c.Get("correlation_id")
			slog.Error("Request error",
				"correlation_id", correlationID,
				"error", err.Error(),
				"path", c.Request.URL.Path,
				"method", c.Request.Method,
			)

			// Return structured error response
			c.JSON(http.StatusInternalServerError, StandardError{
				Error: "Internal server error",
				Code:  "INTERNAL_ERROR",
			})
			return
		}
	})
}

// BadRequest creates a standardized 400 error response
func BadRequest(c *gin.Context, message string, details ...any) {
	response := StandardError{
		Error: message,
		Code:  "BAD_REQUEST",
	}
	if len(details) > 0 {
		response.Details = details[0]
	}
	c.JSON(http.StatusBadRequest, response)
}

// Unauthorized creates a standardized 401 error response
func Unauthorized(c *gin.Context, message string) {
	c.JSON(http.StatusUnauthorized, StandardError{
		Error: message,
		Code:  "UNAUTHORIZED",
	})
}

// Forbidden creates a standardized 403 error response
func Forbidden(c *gin.Context, message string) {
	c.JSON(http.StatusForbidden, StandardError{
		Error: message,
		Code:  "FORBIDDEN",
	})
}

// NotFound creates a standardized 404 error response
func NotFound(c *gin.Context, message string) {
	c.JSON(http.StatusNotFound, StandardError{
		Error: message,
		Code:  "NOT_FOUND",
	})
}

// InternalError creates a standardized 500 error response
func InternalError(c *gin.Context, message string) {
	correlationID, _ := c.Get("correlation_id")
	slog.Error("Internal server error",
		"correlation_id", correlationID,
		"message", message,
		"path", c.Request.URL.Path,
	)

	c.JSON(http.StatusInternalServerError, StandardError{
		Error: "Internal server error",
		Code:  "INTERNAL_ERROR",
	})
}
