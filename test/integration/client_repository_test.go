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

func TestClientRepository_Create(t *testing.T) {
	if !testutil.IsIntegrationTest() {
		t.Skip("Skipping integration test")
	}

	app := testutil.NewTestApp(t)
	defer app.Cleanup()
	defer app.CleanupTables()

	ctx := context.Background()
	userRepo := app.Container.UserRepository
	instanceRepo := app.Container.XrayInstanceRepository
	inboundRepo := app.Container.InboundRepository
	clientRepo := app.Container.ClientRepository

	t.Run("Create_Success", func(t *testing.T) {
		// Create dependencies - server and node first
		testServer, err := testutil.CreateTestServerWithDefaults()
		require.NoError(t, err)
		err = app.Container.ServerRepository.Create(ctx, testServer)
		require.NoError(t, err)

		testNode, err := testutil.CreateTestNodeWithDefaults(testServer.ID)
		require.NoError(t, err)
		err = app.Container.NodeRepository.Create(ctx, testNode)
		require.NoError(t, err)

		// Create user
		user, err := testutil.CreateTestUserWithDefaults()
		require.NoError(t, err)
		err = userRepo.Create(ctx, user)
		require.NoError(t, err)

		// Create instance and inbound
		instance, err := testutil.CreateTestXrayInstanceWithDefaults(testNode.ID)
		require.NoError(t, err)
		err = instanceRepo.Create(ctx, instance)
		require.NoError(t, err)

		inbound, err := testutil.CreateTestInboundWithDefaults(instance.ID)
		require.NoError(t, err)
		err = inboundRepo.Create(ctx, inbound)
		require.NoError(t, err)

		// Create client
		client, err := testutil.CreateTestClientWithDefaults(inbound.ID, user.ID)
		require.NoError(t, err)

		err = clientRepo.Create(ctx, client)
		require.NoError(t, err)

		// Verify creation
		found, err := clientRepo.FindByID(ctx, client.ID)
		require.NoError(t, err)
		assert.Equal(t, client.ID, found.ID)
		assert.Equal(t, client.InboundID, found.InboundID)
		assert.Equal(t, client.UserID, found.UserID)
	})
}

func TestClientRepository_FindByID(t *testing.T) {
	if !testutil.IsIntegrationTest() {
		t.Skip("Skipping integration test")
	}

	app := testutil.NewTestApp(t)
	defer app.Cleanup()
	defer app.CleanupTables()

	ctx := context.Background()
	userRepo := app.Container.UserRepository
	instanceRepo := app.Container.XrayInstanceRepository
	inboundRepo := app.Container.InboundRepository
	clientRepo := app.Container.ClientRepository

	t.Run("FindByID_Success", func(t *testing.T) {
		// Create dependencies - server and node first
		testServer, err := testutil.CreateTestServerWithDefaults()
		require.NoError(t, err)
		err = app.Container.ServerRepository.Create(ctx, testServer)
		require.NoError(t, err)

		testNode, err := testutil.CreateTestNodeWithDefaults(testServer.ID)
		require.NoError(t, err)
		err = app.Container.NodeRepository.Create(ctx, testNode)
		require.NoError(t, err)

		user, err := testutil.CreateTestUserWithDefaults()
		require.NoError(t, err)
		err = userRepo.Create(ctx, user)
		require.NoError(t, err)

		instance, err := testutil.CreateTestXrayInstanceWithDefaults(testNode.ID)
		require.NoError(t, err)
		err = instanceRepo.Create(ctx, instance)
		require.NoError(t, err)

		inbound, err := testutil.CreateTestInboundWithDefaults(instance.ID)
		require.NoError(t, err)
		err = inboundRepo.Create(ctx, inbound)
		require.NoError(t, err)

		client, err := testutil.CreateTestClientWithDefaults(inbound.ID, user.ID)
		require.NoError(t, err)
		err = clientRepo.Create(ctx, client)
		require.NoError(t, err)

		found, err := clientRepo.FindByID(ctx, client.ID)
		require.NoError(t, err)
		assert.Equal(t, client.ID, found.ID)
		assert.Equal(t, client.UUID, found.UUID)
		assert.Equal(t, client.Flow, found.Flow)
	})

	t.Run("FindByID_NotFound", func(t *testing.T) {
		nonExistentID := uuid.New()

		_, err := clientRepo.FindByID(ctx, nonExistentID)
		assert.Error(t, err)
	})
}

