# 🚀 Deployment Guide

## Overview

SuProxy Backend uses a fully automated CI/CD pipeline with GitHub Actions. This document describes the deployment process and architecture.

## Deployment Architecture

### Workflow Pipeline

```
Code Push to main
    ↓
test.yml (Quality Gate)
    ↓ (must pass)
build.yml (Build Docker Image)
    ↓ (automatic)
deploy.yml (Deploy to Servers)
    ↓ (health checks)
Production Running
```

### Multi-Server Support

The deployment system supports multiple servers in different regions:

- **Finland** (Primary) - Current production server
- **Germany** - Future
- **Turkey** - Future
- **USA** - Future
- **Japan** - Future
- **Singapore** - Future

## Automatic Deployment

### Trigger Flow

1. **Push to `main` branch**
2. **Tests run automatically** (`test.yml`)
   - Unit tests
   - Integration tests
   - Linting
   - Security scans
3. **Docker image built** (`build.yml`)
   - Only if tests pass
   - Tagged with `latest`, version, and SHA
4. **Deployment executed** (`deploy.yml`)
   - Only if build succeeds
   - Deploys to all configured servers
   - Runs health checks

## Manual Deployment

### Deploy Specific Version

```bash
# Via GitHub UI:
Actions → Deploy to Production → Run workflow
- Version: v1.0.42 (or 'latest')
- Servers: all (or finland,germany,turkey)
```

### Deploy to Specific Servers

```bash
# Deploy only to Finland and Germany
Servers: finland,germany

# Deploy to all servers
Servers: all
```

## Version Management

### Version Tagging

Every build produces three tags:

1. **`latest`** - Always the most recent build
2. **`v1.0.X`** - Semantic version (auto-incremented)
3. **`sha-abc123`** - Git commit identifier

### Pull Specific Version

```bash
# Latest
docker pull ghcr.io/tuncay005-png/suproxy-backend:latest

# Specific version
docker pull ghcr.io/tuncay005-png/suproxy-backend:v1.0.42

# Specific commit
docker pull ghcr.io/tuncay005-png/suproxy-backend:sha-a1b2c3d
```

## Server Setup

### VPS Requirements

Each VPS must have:
- Docker installed
- Docker Compose installed
- `/opt/suproxy` directory structure
- `.env.production` configuration
- SSH access configured

### Directory Structure

```
/opt/suproxy/
├── .env.production              # Environment configuration
├── docker-compose.production.yml
├── scripts/
│   └── deploy.sh               # Deployment script
├── backups/                    # Database backups
├── prometheus/
│   └── prometheus.yml
└── grafana/
    └── provisioning/
```

## Health Checks

### Automatic Health Verification

After each deployment, the system:

1. Waits 5 seconds for services to start
2. Checks if all Docker containers are running
3. Polls the `/health` endpoint (up to 60 seconds)
4. Verifies HTTP 200 response

### Manual Health Check

```bash
# On VPS
curl http://localhost:8080/health

# Expected response
{"status":"ok","timestamp":"2026-07-19T..."}
```

## Rollback

If a deployment fails:

1. **Automatic**: Health check fails → deployment marked as failed
2. **Manual**: Use previous version

```bash
# Via GitHub UI:
Actions → Deploy to Production → Run workflow
- Version: v1.0.41  # Previous working version
- Servers: all
```

## GitHub Secrets Configuration

### Current Server (Finland)

Uses legacy secrets with fallback:

```
VPS_FINLAND_HOST     (or VPS_HOST)
VPS_FINLAND_USER     (or VPS_USER)
VPS_FINLAND_KEY      (or SSH_PRIVATE_KEY)
VPS_FINLAND_PORT     (or VPS_PORT)
```

### Future Servers

Add these secrets in GitHub repository settings:

**Germany:**
```
VPS_GERMANY_HOST
VPS_GERMANY_USER
VPS_GERMANY_KEY
VPS_GERMANY_PORT
```

