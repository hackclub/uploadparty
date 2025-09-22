package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string

	Port        string
	GinMode     string
	FrontendURL string
	JWTSecret   string

	// Storage (Google Cloud Storage)
	GCPProjectID               string
	GCSBucket                  string
	GoogleApplicationCredsPath string // GOOGLE_APPLICATION_CREDENTIALS

	// External license directory (generic, provider may be hidden)
	LicensesProvider string // e.g., "airtable" or "none"
	LicensesToken    string // generic bearer token or API key (keep secure)
	LicensesDSN      string // opaque DSN string, e.g., "base=...;table=..."
}

func Load() *Config {
	// Load environment variables from .env files if present. We try common locations
	// and ignore errors so production can rely on injected environment variables.
	loaded := tryLoadEnv(
		".env",
		"../.env",
		"../../.env",
		"backend/.env",
		"../backend/.env",
		"../../backend/.env",
	)
	if loaded != "" {
		log.Printf("[env] loaded %s", loaded)
	}

	cfg := &Config{
		DBHost:      getEnv("DB_HOST", "localhost"),
		DBPort:      getEnv("DB_PORT", "5432"),
		DBUser:      getEnv("DB_USER", "uploadparty"),
		DBPassword:  getEnv("DB_PASSWORD", "your_local_db_password"),
		DBName:      getEnv("DB_NAME", "uploadparty_db"),
		DBSSLMode:   getEnv("DB_SSL_MODE", "disable"),
		Port:        getEnv("PORT", "8080"),
		GinMode:     getEnv("GIN_MODE", "debug"),
		FrontendURL: getEnv("FRONTEND_URL", "http://localhost:3000"),
		JWTSecret:   getEnv("JWT_SECRET", "change_me"),
		// Storage (GCS)
		GCPProjectID:               getEnv("GCP_PROJECT_ID", ""),
		GCSBucket:                  getEnv("GCS_BUCKET", "uploadparty-beats"),
		GoogleApplicationCredsPath: getEnv("GOOGLE_APPLICATION_CREDENTIALS", ""),
		// Licenses (generic)
		LicensesProvider: getEnv("LICENSES_PROVIDER", "none"),
		LicensesToken:    getEnv("LICENSES_TOKEN", ""),
		LicensesDSN:      getEnv("LICENSES_DSN", ""),
	}
	if cfg.JWTSecret == "change_me" {
		log.Println("[WARN] Using default JWT secret; set JWT_SECRET in env for non-dev")
	}
	if cfg.DBPassword == "your_local_db_password" || cfg.DBPassword == "" {
		log.Println("[WARN] Using default or empty DB password; set DB_PASSWORD in env for non-dev and production")
	}
	if cfg.GCSBucket == "" {
		log.Println("[WARN] GCS bucket not set; set GCS_BUCKET for media storage when enabling uploads")
	}
	if cfg.GoogleApplicationCredsPath == "" {
		log.Println("[WARN] GOOGLE_APPLICATION_CREDENTIALS not set; Google Cloud SDK default credentials will be used if available")
	}
	if cfg.LicensesProvider != "none" && (cfg.LicensesToken == "" || cfg.LicensesDSN == "") {
		log.Println("[WARN] LICENSES_PROVIDER set but LICENSES_TOKEN or LICENSES_DSN is missing; license lookups will be disabled")
	}
	return cfg
}

func tryLoadEnv(paths ...string) string {
	for _, p := range paths {
		if err := godotenv.Load(p); err == nil {
			return p
		}
	}
	return ""
}

func getEnv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
