package integration_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/suproxy/backend/internal/domain/user"
	"github.com/suproxy/backend/internal/infrastructure/testutil"
)

func TestUserRepository_Create(t *testing.T) {
	if !testutil.IsIntegrationTest() {
		t.Skip("Skipping integration test")
	}

	app := testutil.NewTestApp(t)
	defer app.Cleanup()
	defer app.CleanupTables()

	ctx := context.Background()
	repo := app.Container.UserRepository

	t.Run("Create_Success", func(t *testing.T) {
		testUser, err := testutil.CreateTestUserWithDefaults()
		require.NoError(t, err)

		err = repo.Create(ctx, testUser)
		require.NoError(t, err)

		// Verify creation
		found, err := repo.FindByID(ctx, testUser.ID)
		require.NoError(t, err)
		assert.Equal(t, testUser.ID, found.ID)
		assert.Equal(t, testUser.Email.String(), found.Email.String())
	})

	t.Run("Create_DuplicateEmail_Error", func(t *testing.T) {
		testUser, err := testutil.CreateTestUserWithDefaults()
		require.NoError(t, err)

		err = repo.Create(ctx, testUser)
		require.NoError(t, err)

		// Try to create another user with same email
		duplicateUser, err := testutil.CreateTestUser("anotheruser", testUser.Email.String(), "Pass123!@#")
		require.NoError(t, err)

		err = repo.Create(ctx, duplicateUser)
		assert.Error(t, err)
	})
}

func TestUserRepository_FindByID(t *testing.T) {
	if !testutil.IsIntegrationTest() {
		t.Skip("Skipping integration test")
	}

	app := testutil.NewTestApp(t)
	defer app.Cleanup()
	defer app.CleanupTables()

	ctx := context.Background()
	repo := app.Container.UserRepository

	t.Run("FindByID_Success", func(t *testing.T) {
		testUser, err := testutil.CreateTestUserWithDefaults()
		require.NoError(t, err)

		err = repo.Create(ctx, testUser)
		require.NoError(t, err)

		found, err := repo.FindByID(ctx, testUser.ID)
		require.NoError(t, err)
		assert.Equal(t, testUser.ID, found.ID)
		assert.Equal(t, testUser.Email.String(), found.Email.String())
		assert.Equal(t, testUser.Role, found.Role)
		assert.Equal(t, testUser.Status, found.Status)
	})

	t.Run("FindByID_NotFound", func(t *testing.T) {
		nonExistentID := uuid.New()

		_, err := repo.FindByID(ctx, nonExistentID)
		assert.Error(t, err)
	})
}

func TestUserRepository_FindByEmail(t *testing.T) {
	if !testutil.IsIntegrationTest() {
		t.Skip("Skipping integration test")
	}

	app := testutil.NewTestApp(t)
	defer app.Cleanup()
	defer app.CleanupTables()

	ctx := context.Background()
	repo := app.Container.UserRepository

	t.Run("FindByEmail_Success", func(t *testing.T) {
		testUser, err := testutil.CreateTestUserWithDefaults()
		require.NoError(t, err)

		err = repo.Create(ctx, testUser)
		require.NoError(t, err)

		found, err := repo.FindByEmail(ctx, testUser.Email)
		require.NoError(t, err)
		assert.Equal(t, testUser.ID, found.ID)
		assert.Equal(t, testUser.Email.String(), found.Email.String())
	})

	t.Run("FindByEmail_NotFound", func(t *testing.T) {
		nonExistentEmail, _ := user.NewEmail("nonexistent@example.com")

		_, err := repo.FindByEmail(ctx, nonExistentEmail)
		assert.Error(t, err)
	})
}

func TestUserRepository_Update(t *testing.T) {
	if !testutil.IsIntegrationTest() {
		t.Skip("Skipping integration test")
	}

	app := testutil.NewTestApp(t)
	defer app.Cleanup()
	defer app.CleanupTables()

	ctx := context.Background()
	repo := app.Container.UserRepository

	t.Run("Update_Success", func(t *testing.T) {
		testUser, err := testutil.CreateTestUserWithDefaults()
		require.NoError(t, err)

		err = repo.Create(ctx, testUser)
		require.NoError(t, err)

		// Update user
		testUser.Role = user.RoleAdmin
		newProfile := user.NewProfile("Updated", "Name", "1234567890", "avatar.jpg")
		testUser.UpdateProfile(newProfile)

		err = repo.Update(ctx, testUser)
		require.NoError(t, err)

		// Verify update
		found, err := repo.FindByID(ctx, testUser.ID)
		require.NoError(t, err)
		assert.Equal(t, user.RoleAdmin, found.Role)
		assert.Equal(t, "Updated", found.Profile.FirstName)
		assert.Equal(t, "Name", found.Profile.LastName)
	})

	t.Run("Update_NotFound", func(t *testing.T) {
		testUser, err := testutil.CreateTestUserWithDefaults()
		require.NoError(t, err)
		testUser.ID = uuid.New() // Non-existent ID

		err = repo.Update(ctx, testUser)
		assert.Error(t, err)
	})
}

