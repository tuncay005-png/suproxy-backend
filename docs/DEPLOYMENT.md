# SuProxy Production Deployment Guide

**Current Version:** 1.0.0  
**Last Updated:** 2024-08-21

---

## Quick Start

```bash
# 1. Clone and configure
git clone https://github.com/tuncay005-png/suproxy-backend.git
cd suproxy-backend
cp .env.production.example .env.production

# 2. Configure environment (see Configuration section)
nano .env.production

# 3. Set up HTTPS (see HTTPS Setup section)
./scripts/setup-https-traefik.sh
# or
./scripts/setup-https-nginx.sh

# 4. Deploy
docker-compose -f docker-compose.production-https.yml up -d

# 5. Verify
curl https://api.yourdomain.com/health
```

---

## Prerequisites

### Required
- ✅ Docker 20.10+ and Docker Compose V2
- ✅ Domain name with DNS configured
- ✅ Server with 4GB+ RAM, 2+ CPU cores
- ✅ Ports 80, 443 open in firewall
- ✅ Email address for SSL certificates

### Domain Configuration
```bash
# Required DNS A records:
api.yourdomain.com      → YOUR_SERVER_IP
grafana.yourdomain.com  → YOUR_SERVER_IP
prometheus.yourdomain.com → YOUR_SERVER_IP (optional)
traefik.yourdomain.com  → YOUR_SERVER_IP (if using Traefik)
```

**Verify DNS:**
```bash
nslookup api.yourdomain.com
# Should return your server IP
```

---

## Configuration

### 1. Environment Variables

Copy and edit production environment file:

```bash
cp .env.production.example .env.production
nano .env.production
```

### Required Configuration

#### Domain & SSL
```bash
DOMAIN=yourdomain.com
ACME_EMAIL=admin@yourdomain.com
```

#### Database (CRITICAL)
```bash
DB_USER=suproxy_prod
DB_PASSWORD=$(openssl rand -base64 32)  # Generate secure password
DB_NAME=suproxy_prod
DB_SSLMODE=require
```

#### JWT Authentication (CRITICAL)
```bash
# Generate secure secret
JWT_SECRET=$(openssl rand -base64 32)

# Must be generated before production
# Application will FAIL to start if default value is used
JWT_ACCESS_EXPIRY=15    # minutes
JWT_REFRESH_EXPIRY=168  # hours (7 days)
```

#### Monitoring
```bash
GRAFANA_USER=admin
GRAFANA_PASSWORD=$(openssl rand -base64 16)

# Generate basic auth hashes
# htpasswd -nb admin yourpassword
TRAEFIK_DASHBOARD_AUTH='admin:$apr1$...'
MONITORING_AUTH='admin:$apr1$...'
```

### Security Validation

The application includes production security checks:

**JWT Secret Validation:**
```bash
# In production, these will FAIL startup:
JWT_SECRET=change-me-in-production  # ❌ Rejected
JWT_SECRET=""                       # ❌ Rejected  
JWT_SECRET="   "                    # ❌ Rejected

# Valid production secret:
JWT_SECRET=$(openssl rand -base64 32)  # ✅ Accepted
```

**Cookie Security:**
```bash
COOKIE_SECURE=true
COOKIE_SAMESITE=strict
```

---

## HTTPS Setup

**⚠️ HTTPS is REQUIRED for production.** Do not run production over HTTP.

### Option 1: Traefik (Recommended)

**Advantages:**
- ✅ Automatic SSL certificate management
- ✅ Auto-renewal built-in
- ✅ Docker-native
- ✅ Monitoring dashboard

```bash
# Run automated setup
chmod +x scripts/setup-https-traefik.sh
./scripts/setup-https-traefik.sh

# Follow prompts for staging → production
```

**Deployment:**
```bash
docker-compose -f docker-compose.production-https.yml up -d
```

### Option 2: Nginx

**Advantages:**
- ✅ Traditional and widely known
- ✅ Fine-grained control
- ✅ High performance

```bash
# Run automated setup
chmod +x scripts/setup-https-nginx.sh
./scripts/setup-https-nginx.sh
```

**Deployment:**
```bash
docker-compose -f docker-compose.production-nginx.yml up -d
```

### Complete HTTPS Documentation

For detailed HTTPS setup, troubleshooting, and security configuration:

📖 **See: [docs/HTTPS_SETUP.md](./HTTPS_SETUP.md)**

