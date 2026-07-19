# 🌍 Multi-Server Deployment Guide

## Overview

SuProxy Backend supports deployment to multiple servers across different regions for high availability, performance, and redundancy.

## Architecture

### Current Servers

- **Finland** 🇫🇮 - Primary production server (active)
- **Germany** 🇩🇪 - Future
- **Turkey** 🇹🇷 - Future

### Future Expansion

- **USA** 🇺🇸 - North America coverage
- **Japan** 🇯🇵 - Asia coverage
- **Singapore** 🇸🇬 - Southeast Asia coverage

## Adding a New Server

### Step 1: Server Setup

Set up VPS with required software:

```bash
# Install Docker
curl -fsSL https://get.docker.com | sh

# Install Docker Compose
sudo apt-get update
sudo apt-get install docker-compose-plugin

# Create directory structure
sudo mkdir -p /opt/suproxy/scripts
sudo mkdir -p /opt/suproxy/backups
sudo mkdir -p /opt/suproxy/prometheus
sudo mkdir -p /opt/suproxy/grafana/provisioning

# Set permissions
sudo chown -R $USER:$USER /opt/suproxy
```

### Step 2: Copy Configuration Files

Transfer files to new server:

```bash
# From local machine
scp docker-compose.production.yml user@new-server:/opt/suproxy/
scp .env.production user@new-server:/opt/suproxy/
scp scripts/deploy.sh user@new-server:/opt/suproxy/scripts/
scp -r prometheus/ user@new-server:/opt/suproxy/
scp -r grafana/ user@new-server:/opt/suproxy/

# On new server, make script executable
ssh user@new-server
chmod +x /opt/suproxy/scripts/deploy.sh
```

### Step 3: Configure Environment

Edit `.env.production` on new server:

```bash
ssh user@new-server
cd /opt/suproxy
nano .env.production

# Update server-specific values:
# - Database credentials
# - JWT secrets
# - Grafana passwords
# - Port configurations (if needed)
```

### Step 4: Add GitHub Secrets

Add secrets for the new server in GitHub repository settings:

Navigate to: **Settings** → **Secrets and variables** → **Actions** → **New repository secret**

```
VPS_<COUNTRY>_HOST     = IP address or hostname
VPS_<COUNTRY>_USER     = SSH username
VPS_<COUNTRY>_KEY      = SSH private key
VPS_<COUNTRY>_PORT     = SSH port (usually 22)
```

**Example for Germany:**
```
VPS_GERMANY_HOST     = 192.168.1.100
VPS_GERMANY_USER     = deploy
VPS_GERMANY_KEY      = -----BEGIN OPENSSH PRIVATE KEY-----...
VPS_GERMANY_PORT     = 22
```

### Step 5: Update Deploy Workflow

The workflow is already prepared for new servers! Just add the server mapping:

Edit `.github/workflows/deploy.yml`:

```yaml
- name: Get Server Configuration
  id: config
  run: |
    SERVER="${{ matrix.server }}"
    
    case "$SERVER" in
      # ... existing servers ...
      newcountry)
        echo "host_secret=VPS_NEWCOUNTRY_HOST" >> $GITHUB_OUTPUT
        echo "user_secret=VPS_NEWCOUNTRY_USER" >> $GITHUB_OUTPUT
        echo "key_secret=VPS_NEWCOUNTRY_KEY" >> $GITHUB_OUTPUT
        echo "port_secret=VPS_NEWCOUNTRY_PORT" >> $GITHUB_OUTPUT
        ;;
    esac
```

### Step 6: Test Deployment

Test deployment to new server:

```bash
# Via GitHub UI:
Actions → Deploy to Production → Run workflow
- Version: latest
- Servers: newcountry

# Verify deployment
ssh user@new-server
cd /opt/suproxy
docker-compose -f docker-compose.production.yml ps
curl http://localhost:8080/health
```

## Deployment Strategies

### Sequential Deployment (Current)

Deploys to one server at a time:

```yaml
strategy:
  max-parallel: 1
  fail-fast: false
```

**Advantages:**
- Safer - issues detected before affecting all servers
- Allows monitoring between deployments
- Can stop if first deployment fails

**Disadvantages:**
- Slower - each server waits for previous
- Full deployment takes longer

### Parallel Deployment (Future)

Deploy to multiple servers simultaneously:

```yaml
strategy:
  max-parallel: 3  # Deploy to 3 servers at once
  fail-fast: true  # Stop if any fails
```

**Advantages:**
- Faster - reduced deployment time
- Good for urgent hotfixes

**Disadvantages:**
- Higher risk - issues affect multiple servers
- Harder to monitor

### Gradual Rollout (Future)

Deploy to increasing percentages:

```yaml
# Phase 1: Canary (10%)
Servers: finland

# Phase 2: Partial (50%)
Servers: finland,germany

# Phase 3: Full (100%)
Servers: all
```

## Server Groups

### By Region

```bash
# Europe
Servers: finland,germany

# Asia
Servers: japan,singapore

# All regions
Servers: all
```

### By Environment

```bash
# Production
Servers: finland,germany,turkey

# Staging (if configured)
Servers: staging-finland

# Development
Servers: dev-server
```

### By Priority

```bash
# Critical infrastructure
Servers: finland,germany  # Primary markets

# Secondary expansion
Servers: turkey,usa  # Growing markets

# Experimental
Servers: japan,singapore  # Testing regions
```

## Load Balancing

### DNS-based (Recommended for Future)

Use GeoDNS to route users to nearest server:

