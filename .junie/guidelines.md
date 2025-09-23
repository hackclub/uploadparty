# UploadParty – Project Guidelines for Junie (Full‑Stack: Go Gin + Next.js)

These guidelines tell Junie how to work within this repository as a full‑stack engineer building a music‑focused platform application. The backend is Go (Gin) with Postgres, Redis, and Google Cloud Storage; the frontend is Next.js (App Router). Important: Junie must re-read this guidelines document at the start of every request to ensure all actions follow the latest project rules.

## Tech Stack at a glance
- Backend: Go (Gin), GORM, PostgreSQL, Redis, JWT auth, CORS, file upload limits, Dockerized
- Frontend: Next.js 15 (App Router) with React 19; styling framework is optional. If the user asks, Junie should write unique, handcrafted CSS on demand.
- Storage/Infra: Google Cloud Platform (app hosting and persistent volumes with lots of storage), Google Cloud Storage (primary for object storage/uploads), optional Nginx reverse proxy (prod), Docker Compose for local stack

## Project structure
- backend/
  - cmd/server/main.go – Gin bootstrap, routes, middleware, server
  - config/config.go – Loads env variables (DB, Redis, JWT, GCS, server)
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

Note: The exact design of the authenticated area (after-auth experience) is not finalized yet. When making assumptions or scaffolding routes/components inside site/src/app/(app)/, consult docs/OnePager.md to ground decisions in the current product concept.

## Full‑Stack application requirements (implementation checklist)
- Backend (Gin/Golang)
  - Use Gin with high‑performance routing (radix tree via httprouter) and a small memory footprint.
  - Add middleware chain examples: Logger, Recovery, CORS, Auth (JWT), rate limiters, request size/timeouts. Groups: /api/v1, /auth, and protected routes.
  - Demonstrate JSON validation and binding (query, form, multipart). Return validation errors with clear messages.
  - Support multiple renderers where appropriate (JSON primary; optionally XML/YAML/ProtoBuf for examples).
  - Ensure crash‑free behavior with gin.Recovery() and panic safety in handlers.
- Frontend (Next.js/React)
  - Build with Next.js App Router (React Server Components by default). Use client components only where necessary.
  - Implement pages and UI as React components. Styling is framework-agnostic; if the user asks, write unique, handcrafted CSS. Tailwind is optional.
- Integration (Full‑Stack Connection)
  - The Next.js app should call the Gin API (e.g., GET /health, POST /api/v1/example) to demonstrate SSR/ISR/fetch patterns.
  - Respect CORS by aligning FRONTEND_URL with the running frontend origin.

Note on storage and hosting: We deploy on Google Cloud Platform and use Google Cloud Storage for primary media/upload storage. Google Cloud Storage provides scalable object storage and CDN workflows with cost-effective options.

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
- GOOGLE_CLOUD_PROJECT=uploadparty-project
- GOOGLE_CLOUD_STORAGE_BUCKET=uploadparty-beats
- GOOGLE_APPLICATION_CREDENTIALS=path/to/service-account-key.json
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
- Google Cloud Platform (recommended hosting): Deploy API and frontend as Cloud Run services. Use Google Cloud Storage for primary file storage. Configure env vars via Cloud Run (match those in backend/config/config.go).

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
  - Styling: Framework-agnostic. If the user asks, write unique, handcrafted CSS; Tailwind is optional and not required.
  - Scripts (site/package.json):
    - dev: next dev --turbopack
    - build: next build --turbopack
    - start: next start
    - lint: eslint

