package tests

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/uploadparty/app/internal/controllers"
)

func setupTestDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to test database")
	}
	return db
}

func TestHealthController_Health(t *testing.T) {
	tests := []struct {
		name           string
		setupDB        func() *gorm.DB
		expectedStatus int
		expectedFields []string
	}{
		{
			name: "healthy database connection",
			setupDB: func() *gorm.DB {
				return setupTestDB()
			},
			expectedStatus: http.StatusOK,
			expectedFields: []string{"status", "timestamp", "services"},
		},
		{
			name: "nil database connection",
			setupDB: func() *gorm.DB {
				return nil
			},
			expectedStatus: http.StatusServiceUnavailable,
			expectedFields: []string{"status", "timestamp", "services"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			gin.SetMode(gin.TestMode)
			router := gin.New()

			db := tt.setupDB()
			healthController := controllers.NewHealthController(db)

			router.GET("/health", healthController.Health)

			// Create request
			req, _ := http.NewRequest("GET", "/health", nil)
			w := httptest.NewRecorder()

			// Execute
			router.ServeHTTP(w, req)

			// Assert status code
			assert.Equal(t, tt.expectedStatus, w.Code)

			// Assert response structure
			var response map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)

			// Check required fields
			for _, field := range tt.expectedFields {
				assert.Contains(t, response, field, "Response should contain field: %s", field)
			}

			// Verify services section
			if services, ok := response["services"].(map[string]interface{}); ok {
				assert.Contains(t, services, "database", "Services should contain database status")
			}
		})
	}
}

func BenchmarkHealthCheck(b *testing.B) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	db := setupTestDB()
	healthController := controllers.NewHealthController(db)
	router.GET("/health", healthController.Health)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		req, _ := http.NewRequest("GET", "/health", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
	}
}
