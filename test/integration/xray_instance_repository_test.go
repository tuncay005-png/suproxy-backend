package integration_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/suproxy/backend/internal/domain/xray"
	"github.com/suproxy/backend/internal/infrastructure/testutil"
)

func TestXrayInstanceRepository_Create(t *testing.T) {
	if !testutil.IsIntegrationTest() {
		t.Skip("Skipping integration test")
	}

	app := testutil.NewTestApp(t)
	defer app.Cleanup()
	defer app.CleanupTables()

	ctx := context.Background()
	repo := app.Container.XrayInstanceRepository

	t.Run("Create_Success", func(t *testing.T) {
		instance, err := testutil.CreateTestXrayInstanceWithDefaults()
		require.NoError(t, err)

		err = repo.Create(ctx, instance)
		require.NoError(t, err)

		// Verify creation
		found, err := repo.FindByID(ctx, instance.ID)
		require.NoError(t, err)
		assert.Equal(t, instance.ID, found.ID)
		assert.Equal(t, instance.Name, found.Name)
		assert.Equal(t, instance.Protocol, found.Protocol)
	})
}

func TestXrayInstanceRepository_FindByID(t *testing.T) {
	if !testutil.IsIntegrationTest() {
		t.Skip("Skipping integration test")
	}

	app := testutil.NewTestApp(t)
	defer app.Cleanup()
	defer app.CleanupTables()

	ctx := context.Background()
	repo := app.Container.XrayInstanceRepository

	t.Run("FindByID_Success", func(t *testing.T) {
		instance, err := testutil.CreateTestXrayInstanceWithDefaults()
		require.NoError(t, err)

		err = repo.Create(ctx, instance)
		require.NoError(t, err)

		found, err := repo.FindByID(ctx, instance.ID)
		require.NoError(t, err)
		assert.Equal(t, instance.ID, found.ID)
		assert.Equal(t, instance.Name, found.Name)
		assert.Equal(t, instance.Protocol, found.Protocol)
		assert.Equal(t, instance.APIPort, found.APIPort)
	})

	t.Run("FindByID_NotFound", func(t *testing.T) {
		nonExistentID := uuid.New()

		_, err := repo.FindByID(ctx, nonExistentID)
		assert.Error(t, err)
	})
}

func TestXrayInstanceRepository_FindByNodeID(t *testing.T) {
	if !testutil.IsIntegrationTest() {
		t.Skip("Skipping integration test")
	}

	app := testutil.NewTestApp(t)
	defer app.Cleanup()
	defer app.CleanupTables()

	ctx := context.Background()
	repo := app.Container.XrayInstanceRepository

	t.Run("FindByNodeID_Success", func(t *testing.T) {
		instance, err := testutil.CreateTestXrayInstanceWithDefaults()
		require.NoError(t, err)

		err = repo.Create(ctx, instance)
		require.NoError(t, err)

		found, err := repo.FindByNodeID(ctx, instance.NodeID)
		require.NoError(t, err)
		assert.Equal(t, instance.ID, found.ID)
		assert.Equal(t, instance.NodeID, found.NodeID)
	})

	t.Run("FindByNodeID_NotFound", func(t *testing.T) {
		nonExistentNodeID := uuid.New()

		_, err := repo.FindByNodeID(ctx, nonExistentNodeID)
		assert.Error(t, err)
	})
}

func TestXrayInstanceRepository_FindRunning(t *testing.T) {
	if !testutil.IsIntegrationTest() {
		t.Skip("Skipping integration test")
	}

	app := testutil.NewTestApp(t)
	defer app.Cleanup()
	defer app.CleanupTables()

	ctx := context.Background()
	repo := app.Container.XrayInstanceRepository

	t.Run("FindRunning_MultipleInstances", func(t *testing.T) {
		// Create running instance
		runningInstance, err := testutil.CreateTestXrayInstanceWithDefaults()
		require.NoError(t, err)
		runningInstance.Start()
		err = repo.Create(ctx, runningInstance)
		require.NoError(t, err)

		// Create stopped instance
		stoppedInstance, err := testutil.CreateTestXrayInstanceWithDefaults()
		require.NoError(t, err)
		err = repo.Create(ctx, stoppedInstance)
		require.NoError(t, err)

		// Find running instances
		running, err := repo.FindRunning(ctx)
		require.NoError(t, err)
		assert.GreaterOrEqual(t, len(running), 1)
		
		// Verify all returned instances are running
		for _, inst := range running {
			assert.Equal(t, xray.InstanceStatusRunning, inst.Status)
		}
	})
}

func TestXrayInstanceRepository_Update(t *testing.T) {
	if !testutil.IsIntegrationTest() {
		t.Skip("Skipping integration test")
	}

	app := testutil.NewTestApp(t)
	defer app.Cleanup()
	defer app.CleanupTables()

	ctx := context.Background()
	repo := app.Container.XrayInstanceRepository

	t.Run("Update_Success", func(t *testing.T) {
		instance, err := testutil.CreateTestXrayInstanceWithDefaults()
		require.NoError(t, err)

		err = repo.Create(ctx, instance)
		require.NoError(t, err)

		// Update instance
		instance.Start()
		err = repo.Update(ctx, instance)
		require.NoError(t, err)

		// Verify update
		found, err := repo.FindByID(ctx, instance.ID)
		require.NoError(t, err)
		assert.Equal(t, xray.InstanceStatusRunning, found.Status)
	})

	t.Run("Update_NotFound", func(t *testing.T) {
		instance, err := testutil.CreateTestXrayInstanceWithDefaults()
		require.NoError(t, err)
		instance.ID = uuid.New() // Non-existent ID

		err = repo.Update(ctx, instance)
		assert.Error(t, err)
	})
}

