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

func TestInboundRepository_Create(t *testing.T) {
	if !testutil.IsIntegrationTest() {
		t.Skip("Skipping integration test")
	}

	app := testutil.NewTestApp(t)
	defer app.Cleanup()
	defer app.CleanupTables()

	ctx := context.Background()
	instanceRepo := app.Container.XrayInstanceRepository
	inboundRepo := app.Container.InboundRepository

	t.Run("Create_Success", func(t *testing.T) {
		// Create instance first
		instance, err := testutil.CreateTestXrayInstanceWithDefaults()
		require.NoError(t, err)
		err = instanceRepo.Create(ctx, instance)
		require.NoError(t, err)

		// Create inbound
		inbound, err := testutil.CreateTestInboundWithDefaults(instance.ID)
		require.NoError(t, err)

		err = inboundRepo.Create(ctx, inbound)
		require.NoError(t, err)

		// Verify creation
		found, err := inboundRepo.FindByID(ctx, inbound.ID)
		require.NoError(t, err)
		assert.Equal(t, inbound.ID, found.ID)
		assert.Equal(t, inbound.InstanceID, found.InstanceID)
		assert.Equal(t, inbound.Protocol, found.Protocol)
	})
}

func TestInboundRepository_FindByID(t *testing.T) {
	if !testutil.IsIntegrationTest() {
		t.Skip("Skipping integration test")
	}

	app := testutil.NewTestApp(t)
	defer app.Cleanup()
	defer app.CleanupTables()

	ctx := context.Background()
	instanceRepo := app.Container.XrayInstanceRepository
	inboundRepo := app.Container.InboundRepository

	t.Run("FindByID_Success", func(t *testing.T) {
		instance, err := testutil.CreateTestXrayInstanceWithDefaults()
		require.NoError(t, err)
		err = instanceRepo.Create(ctx, instance)
		require.NoError(t, err)

		inbound, err := testutil.CreateTestInboundWithDefaults(instance.ID)
		require.NoError(t, err)
		err = inboundRepo.Create(ctx, inbound)
		require.NoError(t, err)

		found, err := inboundRepo.FindByID(ctx, inbound.ID)
		require.NoError(t, err)
		assert.Equal(t, inbound.ID, found.ID)
		assert.Equal(t, inbound.Port, found.Port)
		assert.Equal(t, inbound.Network, found.Network)
	})

	t.Run("FindByID_NotFound", func(t *testing.T) {
		nonExistentID := uuid.New()

		_, err := inboundRepo.FindByID(ctx, nonExistentID)
		assert.Error(t, err)
	})
}

func TestInboundRepository_FindByInstanceID(t *testing.T) {
	if !testutil.IsIntegrationTest() {
		t.Skip("Skipping integration test")
	}

	app := testutil.NewTestApp(t)
	defer app.Cleanup()
	defer app.CleanupTables()

	ctx := context.Background()
	instanceRepo := app.Container.XrayInstanceRepository
	inboundRepo := app.Container.InboundRepository

	t.Run("FindByInstanceID_MultipleInbounds", func(t *testing.T) {
		instance, err := testutil.CreateTestXrayInstanceWithDefaults()
		require.NoError(t, err)
		err = instanceRepo.Create(ctx, instance)
		require.NoError(t, err)

		// Create multiple inbounds
		for i := 0; i < 3; i++ {
			inbound, err := testutil.CreateTestInboundWithDefaults(instance.ID)
			require.NoError(t, err)
			err = inboundRepo.Create(ctx, inbound)
			require.NoError(t, err)
		}

		inbounds, err := inboundRepo.FindByInstanceID(ctx, instance.ID)
		require.NoError(t, err)
		assert.Len(t, inbounds, 3)

		for _, inbound := range inbounds {
			assert.Equal(t, instance.ID, inbound.InstanceID)
		}
	})

	t.Run("FindByInstanceID_Empty", func(t *testing.T) {
		nonExistentInstanceID := uuid.New()

		inbounds, err := inboundRepo.FindByInstanceID(ctx, nonExistentInstanceID)
		require.NoError(t, err)
		assert.Empty(t, inbounds)
	})
}

