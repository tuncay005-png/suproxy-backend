#!/bin/bash
# PostgreSQL Restore Script for SuProxy Production
# This script restores database from backup with safety checks

set -e

# Configuration
BACKUP_DIR="/opt/suproxy/backups"
DB_CONTAINER="suproxy-postgres-prod"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Logging functions
log() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

log_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

# Usage information
usage() {
    echo "Usage: $0 [backup_file]"
    echo ""
    echo "Restore PostgreSQL database from backup"
    echo ""
    echo "Options:"
    echo "  backup_file    Path to backup file (optional - will show list if not provided)"
    echo ""
    echo "Examples:"
    echo "  $0                                      # Interactive mode - select from list"
    echo "  $0 suproxy_backup_20240821_120000.sql.gz"
    echo "  $0 /path/to/backup.sql.gz"
    echo ""
    exit 1
}

# List available backups
list_backups() {
    log "Available backups in ${BACKUP_DIR}:"
    echo ""
    
    BACKUPS=$(find "${BACKUP_DIR}" -name "suproxy_backup_*.sql.gz" -type f | sort -r)
    
    if [ -z "$BACKUPS" ]; then
        log_error "No backups found in ${BACKUP_DIR}"
        exit 1
    fi
    
    INDEX=1
    declare -A BACKUP_MAP
    
    while IFS= read -r BACKUP; do
        FILENAME=$(basename "$BACKUP")
        SIZE=$(du -h "$BACKUP" | cut -f1)
        DATE=$(stat -c %y "$BACKUP" | cut -d'.' -f1)
        
        echo "${INDEX}. ${FILENAME}"
        echo "   Size: ${SIZE} | Date: ${DATE}"
        
        # Verify checksum exists
        if [ -f "${BACKUP}.sha256" ]; then
            echo "   ✓ Checksum verified"
        else
            echo "   ⚠ No checksum file"
        fi
        echo ""
        
        BACKUP_MAP[$INDEX]="$BACKUP"
        INDEX=$((INDEX + 1))
    done <<< "$BACKUPS"
    
    echo -n "Select backup number (or 'q' to quit): "
    read SELECTION
    
    if [ "$SELECTION" = "q" ]; then
        log "Restore cancelled"
        exit 0
    fi
    
    SELECTED_BACKUP="${BACKUP_MAP[$SELECTION]}"
    
    if [ -z "$SELECTED_BACKUP" ]; then
        log_error "Invalid selection"
        exit 1
    fi
    
    echo "$SELECTED_BACKUP"
}

# Verify backup integrity
verify_backup() {
    local BACKUP_FILE=$1
    
    log "Verifying backup integrity..."
    
    # Check if file exists
    if [ ! -f "$BACKUP_FILE" ]; then
        log_error "Backup file not found: $BACKUP_FILE"
        exit 1
    fi
    
    # Check gzip integrity
    if ! gunzip -t "$BACKUP_FILE" 2>/dev/null; then
        log_error "Backup file is corrupted (gzip test failed)"
        exit 1
    fi
    
    # Verify checksum if available
    if [ -f "${BACKUP_FILE}.sha256" ]; then
        log "Verifying checksum..."
        cd "$(dirname "$BACKUP_FILE")"
        if sha256sum -c "$(basename "${BACKUP_FILE}.sha256")" 2>/dev/null; then
            log_success "Checksum verified"
        else
            log_error "Checksum verification failed!"
            exit 1
        fi
    else
        log_warning "No checksum file found - skipping verification"
    fi
    
    log_success "Backup integrity verified"
}

# Create pre-restore backup
create_pre_restore_backup() {
    log "Creating pre-restore backup..."
    
    DB_USER=$(docker exec "${DB_CONTAINER}" printenv POSTGRES_USER)
    DB_NAME=$(docker exec "${DB_CONTAINER}" printenv POSTGRES_DB)
    
    PRE_RESTORE_FILE="${BACKUP_DIR}/pre_restore_$(date +%Y%m%d_%H%M%S).sql.gz"
    
    if docker exec "${DB_CONTAINER}" pg_dump -U "${DB_USER}" -d "${DB_NAME}" \
        --format=plain --no-owner --no-acl | gzip > "${PRE_RESTORE_FILE}"; then
        log_success "Pre-restore backup created: $(basename ${PRE_RESTORE_FILE})"
    else
        log_error "Failed to create pre-restore backup"
        exit 1
    fi
}