func TestXrayInstanceRepository_Delete(t *testing.T) {
	if !testutil.IsIntegrationTest() {
		t.Skip("Skipping integration test")
	}

	app := testutil.NewTestApp(t)
	defer app.Cleanup()
	defer app.CleanupTables()

	ctx := context.Background()
	repo := app.Container.XrayInstanceRepository

	t.Run("Delete_Success", func(t *testing.T) {
		instance, err := testutil.CreateTestXrayInstanceWithDefaults()
		require.NoError(t, err)

		err = repo.Create(ctx, instance)
		require.NoError(t, err)

		err = repo.Delete(ctx, instance.ID)
		require.NoError(t, err)

		// Verify deletion
		_, err = repo.FindByID(ctx, instance.ID)
		assert.Error(t, err)
	})

	t.Run("Delete_NotFound", func(t *testing.T) {
		nonExistentID := uuid.New()

		err := repo.Delete(ctx, nonExistentID)
		assert.Error(t, err)
	})
}

func TestXrayInstanceRepository_List(t *testing.T) {
	if !testutil.IsIntegrationTest() {
		t.Skip("Skipping integration test")
	}

	app := testutil.NewTestApp(t)
	defer app.Cleanup()
	defer app.CleanupTables()

	ctx := context.Background()
	repo := app.Container.XrayInstanceRepository

	tests := []struct {
		name          string
		createCount   int
		offset        int
		limit         int
		expectedCount int
	}{
		{"List_All", 3, 0, 10, 3},
		{"List_WithOffset", 3, 1, 10, 2},
		{"List_WithLimit", 3, 0, 2, 2},
		{"List_WithOffsetAndLimit", 5, 2, 2, 2},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer app.CleanupTables()

			// Create instances
			for i := 0; i < tt.createCount; i++ {
				instance, err := testutil.CreateTestXrayInstanceWithDefaults()
				require.NoError(t, err)
				err = repo.Create(ctx, instance)
				require.NoError(t, err)
			}

			// List instances
			instances, err := repo.List(ctx, tt.offset, tt.limit)
			require.NoError(t, err)
			assert.Equal(t, tt.expectedCount, len(instances))
		})
	}
}

func TestXrayInstanceRepository_Count(t *testing.T) {
	if !testutil.IsIntegrationTest() {
		t.Skip("Skipping integration test")
	}

	app := testutil.NewTestApp(t)
	defer app.Cleanup()
	defer app.CleanupTables()

	ctx := context.Background()
	repo := app.Container.XrayInstanceRepository

	t.Run("Count_EmptyDatabase", func(t *testing.T) {
		count, err := repo.Count(ctx)
		require.NoError(t, err)
		assert.Equal(t, int64(0), count)
	})

	t.Run("Count_WithInstances", func(t *testing.T) {
		defer app.CleanupTables()

		for i := 0; i < 3; i++ {
			instance, err := testutil.CreateTestXrayInstanceWithDefaults()
			require.NoError(t, err)
			err = repo.Create(ctx, instance)
			require.NoError(t, err)
		}

		count, err := repo.Count(ctx)
		require.NoError(t, err)
		assert.Equal(t, int64(3), count)
	})
}

func TestXrayInstanceRepository_ListWithFilters(t *testing.T) {
	if !testutil.IsIntegrationTest() {
		t.Skip("Skipping integration test")
	}

	app := testutil.NewTestApp(t)
	defer app.Cleanup()
	defer app.CleanupTables()

	ctx := context.Background()
	repo := app.Container.XrayInstanceRepository

	// Create test instances
	runningInstance, err := testutil.CreateTestXrayInstanceWithDefaults()
	require.NoError(t, err)
	runningInstance.Start()
	err = repo.Create(ctx, runningInstance)
	require.NoError(t, err)

	stoppedInstance, err := testutil.CreateTestXrayInstanceWithDefaults()
	require.NoError(t, err)
	err = repo.Create(ctx, stoppedInstance)
	require.NoError(t, err)

	t.Run("Filter_ByStatus", func(t *testing.T) {
		runningStatus := xray.InstanceStatusRunning
		filters := xray.XrayInstanceFilters{
			Offset: 0,
			Limit:  10,
			Status: &runningStatus,
		}

		instances, total, err := repo.ListWithFilters(ctx, filters)
		require.NoError(t, err)
		assert.GreaterOrEqual(t, total, int64(1))
		assert.GreaterOrEqual(t, len(instances), 1)

		for _, inst := range instances {
			assert.Equal(t, xray.InstanceStatusRunning, inst.Status)
		}
	})

	t.Run("Filter_WithPagination", func(t *testing.T) {
		filters := xray.XrayInstanceFilters{
			Offset: 0,
			Limit:  1,
		}

		instances, total, err := repo.ListWithFilters(ctx, filters)
		require.NoError(t, err)
		assert.GreaterOrEqual(t, total, int64(2))
		assert.Len(t, instances, 1)
	})
}

