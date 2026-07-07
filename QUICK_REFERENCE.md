# Quick Reference Guide

## Common Operations

### Start Services

```bash
# Development
docker-compose up -d

# Production
docker-compose -f docker-compose.production.yml up -d

# With rebuild
docker-compose up -d --build
```

### Stop Services

```bash
# Development
docker-compose down

# Production
docker-compose -f docker-compose.production.yml down

# With volume cleanup
docker-compose down -v
```

### View Logs

```bash
# All services
docker-compose logs -f

# Specific service
docker-compose logs -f api

# Last 100 lines
docker-compose logs --tail=100 api

# Production
docker-compose -f docker-compose.production.yml logs -f api
```

### Check Status

```bash
# Container status
docker-compose ps

# Health check
curl http://localhost:8080/health

# Detailed health (with auth)
curl -H "Authorization: Bearer YOUR_TOKEN" http://localhost:8080/api/v1/admin/system/health
```

### Database Operations

```bash
# Connect to database
docker-compose exec postgres psql -U suproxy -d suproxy

# Production database
docker-compose -f docker-compose.production.yml exec postgres psql -U suproxy_prod -d suproxy_prod

# Run SQL file
docker-compose exec -T postgres psql -U suproxy -d suproxy < query.sql

# Backup
./scripts/backup.sh  # Linux/Mac
.\scripts\backup.ps1  # Windows

# Restore
docker-compose exec -T postgres psql -U suproxy -d suproxy < backup.sql
```

### Restart Services

```bash
# Restart single service
docker-compose restart api

# Restart all services
docker-compose restart

# Production
docker-compose -f docker-compose.production.yml restart api
```

### Update & Redeploy

```bash
# Pull latest code
git pull

# Rebuild and restart
docker-compose up -d --build

# Production
./scripts/deploy.sh  # Linux/Mac
.\scripts\deploy.ps1  # Windows
```

### Cleanup

```bash
# Remove stopped containers
docker-compose down

# Remove with volumes (⚠️ deletes data)
docker-compose down -v

# Remove unused Docker resources
docker system prune

# Remove all unused (⚠️ aggressive)
docker system prune -a --volumes
```

## Troubleshooting

### API Not Starting

```bash
# Check logs
docker-compose logs api

# Check if port is in use
netstat -ano | findstr :8080  # Windows
lsof -i :8080                 # Linux/Mac

# Restart service
docker-compose restart api
```

### Database Connection Issues

```bash
# Check database status
docker-compose ps postgres

# Check database logs
docker-compose logs postgres

# Test connection
docker-compose exec postgres pg_isready -U suproxy

# Verify environment variables
docker-compose exec api env | grep DATABASE
```

### Memory Issues

```bash
# Check resource usage
docker stats

# Restart service
docker-compose restart api

# Increase limits in docker-compose.production.yml
deploy:
  resources:
    limits:
      memory: 8G
```

### Disk Space Issues

```bash
# Check disk usage
df -h                    # Linux/Mac
Get-PSDrive C | fl       # Windows (PowerShell)

# Docker disk usage
docker system df

# Clean up Docker
docker system prune -a

# Remove old backups
find ./backups -name "*.sql.gz" -mtime +30 -delete  # Linux/Mac
Get-ChildItem ./backups -Filter *.sql.zip | Where-Object {$_.LastWriteTime -lt (Get-Date).AddDays(-30)} | Remove-Item  # Windows
```

### Performance Issues

```bash
# Check metrics
curl http://localhost:8080/metrics

# View Prometheus
open http://localhost:9090

# View Grafana
open http://localhost:3000

# Check database connections
docker-compose exec postgres psql -U suproxy -d suproxy -c "SELECT count(*) FROM pg_stat_activity;"
```

## Development Workflow

### Local Development

```bash
# Start database only
docker-compose up -d postgres

# Run API locally
go run ./cmd/api

# Run with hot reload (if using Air)
air
```

### Run Tests

```bash
# All tests
go test ./...

# Specific package
go test ./internal/application/usecase/auth/...

# With coverage
go test -cover ./...

# Verbose
go test -v ./...
```

### Code Quality

```bash
# Format code
go fmt ./...

# Vet code
go vet ./...

# Lint (requires golangci-lint)
golangci-lint run ./...

# Security scan (requires gosec)
gosec ./...
```

### Build

```bash
# Build binary
go build -o suproxy-api ./cmd/api

# Build for Linux (from Windows/Mac)
GOOS=linux GOARCH=amd64 go build -o suproxy-api ./cmd/api

# Build Docker image
docker build -t suproxy/backend:dev .

# Multi-platform build
docker buildx build --platform linux/amd64,linux/arm64 -t suproxy/backend:latest .
```