**Turkey:**
```
VPS_TURKEY_HOST
VPS_TURKEY_USER
VPS_TURKEY_KEY
VPS_TURKEY_PORT
```

**Additional servers:** Follow the same pattern with country name.

## Deployment Script (`deploy.sh`)

### What It Does

```bash
1. Validates environment variables
2. Pulls Docker image from GHCR
3. Stops existing containers
4. Starts new containers
5. Waits for health checks
6. Reports status
```

### Important: No Local Building

The VPS **never builds Docker images**. It only:
- Pulls pre-built images from GHCR
- Runs `docker-compose up -d`
- Verifies health

## Monitoring

### Production Monitoring Stack

Each VPS runs:

- **Prometheus** - Metrics collection (port 9090)
- **Grafana** - Visualization (port 3000)
- **API** - Application (port 8080)

### Access Monitoring

```bash
# Prometheus
http://<VPS_IP>:9090

# Grafana
http://<VPS_IP>:3000
Username: admin
Password: (from .env.production)
```

## Troubleshooting

### Deployment Failed

```bash
# SSH to VPS
ssh user@vps-host

# Check logs
cd /opt/suproxy
docker-compose -f docker-compose.production.yml logs -f api

# Check container status
docker-compose -f docker-compose.production.yml ps

# Restart manually
./scripts/deploy.sh
```

### Health Check Failed

```bash
# Check API logs
docker-compose -f docker-compose.production.yml logs api

# Check all services
docker-compose -f docker-compose.production.yml ps

# Verify database connection
docker-compose -f docker-compose.production.yml logs postgres
```

### Image Pull Failed

```bash
# Verify image exists in GHCR
docker manifest inspect ghcr.io/tuncay005-png/suproxy-backend:latest

# Re-authenticate if needed (on VPS)
echo $GITHUB_TOKEN | docker login ghcr.io -u USERNAME --password-stdin
```

## Best Practices

### Before Deploying

1. ✅ Ensure all tests pass locally
2. ✅ Review code changes
3. ✅ Check build logs in GitHub Actions
4. ✅ Verify the version to deploy

### After Deploying

1. ✅ Monitor health checks
2. ✅ Check application logs
3. ✅ Verify API endpoints
4. ✅ Monitor Grafana dashboards

### Emergency Procedures

1. **Quick Rollback**: Deploy previous version immediately
2. **Stop Deployment**: Cancel workflow in GitHub Actions
3. **Manual Intervention**: SSH to VPS and run commands directly

## Adding New Servers

To add a new server (e.g., "australia"):

1. **Set up VPS** with required software
2. **Add secrets** to GitHub:
   - `VPS_AUSTRALIA_HOST`
   - `VPS_AUSTRALIA_USER`
   - `VPS_AUSTRALIA_KEY`
   - `VPS_AUSTRALIA_PORT`
3. **Update deploy.yml** (add to server mapping)
4. **Deploy**: Use `servers: australia` or `all`

No workflow logic changes required after initial setup!

## Security

### Secrets Management

- Never commit secrets to repository
- Use GitHub Secrets for all credentials
- Rotate SSH keys regularly
- Use strong passwords for services

### Network Security

- Configure firewall rules on VPS
- Use SSH key authentication only
- Disable password authentication
- Keep ports 8080, 9090, 3000 behind reverse proxy

### Docker Security

- Images run as non-root user
- Read-only filesystem where possible
- Resource limits configured
- Security hardening enabled

## Related Documentation

- [CI/CD Architecture](./CI_CD_ARCHITECTURE.md) - Workflow details
- [SERVER_SETUP.md](./SERVER_SETUP.md) - VPS configuration (future)
- [ROLLBACK.md](./ROLLBACK.md) - Rollback procedures (future)
- [MULTISERVER.md](./MULTISERVER.md) - Multi-region setup (future)
- [BACKUP.md](./BACKUP.md) - Backup strategies (future)
