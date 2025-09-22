# Security Guidelines for Upload Party

## Overview

Upload Party handles sensitive user data including personal information, audio files, and authentication credentials. This document outlines the security measures implemented to protect users and the platform.

## Security Features Implemented

### 1. Authentication & Authorization
- JWT-based authentication with secure token generation
- Password hashing using bcrypt with appropriate salt rounds
- Role-based access control for different user types
- Session management with secure token storage

### 2. Input Validation & Sanitization
- Request size limits (50MB for general requests, 100MB for audio uploads)
- File type validation for audio uploads
- SQL injection prevention through GORM ORM
- XSS protection via security headers
- Input validation using struct tags and validators

### 3. Rate Limiting
- Global rate limiting: 10 requests per second with burst of 20
- IP-based rate limiting to prevent abuse
- Stricter limits for file uploads (1 request per second)
- DDoS protection at Nginx level

### 4. Security Headers
- `X-Frame-Options: DENY` - Prevents clickjacking
- `X-Content-Type-Options: nosniff` - Prevents MIME sniffing
- `X-XSS-Protection: 1; mode=block` - XSS protection
- `Strict-Transport-Security` - Forces HTTPS
- Content Security Policy (CSP) - Prevents code injection
- `Referrer-Policy` - Controls referrer information

### 5. HTTPS/TLS
- TLS 1.2+ enforcement
- Strong cipher suites
- HSTS headers for browser security
- SSL certificate validation

### 6. Database Security
- Connection pooling with secure configurations
- Row-level security policies
- Database user with minimal required privileges
- SQL injection prevention
- Prepared statements for all queries

### 7. File Upload Security
- File type validation (audio files only)
- File size limits
- Virus scanning (to be implemented with Cloud Storage or a scanning service)
- Secure file naming to prevent path traversal
- Direct upload to Google Cloud Storage (GCS) with signed URLs

### 8. Container Security
- Non-root user in Docker containers
- Read-only filesystem where possible
- Minimal base images (Alpine Linux)
- Security updates applied
- No secrets in container images

### 9. Environment Security
- Environment variables for all sensitive data
- No hardcoded secrets in code (never commit credentials or API keys)
- Use .env files for local/dev only; production secrets must be injected via the platform (Coolify, Docker secrets, etc.)
- The backend automatically loads .env files if present (backend/.env and/or ./.env)
- Separate production configurations
- Secrets management best practices

## Environment Variables Security

### Required for Production

```bash
# Strong JWT secret (minimum 256 bits)
JWT_SECRET=your-super-secure-jwt-secret-key-min-256-bits

# Database credentials
DB_PASSWORD=your-super-secure-production-db-password
REDIS_PASSWORD=your-secure-redis-password

# Google Cloud (GCS) configuration
GCP_PROJECT_ID=your-gcp-project-id
GCS_BUCKET=your-prod-gcs-bucket
# Path inside container or host to service account JSON
GOOGLE_APPLICATION_CREDENTIALS=/path/to/service-account.json

# External license directory (generic)
# Avoid provider-specific names in env to reduce information exposure
LICENSES_PROVIDER=none
LICENSES_TOKEN=
LICENSES_DSN=
```

## Docker Security Configuration

### Multi-stage builds
- Separate build and runtime environments
- Minimal attack surface in production image

### User permissions
- Non-root user (`appuser`) for running application
- Proper file ownership and permissions

### Read-only filesystem
- Application runs with read-only root filesystem
- Temporary directories mounted as tmpfs

## Network Security

### Nginx Configuration
- SSL/TLS termination at proxy level
- Rate limiting and DDoS protection
- Security headers injection
- Request size limits
- Timeout configurations

### Internal Communications
- Services communicate over internal Docker network
- Database not exposed externally
- Redis secured with password authentication

## Monitoring & Logging

### Security Logging
- Failed authentication attempts
- Rate limit violations
- File upload attempts
- Error conditions and exceptions

### Monitoring
- Health checks for all services
- Performance metrics
- Security event alerting

## Data Protection

### Personal Data
- GDPR compliance considerations
- Data retention policies
- User data encryption at rest
- Secure data transmission

### Audio Files
- S3 bucket security policies
- Access logging
- File integrity checks
- Backup and recovery procedures

## Incident Response

### Security Incident Handling
1. Immediate containment
2. Impact assessment
3. Evidence collection
4. Recovery procedures
5. Post-incident review

### Contact Information
- Security team: security@uploadparty.com
- Emergency contact: +1-XXX-XXX-XXXX

## Security Testing

### Regular Testing
- Dependency vulnerability scanning
- Container image scanning
- Penetration testing
- Code security reviews

### Automated Security
- CI/CD pipeline security checks
- Automated dependency updates
- Security linting and analysis

## Compliance

### Standards
- OWASP Top 10 compliance
- GDPR data protection requirements
- Industry security best practices

### Audit Trail
- All security-relevant events logged
- Audit log integrity protection
- Regular security assessments

## Deployment Security Checklist

### Pre-deployment
- [ ] All secrets properly configured
- [ ] SSL certificates installed
- [ ] Database security configured
- [ ] Rate limiting tested
- [ ] Security headers verified

### Post-deployment
- [ ] Security monitoring active
- [ ] Backup procedures tested
- [ ] Incident response plan updated
- [ ] Security documentation current

## Updates and Maintenance

### Security Updates
- Regular dependency updates
- Operating system security patches
- Container base image updates
- Security configuration reviews

### Version Control
- No secrets committed to repository
- Signed commits for security-critical changes
- Regular security-focused code reviews

---

For security questions or to report vulnerabilities, please contact: security@uploadparty.com

**Last Updated:** 2025-01-01
**Version:** 1.0


## Secrets folder (local/dev)

- This public repository includes a secrets/ directory that is git-ignored by default to help you keep sensitive files out of version control.
- Place local-only secret files here, such as gcp-service-account.json and sample license files (licenses.token, licenses.dsn).
- Docker Compose mounts ./secrets into the API container at /app/secrets (read-only). The default GOOGLE_APPLICATION_CREDENTIALS path is /app/secrets/gcp-service-account.json.
- Never commit real secrets. Only *.example files and README.md inside secrets/ are tracked.
