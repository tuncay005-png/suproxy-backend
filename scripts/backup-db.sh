#!/bin/bash
# PostgreSQL Automated Backup Script for SuProxy Production
# This script creates compressed database backups with rotation

set -e

# Configuration
BACKUP_DIR="/opt/suproxy/backups"
BACKUP_RETENTION_DAYS=30
DB_CONTAINER="suproxy-postgres-prod"
TIMESTAMP=$(date +%Y%m%d_%H%M%S)
BACKUP_FILE="suproxy_backup_${TIMESTAMP}.sql.gz"
LOG_FILE="${BACKUP_DIR}/backup.log"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Logging function
log() {
    echo "[$(date +'%Y-%m-%d %H:%M:%S')] $1" | tee -a "${LOG_FILE}"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1" | tee -a "${LOG_FILE}"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1" | tee -a "${LOG_FILE}"
}

log_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1" | tee -a "${LOG_FILE}"
}

# Create backup directory if it doesn't exist
mkdir -p "${BACKUP_DIR}"

log "Starting PostgreSQL backup..."

# Check if container is running
if ! docker ps | grep -q "${DB_CONTAINER}"; then
    log_error "Database container ${DB_CONTAINER} is not running!"
    exit 1
fi

# Get database credentials from environment or container
DB_USER=$(docker exec "${DB_CONTAINER}" printenv POSTGRES_USER)
DB_NAME=$(docker exec "${DB_CONTAINER}" printenv POSTGRES_DB)

if [ -z "$DB_USER" ] || [ -z "$DB_NAME" ]; then
    log_error "Failed to retrieve database credentials"
    exit 1
fi

log "Database: ${DB_NAME}, User: ${DB_USER}"

# Create backup
log "Creating backup: ${BACKUP_FILE}"
if docker exec "${DB_CONTAINER}" pg_dump -U "${DB_USER}" -d "${DB_NAME}" \
    --format=plain --no-owner --no-acl | gzip > "${BACKUP_DIR}/${BACKUP_FILE}"; then
    
    BACKUP_SIZE=$(du -h "${BACKUP_DIR}/${BACKUP_FILE}" | cut -f1)
    log_success "Backup created successfully: ${BACKUP_FILE} (${BACKUP_SIZE})"
else
    log_error "Backup failed!"
    exit 1
fi

# Verify backup integrity
log "Verifying backup integrity..."
if gunzip -t "${BACKUP_DIR}/${BACKUP_FILE}"; then
    log_success "Backup integrity verified"
else
    log_error "Backup integrity check failed!"
    exit 1
fi

# Create checksum
log "Creating checksum..."
cd "${BACKUP_DIR}"
sha256sum "${BACKUP_FILE}" > "${BACKUP_FILE}.sha256"
log_success "Checksum created: ${BACKUP_FILE}.sha256"

# Remove old backups (retention policy)
log "Applying retention policy (${BACKUP_RETENTION_DAYS} days)..."
DELETED_COUNT=0
find "${BACKUP_DIR}" -name "suproxy_backup_*.sql.gz" -type f -mtime +${BACKUP_RETENTION_DAYS} | while read OLD_BACKUP; do
    log "Deleting old backup: $(basename ${OLD_BACKUP})"
    rm -f "${OLD_BACKUP}"
    rm -f "${OLD_BACKUP}.sha256"
    DELETED_COUNT=$((DELETED_COUNT + 1))
done

if [ $DELETED_COUNT -eq 0 ]; then
    log "No old backups to delete"
else
    log_success "Deleted ${DELETED_COUNT} old backup(s)"
fi

# Summary
TOTAL_BACKUPS=$(find "${BACKUP_DIR}" -name "suproxy_backup_*.sql.gz" -type f | wc -l)
TOTAL_SIZE=$(du -sh "${BACKUP_DIR}" | cut -f1)

log "═══════════════════════════════════════════════════"
log "Backup Summary:"
log "  - Latest backup: ${BACKUP_FILE}"
log "  - Backup size: ${BACKUP_SIZE}"
log "  - Total backups: ${TOTAL_BACKUPS}"
log "  - Total backup size: ${TOTAL_SIZE}"
log "  - Retention: ${BACKUP_RETENTION_DAYS} days"
log "═══════════════════════════════════════════════════"

log_success "Backup completed successfully!"

# Send notification (optional - can be configured later)
# curl -X POST https://your-monitoring-service.com/webhook \
#   -H "Content-Type: application/json" \
#   -d "{\"status\":\"success\",\"backup\":\"${BACKUP_FILE}\",\"size\":\"${BACKUP_SIZE}\"}"

exit 0
