# Development Guide

## Prerequisites

- **Go**: 1.23 or later
- **Node.js**: 18 or later
- **Docker**: For database and Redis
- **Make**: For build automation

## Quick Start

1. **Clone and setup environment**:
   ```bash
   git clone <repository-url>
   cd uploadparty
   cp .env.example .env
   # Edit .env with your settings
   ```

2. **Start infrastructure**:
   ```bash
   make db-up  # Starts PostgreSQL and Redis
   ```

3. **Install dependencies**:
   ```bash
   # Backend
   cd backend && go mod tidy

   # Frontend  
   cd site && npm install
   ```

4. **Run migrations**:
   ```bash
   make migrate
   ```

5. **Start development servers**:
   ```bash
   # Terminal 1: Backend with hot reload
   make api-air

   # Terminal 2: Frontend
   cd site && npm run dev
   ```

6. **Access the application**:
   - Frontend: http://localhost:3000
   - Backend API: http://localhost:8080
   - Health check: http://localhost:8080/health

## Available Make Commands

### Database
- `make db-up`: Start PostgreSQL and Redis containers
- `make db-down`: Stop database containers
- `make migrate`: Run database migrations
- `make migrate-dev`: Run migrations including dev seed data

### Backend
- `make api`: Run backend server (no hot reload)
- `make api-air`: Run backend with hot reload using Air
- `make air-install`: Install Air hot reload tool
- `make test-backend`: Run backend tests

### Full Stack
- `make dev`: Start all development services
- `make clean`: Clean build artifacts

## Development Workflow

### Backend Development

1. **File Structure**:
   - Controllers: Handle HTTP requests
   - Services: Business logic
   - Models: Database entities
   - Middlewares: Cross-cutting concerns (auth, logging)

2. **Adding New Endpoints**:
   ```go
   // 1. Add route in cmd/server/main.go
   api.GET("/projects", projCtl.ListMine)

   // 2. Implement controller method
   func (pc *ProjectController) ListMine(c *gin.Context) {
       // Implementation
   }

   // 3. Add tests
   func TestProjectController_ListMine(t *testing.T) {
       // Test implementation
   }
   ```

3. **Database Changes**:
   - Create migration file in `backend/migrations/`
   - Use sequential numbering: `002_add_projects_table.sql`
   - Test both up and down migrations

### Frontend Development

1. **Route Structure**:
   - Public routes: `src/app/(public)/`
   - Authenticated routes: `src/app/(app)/`
   - API routes: `src/app/api/`

2. **Component Guidelines**:
   - Use Server Components by default
   - Add `'use client'` only when needed for interactivity
   - Colocate tests in `__tests__/` directories

3. **API Integration**:
   ```javascript
   // Use the secure API client
   import secureApiClient from '../lib/secure-api';

   // Server-side data fetching
   const data = await secureApiClient.getProjects(req, res);

   // Client-side with hooks
   const { projects, loading, error } = useProjects();
   ```

## Testing

### Backend Tests
```bash
# Run all tests
cd backend && go test ./...

# Run with coverage
go test ./... -cover

# Run specific test
go test ./tests -run TestHealthController

# Benchmarks
go test ./... -bench=.
```

### Frontend Tests
```bash
cd site

# Unit tests
npm test

# Watch mode
npm run test:watch

# Coverage
npm run test:coverage

# E2E tests
npm run test:e2e

# Type checking
npm run typecheck
```

## Code Quality

### Backend Standards
- Use `golangci-lint` for linting
- Follow Go naming conventions
- Document exported functions
- Use table-driven tests
- Handle errors explicitly

### Frontend Standards
- Use ESLint with Next.js config
- Follow React/Next.js best practices
- Use TypeScript for type safety
- Test user interactions, not implementation details

## Environment Configuration

### Development (.env)
```bash
# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=uploadparty
DB_PASSWORD=your_password
DB_NAME=uploadparty_db

# API
PORT=8080
JWT_SECRET=your_jwt_secret
FRONTEND_URL=http://localhost:3000

# Frontend
NEXT_PUBLIC_API_URL=http://localhost:8080
NEXT_PUBLIC_SITE_URL=http://localhost:3000

# Auth0
AUTH0_BASE_URL=http://localhost:3000
AUTH0_ISSUER_BASE_URL=https://your-tenant.auth0.com
AUTH0_CLIENT_ID=your_client_id
AUTH0_CLIENT_SECRET=your_client_secret
```

### Production
- Use platform environment variables
- Never commit secrets to version control
- Use strong, randomly generated secrets
- Configure CORS for production domains

## Debugging

### Backend Debugging
- Use built-in logging with structured output
- Debug with Delve: `dlv debug ./cmd/server`
- Check health endpoint: `/health`

### Frontend Debugging
- Use React Developer Tools
- Check browser network tab for API calls
- Use Next.js built-in error overlay

## Performance

### Backend
- Monitor response times with logging middleware
- Use database indexes appropriately
- Implement caching where beneficial
- Profile with `go tool pprof`

### Frontend
- Use Next.js Image component for images
- Implement proper loading states
- Monitor Core Web Vitals
- Use React Profiler for component performance

## Security

### Development Security
- Use HTTPS in production
- Validate all inputs
- Sanitize database queries
- Keep dependencies updated
- Review security headers

### Authentication
- Use Auth0 for production authentication
- Set `BYPASS_AUTH=true` only in development
- Implement proper session management
- Use secure cookie settings

## Troubleshooting

### Common Issues

1. **Database connection fails**:
   - Check if PostgreSQL container is running
   - Verify connection string in .env
   - Ensure database exists

2. **Frontend can't reach backend**:
   - Check CORS configuration
   - Verify NEXT_PUBLIC_API_URL
   - Ensure backend is running on correct port

3. **Hot reload not working**:
   - Restart Air process
   - Check .air.toml configuration
   - Verify file watching permissions

4. **Tests failing**:
   - Clean test database
   - Check test environment configuration
   - Verify all dependencies are installed

### Getting Help

1. Check error logs in terminal
2. Review health check endpoint
3. Verify environment configuration
4. Check Docker container status
5. Review recent changes in git