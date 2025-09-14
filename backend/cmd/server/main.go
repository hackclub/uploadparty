package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"uploadparty/config"
	"uploadparty/internal/controllers"
	"uploadparty/internal/middlewares"
	"uploadparty/internal/services"
	"uploadparty/pkg/db"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Set Gin mode
	gin.SetMode(cfg.Server.GinMode)

	// Connect to database
	database, err := db.Connect(cfg)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Auto migrate database schemas
	if err := db.AutoMigrate(database); err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	// Initialize services
	authService := services.NewAuthService(database, cfg.JWT.Secret)

	// Initialize controllers
	authController := controllers.NewAuthController(authService)

	// Setup router with security middlewares
	router := gin.New()

	// Recovery middleware
	router.Use(gin.Recovery())

	// Custom logger
	router.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		return fmt.Sprintf("%s - [%s] \"%s %s %s %d %s \"%s\" %s\"\n",
			param.ClientIP,
			param.TimeStamp.Format(time.RFC1123),
			param.Method,
			param.Path,
			param.Request.Proto,
			param.StatusCode,
			param.Latency,
			param.Request.UserAgent(),
			param.ErrorMessage,
		)
	}))

	// Security middlewares
	router.Use(middlewares.SecurityHeaders())
	router.Use(middlewares.CORSMiddleware(cfg.Server.FrontendURL))
	router.Use(middlewares.RequestSizeLimit(50 << 20)) // 50MB limit
	router.Use(middlewares.Timeout(30 * time.Second))

	// Rate limiting
	router.Use(middlewares.IPRateLimiter(rate.Limit(10), 20)) // 10 requests per second, burst of 20

	// File upload limits for specific routes
	audioTypes := []string{"audio/mpeg", "audio/wav", "audio/mp3", "audio/x-wav"}
	uploadLimiter := middlewares.FileUploadLimiter(100<<20, audioTypes) // 100MB for audio files

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":    "healthy",
			"timestamp": time.Now().Unix(),
		})
	})

	// API routes
	api := router.Group("/api/v1")
	{
		// Auth routes (no authentication required)
		auth := api.Group("/auth")
		{
			auth.POST("/register", authController.Register)
			auth.POST("/login", authController.Login)
		}

		// Protected routes (authentication required)
		protected := api.Group("")
		protected.Use(middlewares.AuthMiddleware(cfg.JWT.Secret))
		{
			// User profile
			protected.GET("/profile", authController.GetProfile)

			// Beat routes will be added here
			// beats := protected.Group("/beats")
			// {
			//     beats.POST("/upload", uploadLimiter, beatController.Upload)
			//     beats.GET("", beatController.GetBeats)
			//     beats.GET("/:id", beatController.GetBeat)
			// }
		}

		// Public routes (optional authentication)
		public := api.Group("")
		public.Use(middlewares.OptionalAuthMiddleware(cfg.JWT.Secret))
		{
			// Public beat browsing, challenge listings, etc.
		}
	}

	// Setup HTTP server with security configurations
	srv := &http.Server{
		Addr:           ":" + cfg.Server.Port,
		Handler:        router,
		ReadTimeout:    15 * time.Second,
		WriteTimeout:   15 * time.Second,
		IdleTimeout:    60 * time.Second,
		MaxHeaderBytes: 1 << 20, // 1MB
	}

	// Start server in a goroutine
	go func() {
		log.Printf("Server starting on port %s", cfg.Server.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// Give the server 5 seconds to finish current requests
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exited")
}