Covers:
- Step-by-step setup for both Traefik and Nginx
- SSL certificate management
- Security headers configuration
- HTTP → HTTPS redirect
- Certificate renewal
- Troubleshooting
- Monitoring

---

## Deployment Architecture

### Production Stack

```
┌─────────────────────────────────────────────────────┐
│                  Internet (HTTPS)                    │
└──────────────────────┬──────────────────────────────┘
                       │
                       ▼
┌─────────────────────────────────────────────────────┐
│            Reverse Proxy (Traefik/Nginx)            │
│  • SSL Termination                                  │
│  • Security Headers                                 │
│  • Rate Limiting                                    │
│  • HTTP → HTTPS Redirect                            │
└──────────────┬──────────────────────────────────────┘
               │
        ┌──────┴──────┬──────────┬──────────┐
        ▼             ▼          ▼          ▼
   ┌────────┐   ┌──────────┐ ┌──────┐ ┌─────────┐
   │  API   │   │PostgreSQL│ │Prome-│ │ Grafana │
   │ :8080  │   │  :5432   │ │theus │ │  :3000  │
   └────────┘   └──────────┘ └──────┘ └─────────┘
```

### Services

| Service | Purpose | Exposed |
|---------|---------|---------|
| Traefik/Nginx | Reverse proxy, SSL | Ports 80, 443 |
| API | Backend application | Internal :8080 |
| PostgreSQL | Database | Internal :5432 |
| Prometheus | Metrics collection | Via proxy (optional) |
| Grafana | Monitoring dashboards | Via proxy |

---

## Database Setup

### Migrations

Migrations run automatically on startup. To run manually:

```bash
# Check current version
docker-compose -f docker-compose.production-https.yml exec api \
  migrate -path /app/migrations -database "postgresql://..." version

# Run migrations manually (if needed)
docker-compose -f docker-compose.production-https.yml exec api \
  migrate -path /app/migrations -database "postgresql://..." up
```

### Backup

**Automated backup (recommended):**
```bash
# Create backup script
nano /opt/suproxy/scripts/backup-db.sh

# Add to cron
crontab -e
0 2 * * * /opt/suproxy/scripts/backup-db.sh
```

**Manual backup:**
```bash
docker exec suproxy-postgres-prod pg_dump -U suproxy_prod suproxy_prod > backup.sql
```

**Restore:**
```bash
docker exec -i suproxy-postgres-prod psql -U suproxy_prod suproxy_prod < backup.sql
```

---

## Monitoring

### Health Checks

```bash
# API health
curl https://api.yourdomain.com/health

# Database health
docker exec suproxy-postgres-prod pg_isready -U suproxy_prod

# All container health
docker-compose -f docker-compose.production-https.yml ps
```

### Grafana Dashboards

Access: `https://grafana.yourdomain.com`

Pre-configured dashboards for:
- API performance metrics
- Database connections
- Request rates and latencies
- Error rates
- Resource usage

### Prometheus Metrics

Access: `https://prometheus.yourdomain.com` (if exposed)

Metrics available:
- `http_requests_total` - Total requests
- `http_request_duration_seconds` - Request latency
- `database_connections_open` - Active DB connections
- `jwt_tokens_generated_total` - JWT generation count

---

## Security Features

### Implemented Security

✅ **HTTPS/TLS**
- TLS 1.2 and 1.3 only
- Strong cipher suites
- Perfect Forward Secrecy

✅ **Security Headers**
- HSTS (Strict-Transport-Security)
- X-Frame-Options: DENY
- X-Content-Type-Options: nosniff
- X-XSS-Protection
- Content-Security-Policy

✅ **Authentication**
- JWT with secure secret validation
- Production startup fails if default secret used
- Secure cookie configuration (HttpOnly, Secure, SameSite)

✅ **Rate Limiting**
- Auth endpoints: 5 req/min per IP
- Admin endpoints: 100 req/min per user
- Additional reverse proxy limits

✅ **CORS**
- Configurable allowed origins
- No wildcard (*) in production

✅ **Database**
- SSL/TLS connections required
- Prepared statements (SQL injection protection)
- Connection pooling with limits

### Security Checklist

Before production:
- [ ] HTTPS configured and working
- [ ] JWT secret generated (not default)
- [ ] Database password changed
- [ ] Grafana password changed
- [ ] CORS origins configured
- [ ] Firewall configured (only 80, 443 open)
- [ ] Security headers verified
- [ ] SSL Labs test shows A/A+

