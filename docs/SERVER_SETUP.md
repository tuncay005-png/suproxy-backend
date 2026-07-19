# 🖥️ Server Setup Guide

## Overview

This guide explains how to set up a new VPS for running SuProxy Backend in production.

## Prerequisites

- Ubuntu 22.04 LTS or newer
- Minimum 2GB RAM, 2 vCPU
- 20GB disk space
- Root or sudo access
- Public IP address
- SSH access

## Initial Server Setup

### 1. Update System

```bash
# Update package list
sudo apt-get update

# Upgrade packages
sudo apt-get upgrade -y

# Install essential tools
sudo apt-get install -y \
  curl \
  wget \
  git \
  vim \
  htop \
  net-tools \
  ufw
```

### 2. Configure Firewall

```bash
# Allow SSH
sudo ufw allow 22/tcp

# Allow HTTP/HTTPS (for reverse proxy)
sudo ufw allow 80/tcp
sudo ufw allow 443/tcp

# Allow application ports (adjust as needed)
sudo ufw allow 8080/tcp  # API
sudo ufw allow 9090/tcp  # Prometheus
sudo ufw allow 3000/tcp  # Grafana

# Enable firewall
sudo ufw --force enable

# Check status
sudo ufw status
```

### 3. Create Deployment User

```bash
# Create user
sudo adduser deploy

# Add to docker group (will be created later)
sudo usermod -aG docker deploy

# Add sudo privileges (optional)
sudo usermod -aG sudo deploy

# Switch to deploy user
su - deploy
```

## Install Docker

### Docker Engine

```bash
# Install Docker
curl -fsSL https://get.docker.com | sh

# Start Docker service
sudo systemctl start docker
sudo systemctl enable docker

# Verify installation
docker --version
```

### Docker Compose

```bash
# Docker Compose is included in Docker
docker compose version

# If not available, install manually
sudo apt-get install docker-compose-plugin
```

## Install Additional Tools

### PostgreSQL Client (for debugging)

```bash
sudo apt-get install -y postgresql-client
```

### AWS CLI (for S3 backups)

```bash
curl "https://awscli.amazonaws.com/awscli-exe-linux-x86_64.zip" -o "awscliv2.zip"
unzip awscliv2.zip
sudo ./aws/install
```

## Setup Project Directory

### Create Directory Structure

```bash
# Create main directory
sudo mkdir -p /opt/suproxy
sudo chown -R deploy:deploy /opt/suproxy
cd /opt/suproxy

# Create subdirectories
mkdir -p scripts
mkdir -p backups
mkdir -p prometheus
mkdir -p grafana/provisioning
mkdir -p logs
```

### Set Permissions

```bash
chmod 755 /opt/suproxy
chmod 755 /opt/suproxy/scripts
chmod 700 /opt/suproxy/backups  # Restrict backup access
```

## Configure SSH

### Generate SSH Key for GitHub

```bash
# On your local machine
ssh-keygen -t ed25519 -C "deploy@suproxy"

# Copy public key
cat ~/.ssh/id_ed25519.pub

# Add to GitHub: Settings → SSH and GPG keys → New SSH key
```

### Configure SSH for VPS

```bash
# On VPS
mkdir -p ~/.ssh
chmod 700 ~/.ssh

# Add your public key
nano ~/.ssh/authorized_keys
# Paste your public key here

chmod 600 ~/.ssh/authorized_keys
```

### Harden SSH

```bash
# Edit SSH config
sudo nano /etc/ssh/sshd_config

# Recommended settings:
PermitRootLogin no
PasswordAuthentication no
PubkeyAuthentication yes
```



ChallengeResponseAuthentication no
UsePAM yes
X11Forwarding no
AllowUsers deploy

# Restart SSH
sudo systemctl restart sshd
```

## Clone Repository

```bash
cd /opt/suproxy

# If repository is private, set up SSH key first
git clone git@github.com:tuncay005-png/suproxy-backend.git .

# Or use HTTPS with token
git clone https://github.com/tuncay005-png/suproxy-backend.git .
```

## Configure Environment

### Create Production Environment File

```bash
cd /opt/suproxy
cp .env.example .env.production

# Edit with actual values
nano .env.production
```

### Set Secure Values

```bash
# Generate secure JWT secret
openssl rand -base64 32

# Generate secure database password
openssl rand -base64 24

# Generate secure Grafana password
openssl rand -base64 16

# Update .env.production with these values
```

### Protect Environment File

```bash
chmod 600 .env.production
```

## Configure Docker Registry Authentication

### Login to GHCR

```bash
# Create GitHub personal access token
# GitHub → Settings → Developer settings → Personal access tokens
# Scopes: read:packages

# Save token
export GITHUB_TOKEN="ghp_your_token_here"

# Login to GHCR
echo $GITHUB_TOKEN | docker login ghcr.io -u tuncay005-png --password-stdin

# Verify
docker pull ghcr.io/tuncay005-png/suproxy-backend:latest
```

## Setup Monitoring

### Prometheus Configuration

```bash
cd /opt/suproxy

# Create Prometheus config
nano prometheus/prometheus.yml
```

```yaml
global:
  scrape_interval: 15s
  evaluation_interval: 15s

scrape_configs:
  - job_name: 'suproxy-api'
    static_configs:
      - targets: ['api:8080']

  - job_name: 'prometheus'
    static_configs:
      - targets: ['localhost:9090']
```

### Grafana Provisioning

```bash
# Create datasource config
mkdir -p grafana/provisioning/datasources
nano grafana/provisioning/datasources/prometheus.yml
```

```yaml
apiVersion: 1

datasources:
  - name: Prometheus
    type: prometheus
    access: proxy
    url: http://prometheus:9090
    isDefault: true
