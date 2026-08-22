#!/bin/bash
# PostgreSQL Backup Verification Script
# Tests backup integrity and restore capability

set -e

# Configuration
BACKUP_DIR="/opt/suproxy/backups"
DB_CONTAINER="suproxy-postgres-prod"
TEST_DB="suproxy_restore_test"

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

log() { echo -e "${BLUE}[INFO]${NC} $1"; }
log_error() { echo -e "${RED}[ERROR]${NC} $1"; }
log_success() { echo -e "${GREEN}[SUCCESS]${NC} $1"; }
log_warning() { echo -e "${YELLOW}[WARNING]${NC} $1"; }

echo "╔════════════════════════════════════════════════════════════╗"
echo "║       SuProxy Backup Verification                          ║"
echo "╚════════════════════════════════════════════════════════════╝"
echo ""

# Check if container is running
if ! docker ps | grep -q "${DB_CONTAINER}"; then
    log_error "Database container ${DB_CONTAINER} is not running!"
    exit 1
fi

# Get latest backup
LATEST_BACKUP=$(find "${BACKUP_DIR}" -name "suproxy_backup_*.sql.gz" -type f | sort -r | head -n 1)

if [ -z "$LATEST_BACKUP" ]; then
    log_error "No backups found in ${BACKUP_DIR}"
    exit 1
fi

BACKUP_NAME=$(basename "$LATEST_BACKUP")
log "Testing backup: ${BACKUP_NAME}"

# Test 1: File integrity
log "Test 1: Checking file integrity..."
if gunzip -t "$LATEST_BACKUP" 2>/dev/null; then
    log_success "✓ File integrity check passed"
else
    log_error "✗ File integrity check failed"
    exit 1
fi

# Test 2: Checksum verification
log "Test 2: Verifying checksum..."
if [ -f "${LATEST_BACKUP}.sha256" ]; then
    cd "$(dirname "$LATEST_BACKUP")"
    if sha256sum -c "$(basename "${LATEST_BACKUP}.sha256")" 2>/dev/null; then
        log_success "✓ Checksum verification passed"
    else
        log_error "✗ Checksum verification failed"
        exit 1
    fi
else
    log_warning "⚠ No checksum file found"
fi

# Test 3: SQL content validation
log "Test 3: Validating SQL content..."
SQL_LINES=$(gunzip -c "$LATEST_BACKUP" | wc -l)
if [ $SQL_LINES -gt 0 ]; then
    log_success "✓ SQL content valid (${SQL_LINES} lines)"
else
    log_error "✗ SQL content appears empty"
    exit 1
fi

# Test 4: Restore to test database
log "Test 4: Testing restore capability..."

DB_USER=$(docker exec "${DB_CONTAINER}" printenv POSTGRES_USER)

# Drop test database if exists
docker exec "${DB_CONTAINER}" psql -U "${DB_USER}" -d postgres -c \
    "DROP DATABASE IF EXISTS ${TEST_DB};" > /dev/null 2>&1

# Create test database
log "Creating test database: ${TEST_DB}"
docker exec "${DB_CONTAINER}" psql -U "${DB_USER}" -d postgres -c \
    "CREATE DATABASE ${TEST_DB};" > /dev/null

# Restore to test database
log "Restoring to test database..."
if gunzip -c "$LATEST_BACKUP" | docker exec -i "${DB_CONTAINER}" psql -U "${DB_USER}" -d "${TEST_DB}" > /dev/null 2>&1; then
    log_success "✓ Restore successful"
else
    log_error "✗ Restore failed"
    # Cleanup
    docker exec "${DB_CONTAINER}" psql -U "${DB_USER}" -d postgres -c \
        "DROP DATABASE IF EXISTS ${TEST_DB};" > /dev/null 2>&1
    exit 1
fi

# Test 5: Verify restored data
log "Test 5: Verifying restored data..."

# Check tables
TABLE_COUNT=$(docker exec "${DB_CONTAINER}" psql -U "${DB_USER}" -d "${TEST_DB}" -t -c \
    "SELECT COUNT(*) FROM information_schema.tables WHERE table_schema='public';" | tr -d ' ')

if [ $TABLE_COUNT -gt 0 ]; then
    log_success "✓ Found ${TABLE_COUNT} tables in restored database"
else
    log_error "✗ No tables found in restored database"
    docker exec "${DB_CONTAINER}" psql -U "${DB_USER}" -d postgres -c \
        "DROP DATABASE IF EXISTS ${TEST_DB};" > /dev/null 2>&1
    exit 1
fi

# Check for users table (critical table)
USER_TABLE=$(docker exec "${DB_CONTAINER}" psql -U "${DB_USER}" -d "${TEST_DB}" -t -c \
    "SELECT COUNT(*) FROM information_schema.tables WHERE table_schema='public' AND table_name='users';" | tr -d ' ')

if [ $USER_TABLE -eq 1 ]; then
    log_success "✓ Critical table 'users' found"
else
    log_warning "⚠ Critical table 'users' not found"
fi

# Cleanup test database
log "Cleaning up test database..."
docker exec "${DB_CONTAINER}" psql -U "${DB_USER}" -d postgres -c \
    "DROP DATABASE ${TEST_DB};" > /dev/null

log_success "✓ Cleanup completed"

echo ""
echo "╔════════════════════════════════════════════════════════════╗"
echo "║         All Backup Verification Tests Passed! ✓           ║"
echo "╚════════════════════════════════════════════════════════════╝"
echo ""
log "Backup Summary:"
log "  - File: ${BACKUP_NAME}"
log "  - Size: $(du -h "$LATEST_BACKUP" | cut -f1)"
log "  - Tables: ${TABLE_COUNT}"
log "  - Status: Ready for production use"
echo ""

exit 0
