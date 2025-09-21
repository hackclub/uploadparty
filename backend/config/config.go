package config

import (
	"log"
	"os"
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
}

func Load() *Config {
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
	}
	if cfg.JWTSecret == "change_me" {
		log.Println("[WARN] Using default JWT secret; set JWT_SECRET in env for non-dev")
	}
	return cfg
}

func getEnv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
