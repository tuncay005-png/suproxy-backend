#!/bin/bash
# Script to configure automated daily backups via cron

set -e

SCRIPT_DIR="/opt/suproxy/scripts"
BACKUP_SCRIPT="${SCRIPT_DIR}/backup-db.sh"
CRON_TIME="0 2 * * *"  # 2 AM daily

echo "╔════════════════════════════════════════════════════════════╗"
echo "║       Configure Automated Database Backups                 ║"
echo "╚════════════════════════════════════════════════════════════╝"
echo ""

# Check if backup script exists
if [ ! -f "$BACKUP_SCRIPT" ]; then
    echo "Error: Backup script not found at $BACKUP_SCRIPT"
    exit 1
fi

# Make scripts executable
echo "Making scripts executable..."
chmod +x "${SCRIPT_DIR}/backup-db.sh"
chmod +x "${SCRIPT_DIR}/restore-db.sh"
chmod +x "${SCRIPT_DIR}/verify-backup.sh"
echo "✓ Scripts are now executable"
echo ""

# Check if cron entry already exists
if crontab -l 2>/dev/null | grep -q "backup-db.sh"; then
    echo "⚠ Cron job already exists"
    echo ""
    echo "Current backup schedule:"
    crontab -l | grep "backup-db.sh"
    echo ""
    echo -n "Do you want to update it? (yes/no): "
    read UPDATE
    
    if [ "$UPDATE" != "yes" ]; then
        echo "Setup cancelled"
        exit 0
    fi
    
    # Remove old entry
    crontab -l | grep -v "backup-db.sh" | crontab -
    echo "✓ Old cron job removed"
fi

# Add new cron entry
echo "Adding cron job for daily backups at 2:00 AM..."
(crontab -l 2>/dev/null; echo "${CRON_TIME} ${BACKUP_SCRIPT} >> /opt/suproxy/backups/backup.log 2>&1") | crontab -

echo "✓ Cron job added successfully"
echo ""

# Display configuration
echo "Backup Configuration:"
echo "  - Schedule: Daily at 2:00 AM (server time)"
echo "  - Backup script: ${BACKUP_SCRIPT}"
echo "  - Log file: /opt/suproxy/backups/backup.log"
echo "  - Retention: 30 days"
echo ""

# Verify cron entry
echo "Current cron schedule:"
crontab -l | grep "backup-db.sh"
echo ""

# Test backup script
echo -n "Do you want to run a test backup now? (yes/no): "
read TEST

if [ "$TEST" = "yes" ]; then
    echo ""
    echo "Running test backup..."
    echo "════════════════════════════════════════════════════════════"
    ${BACKUP_SCRIPT}
    echo ""
    echo "════════════════════════════════════════════════════════════"
    echo "✓ Test backup completed!"
    echo ""
    echo "Verify backup with: ${SCRIPT_DIR}/verify-backup.sh"
fi

echo ""
echo "╔════════════════════════════════════════════════════════════╗"
echo "║      Automated Backup Configuration Complete!             ║"
echo "╚════════════════════════════════════════════════════════════╝"
echo ""
echo "Next steps:"
echo "  1. Monitor backup logs: tail -f /opt/suproxy/backups/backup.log"
echo "  2. Verify backups: ${SCRIPT_DIR}/verify-backup.sh"
echo "  3. Test restore: ${SCRIPT_DIR}/restore-db.sh"
echo ""

exit 0