```
api.suproxy.com
  → Finland (EU users)
  → Germany (EU users)
  → Turkey (MENA users)
  → USA (Americas users)
  → Japan (Asia users)
```

### Reverse Proxy

Use Nginx/HAProxy to distribute load:

```nginx
upstream suproxy_backend {
  server finland.suproxy.com:8080 weight=3;
  server germany.suproxy.com:8080 weight=2;
  server turkey.suproxy.com:8080 weight=1;
}
```

### CDN Integration

Use Cloudflare/AWS CloudFront for:
- Geographic load balancing
- DDoS protection
- SSL termination
- Caching

## Monitoring Multi-Server Setup

### Centralized Monitoring (Future)

**Option 1: Prometheus Federation**

Central Prometheus scrapes all server Prometheus instances:

```yaml
# Central Prometheus config
scrape_configs:
  - job_name: 'finland'
    static_configs:
      - targets: ['finland:9090']
  - job_name: 'germany'
    static_configs:
      - targets: ['germany:9090']
```

**Option 2: Grafana Multi-Source**

Single Grafana instance with multiple data sources:

```
Grafana Dashboard
  ├── Data Source: Finland Prometheus
  ├── Data Source: Germany Prometheus
  └── Data Source: Turkey Prometheus
```

**Option 3: Observability Platform**

Use SaaS monitoring:
- Datadog
- New Relic
- Grafana Cloud

### Health Check Dashboard

Monitor all servers from single location:

```bash
# Create monitoring script
#!/bin/bash

SERVERS=("finland" "germany" "turkey")

for server in "${SERVERS[@]}"; do
  echo "Checking $server..."
  curl -s http://$server.suproxy.com:8080/health || echo "❌ $server is down"
done
```

## Database Strategy

### Option 1: Database Per Server (Current)

Each server has its own PostgreSQL:

**Advantages:**
- Simple setup
- No cross-region latency
- Independent failures

**Disadvantages:**
- Data inconsistency between servers
- No shared user accounts
- Requires data synchronization

### Option 2: Primary-Replica Setup

One primary database, multiple replicas:

```
Finland (Primary - Read/Write)
   ↓ Replication
Germany (Replica - Read Only)
   ↓ Replication
Turkey (Replica - Read Only)
```

### Option 3: Distributed Database

Use PostgreSQL with Citus or CockroachDB:

**Advantages:**
- Automatic replication
- Geographic distribution
- Strong consistency

**Disadvantages:**
- Complex setup
- Higher cost
- Requires expertise

## Data Synchronization

If using per-server databases, sync data:

### User Data Sync

```bash
# Replicate user table from primary to secondaries
pg_dump -h finland -U user -t users | psql -h germany -U user
```

### Configuration Sync

Store shared config in:
- Redis (centralized)
- S3/Object Storage
- etcd/Consul

## Failure Handling

### Single Server Failure

```yaml
strategy:
  fail-fast: false  # Continue deploying to other servers
```

If one server fails:
- ✅ Other servers continue deploying
- ⚠️  Failed server marked in summary
- 🔧 Manual intervention required

### Multiple Server Failure

If more than 50% fail:
- ❌ Deployment marked as failed
- 🚨 Alert team immediately
- 🔄 Consider rollback

### Network Partition

If SSH fails:
- Retry 3 times with backoff
- If still fails, skip server
- Log failure for manual review

## Cost Optimization

### Right-sizing Servers

```bash
# Small regions (testing)
1 vCPU, 2GB RAM - $10/month

# Medium regions (growing)
2 vCPU, 4GB RAM - $20/month

# Large regions (primary)
4 vCPU, 8GB RAM - $40/month
```

### Scaling Strategy

```bash
# Start small
Phase 1: 1 server (Finland)

# Add redundancy
Phase 2: 2 servers (Finland + Germany)

# Expand regions
Phase 3: 3-5 servers (EU + MENA)

# Global coverage
Phase 4: 6+ servers (Worldwide)
```

## Security

### Firewall Rules

```bash
# On each server
sudo ufw allow 22/tcp    # SSH
sudo ufw allow 80/tcp    # HTTP (reverse proxy)
sudo ufw allow 443/tcp   # HTTPS (reverse proxy)
sudo ufw allow 8080/tcp  # API (internal only)
sudo ufw allow 9090/tcp  # Prometheus (internal only)
sudo ufw allow 3000/tcp  # Grafana (internal only)
sudo ufw enable
```

### VPN/Private Network

Connect servers via VPN:
- Tailscale
- WireGuard
- Cloud provider VPC

### SSH Hardening

```bash
# Disable password authentication
PasswordAuthentication no

# Use SSH keys only
PubkeyAuthentication yes

# Limit user access
AllowUsers deploy
```

## Testing Multi-Server Deployment

### Dry Run

Test deployment workflow without affecting production:

```bash
# Test with staging servers
Servers: staging-finland,staging-germany
```

### Canary Server

Use one server as canary:

```bash
# Deploy to canary first
Servers: finland

# Wait and monitor (30 minutes)

# If successful, deploy to others
Servers: germany,turkey
```

### Rollback Testing

Periodically test rollback across servers:

```bash
# Deploy version N+1
Servers: all
Version: v1.0.42

# Rollback to N
Servers: all
Version: v1.0.41
```

## Related Documentation

- [DEPLOYMENT.md](./DEPLOYMENT.md) - Deployment procedures
- [ROLLBACK.md](./ROLLBACK.md) - Rollback strategies
- [CI_CD_ARCHITECTURE.md](./CI_CD_ARCHITECTURE.md) - Workflow details
- [SERVER_SETUP.md](./SERVER_SETUP.md) - Server configuration (future)