# Perform restore
restore_database() {
    local BACKUP_FILE=$1
    
    log "Starting database restore..."
    
    # Get database credentials
    DB_USER=$(docker exec "${DB_CONTAINER}" printenv POSTGRES_USER)
    DB_NAME=$(docker exec "${DB_CONTAINER}" printenv POSTGRES_DB)
    
    if [ -z "$DB_USER" ] || [ -z "$DB_NAME" ]; then
        log_error "Failed to retrieve database credentials"
        exit 1
    fi
    
    log "Database: ${DB_NAME}, User: ${DB_USER}"
    
    # Terminate existing connections
    log "Terminating existing database connections..."
    docker exec "${DB_CONTAINER}" psql -U "${DB_USER}" -d postgres -c \
        "SELECT pg_terminate_backend(pid) FROM pg_stat_activity WHERE datname='${DB_NAME}' AND pid <> pg_backend_pid();" \
        > /dev/null 2>&1 || true
    
    # Drop and recreate database
    log_warning "Dropping database ${DB_NAME}..."
    docker exec "${DB_CONTAINER}" psql -U "${DB_USER}" -d postgres -c \
        "DROP DATABASE IF EXISTS ${DB_NAME};" > /dev/null
    
    log "Recreating database ${DB_NAME}..."
    docker exec "${DB_CONTAINER}" psql -U "${DB_USER}" -d postgres -c \
        "CREATE DATABASE ${DB_NAME};" > /dev/null
    
    # Restore from backup
    log "Restoring data from backup..."
    if gunzip -c "$BACKUP_FILE" | docker exec -i "${DB_CONTAINER}" psql -U "${DB_USER}" -d "${DB_NAME}" > /dev/null 2>&1; then
        log_success "Database restored successfully"
    else
        log_error "Database restore failed!"
        exit 1
    fi
    
    # Verify restore
    log "Verifying restore..."
    TABLE_COUNT=$(docker exec "${DB_CONTAINER}" psql -U "${DB_USER}" -d "${DB_NAME}" -t -c \
        "SELECT COUNT(*) FROM information_schema.tables WHERE table_schema='public';")
    
    log_success "Restore verified - ${TABLE_COUNT} tables found"
}

# Main script
main() {
    echo "╔════════════════════════════════════════════════════════════╗"
    echo "║        SuProxy Database Restore                            ║"
    echo "╚════════════════════════════════════════════════════════════╝"
    echo ""
    
    # Check if container is running
    if ! docker ps | grep -q "${DB_CONTAINER}"; then
        log_error "Database container ${DB_CONTAINER} is not running!"
        exit 1
    fi
    
    # Determine backup file
    if [ -z "$1" ]; then
        BACKUP_FILE=$(list_backups)
    else
        if [ -f "$1" ]; then
            BACKUP_FILE="$1"
        elif [ -f "${BACKUP_DIR}/$1" ]; then
            BACKUP_FILE="${BACKUP_DIR}/$1"
        else
            log_error "Backup file not found: $1"
            exit 1
        fi
    fi
    
    BACKUP_NAME=$(basename "$BACKUP_FILE")
    log "Selected backup: ${BACKUP_NAME}"
    
    # Verify backup
    verify_backup "$BACKUP_FILE"
    
    # Confirmation
    log_warning "╔════════════════════════════════════════════════════════════╗"
    log_warning "║                    ⚠️  WARNING  ⚠️                          ║"
    log_warning "╠════════════════════════════════════════════════════════════╣"
    log_warning "║  This will:                                                ║"
    log_warning "║  1. Create a backup of current database                   ║"
    log_warning "║  2. Drop the existing database                             ║"
    log_warning "║  3. Restore from: ${BACKUP_NAME}"
    log_warning "╚════════════════════════════════════════════════════════════╝"
    echo ""
    echo -n "Are you sure you want to continue? (yes/no): "
    read CONFIRMATION
    
    if [ "$CONFIRMATION" != "yes" ]; then
        log "Restore cancelled"
        exit 0
    fi
    
    # Create pre-restore backup
    create_pre_restore_backup
    
    # Perform restore
    restore_database "$BACKUP_FILE"
    
    echo ""
    log_success "╔════════════════════════════════════════════════════════════╗"
    log_success "║         Database Restore Completed Successfully!          ║"
    log_success "╚════════════════════════════════════════════════════════════╝"
    echo ""
    log "Backup used: ${BACKUP_NAME}"
    log "Pre-restore backup saved in ${BACKUP_DIR}"
    echo ""
}

# Handle arguments
if [ "$1" = "-h" ] || [ "$1" = "--help" ]; then
    usage
fi

main "$1"
