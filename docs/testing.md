# Testing Guide

## Overview

SuProxy Backend uses a comprehensive testing strategy with both unit tests and integration tests.

## Test Structure

```
.
├── internal/
│   ├── application/
│   │   └── service/
│   │       └── *_test.go          # Unit tests
│   ├── infrastructure/
│   │   └── testutil/               # Test utilities
│   │       ├── config.go
│   │       ├── database.go
│   │       ├── bootstrap.go
│   │       ├── fixtures.go
│   │       ├── http.go
│   │       ├── auth.go
│   │       ├── assert.go
│   │       ├── mock_helpers.go
│   │       ├── container.go
│   │       └── README.md
│   └── domain/
│       └── user/
│           └── *_test.go           # Domain unit tests
└── test/
    └── integration/
        └── *_test.go                # Integration tests
```

## Test Types

### Unit Tests

Unit tests test individual components in isolation using mocks.

**Location**: Next to the code being tested (`*_test.go`)

**Example**:
```go
package service

import (
	"testing"
	"github.com/stretchr/testify/mock"
)

func TestMyService(t *testing.T) {
	// Setup mocks
	mockRepo := new(MockRepository)
	mockRepo.On("FindByID", mock.Anything, mock.Anything).Return(nil, nil)

	// Test
	service := NewMyService(mockRepo)
	// ... assertions
}
```

**Run**:
```bash
make test
# or
go test -short ./...
```

### Integration Tests

Integration tests test the entire system with real database and dependencies.

**Location**: `test/integration/`

**Example**:
```go
func TestUserRepository_Integration(t *testing.T) {
	if !testutil.IsIntegrationTest() {
		t.Skip("Skipping integration test")
	}

	app := testutil.NewTestApp(t)
	defer app.Cleanup()
	defer app.CleanupTables()

	// Test with real database
	// ...
}
```

**Run**:
```bash
make test-integration
# or
INTEGRATION_TEST=true go test ./...
```

## Test Utilities

The `testutil` package provides comprehensive testing utilities:

### Configuration
```go
cfg := testutil.TestConfig()        // Test configuration
minCfg := testutil.MinimalConfig()  // Minimal config for unit tests
```

### Database
```go
testDB := testutil.NewTestDatabase(t)
defer testDB.Close()
defer testDB.Cleanup()  // Truncate all tables

testDB.ExecSQL("INSERT INTO ...")
count := testDB.CountRows("users")
```

### Application Bootstrap
```go
app := testutil.NewTestApp(t)
defer app.Cleanup()
defer app.CleanupTables()

// Access all repositories
userRepo := app.Container.UserRepository
```

### Fixtures
```go
// Create test entities
user, err := testutil.CreateTestUserWithDefaults()
admin, err := testutil.CreateTestAdminUser()
instance, err := testutil.CreateTestXrayInstanceWithDefaults()
inbound, err := testutil.CreateTestInboundWithDefaults(instanceID)
client, err := testutil.CreateTestClientWithDefaults(inboundID, userID)
```

### HTTP Testing
```go
httpCtx := testutil.NewHTTPTestContext(t)

// Make requests
resp := httpCtx.GET("/api/v1/users", nil)
resp := httpCtx.POST("/api/v1/users", body, headers)

// Assert response
httpCtx.AssertStatusCode(200)
httpCtx.AssertJSONResponse(200, &result)
```

### Authentication
```go
authHelper := testutil.NewAuthHelper(app.JWT, t)

// Generate tokens
accessToken := authHelper.GenerateUserToken(userID)
adminToken := authHelper.GenerateAdminToken(userID)

// Create authenticated users
user, accessToken, refreshToken := authHelper.CreateAuthenticatedUser(userRepo)
admin, accessToken, refreshToken := authHelper.CreateAuthenticatedAdmin(userRepo)
```

### Custom Assertions
```go
testutil.AssertTimeNow(t, timestamp, 1*time.Second)
testutil.AssertUUIDValid(t, id)
testutil.AssertErrorContains(t, err, "expected message")
testutil.AssertJSONFieldValue(t, jsonMap, "field", "value")
```

## Running Tests

### All Commands

```bash
# Unit tests only
make test
make test-unit

# Integration tests only
make test-integration

# All tests (unit + integration)
make test-all

# With coverage
make test-coverage
make test-coverage-integration
make test-coverage-all

# Verbose output
make test-verbose

# Specific package
make test-package PKG=./internal/domain/user

# Clean test cache
make test-clean
```

### Environment Variables

```bash
# Enable integration tests
INTEGRATION_TEST=true go test ./...

# Set test database
TEST_DB_HOST=localhost
TEST_DB_PORT=5432
TEST_DB_USER=suproxy_test
TEST_DB_PASSWORD=suproxy_test
TEST_DB_NAME=suproxy_test

# CI environment
CI=true INTEGRATION_TEST=true go test ./...
```

### Docker-based Testing

```bash
# Start test database
docker-compose up -d postgres

# Setup test database
make test-db-setup

# Run integration tests
make test-integration

# Teardown test database
make test-db-teardown
```

## Test Database Setup

### Local Setup

```bash
# Create test database and user
docker-compose exec postgres psql -U suproxy << EOF
CREATE DATABASE suproxy_test;
CREATE USER suproxy_test WITH PASSWORD 'suproxy_test';
GRANT ALL PRIVILEGES ON DATABASE suproxy_test TO suproxy_test;
EOF
```

