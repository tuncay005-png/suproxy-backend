# ✅ Deployment Checklist

## Pre-Deployment Verification

### GitHub Secrets

Verify all required secrets are configured in GitHub repository settings:

**Settings → Secrets and variables → Actions → Repository secrets**

#### Current Production (Finland)
- [ ] `VPS_FINLAND_HOST` (or legacy `VPS_HOST`)
- [ ] `VPS_FINLAND_USER` (or legacy `VPS_USER`)
- [ ] `VPS_FINLAND_KEY` (or legacy `SSH_PRIVATE_KEY`)
- [ ] `VPS_FINLAND_PORT` (or legacy `VPS_PORT`)
- [ ] `GITHUB_TOKEN` (automatically provided)

#### Future Servers (Optional)
- [ ] `VPS_GERMANY_HOST`, `VPS_GERMANY_USER`, `VPS_GERMANY_KEY`, `VPS_GERMANY_PORT`
- [ ] `VPS_TURKEY_HOST`, `VPS_TURKEY_USER`, `VPS_TURKEY_KEY`, `VPS_TURKEY_PORT`

### VPS Configuration

SSH to your VPS and verify:

```bash
ssh user@vps-host
cd /opt/suproxy
```

- [ ] Docker installed: `docker --version`
- [ ] Docker Compose installed: `docker compose version`
- [ ] Directory structure exists: `/opt/suproxy/`
- [ ] Deploy script exists: `/opt/suproxy/scripts/deploy.sh`
- [ ] Deploy script is executable: `chmod +x scripts/deploy.sh`
- [ ] `.env.production` exists and configured
- [ ] `docker-compose.production.yml` exists
- [ ] Logged into GHCR: `docker pull ghcr.io/tuncay005-png/suproxy-backend:latest`

### GitHub Container Registry (GHCR)

Verify GHCR permissions:

- [ ] Repository has packages write permission
- [ ] GHCR images are visible: https://github.com/tuncay005-png?tab=packages
- [ ] Can pull images manually: `docker pull ghcr.io/tuncay005-png/suproxy-backend:latest`

## First Deployment Test

### Step 1: Commit and Push

```bash
# Make a small change (e.g., update README)
git add .
git commit -m "test: verify CI/CD pipeline"
git push origin main
```

### Step 2: Monitor GitHub Actions

Go to: **GitHub → Actions**

Expected workflow sequence:

1. **Tests** (test.yml)
   - [ ] Unit tests pass
   - [ ] Integration tests pass
   - [ ] Linting passes
   - [ ] Security scan passes
   - **Status:** ✅ Success

2. **Build Docker Image** (build.yml)
   - [ ] Triggered after tests pass
   - [ ] Docker image built
   - [ ] Tags pushed to GHCR:
     - `latest`
     - `v1.0.X`
     - `sha-abc123`
   - **Status:** ✅ Success

3. **Deploy to Production** (deploy.yml)
   - [ ] Triggered after build succeeds
   - [ ] SSH connection successful
   - [ ] Docker image pulled
   - [ ] Containers started
   - [ ] Health check passes
   - **Status:** ✅ Success

### Step 3: Verify Deployment on VPS

```bash
# SSH to VPS
ssh user@vps-host
cd /opt/suproxy

# Check containers are running
docker-compose -f docker-compose.production.yml ps

# Expected output: All containers "Up"
# - suproxy-api-prod
# - suproxy-postgres-prod
# - suproxy-prometheus
# - suproxy-grafana

# Check health endpoint
curl http://localhost:8080/health

# Expected: {"status":"ok","timestamp":"..."}

# Check logs (no errors)
docker-compose -f docker-compose.production.yml logs --tail=50 api
```

### Step 4: Access Services

- [ ] API: `http://<VPS-IP>:8080/health`
- [ ] Prometheus: `http://<VPS-IP>:9090`
- [ ] Grafana: `http://<VPS-IP>:3000` (admin / password from .env.production)

## Manual Deployment Test

### Deploy Specific Version

1. Go to: **Actions → Deploy to Production → Run workflow**
2. Inputs:
   - Version: `latest` (or specific version like `v1.0.42`)
   - Servers: `finland` (or `all`)
3. Click **Run workflow**
4. Monitor deployment:
   - [ ] Prepare step succeeds
   - [ ] Deploy to finland succeeds
   - [ ] Health check passes
   - [ ] Deployment summary shows success

## Release Creation Test

1. Go to: **Actions → Create Release → Run workflow**
2. Inputs:
   - Version: `v1.0.0` (increment as needed)
   - Prerelease: `false`
   - Draft: `false`
3. Click **Run workflow**
4. Verify:
   - [ ] Git tag created
   - [ ] GitHub release created: **Releases** tab
   - [ ] Changelog generated
   - [ ] Deployment triggered automatically
   - [ ] New version deployed to production

## Security Scan Test

1. Go to: **Actions → Security Scan → Run workflow**
2. Click **Run workflow**
3. Verify:
   - [ ] Dependency scan passes
   - [ ] Code scan passes
   - [ ] Docker scan passes
   - [ ] Secret scan passes
   - [ ] Results visible in **Security** tab

## Health Check Test

1. Go to: **Actions → Health Check → Run workflow**
2. Inputs:
   - Servers: `all`
