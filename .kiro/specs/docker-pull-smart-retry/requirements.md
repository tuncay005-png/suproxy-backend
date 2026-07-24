# Requirements Document

## Introduction

This document specifies requirements for implementing a production-grade, smart retry mechanism for Docker pull operations in the deployment script (`scripts/deploy.sh`). The current retry mechanism retries ALL Docker pull failures blindly, including permanent errors like authentication failures or missing manifests. This enhancement adds intelligent error classification to retry only transient network errors and fail fast on permanent errors, reducing deployment time and providing clearer feedback.

## Glossary

- **Deploy_Script**: The bash script located at `scripts/deploy.sh` responsible for deploying the application
- **Docker_Pull**: The `docker pull` command that downloads container images from GHCR
- **GHCR**: GitHub Container Registry (ghcr.io)
- **Network_Error**: Transient network failures that should be retried - specifically: i/o timeout, TLS handshake timeout, connection reset, context deadline exceeded, EOF
- **Permanent_Error**: Non-recoverable errors that should not be retried - specifically: manifest unknown, unauthorized, denied, pull access denied
- **Error_Classifier**: A bash function that analyzes docker pull error output to determine error type
- **Retry_Handler**: A bash function that performs docker pull with exponential backoff
- **Exponential_Backoff**: A retry strategy where delay increases exponentially between attempts (5s → 10s → 20s)
- **Attempt**: A single execution of the docker pull command - total of 4 attempts (Attempt 1 = initial pull with no delay, Attempt 2 = retry after 5s, Attempt 3 = retry after 10s, Attempt 4 = retry after 20s)

## Requirements

### Requirement 1: Smart Error Classification

**User Story:** As a DevOps engineer, I want the deployment script to distinguish between transient and permanent Docker pull errors, so that I get fast feedback on configuration issues and automatic recovery from network glitches.

#### Acceptance Criteria

1. WHEN Docker pull fails with an error message, THE Error_Classifier SHALL analyze the stderr output to determine error type
2. WHEN the error message contains "i/o timeout", THE Error_Classifier SHALL classify it as a Network_Error
3. WHEN the error message contains "TLS handshake timeout", THE Error_Classifier SHALL classify it as a Network_Error
4. WHEN the error message contains "connection reset", THE Error_Classifier SHALL classify it as a Network_Error
5. WHEN the error message contains "context deadline exceeded", THE Error_Classifier SHALL classify it as a Network_Error
6. WHEN the error message contains "EOF", THE Error_Classifier SHALL classify it as a Network_Error
7. WHEN the error message contains "manifest unknown", THE Error_Classifier SHALL classify it as a Permanent_Error
8. WHEN the error message contains "unauthorized", THE Error_Classifier SHALL classify it as a Permanent_Error
9. WHEN the error message contains "denied", THE Error_Classifier SHALL classify it as a Permanent_Error
10. WHEN the error message contains "pull access denied", THE Error_Classifier SHALL classify it as a Permanent_Error
11. WHEN the error message does not match any known pattern, THE Error_Classifier SHALL classify it as a Network_Error by default

### Requirement 2: Fail Fast on Permanent Errors

**User Story:** As a DevOps engineer, I want deployment to fail immediately on authentication or configuration errors, so that I can fix the problem quickly without waiting for unnecessary retry attempts.

#### Acceptance Criteria

1. WHEN Docker_Pull fails AND Error_Classifier determines a Permanent_Error, THE Deploy_Script SHALL log the error with RED color formatting
2. WHEN a Permanent_Error is detected, THE Deploy_Script SHALL log "Permanent error detected." with RED color
3. WHEN a Permanent_Error is detected, THE Deploy_Script SHALL log "Retry skipped." with RED color
4. WHEN a Permanent_Error is detected, THE Deploy_Script SHALL log "Deployment aborted." with RED color
5. WHEN a Permanent_Error is detected, THE Deploy_Script SHALL exit with code 1 immediately
6. WHEN a Permanent_Error is detected, THE Deploy_Script SHALL NOT sleep or wait before exiting

