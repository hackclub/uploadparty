# UploadParty – Project Guidelines for Junie (Full‑Stack: Go Gin + Next.js)

These guidelines tell Junie how to work within this repository as a full‑stack engineer building a music‑focused platform application. The backend is Go (Gin) with Postgres, Redis, and AWS S3; the frontend is Next.js (App Router).

## Tech Stack at a glance
- Backend: Go (Gin), GORM, PostgreSQL, Redis, JWT auth, CORS, file upload limits, Dockerized
- Frontend: Next.js 15 (App Router) with React 19, Tailwind CSS 4
- Storage/Infra: AWS S3 (uploads), optional Nginx reverse proxy (prod), Docker Compose for local stack

## Project structure
- backend/
  - cmd/server/main.go – Gin bootstrap, routes, middleware, server
  - config/config.go – Loads env variables (DB, Redis, JWT, S3, server)
  - internal/ – controllers, services, middlewares, models
  - pkg/db/connection.go – Postgres connection via GORM + AutoMigrate
- site/ (Next.js app)
  - src/app/ – App Router pages and layouts
    - layout.js – global layout
    - page.js – marketing/landing page
    - (app)/ – authenticated area (e.g., dashboard, app)
  - package.json – Next.js 15 scripts (uses Turbopack)
- docker-compose.yml – Postgres, Redis, API, optional frontend, Nginx
- nginx/ – reverse proxy config (prod)
- Dockerfile – API build file (referenced by docker-compose)
- .junie/guidelines.md – this document