3. Click **Run workflow**
4. Verify:
   - [ ] All servers report healthy
   - [ ] No issues created
   - [ ] Summary shows success

## Rollback Test

### Deploy New Version

```bash
Actions → Deploy to Production
  Version: latest
  Servers: finland
```

### Rollback to Previous

```bash
Actions → Deploy to Production
  Version: v1.0.99  # Previous version
  Servers: finland
```

### Verify Rollback

```bash
# SSH to VPS
ssh user@vps-host

# Check running version
docker images | grep suproxy-backend

# Should show v1.0.99
```

## Automated Workflows Verification

### Daily Security Scan

- [ ] Scheduled to run daily at 3 AM
- [ ] Check next scheduled run: **Actions → Security Scan**
- [ ] Review past runs for any failures

### Periodic Health Check

- [ ] Scheduled to run every 15 minutes
- [ ] Check next scheduled run: **Actions → Health Check**
- [ ] Verify no open health-check-failure issues

## Documentation Verification

Ensure all documentation is accessible and accurate:

- [ ] `README.md` - Project overview
- [ ] `IMPLEMENTATION_SUMMARY.md` - Implementation details
- [ ] `docs/CI_CD_ARCHITECTURE.md` - Workflow architecture
- [ ] `docs/DEPLOYMENT.md` - Deployment guide
- [ ] `docs/ROLLBACK.md` - Rollback procedures
- [ ] `docs/MULTISERVER.md` - Multi-server setup
- [ ] `docs/BACKUP.md` - Backup strategies
- [ ] `docs/SERVER_SETUP.md` - VPS setup guide

## Backup Verification

### Manual Backup Test

```bash
# SSH to VPS
ssh user@vps-host
cd /opt/suproxy

# Run backup
./scripts/backup.sh

# Verify backup created
ls -lh backups/

# Should see: postgres_backup_YYYYMMDD_HHMMSS.sql.gz
```

### Automated Backup (Cron)

```bash
# Check crontab
crontab -l

# Should see:
# 0 2 * * * /opt/suproxy/scripts/backup.sh >> /var/log/suproxy-backup.log 2>&1

# Check backup log
cat /var/log/suproxy-backup.log
```

## Monitoring Setup

### Prometheus

1. Access: `http://<VPS-IP>:9090`
2. Verify:
   - [ ] UI loads
   - [ ] Targets are up: Status → Targets
   - [ ] Metrics collecting: Graph → Query

### Grafana

1. Access: `http://<VPS-IP>:3000`
2. Login: admin / (from .env.production)
3. Verify:
   - [ ] UI loads
   - [ ] Prometheus data source configured
   - [ ] Can query metrics

## Troubleshooting Common Issues

### Build Fails

- **Check:** GitHub Actions logs
- **Fix:** Review test failures, fix code, push again

### Deployment Fails - SSH Error

- **Check:** GitHub Secrets are correct
- **Fix:** Verify VPS_*_HOST, _USER, _KEY, _PORT

### Deployment Fails - Docker Pull Error

- **Check:** Image exists in GHCR
- **Fix:** Verify build.yml completed successfully

### Health Check Fails

- **Check:** Container logs on VPS
```bash
docker-compose -f docker-compose.production.yml logs api
```
- **Fix:** Review logs, fix issues, redeploy

### Containers Not Starting

- **Check:** `.env.production` configuration
- **Check:** Database connection
- **Fix:** Verify all environment variables

## Success Criteria

All checkboxes should be ✅ before considering deployment complete:

### GitHub Actions
- [ ] All workflows visible in Actions tab
- [ ] No failing workflows
- [ ] Artifacts accessible
- [ ] Secrets configured

### VPS
- [ ] All containers running
- [ ] Health endpoint returns 200
- [ ] No errors in logs
- [ ] Services accessible

### Workflows
- [ ] Automatic deployment works
- [ ] Manual deployment works
- [ ] Release creation works
- [ ] Rollback tested and works
- [ ] Security scans running
- [ ] Health checks running

### Documentation
- [ ] All docs created and accurate
- [ ] Team can follow procedures
- [ ] Troubleshooting guides helpful

## Next Actions After Verification

1. **Add More Servers**
   - Setup new VPS (see `docs/SERVER_SETUP.md`)
   - Add GitHub Secrets
   - Deploy to new server

2. **Configure Monitoring Alerts**
   - Setup Prometheus alerts
   - Configure notification channels
   - Test alerting

3. **Implement Automated Backups to Cloud**
   - Configure S3/GCS
   - Upload backups offsite
   - Test recovery from cloud backups

4. **Setup Blue/Green or Canary Deployment**
   - Plan deployment strategy
   - Implement traffic switching
   - Test deployment scenarios

## Emergency Contacts

Document who to contact for issues:

- **Infrastructure:** _________________
- **Application:** _________________
- **Database:** _________________
- **Security:** _________________

## Sign-Off

- [ ] All tests passed
- [ ] Production deployment successful
- [ ] Rollback tested
- [ ] Documentation complete
- [ ] Team trained

**Deployed by:** _________________ **Date:** _________________

**Verified by:** _________________ **Date:** _________________

---

**🎉 Your enterprise-grade CI/CD is now live!**
