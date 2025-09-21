# UploadParty A modern, music‑focused platform for uploading beats and audio, running community challenges.

## Website
- Live site: coming soon
- Note: Set the repository Description, Website, and Topics in GitHub → Settings → General and Settings → Topics for better discoverability.

## Tech stack
- Backend: Go (Gin), GORM, PostgreSQL, Redis, JWT auth, CORS, upload limits
- Frontend: Next.js 15 (App Router) + React 19, Tailwind CSS 4
- Storage/Infra: AWS S3, optional Nginx reverse proxy, Docker Compose for local
# UploadParty A modern, music‑focused platform for uploading beats and audio, running community challenges.

## Website
- Live site: coming soon
- Note: Set the repository Description, Website, and Topics in GitHub → Settings → General and Settings → Topics for better discoverability.

## Tech stack
- Backend: Go (Gin), GORM, PostgreSQL, Redis, JWT auth, CORS, upload limits
- Frontend: Next.js 15 (App Router) + React 19, Tailwind CSS 4
- Storage/Infra: AWS S3, optional Nginx reverse proxy, Docker Compose for local

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