### Requirement 3: Exponential Backoff for Network Errors

**User Story:** As a DevOps engineer, I want transient network errors to be retried with increasing delays, so that temporary GHCR outages don't cause deployment failures.

#### Acceptance Criteria

1. THE Retry_Handler SHALL attempt Docker_Pull a maximum of 4 times total (Attempt 1, Attempt 2, Attempt 3, Attempt 4)
2. WHEN Attempt 1 fails with a Network_Error, THE Retry_Handler SHALL wait 5 seconds before Attempt 2
3. WHEN Attempt 2 fails with a Network_Error, THE Retry_Handler SHALL wait 10 seconds before Attempt 3
4. WHEN Attempt 3 fails with a Network_Error, THE Retry_Handler SHALL wait 20 seconds before Attempt 4
5. WHEN Docker_Pull succeeds on any attempt, THE Retry_Handler SHALL break the retry loop immediately
6. WHEN all 4 attempts fail, THE Retry_Handler SHALL exit with code 1

### Requirement 4: Configurable Retry Parameters

**User Story:** As a DevOps engineer, I want to easily adjust retry behavior, so that I can tune the mechanism for different deployment environments.

#### Acceptance Criteria

1. THE Deploy_Script SHALL define MAX_PULL_RETRIES as a configurable variable at the top of the script with default value 3
2. THE Deploy_Script SHALL define BACKOFF_BASE_DELAY as a configurable variable at the top of the script with default value 5
3. WHEN calculating backoff delay for attempt N, THE Retry_Handler SHALL use formula: BACKOFF_BASE_DELAY * (2^(N-1))
4. THE Deploy_Script SHALL use these variables throughout the retry logic consistently

### Requirement 5: Detailed Logging

**User Story:** As a DevOps engineer, I want detailed logs during retry operations, so that I can diagnose deployment issues and understand retry behavior.

#### Acceptance Criteria

1. WHEN starting a Docker_Pull attempt, THE Retry_Handler SHALL log "Attempt N/4" where N is the current attempt number
2. WHEN starting a Docker_Pull attempt, THE Retry_Handler SHALL log "docker pull <image>" showing the full docker pull command
3. WHEN Docker_Pull fails, THE Retry_Handler SHALL log the error type detected (Network_Error or Permanent_Error) with appropriate messaging
4. WHEN retrying after a Network_Error, THE Retry_Handler SHALL log "Network timeout detected." with YELLOW color
5. WHEN retrying after a Network_Error, THE Retry_Handler SHALL log "Retrying in Ns..." with YELLOW color where N is the backoff delay
6. WHEN Docker_Pull succeeds, THE Retry_Handler SHALL log "Docker image pulled successfully: <image>" with GREEN color
7. WHEN all retries are exhausted, THE Retry_Handler SHALL log "Docker pull failed after N attempts" with RED color
8. THE Deploy_Script SHALL preserve existing color formatting constants (RED, GREEN, YELLOW, NC)

### Requirement 6: Backward Compatibility

**User Story:** As a DevOps engineer, I want the enhanced retry mechanism to preserve existing deployment behavior, so that other deployment logic remains unaffected.

#### Acceptance Criteria

1. THE Deploy_Script SHALL NOT modify VERSION environment variable handling
2. THE Deploy_Script SHALL NOT modify DOCKER_REGISTRY environment variable handling
3. THE Deploy_Script SHALL NOT modify GITHUB_REPOSITORY_OWNER environment variable handling
4. THE Deploy_Script SHALL NOT modify docker-compose commands or health check logic
5. THE Deploy_Script SHALL maintain the same exit codes as the current implementation (0 for success, 1 for failure)
6. WHEN Docker_Pull succeeds on first attempt, THE Retry_Handler SHALL complete without any backoff delays

