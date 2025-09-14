# Multi-stage build for security and smaller image size
FROM golang:1.21-alpine AS builder

# Install security updates and required packages
RUN apk update && apk add --no-cache \
    ca-certificates \
    git \
    tzdata \
    && rm -rf /var/cache/apk/*

# Create appuser for security
RUN adduser -D -g '' appuser

# Set working directory
WORKDIR /app

# Copy go mod and sum files for dependency caching
COPY backend/go.mod backend/go.sum ./

# Download dependencies
RUN go mod download && go mod verify

# Copy backend source code
COPY backend/ .

# Build the application with security flags
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags='-w -s -extldflags "-static"' \
    -a -installsuffix cgo \
    -o main cmd/server/main.go

# Final stage - minimal runtime image
FROM alpine:3.18

# Install security updates and CA certificates
RUN apk update && apk add --no-cache \
    ca-certificates \
    tzdata \
    && rm -rf /var/cache/apk/* \
    && update-ca-certificates

# Create appuser in final image
RUN adduser -D -g '' appuser

# Set working directory
WORKDIR /app

# Copy the binary from builder stage
COPY --from=builder /app/main .

# Copy any required config files (if needed)
# COPY --from=builder /app/config ./config

# Change ownership to appuser
RUN chown -R appuser:appuser /app

# Switch to non-root user
USER appuser

# Expose port (use non-privileged port)
EXPOSE 8080

# Add health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

# Run the application
CMD ["./main"]
