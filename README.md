# Upload Party 🎵

Upload Party is a gamified online music creation platform designed for 500 teenagers, focused on beat sharing, collaboration, and competition. Supported by Native Instruments (500 Komplete 15 licenses) and Hack Club.

## 🎯 Mission

Help 50,000 high school students build technical skills through music creation, gamification, and community collaboration.

## 🏗️ Architecture

### Backend (Go + Gin)
- **Language**: Golang with Gin framework
- **Database**: PostgreSQL with GORM
- **Cache**: Redis for sessions and leaderboards
- **Storage**: AWS S3 for audio files
- **Authentication**: JWT with secure middleware
- **Security**: Rate limiting, CORS, file validation, encryption

### Frontend (React + TypeScript)
- **Framework**: React 18 with TypeScript
- **Routing**: React Router v6
- **Styling**: Tailwind CSS with custom components
- **State Management**: Zustand
- **API Client**: Axios with interceptors
- **UI**: Headless UI, Framer Motion, Heroicons

### Infrastructure
- **Containerization**: Docker with multi-stage builds
- **Reverse Proxy**: Nginx with SSL termination
- **Orchestration**: Docker Compose
- **Security**: Non-root containers, read-only filesystems

## 📁 Project Structure

```
/uploadparty
├── backend/                    # Go backend application
│   ├── cmd/server/            # Application entry point
│   ├── config/                # Configuration management
│   ├── internal/              # Private application code
│   │   ├── controllers/       # HTTP handlers
│   │   ├── services/          # Business logic
│   │   ├── repositories/      # Data access layer
│   │   ├── models/           # Data structures
│   │   ├── middlewares/      # HTTP middleware
│   │   └── utils/            # Utility functions
│   └── pkg/                   # Public packages
├── frontend/                   # React frontend application
│   ├── src/
│   │   ├── components/       # Reusable UI components
│   │   ├── pages/           # Route components
│   │   ├── services/        # API clients
│   │   ├── hooks/           # Custom React hooks
│   │   ├── store/           # State management
│   │   ├── types/           # TypeScript definitions
│   │   └── utils/           # Utility functions
│   ├── public/              # Static assets
│   └── Dockerfile           # Frontend container build
├── nginx/                     # Reverse proxy configuration
├── scripts/                   # Database and deployment scripts
├── docker-compose.yml         # Multi-service orchestration
├── Dockerfile                 # Backend container build
└── SECURITY.md               # Security documentation
```

## 🚀 Quick Start

### Prerequisites
- Docker and Docker Compose
- Node.js 18+ (for local frontend development)
- Go 1.21+ (for local backend development)
- PostgreSQL 15+ (if running locally)

### Development Setup

1. **Clone and Setup Environment**
   ```bash
   git clone <repository-url>
   cd uploadparty
   cp .env.example .env
   # Edit .env with your configuration
   ```

2. **Start with Docker Compose**
   ```bash
   docker-compose up -d
   ```

3. **Local Frontend Development**
   ```bash
   cd frontend
   npm install
   npm start
   ```

4. **Local Backend Development**
   ```bash
   cd backend
   go mod tidy
   go run cmd/server/main.go
   ```

### Production Deployment

1. **Configure Production Environment**
   ```bash
   cp .env.example .env.production
   # Configure production values
   ```

2. **Deploy Stack**
   ```bash
   docker-compose -f docker-compose.yml -f docker-compose.prod.yml up -d
   ```

## 🔧 Key Features

### 🎵 Music Creation
- Audio file upload (MP3, WAV support)
- Beat metadata (BPM, key, genre, tags)
- Audio visualization and playback
- File size and format validation

### 🏆 Gamification
- Point-based scoring system
- User levels and rankings
- Achievement badges
- Real-time leaderboards

### 🤝 Community Features
- User profiles and portfolios
- Beat commenting and feedback
- Collaboration requests
- Social following system

### 🎯 Challenges
- Time-limited music challenges
- Genre-specific competitions
- Native Instruments prize distribution
- Submission judging and scoring

### 🔐 Security
- JWT authentication with secure headers
- Rate limiting and DDoS protection
- File upload validation
- Row-level database security
- HTTPS/TLS encryption

## 📊 Database Models

### Core Entities
- **Users**: Profiles, authentication, gamification data
- **Beats**: Audio files, metadata, engagement metrics
- **Challenges**: Competitions, rules, rewards
- **Comments**: Timestamped feedback on beats
- **Scores**: Point tracking and leaderboard data

### Relationships
- Users create beats and participate in challenges
- Beats can have multiple collaborators
- Challenges accept beat submissions
- Comments provide timestamped feedback

## 🌐 API Endpoints

### Authentication
- `POST /api/v1/auth/register` - User registration
- `POST /api/v1/auth/login` - User login
- `GET /api/v1/profile` - Get user profile

### Beats
- `GET /api/v1/beats` - List beats with filters
- `POST /api/v1/beats/upload` - Upload new beat
- `GET /api/v1/beats/:id` - Get beat details
- `POST /api/v1/beats/:id/like` - Like/unlike beat

### Challenges
- `GET /api/v1/challenges` - List active challenges
- `GET /api/v1/challenges/:id` - Challenge details
- `POST /api/v1/challenges/:id/submit` - Submit to challenge

### Community
- `GET /api/v1/leaderboard` - Get leaderboard
- `GET /api/v1/users/:id` - User profile
- `POST /api/v1/users/:id/follow` - Follow user

## 🛡️ Security Features

### Backend Security
- JWT token authentication
- bcrypt password hashing
- Rate limiting (10 req/s, burst 20)
- File type and size validation
- SQL injection prevention via ORM
- CORS protection
- Security headers (HSTS, CSP, etc.)

### Infrastructure Security
- Non-root container execution
- Read-only filesystems
- Multi-stage Docker builds
- Secrets management
- SSL/TLS termination
- Database connection encryption

## 📈 Performance Optimizations

### Backend
- Database connection pooling
- Redis caching for sessions
- Optimized database indexes
- File upload streaming
- Response compression

### Frontend
- Code splitting and lazy loading
- Asset optimization and caching
- Service worker for offline support
- Optimized bundle sizes
- CDN-ready static assets

## 🧪 Testing

### Backend Testing
```bash
cd backend
go test ./...
```

### Frontend Testing
```bash
cd frontend
npm test
```

### Integration Testing
```bash
docker-compose -f docker-compose.test.yml up --build
```

## 📚 Documentation

- [Security Guidelines](SECURITY.md)
- [API Documentation](docs/api.md)
- [Deployment Guide](docs/deployment.md)
- [Contributing Guidelines](CONTRIBUTING.md)

## 🤝 Contributing

1. Fork the repository
2. Create feature branch (`git checkout -b feature/amazing-feature`)
3. Commit changes (`git commit -m 'Add amazing feature'`)
4. Push to branch (`git push origin feature/amazing-feature`)
5. Open Pull Request

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- **Native Instruments** for providing 500 Komplete 15 licenses
- **Hack Club** for supporting high school technical education
- The open-source community for excellent tools and libraries

## 📞 Support

- Email: support@uploadparty.com
- Discord: [Upload Party Community](https://discord.gg/uploadparty)
- Documentation: [docs.uploadparty.com](https://docs.uploadparty.com)

---

**Built with ❤️ for the next generation of music producers**
