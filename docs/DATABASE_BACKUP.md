# Database Backup and Restore Guide

**Last Updated:** 2024-08-21  
**Version:** 1.0.0

---

## Overview

This guide covers the automated PostgreSQL backup system for SuProxy production, including backup procedures, restoration, verification, and disaster recovery.

---

## Table of Contents

1. [Backup Strategy](#backup-strategy)
2. [Automated Backups](#automated-backups)
3. [Manual Backup](#manual-backup)
4. [Restore Procedures](#restore-procedures)
5. [Backup Verification](#backup-verification)
6. [Monitoring](#monitoring)
7. [Disaster Recovery](#disaster-recovery)
8. [Troubleshooting](#troubleshooting)

---

## Backup Strategy

### Schedule
- **Frequency:** Daily
- **Time:** 2:00 AM (server timezone)
- **Retention:** 30 days
- **Format:** Compressed SQL dumps (`.sql.gz`)

### What's Backed Up
- All database tables and data
- Database schema
- Sequences and indexes
- Constraints and triggers
- User permissions (excluded for portability)

### Storage
- **Location:** `/opt/suproxy/backups/`
- **Naming:** `suproxy_backup_YYYYMMDD_HHMMSS.sql.gz`
- **Checksums:** SHA-256 for each backup
- **Logs:** `/opt/suproxy/backups/backup.log`

---

## Automated Backups

### Initial Setup

#### 1. Configure Cron Job

```bash
cd /opt/suproxy
chmod +x scripts/setup-backup-cron.sh
./scripts/setup-backup-cron.sh
```

This script will:
- Make all backup scripts executable
- Add daily cron job (2:00 AM)
- Optionally run a test backup

#### 2. Verify Configuration

```bash
# Check cron schedule
crontab -l | grep backup

# Expected output:
# 0 2 * * * /opt/suproxy/scripts/backup-db.sh >> /opt/suproxy/backups/backup.log 2>&1
```

#### 3. Test Backup

```bash
# Run manual backup
/opt/suproxy/scripts/backup-db.sh

# Verify backup was created
ls -lh /opt/suproxy/backups/
```

### Backup Process

The automated backup performs these steps:

1. **Pre-flight checks**
   - Verify database container is running
   - Check disk space
   - Retrieve database credentials

2. **Create backup**
   - Dump database using `pg_dump`
   - Compress with gzip
   - Name with timestamp

3. **Verification**
   - Test gzip integrity
   - Generate SHA-256 checksum

4. **Retention management**
   - Remove backups older than 30 days
   - Keep last 30 daily backups minimum

5. **Logging**
   - Record all operations
   - Log errors and warnings
   - Generate summary

### Backup Output

```
[2024-08-21 02:00:01] Starting PostgreSQL backup...
[2024-08-21 02:00:01] Database: suproxy_prod, User: suproxy_prod
[2024-08-21 02:00:01] Creating backup: suproxy_backup_20240821_020001.sql.gz
[SUCCESS] Backup created successfully: suproxy_backup_20240821_020001.sql.gz (15M)
[SUCCESS] Backup integrity verified
[SUCCESS] Checksum created: suproxy_backup_20240821_020001.sql.gz.sha256
═══════════════════════════════════════════════════
Backup Summary:
  - Latest backup: suproxy_backup_20240821_020001.sql.gz
  - Backup size: 15M
  - Total backups: 30
  - Total backup size: 450M
  - Retention: 30 days
═══════════════════════════════════════════════════
[SUCCESS] Backup completed successfully!
```

---

## Manual Backup

### Create Backup

```bash
# Standard backup
/opt/suproxy/scripts/backup-db.sh

# Background backup (returns immediately)
/opt/suproxy/scripts/backup-db.sh &

# Monitor progress
tail -f /opt/suproxy/backups/backup.log
```

### Backup to Custom Location

```bash
# Modify BACKUP_DIR in script temporarily
BACKUP_DIR="/path/to/custom/location" /opt/suproxy/scripts/backup-db.sh
```

### Backup Before Critical Operations

```bash
# Before database migration
/opt/suproxy/scripts/backup-db.sh

# Before major update
/opt/suproxy/scripts/backup-db.sh

# Before schema changes
/opt/suproxy/scripts/backup-db.sh
```

---

## Restore Procedures

### Quick Start

```bash
# Interactive mode (select from list)
/opt/suproxy/scripts/restore-db.sh

# Restore specific backup
/opt/suproxy/scripts/restore-db.sh suproxy_backup_20240821_020001.sql.gz

# Restore from full path
/opt/suproxy/scripts/restore-db.sh /path/to/backup.sql.gz
```

### Step-by-Step Restore

#### 1. List Available Backups

```bash
ls -lh /opt/suproxy/backups/suproxy_backup_*.sql.gz
```

#### 2. Verify Backup Integrity

```bash
# Test gzip integrity
gunzip -t /opt/suproxy/backups/suproxy_backup_20240821_020001.sql.gz

# Verify checksum
cd /opt/suproxy/backups
sha256sum -c suproxy_backup_20240821_020001.sql.gz.sha256
```

#### 3. Stop API Service (Recommended)

```bash
docker-compose -f docker-compose.production-https.yml stop api

# Or for nginx setup
docker-compose -f docker-compose.production-nginx.yml stop api
```

#### 4. Run Restore

```bash
/opt/suproxy/scripts/restore-db.sh suproxy_backup_20240821_020001.sql.gz
```

**The script will:**
- Show backup details
- Verify integrity
- Create pre-restore backup
- Request confirmation
- Drop existing database
- Restore from backup
- Verify restoration

#### 5. Restart API Service

```bash
docker-compose -f docker-compose.production-https.yml start api

# Verify health
curl https://api.yourdomain.com/health
```

### Restore Output

```
╔════════════════════════════════════════════════════════════╗
║        SuProxy Database Restore                            ║
╚════════════════════════════════════════════════════════════╝

[INFO] Selected backup: suproxy_backup_20240821_020001.sql.gz
[INFO] Verifying backup integrity...
[SUCCESS] Backup integrity verified
[SUCCESS] Checksum verified

╔════════════════════════════════════════════════════════════╗
║                    ⚠️  WARNING  ⚠️                          ║
╠════════════════════════════════════════════════════════════╣
║  This will:                                                ║
║  1. Create a backup of current database                   ║
║  2. Drop the existing database                             ║
║  3. Restore from: suproxy_backup_20240821_020001.sql.gz   ║
╚════════════════════════════════════════════════════════════╝

Are you sure you want to continue? (yes/no): yes

[INFO] Creating pre-restore backup...
[SUCCESS] Pre-restore backup created: pre_restore_20240821_030000.sql.gz
[INFO] Starting database restore...
[WARNING] Dropping database suproxy_prod...
[INFO] Recreating database suproxy_prod...
[INFO] Restoring data from backup...
[SUCCESS] Database restored successfully
[INFO] Verifying restore...
[SUCCESS] Restore verified - 25 tables found

╔════════════════════════════════════════════════════════════╗
║         Database Restore Completed Successfully!          ║
╚════════════════════════════════════════════════════════════╝

Backup used: suproxy_backup_20240821_020001.sql.gz
Pre-restore backup saved in /opt/suproxy/backups
```

---

## Backup Verification

### Automated Verification

```bash
# Verify latest backup
/opt/suproxy/scripts/verify-backup.sh
```

### What Gets Verified

1. **File Integrity**
   - Gzip compression integrity
   - File not corrupted

2. **Checksum Verification**
   - SHA-256 checksum matches

3. **SQL Content Validation**
   - File contains valid SQL
   - Not empty or truncated

4. **Restore Test**
   - Creates temporary test database
   - Restores backup to test DB
   - Verifies data integrity

5. **Table Verification**
   - Counts restored tables
   - Checks critical tables exist
   - Validates schema

### Verification Output

```
╔════════════════════════════════════════════════════════════╗
║       SuProxy Backup Verification                          ║
╚════════════════════════════════════════════════════════════╝

[INFO] Testing backup: suproxy_backup_20240821_020001.sql.gz
[INFO] Test 1: Checking file integrity...
[SUCCESS] ✓ File integrity check passed
[INFO] Test 2: Verifying checksum...
[SUCCESS] ✓ Checksum verification passed
[INFO] Test 3: Validating SQL content...
[SUCCESS] ✓ SQL content valid (15000 lines)
[INFO] Test 4: Testing restore capability...
[INFO] Creating test database: suproxy_restore_test
[INFO] Restoring to test database...
[SUCCESS] ✓ Restore successful
[INFO] Test 5: Verifying restored data...
[SUCCESS] ✓ Found 25 tables in restored database
[SUCCESS] ✓ Critical table 'users' found
[INFO] Cleaning up test database...
[SUCCESS] ✓ Cleanup completed

╔════════════════════════════════════════════════════════════╗
║         All Backup Verification Tests Passed! ✓           ║
╚════════════════════════════════════════════════════════════╝

Backup Summary:
  - File: suproxy_backup_20240821_020001.sql.gz
  - Size: 15M
  - Tables: 25
  - Status: Ready for production use
```

### Schedule Verification

Add to cron for weekly verification:

```bash
# Verify backups every Sunday at 3 AM
0 3 * * 0 /opt/suproxy/scripts/verify-backup.sh >> /opt/suproxy/backups/verify.log 2>&1
```

---

## Monitoring

### Check Backup Status

```bash
# View recent backups
ls -lht /opt/suproxy/backups/suproxy_backup_*.sql.gz | head -5

# Count backups
find /opt/suproxy/backups -name "suproxy_backup_*.sql.gz" | wc -l

# Total backup size
du -sh /opt/suproxy/backups/
```

### Monitor Backup Logs

```bash
# View latest backup log
tail -50 /opt/suproxy/backups/backup.log

# Follow backup in real-time
tail -f /opt/suproxy/backups/backup.log

# Check for errors
grep ERROR /opt/suproxy/backups/backup.log

# Last successful backup
grep "Backup completed successfully" /opt/suproxy/backups/backup.log | tail -1
```

### Check Disk Space

```bash
# Backup directory disk usage
df -h /opt/suproxy/backups

# Warn if less than 10GB free
AVAILABLE=$(df /opt/suproxy/backups | tail -1 | awk '{print $4}')
if [ $AVAILABLE -lt 10000000 ]; then
    echo "WARNING: Low disk space"
fi
```

### Alert on Backup Failure

```bash
# Check last backup status
if ! grep -q "Backup completed successfully" /opt/suproxy/backups/backup.log; then
    # Send alert (configure your alerting method)
    echo "Backup failed!" | mail -s "Backup Alert" admin@yourdomain.com
fi
```

---

## Disaster Recovery

### Scenario 1: Database Corruption

```bash
# 1. Stop API
docker-compose -f docker-compose.production-https.yml stop api

# 2. Restore from latest backup
/opt/suproxy/scripts/restore-db.sh

# 3. Verify restoration
docker exec suproxy-postgres-prod psql -U suproxy_prod -d suproxy_prod -c "SELECT COUNT(*) FROM users;"

# 4. Restart API
docker-compose -f docker-compose.production-https.yml start api
```

### Scenario 2: Accidental Data Deletion

```bash
# 1. Identify when data was lost
ls -lt /opt/suproxy/backups/

# 2. Find backup before deletion
# (e.g., deletion happened at 10 AM, use 2 AM backup)

# 3. Restore from that backup
/opt/suproxy/scripts/restore-db.sh suproxy_backup_YYYYMMDD_020001.sql.gz
```

### Scenario 3: Server Failure

```bash
# On new server:

# 1. Install Docker and Docker Compose
curl -fsSL https://get.docker.com | sh

# 2. Clone repository
git clone https://github.com/your-org/suproxy-backend.git
cd suproxy-backend

# 3. Copy backup files to new server
scp old-server:/opt/suproxy/backups/*.sql.gz /opt/suproxy/backups/

# 4. Start database only
docker-compose -f docker-compose.production-https.yml up -d postgres

# 5. Wait for postgres to be ready
sleep 10

# 6. Restore database
/opt/suproxy/scripts/restore-db.sh

# 7. Start all services
docker-compose -f docker-compose.production-https.yml up -d
```

### Scenario 4: Point-in-Time Recovery

```bash
# Combine latest backup + transaction logs (if available)

# 1. Restore from latest backup
/opt/suproxy/scripts/restore-db.sh

# 2. If WAL archiving is enabled, replay transaction logs
# (Advanced - requires WAL configuration)
```

---

## Troubleshooting

### Backup Fails

**Problem:** Backup script exits with error

**Solutions:**

```bash
# Check if database container is running
docker ps | grep postgres

# Check disk space
df -h /opt/suproxy/backups

# Check database connectivity
docker exec suproxy-postgres-prod psql -U suproxy_prod -l

# Check backup script permissions
ls -l /opt/suproxy/scripts/backup-db.sh
chmod +x /opt/suproxy/scripts/backup-db.sh

# Run backup manually with verbose output
bash -x /opt/suproxy/scripts/backup-db.sh
```

### Restore Fails

**Problem:** Restore script can't restore backup

**Solutions:**

```bash
# Verify backup integrity
gunzip -t /opt/suproxy/backups/suproxy_backup_*.sql.gz

# Check backup file permissions
ls -l /opt/suproxy/backups/suproxy_backup_*.sql.gz

# Manually test restore
gunzip -c /opt/suproxy/backups/suproxy_backup_*.sql.gz | \
    docker exec -i suproxy-postgres-prod psql -U suproxy_prod -d suproxy_prod

# Check for conflicting connections
docker exec suproxy-postgres-prod psql -U suproxy_prod -d postgres -c \
    "SELECT * FROM pg_stat_activity WHERE datname='suproxy_prod';"
```

### Backup Too Large

**Problem:** Backup files consuming too much disk space

**Solutions:**

```bash
# Check backup sizes
du -h /opt/suproxy/backups/suproxy_backup_*.sql.gz

# Reduce retention period (edit backup script)
nano /opt/suproxy/scripts/backup-db.sh
# Change: BACKUP_RETENTION_DAYS=30 to BACKUP_RETENTION_DAYS=14

# Use higher compression
# Change pg_dump to use --compress=9

# Move old backups to archive storage
find /opt/suproxy/backups -name "suproxy_backup_*.sql.gz" -mtime +7 \
    -exec mv {} /archive/location/ \;
```

### Cron Not Running

**Problem:** Backups not running automatically

**Solutions:**

```bash
# Check cron service
systemctl status cron

# Check cron logs
grep CRON /var/log/syslog

# Verify cron entry
crontab -l | grep backup

# Test cron entry manually
/opt/suproxy/scripts/backup-db.sh

# Check cron permissions
ls -l /opt/suproxy/scripts/backup-db.sh

# Ensure full paths in cron
# Edit: crontab -e
# Use: /opt/suproxy/scripts/backup-db.sh
# Not: ./scripts/backup-db.sh
```

---

## Best Practices

### Security

1. **Protect Backup Files**
   ```bash
   chmod 600 /opt/suproxy/backups/*.sql.gz
   chown root:root /opt/suproxy/backups/*.sql.gz
   ```

2. **Encrypt Sensitive Backups**
   ```bash
   # Encrypt backup
   gpg --encrypt --recipient admin@yourdomain.com backup.sql.gz
   
   # Decrypt for restore
   gpg --decrypt backup.sql.gz.gpg > backup.sql.gz
   ```

3. **Off-site Backup Copy**
   ```bash
   # Sync to remote server
   rsync -avz /opt/suproxy/backups/ backup-server:/backups/suproxy/
   
   # Or upload to S3
   aws s3 sync /opt/suproxy/backups/ s3://your-bucket/backups/
   ```

### Maintenance

1. **Monthly Verification**
   - Run verification script
   - Test restore procedure
   - Review backup logs

2. **Quarterly Disaster Recovery Drill**
   - Simulate server failure
   - Practice full restoration
   - Document time to recover

3. **Annual Review**
   - Review retention policy
   - Check disk space trends
   - Update documentation

---

## Configuration Reference

### Backup Script Variables

```bash
BACKUP_DIR="/opt/suproxy/backups"        # Backup storage location
BACKUP_RETENTION_DAYS=30                  # Keep backups for 30 days
DB_CONTAINER="suproxy-postgres-prod"      # Database container name
```

### Cron Schedule Examples

```bash
# Daily at 2 AM
0 2 * * * /opt/suproxy/scripts/backup-db.sh

# Every 6 hours
0 */6 * * * /opt/suproxy/scripts/backup-db.sh

# Weekly on Sunday at 2 AM
0 2 * * 0 /opt/suproxy/scripts/backup-db.sh

# Every 12 hours
0 */12 * * * /opt/suproxy/scripts/backup-db.sh
```

---

## Support

### Quick Links

- Backup Issues: Check `/opt/suproxy/backups/backup.log`
- Restore Help: Run with `-h` flag
- Verification: `/opt/suproxy/scripts/verify-backup.sh`

### Emergency Contacts

- Database Administrator: [Contact Info]
- System Administrator: [Contact Info]
- On-call Support: [Contact Info]

---

**Document Version:** 1.0.0  
**Last Reviewed:** 2024-08-21  
**Next Review:** 2024-09-21
