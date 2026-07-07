# Test Utilities Package

This package provides comprehensive utilities for integration and unit testing in the SuProxy Backend project.

## Features

- ✅ Test configuration management
- ✅ Test database utilities with cleanup
- ✅ Test application bootstrap
- ✅ Test data fixtures
- ✅ HTTP test utilities
- ✅ Authentication helpers
- ✅ Custom assertions
- ✅ Mock helpers
- ✅ Container support (placeholder for testcontainers-go)

## Usage

### Basic Integration Test

```go
package mypackage_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/suproxy/backend/internal/infrastructure/testutil"
)

func TestMyFeature(t *testing.T) {
	// Skip if not running integration tests
	if !testutil.IsIntegrationTest() {
		t.Skip("Skipping integration test")
	}

	// Setup test application
	app := testutil.NewTestApp(t)
	defer app.Cleanup()
	defer app.CleanupTables()

	ctx := context.Background()

	// Use repositories from container
	userRepo := app.Container.UserRepository

	// Create test data using fixtures
	user, err := testutil.CreateTestUserWithDefaults()
	assert.NoError(t, err)

	// Test your feature
	err = userRepo.Create(ctx, user)
	assert.NoError(t, err)

	// Verify
	found, err := userRepo.FindByID(ctx, user.ID)
	assert.NoError(t, err)
	assert.Equal(t, user.ID, found.ID)
}
```

### HTTP API Integration Test

```go
func TestUserAPI(t *testing.T) {
	if !testutil.IsIntegrationTest() {
		t.Skip("Skipping integration test")
	}

	app := testutil.NewTestApp(t)
	defer app.Cleanup()
	defer app.CleanupTables()

	// Setup HTTP context
	httpCtx := testutil.NewHTTPTestContext(t)
	
	// Setup your routes
	// router := setupRoutes(app)
	// httpCtx.Router = router

	// Create authenticated user
	authHelper := testutil.NewAuthHelper(app.JWT, t)
	user, accessToken, _ := authHelper.CreateAuthenticatedUser(app.Container.UserRepository)

	// Make authenticated request
	headers := testutil.AuthHeader(accessToken)
	resp := httpCtx.GET("/api/v1/auth/me", headers)

	// Assert response
	httpCtx.AssertStatusCode(200)

	var response map[string]interface{}
	httpCtx.GetResponseJSON(&response)
	testutil.AssertJSONFieldValue(t, response, "id", user.ID.String())
}
```

### Database Test with Cleanup

```go
func TestDatabaseOperation(t *testing.T) {
	// Setup test database
	testDB := testutil.NewTestDatabase(t)
	defer testDB.Close()
	defer testDB.Cleanup() // Truncate all tables

	ctx := context.Background()

	// Run migrations if needed
	testDB.RunMigrations()

	// Your database tests here
	testDB.ExecSQL("INSERT INTO users ...")
	
	count := testDB.CountRows("users")
	assert.Equal(t, 1, count)
}
```

### Using Fixtures

```go
func TestWithFixtures(t *testing.T) {
	app := testutil.NewTestApp(t)
	defer app.Cleanup()

	ctx := context.Background()

	// Create user with fixture
	user, err := testutil.CreateTestUserWithDefaults()
	assert.NoError(t, err)

	// Create admin user
	admin, err := testutil.CreateTestAdminUser()
	assert.NoError(t, err)

	// Create xray instance
	instance, err := testutil.CreateTestXrayInstanceWithDefaults()
	assert.NoError(t, err)

	// Create inbound for instance
	inbound, err := testutil.CreateTestInboundWithDefaults(instance.ID)
	assert.NoError(t, err)

	// Use them in your tests
	app.Container.UserRepository.Create(ctx, user)
	app.Container.UserRepository.Create(ctx, admin)
}
```

### Authentication Testing

