# 💾 Backup and Recovery Guide

## Overview

This document describes backup strategies, automated backup workflows, and recovery procedures for SuProxy Backend.

## Backup Strategy

### What to Backup

1. **PostgreSQL Database** - Critical user and configuration data
2. **Configuration Files** - .env.production, docker-compose.yml
3. **Application Logs** - For audit and debugging
4. **Docker Volumes** - Persistent data (optional)

### Backup Frequency

- **Database**: Daily automated backups
- **Configuration**: On every deployment
- **Logs**: Continuous (rotate weekly)
- **Full System**: Weekly snapshots

### Retention Policy

- **Daily backups**: Keep for 7 days
- **Weekly backups**: Keep for 4 weeks
- **Monthly backups**: Keep for 12 months
- **Critical backups**: Keep indefinitely

## Automated Database Backup

### Backup Script

The project includes `scripts/backup.sh`:

```bash
#!/bin/bash
# Location: /opt/suproxy/scripts/backup.sh

BACKUP_DIR="/opt/suproxy/backups"
DATE=$(date +%Y%m%d_%H%M%S)
BACKUP_FILE="$BACKUP_DIR/postgres_backup_$DATE.sql.gz"

# Create backup
docker-compose -f /opt/suproxy/docker-compose.production.yml exec -T postgres \
  pg_dump -U $DB_USER $DB_NAME | gzip > "$BACKUP_FILE"

# Keep only last 7 days
find $BACKUP_DIR -name "postgres_backup_*.sql.gz" -mtime +7 -delete

echo "✅ Backup completed: $BACKUP_FILE"
```

### Automated Schedule

Set up cron job on VPS:

```bash
# Edit crontab
crontab -e

# Add daily backup at 2 AM
0 2 * * * /opt/suproxy/scripts/backup.sh >> /var/log/backup.log 2>&1

# Add weekly full backup at 3 AM Sunday
0 3 * * 0 /opt/suproxy/scripts/backup_full.sh >> /var/log/backup.log 2>&1
```

## Pre-Deployment Backup

### Automatic Backup Workflow (Future)

Create `.github/workflows/backup.yml`:

```yaml
name: Pre-Deployment Backup

on:
  workflow_call:
    inputs:
      server:
        required: true
        type: string
```


### Manual Pre-Deployment Backup

Before critical deployments:

```bash
# SSH to server
ssh user@vps-host

# Run backup script
cd /opt/suproxy
./scripts/backup.sh

# Verify backup created
ls -lh backups/ | tail -1
```

## Database Recovery

### Restore from Backup

```bash
# SSH to server
ssh user@vps-host
cd /opt/suproxy

# List available backups
ls -lh backups/

# Stop application
docker-compose -f docker-compose.production.yml stop api

# Restore database
gunzip -c backups/postgres_backup_20260719_020000.sql.gz | \
  docker-compose -f docker-compose.production.yml exec -T postgres \
  psql -U $DB_USER $DB_NAME

# Restart application
docker-compose -f docker-compose.production.yml start api

# Verify health
curl http://localhost:8080/health
```

### Point-in-Time Recovery

PostgreSQL supports WAL-based recovery:

```bash
# Configure in docker-compose.production.yml
services:
  postgres:
    command:
      - postgres
      - -c
      - wal_level=replica
      - -c
      - archive_mode=on
      - -c
      - archive_command='cp %p /backups/wal/%f'
```

## Backup Storage

### Local Storage (Current)

Backups stored on VPS:

```
/opt/suproxy/backups/
├── postgres_backup_20260719_020000.sql.gz
├── postgres_backup_20260718_020000.sql.gz
└── postgres_backup_20260717_020000.sql.gz
```

### Remote Storage (Recommended for Future)

#### Option 1: AWS S3

```bash
# Install AWS CLI
apt-get install awscli

# Configure
aws configure

# Upload backup
aws s3 cp backups/postgres_backup_*.sql.gz \
  s3://suproxy-backups/$(hostname)/
```



#### Option 2: Google Cloud Storage

```bash
# Install gsutil
curl https://sdk.cloud.google.com | bash

# Authenticate
gcloud auth login

# Upload backup
gsutil cp backups/postgres_backup_*.sql.gz \
  gs://suproxy-backups/$(hostname)/
```

#### Option 3: Backblaze B2

```bash
# Install B2 CLI
pip install b2

# Authenticate
b2 authorize-account <keyID> <applicationKey>

# Upload backup
b2 upload-file suproxy-backups \
  backups/postgres_backup_*.sql.gz \
  $(hostname)/postgres_backup_*.sql.gz
```

## Configuration Backup

### Environment Files

```bash
# Backup before changes
cp .env.production .env.production.backup

# Restore if needed
cp .env.production.backup .env.production
```

### Docker Compose

```bash
# Backup docker-compose
cp docker-compose.production.yml docker-compose.production.yml.backup

# Restore
cp docker-compose.production.yml.backup docker-compose.production.yml
```

## Application Logs

### Log Collection

```bash
# Export recent logs
docker-compose -f docker-compose.production.yml logs --since 24h api > logs/api_$(date +%Y%m%d).log

# Compress old logs
gzip logs/api_*.log
```