func TestClientRepository_FindByInboundID(t *testing.T) {
	if !testutil.IsIntegrationTest() {
		t.Skip("Skipping integration test")
	}

	app := testutil.NewTestApp(t)
	defer app.Cleanup()
	defer app.CleanupTables()

	ctx := context.Background()
	userRepo := app.Container.UserRepository
	instanceRepo := app.Container.XrayInstanceRepository
	inboundRepo := app.Container.InboundRepository
	clientRepo := app.Container.ClientRepository

	t.Run("FindByInboundID_MultipleClients", func(t *testing.T) {
		// Create dependencies - server and node first
		testServer, err := testutil.CreateTestServerWithDefaults()
		require.NoError(t, err)
		err = app.Container.ServerRepository.Create(ctx, testServer)
		require.NoError(t, err)

		testNode, err := testutil.CreateTestNodeWithDefaults(testServer.ID)
		require.NoError(t, err)
		err = app.Container.NodeRepository.Create(ctx, testNode)
		require.NoError(t, err)

		user, err := testutil.CreateTestUserWithDefaults()
		require.NoError(t, err)
		err = userRepo.Create(ctx, user)
		require.NoError(t, err)

		instance, err := testutil.CreateTestXrayInstanceWithDefaults(testNode.ID)
		require.NoError(t, err)
		err = instanceRepo.Create(ctx, instance)
		require.NoError(t, err)

		inbound, err := testutil.CreateTestInboundWithDefaults(instance.ID)
		require.NoError(t, err)
		err = inboundRepo.Create(ctx, inbound)
		require.NoError(t, err)

		// Create multiple clients
		for i := 0; i < 3; i++ {
			client, err := testutil.CreateTestClientWithDefaults(inbound.ID, user.ID)
			require.NoError(t, err)
			err = clientRepo.Create(ctx, client)
			require.NoError(t, err)
		}

		clients, err := clientRepo.FindByInboundID(ctx, inbound.ID)
		require.NoError(t, err)
		assert.Len(t, clients, 3)

		for _, client := range clients {
			assert.Equal(t, inbound.ID, client.InboundID)
		}
	})
}

func TestClientRepository_FindByUserID(t *testing.T) {
	if !testutil.IsIntegrationTest() {
		t.Skip("Skipping integration test")
	}

	app := testutil.NewTestApp(t)
	defer app.Cleanup()
	defer app.CleanupTables()

	ctx := context.Background()
	userRepo := app.Container.UserRepository
	instanceRepo := app.Container.XrayInstanceRepository
	inboundRepo := app.Container.InboundRepository
	clientRepo := app.Container.ClientRepository

	t.Run("FindByUserID_MultipleClients", func(t *testing.T) {
		user, err := testutil.CreateTestUserWithDefaults()
		require.NoError(t, err)
		err = userRepo.Create(ctx, user)
		require.NoError(t, err)

		// Create dependencies - server and node first
		testServer, err := testutil.CreateTestServerWithDefaults()
		require.NoError(t, err)
		err = app.Container.ServerRepository.Create(ctx, testServer)
		require.NoError(t, err)

		testNode, err := testutil.CreateTestNodeWithDefaults(testServer.ID)
		require.NoError(t, err)
		err = app.Container.NodeRepository.Create(ctx, testNode)
		require.NoError(t, err)

		instance, err := testutil.CreateTestXrayInstanceWithDefaults(testNode.ID)
		require.NoError(t, err)
		err = instanceRepo.Create(ctx, instance)
		require.NoError(t, err)

		inbound, err := testutil.CreateTestInboundWithDefaults(instance.ID)
		require.NoError(t, err)
		err = inboundRepo.Create(ctx, inbound)
		require.NoError(t, err)

		// Create multiple clients for same user
		for i := 0; i < 2; i++ {
			client, err := testutil.CreateTestClientWithDefaults(inbound.ID, user.ID)
			require.NoError(t, err)
			err = clientRepo.Create(ctx, client)
			require.NoError(t, err)
		}

		clients, err := clientRepo.FindByUserID(ctx, user.ID)
		require.NoError(t, err)
		assert.Len(t, clients, 2)

		for _, client := range clients {
			assert.Equal(t, user.ID, client.UserID)
		}
	})
}