---

## Operations

### Starting Services

```bash
# Start all services
docker-compose -f docker-compose.production-https.yml up -d

# Start specific service
docker-compose -f docker-compose.production-https.yml up -d api

# View logs
docker-compose -f docker-compose.production-https.yml logs -f

# View logs for specific service
docker-compose -f docker-compose.production-https.yml logs -f api
```

### Stopping Services

```bash
# Stop all services
docker-compose -f docker-compose.production-https.yml down

# Stop but keep volumes
docker-compose -f docker-compose.production-https.yml stop
```

### Updating Application

```bash
# Pull latest image
docker-compose -f docker-compose.production-https.yml pull api

# Restart with new image
docker-compose -f docker-compose.production-https.yml up -d api

# Verify
curl https://api.yourdomain.com/health
```

### Scaling (Vertical)

Edit resource limits in `docker-compose.production-https.yml`:

```yaml
api:
  deploy:
    resources:
      limits:
        cpus: '8'      # Increase CPUs
        memory: 8G     # Increase memory
      reservations:
        cpus: '4'
        memory: 4G
```

Apply changes:
```bash
docker-compose -f docker-compose.production-https.yml up -d
```

---

## Troubleshooting

### Service Won't Start

```bash
# Check logs
docker-compose -f docker-compose.production-https.yml logs api

# Common issues:
# 1. JWT secret not set (SECURITY ERROR in logs)
# 2. Database connection failed (check DB_PASSWORD)
# 3. Port already in use
# 4. SSL certificate not obtained
```

### Database Connection Issues

```bash
# Test database connection
docker exec suproxy-postgres-prod psql -U suproxy_prod -c "SELECT 1;"

# Check database logs
docker logs suproxy-postgres-prod

# Verify database is running
docker ps | grep postgres
```

### SSL Certificate Issues

See [docs/HTTPS_SETUP.md](./HTTPS_SETUP.md) for comprehensive SSL troubleshooting.

### Performance Issues

```bash
# Check resource usage
docker stats

# Check database connections
docker exec suproxy-postgres-prod psql -U suproxy_prod -c \
  "SELECT count(*) FROM pg_stat_activity;"

# Check slow queries
docker exec suproxy-postgres-prod psql -U suproxy_prod -c \
  "SELECT * FROM pg_stat_statements ORDER BY total_time DESC LIMIT 10;"
```

---

## Production Checklist

### Before Deployment

- [ ] All prerequisites met
- [ ] DNS configured and propagated
- [ ] `.env.production` configured with secure values
- [ ] SSL certificates obtained
- [ ] Firewall configured
- [ ] Database backup strategy in place

### After Deployment

- [ ] Health check passes
- [ ] HTTPS working correctly
- [ ] Security headers present
- [ ] Monitoring dashboards accessible
- [ ] Load testing completed
- [ ] Frontend can connect
- [ ] Authentication working
- [ ] Database backups tested

### First Week

- [ ] Monitor logs daily
- [ ] Check certificate auto-renewal
- [ ] Verify backup completion
- [ ] Review metrics and alerts
- [ ] Test disaster recovery

---

## CI/CD Integration

The project includes GitHub Actions workflows:

1. **Tests** - Run on every push
2. **Build** - Build and push Docker image (after tests pass)
3. **Deploy** - Deploy to production server

See [PRODUCTION_READY.md](../PRODUCTION_READY.md) for CI/CD documentation.

---

## Support & Resources

### Documentation
- [HTTPS Setup Guide](./HTTPS_SETUP.md) - Detailed SSL/TLS configuration
- [Load Testing Guide](./LOAD_TESTING.md) - Performance testing procedures
- [Production Hardening Audit](../PRODUCTION_HARDENING_AUDIT.md) - Security assessment

### External Resources
- [Docker Docs](https://docs.docker.com/)
- [Traefik Docs](https://doc.traefik.io/traefik/)
- [Nginx Docs](https://nginx.org/en/docs/)
- [Let's Encrypt](https://letsencrypt.org/)

### Getting Help

1. Check logs: `docker-compose logs -f`
2. Review documentation in `docs/` directory
3. Check GitHub Issues
4. Review security audit: `PRODUCTION_HARDENING_AUDIT.md`

---

**Production Status:** ✅ Ready for Deployment  
**Security:** ✅ Hardened  
**HTTPS:** ✅ Configured  
**Monitoring:** ✅ Enabled
