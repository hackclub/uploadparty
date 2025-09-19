# Multi-stage build for backend (Go) and frontend (Next.js)

# --- Backend builder ---
FROM golang:1.21-alpine AS gobuilder

# Install required packages for Go build
RUN apk update && apk add --no-cache \
    ca-certificates \
    git \
    tzdata \
    && rm -rf /var/cache/apk/*

# Non-root user (mirrored later)
RUN adduser -D -g '' appuser

WORKDIR /app

# Copy go mod/sum and download deps (expects repo-root build context)
COPY backend/go.mod backend/go.sum ./
RUN go mod download && go mod verify

# Copy backend source code and build
COPY backend/ .
# Ensure module graph and go.sum are up to date after copying sources
RUN go mod tidy && go mod download && go mod verify
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags='-w -s -extldflags "-static"' \
    -a -installsuffix cgo \
    -o main cmd/server/main.go

# --- Frontend dependencies (for build) ---
FROM node:20-alpine AS webdeps
WORKDIR /web
COPY site/package*.json ./
# Full deps for build to avoid missing dev-time tools (like Tailwind/PostCSS)
RUN npm ci

# --- Frontend builder ---
FROM node:20-alpine AS webbuilder
WORKDIR /web
COPY --from=webdeps /web/node_modules ./node_modules
COPY site/ .
RUN npm run build

# --- Frontend production deps (runtime only) ---
FROM node:20-alpine AS webproddeps
WORKDIR /web
COPY site/package*.json ./
RUN npm ci --omit=dev

# --- Final runtime image: Node (for Next.js) + Go binary ---
FROM node:20-alpine

# Install CA certs and tzdata for consistency
RUN apk update && apk add --no-cache \
    ca-certificates \
    tzdata \
    wget \
    && rm -rf /var/cache/apk/* \
    && update-ca-certificates

# Create non-root user
RUN adduser -D -g '' appuser

# Working directory
WORKDIR /app

# Backend binary
COPY --from=gobuilder /app/main /app/main

# Frontend app (copy only what's needed to run `next start`)
RUN mkdir -p /app/site
COPY --from=webbuilder /web/.next /app/site/.next
COPY --from=webbuilder /web/public /app/site/public
COPY --from=webbuilder /web/package.json /app/site/package.json
COPY --from=webproddeps /web/node_modules /app/site/node_modules

# Ownership
RUN chown -R appuser:appuser /app

# Switch to non-root
USER appuser

# Expose backend and frontend ports
EXPOSE 8080 3000

# Healthcheck for backend (frontend served on 3000)
HEALTHCHECK --interval=30s --timeout=5s --start-period=10s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

# Start both processes: backend (8080) and Next.js (3000)
CMD ["sh", "-c", "PORT=8080 ./main & (cd site && node node_modules/next/dist/bin/next start -p 3000)"]
