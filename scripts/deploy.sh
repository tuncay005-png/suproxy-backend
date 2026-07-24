#!/bin/bash

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Configuration constants for Docker pull retry mechanism
MAX_PULL_RETRIES=3
BACKOFF_BASE_DELAY=5

# Function: classify_docker_error
# Classifies Docker pull errors as "network" (retryable) or "permanent" (fail-fast)
# Arguments:
#   $1 - Error output from docker pull command
# Returns:
#   Echoes "network" or "permanent"
classify_docker_error() {
    local error_output="$1"
    
    # Check for permanent errors FIRST (CRITICAL: check permanent before network)
    if echo "$error_output" | grep -qi "manifest unknown"; then
        echo "permanent"
        return
    fi
    
    if echo "$error_output" | grep -qi "unauthorized"; then
        echo "permanent"
        return
    fi
    
    if echo "$error_output" | grep -qi "denied"; then
        echo "permanent"
        return
    fi
    
    if echo "$error_output" | grep -qi "pull access denied"; then
        echo "permanent"
        return
    fi
    
    # Check for network errors (case-insensitive)
    if echo "$error_output" | grep -qi "i/o timeout"; then
        echo "network"
        return
    fi
    
    if echo "$error_output" | grep -qi "TLS handshake timeout"; then
        echo "network"
        return
    fi
    
    if echo "$error_output" | grep -qi "connection reset"; then
        echo "network"
        return
    fi
    
    if echo "$error_output" | grep -qi "context deadline exceeded"; then
        echo "network"
        return
    fi
    
    if echo "$error_output" | grep -qi "EOF"; then
        echo "network"
        return
    fi
    
    # Default to "network" for unknown errors (conservative default)
    echo "network"
}

# Function: pull_docker_image_with_retry
# Pulls Docker image with smart retry logic and exponential backoff
# Arguments:
#   $1 - Full image name (e.g., ghcr.io/org/repo:tag)
# Returns:
#   Exit code 0 on success, 1 on failure
pull_docker_image_with_retry() {
    local image="$1"
    
    # Loop through attempts 1 to 4 (1 initial + 3 retries)
    for ATTEMPT in {1..4}; do
        # Log attempt number and command
        echo "Attempt $ATTEMPT/4"
        echo "docker pull ${image}"
        
        # Execute docker pull and capture stderr
        ERROR_OUTPUT=$(docker pull "${image}" 2>&1)
        PULL_EXIT_CODE=$?
        
        # Check if pull succeeded
        if [ $PULL_EXIT_CODE -eq 0 ]; then
            echo -e "${GREEN}Docker image pulled successfully: ${image}${NC}"
            return 0
        fi
        
        # Pull failed - classify the error
        ERROR_TYPE=$(classify_docker_error "$ERROR_OUTPUT")
        
        # Log error classification
        echo -e "${YELLOW}Error type: ${ERROR_TYPE}${NC}"
        
        # Handle permanent errors - fail immediately
        if [ "$ERROR_TYPE" = "permanent" ]; then
            echo -e "${RED}Permanent error detected.${NC}"
            echo -e "${RED}Retry skipped.${NC}"
            echo -e "${RED}Deployment aborted.${NC}"
            return 1
        fi
        
        # Handle network errors - retry with backoff if attempts remain
        if [ "$ERROR_TYPE" = "network" ]; then
            # Check if this was the last attempt
            if [ $ATTEMPT -eq 4 ]; then
                echo -e "${RED}Docker pull failed after 4 attempts${NC}"
                echo -e "${YELLOW}Make sure the image exists in GHCR: ${image}${NC}"
                return 1
            fi
            
            # Calculate exponential backoff delay
            # Attempt 1: no delay (not reached since success exits above)
            # Attempt 2: 5 * 2^0 = 5s
            # Attempt 3: 5 * 2^1 = 10s
            # Attempt 4: 5 * 2^2 = 20s
            BACKOFF_DELAY=$((BACKOFF_BASE_DELAY * (2 ** (ATTEMPT - 2))))
            
            # Log network error and backoff
            echo -e "${YELLOW}Network timeout detected.${NC}"
            echo -e "${YELLOW}Retrying in ${BACKOFF_DELAY}s...${NC}"
            
            # Sleep for the calculated delay
            sleep $BACKOFF_DELAY
        fi
    done
    
    # Should not reach here due to return 1 in loop, but safety fallback
    return 1
}

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

# Pull Docker image from registry with smart retry mechanism
if ! pull_docker_image_with_retry "${FULL_IMAGE}"; then
    echo -e "${RED}Failed to pull Docker image. Deployment aborted.${NC}"
    exit 1
fi

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