### Requirement 7: Pure Bash Implementation

**User Story:** As a DevOps engineer, I want the retry mechanism implemented in pure bash, so that deployments don't depend on external tools or libraries.

#### Acceptance Criteria

1. THE Error_Classifier SHALL be implemented as a bash function without external dependencies
2. THE Retry_Handler SHALL be implemented as a bash function without external dependencies
3. THE Deploy_Script SHALL use bash built-in commands for string matching (grep, if statements)
4. THE Deploy_Script SHALL use bash arithmetic for calculating backoff delays
5. THE Deploy_Script SHALL NOT require installation of additional packages or tools

## Implementation Notes

### Current Retry Mechanism (Lines 47-73)

The existing retry logic in `deploy.sh` looks like this:

```bash
MAX_PULL_RETRIES=3
PULL_RETRY=0
PULL_SUCCESS=false

while [ $PULL_RETRY -lt $MAX_PULL_RETRIES ]; do
    if docker pull "${FULL_IMAGE}"; then
        PULL_SUCCESS=true
        break
    fi
    
    PULL_RETRY=$((PULL_RETRY+1))
    
    if [ $PULL_RETRY -lt $MAX_PULL_RETRIES ]; then
        BACKOFF_DELAY=$((PULL_RETRY * 10))
        echo -e "${YELLOW}Pull attempt $PULL_RETRY failed, retrying in ${BACKOFF_DELAY}s... ($PULL_RETRY/$MAX_PULL_RETRIES)${NC}"
        sleep $BACKOFF_DELAY
    fi
done

if [ "$PULL_SUCCESS" = false ]; then
    echo -e "${RED}Docker pull failed after $MAX_PULL_RETRIES attempts${NC}"
    exit 1
fi
```

**Problems with Current Implementation:**
- Retries ALL errors (authentication, manifest not found, network timeouts)
- Linear backoff (10s, 20s, 30s) instead of exponential (5s, 10s, 20s)
- No error classification or logging of error types
- Wastes time retrying permanent errors

### Target Error Messages

Based on production logs and Docker documentation, here are the error patterns to detect:

**Network Errors (RETRY):**
- `i/o timeout`
- `TLS handshake timeout`
- `connection reset`
- `context deadline exceeded`
- `EOF`

**Permanent Errors (FAIL FAST):**
- `manifest unknown` (image doesn't exist)
- `unauthorized` (authentication failed)
- `denied` (access denied)
- `pull access denied` (access denied)

### Example Enhanced Output

**Scenario 1: Network error with retry**
```
Attempt 1/4
docker pull ghcr.io/tuncay005-png/suproxy-backend:latest
Network timeout detected.
Retrying in 5 seconds...
Attempt 2/4
docker pull ghcr.io/tuncay005-png/suproxy-backend:latest
Docker image pulled successfully: ghcr.io/tuncay005-png/suproxy-backend:latest
```

**Scenario 2: Permanent error - fail fast**
```
Attempt 1/4
docker pull ghcr.io/tuncay005-png/suproxy-backend:latest
Permanent error detected.
Retry skipped.
Deployment aborted.
```

**Scenario 3: All retries exhausted**
```
Attempt 1/4
docker pull ghcr.io/tuncay005-png/suproxy-backend:latest
Network timeout detected.
Retrying in 5 seconds...
Attempt 2/4
docker pull ghcr.io/tuncay005-png/suproxy-backend:latest
Network timeout detected.
Retrying in 10 seconds...
Attempt 3/4
docker pull ghcr.io/tuncay005-png/suproxy-backend:latest
Network timeout detected.
Retrying in 20 seconds...
Attempt 4/4
docker pull ghcr.io/tuncay005-png/suproxy-backend:latest
Docker pull failed after 4 attempts
Make sure the image exists in GHCR: ghcr.io/tuncay005-png/suproxy-backend:latest
```