## Media and uploads
- Large file uploads are expected (music). Middleware already limits request sizes and content types.
- Storage strategy: We use Google Cloud Storage as the primary durable storage for media in production. For local development, use the ./uploads directory to emulate production storage behavior. When enabled, stream directly to Google Cloud Storage or stage to ./uploads and then upload, followed by cleanup.
- Local/dev: Use the ./uploads directory (created automatically or mounted via Docker) to emulate production storage behavior.
- Keep max body size and allowed MIME types in sync between frontend and backend.
- Prefer streaming and non‑blocking I/O; avoid holding entire files in memory; ensure temporary files are removed after successful uploads.

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
- Uploads failing: check request size limits and MIME types in middlewares; verify Google Cloud Storage credentials and bucket name.
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
- Scalability: Treat this as a platform application. Design for horizontal scalability of the API (stateless services, shared nothing), durable object storage via Google Cloud Storage (multipart uploads, lifecycle policies), and the option to front media with a CDN.
- Storage baseline: Expect to hold audio files for approximately 100 daily active users (DAU) initially; plan storage, egress, and cost controls accordingly.
- When adding or modifying features, ensure middleware limits (timeouts, request size, rate limits) are respected and do not unnecessarily throttle legitimate traffic under this concurrency.
- Prefer non-blocking I/O and streaming for uploads; avoid long CPU-bound work in request handlers. Offload heavy tasks to background workers where feasible.
- Monitor and optimize DB usage (connection pooling via GORM/Postgres) and cache hot paths with Redis when appropriate.
- For local/dev load checks, use lightweight tools (e.g., hey, autocannon) to simulate ~100 concurrent requests to critical endpoints (health, auth/login, representative GET/POST flows).


## Frontend routing architecture and auth boundary
To keep the frontend clean and predictable, we enforce a clear separation between the public (pre‑auth) experience and the authenticated application area.

- Public (pre‑auth) landing:
  - Path: site/src/app/page.js
  - Purpose: Marketing/landing page and any public content visible before authentication.
  - Notes: Prefer React Server Components where possible. Avoid unnecessary "use client". Keep bundles small and cacheable. No auth required.

- Authenticated application:
  - Path: site/src/app/(app)/
  - Purpose: The signed‑in experience (e.g., dashboard, settings, uploads). This group has its own layout file at site/src/app/(app)/layout.js.
  - Routing: All routes inside (app) are considered protected and should assume an authenticated user context.
  - Data fetching: Use server components by default; elevate to client components only for interactive UI pieces.

- Directory and routing conventions:
  - Root route / is public (site/src/app/page.js).
  - Grouped routes under (app) reflect the authenticated area with its own layout and navigation.
  - You may optionally introduce a (public) group for additional public sections if needed, but the root landing remains at page.js.

- Auth boundary (recommended pattern):
  - Implement auth checks at the boundary of (app). Common options:
    - Next.js Middleware (site/src/middleware.ts) to redirect unauthenticated users away from /(app) to / (or /auth/login).
    - Layout‑level guard within site/src/app/(app)/layout.js that checks session/JWT on the server and redirects if missing.
  - Ensure redirects are fast and avoid client‑side flashes. Prefer server‑side checks.

- API integration from both sides:
  - Public pages can call read‑only endpoints (e.g., GET /health) and should respect CORS.
  - Authenticated routes call protected API endpoints with credentials (cookies or Authorization: Bearer <token>), following the backend’s /api/v1 and /auth grouping.

- Styling conventions:
  - Framework‑agnostic; when the user asks, write unique, handcrafted CSS. Tailwind is optional, not required.

- Example mental model:
  - Public: “What is UploadParty?” → site/src/app/page.js
  - Authenticated: “Do the work” → site/src/app/(app)/* (dashboard, settings, uploads)

- Status of post‑auth UX:
  - The exact after‑auth experience is not finalized yet. Use docs/OnePager.md to ground assumptions about user goals, flows, and scope until detailed specs are provided.

## FAQ: Gin performance and reflection
Q: Does the Gin Web Framework achieve its improved performance by avoiding the use of reflection?
A: Gin’s high performance primarily comes from its use of a radix tree router (via httprouter) and keeping the hot path allocation‑free. Routing and handler dispatch do not rely on reflection. However, Gin’s request binding/validation helpers (e.g., ShouldBindJSON, form/multipart binding) use Go reflection to map payloads to structs. In short: Gin minimizes reflection on the hot path (routing), but it does use reflection for binding/validation when you opt into those features.