func TestInboundRepository_FindEnabledByInstanceID(t *testing.T) {
	if !testutil.IsIntegrationTest() {
		t.Skip("Skipping integration test")
	}

	app := testutil.NewTestApp(t)
	defer app.Cleanup()
	defer app.CleanupTables()

	ctx := context.Background()
	instanceRepo := app.Container.XrayInstanceRepository
	inboundRepo := app.Container.InboundRepository

	t.Run("FindEnabledByInstanceID_FilterDisabled", func(t *testing.T) {
		instance, err := testutil.CreateTestXrayInstanceWithDefaults()
		require.NoError(t, err)
		err = instanceRepo.Create(ctx, instance)
		require.NoError(t, err)

		// Create enabled inbound
		enabledInbound, err := testutil.CreateTestInboundWithDefaults(instance.ID)
		require.NoError(t, err)
		err = inboundRepo.Create(ctx, enabledInbound)
		require.NoError(t, err)

		// Create disabled inbound
		disabledInbound, err := testutil.CreateTestInboundWithDefaults(instance.ID)
		require.NoError(t, err)
		disabledInbound.Disable()
		err = inboundRepo.Create(ctx, disabledInbound)
		require.NoError(t, err)

		// Find only enabled
		inbounds, err := inboundRepo.FindEnabledByInstanceID(ctx, instance.ID)
		require.NoError(t, err)
		assert.GreaterOrEqual(t, len(inbounds), 1)

		for _, inbound := range inbounds {
			assert.True(t, inbound.IsEnabled)
		}
	})
}

func TestInboundRepository_Update(t *testing.T) {
	if !testutil.IsIntegrationTest() {
		t.Skip("Skipping integration test")
	}

	app := testutil.NewTestApp(t)
	defer app.Cleanup()
	defer app.CleanupTables()

	ctx := context.Background()
	instanceRepo := app.Container.XrayInstanceRepository
	inboundRepo := app.Container.InboundRepository

	t.Run("Update_Success", func(t *testing.T) {
		instance, err := testutil.CreateTestXrayInstanceWithDefaults()
		require.NoError(t, err)
		err = instanceRepo.Create(ctx, instance)
		require.NoError(t, err)

		inbound, err := testutil.CreateTestInboundWithDefaults(instance.ID)
		require.NoError(t, err)
		err = inboundRepo.Create(ctx, inbound)
		require.NoError(t, err)

		// Update inbound
		inbound.Disable()
		err = inboundRepo.Update(ctx, inbound)
		require.NoError(t, err)

		// Verify update
		found, err := inboundRepo.FindByID(ctx, inbound.ID)
		require.NoError(t, err)
		assert.False(t, found.IsEnabled)
	})

	t.Run("Update_NotFound", func(t *testing.T) {
		instance, err := testutil.CreateTestXrayInstanceWithDefaults()
		require.NoError(t, err)

		inbound, err := testutil.CreateTestInboundWithDefaults(instance.ID)
		require.NoError(t, err)
		inbound.ID = uuid.New() // Non-existent ID

		err = inboundRepo.Update(ctx, inbound)
		assert.Error(t, err)
	})
}

func TestInboundRepository_Delete(t *testing.T) {
	if !testutil.IsIntegrationTest() {
		t.Skip("Skipping integration test")
	}

	app := testutil.NewTestApp(t)
	defer app.Cleanup()
	defer app.CleanupTables()

	ctx := context.Background()
	instanceRepo := app.Container.XrayInstanceRepository
	inboundRepo := app.Container.InboundRepository

	t.Run("Delete_Success", func(t *testing.T) {
		instance, err := testutil.CreateTestXrayInstanceWithDefaults()
		require.NoError(t, err)
		err = instanceRepo.Create(ctx, instance)
		require.NoError(t, err)

		inbound, err := testutil.CreateTestInboundWithDefaults(instance.ID)
		require.NoError(t, err)
		err = inboundRepo.Create(ctx, inbound)
		require.NoError(t, err)

		err = inboundRepo.Delete(ctx, inbound.ID)
		require.NoError(t, err)

		// Verify deletion
		_, err = inboundRepo.FindByID(ctx, inbound.ID)
		assert.Error(t, err)
	})

	t.Run("Delete_NotFound", func(t *testing.T) {
		nonExistentID := uuid.New()

		err := inboundRepo.Delete(ctx, nonExistentID)
		assert.Error(t, err)
	})
}

