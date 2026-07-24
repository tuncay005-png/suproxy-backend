# Implementation Plan: Smart Docker Pull Retry Mechanism

## Overview

This implementation plan converts the smart Docker pull retry mechanism design into actionable coding tasks. The implementation will replace the existing blind retry logic in `scripts/deploy.sh` (lines 47-73) with intelligent error classification and exponential backoff.

**Retry Structure:**
- **Total Attempts**: 4 (1 initial + 3 retries)
- **Attempt 1**: Initial pull with no delay
- **Attempt 2**: First retry after 5-second delay
- **Attempt 3**: Second retry after 10-second delay  
- **Attempt 4**: Third retry after 20-second delay

**Error Classification:**

**Network errors (retry):**
- i/o timeout
- TLS handshake timeout
- connection reset
- context deadline exceeded
- EOF

**Permanent errors (fail immediately):**
- manifest unknown
- unauthorized
- denied
- pull access denied

**Logging Format:**

Retry attempts:
```
Attempt 1/4
docker pull ...
Network timeout detected.
Retrying in 5 seconds...
Attempt 2/4
...
```

Permanent errors:
```
Permanent error detected.
Retry skipped.
Deployment aborted.
```

## Tasks

- [x] 1. Set up configuration constants and function structure
  - Add `MAX_PULL_RETRIES=3` and `BACKOFF_BASE_DELAY=5` configuration constants at the top of `scripts/deploy.sh` (after color constants)
  - Define function placeholders for `classify_docker_error()` and `pull_docker_image_with_retry()` after the configuration section
  - Ensure proper bash function syntax with local variable declarations
  - _Requirements: 4.1, 4.2, 7.1, 7.2_

- [x] 2. Implement error classifier function
  - [x] 2.1 Create `classify_docker_error()` function with error pattern matching
    - Accept error output as parameter: `local error_output="$1"`
    - Use `grep -qi` (case-insensitive) to check for permanent error patterns: "manifest unknown", "unauthorized", "denied", "pull access denied"
    - Use `grep -qi` to check for network error patterns: "i/o timeout", "TLS handshake timeout", "connection reset", "context deadline exceeded", "EOF"
    - Return "permanent" if permanent pattern matches, otherwise return "network" (conservative default)
    - Use `echo` to return the classification result
    - _Requirements: 1.1, 1.2, 1.3, 1.4, 1.5, 1.6, 1.7, 1.8, 1.9, 1.10, 1.11, 7.3_

  - [ ]* 2.2 Write property test for network error pattern classification
    - **Property 1: Network Error Patterns Are Classified as Retryable**
    - **Validates: Requirements 1.2, 1.3, 1.4, 1.5, 1.6**
    - Create test script that generates 100 random error messages containing network error patterns
    - Verify each message is classified as "network"
    - Test all network patterns: "i/o timeout", "TLS handshake timeout", "connection reset", "context deadline exceeded", "EOF"

  - [ ]* 2.3 Write property test for permanent error pattern classification
    - **Property 2: Permanent Error Patterns Are Classified as Non-Retryable**
    - **Validates: Requirements 1.7, 1.8, 1.9, 1.10**
    - Create test script that generates 100 random error messages containing permanent error patterns
    - Verify each message is classified as "permanent"
    - Test all permanent patterns: "manifest unknown", "unauthorized", "denied", "pull access denied"

  - [ ]* 2.4 Write property test for case-insensitivity
    - **Property 5: Error Classification Is Case-Insensitive**
    - **Validates: Requirements 1.1**
    - Create test script that generates known error patterns with random case variations
    - Verify classification result is consistent across case variations (uppercase, lowercase, mixed case)
    - Test with 100 iterations

  - [ ]* 2.5 Write property test for unknown error default behavior
    - **Property 3: Unknown Error Patterns Default to Network Classification**
    - **Validates: Requirements 1.11**
    - Create test script that generates 100 random error messages without any known patterns
    - Verify all unknown errors default to "network" classification

- [x] 3. Implement retry handler function with exponential backoff
  - [x] 3.1 Create `pull_docker_image_with_retry()` function skeleton
    - Accept image name as parameter: `local image="$1"`
    - Initialize attempt counter and loop structure (1 to 4 attempts)
    - Set up stderr capture mechanism using `2>&1` redirection
    - _Requirements: 3.1, 7.4_

  - [x] 3.2 Implement docker pull execution and error capture
    - Execute `docker pull "${image}"` and capture both stdout and stderr
    - Store exit code in variable for checking success/failure
    - On success (exit code 0), log success message with GREEN color and return 0
    - On failure, capture stderr output for classification
    - _Requirements: 3.5, 5.4_

  - [x] 3.3 Implement error classification and branching logic
    - Call `classify_docker_error` function with captured stderr
    - If classification is "permanent", log "Permanent error detected." with RED color
    - If classification is "permanent", log "Retry skipped." with RED color
    - If classification is "permanent", log "Deployment aborted." with RED color and exit 1 immediately
    - If classification is "network" and attempts remain, proceed to backoff calculation
    - If classification is "network" and no attempts remain, log exhausted retries message with RED color and exit 1
    - _Requirements: 1.1, 2.1, 2.2, 2.3, 2.4, 2.5, 2.6, 5.2, 5.3, 5.7_

  - [x] 3.4 Implement exponential backoff delay calculation and logging
    - Calculate backoff delay using formula: `BACKOFF_DELAY=$((BACKOFF_BASE_DELAY * (2 ** (ATTEMPT - 2))))` for ATTEMPT > 1
    - For attempt 1: no delay (initial attempt)
    - For attempt 2: 5 seconds delay (5 * 2^0)
    - For attempt 3: 10 seconds delay (5 * 2^1)
    - For attempt 4: 20 seconds delay (5 * 2^2)
    - Log "Attempt N/4" with current attempt number
    - Log "docker pull <image>" showing the full command
    - Log "Network timeout detected." with YELLOW color before sleeping
    - Log "Retrying in Ns..." with YELLOW color where N is the backoff delay
    - Use bash `sleep` command with calculated delay
    - _Requirements: 3.1, 3.2, 3.3, 3.4, 4.3, 5.1, 5.2, 5.4, 5.5, 7.4_

  - [ ]* 3.5 Write property test for exponential backoff calculation
    - **Property 4: Exponential Backoff Delay Calculation**
    - **Validates: Requirements 3.1, 3.2, 3.3, 3.4**
    - Create test script that verifies backoff formula for various BACKOFF_BASE_DELAY values
    - Test with 100 random base delay values (1-30 seconds)
    - Verify backoff(N) = BACKOFF_BASE_DELAY * (2^(N-2)) for attempts 2, 3, 4
    - Verify backoff(1) = 0 (no delay on first attempt)

