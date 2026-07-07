# Integration Tests

This directory contains integration tests for SuProxy Backend.

## Quick Start

### Prerequisites

1. PostgreSQL running (via docker-compose)
2. Test database created

### Setup Test Database

```bash
# Start PostgreSQL
docker-compose up -d postgres

# Create test database
make test-db-setup
```

### Run Integration Tests

```bash
# All integration tests
make test-integration

# Specific test
INTEGRATION_TEST=true go test -v ./test/integration -run TestUserRepository

# With coverage
make test-coverage-integration
```

## Writing Integration Tests

### Basic Structure

```go
package integration_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/suproxy/backend/internal/infrastructure/testutil"
)

func TestMyFeature_Integration(t *testing.T) {
	// Skip if not running integration tests
	if !testutil.IsIntegrationTest() {
		t.Skip("Skipping integration test")
	}

	// Setup
	app := testutil.NewTestApp(t)
	defer app.Cleanup()
	defer app.CleanupTables()

	ctx := context.Background()

	// Test logic
	// Use app.Container for repositories
	// Use testutil fixtures for test data
	
	// Assertions
	assert.NoError(t, err)
	require.NotNil(t, result)
}
```

### Using Fixtures

```go
// Create test user
user, err := testutil.CreateTestUserWithDefaults()
require.NoError(t, err)

// Save to database
err = app.Container.UserRepository.Create(ctx, user)
require.NoError(t, err)

// Create admin user
admin, err := testutil.CreateTestAdminUser()
require.NoError(t, err)
```

### Testing with Authentication

```go
authHelper := testutil.NewAuthHelper(app.JWT, t)

// Create authenticated user
user, accessToken, refreshToken := authHelper.CreateAuthenticatedUser(
	app.Container.UserRepository,
)

// Use tokens in tests
headers := testutil.AuthHeader(accessToken)
```

### HTTP API Testing

```go
httpCtx := testutil.NewHTTPTestContext(t)

// Setup routes
// router := setupYourRoutes(app)
// httpCtx.Router = router

// Make request
resp := httpCtx.GET("/api/v1/users", headers)

// Assert response
httpCtx.AssertStatusCode(200)

var result map[string]interface{}
httpCtx.GetResponseJSON(&result)
```

## Test Organization

Each integration test file should:
1. Test a specific feature or component
2. Be independent (can run in isolation)
3. Clean up after itself
4. Use meaningful test names

## Best Practices

1. **Always check integration flag**:
   ```go
   if !testutil.IsIntegrationTest() {
       t.Skip("Skipping integration test")
   }
   ```

2. **Always cleanup**:
   ```go
   defer app.Cleanup()
   defer app.CleanupTables()
   ```

3. **Use fixtures**:
   ```go
   user, _ := testutil.CreateTestUserWithDefaults()
   ```

4. **Test isolation**:
   - Each test should be independent
   - Don't rely on test execution order
   - Clean database state between tests

5. **Meaningful names**:
   ```go
   func TestUserRepository_CreateAndFind_Success(t *testing.T)
   ```

## Environment Variables

Set these in your environment or `.env.test`:

```bash
INTEGRATION_TEST=true
TEST_DB_HOST=localhost
TEST_DB_PORT=5432
TEST_DB_USER=suproxy_test
TEST_DB_PASSWORD=suproxy_test
TEST_DB_NAME=suproxy_test
```

## Troubleshooting

### Database Connection Error

```bash
# Ensure PostgreSQL is running
docker-compose ps

# Recreate test database
make test-db-teardown
make test-db-setup
```

### Tests Skipped

Make sure `INTEGRATION_TEST=true` is set:
```bash
INTEGRATION_TEST=true go test ./test/integration
```

### Cleanup Issues

Manually truncate tables:
```bash
docker-compose exec postgres psql -U suproxy_test -d suproxy_test -c "
TRUNCATE TABLE sessions, audit_logs, clients, reality_configs, 
inbounds, xray_instances, devices, subscriptions, users, nodes CASCADE;
"
```

## Documentation

- [Testing Guide](../docs/testing.md) - Comprehensive testing documentation
- [Test Utilities](../internal/infrastructure/testutil/README.md) - Test utility documentation

## Examples

See existing integration tests for examples:
- `user_test.go` - User repository and authentication tests

Add more integration tests following the same patterns!
