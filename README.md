# UploadParty
A modern, music‑focused platform for uploading beats and audio, running community challenges.

## Website
- Live site: coming soon
- Note: Set the repository Description, Website, and Topics in GitHub → Settings → General and Settings → Topics for better discoverability.

## Tech stack
- Backend: Go (Gin), GORM, PostgreSQL, Redis, JWT auth, CORS, upload limits
- Frontend: Next.js 15 (App Router) + React 19, Tailwind CSS 4
- Storage/Infra: Google Cloud Storage (GCS), optional Nginx reverse proxy, Docker Compose for local

## Environment configuration (no hardcoded secrets)

### Secrets and external integrations (public repo safe)
- For optional external integrations (such as a license directory), see secrets/README.md for local/dev guidance.
- This public README intentionally omits provider details. Never commit secrets; use environment variables and the secrets/ folder for local files.

### Secrets folder (public repo safe)
- For local files like Google service account JSON or sample license DSNs/tokens, use the ./secrets folder.
- The folder is git-ignored by default; only README.md and *.example files are tracked.
- Docker Compose mounts ./secrets into the API container at /app/secrets (read-only) and defaults GOOGLE_APPLICATION_CREDENTIALS to /app/secrets/gcp-service-account.json.
- Examples are provided in secrets/*.example — copy them, fill locally, and never commit real secrets.

Quick start (local dev without Docker):
1) Copy backend/.env.example to backend/.env and edit values (especially DB_PASSWORD and JWT_SECRET).
2) From backend/: go run ./cmd/server
3) Health endpoint: http://localhost:8080/health

Hot reload for backend (Air):
- We use the maintained fork github.com/air-verse/air for Go hot reloading during development.
- Install once: make air-install  (runs: go install github.com/air-verse/air@latest)
- Run with reload: make api-air   (uses backend/.air.toml, loads backend/.env)
- Notes:
  - Ensure your DB/Redis are running (make db-up) and backend/.env is configured.
  - Air watches backend/ and rebuilds on changes; it ignores tmp, vendor, node_modules.

Docker Compose (multi-file):
- Base file: docker-compose.yml (shared: redis, frontend, nginx, networks, volumes)
- Development overrides: docker-compose.dev.yml (adds postgres-local, api wired to it)
- Production overrides: docker-compose.prod.yml (adds cloudsql-proxy, api wired to it)

Usage:
- Development:
  docker compose -f docker-compose.dev.yml up

- Production:
  docker compose -f docker-compose.prod.yml up

Notes:
- Dev API connects to postgres-local with DB_HOST=postgres-local and GIN_MODE=debug.
- Prod API connects via cloudsql-proxy with DB_HOST=cloudsql-proxy and GIN_MODE=release.
- Secrets and credentials should be provided via env/.env and the ./secrets folder (mounted read-only).

Database migrations:
- Migrations live under backend/migrations (e.g., 001_init.sql). Dev-only seed files use the suffix .dev.sql and are excluded by default.
- Rails-like local flow for juniors:
  1) Copy backend/.env.example to backend/.env and set DB_* for local (host=localhost, port=5432).
  2) Start local DB and Redis: make db-up
  3) Apply migrations: make migrate
  4) (Optional) Seed dev data: make migrate-dev (includes *.dev.sql like 090_dev_seed.dev.sql)
  5) Run API: make api
- Migration CLI (advanced):
  - From backend/: go run ./cmd/migrate
  - Flags:
    - -dir: path to migrations directory (default: migrations)
    - -list: list .sql files without executing
    - -dry-run: print SQL without executing
    - -env: set to dev to include *.dev.sql (or set MIGRATIONS_ENV=dev)
- You can build a binary too: from backend/: go build -o bin/migrate ./cmd/migrate
- Note: The old scripts/ folder is deprecated; migrations have been moved to backend/migrations.

Troubleshooting (Cloud SQL connection):
- Ensure the Cloud SQL Auth Proxy is running (docker ps shows uploadparty-cloudsql-proxy) and INSTANCE_CONNECTION_NAME is correct.
- With the proxy, DB_HOST should be cloudsql-proxy, DB_PORT=5432, and DB_SSL_MODE=disable (TLS is terminated by the proxy).
- For direct/public IP connections (no proxy), set DB_HOST to your instance address and DB_SSL_MODE=require or verify-full, and open the appropriate firewall rules.
- The server logs which .env file it loaded: [env] loaded <path>. If none is logged, ensure your .env exists in one of: .env, ../.env, ../../.env, backend/.env, ../backend/.env, ../../backend/.env.
- Verify that DB_USER/DB_PASSWORD have access to the target database in the Cloud SQL instance.

Security notes:
- If defaults are used for JWT_SECRET or DB_PASSWORD, the server logs a warning at startup.
- Do not commit real secrets. In production, set env vars through your platform (e.g., Coolify, Docker secrets, etc.).

## API routing separation
To make it clear which clients call which endpoints, API v1 is split by client type. Authentication uses the same JWT middleware for now.

- VST ingestion (plugin/DAW):
  - Base: /api/v1/ingest
  - POST /projects — Upsert project by title with heartbeat/metadata (used by VST)
  - POST /projects/:id/plugins — Upsert or attach plugin metadata to a project
  - PATCH /projects/:id/complete — Mark a project complete from the DAW

- Frontend application (Next.js):
  - Base: /api/v1/app
  - GET /projects — List my projects (includes attached plugins)
  - GET /projects/:id/plugins — List plugins for a project
  - PATCH /projects/:id/complete — Mark a project complete from the app

- Public (no auth):
  - GET /profiles/:handle — Public profile and public projects

