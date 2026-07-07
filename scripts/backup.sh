#!/bin/bash

set -e

# Colors for output
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

BACKUP_DIR="./backups"
TIMESTAMP=$(date +%Y%m%d_%H%M%S)
DB_BACKUP_FILE="suproxy_db_${TIMESTAMP}.sql"

echo -e "${GREEN}========================================${NC}"
echo -e "${GREEN}SuProxy Database Backup Script${NC}"
echo -e "${GREEN}========================================${NC}"

# Create backup directory if it doesn't exist
mkdir -p "$BACKUP_DIR"

# Load environment variables
if [ -f ".env.production" ]; then
    set -a
    source .env.production
    set +a
elif [ -f ".env" ]; then
    set -a
    source .env
    set +a
fi

# Backup database
echo -e "${YELLOW}Creating database backup...${NC}"
docker-compose -f docker-compose.production.yml exec -T postgres pg_dump \
    -U ${DB_USER} \
    -d ${DB_NAME} \
    --format=plain \
    --no-owner \
    --no-acl > "${BACKUP_DIR}/${DB_BACKUP_FILE}"

if [ $? -eq 0 ]; then
    echo -e "${GREEN}Database backup created: ${BACKUP_DIR}/${DB_BACKUP_FILE}${NC}"
    
    # Compress backup
    gzip "${BACKUP_DIR}/${DB_BACKUP_FILE}"
    echo -e "${GREEN}Backup compressed: ${BACKUP_DIR}/${DB_BACKUP_FILE}.gz${NC}"
    
    # Delete backups older than 30 days
    find "$BACKUP_DIR" -name "suproxy_db_*.sql.gz" -mtime +30 -delete
    echo -e "${GREEN}Old backups cleaned up (>30 days)${NC}"
else
    echo -e "${RED}Database backup failed${NC}"
    exit 1
fi

echo -e "${GREEN}========================================${NC}"
echo -e "${GREEN}Backup completed successfully${NC}"
echo -e "${GREEN}========================================${NC}"
