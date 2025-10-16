package db

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/uploadparty/app/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func Connect(cfg *config.Config) (*gorm.DB, error) {
	if DB != nil {
		return DB, nil
	}

	var dsn string

	// Check if we're using Cloud SQL connection name (for Cloud SQL Auth Proxy)
	if cfg.CloudSQLConnectionName != "" {
		log.Printf("[DB] Using Cloud SQL connection name: %s", cfg.CloudSQLConnectionName)
		// For Cloud SQL Auth Proxy, use unix socket connection
		dsn = fmt.Sprintf("host=/cloudsql/%s user=%s password=%s dbname=%s sslmode=disable TimeZone=UTC",
			cfg.CloudSQLConnectionName, cfg.DBUser, cfg.DBPassword, cfg.DBName)
	} else {
		log.Printf("[DB] Using direct connection to: %s:%s", cfg.DBHost, cfg.DBPort)
		// Direct TCP connection (for local dev or public IP)
		dsn = fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=UTC",
			cfg.DBHost, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBPort, cfg.DBSSLMode)
	}

	log.Printf("[DB] Connecting with DSN: %s", sanitizeDSN(dsn))

	// Configure GORM with better logging for production
	gormConfig := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	}

	// In production, reduce log verbosity
	if os.Getenv("GIN_MODE") == "release" {
		gormConfig.Logger = logger.Default.LogMode(logger.Error)
	}

	// Use PreferSimpleProtocol to allow multi-statement SQL execution (needed for migrations with DO $$ ... $$ blocks)
	pgcfg := postgres.Config{
		DSN:                  dsn,
		PreferSimpleProtocol: true,
	}

	// Retry connection with exponential backoff
	var db *gorm.DB
	var err error

	for attempt := 1; attempt <= 5; attempt++ {
		db, err = gorm.Open(postgres.New(pgcfg), gormConfig)
		if err == nil {
			break
		}

		log.Printf("[DB] Connection attempt %d failed: %v", attempt, err)
		if attempt < 5 {
			backoff := time.Duration(attempt*attempt) * time.Second
			log.Printf("[DB] Retrying in %v...", backoff)
			time.Sleep(backoff)
		}
	}

	if err != nil {
		return nil, fmt.Errorf("failed to connect to database after 5 attempts: %w", err)
	}

	// Test the connection
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	// Configure connection pool
	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetConnMaxLifetime(5 * time.Minute)

	// Test ping
	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	DB = db
	log.Println("[DB] connected successfully")
	return db, nil
}

// CreateDatabaseIfNotExists creates the database if it doesn't exist (development only)
func CreateDatabaseIfNotExists(cfg *config.Config) error {
	if !cfg.IsDevelopment() {
		return nil // Only create DB in development
	}

	log.Println("[DEV] Checking if database exists...")

	// Connect to postgres database first to create the target database
	var dsn string
	if cfg.CloudSQLConnectionName != "" {
		dsn = fmt.Sprintf("host=/cloudsql/%s user=%s password=%s dbname=postgres sslmode=disable TimeZone=UTC",
			cfg.CloudSQLConnectionName, cfg.DBUser, cfg.DBPassword)
	} else {
		dsn = fmt.Sprintf("host=%s user=%s password=%s dbname=postgres port=%s sslmode=%s TimeZone=UTC",
			cfg.DBHost, cfg.DBUser, cfg.DBPassword, cfg.DBPort, cfg.DBSSLMode)
	}

	pgcfg := postgres.Config{
		DSN:                  dsn,
		PreferSimpleProtocol: true,
	}

	db, err := gorm.Open(postgres.New(pgcfg), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Error), // Less verbose
	})
	if err != nil {
		return fmt.Errorf("failed to connect to postgres database: %w", err)
	}

	// Check if database exists
	var exists bool
	err = db.Raw("SELECT EXISTS(SELECT datname FROM pg_catalog.pg_database WHERE datname = ?)", cfg.DBName).Scan(&exists).Error
	if err != nil {
		return fmt.Errorf("failed to check if database exists: %w", err)
	}

	if !exists {
		log.Printf("[DEV] Creating database: %s", cfg.DBName)
		err = db.Exec(fmt.Sprintf("CREATE DATABASE %s", cfg.DBName)).Error
		if err != nil {
			return fmt.Errorf("failed to create database: %w", err)
		}
		log.Printf("[DEV] Database %s created successfully", cfg.DBName)
	} else {
		log.Printf("[DEV] Database %s already exists", cfg.DBName)
	}

	return nil
}

// sanitizeDSN removes password from DSN for logging
func sanitizeDSN(dsn string) string {
	// Simple password removal for logging (not perfect but good enough)
	return "host=*** user=*** password=*** dbname=*** ..."
}