```go
func TestAuthentication(t *testing.T) {
	app := testutil.NewTestApp(t)
	defer app.Cleanup()

	authHelper := testutil.NewAuthHelper(app.JWT, t)

	// Generate tokens
	userID := uuid.New()
	accessToken := authHelper.GenerateUserToken(userID)
	adminToken := authHelper.GenerateAdminToken(userID)

	// Validate token
	claims, err := authHelper.ValidateToken(accessToken)
	assert.NoError(t, err)
	assert.Equal(t, userID, claims.UserID)

	// Test invalid tokens
	invalidToken := authHelper.InvalidToken()
	_, err = authHelper.ValidateToken(invalidToken)
	assert.Error(t, err)
}
```

### Custom Assertions

```go
func TestAssertions(t *testing.T) {
	// Time assertions
	now := time.Now()
	testutil.AssertTimeNow(t, now, 1*time.Second)

	// UUID assertions
	id := uuid.New()
	testutil.AssertUUIDValid(t, id)

	// Error assertions
	err := errors.New("test error")
	testutil.AssertErrorContains(t, err, "test")

	// JSON assertions
	jsonMap := map[string]interface{}{"key": "value"}
	testutil.AssertJSONFieldValue(t, jsonMap, "key", "value")
}
```

### Mock Helpers

```go
func TestWithMocks(t *testing.T) {
	mockRepo := new(MockUserRepository)

	// Use helper matchers
	mockRepo.On("Create", testutil.AnyContext(), testutil.AnyUUID()).
		Return(nil)

	// Use mock call builder
	builder := testutil.NewMockCallBuilder(
		&mockRepo.Mock,
		"FindByID",
		testutil.AnyContext(),
		testutil.AnyUUID(),
	)
	builder.Return(nil, nil).Once()

	// Test your code
	// ...

	mockRepo.AssertExpectations(t)
}
```

## Environment Variables

### Test Database
- `TEST_DB_HOST` - Database host (default: localhost)
- `TEST_DB_PORT` - Database port (default: 5432)
- `TEST_DB_USER` - Database user (default: suproxy_test)
- `TEST_DB_PASSWORD` - Database password (default: suproxy_test)
- `TEST_DB_NAME` - Database name (default: suproxy_test)

### Test Control
- `INTEGRATION_TEST` - Enable integration tests (default: false)
- `CI` - Running in CI environment (default: false)

## Running Tests

### Unit Tests Only
```bash
go test ./... -short
```

### Integration Tests
```bash
INTEGRATION_TEST=true go test ./...
```

### Specific Package
```bash
INTEGRATION_TEST=true go test ./internal/application/usecase/auth
```

### With Coverage
```bash
INTEGRATION_TEST=true go test ./... -cover -coverprofile=coverage.out
```

### Verbose Output
```bash
INTEGRATION_TEST=true go test ./... -v
```

## Test Database Setup

For local development:

```bash
# Create test database
docker-compose exec postgres psql -U suproxy -c "CREATE DATABASE suproxy_test;"
docker-compose exec postgres psql -U suproxy -c "CREATE USER suproxy_test WITH PASSWORD 'suproxy_test';"
docker-compose exec postgres psql -U suproxy -c "GRANT ALL PRIVILEGES ON DATABASE suproxy_test TO suproxy_test;"
```

For CI/CD, use docker-compose with test configuration.

## Best Practices

1. **Test Isolation**: Always cleanup between tests using `CleanupTables()`
2. **Use Fixtures**: Prefer fixtures over manual entity creation
3. **Integration Test Flag**: Always check `IsIntegrationTest()` for integration tests
4. **Parallel Tests**: Ensure tests can run in parallel when possible
5. **Meaningful Assertions**: Use custom assertions for better error messages
6. **Mock Cleanup**: Always call `AssertExpectations(t)` on mocks
7. **Context Usage**: Always use context for cancellation support
8. **Resource Cleanup**: Use `defer` for cleanup functions

## Future Enhancements

- [ ] Full testcontainers-go integration
- [ ] Redis test container support
- [ ] API test client generator
- [ ] Database seeding utilities
- [ ] Performance test utilities
- [ ] Load test helpers
- [ ] Snapshot testing utilities

## Contributing

When adding new test utilities:
1. Keep functions focused and reusable
2. Add comprehensive documentation
3. Include usage examples
4. Follow existing patterns
5. Add to this README

