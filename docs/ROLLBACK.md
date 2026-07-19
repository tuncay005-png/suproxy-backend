# 🔄 Rollback Guide

## Overview

This document describes how to rollback a deployment to a previous version when issues are detected in production.

## Quick Rollback

### Via GitHub UI (Recommended)

1. Go to **Actions** → **Deploy to Production**
2. Click **Run workflow**
3. Set **version** to previous working version (e.g., `v1.0.41`)
4. Set **servers** to `all` or specific servers
5. Click **Run workflow**

### Via Command Line

```bash
# Trigger workflow via GitHub CLI
gh workflow run deploy.yml \
  -f version=v1.0.41 \
  -f servers=all
```

## Finding Previous Versions

### Method 1: GitHub Releases

View all releases at: `https://github.com/<owner>/suproxy-backend/releases`

### Method 2: GitHub Container Registry

```bash
# List all tags
docker image ls ghcr.io/tuncay005-png/suproxy-backend

# Or view on GitHub
# Navigate to: Packages → suproxy-backend → Tags
```

### Method 3: Git Tags

```bash
git fetch --tags
git tag -l | sort -V | tail -10
```

## Rollback Strategies

### Strategy 1: Specific Version Rollback

Best when you know which version was working:

```bash
# Deploy previous stable version
Version: v1.0.41
Servers: all
```

### Strategy 2: SHA-based Rollback

When you know the exact commit that was working:

```bash
# Find commit SHA from git history
git log --oneline -10

# Deploy specific commit
Version: sha-a1b2c3d
Servers: all
```

### Strategy 3: Latest Tag Rollback

Use the `latest` tag if it's been confirmed working:

```bash
Version: latest
Servers: all
```

## Rollback to Specific Servers

If only one region has issues:

```bash
# Rollback only Finland
Version: v1.0.41
Servers: finland

# Rollback multiple regions
Version: v1.0.41
Servers: finland,germany
```

## Emergency Manual Rollback

If GitHub Actions is unavailable, rollback directly on VPS:

### Step 1: SSH to affected server

```bash
ssh user@vps-host
cd /opt/suproxy
```

### Step 2: Pull previous version

```bash
# Edit .env.production
nano .env.production

# Change VERSION to previous version
VERSION=v1.0.41

# Save and exit (Ctrl+X, Y, Enter)
```

### Step 3: Redeploy

```bash
# Run deploy script
./scripts/deploy.sh
```

### Step 4: Verify

```bash
# Check health
curl http://localhost:8080/health

# Check logs
docker-compose -f docker-compose.production.yml logs -f api
```

## Automated Rollback (Future)

### Health Check Failure Rollback

Future implementation will automatically rollback if:
- Health checks fail after deployment
- Error rate exceeds threshold
- Response time degrades significantly

### Monitoring Alert Rollback

Future implementation will trigger rollback when:
- Prometheus alerts fire
- Error rate > 5%
- Downtime detected

## Rollback Verification

After rollback, verify:

### 1. Health Check

```bash
curl http://<VPS_IP>:8080/health
```

Expected response:
```json
{"status":"ok","timestamp":"2026-07-19T..."}
```

### 2. Container Status

```bash
docker-compose -f docker-compose.production.yml ps
```

All containers should show `Up` status.

### 3. Application Logs

```bash
docker-compose -f docker-compose.production.yml logs --tail=100 api
```

No errors should appear.

### 4. Monitoring Dashboards

Check Grafana for:
- Request rate
- Error rate
- Response times
- Resource usage

## Rollback Testing

### Test Rollback Procedure

Periodically test rollback:

```bash
# 1. Deploy current version
Version: latest
Servers: finland

# 2. Immediately rollback to previous
Version: v1.0.40
Servers: finland

# 3. Verify both deployments succeed
```

## Common Rollback Scenarios

### Scenario 1: Bad Database Migration

```bash
# Rollback to version before migration
Version: v1.0.40  # Before migration
Servers: all

# Then manually fix database if needed
ssh user@vps-host
# Run manual migration rollback
```

### Scenario 2: Performance Degradation

```bash
# Rollback to last known good performance
Version: v1.0.38
Servers: all

# Monitor performance in Grafana
```

### Scenario 3: Security Issue

```bash
# Immediately rollback to patched version
Version: v1.0.35  # Known secure version
Servers: all

# Then deploy fixed version when ready
```

### Scenario 4: Third-party Dependency Issue

```bash
# Rollback to version with working dependencies
Version: v1.0.39
Servers: all
```

## Rollback Best Practices

### Before Rollback

1. ✅ Identify the issue clearly
2. ✅ Determine last working version
3. ✅ Verify the target version image exists
4. ✅ Notify team of rollback

### During Rollback

1. ✅ Monitor logs during rollback
2. ✅ Check health checks pass
3. ✅ Verify all regions if rolling back all
4. ✅ Keep communication channel open

### After Rollback

1. ✅ Confirm application is stable
2. ✅ Review logs for root cause
3. ✅ Document the issue
4. ✅ Create fix in new version
5. ✅ Test fix thoroughly before redeployment

## Version Pinning

### Pin to Stable Version

In `.env.production`:

```bash
# Pin to known stable version instead of 'latest'
VERSION=v1.0.41
```

This prevents automatic updates to potentially unstable versions.

### Gradual Rollout

Deploy new versions gradually:

```bash
# Step 1: Deploy to one region first
Version: v1.0.42
Servers: finland

# Step 2: Monitor for issues
# Wait 30-60 minutes

# Step 3: Deploy to remaining regions
Version: v1.0.42
Servers: germany,turkey
```

## Rollback Metrics

Track rollback events:

- Rollback frequency
- Time to rollback
- Rollback success rate
- Issues requiring rollback

## Database Rollback

### Database Migration Rollback

If deployment includes database migrations:

```bash
# SSH to VPS
ssh user@vps-host
cd /opt/suproxy

# Access database
docker-compose -f docker-compose.production.yml exec postgres psql -U suproxy_prod

# Check migration version
SELECT version FROM schema_migrations ORDER BY version DESC LIMIT 1;

# Manually rollback migration if needed
# (Application should handle this automatically)
```

### Database Backup Restoration

If database corruption occurs:

```bash
# Restore from backup (see BACKUP.md)
docker-compose -f docker-compose.production.yml exec postgres \
  pg_restore -U suproxy_prod -d suproxy_prod /backups/backup_YYYYMMDD.dump
```

## Preventing Rollbacks

### Pre-deployment Checklist

- [ ] All tests pass locally
- [ ] Code review completed
- [ ] Staging deployment successful
- [ ] Database migrations tested
- [ ] Dependencies verified
- [ ] Performance tested

### Canary Deployment (Future)

Deploy to small percentage of traffic first:

```bash
# Deploy to 10% of servers
# Monitor for issues
# Gradually increase to 100%
```

### Blue/Green Deployment (Future)

Maintain two identical production environments:

```bash
# Deploy to Blue environment
# Switch traffic to Blue
# Keep Green as instant rollback option
```

## Related Documentation

- [DEPLOYMENT.md](./DEPLOYMENT.md) - Deployment procedures
- [BACKUP.md](./BACKUP.md) - Backup strategies
- [MULTISERVER.md](./MULTISERVER.md) - Multi-region setup
- [CI_CD_ARCHITECTURE.md](./CI_CD_ARCHITECTURE.md) - Workflow details
