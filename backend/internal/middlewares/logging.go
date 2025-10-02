package middlewares

import (
	"bytes"
	"io"
	"log/slog"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type responseBodyWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (r responseBodyWriter) Write(b []byte) (int, error) {
	r.body.Write(b)
	return r.ResponseWriter.Write(b)
}

// StructuredLogging provides comprehensive request/response logging with correlation IDs
func StructuredLogging() gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		// Generate correlation ID for request tracing
		correlationID := uuid.New().String()
		c.Header("X-Correlation-ID", correlationID)
		c.Set("correlation_id", correlationID)

		start := time.Now()
		path := c.Request.URL.Path
		method := c.Request.Method

		// Read request body for logging (if needed)
		var requestBody []byte
		if c.Request.Body != nil {
			requestBody, _ = io.ReadAll(c.Request.Body)
			c.Request.Body = io.NopCloser(bytes.NewBuffer(requestBody))
		}

		// Wrap response writer to capture response body
		w := &responseBodyWriter{
			ResponseWriter: c.Writer,
			body:           bytes.NewBufferString(""),
		}
		c.Writer = w

		// Process request
		c.Next()

		// Calculate duration
		duration := time.Since(start)
		statusCode := c.Writer.Status()

		// Log request details
		logLevel := slog.LevelInfo
		if statusCode >= 400 {
			logLevel = slog.LevelError
		} else if statusCode >= 300 {
			logLevel = slog.LevelWarn
		}

		slog.Log(c.Request.Context(), logLevel, "HTTP Request",
			"correlation_id", correlationID,
			"method", method,
			"path", path,
			"status_code", statusCode,
			"duration_ms", duration.Milliseconds(),
			"client_ip", c.ClientIP(),
			"user_agent", c.Request.UserAgent(),
			"request_size", len(requestBody),
			"response_size", w.body.Len(),
		)

		// Log errors if any
		if len(c.Errors) > 0 {
			for _, err := range c.Errors {
				slog.Error("Request error",
					"correlation_id", correlationID,
					"error", err.Error(),
					"type", err.Type,
				)
			}
		}
	})
}
