package integration_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/suproxy/backend/internal/domain/audit"
	"github.com/suproxy/backend/internal/infrastructure/testutil"
)

func createTestAuditLog(userID uuid.UUID) *audit.Log {
	return audit.NewLog(
		userID,
		audit.ActionCreate,
		"user",
		uuid.New(),
		"127.0.0.1",
		"test-agent",
	)
}

func TestAuditRepository_Create(t *testing.T) {
	if !testutil.IsIntegrationTest() {
		t.Skip("Skipping integration test")
	}

	app := testutil.NewTestApp(t)
	defer app.Cleanup()
	defer app.CleanupTables()

	ctx := context.Background()
	repo := app.Container.AuditLogRepository

	t.Run("Create_Success", func(t *testing.T) {
		userID := uuid.New()
		log := createTestAuditLog(userID)

		err := repo.Create(ctx, log)
		require.NoError(t, err)

		// Verify creation
		found, err := repo.FindByID(ctx, log.ID)
		require.NoError(t, err)
		assert.Equal(t, log.ID, found.ID)
		assert.Equal(t, log.UserID, found.UserID)
		assert.Equal(t, log.Action, found.Action)
	})
}

func TestAuditRepository_FindByID(t *testing.T) {
	if !testutil.IsIntegrationTest() {
		t.Skip("Skipping integration test")
	}

	app := testutil.NewTestApp(t)
	defer app.Cleanup()
	defer app.CleanupTables()

	ctx := context.Background()
	repo := app.Container.AuditLogRepository

	t.Run("FindByID_Success", func(t *testing.T) {
		userID := uuid.New()
		log := createTestAuditLog(userID)

		err := repo.Create(ctx, log)
		require.NoError(t, err)

		found, err := repo.FindByID(ctx, log.ID)
		require.NoError(t, err)
		assert.Equal(t, log.ID, found.ID)
		assert.Equal(t, log.UserID, found.UserID)
		assert.Equal(t, log.Action, found.Action)
		assert.Equal(t, log.EntityType, found.EntityType)
		assert.Equal(t, log.EntityID, found.EntityID)
		assert.Equal(t, log.IPAddress, found.IPAddress)
	})

	t.Run("FindByID_NotFound", func(t *testing.T) {
		nonExistentID := uuid.New()

		_, err := repo.FindByID(ctx, nonExistentID)
		assert.Error(t, err)
	})
}

func TestAuditRepository_FindByUserID(t *testing.T) {
	if !testutil.IsIntegrationTest() {
		t.Skip("Skipping integration test")
	}

	app := testutil.NewTestApp(t)
	defer app.Cleanup()
	defer app.CleanupTables()

	ctx := context.Background()
	repo := app.Container.AuditLogRepository

	t.Run("FindByUserID_Success", func(t *testing.T) {
		userID := uuid.New()

		// Create multiple logs for same user
		for i := 0; i < 3; i++ {
			log := createTestAuditLog(userID)
			err := repo.Create(ctx, log)
			require.NoError(t, err)
		}

		from := time.Now().Add(-1 * time.Hour)
		to := time.Now().Add(1 * time.Hour)

		logs, err := repo.FindByUserID(ctx, userID, from, to)
		require.NoError(t, err)
		assert.GreaterOrEqual(t, len(logs), 3)

		for _, log := range logs {
			assert.Equal(t, userID, log.UserID)
		}
	})

	t.Run("FindByUserID_EmptyResult", func(t *testing.T) {
		nonExistentUserID := uuid.New()
		from := time.Now().Add(-1 * time.Hour)
		to := time.Now().Add(1 * time.Hour)

		logs, err := repo.FindByUserID(ctx, nonExistentUserID, from, to)
		require.NoError(t, err)
		assert.Empty(t, logs)
	})
}