- [x] 4. Checkpoint - Ensure functions are tested and working
  - Ensure all tests pass, ask the user if questions arise.

- [x] 5. Integrate retry handler into deploy.sh
  - [x] 5.1 Replace existing retry logic with new function call
    - Locate lines 47-73 in `scripts/deploy.sh` (current retry block)
    - Remove the old while loop, PULL_RETRY, PULL_SUCCESS variables
    - Replace with single function call: `pull_docker_image_with_retry "${FULL_IMAGE}"`
    - Verify the script still exits with correct codes (0 on success, 1 on failure)
    - _Requirements: 6.5, 6.6_

  - [x] 5.2 Verify backward compatibility with environment variables
    - Ensure VERSION, DOCKER_REGISTRY, GITHUB_REPOSITORY_OWNER handling remains unchanged
    - Verify FULL_IMAGE construction logic is not modified
    - Verify docker-compose commands after pull operation remain unchanged
    - Verify health check logic remains unchanged
    - _Requirements: 6.1, 6.2, 6.3, 6.4_

  - [x] 5.3 Preserve color formatting constants
    - Verify RED, GREEN, YELLOW, NC color constants are still defined
    - Ensure new log messages use existing color constants consistently
    - Verify no changes to color usage in non-retry sections
    - _Requirements: 5.6, 6.6_

- [ ]* 6. Write integration tests using BATS
  - [ ]* 6.1 Set up BATS test framework
    - Create `test/deploy_retry_test.bats` file
    - Add BATS shebang and setup/teardown functions
    - Source `scripts/deploy.sh` functions in test setup

  - [ ]* 6.2 Write unit tests for error classifier
    - Test all network error patterns return "network" (i/o timeout, TLS handshake timeout, connection reset, context deadline exceeded, EOF)
    - Test all permanent error patterns return "permanent" (manifest unknown, unauthorized, denied, pull access denied)
    - Test unknown patterns default to "network"
    - Test case-insensitive matching
    - Test multi-line error messages
    - _Requirements: 1.2, 1.3, 1.4, 1.5, 1.6, 1.7, 1.8, 1.9, 1.10, 1.11_

  - [ ]* 6.3 Write unit tests for backoff calculation
    - Test exponential backoff formula for attempts 1-4
    - Test BACKOFF_BASE_DELAY configuration is respected
    - Test boundary conditions (attempt 0, attempt 5)
    - _Requirements: 3.2, 3.3, 3.4, 4.3_

  - [ ]* 6.4 Write integration tests with mocked docker command
    - Test success on first attempt (no retries)
    - Test success on second attempt (one 5s retry)
    - Test permanent error exits immediately without retry
    - Test network error exhausts all retries (5s, 10s, 20s delays)
    - Test correct exit codes (0 on success, 1 on failure)
    - _Requirements: 2.3, 2.4, 3.1, 3.2, 3.3, 3.4, 3.5, 3.6, 6.5_

- [x] 7. Final checkpoint and manual verification
  - Ensure all tests pass, ask the user if questions arise.

## Notes

- Tasks marked with `*` are optional and can be skipped for faster MVP
- Each task references specific requirements for traceability
- The implementation is pure bash with no external dependencies (except BATS for testing, which is optional)
- Property tests validate correctness properties defined in the design document
- Integration tests ensure the mechanism works correctly with actual Docker CLI behavior
- The retry mechanism preserves all existing deployment script functionality
- Configuration constants (MAX_PULL_RETRIES, BACKOFF_BASE_DELAY) are easily tunable at the top of the script
- Exponential backoff: Attempt 1 (0s) → Attempt 2 (5s) → Attempt 3 (10s) → Attempt 4 (20s)
- Total maximum wait time: 35 seconds across all retries
- Conservative default: Unknown errors are classified as "network" to maximize deployment success rate

## Task Dependency Graph

```json
{
  "waves": [
    { "id": 0, "tasks": ["1.1"] },
    { "id": 1, "tasks": ["2.1"] },
    { "id": 2, "tasks": ["2.2", "2.3", "2.4", "2.5", "3.1"] },
    { "id": 3, "tasks": ["3.2"] },
    { "id": 4, "tasks": ["3.3", "3.4"] },
    { "id": 5, "tasks": ["3.5"] },
    { "id": 6, "tasks": ["5.1"] },
    { "id": 7, "tasks": ["5.2", "5.3", "6.1"] },
    { "id": 8, "tasks": ["6.2", "6.3"] },
    { "id": 9, "tasks": ["6.4"] }
  ]
}
```