func TestClientRepository_FindByUUID(t *testing.T) {
	if !testutil.IsIntegrationTest() {
		t.Skip("Skipping integration test")
	}

	app := testutil.NewTestApp(t)
	defer app.Cleanup()
	defer app.CleanupTables()

	ctx := context.Background()
	userRepo := app.Container.UserRepository
	instanceRepo := app.Container.XrayInstanceRepository
	inboundRepo := app.Container.InboundRepository
	clientRepo := app.Container.ClientRepository

	t.Run("FindByUUID_Success", func(t *testing.T) {
		user, err := testutil.CreateTestUserWithDefaults()
		require.NoError(t, err)
		err = userRepo.Create(ctx, user)
		require.NoError(t, err)

		// Create dependencies - server and node first
		testServer, err := testutil.CreateTestServerWithDefaults()
		require.NoError(t, err)
		err = app.Container.ServerRepository.Create(ctx, testServer)
		require.NoError(t, err)

		testNode, err := testutil.CreateTestNodeWithDefaults(testServer.ID)
		require.NoError(t, err)
		err = app.Container.NodeRepository.Create(ctx, testNode)
		require.NoError(t, err)

		instance, err := testutil.CreateTestXrayInstanceWithDefaults(testNode.ID)
		require.NoError(t, err)
		err = instanceRepo.Create(ctx, instance)
		require.NoError(t, err)

		inbound, err := testutil.CreateTestInboundWithDefaults(instance.ID)
		require.NoError(t, err)
		err = inboundRepo.Create(ctx, inbound)
		require.NoError(t, err)

		client, err := testutil.CreateTestClientWithDefaults(inbound.ID, user.ID)
		require.NoError(t, err)
		err = clientRepo.Create(ctx, client)
		require.NoError(t, err)

		found, err := clientRepo.FindByUUID(ctx, client.UUID)
		require.NoError(t, err)
		assert.Equal(t, client.ID, found.ID)
		assert.Equal(t, client.UUID, found.UUID)
	})

	t.Run("FindByUUID_NotFound", func(t *testing.T) {
		nonExistentUUID := uuid.New().String()

		_, err := clientRepo.FindByUUID(ctx, nonExistentUUID)
		assert.Error(t, err)
	})
}

func TestClientRepository_FindEnabledByInboundID(t *testing.T) {
	if !testutil.IsIntegrationTest() {
		t.Skip("Skipping integration test")
	}

	app := testutil.NewTestApp(t)
	defer app.Cleanup()
	defer app.CleanupTables()

	ctx := context.Background()
	userRepo := app.Container.UserRepository
	instanceRepo := app.Container.XrayInstanceRepository
	inboundRepo := app.Container.InboundRepository
	clientRepo := app.Container.ClientRepository

	t.Run("FindEnabledByInboundID_FilterDisabled", func(t *testing.T) {
		user, err := testutil.CreateTestUserWithDefaults()
		require.NoError(t, err)
		err = userRepo.Create(ctx, user)
		require.NoError(t, err)

		instance, err := testutil.CreateTestXrayInstanceWithDefaults()
		require.NoError(t, err)
		err = instanceRepo.Create(ctx, instance)
		require.NoError(t, err)

		inbound, err := testutil.CreateTestInboundWithDefaults(instance.ID)
		require.NoError(t, err)
		err = inboundRepo.Create(ctx, inbound)
		require.NoError(t, err)

		// Create enabled client
		enabledClient, err := testutil.CreateTestClientWithDefaults(inbound.ID, user.ID)
		require.NoError(t, err)
		err = clientRepo.Create(ctx, enabledClient)
		require.NoError(t, err)

		// Create disabled client
		disabledClient, err := testutil.CreateTestClientWithDefaults(inbound.ID, user.ID)
		require.NoError(t, err)
		disabledClient.Disable()
		err = clientRepo.Create(ctx, disabledClient)
		require.NoError(t, err)

		// Find only enabled
		clients, err := clientRepo.FindEnabledByInboundID(ctx, inbound.ID)
		require.NoError(t, err)
		assert.GreaterOrEqual(t, len(clients), 1)

		for _, client := range clients {
			assert.True(t, client.IsEnabled())
		}
	})
}