### Log Rotation

Configure in docker-compose:

```yaml
services:
  api:
    logging:
      driver: "json-file"
      options:
        max-size: "10m"
        max-file: "5"
```



## Disaster Recovery

### Full System Backup

```bash
#!/bin/bash
# Full backup including database, configs, and logs

BACKUP_DATE=$(date +%Y%m%d)
BACKUP_DIR="/opt/suproxy/backups/full_$BACKUP_DATE"

mkdir -p $BACKUP_DIR

# Database
docker-compose -f /opt/suproxy/docker-compose.production.yml exec -T postgres \
  pg_dump -U $DB_USER $DB_NAME | gzip > "$BACKUP_DIR/database.sql.gz"

# Configurations
cp /opt/suproxy/.env.production $BACKUP_DIR/
cp /opt/suproxy/docker-compose.production.yml $BACKUP_DIR/

# Application logs (last 7 days)
docker-compose -f /opt/suproxy/docker-compose.production.yml logs --since 168h \
  > $BACKUP_DIR/logs.txt

# Compress everything
tar -czf /opt/suproxy/backups/full_backup_$BACKUP_DATE.tar.gz -C /opt/suproxy/backups full_$BACKUP_DATE
rm -rf $BACKUP_DIR

echo "✅ Full backup: full_backup_$BACKUP_DATE.tar.gz"
```

### System Recovery

```bash
# On new server
# 1. Install Docker and Docker Compose

# 2. Extract backup
tar -xzf full_backup_20260719.tar.gz -C /opt/suproxy/

# 3. Restore configurations
cp full_20260719/.env.production /opt/suproxy/
cp full_20260719/docker-compose.production.yml /opt/suproxy/

# 4. Start services
docker-compose -f /opt/suproxy/docker-compose.production.yml up -d

# 5. Restore database
gunzip -c full_20260719/database.sql.gz | \
  docker-compose -f docker-compose.production.yml exec -T postgres \
  psql -U $DB_USER $DB_NAME

# 6. Verify
curl http://localhost:8080/health
```

## Backup Monitoring

### Verify Backups

```bash
#!/bin/bash
# Check backup integrity

LATEST_BACKUP=$(ls -t /opt/suproxy/backups/postgres_backup_*.sql.gz | head -1)

if [ -z "$LATEST_BACKUP" ]; then
  echo "❌ No backups found"
  exit 1
fi

# Check age
BACKUP_AGE=$(find "$LATEST_BACKUP" -mtime +1)
if [ -n "$BACKUP_AGE" ]; then
  echo "⚠️  Latest backup is older than 24 hours"
fi

# Check size
BACKUP_SIZE=$(stat -f%z "$LATEST_BACKUP" 2>/dev/null || stat -c%s "$LATEST_BACKUP")
if [ "$BACKUP_SIZE" -lt 1000 ]; then
  echo "❌ Backup file is too small"
  exit 1
fi

echo "✅ Backup verification passed"
```



### Backup Alerts

Set up monitoring:

```bash
# Add to crontab - verify backups daily
0 3 * * * /opt/suproxy/scripts/verify_backup.sh || \
  echo "Backup verification failed" | mail -s "Backup Alert" admin@example.com
```

## Multi-Server Backup

### Centralized Backup Storage

For multiple servers:

```bash
# Each server uploads to central location
aws s3 sync /opt/suproxy/backups/ \
  s3://suproxy-backups/$(hostname)/ \
  --delete
```

### Cross-Server Replication

Replicate backups between servers:

```bash
# From Finland to Germany
rsync -avz -e ssh /opt/suproxy/backups/ \
  user@germany-server:/opt/suproxy/backups-finland/
```

## Backup Testing

### Regular Restore Tests

Test recovery quarterly:

```bash
# 1. Create test environment
docker run -d --name test-postgres postgres:15-alpine

# 2. Restore backup to test environment
gunzip -c backups/postgres_backup_latest.sql.gz | \
  docker exec -i test-postgres psql -U postgres

# 3. Verify data
docker exec test-postgres psql -U postgres -c "SELECT COUNT(*) FROM users;"

# 4. Cleanup
docker rm -f test-postgres
```

## Best Practices

### Before Deployment

1. ✅ Create backup
2. ✅ Verify backup integrity
3. ✅ Test backup restoration (periodically)
4. ✅ Document recovery procedure

### Backup Security

1. ✅ Encrypt backups at rest
2. ✅ Encrypt backups in transit
3. ✅ Restrict backup access
4. ✅ Audit backup access logs

### Backup Verification

1. ✅ Automate verification
2. ✅ Test restoration regularly
3. ✅ Monitor backup age
4. ✅ Check backup size

## Related Documentation

- [DEPLOYMENT.md](./DEPLOYMENT.md) - Deployment procedures
- [ROLLBACK.md](./ROLLBACK.md) - Rollback strategies
- [MULTISERVER.md](./MULTISERVER.md) - Multi-server setup
- [CI_CD_ARCHITECTURE.md](./CI_CD_ARCHITECTURE.md) - Workflow details
