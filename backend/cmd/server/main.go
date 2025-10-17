package main

import (
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"

	"github.com/uploadparty/app/config"
	"github.com/uploadparty/app/internal/controllers"
	"github.com/uploadparty/app/internal/integrations/licenses"
	"github.com/uploadparty/app/internal/middlewares"
	"github.com/uploadparty/app/internal/models"
	"github.com/uploadparty/app/internal/services"
	"github.com/uploadparty/app/pkg/db"
)

func main() {
	cfg := config.Load()

	// Print environment information
	log.Printf("=== UploadParty Backend Starting ===")
	log.Printf("Environment: %s", cfg.Environment)
	if cfg.IsProduction() {
		log.Printf("Running in PRODUCTION mode")
	} else {
		log.Printf("Running in DEVELOPMENT mode")
	}

	gin.SetMode(cfg.GinMode)
	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())

	corsCfg := cors.Config{
		AllowOrigins:     []string{cfg.FrontendURL},
		AllowMethods:     []string{"GET", "POST", "PATCH", "OPTIONS"},
		AllowHeaders:     []string{"Authorization", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}
	r.Use(cors.New(corsCfg))

	// Basic rate limit 5 rps with burst 10
	rl := middlewares.NewRateLimiter(rate.Limit(5), 10)
	go rl.Cleanup(10 * time.Minute)
	r.Use(rl.Middleware())

	// In development, create database if it doesn't exist
	if cfg.IsDevelopment() {
		if err := db.CreateDatabaseIfNotExists(cfg); err != nil {
			log.Printf("[DEV] Failed to create database: %v", err)
		}
	}

	// Database connection with better error handling for Cloud Run
	log.Printf("Attempting database connection...")
	database, err := db.Connect(cfg)
	if err != nil {
		log.Printf("[ERROR] Database connection failed: %v", err)
		log.Printf("[INFO] Continuing without database for health checks...")
		database = nil
	} else {
		// Ensure DB is actually reachable by executing a simple query.
		var one int
		if err := database.Raw("SELECT 1").Scan(&one).Error; err != nil {
			log.Printf("[ERROR] Database ping failed: %v", err)
			log.Printf("[INFO] Continuing without database for health checks...")
			database = nil
		} else {
			log.Printf("[SUCCESS] Database connection established")
			// Run migrations - required in development, optional in production
			if err := database.AutoMigrate(&models.RSVP{}, &models.User{}, &models.Project{}, &models.Plugin{}); err != nil {
				log.Printf("[ERROR] Database migration failed: %v", err)
				if cfg.IsDevelopment() {
					log.Fatalf("[DEV] Migration failure in development - exiting")
				} else {
					log.Printf("[INFO] Continuing without migrations in production...")
				}
			} else {
				log.Printf("[SUCCESS] Database migrations completed")
			}
		}
	}

	// Initialize external license store (generic). Fail closed to no-op if misconfigured.
	if err := licenses.Init(cfg); err != nil {
		log.Println("[licenses] init failed; external license lookups disabled")
	}

	jwt := middlewares.NewJWT(cfg.JWTSecret)
	auth0 := middlewares.NewAuth0(cfg.Auth0Domain, cfg.Auth0Audience)

	// Initialize email service
	emailService, err := services.NewEmailService(cfg)
	if err != nil {
		log.Printf("[EMAIL] Failed to initialize email service: %v", err)
		emailService = nil
	}

	healthCtl := controllers.NewHealthController(database)
	authCtl := controllers.NewAuthController(database, cfg.JWTSecret)
	projCtl := controllers.NewProjectController(database)
	pluginCtl := controllers.NewPluginController(database)
	profCtl := controllers.NewProfileController(database, cfg.JWTSecret)
	rsvpCtl := controllers.NewRSVPController(database, emailService)

	// Health
	r.GET("/health", healthCtl.Health)
	// Public alias for health under /api to work behind proxies
	r.GET("/api/health", healthCtl.Health)

	// Auth group (public endpoints)
	auth := r.Group("/auth")
	{
		auth.POST("/register", authCtl.Register)
		auth.POST("/login", authCtl.Login)
	}

	// RSVP (public endpoints)
	r.POST("/rsvp", rsvpCtl.Create)
	r.GET("/rsvp/count", rsvpCtl.Count)
	r.GET("/rsvp/:id/referrals", rsvpCtl.GetReferrals)
	r.PATCH("/rsvp/:id/referral-code", rsvpCtl.UpdateReferralCode)

	// Auth0 sync endpoint (protected by Auth0 JWT)
	// This endpoint is called by the frontend after Auth0 login to sync user info to our DB
	authSync := r.Group("/api/v1/auth")
	if cfg.Auth0Domain != "" {
		authSync.Use(auth0.RequireAuth0())
	}
	{
		authSync.POST("/sync", authCtl.SyncUser)
	}

	// API v1 protected (JWT). We now split routes by client type: /app (frontend) vs /ingest (VST/plugin)
	api := r.Group("/api/v1")
	api.Use(jwt.RequireAuth())
	{
		// --- Separated groups ---
		// VST/plugin ingestion endpoints: heartbeat/metadata and plugin upserts.
		ingest := api.Group("/ingest")
		{
			ingest.POST("/projects", projCtl.Upsert) // upsert by title; used by VST heartbeat/metadata capture
			ingest.POST("/projects/:id/plugins", pluginCtl.UpsertForProject)
			ingest.PATCH("/projects/:id/complete", projCtl.MarkComplete)
		}

		// Frontend application endpoints: listing, reading, user-triggered updates.
		app := api.Group("/app")
		{
			app.GET("/projects", projCtl.ListMine)
			app.GET("/projects/:id/plugins", pluginCtl.ListByProject)
			app.PATCH("/projects/:id/complete", projCtl.MarkComplete)
		}
	}

	// Public profiles
	r.GET("/profiles/:handle", profCtl.GetPublicProfile)

	// Serve static frontend files (for Cloud Run single-service deployment)
	r.Static("/static", "./app/.next/static")
	r.StaticFile("/favicon.ico", "./app/public/favicon.ico")
	// Catch-all for Next.js routes - serve index.html
	r.NoRoute(func(c *gin.Context) {
		path := c.Request.URL.Path
		// Only serve HTML for non-API routes
		if path != "/health" && !strings.HasPrefix(path, "/api/") && !strings.HasPrefix(path, "/auth/") {
			c.File("./app/public/index.html")
		} else {
			c.JSON(404, gin.H{"error": "route not found"})
		}
	})

	port := cfg.Port
	if p := os.Getenv("PORT"); p != "" {
		port = p
	}
	addr := ":" + port

	// Enhanced logging for Cloud Run debugging
	log.Printf("=== UploadParty Backend Starting ===")
	log.Printf("PORT environment variable: %s", os.Getenv("PORT"))
	log.Printf("Config port: %s", cfg.Port)
	log.Printf("Final port: %s", port)
	log.Printf("Listen address: %s", addr)
	log.Printf("GIN_MODE: %s", cfg.GinMode)
	log.Printf("Frontend URL: %s", cfg.FrontendURL)
	log.Printf("=== Server Starting ===")

	BackendURL := "http://localhost" + addr
	log.Println("Server listening on", BackendURL)

	// Bind to all interfaces for Cloud Run (0.0.0.0)
	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