func TestClientRepository_Update(t *testing.T) {
	if !testutil.IsIntegrationTest() {
		t.Skip("Skipping integration test")
	}

	app := testutil.NewTestApp(t)
	defer app.Cleanup()
	defer app.CleanupTables()

	ctx := context.Background()
	userRepo := app.Container.UserRepository
	instanceRepo := app.Container.XrayInstanceRepository
	inboundRepo := app.Container.InboundRepository
	clientRepo := app.Container.ClientRepository

	t.Run("Update_Success", func(t *testing.T) {
		user, err := testutil.CreateTestUserWithDefaults()
		require.NoError(t, err)
		err = userRepo.Create(ctx, user)
		require.NoError(t, err)

		instance, err := testutil.CreateTestXrayInstanceWithDefaults()
		require.NoError(t, err)
		err = instanceRepo.Create(ctx, instance)
		require.NoError(t, err)

		inbound, err := testutil.CreateTestInboundWithDefaults(instance.ID)
		require.NoError(t, err)
		err = inboundRepo.Create(ctx, inbound)
		require.NoError(t, err)

		client, err := testutil.CreateTestClientWithDefaults(inbound.ID, user.ID)
		require.NoError(t, err)
		err = clientRepo.Create(ctx, client)
		require.NoError(t, err)

		// Update client
		client.Disable()
		err = clientRepo.Update(ctx, client)
		require.NoError(t, err)

		// Verify update
		found, err := clientRepo.FindByID(ctx, client.ID)
		require.NoError(t, err)
		assert.False(t, found.IsEnabled())
	})

	t.Run("Update_NotFound", func(t *testing.T) {
		inboundID := uuid.New()
		userID := uuid.New()

		client, err := testutil.CreateTestClientWithDefaults(inboundID, userID)
		require.NoError(t, err)
		client.ID = uuid.New() // Non-existent ID

		err = clientRepo.Update(ctx, client)
		assert.Error(t, err)
	})
}

func TestClientRepository_Delete(t *testing.T) {
	if !testutil.IsIntegrationTest() {
		t.Skip("Skipping integration test")
	}

	app := testutil.NewTestApp(t)
	defer app.Cleanup()
	defer app.CleanupTables()

	ctx := context.Background()
	userRepo := app.Container.UserRepository
	instanceRepo := app.Container.XrayInstanceRepository
	inboundRepo := app.Container.InboundRepository
	clientRepo := app.Container.ClientRepository

	t.Run("Delete_Success", func(t *testing.T) {
		user, err := testutil.CreateTestUserWithDefaults()
		require.NoError(t, err)
		err = userRepo.Create(ctx, user)
		require.NoError(t, err)

		instance, err := testutil.CreateTestXrayInstanceWithDefaults()
		require.NoError(t, err)
		err = instanceRepo.Create(ctx, instance)
		require.NoError(t, err)

		inbound, err := testutil.CreateTestInboundWithDefaults(instance.ID)
		require.NoError(t, err)
		err = inboundRepo.Create(ctx, inbound)
		require.NoError(t, err)

		client, err := testutil.CreateTestClientWithDefaults(inbound.ID, user.ID)
		require.NoError(t, err)
		err = clientRepo.Create(ctx, client)
		require.NoError(t, err)

		err = clientRepo.Delete(ctx, client.ID)
		require.NoError(t, err)

		// Verify deletion
		_, err = clientRepo.FindByID(ctx, client.ID)
		assert.Error(t, err)
	})

	t.Run("Delete_NotFound", func(t *testing.T) {
		nonExistentID := uuid.New()

		err := clientRepo.Delete(ctx, nonExistentID)
		assert.Error(t, err)
	})
}

func TestClientRepository_List(t *testing.T) {
	if !testutil.IsIntegrationTest() {
		t.Skip("Skipping integration test")
	}

	app := testutil.NewTestApp(t)
	defer app.Cleanup()
	defer app.CleanupTables()

	ctx := context.Background()
	userRepo := app.Container.UserRepository
	instanceRepo := app.Container.XrayInstanceRepository
	inboundRepo := app.Container.InboundRepository
	clientRepo := app.Container.ClientRepository

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

			user, err := testutil.CreateTestUserWithDefaults()
			require.NoError(t, err)
			err = userRepo.Create(ctx, user)
			require.NoError(t, err)

			instance, err := testutil.CreateTestXrayInstanceWithDefaults()
			require.NoError(t, err)
			err = instanceRepo.Create(ctx, instance)
			require.NoError(t, err)

			inbound, err := testutil.CreateTestInboundWithDefaults(instance.ID)
			require.NoError(t, err)
			err = inboundRepo.Create(ctx, inbound)
			require.NoError(t, err)

			// Create clients
			for i := 0; i < tt.createCount; i++ {
				client, err := testutil.CreateTestClientWithDefaults(inbound.ID, user.ID)
				require.NoError(t, err)
				err = clientRepo.Create(ctx, client)
				require.NoError(t, err)
			}

			// List clients
			clients, err := clientRepo.List(ctx, tt.offset, tt.limit)
			require.NoError(t, err)
			assert.Equal(t, tt.expectedCount, len(clients))
		})
	}
}

