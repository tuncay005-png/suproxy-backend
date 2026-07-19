#!/bin/bash

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${GREEN}========================================${NC}"
echo -e "${GREEN}SuProxy Backend Deployment Script${NC}"
echo -e "${GREEN}========================================${NC}"

# Check if .env.production exists
if [ ! -f ".env.production" ]; then
    echo -e "${RED}Error: .env.production file not found${NC}"
    echo -e "${YELLOW}Please create .env.production from .env.example${NC}"
    exit 1
fi

# Load environment variables
set -a
source .env.production
set +a

# Validate required environment variables
REQUIRED_VARS=("DB_USER" "DB_PASSWORD" "JWT_SECRET" "GRAFANA_PASSWORD")
for VAR in "${REQUIRED_VARS[@]}"; do
    if [ -z "${!VAR}" ] || [[ "${!VAR}" == *"CHANGE_ME"* ]]; then
        echo -e "${RED}Error: $VAR is not set or contains default value${NC}"
        echo -e "${YELLOW}Please update .env.production with secure values${NC}"
        exit 1
    fi
done

echo -e "${GREEN}Environment variables validated${NC}"

# Determine image to pull
IMAGE_REGISTRY="${DOCKER_REGISTRY:-ghcr.io/${GITHUB_REPOSITORY_OWNER:-tuncay005-png}/suproxy-backend}"
IMAGE_TAG="${VERSION:-latest}"
FULL_IMAGE="${IMAGE_REGISTRY}:${IMAGE_TAG}"

# Pull Docker image from registry
echo -e "${GREEN}Pulling Docker image: ${FULL_IMAGE}${NC}"
docker pull "${FULL_IMAGE}"

if [ $? -ne 0 ]; then
    echo -e "${RED}Docker pull failed${NC}"
    echo -e "${YELLOW}Make sure the image exists in GHCR: ${FULL_IMAGE}${NC}"
    exit 1
fi

echo -e "${GREEN}Docker image pulled successfully: ${FULL_IMAGE}${NC}"

# Stop existing containers
echo -e "${YELLOW}Stopping existing containers...${NC}"
docker-compose -f docker-compose.production.yml down

# Start services
echo -e "${GREEN}Starting services...${NC}"
docker-compose -f docker-compose.production.yml up -d

# Wait for health checks
echo -e "${YELLOW}Waiting for services to be healthy...${NC}"
sleep 10

# Check if API is healthy
MAX_RETRIES=30
RETRY_COUNT=0
while [ $RETRY_COUNT -lt $MAX_RETRIES ]; do
    if curl -f http://localhost:${API_PORT:-8080}/health > /dev/null 2>&1; then
        echo -e "${GREEN}API is healthy!${NC}"
        break
    fi
    RETRY_COUNT=$((RETRY_COUNT+1))
    echo -e "${YELLOW}Waiting for API... ($RETRY_COUNT/$MAX_RETRIES)${NC}"
    sleep 2
done

if [ $RETRY_COUNT -eq $MAX_RETRIES ]; then
    echo -e "${RED}API health check failed${NC}"
    echo -e "${YELLOW}Checking logs...${NC}"
    docker-compose -f docker-compose.production.yml logs api
    exit 1
fi

# Show running containers
echo -e "${GREEN}========================================${NC}"
echo -e "${GREEN}Deployment Successful!${NC}"
echo -e "${GREEN}========================================${NC}"
docker-compose -f docker-compose.production.yml ps

echo -e "\n${GREEN}Service URLs:${NC}"
echo -e "API: http://localhost:${API_PORT:-8080}"
echo -e "Prometheus: http://localhost:${PROMETHEUS_PORT:-9090}"
echo -e "Grafana: http://localhost:${GRAFANA_PORT:-3000}"
echo -e "\n${YELLOW}To view logs: docker-compose -f docker-compose.production.yml logs -f${NC}"
