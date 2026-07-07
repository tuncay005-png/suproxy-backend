package integration_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/suproxy/backend/internal/domain/user"
	"github.com/suproxy/backend/internal/infrastructure/testutil"
)

// TestUserRepository_Integration tests user repository operations
func TestUserRepository_Integration(t *testing.T) {
	if !testutil.IsIntegrationTest() {
		t.Skip("Skipping integration test (set INTEGRATION_TEST=true to run)")
	}

	app := testutil.NewTestApp(t)
	defer app.Cleanup()
	defer app.CleanupTables()

	ctx := context.Background()
	userRepo := app.Container.UserRepository

	t.Run("Create and Find User", func(t *testing.T) {
		// Create test user
		testUser, err := testutil.CreateTestUserWithDefaults()
		require.NoError(t, err)

		// Save to database
		err = userRepo.Create(ctx, testUser)
		require.NoError(t, err)

		// Find by ID
		found, err := userRepo.FindByID(ctx, testUser.ID)
		require.NoError(t, err)
		assert.Equal(t, testUser.ID, found.ID)
		assert.Equal(t, testUser.Email.String(), found.Email.String())
	})

	t.Run("Find User by Email", func(t *testing.T) {
		// Create test user
		testUser, err := testutil.CreateTestUserWithDefaults()
		require.NoError(t, err)

		err = userRepo.Create(ctx, testUser)
		require.NoError(t, err)

		// Find by email
		found, err := userRepo.FindByEmail(ctx, testUser.Email.String())
		require.NoError(t, err)
		assert.Equal(t, testUser.ID, found.ID)
	})

	t.Run("Update User", func(t *testing.T) {
		// Create test user
		testUser, err := testutil.CreateTestUserWithDefaults()
		require.NoError(t, err)

		err = userRepo.Create(ctx, testUser)
		require.NoError(t, err)

		// Update user
		testUser.PromoteToAdmin()
		err = userRepo.Update(ctx, testUser)
		require.NoError(t, err)

		// Verify update
		found, err := userRepo.FindByID(ctx, testUser.ID)
		require.NoError(t, err)
		assert.Equal(t, user.RoleAdmin, found.Role)
	})

	t.Run("Delete User", func(t *testing.T) {
		// Create test user
		testUser, err := testutil.CreateTestUserWithDefaults()
		require.NoError(t, err)

		err = userRepo.Create(ctx, testUser)
		require.NoError(t, err)

		// Delete user
		err = userRepo.Delete(ctx, testUser.ID)
		require.NoError(t, err)

		// Verify deletion
		_, err = userRepo.FindByID(ctx, testUser.ID)
		assert.Error(t, err)
	})

	t.Run("List Users with Pagination", func(t *testing.T) {
		// Create multiple users
		authHelper := testutil.NewAuthHelper(app.JWT, t)
		users := authHelper.CreateMultipleUsers(userRepo, 5)

		// List users
		found, err := userRepo.List(ctx, 0, 10)
		require.NoError(t, err)
		assert.GreaterOrEqual(t, len(found), len(users))
	})

	t.Run("Count Users", func(t *testing.T) {
		// Get initial count
		initialCount, err := userRepo.Count(ctx)
		require.NoError(t, err)

		// Create test user
		testUser, err := testutil.CreateTestUserWithDefaults()
		require.NoError(t, err)

		err = userRepo.Create(ctx, testUser)
		require.NoError(t, err)

		// Verify count increased
		newCount, err := userRepo.Count(ctx)
		require.NoError(t, err)
		assert.Equal(t, initialCount+1, newCount)
	})
}

// TestUserAuthentication_Integration tests user authentication flow
func TestUserAuthentication_Integration(t *testing.T) {
	if !testutil.IsIntegrationTest() {
		t.Skip("Skipping integration test (set INTEGRATION_TEST=true to run)")
	}

	app := testutil.NewTestApp(t)
	defer app.Cleanup()
	defer app.CleanupTables()

	ctx := context.Background()
	userRepo := app.Container.UserRepository
	authHelper := testutil.NewAuthHelper(app.JWT, t)

	t.Run("Generate and Validate Access Token", func(t *testing.T) {
		// Create user
		testUser, err := testutil.CreateTestUserWithDefaults()
		require.NoError(t, err)

		err = userRepo.Create(ctx, testUser)
		require.NoError(t, err)

		// Generate token
		accessToken := authHelper.GenerateUserToken(testUser.ID)
		require.NotEmpty(t, accessToken)

		// Validate token
		claims, err := authHelper.ValidateToken(accessToken)
		require.NoError(t, err)
		assert.Equal(t, testUser.ID, claims.UserID)
		assert.Equal(t, user.RoleUser, claims.Role)
	})

	t.Run("Generate Admin Token", func(t *testing.T) {
		// Create admin user
		adminUser, err := testutil.CreateTestAdminUser()
		require.NoError(t, err)

		err = userRepo.Create(ctx, adminUser)
		require.NoError(t, err)

		// Generate admin token
		adminToken := authHelper.GenerateAdminToken(adminUser.ID)
		require.NotEmpty(t, adminToken)

		// Validate token
		claims, err := authHelper.ValidateToken(adminToken)
		require.NoError(t, err)
		assert.Equal(t, adminUser.ID, claims.UserID)
		assert.Equal(t, user.RoleAdmin, claims.Role)
	})

	t.Run("Generate Token Pair", func(t *testing.T) {
		// Create user
		testUser, err := testutil.CreateTestUserWithDefaults()
		require.NoError(t, err)

		err = userRepo.Create(ctx, testUser)
		require.NoError(t, err)

		// Generate token pair
		tokenPair := authHelper.GenerateTokenPair(testUser.ID, user.RoleUser)
		require.NotEmpty(t, tokenPair.AccessToken)
		require.NotEmpty(t, tokenPair.RefreshToken)

		// Validate both tokens
		accessClaims, err := authHelper.ValidateToken(tokenPair.AccessToken)
		require.NoError(t, err)
		assert.Equal(t, testUser.ID, accessClaims.UserID)

		refreshClaims, err := authHelper.ValidateToken(tokenPair.RefreshToken)
		require.NoError(t, err)
		assert.Equal(t, testUser.ID, refreshClaims.UserID)
	})

	t.Run("Reject Invalid Token", func(t *testing.T) {
		invalidToken := authHelper.InvalidToken()

		_, err := authHelper.ValidateToken(invalidToken)
		assert.Error(t, err)
	})
}

