# 🚀 SuProxy Backend

Enterprise-grade VPN management backend with automated CI/CD, multi-server deployment, and comprehensive monitoring.

## 🏗️ Architecture

```
GitHub Actions CI/CD
    ↓
Tests → Build → Deploy
    ↓
Multi-Region Servers
    ↓
High Availability
```

## ✨ Features

- 🔐 JWT Authentication & Authorization
- 📊 Real-time Metrics & Monitoring
- 🌍 Multi-Server Deployment
- 🔄 Automated Rollback
- 💾 Automated Backups
- 🐳 Docker Containerized
- 📈 Prometheus + Grafana
- 🚨 Health Checks
- 🔒 Security Scanning
- 📝 Comprehensive Logging

## 🚀 Quick Start

### Prerequisites

- Go 1.23+
- Docker & Docker Compose
- PostgreSQL 15+

### Local Development

```bash
# Clone repository
git clone https://github.com/tuncay005-png/suproxy-backend.git
cd suproxy-backend

# Install dependencies
go mod download

# Copy environment file
cp .env.example .env

# Run tests
go test -v ./...

# Run locally
go run cmd/api/main.go
```

### Docker Development

```bash
# Start all services
docker-compose up -d

# View logs
docker-compose logs -f

# Stop services
docker-compose down
```

## 📚 Documentation

### Core Documentation

- [API Quick Reference](./API_QUICK_REFERENCE.md) - API endpoints overview
- [Authentication](./docs/authentication.md) - Auth implementation
- [Testing](./docs/testing.md) - Testing strategy

### Infrastructure Documentation

- [CI/CD Architecture](./docs/CI_CD_ARCHITECTURE.md) - Workflow details
- [Deployment Guide](./docs/DEPLOYMENT.md) - How to deploy
- [Rollback Procedures](./docs/ROLLBACK.md) - Recovery strategies
- [Multi-Server Setup](./docs/MULTISERVER.md) - Geographic distribution
- [Backup & Recovery](./docs/BACKUP.md) - Data protection

## 🔄 CI/CD Pipeline

### Automated Workflow

```
Push to main
    ↓
Run Tests (test.yml)
    ↓ (pass required)
Build Docker Image (build.yml)
    ↓ (auto-trigger)
Deploy to Servers (deploy.yml)
    ↓
Health Checks
    ↓
Production Live ✅
```

### Workflow Files

- **test.yml** - Quality assurance (tests, lint, security)
- **build.yml** - Docker image building
- **deploy.yml** - Multi-server deployment
- **release.yml** - GitHub releases

## 🌍 Deployment

### Automatic Deployment

Deployments trigger automatically on push to `main`:

```bash
git push origin main
# → Tests run automatically
# → Docker image built
# → Deployed to all servers
```

### Manual Deployment

Deploy specific version:

```bash
# Via GitHub UI:
Actions → Deploy to Production → Run workflow
- Version: v1.0.42
- Servers: all
```

### Server Regions

- **Finland** 🇫🇮 - Primary (active)
- **Germany** 🇩🇪 - Future
- **Turkey** 🇹🇷 - Future



## 🐳 Docker Images

### Image Tags

Every build produces three tags:

```bash
# Latest
ghcr.io/tuncay005-png/suproxy-backend:latest

# Semantic version
ghcr.io/tuncay005-png/suproxy-backend:v1.0.42

# Commit SHA
ghcr.io/tuncay005-png/suproxy-backend:sha-a1b2c3d
```

### Pull Image

```bash
docker pull ghcr.io/tuncay005-png/suproxy-backend:latest
```

## 🧪 Testing

### Run All Tests

```bash
# Unit tests
go test -v -short ./...

# Integration tests
INTEGRATION_TEST=true go test -v ./...

# With coverage
go test -v -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### Linting

```bash
# Run linter
golangci-lint run
```

### Security Scan

```bash
# Run security scanner
gosec ./...
```

## 📊 Monitoring

### Access Monitoring

```bash
# Prometheus
http://localhost:9090