func TestAuditRepository_FindByEntityID(t *testing.T) {
	if !testutil.IsIntegrationTest() {
		t.Skip("Skipping integration test")
	}

	app := testutil.NewTestApp(t)
	defer app.Cleanup()
	defer app.CleanupTables()

	ctx := context.Background()
	repo := app.Container.AuditLogRepository

	t.Run("FindByEntityID_Success", func(t *testing.T) {
		entityID := uuid.New()
		entityType := "user"

		// Create logs for same entity
		for i := 0; i < 2; i++ {
			log := audit.NewLog(
				uuid.New(),
				audit.ActionUpdate,
				entityType,
				entityID,
				"127.0.0.1",
				"test-agent",
			)
			err := repo.Create(ctx, log)
			require.NoError(t, err)
		}

		logs, err := repo.FindByEntityID(ctx, entityType, entityID)
		require.NoError(t, err)
		assert.GreaterOrEqual(t, len(logs), 2)

		for _, log := range logs {
			assert.Equal(t, entityID, log.EntityID)
			assert.Equal(t, entityType, log.EntityType)
		}
	})
}

func TestAuditRepository_List(t *testing.T) {
	if !testutil.IsIntegrationTest() {
		t.Skip("Skipping integration test")
	}

	app := testutil.NewTestApp(t)
	defer app.Cleanup()
	defer app.CleanupTables()

	ctx := context.Background()
	repo := app.Container.AuditLogRepository

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

			userID := uuid.New()

			// Create audit logs
			for i := 0; i < tt.createCount; i++ {
				log := createTestAuditLog(userID)
				err := repo.Create(ctx, log)
				require.NoError(t, err)
			}

			// List logs
			logs, err := repo.List(ctx, tt.offset, tt.limit)
			require.NoError(t, err)
			assert.Equal(t, tt.expectedCount, len(logs))
		})
	}
}

func TestAuditRepository_ListWithFilters(t *testing.T) {
	if !testutil.IsIntegrationTest() {
		t.Skip("Skipping integration test")
	}

	app := testutil.NewTestApp(t)
	defer app.Cleanup()
	defer app.CleanupTables()

	ctx := context.Background()
	repo := app.Container.AuditLogRepository

	// Create test data
	user1ID := uuid.New()
	user2ID := uuid.New()

	// User1 logs
	for i := 0; i < 3; i++ {
		log := createTestAuditLog(user1ID)
		err := repo.Create(ctx, log)
		require.NoError(t, err)
	}

	// User2 logs
	for i := 0; i < 2; i++ {
		log := createTestAuditLog(user2ID)
		err := repo.Create(ctx, log)
		require.NoError(t, err)
	}

	t.Run("Filter_ByUserID", func(t *testing.T) {
		filters := audit.AuditFilters{
			UserID: &user1ID,
			Offset: 0,
			Limit:  10,
		}

		logs, total, err := repo.ListWithFilters(ctx, filters)
		require.NoError(t, err)
		assert.Equal(t, int64(3), total)
		assert.Len(t, logs, 3)

		for _, log := range logs {
			assert.Equal(t, user1ID, log.UserID)
		}
	})

	t.Run("Filter_ByAction", func(t *testing.T) {
		action := string(audit.ActionCreate)
		filters := audit.AuditFilters{
			Action: &action,
			Offset: 0,
			Limit:  10,
		}

		logs, total, err := repo.ListWithFilters(ctx, filters)
		require.NoError(t, err)
		assert.GreaterOrEqual(t, total, int64(5))

		for _, log := range logs {
			assert.Equal(t, audit.Action(action), log.Action)
		}
	})

	t.Run("Filter_ByEntityType", func(t *testing.T) {
		entityType := "user"
		filters := audit.AuditFilters{
			EntityType: &entityType,
			Offset:     0,
			Limit:      10,
		}

		logs, total, err := repo.ListWithFilters(ctx, filters)
		require.NoError(t, err)
		assert.GreaterOrEqual(t, total, int64(5))

		for _, log := range logs {
			assert.Equal(t, entityType, log.EntityType)
		}
	})

	t.Run("Filter_ByDateRange", func(t *testing.T) {
		from := time.Now().Add(-1 * time.Hour)
		to := time.Now().Add(1 * time.Hour)

		filters := audit.AuditFilters{
			DateFrom: &from,
			DateTo:   &to,
			Offset:   0,
			Limit:    10,
		}

		logs, total, err := repo.ListWithFilters(ctx, filters)
		require.NoError(t, err)
		assert.GreaterOrEqual(t, total, int64(5))

		for _, log := range logs {
			assert.True(t, log.CreatedAt.After(from))
			assert.True(t, log.CreatedAt.Before(to))
		}
	})

	t.Run("Filter_WithPagination", func(t *testing.T) {
		filters := audit.AuditFilters{
			Offset: 0,
			Limit:  2,
		}

		logs, total, err := repo.ListWithFilters(ctx, filters)
		require.NoError(t, err)
		assert.GreaterOrEqual(t, total, int64(5))
		assert.Len(t, logs, 2)
	})
}

