# SuProxy Backend

Modern, production-ready VPN proxy management backend built with Go and Clean Architecture.

## Features

- 🏗️ Clean Architecture (Domain, Application, Infrastructure, Interface)
- 🔐 JWT Authentication & Authorization
- 👥 User Management (Admin & User roles)
- 🚀 Xray Instance Management
- 📊 Real-time Metrics & Monitoring (Prometheus + Grafana)
- 🔍 Audit Logging
- 🐳 Docker & Docker Compose support
- 📝 Structured Logging
- ✅ Health Checks
- 🔄 Graceful Shutdown
- 🛡️ Production-ready security hardening

## Quick Start

### Development

```bash
# Copy environment file
cp .env.example .env

# Start services
docker-compose up -d

# View logs
docker-compose logs -f api
```

Access the API at http://localhost:8080

### Production

```bash
# Setup production environment
cp .env.example .env.production
# Edit .env.production with secure values

# Deploy (Linux/Mac)
./scripts/deploy.sh

# Deploy (Windows)
.\scripts\deploy.ps1
```

See [DEPLOYMENT.md](DEPLOYMENT.md) for detailed production deployment instructions.

## API Endpoints

### Public Endpoints
- `POST /api/v1/auth/register` - User registration
- `POST /api/v1/auth/login` - User login
- `GET /health` - Health check
- `GET /metrics` - Prometheus metrics

### User Endpoints (Authenticated)
- `GET /api/v1/auth/me` - Current user info
- `GET /api/v1/auth/sessions` - User sessions
- `POST /api/v1/auth/refresh` - Refresh token
- `POST /api/v1/auth/logout` - Logout
- `GET /api/v1/subscriptions` - User subscriptions
- `GET /api/v1/devices` - User devices
- `GET /api/v1/servers` - Available servers

### Admin Endpoints (Admin Role)
- User Management: 4 endpoints
- Xray Instance Management: 8 endpoints
- Inbound Management: 7 endpoints
- Client Management: 8 endpoints
- Audit Logs: 3 endpoints
- System Admin: 5 endpoints

Total: 47 production-ready endpoints

## Project Structure

```
.
├── cmd/
│   └── api/              # Application entry point
├── internal/
│   ├── domain/           # Domain entities & business logic
│   ├── application/      # Use cases, DTOs, mappers
│   │   ├── dto/         # Data Transfer Objects
│   │   ├── mapper/      # Entity-DTO mappers
│   │   ├── service/     # Application services
│   │   └── usecase/     # Use case implementations
│   ├── infrastructure/   # External concerns
│   │   ├── config/      # Configuration
│   │   ├── database/    # Database & migrations
│   │   ├── jwt/         # JWT implementation
│   │   ├── logger/      # Structured logging
│   │   ├── metrics/     # Prometheus metrics
│   │   └── xray/        # Xray integration
│   └── interfaces/       # External interfaces
│       └── http/         # HTTP handlers & middleware
├── configs/              # Configuration files
├── scripts/              # Deployment & maintenance scripts
├── prometheus/           # Prometheus configuration
├── grafana/              # Grafana provisioning
└── docs/                 # Documentation
```

## Technology Stack

- **Language**: Go 1.21+
- **Framework**: Gin (HTTP router)
- **Database**: PostgreSQL 15
- **Metrics**: Prometheus
- **Visualization**: Grafana
- **Migration**: golang-migrate
- **Container**: Docker & Docker Compose
- **Logging**: Structured JSON logging

## Architecture Principles

### Clean Architecture Layers

1. **Domain Layer**: Core business entities and interfaces
2. **Application Layer**: Use cases, business logic orchestration
3. **Infrastructure Layer**: External dependencies (DB, Xray, etc.)
4. **Interface Layer**: HTTP handlers, middleware

### Design Patterns

- Repository pattern for data access
- Service pattern for business logic
- DTO pattern for data transfer
- Dependency injection via constructor
- Interface-based abstractions

## Monitoring & Observability

### Metrics (Prometheus)
- HTTP request duration & count
- Active users, clients, xray instances
- Database connection pool metrics
- Go runtime metrics
- Custom business metrics

### Logging
- Structured JSON logs
- Request ID tracking
- Correlation ID support
- Log levels: debug, info, warn, error

### Health Checks
- Database connectivity
- Xray process status
- Application health
- System statistics

## Security Features

- 🔐 JWT-based authentication
- 🛡️ Role-based access control (RBAC)
- 🔒 Non-root container execution
- 📦 Read-only container filesystem
- 🚫 No new privileges security option
- 🔑 Environment-based secrets
- 📝 Comprehensive audit logging
- 🌐 SSL/TLS support (via reverse proxy)

## Development

### Prerequisites
- Go 1.21+
- Docker & Docker Compose
- PostgreSQL 15 (optional, Docker provides)
- Make (optional)

### Running Locally

```bash
# Install dependencies
go mod download

# Run database (Docker)
docker-compose up -d postgres

# Run migrations
make migrate-up

# Run application
make run
```

### Testing

```bash
# Run tests
make test

# Run tests with coverage
make test-coverage

# Run linter
make lint
```

### Building

```bash
# Build binary
make build

# Build Docker image
make docker-build
```

## Deployment

### Environment Variables

See `.env.example` for all available configuration options.

**Critical Production Settings:**
- `SUPROXY_ENVIRONMENT=production`
- `SUPROXY_JWT_SECRET_KEY` (64+ chars)
- `SUPROXY_DATABASE_PASSWORD` (strong password)
- `SUPROXY_DATABASE_SSLMODE=require`
- `SUPROXY_XRAY_USE_MOCK=false`

### Docker Compose

**Development:**
```bash
docker-compose up -d
```

**Production:**
```bash
docker-compose -f docker-compose.production.yml up -d
```

### Resource Requirements

**Minimum:**
- CPU: 2 cores
- RAM: 4GB
- Disk: 20GB

**Recommended (Production):**
- CPU: 4+ cores
- RAM: 8GB+
- Disk: 50GB+ SSD

## Backup & Restore

### Database Backup
```bash
# Linux/Mac
./scripts/backup.sh

# Windows
.\scripts\backup.ps1
```

Backups are stored in `./backups` with automatic 30-day retention.

### Restore
```bash
docker-compose -f docker-compose.production.yml exec -T postgres psql \
    -U suproxy_prod -d suproxy_prod < backup.sql
```

## Contributing

1. Follow Clean Architecture principles
2. Write unit tests for new features
3. Update documentation
4. Follow Go best practices
5. Use structured logging
6. Add metrics for monitoring

## License

Copyright © 2024 SuProxy. All rights reserved.

## Support

For production deployment assistance, see [DEPLOYMENT.md](DEPLOYMENT.md)

For API documentation, see [docs/swagger.md](docs/swagger.md)

For authentication details, see [docs/authentication.md](docs/authentication.md)