# Grafana
http://localhost:3000
Username: admin
Password: (see .env.production)
```

### Health Check

```bash
curl http://localhost:8080/health
```

## 🔐 Security

- ✅ JWT-based authentication
- ✅ HTTPS/TLS support
- ✅ Security scanning (gosec)
- ✅ Dependency scanning
- ✅ Non-root Docker containers
- ✅ Read-only filesystems
- ✅ Resource limits
- ✅ Secret management

## 🛠️ Development

### Project Structure

```
suproxy-backend/
├── cmd/
│   └── api/              # Application entry point
├── docs/                 # Documentation
├── internal/             # Private application code
├── migrations/           # Database migrations
├── scripts/              # Deployment & utility scripts
├── .github/
│   └── workflows/        # CI/CD workflows
├── docker-compose.yml    # Development compose
└── docker-compose.production.yml  # Production compose
```

### Environment Variables

```bash
# Server
SUPROXY_SERVER_ADDRESS=:8080

# Database
SUPROXY_DATABASE_HOST=localhost
SUPROXY_DATABASE_PORT=5432
SUPROXY_DATABASE_USER=suproxy
SUPROXY_DATABASE_PASSWORD=password
SUPROXY_DATABASE_DBNAME=suproxy

# JWT
SUPROXY_JWT_SECRET_KEY=your-secret-key
SUPROXY_JWT_ACCESS_TOKEN_EXPIRY=15
SUPROXY_JWT_REFRESH_TOKEN_EXPIRY=168
```



## 🚨 Troubleshooting

### Deployment Failed

```bash
# Check workflow logs
# GitHub → Actions → Failed workflow → View logs

# SSH to server
ssh user@vps-host

# Check container status
docker-compose -f docker-compose.production.yml ps

# View logs
docker-compose -f docker-compose.production.yml logs -f api
```

### Rollback

```bash
# Via GitHub UI:
Actions → Deploy to Production → Run workflow
- Version: v1.0.41  # Previous version
- Servers: all
```

### Database Issues

```bash
# Check database logs
docker-compose -f docker-compose.production.yml logs postgres

# Access database
docker-compose -f docker-compose.production.yml exec postgres psql -U suproxy_prod

# Run migrations
docker-compose -f docker-compose.production.yml exec api /app/suproxy-api migrate
```

## 📈 Performance

### Resource Usage

- **API**: 2 vCPU, 2GB RAM (typical)
- **PostgreSQL**: 1 vCPU, 1GB RAM (typical)
- **Prometheus**: 0.5 vCPU, 512MB RAM
- **Grafana**: 0.5 vCPU, 256MB RAM

### Scaling

```bash
# Vertical scaling (increase resources)
# Edit docker-compose.production.yml:
deploy:
  resources:
    limits:
      cpus: '4'
      memory: 4G

# Horizontal scaling (add servers)
# See docs/MULTISERVER.md
```

## 🤝 Contributing

1. Fork the repository
2. Create feature branch (`git checkout -b feature/amazing-feature`)
3. Commit changes (`git commit -m 'Add amazing feature'`)
4. Push to branch (`git push origin feature/amazing-feature`)
5. Open Pull Request

### Development Guidelines

- Write tests for new features
- Follow existing code style
- Update documentation
- Ensure CI passes

## 📄 License

This project is proprietary and confidential.

## 👥 Team

- **Maintainer**: Tuncay
- **Repository**: https://github.com/tuncay005-png/suproxy-backend

## 🔗 Links

- [GitHub Repository](https://github.com/tuncay005-png/suproxy-backend)
- [Container Registry](https://github.com/tuncay005-png/suproxy-backend/pkgs/container/suproxy-backend)
- [Issues](https://github.com/tuncay005-png/suproxy-backend/issues)
- [Releases](https://github.com/tuncay005-png/suproxy-backend/releases)

## 📞 Support

For support, please open an issue on GitHub or contact the maintainers.

---

**Built with ❤️ for enterprise-grade VPN management**
