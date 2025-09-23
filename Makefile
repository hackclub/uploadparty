# Default environment file
ENV_FILE ?= backend/.env

.PHONY: help db-up db-down migrate migrate-dev migrate-list migrate-dry api api-air air-install dev prod deploy

help:
	@echo "Useful commands:"
	@echo "  make db-up         # Start local Postgres (postgres-local) and Redis"
	@echo "  make db-down       # Stop local Postgres and Redis"
	@echo "  make migrate       # Apply SQL migrations"
	@echo "  make migrate-dev   # Apply SQL migrations incl. *.dev.sql (seeds)"
	@echo "  make migrate-list  # List migrations that would run"
	@echo "  make migrate-dry   # Print SQL without executing"
	@echo "  make api           # Run the Go API (from backend/)"
	@echo "  make api-air       # Run the Go API with Air hot reload (from backend/)"
	@echo "  make air-install   # Install air-verse/air locally"

# Local DB using docker compose
 db-up:
	docker compose up -d postgres-local redis

 db-down:
	docker compose stop postgres-local redis

# Migrations (run from backend/ so config loads backend/.env automatically)
 migrate:
	cd backend && go run ./cmd/migrate

 migrate-dev:
	cd backend && MIGRATIONS_ENV=dev go run ./cmd/migrate -env=dev

 migrate-list:
	cd backend && go run ./cmd/migrate -list

 migrate-dry:
	cd backend && go run ./cmd/migrate -dry-run

# Run API locally
 api:
	cd backend && go run ./cmd/server

# Hot reload using Air (air-verse)
 air-install:
	go install github.com/air-verse/air@latest

 api-air:
	cd backend && air -c .air.toml

# Compose wrappers
# Development environment
 dev:
	docker compose -f docker-compose.dev.yml up

# Production environment
 prod:
	docker compose -f docker-compose.prod.yml up -d

# Deploy alias (uses prod)
 deploy: prod
