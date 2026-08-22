# Database Backup Quick Reference

## Daily Operations

### Check Latest Backup
```bash
ls -lht /opt/suproxy/backups/suproxy_backup_*.sql.gz | head -1
```

### View Backup Logs
```bash
tail -50 /opt/suproxy/backups/backup.log
```

### Manual Backup
```bash
/opt/suproxy/scripts/backup-db.sh
```

### Verify Backup
```bash
/opt/suproxy/scripts/verify-backup.sh
```

## Restore Operations

### Interactive Restore
```bash
/opt/suproxy/scripts/restore-db.sh
```

### Restore Specific Backup
```bash
/opt/suproxy/scripts/restore-db.sh suproxy_backup_20240821_020001.sql.gz
```

## Emergency Procedures

### Quick Restore (Last Backup)
```bash
# Stop API
docker-compose -f docker-compose.production-https.yml stop api

# Restore
/opt/suproxy/scripts/restore-db.sh

# Start API
docker-compose -f docker-compose.production-https.yml start api
```

### Check Backup Health
```bash
# Last successful backup
grep "Backup completed successfully" /opt/suproxy/backups/backup.log | tail -1

# Count backups
find /opt/suproxy/backups -name "suproxy_backup_*.sql.gz" | wc -l

# Total size
du -sh /opt/suproxy/backups/
```

## Monitoring Commands

```bash
# Backup status
systemctl status cron | grep backup

# Cron schedule
crontab -l | grep backup

# Disk space
df -h /opt/suproxy/backups

# Recent errors
grep ERROR /opt/suproxy/backups/backup.log | tail -5
```

## File Locations

- **Backups:** `/opt/suproxy/backups/`
- **Scripts:** `/opt/suproxy/scripts/`
- **Logs:** `/opt/suproxy/backups/backup.log`
- **Docs:** `/opt/suproxy/docs/DATABASE_BACKUP.md`

## Important Notes

⚠️ **Always verify backup before critical operations**  
⚠️ **Test restore procedure quarterly**  
⚠️ **Monitor disk space regularly**  
⚠️ **Keep off-site backup copy**