func TestAuditRepository_Count(t *testing.T) {
	if !testutil.IsIntegrationTest() {
		t.Skip("Skipping integration test")
	}

	app := testutil.NewTestApp(t)
	defer app.Cleanup()
	defer app.CleanupTables()

	ctx := context.Background()
	repo := app.Container.AuditLogRepository

	t.Run("Count_EmptyDatabase", func(t *testing.T) {
		count, err := repo.Count(ctx)
		require.NoError(t, err)
		assert.Equal(t, int64(0), count)
	})

	t.Run("Count_WithLogs", func(t *testing.T) {
		defer app.CleanupTables()

		userID := uuid.New()
		for i := 0; i < 5; i++ {
			log := createTestAuditLog(userID)
			err := repo.Create(ctx, log)
			require.NoError(t, err)
		}

		count, err := repo.Count(ctx)
		require.NoError(t, err)
		assert.Equal(t, int64(5), count)
	})
}

func TestAuditRepository_CountByAction(t *testing.T) {
	if !testutil.IsIntegrationTest() {
		t.Skip("Skipping integration test")
	}

	app := testutil.NewTestApp(t)
	defer app.Cleanup()
	defer app.CleanupTables()

	ctx := context.Background()
	repo := app.Container.AuditLogRepository

	t.Run("CountByAction_Success", func(t *testing.T) {
		userID := uuid.New()

		// Create different action types
		actions := []audit.Action{audit.ActionCreate, audit.ActionUpdate, audit.ActionDelete}
		for _, action := range actions {
			for i := 0; i < 2; i++ {
				log := audit.NewLog(
					userID,
					action,
					"user",
					uuid.New(),
					"127.0.0.1",
					"test-agent",
				)
				err := repo.Create(ctx, log)
				require.NoError(t, err)
			}
		}

		counts, err := repo.CountByAction(ctx)
		require.NoError(t, err)
		assert.GreaterOrEqual(t, len(counts), 3)

		for action, count := range counts {
			if action == string(audit.ActionCreate) || action == string(audit.ActionUpdate) || action == string(audit.ActionDelete) {
				assert.GreaterOrEqual(t, count, int64(2))
			}
		}
	})
}

func TestAuditRepository_CountByEntityType(t *testing.T) {
	if !testutil.IsIntegrationTest() {
		t.Skip("Skipping integration test")
	}

	app := testutil.NewTestApp(t)
	defer app.Cleanup()
	defer app.CleanupTables()

	ctx := context.Background()
	repo := app.Container.AuditLogRepository

	t.Run("CountByEntityType_Success", func(t *testing.T) {
		userID := uuid.New()

		// Create different entity types
		entityTypes := []string{"user", "xray_instance", "inbound"}
		for _, entityType := range entityTypes {
			for i := 0; i < 2; i++ {
				log := audit.NewLog(
					userID,
					audit.ActionCreate,
					entityType,
					uuid.New(),
					"127.0.0.1",
					"test-agent",
				)
				err := repo.Create(ctx, log)
				require.NoError(t, err)
			}
		}

		counts, err := repo.CountByEntityType(ctx)
		require.NoError(t, err)
		assert.GreaterOrEqual(t, len(counts), 3)

		for entityType, count := range counts {
			if entityType == "user" || entityType == "xray_instance" || entityType == "inbound" {
				assert.GreaterOrEqual(t, count, int64(2))
			}
		}
	})
}