## Monitoring

### Metrics

```bash
# Prometheus metrics
curl http://localhost:8080/metrics

# Prometheus UI
open http://localhost:9090

# Sample query (HTTP request rate)
rate(http_requests_total[5m])

# Sample query (Error rate)
rate(http_requests_total{status=~"5.."}[5m])
```

### Grafana

```bash
# Access Grafana
open http://localhost:3000

# Default credentials
Username: admin
Password: (from .env.production)

# Import dashboard
Dashboard > Import > Upload JSON
```

### Logs

```bash
# Follow logs
docker-compose logs -f api

# Search logs
docker-compose logs api | grep ERROR

# JSON logs with jq
docker-compose logs api | jq '.level, .msg'
```

## Environment Variables

### View Current Config

```bash
# In container
docker-compose exec api env | grep SUPROXY

# Local
printenv | grep SUPROXY  # Linux/Mac
Get-ChildItem Env: | Where-Object {$_.Name -like "SUPROXY*"}  # Windows
```

### Update Config

```bash
# Edit .env or .env.production
vi .env.production

# Restart to apply changes
docker-compose -f docker-compose.production.yml restart api
```

## Security

### Generate Secure Secret

```bash
# JWT Secret (64+ characters)
openssl rand -base64 64  # Linux/Mac

# Windows (PowerShell)
[Convert]::ToBase64String((1..64 | ForEach-Object { Get-Random -Minimum 0 -Maximum 256 }))
```

### Check Security

```bash
# Container security scan
docker scan suproxy/backend:latest

# Go security scan
gosec ./...

# Dependency audit
go list -json -m all | nancy sleuth
```

## Backup & Restore

### Manual Backup

```bash
# Database backup
docker-compose exec postgres pg_dump -U suproxy -d suproxy > backup_$(date +%Y%m%d).sql

# With compression
docker-compose exec postgres pg_dump -U suproxy -d suproxy | gzip > backup_$(date +%Y%m%d).sql.gz
```

### Manual Restore

```bash
# From SQL file
docker-compose exec -T postgres psql -U suproxy -d suproxy < backup.sql

# From compressed file
gunzip -c backup.sql.gz | docker-compose exec -T postgres psql -U suproxy -d suproxy
```

### Automated Backup

```bash
# Run backup script
./scripts/backup.sh  # Linux/Mac
.\scripts\backup.ps1  # Windows

# Add to cron (Linux)
0 2 * * * /path/to/suproxy-backend/scripts/backup.sh

# Add to Task Scheduler (Windows)
# Use Task Scheduler GUI to schedule .\scripts\backup.ps1
```

## API Testing

### Using curl

```bash
# Health check
curl http://localhost:8080/health

# Register
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"username":"testuser","email":"test@example.com","password":"Test123!@#"}'

# Login
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"testuser","password":"Test123!@#"}'

# Get current user (with token)
curl http://localhost:8080/api/v1/auth/me \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"

# Admin endpoint
curl http://localhost:8080/api/v1/admin/users \
  -H "Authorization: Bearer ADMIN_ACCESS_TOKEN"
```

### Using httpie (if installed)

```bash
# Register
http POST :8080/api/v1/auth/register username=testuser email=test@example.com password=Test123!@#

# Login
http POST :8080/api/v1/auth/login username=testuser password=Test123!@#

# With auth
http :8080/api/v1/auth/me Authorization:"Bearer YOUR_TOKEN"
```

## Service URLs

- **API**: http://localhost:8080
- **Health**: http://localhost:8080/health
- **Metrics**: http://localhost:8080/metrics
- **Prometheus**: http://localhost:9090
- **Grafana**: http://localhost:3000

## Important Files

- `.env` - Development environment variables
- `.env.production` - Production environment variables
- `docker-compose.yml` - Development compose file
- `docker-compose.production.yml` - Production compose file
- `Dockerfile` - Container image definition
- `configs/config.yaml` - Application config template
- `DEPLOYMENT.md` - Detailed deployment guide
- `PRODUCTION_CHECKLIST.md` - Pre-deployment checklist

## Support

For detailed information, see:
- [DEPLOYMENT.md](DEPLOYMENT.md) - Full deployment guide
- [README.md](README.md) - Project overview
- [docs/authentication.md](docs/authentication.md) - Auth documentation
- [docs/swagger.md](docs/swagger.md) - API documentation