### Automated Setup

```bash
make test-db-setup
```

This creates:
- Database: `suproxy_test`
- User: `suproxy_test`
- Password: `suproxy_test`

## Writing Tests

### Unit Test Best Practices

1. **Test in isolation**: Use mocks for dependencies
2. **One assertion per test**: Keep tests focused
3. **Descriptive names**: `TestServiceName_Scenario_ExpectedResult`
4. **Table-driven tests**: For testing multiple scenarios
5. **Setup and teardown**: Use `t.Cleanup()` for cleanup

**Example**:
```go
func TestUserService_CreateUser_Success(t *testing.T) {
	// Arrange
	mockRepo := new(MockUserRepository)
	mockRepo.On("Create", mock.Anything, mock.Anything).Return(nil)
	service := NewUserService(mockRepo)

	// Act
	err := service.CreateUser(ctx, userData)

	// Assert
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}
```

### Integration Test Best Practices

1. **Check integration flag**: Always check `IsIntegrationTest()`
2. **Cleanup after tests**: Use `defer app.CleanupTables()`
3. **Test isolation**: Each test should be independent
4. **Real dependencies**: Use real database, not mocks
5. **Test transactions**: Ensure ACID properties

**Example**:
```go
func TestUserCreation_Integration(t *testing.T) {
	if !testutil.IsIntegrationTest() {
		t.Skip("Skipping integration test")
	}

	app := testutil.NewTestApp(t)
	defer app.Cleanup()
	defer app.CleanupTables()

	ctx := context.Background()
	
	// Test with real database
	user, err := testutil.CreateTestUserWithDefaults()
	require.NoError(t, err)

	err = app.Container.UserRepository.Create(ctx, user)
	require.NoError(t, err)

	// Verify
	found, err := app.Container.UserRepository.FindByID(ctx, user.ID)
	require.NoError(t, err)
	assert.Equal(t, user.ID, found.ID)
}
```

### Table-Driven Tests

```go
func TestValidation(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{"valid email", "test@example.com", false},
		{"invalid email", "invalid", true},
		{"empty email", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateEmail(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
```

## CI/CD Integration

### GitHub Actions Example

```yaml
name: Tests

on: [push, pull_request]

jobs:
  test:
    runs-on: ubuntu-latest

    services:
      postgres:
        image: postgres:15-alpine
        env:
          POSTGRES_USER: suproxy_test
          POSTGRES_PASSWORD: suproxy_test
          POSTGRES_DB: suproxy_test
        ports:
          - 5432:5432
        options: >-
          --health-cmd pg_isready
          --health-interval 10s
          --health-timeout 5s
          --health-retries 5

    steps:
      - uses: actions/checkout@v3
      
      - name: Set up Go
        uses: actions/setup-go@v4
        with:
          go-version: '1.21'
      
      - name: Run unit tests
        run: make test
      
      - name: Run integration tests
        env:
          INTEGRATION_TEST: true
          TEST_DB_HOST: localhost
          TEST_DB_PORT: 5432
        run: make test-integration
      
      - name: Upload coverage
        uses: codecov/codecov-action@v3
        with:
          files: ./coverage-all.out
```

## Coverage

### Generate Coverage Report

```bash
# Unit tests coverage
make test-coverage

# Integration tests coverage
make test-coverage-integration

# All tests coverage
make test-coverage-all
```

### Coverage Goals

- **Overall**: 80%+
- **Business Logic**: 90%+
- **Handlers**: 70%+
- **Infrastructure**: 60%+

### View Coverage in Browser

```bash
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

## Debugging Tests

### Run Single Test

```bash
go test -v -run TestName ./path/to/package
```

### Debug with Delve

```bash
dlv test ./path/to/package -- -test.run TestName
```

### Print Debug Info

```go
t.Logf("Debug info: %+v", value)
```

### Skip Cleanup for Debugging

```go
// Comment out cleanup to inspect database state
// defer app.CleanupTables()
```

## Performance Testing

### Benchmark Tests

```go
func BenchmarkUserCreate(b *testing.B) {
	app := testutil.NewTestApp(b)
	defer app.Cleanup()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		user, _ := testutil.CreateTestUserWithDefaults()
		_ = app.Container.UserRepository.Create(context.Background(), user)
	}
}
```

**Run**:
```bash
go test -bench=. -benchmem ./...
```

## Troubleshooting

### Tests Fail with Database Connection Error

- Ensure PostgreSQL is running: `docker-compose ps`
- Check test database exists: `make test-db-setup`
- Verify environment variables in `.env.test`

### Integration Tests Skipped

- Set `INTEGRATION_TEST=true`
- Or run with: `make test-integration`

### Test Cache Issues

- Clean test cache: `make test-clean`
- Force rerun: `go test -count=1 ./...`

### Database State Issues

- Ensure `CleanupTables()` is called
- Check for transaction issues
- Manually truncate: `make test-db-teardown && make test-db-setup`

## Resources

- [Go Testing Package](https://pkg.go.dev/testing)
- [Testify Documentation](https://github.com/stretchr/testify)
- [Test Utilities README](../internal/infrastructure/testutil/README.md)

