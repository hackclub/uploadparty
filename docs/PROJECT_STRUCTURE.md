# UploadParty Project Structure

## Overview
UploadParty is a full-stack application for music producers to upload beats and participate in community challenges.

## Directory Structure

```
uploadparty/
├── .junie/                    # Junie AI assistant guidelines and documentation
│   └── guidelines.md          # Production standards and best practices
├── backend/                   # Go API server
│   ├── cmd/                   # Application entry points
│   │   ├── migrate/           # Database migration tool
│   │   └── server/            # Main API server
│   ├── config/                # Configuration management
│   ├── internal/              # Private application code
│   │   ├── controllers/       # HTTP handlers
│   │   ├── integrations/      # External service integrations
│   │   ├── middlewares/       # HTTP middlewares (auth, logging, validation)
│   │   ├── models/            # Database models
│   │   ├── services/          # Business logic
│   │   └── utils/             # Utility functions
│   ├── migrations/            # SQL migration files
│   ├── pkg/                   # Public packages
│   │   └── db/                # Database connection utilities
│   └── tests/                 # Test files
├── docs/                      # Project documentation
│   ├── frontend/              # Frontend-specific documentation
│   │   └── SEO_TESTING.md     # SEO testing checklist
│   ├── PROJECT_STRUCTURE.md   # This file
│   └── TODO_UNLEASH.md        # Feature flag implementation guide
├── nginx/                     # Nginx configuration
├── secrets/                   # Local development secrets (git-ignored)
├── site/                      # Next.js frontend application
│   ├── e2e/                   # End-to-end tests (Playwright)
│   ├── public/                # Static assets
│   └── src/                   # Source code
│       └── app/               # Next.js App Router
│           ├── (app)/         # Authenticated routes
│           ├── (public)/      # Public routes
│           ├── __tests__/     # Unit tests
│           ├── api/           # API routes (Auth0)
│           ├── components/    # Shared React components
│           ├── hooks/         # Custom React hooks
│           └── lib/           # Utility functions and API clients
└── Docker files and configs
```

## Key Components

### Backend (Go/Gin)
- **Framework**: Gin HTTP framework
- **Database**: PostgreSQL with GORM ORM
- **Authentication**: JWT tokens
- **Features**: Rate limiting, CORS, structured logging
- **Testing**: Unit tests with testify, SQLite for test database

### Frontend (Next.js)
- **Framework**: Next.js 15 with App Router
- **React**: React 19
- **Styling**: Tailwind CSS 4
- **Authentication**: Auth0 integration
- **Testing**: Jest + React Testing Library + Playwright

### Development Tools
- **Hot Reload**: Air for Go backend development
- **Linting**: ESLint for frontend, golangci-lint for backend
- **Testing**: Comprehensive test suite with coverage requirements
- **Type Safety**: TypeScript for frontend type checking

## Development Workflow

### Getting Started
1. Copy `.env.example` to `.env` at project root
2. Start services: `make db-up` (database and Redis)
3. Run backend: `make api-air` (with hot reload)
4. Run frontend: `cd site && npm run dev`

### Testing
- **Backend**: `cd backend && go test ./...`
- **Frontend Unit**: `cd site && npm test`
- **Frontend E2E**: `cd site && npm run test:e2e`

### Production Deployment
- **Docker**: Multi-service Docker Compose setup
- **Environment**: Centralized environment configuration
- **Monitoring**: Health checks and structured logging
- **Security**: Production-ready security headers and validation

## File Organization Principles

### Backend
- **cmd/**: Executable entry points
- **internal/**: Private application code (cannot be imported by other projects)
- **pkg/**: Public packages (can be imported by other projects)
- **tests/**: All test files

### Frontend
- **Route Groups**: `(app)` for authenticated, `(public)` for public routes
- **Colocation**: Tests alongside source files in `__tests__/`
- **Shared Code**: Common components, hooks, and utilities in dedicated directories

## Documentation Guidelines

- **README files**: Provide context for each major directory
- **Code comments**: Document exported functions and complex logic
- **API documentation**: Maintain up-to-date API route documentation
- **Deployment guides**: Step-by-step production deployment instructions

## Security Considerations

- **Secrets Management**: No secrets in code, use environment variables
- **Authentication**: Secure JWT handling and Auth0 integration
- **Input Validation**: Comprehensive validation on all endpoints
- **Security Headers**: Production-ready security middleware
- **Dependencies**: Regular security audits and updates