func TestUserRepository_Delete(t *testing.T) {
	if !testutil.IsIntegrationTest() {
		t.Skip("Skipping integration test")
	}

	app := testutil.NewTestApp(t)
	defer app.Cleanup()
	defer app.CleanupTables()

	ctx := context.Background()
	repo := app.Container.UserRepository

	t.Run("Delete_Success", func(t *testing.T) {
		testUser, err := testutil.CreateTestUserWithDefaults()
		require.NoError(t, err)

		err = repo.Create(ctx, testUser)
		require.NoError(t, err)

		err = repo.Delete(ctx, testUser.ID)
		require.NoError(t, err)

		// Verify deletion
		_, err = repo.FindByID(ctx, testUser.ID)
		assert.Error(t, err)
	})

	t.Run("Delete_NotFound", func(t *testing.T) {
		nonExistentID := uuid.New()

		err := repo.Delete(ctx, nonExistentID)
		assert.Error(t, err)
	})
}

func TestUserRepository_List(t *testing.T) {
	if !testutil.IsIntegrationTest() {
		t.Skip("Skipping integration test")
	}

	app := testutil.NewTestApp(t)
	defer app.Cleanup()
	defer app.CleanupTables()

	ctx := context.Background()
	repo := app.Container.UserRepository
	authHelper := testutil.NewAuthHelper(app.JWT, t)

	tests := []struct {
		name          string
		createCount   int
		offset        int
		limit         int
		expectedCount int
	}{
		{"List_All", 5, 0, 10, 5},
		{"List_WithOffset", 5, 2, 10, 3},
		{"List_WithLimit", 5, 0, 3, 3},
		{"List_WithOffsetAndLimit", 5, 1, 2, 2},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer app.CleanupTables()

			// Create users
			authHelper.CreateMultipleUsers(repo, tt.createCount)

			// List users
			users, err := repo.List(ctx, tt.offset, tt.limit)
			require.NoError(t, err)
			assert.Equal(t, tt.expectedCount, len(users))
		})
	}
}

func TestUserRepository_Count(t *testing.T) {
	if !testutil.IsIntegrationTest() {
		t.Skip("Skipping integration test")
	}

	app := testutil.NewTestApp(t)
	defer app.Cleanup()
	defer app.CleanupTables()

	ctx := context.Background()
	repo := app.Container.UserRepository
	authHelper := testutil.NewAuthHelper(app.JWT, t)

	t.Run("Count_EmptyDatabase", func(t *testing.T) {
		count, err := repo.Count(ctx)
		require.NoError(t, err)
		assert.Equal(t, int64(0), count)
	})

	t.Run("Count_WithUsers", func(t *testing.T) {
		defer app.CleanupTables()

		authHelper.CreateMultipleUsers(repo, 5)

		count, err := repo.Count(ctx)
		require.NoError(t, err)
		assert.Equal(t, int64(5), count)
	})
}

func TestUserRepository_ExistsByEmail(t *testing.T) {
	if !testutil.IsIntegrationTest() {
		t.Skip("Skipping integration test")
	}

	app := testutil.NewTestApp(t)
	defer app.Cleanup()
	defer app.CleanupTables()

	ctx := context.Background()
	repo := app.Container.UserRepository

	t.Run("ExistsByEmail_True", func(t *testing.T) {
		testUser, err := testutil.CreateTestUserWithDefaults()
		require.NoError(t, err)

		err = repo.Create(ctx, testUser)
		require.NoError(t, err)

		exists, err := repo.ExistsByEmail(ctx, testUser.Email)
		require.NoError(t, err)
		assert.True(t, exists)
	})

	t.Run("ExistsByEmail_False", func(t *testing.T) {
		nonExistentEmail, _ := user.NewEmail("nonexistent@example.com")

		exists, err := repo.ExistsByEmail(ctx, nonExistentEmail)
		require.NoError(t, err)
		assert.False(t, exists)
	})
}

func TestUserRepository_ListWithFilters(t *testing.T) {
	if !testutil.IsIntegrationTest() {
		t.Skip("Skipping integration test")
	}

	app := testutil.NewTestApp(t)
	defer app.Cleanup()
	defer app.CleanupTables()

	ctx := context.Background()
	repo := app.Container.UserRepository

	// Create test users
	adminUser, err := testutil.CreateTestAdminUser()
	require.NoError(t, err)
	err = repo.Create(ctx, adminUser)
	require.NoError(t, err)

	regularUser, err := testutil.CreateTestUserWithDefaults()
	require.NoError(t, err)
	err = repo.Create(ctx, regularUser)
	require.NoError(t, err)

	t.Run("Filter_ByRole", func(t *testing.T) {
		adminRole := user.RoleAdmin
		filters := user.UserFilters{
			Offset: 0,
			Limit:  10,
			Role:   &adminRole,
		}

		users, total, err := repo.ListWithFilters(ctx, filters)
		require.NoError(t, err)
		assert.Equal(t, int64(1), total)
		assert.Len(t, users, 1)
		assert.Equal(t, user.RoleAdmin, users[0].Role)
	})

	t.Run("Filter_ByStatus", func(t *testing.T) {
		activeStatus := user.StatusActive
		filters := user.UserFilters{
			Offset: 0,
			Limit:  10,
			Status: &activeStatus,
		}

		users, total, err := repo.ListWithFilters(ctx, filters)
		require.NoError(t, err)
		assert.GreaterOrEqual(t, total, int64(2))
		assert.GreaterOrEqual(t, len(users), 2)
	})

	t.Run("Filter_WithPagination", func(t *testing.T) {
		filters := user.UserFilters{
			Offset: 0,
			Limit:  1,
		}

		users, total, err := repo.ListWithFilters(ctx, filters)
		require.NoError(t, err)
		assert.GreaterOrEqual(t, total, int64(2))
		assert.Len(t, users, 1)
	})
}