## Read this first — Documentation checklist
Before writing any code, read these docs end-to-end to understand intent, scope, and constraints:
- README.md (root) — architecture, quick start, features, security
- docs/OnePager.md — executive overview and value props
- docs/Stakeholder_Project_Outline.md — goals, scope, risks, standards, change mgmt
- docker-compose.yml — local stack, env defaults, healthchecks
- backend/config/config.go — environment variables and defaults
- backend/cmd/server/main.go — routes, middleware, limits, API base path
- backend/pkg/db/connection.go — DB connection and migrations
- site/package.json — scripts and framework versions
- site/src/app/* — layouts/pages to understand routing boundaries
If anything conflicts or is unclear, pause coding and ask for clarification before making changes.

## Environment variables
The backend loads from .env (if present) or process env.

Minimum required for local dev without Docker:
- DB_HOST=localhost
- DB_PORT=5432
- DB_USER=uploadparty
- DB_PASSWORD=your_local_db_password
- DB_NAME=uploadparty_db
- DB_SSL_MODE=disable
- PORT=8080
- GIN_MODE=debug
- FRONTEND_URL=http://localhost:3000
- JWT_SECRET=change_me
- REDIS_URL=redis://localhost:6379
- AWS_REGION=us-east-1
- AWS_ACCESS_KEY_ID=…
- AWS_SECRET_ACCESS_KEY=…
- AWS_S3_BUCKET=uploadparty-beats
- NI_API_KEY=… (if used)
- NI_API_URL=https://api.native-instruments.com

When using Docker Compose, many defaults are set in docker-compose.yml. You can still override via your shell env.

## Running locally

Option A: Full stack via Docker Compose (recommended for parity)
1) Ensure Docker is running.
2) From repo root:
   - docker compose up -d postgres redis
   - Wait for healthchecks to pass.
   - docker compose up -d api
   - (Optional) docker compose up -d nginx
3) API will listen on http://localhost:8080; Postgres on 5432; Redis on 6379.
4) Frontend is developed outside compose (see Option B). There is a compose “frontend” service placeholder but the active app lives in site/.

Option B: API via Go, Frontend via Next.js (manual dev mode)
- Backend
  - Create .env in backend/ or repo root with the vars above.
  - From backend/:
    - go run ./cmd/server
  - Health endpoint: http://localhost:8080/health
- Frontend
  - From site/:
    - npm install
    - npm run dev (Next.js 15 with Turbopack)
  - App served at http://localhost:3000

CORS: The API reads FRONTEND_URL for allowed origins. Ensure it matches your frontend URL during development.

## Build and production
- Backend (Docker): docker compose build api; docker compose up -d api
- Backend (binary): from backend/, run `go build -o bin/server ./cmd/server`
- Frontend: from site/, run `npm run build` then `npm start`
- Nginx (optional): docker compose up -d nginx to serve combined frontend/API behind reverse proxy

## Testing
- Go: place tests alongside code (xxx_test.go). Run from backend/ with `go test ./...`
- Frontend: No test runner is configured. If you add tests, prefer Vitest or Jest and document scripts in site/package.json.
- For PRs that modify backend logic, run `go test ./...` at minimum.

## Code style and conventions
- Go
  - Use Go 1.21+ (match local toolchain to CI if defined). Run `go fmt ./...` and `go vet ./...`.
  - Keep handlers lean; push logic to services. Use context-aware DB ops where relevant.
  - Use structured logging where feasible; current setup prints formatted logs via Gin.
  - GORM migrations: rely on db.AutoMigrate in pkg/db/connection.go. For destructive changes, add explicit migrations and back them up.
- API
  - Base path: /api/v1
  - Auth: JWT (see controllers.NewAuthController). Public endpoints under /auth; protected require Authorization: Bearer <token>.
  - Middlewares in place: security headers, CORS, request size limits, timeouts, IP rate limits, and an audio upload limiter pattern (audio mpeg/wav/mp3/x-wav). Respect these when adding routes.
  - Health: GET /health
- Frontend
  - Next.js App Router. Keep server/client component boundaries clear. Only use "use client" where required.
  - Tailwind CSS 4; avoid custom CSS unless needed.
  - Scripts (site/package.json):
    - dev: next dev --turbopack
    - build: next build --turbopack
    - start: next start
    - lint: eslint

## Media and uploads
- Large file uploads are expected (music). Middleware already limits request sizes and content types.
- For persistent storage, prefer streaming directly to S3 from the API, or use a temp dir (mounted at ./uploads) then upload to S3 and delete local temp files.
- Keep max body size and allowed MIME types in sync between frontend and backend.

## Security
- Never commit real secrets. Use .env locally; use environment injection in CI/CD.
- Keep JWT_SECRET strong in non-dev.
- Validate and sanitize all user inputs.
- Ensure CORS FRONTEND_URL matches deployed origin.
- Nginx config provided for TLS termination; place certs in nginx/ssl for prod.

## Branching and PR workflow
- Create feature branches from main.
- Keep changes minimal and scoped.
- Update this guidelines file if you introduce new required steps or tools.
- Before PR:
  - Backend: go fmt, go vet, go test ./...
  - Frontend: npm run build, npm run lint
  - Run the app locally (Option A or B) and smoke-test flows you changed.

## Common troubleshooting
- API can’t connect to DB: confirm DB_HOST/PORT/USER/PASSWORD/DB_NAME and that Postgres is healthy (docker ps, docker logs uploadparty-db).
- CORS errors: FRONTEND_URL mismatch. Update env or config.
- Uploads failing: check request size limits and MIME types in middlewares; verify S3 credentials and bucket name.
- 401 on protected routes: ensure Authorization header has Bearer token from /auth/login response.

## What Junie should do before submitting any change
0) Before coding:
   - Read all docs listed in 'Read this first — Documentation checklist' and ensure you understand the project context. If unclear, ask questions before making changes.
1) If backend code is touched:
   - go fmt ./...; go vet ./...; go test ./...
   - If DB or env settings changed, ensure docker compose still brings up services and the healthcheck passes.
2) If frontend code is touched:
   - npm run lint; npm run build
3) If both:
   - Run Option A or B locally and verify the changed flow in the browser.
4) Keep edits minimal and document anything noteworthy in PR description.


## Performance and capacity targets
- Concurrency baseline: The app must support at least 100 concurrent users at a time during normal operation without degradation of core flows (auth, uploads, feed, challenges).
- Scalability: Treat this as a platform application. Design for horizontal scalability of the API (stateless services, shared nothing), durable object storage via S3 (multipart uploads, lifecycle policies), and the option to front media with a CDN.
- Storage baseline: Expect to hold audio files for approximately 100 daily active users (DAU) initially; plan storage, egress, and cost controls accordingly.
- When adding or modifying features, ensure middleware limits (timeouts, request size, rate limits) are respected and do not unnecessarily throttle legitimate traffic under this concurrency.
- Prefer non-blocking I/O and streaming for uploads; avoid long CPU-bound work in request handlers. Offload heavy tasks to background workers where feasible.
- Monitor and optimize DB usage (connection pooling via GORM/Postgres) and cache hot paths with Redis when appropriate.
- For local/dev load checks, use lightweight tools (e.g., hey, autocannon) to simulate ~100 concurrent requests to critical endpoints (health, auth/login, representative GET/POST flows).