```

## Setup Automated Backups

### Create Backup Script

```bash
nano /opt/suproxy/scripts/backup.sh
```

```bash
#!/bin/bash

BACKUP_DIR="/opt/suproxy/backups"
DATE=$(date +%Y%m%d_%H%M%S)
BACKUP_FILE="$BACKUP_DIR/postgres_backup_$DATE.sql.gz"

# Load environment
source /opt/suproxy/.env.production

# Create backup
docker-compose -f /opt/suproxy/docker-compose.production.yml exec -T postgres \
  pg_dump -U $DB_USER $DB_NAME | gzip > "$BACKUP_FILE"

# Keep only last 7 days
find $BACKUP_DIR -name "postgres_backup_*.sql.gz" -mtime +7 -delete

echo "✅ Backup completed: $BACKUP_FILE"
```

```bash
chmod +x /opt/suproxy/scripts/backup.sh
```

### Schedule Backups with Cron

```bash
# Edit crontab
crontab -e

# Add daily backup at 2 AM
0 2 * * * /opt/suproxy/scripts/backup.sh >> /var/log/suproxy-backup.log 2>&1
```

## Initial Deployment

### Pull and Start Services

```bash
cd /opt/suproxy

# Pull latest image
docker pull ghcr.io/tuncay005-png/suproxy-backend:latest

# Start services
docker-compose -f docker-compose.production.yml up -d

# Check status
docker-compose -f docker-compose.production.yml ps

# View logs
docker-compose -f docker-compose.production.yml logs -f
```

### Verify Deployment

```bash
# Health check
curl http://localhost:8080/health

# Expected response
{"status":"ok","timestamp":"2026-07-19T..."}

# Check all containers
docker ps

# Check resource usage
docker stats
```

## Configure System Services

### Enable Docker on Boot

```bash
sudo systemctl enable docker
```

### Auto-restart Containers

Already configured in docker-compose.production.yml:

```yaml
services:
  api:
    restart: always
  postgres:
    restart: always
```

## Setup Logging

### Configure Log Rotation

```bash
sudo nano /etc/logrotate.d/suproxy
```

```
/opt/suproxy/logs/*.log {
    daily
    rotate 7
    compress
    delaycompress
    missingok
    notifempty
    create 0640 deploy deploy
}
```

### Docker Log Limits

Already configured in docker-compose.production.yml:

```yaml
logging:
  driver: "json-file"
  options:
    max-size: "10m"
    max-file: "5"
```

## Security Hardening

### Fail2Ban

```bash
# Install
sudo apt-get install -y fail2ban

# Configure
sudo nano /etc/fail2ban/jail.local
```

```ini
[sshd]
enabled = true
port = 22
maxretry = 3
bantime = 3600
```

```bash
# Restart
sudo systemctl restart fail2ban
```

### Automatic Updates

```bash
# Install unattended-upgrades
sudo apt-get install -y unattended-upgrades

# Configure
sudo dpkg-reconfigure -plow unattended-upgrades
```

### System Monitoring

```bash
# Install monitoring tools
sudo apt-get install -y sysstat

# Enable
sudo systemctl enable sysstat
sudo systemctl start sysstat
```

## Testing

### Test Deploy Script

```bash
cd /opt/suproxy
./scripts/deploy.sh
```

### Test Backup Script

```bash
./scripts/backup.sh

# Verify backup created
ls -lh backups/
```

### Test Health Checks

```bash
# API health
curl http://localhost:8080/health

# Prometheus
curl http://localhost:9090/-/healthy

# Grafana
curl http://localhost:3000/api/health
```

## Monitoring Setup

### Access Grafana

```bash
# Get IP address
ip addr show

# Access Grafana
http://<server-ip>:3000

# Default credentials
Username: admin
Password: (from .env.production GRAFANA_PASSWORD)
```

### Import Dashboards

1. Login to Grafana
2. Go to Dashboards → Import
3. Use dashboard IDs:
   - 1860 (Node Exporter Full)
   - 3662 (Prometheus 2.0 Stats)

## Troubleshooting

### Container Won't Start

```bash
# Check logs
docker-compose -f docker-compose.production.yml logs api

# Check environment
docker-compose -f docker-compose.production.yml config

# Restart service
docker-compose -f docker-compose.production.yml restart api
```

### Database Connection Issues

```bash
# Check PostgreSQL
docker-compose -f docker-compose.production.yml logs postgres

# Connect to database
docker-compose -f docker-compose.production.yml exec postgres \
  psql -U $DB_USER $DB_NAME
```

### Disk Space Issues

```bash
# Check disk usage
df -h

# Clean Docker images
docker system prune -a

# Clean old backups
find /opt/suproxy/backups -mtime +7 -delete
```

## Maintenance

### Update Application

```bash
# Pull latest
docker pull ghcr.io/tuncay005-png/suproxy-backend:latest

# Restart
docker-compose -f docker-compose.production.yml up -d
```

### View Logs

```bash
# All logs
docker-compose -f docker-compose.production.yml logs

# Specific service
docker-compose -f docker-compose.production.yml logs api

# Follow logs
docker-compose -f docker-compose.production.yml logs -f
```

### Check Resource Usage

```bash
# System resources
htop

# Docker resources
docker stats

# Disk usage
du -sh /opt/suproxy/*
```

## Related Documentation

- [DEPLOYMENT.md](./DEPLOYMENT.md) - Deployment procedures
- [MULTISERVER.md](./MULTISERVER.md) - Multi-server setup
- [BACKUP.md](./BACKUP.md) - Backup strategies
- [CI_CD_ARCHITECTURE.md](./CI_CD_ARCHITECTURE.md) - Workflow details
