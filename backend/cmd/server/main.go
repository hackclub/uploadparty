package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"

	"github.com/uploadparty/app/config"
	"github.com/uploadparty/app/internal/controllers"
	"github.com/uploadparty/app/internal/integrations/licenses"
	"github.com/uploadparty/app/internal/middlewares"
	"github.com/uploadparty/app/internal/models"
	"github.com/uploadparty/app/pkg/db"
)

func main() {
	cfg := config.Load()
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

	database, err := db.Connect(cfg)
	if err != nil {
		log.Fatalf("db connect error: %v", err)
	}
	if err := database.AutoMigrate(&models.User{}, &models.Project{}, &models.Plugin{}); err != nil {
		log.Fatalf("db migrate error: %v", err)
	}

	// Initialize external license store (generic). Fail closed to no-op if misconfigured.
	if err := licenses.Init(cfg); err != nil {
		log.Println("[licenses] init failed; external license lookups disabled")
	}

	jwt := middlewares.NewJWT(cfg.JWTSecret)

	healthCtl := controllers.NewHealthController(database)
	authCtl := controllers.NewAuthController(database, cfg.JWTSecret)
	projCtl := controllers.NewProjectController(database)
	pluginCtl := controllers.NewPluginController(database)
	profCtl := controllers.NewProfileController(database, cfg.JWTSecret)

	// Health
	r.GET("/health", healthCtl.Health)
	// Public alias for health under /api to work behind proxies
	r.GET("/api/health", healthCtl.Health)

	// Auth group
	auth := r.Group("/auth")
	{
		auth.POST("/register", authCtl.Register)
		auth.POST("/login", authCtl.Login)
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

	port := cfg.Port
	if p := os.Getenv("PORT"); p != "" {
		port = p
	}
	addr := ":" + port
	log.Printf("listening on %s", addr)
	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatal(err)
	}
}
