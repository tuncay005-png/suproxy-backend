# PowerShell deployment script for Windows

$ErrorActionPreference = "Stop"

Write-Host "========================================" -ForegroundColor Green
Write-Host "SuProxy Backend Deployment Script" -ForegroundColor Green
Write-Host "========================================" -ForegroundColor Green

# Check if .env.production exists
if (-not (Test-Path ".env.production")) {
    Write-Host "Error: .env.production file not found" -ForegroundColor Red
    Write-Host "Please create .env.production from .env.example" -ForegroundColor Yellow
    exit 1
}

# Load environment variables
Get-Content .env.production | ForEach-Object {
    if ($_ -match '^([^#][^=]+)=(.*)$') {
        $name = $matches[1].Trim()
        $value = $matches[2].Trim()
        Set-Item -Path "env:$name" -Value $value
    }
}

# Validate required environment variables
$requiredVars = @("DB_USER", "DB_PASSWORD", "JWT_SECRET", "GRAFANA_PASSWORD")
foreach ($var in $requiredVars) {
    $value = [Environment]::GetEnvironmentVariable($var)
    if ([string]::IsNullOrEmpty($value) -or $value.Contains("CHANGE_ME")) {
        Write-Host "Error: $var is not set or contains default value" -ForegroundColor Red
        Write-Host "Please update .env.production with secure values" -ForegroundColor Yellow
        exit 1
    }
}

Write-Host "Environment variables validated" -ForegroundColor Green

# Determine image to pull
$imageRegistry = if ($env:DOCKER_REGISTRY) { $env:DOCKER_REGISTRY } else { if ($env:GITHUB_REPOSITORY_OWNER) { "ghcr.io/$($env:GITHUB_REPOSITORY_OWNER)/suproxy-backend" } else { "ghcr.io/tuncay005-png/suproxy-backend" } }
$imageTag = if ($env:VERSION) { $env:VERSION } else { "latest" }
$fullImage = "${imageRegistry}:${imageTag}"

# Pull Docker image from registry
Write-Host "Pulling Docker image: $fullImage" -ForegroundColor Green
docker pull $fullImage

if ($LASTEXITCODE -ne 0) {
    Write-Host "Docker pull failed" -ForegroundColor Red
    Write-Host "Make sure the image exists in GHCR: $fullImage" -ForegroundColor Yellow
    exit 1
}

Write-Host "Docker image pulled successfully: $fullImage" -ForegroundColor Green

# Stop existing containers
Write-Host "Stopping existing containers..." -ForegroundColor Yellow
docker-compose -f docker-compose.production.yml down

# Start services
Write-Host "Starting services..." -ForegroundColor Green
docker-compose -f docker-compose.production.yml up -d

# Wait for health checks
Write-Host "Waiting for services to be healthy..." -ForegroundColor Yellow
Start-Sleep -Seconds 10

# Check if API is healthy
$maxRetries = 30
$retryCount = 0
$apiPort = if ($env:API_PORT) { $env:API_PORT } else { "8080" }

while ($retryCount -lt $maxRetries) {
    try {
        $response = Invoke-WebRequest -Uri "http://localhost:${apiPort}/health" -UseBasicParsing -TimeoutSec 2
        if ($response.StatusCode -eq 200) {
            Write-Host "API is healthy!" -ForegroundColor Green
            break
        }
    }
    catch {
        $retryCount++
        Write-Host "Waiting for API... ($retryCount/$maxRetries)" -ForegroundColor Yellow
        Start-Sleep -Seconds 2
    }
}

if ($retryCount -eq $maxRetries) {
    Write-Host "API health check failed" -ForegroundColor Red
    Write-Host "Checking logs..." -ForegroundColor Yellow
    docker-compose -f docker-compose.production.yml logs api
    exit 1
}

# Show running containers
Write-Host "========================================" -ForegroundColor Green
Write-Host "Deployment Successful!" -ForegroundColor Green
Write-Host "========================================" -ForegroundColor Green
docker-compose -f docker-compose.production.yml ps

$prometheusPort = if ($env:PROMETHEUS_PORT) { $env:PROMETHEUS_PORT } else { "9090" }
$grafanaPort = if ($env:GRAFANA_PORT) { $env:GRAFANA_PORT } else { "3000" }

Write-Host ""
Write-Host "Service URLs:" -ForegroundColor Green
Write-Host "API: http://localhost:${apiPort}"
Write-Host "Prometheus: http://localhost:${prometheusPort}"
Write-Host "Grafana: http://localhost:${grafanaPort}"
Write-Host ""
Write-Host "To view logs: docker-compose -f docker-compose.production.yml logs -f" -ForegroundColor Yellow