func TestAuditRepository_CountUniqueUsers(t *testing.T) {
	if !testutil.IsIntegrationTest() {
		t.Skip("Skipping integration test")
	}

	app := testutil.NewTestApp(t)
	defer app.Cleanup()
	defer app.CleanupTables()

	ctx := context.Background()
	repo := app.Container.AuditLogRepository

	t.Run("CountUniqueUsers_Success", func(t *testing.T) {
		// Create logs for 3 different users
		users := []uuid.UUID{uuid.New(), uuid.New(), uuid.New()}
		for _, userID := range users {
			for i := 0; i < 2; i++ {
				log := createTestAuditLog(userID)
				err := repo.Create(ctx, log)
				require.NoError(t, err)
			}
		}

		count, err := repo.CountUniqueUsers(ctx)
		require.NoError(t, err)
		assert.Equal(t, int64(3), count)
	})
}

func TestAuditRepository_CountUniqueIPs(t *testing.T) {
	if !testutil.IsIntegrationTest() {
		t.Skip("Skipping integration test")
	}

	app := testutil.NewTestApp(t)
	defer app.Cleanup()
	defer app.CleanupTables()

	ctx := context.Background()
	repo := app.Container.AuditLogRepository

	t.Run("CountUniqueIPs_Success", func(t *testing.T) {
		userID := uuid.New()
		ips := []string{"192.168.1.1", "192.168.1.2", "192.168.1.3"}

		for _, ip := range ips {
			for i := 0; i < 2; i++ {
				log := audit.NewLog(
					userID,
					audit.ActionLogin,
					"user",
					uuid.New(),
					ip,
					"test-agent",
				)
				err := repo.Create(ctx, log)
				require.NoError(t, err)
			}
		}

		count, err := repo.CountUniqueIPs(ctx)
		require.NoError(t, err)
		assert.Equal(t, int64(3), count)
	})
}

func TestAuditRepository_GetOldestAndNewestLogDate(t *testing.T) {
	if !testutil.IsIntegrationTest() {
		t.Skip("Skipping integration test")
	}

	app := testutil.NewTestApp(t)
	defer app.Cleanup()
	defer app.CleanupTables()

	ctx := context.Background()
	repo := app.Container.AuditLogRepository

	t.Run("GetLogDates_Success", func(t *testing.T) {
		// Clean up first to ensure empty state
		app.CleanupTables()
		
		userID := uuid.New()

		// Create logs
		for i := 0; i < 3; i++ {
			log := createTestAuditLog(userID)
			err := repo.Create(ctx, log)
			require.NoError(t, err)
			time.Sleep(10 * time.Millisecond) // Small delay to ensure different timestamps
		}

		oldest, err := repo.GetOldestLogDate(ctx)
		require.NoError(t, err)
		require.NotNil(t, oldest)

		newest, err := repo.GetNewestLogDate(ctx)
		require.NoError(t, err)
		require.NotNil(t, newest)

		assert.True(t, oldest.Before(*newest) || oldest.Equal(*newest))
	})

	t.Run("GetLogDates_EmptyDatabase", func(t *testing.T) {
		// Clean up to ensure empty state
		app.CleanupTables()

		oldest, err := repo.GetOldestLogDate(ctx)
		require.NoError(t, err)
		assert.Nil(t, oldest)

		newest, err := repo.GetNewestLogDate(ctx)
		require.NoError(t, err)
		assert.Nil(t, newest)
	})
}

func TestAuditRepository_DeleteOlderThan(t *testing.T) {
	if !testutil.IsIntegrationTest() {
		t.Skip("Skipping integration test")
	}

	app := testutil.NewTestApp(t)
	defer app.Cleanup()
	defer app.CleanupTables()

	ctx := context.Background()
	repo := app.Container.AuditLogRepository

	t.Run("DeleteOlderThan_Success", func(t *testing.T) {
		userID := uuid.New()

		// Create logs
		for i := 0; i < 5; i++ {
			log := createTestAuditLog(userID)
			err := repo.Create(ctx, log)
			require.NoError(t, err)
		}

		initialCount, err := repo.Count(ctx)
		require.NoError(t, err)
		assert.Equal(t, int64(5), initialCount)

		// Delete logs older than future date (should delete all)
		futureDate := time.Now().Add(1 * time.Hour)
		err = repo.DeleteOlderThan(ctx, futureDate)
		require.NoError(t, err)

		// Verify deletion
		finalCount, err := repo.Count(ctx)
		require.NoError(t, err)
		assert.Equal(t, int64(0), finalCount)
	})
}