func TestClientRepository_Count(t *testing.T) {
	if !testutil.IsIntegrationTest() {
		t.Skip("Skipping integration test")
	}

	app := testutil.NewTestApp(t)
	defer app.Cleanup()
	defer app.CleanupTables()

	ctx := context.Background()
	userRepo := app.Container.UserRepository
	instanceRepo := app.Container.XrayInstanceRepository
	inboundRepo := app.Container.InboundRepository
	clientRepo := app.Container.ClientRepository

	t.Run("Count_WithClients", func(t *testing.T) {
		user, err := testutil.CreateTestUserWithDefaults()
		require.NoError(t, err)
		err = userRepo.Create(ctx, user)
		require.NoError(t, err)

		instance, err := testutil.CreateTestXrayInstanceWithDefaults()
		require.NoError(t, err)
		err = instanceRepo.Create(ctx, instance)
		require.NoError(t, err)

		inbound, err := testutil.CreateTestInboundWithDefaults(instance.ID)
		require.NoError(t, err)
		err = inboundRepo.Create(ctx, inbound)
		require.NoError(t, err)

		for i := 0; i < 3; i++ {
			client, err := testutil.CreateTestClientWithDefaults(inbound.ID, user.ID)
			require.NoError(t, err)
			err = clientRepo.Create(ctx, client)
			require.NoError(t, err)
		}

		count, err := clientRepo.Count(ctx)
		require.NoError(t, err)
		assert.Equal(t, int64(3), count)
	})
}

func TestClientRepository_ListWithFilters(t *testing.T) {
	if !testutil.IsIntegrationTest() {
		t.Skip("Skipping integration test")
	}

	app := testutil.NewTestApp(t)
	defer app.Cleanup()
	defer app.CleanupTables()

	ctx := context.Background()
	userRepo := app.Container.UserRepository
	instanceRepo := app.Container.XrayInstanceRepository
	inboundRepo := app.Container.InboundRepository
	clientRepo := app.Container.ClientRepository

	user, err := testutil.CreateTestUserWithDefaults()
	require.NoError(t, err)
	err = userRepo.Create(ctx, user)
	require.NoError(t, err)

	instance, err := testutil.CreateTestXrayInstanceWithDefaults()
	require.NoError(t, err)
	err = instanceRepo.Create(ctx, instance)
	require.NoError(t, err)

	inbound, err := testutil.CreateTestInboundWithDefaults(instance.ID)
	require.NoError(t, err)
	err = inboundRepo.Create(ctx, inbound)
	require.NoError(t, err)

	// Create enabled client
	enabledClient, err := testutil.CreateTestClientWithDefaults(inbound.ID, user.ID)
	require.NoError(t, err)
	err = clientRepo.Create(ctx, enabledClient)
	require.NoError(t, err)

	// Create disabled client
	disabledClient, err := testutil.CreateTestClientWithDefaults(inbound.ID, user.ID)
	require.NoError(t, err)
	disabledClient.Disable()
	err = clientRepo.Create(ctx, disabledClient)
	require.NoError(t, err)

	t.Run("Filter_ByEnabled", func(t *testing.T) {
		enabled := true
		filters := xray.ClientFilters{
			Offset:  0,
			Limit:   10,
			Enabled: &enabled,
		}

		clients, total, err := clientRepo.ListWithFilters(ctx, filters)
		require.NoError(t, err)
		assert.GreaterOrEqual(t, total, int64(1))

		for _, client := range clients {
			assert.True(t, client.IsEnabled())
		}
	})

	t.Run("Filter_ByUserID", func(t *testing.T) {
		filters := xray.ClientFilters{
			Offset: 0,
			Limit:  10,
			UserID: &user.ID,
		}

		clients, total, err := clientRepo.ListWithFilters(ctx, filters)
		require.NoError(t, err)
		assert.GreaterOrEqual(t, total, int64(2))

		for _, client := range clients {
			assert.Equal(t, user.ID, client.UserID)
		}
	})

	t.Run("Filter_ByInboundID", func(t *testing.T) {
		filters := xray.ClientFilters{
			Offset:    0,
			Limit:     10,
			InboundID: &inbound.ID,
		}

		clients, total, err := clientRepo.ListWithFilters(ctx, filters)
		require.NoError(t, err)
		assert.GreaterOrEqual(t, total, int64(2))

		for _, client := range clients {
			assert.Equal(t, inbound.ID, client.InboundID)
		}
	})
}

