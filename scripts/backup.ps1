# PowerShell backup script for Windows

$ErrorActionPreference = "Stop"

$BACKUP_DIR = ".\backups"
$TIMESTAMP = Get-Date -Format "yyyyMMdd_HHmmss"
$DB_BACKUP_FILE = "suproxy_db_${TIMESTAMP}.sql"

Write-Host "========================================" -ForegroundColor Green
Write-Host "SuProxy Database Backup Script" -ForegroundColor Green
Write-Host "========================================" -ForegroundColor Green

# Create backup directory if it doesn't exist
if (-not (Test-Path $BACKUP_DIR)) {
    New-Item -ItemType Directory -Path $BACKUP_DIR | Out-Null
}

# Load environment variables
if (Test-Path ".env.production") {
    Get-Content .env.production | ForEach-Object {
        if ($_ -match '^([^#][^=]+)=(.*)$') {
            $name = $matches[1].Trim()
            $value = $matches[2].Trim()
            Set-Item -Path "env:$name" -Value $value
        }
    }
}
elseif (Test-Path ".env") {
    Get-Content .env | ForEach-Object {
        if ($_ -match '^([^#][^=]+)=(.*)$') {
            $name = $matches[1].Trim()
            $value = $matches[2].Trim()
            Set-Item -Path "env:$name" -Value $value
        }
    }
}

# Backup database
Write-Host "Creating database backup..." -ForegroundColor Yellow
$dbUser = $env:DB_USER
$dbName = $env:DB_NAME

docker-compose -f docker-compose.production.yml exec -T postgres pg_dump `
    -U $dbUser `
    -d $dbName `
    --format=plain `
    --no-owner `
    --no-acl | Out-File -Encoding UTF8 -FilePath "$BACKUP_DIR\$DB_BACKUP_FILE"

if ($LASTEXITCODE -eq 0) {
    Write-Host "Database backup created: $BACKUP_DIR\$DB_BACKUP_FILE" -ForegroundColor Green
    
    # Compress backup
    Compress-Archive -Path "$BACKUP_DIR\$DB_BACKUP_FILE" -DestinationPath "$BACKUP_DIR\$DB_BACKUP_FILE.zip"
    Remove-Item "$BACKUP_DIR\$DB_BACKUP_FILE"
    Write-Host "Backup compressed: $BACKUP_DIR\$DB_BACKUP_FILE.zip" -ForegroundColor Green
    
    # Delete backups older than 30 days
    $cutoffDate = (Get-Date).AddDays(-30)
    Get-ChildItem -Path $BACKUP_DIR -Filter "suproxy_db_*.sql.zip" | 
        Where-Object { $_.LastWriteTime -lt $cutoffDate } | 
        Remove-Item
    Write-Host "Old backups cleaned up (>30 days)" -ForegroundColor Green
}
else {
    Write-Host "Database backup failed" -ForegroundColor Red
    exit 1
}

Write-Host "========================================" -ForegroundColor Green
Write-Host "Backup completed successfully" -ForegroundColor Green
Write-Host "========================================" -ForegroundColor Green