func TestInboundRepository_List(t *testing.T) {
	if !testutil.IsIntegrationTest() {
		t.Skip("Skipping integration test")
	}

	app := testutil.NewTestApp(t)
	defer app.Cleanup()
	defer app.CleanupTables()

	ctx := context.Background()
	instanceRepo := app.Container.XrayInstanceRepository
	inboundRepo := app.Container.InboundRepository

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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer app.CleanupTables()

			instance, err := testutil.CreateTestXrayInstanceWithDefaults()
			require.NoError(t, err)
			err = instanceRepo.Create(ctx, instance)
			require.NoError(t, err)

			// Create inbounds
			for i := 0; i < tt.createCount; i++ {
				inbound, err := testutil.CreateTestInboundWithDefaults(instance.ID)
				require.NoError(t, err)
				err = inboundRepo.Create(ctx, inbound)
				require.NoError(t, err)
			}

			// List inbounds
			inbounds, err := inboundRepo.List(ctx, tt.offset, tt.limit)
			require.NoError(t, err)
			assert.Equal(t, tt.expectedCount, len(inbounds))
		})
	}
}

func TestInboundRepository_Count(t *testing.T) {
	if !testutil.IsIntegrationTest() {
		t.Skip("Skipping integration test")
	}

	app := testutil.NewTestApp(t)
	defer app.Cleanup()
	defer app.CleanupTables()

	ctx := context.Background()
	instanceRepo := app.Container.XrayInstanceRepository
	inboundRepo := app.Container.InboundRepository

	t.Run("Count_WithInbounds", func(t *testing.T) {
		instance, err := testutil.CreateTestXrayInstanceWithDefaults()
		require.NoError(t, err)
		err = instanceRepo.Create(ctx, instance)
		require.NoError(t, err)

		for i := 0; i < 3; i++ {
			inbound, err := testutil.CreateTestInboundWithDefaults(instance.ID)
			require.NoError(t, err)
			err = inboundRepo.Create(ctx, inbound)
			require.NoError(t, err)
		}

		count, err := inboundRepo.Count(ctx)
		require.NoError(t, err)
		assert.Equal(t, int64(3), count)
	})
}

func TestInboundRepository_ListWithFilters(t *testing.T) {
	if !testutil.IsIntegrationTest() {
		t.Skip("Skipping integration test")
	}

	app := testutil.NewTestApp(t)
	defer app.Cleanup()
	defer app.CleanupTables()

	ctx := context.Background()
	instanceRepo := app.Container.XrayInstanceRepository
	inboundRepo := app.Container.InboundRepository

	instance, err := testutil.CreateTestXrayInstanceWithDefaults()
	require.NoError(t, err)
	err = instanceRepo.Create(ctx, instance)
	require.NoError(t, err)

	// Create enabled inbound
	enabledInbound, err := testutil.CreateTestInboundWithDefaults(instance.ID)
	require.NoError(t, err)
	err = inboundRepo.Create(ctx, enabledInbound)
	require.NoError(t, err)

	// Create disabled inbound
	disabledInbound, err := testutil.CreateTestInboundWithDefaults(instance.ID)
	require.NoError(t, err)
	disabledInbound.Disable()
	err = inboundRepo.Create(ctx, disabledInbound)
	require.NoError(t, err)

	t.Run("Filter_ByEnabled", func(t *testing.T) {
		enabled := true
		filters := xray.InboundFilters{
			Offset:  0,
			Limit:   10,
			Enabled: &enabled,
		}

		inbounds, total, err := inboundRepo.ListWithFilters(ctx, filters)
		require.NoError(t, err)
		assert.GreaterOrEqual(t, total, int64(1))

		for _, inbound := range inbounds {
			assert.True(t, inbound.IsEnabled)
		}
	})

	t.Run("Filter_ByInstanceID", func(t *testing.T) {
		filters := xray.InboundFilters{
			Offset:     0,
			Limit:      10,
			InstanceID: &instance.ID,
		}

		inbounds, total, err := inboundRepo.ListWithFilters(ctx, filters)
		require.NoError(t, err)
		assert.GreaterOrEqual(t, total, int64(2))

		for _, inbound := range inbounds {
			assert.Equal(t, instance.ID, inbound.InstanceID)
		}
	})
}

