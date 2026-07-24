#!/bin/bash

# Source the deploy script to get the functions
source scripts/deploy.sh

echo "=========================================="
echo "Testing Error Classification"
echo "=========================================="

# Test network errors
echo "Test 1: i/o timeout -> should be 'network'"
result=$(classify_docker_error "Error: i/o timeout")
echo "Result: $result"
[ "$result" = "network" ] && echo "✓ PASS" || echo "✗ FAIL"

echo ""
echo "Test 2: TLS handshake timeout -> should be 'network'"
result=$(classify_docker_error "Error: TLS handshake timeout")
echo "Result: $result"
[ "$result" = "network" ] && echo "✓ PASS" || echo "✗ FAIL"

echo ""
echo "Test 3: connection reset -> should be 'network'"
result=$(classify_docker_error "Error: connection reset")
echo "Result: $result"
[ "$result" = "network" ] && echo "✓ PASS" || echo "✗ FAIL"

echo ""
echo "Test 4: context deadline exceeded -> should be 'network'"
result=$(classify_docker_error "Error: context deadline exceeded")
echo "Result: $result"
[ "$result" = "network" ] && echo "✓ PASS" || echo "✗ FAIL"

echo ""
echo "Test 5: EOF -> should be 'network'"
result=$(classify_docker_error "Error: unexpected EOF")
echo "Result: $result"
[ "$result" = "network" ] && echo "✓ PASS" || echo "✗ FAIL"

# Test permanent errors
echo ""
echo "Test 6: manifest unknown -> should be 'permanent'"
result=$(classify_docker_error "Error: manifest unknown")
echo "Result: $result"
[ "$result" = "permanent" ] && echo "✓ PASS" || echo "✗ FAIL"

echo ""
echo "Test 7: unauthorized -> should be 'permanent'"
result=$(classify_docker_error "Error: unauthorized")
echo "Result: $result"
[ "$result" = "permanent" ] && echo "✓ PASS" || echo "✗ FAIL"

echo ""
echo "Test 8: denied -> should be 'permanent'"
result=$(classify_docker_error "Error: access denied")
echo "Result: $result"
[ "$result" = "permanent" ] && echo "✓ PASS" || echo "✗ FAIL"

echo ""
echo "Test 9: pull access denied -> should be 'permanent'"
result=$(classify_docker_error "Error: pull access denied")
echo "Result: $result"
[ "$result" = "permanent" ] && echo "✓ PASS" || echo "✗ FAIL"

# Test unknown error (should default to network)
echo ""
echo "Test 10: unknown error -> should default to 'network'"
result=$(classify_docker_error "Error: something completely unknown")
echo "Result: $result"
[ "$result" = "network" ] && echo "✓ PASS" || echo "✗ FAIL"

# Test case insensitivity
echo ""
echo "Test 11: MANIFEST UNKNOWN (uppercase) -> should be 'permanent'"
result=$(classify_docker_error "Error: MANIFEST UNKNOWN")
echo "Result: $result"
[ "$result" = "permanent" ] && echo "✓ PASS" || echo "✗ FAIL"

echo ""
echo "=========================================="
echo "Testing Exponential Backoff Calculation"
echo "=========================================="

BACKOFF_BASE_DELAY=5

# Test backoff calculations
for ATTEMPT in 1 2 3 4; do
    if [ $ATTEMPT -eq 1 ]; then
        BACKOFF_DELAY=0
    else
        BACKOFF_DELAY=$((BACKOFF_BASE_DELAY * (2 ** (ATTEMPT - 2))))
    fi
    echo "Attempt $ATTEMPT: Backoff = ${BACKOFF_DELAY}s"
done

echo ""
echo "Expected:"
echo "Attempt 1: Backoff = 0s"
echo "Attempt 2: Backoff = 5s"
echo "Attempt 3: Backoff = 10s"
echo "Attempt 4: Backoff = 20s"

echo ""
echo "=========================================="
echo "All tests completed!"
echo "=========================================="
