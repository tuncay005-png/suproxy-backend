package integration_test

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/suproxy/backend/internal/domain/xray"
	"github.com/suproxy/backend/internal/infrastructure/testutil"
	xrayConfig "github.com/suproxy/backend/internal/infrastructure/xray/config"
)

func TestConfigGenerator_Generate(t *testing.T) {
	if !testutil.IsIntegrationTest() {
		t.Skip("Skipping integration test")
	}

	app := testutil.NewTestApp(t)
	defer app.Cleanup()
	defer app.CleanupTables()

	ctx := context.Background()
	instanceRepo := app.Container.XrayInstanceRepository
	inboundRepo := app.Container.InboundRepository
	clientRepo := app.Container.ClientRepository
	realityRepo := app.Container.RealityConfigRepository
	userRepo := app.Container.UserRepository

	generator := xrayConfig.NewGenerator(instanceRepo, inboundRepo, clientRepo, realityRepo)

	t.Run("Generate_Success_WithClients", func(t *testing.T) {
		// Create test data
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
		// Enable client so it appears in config
		err = client.Enable()
		require.NoError(t, err)
		err = clientRepo.Create(ctx, client)
		require.NoError(t, err)

		// Generate config
		config, err := generator.Generate(ctx, instance.ID)
		require.NoError(t, err)
		assert.NotNil(t, config)
		assert.NotNil(t, config.Log)
		assert.NotNil(t, config.API)
		assert.NotNil(t, config.DNS)
		assert.NotNil(t, config.Stats)
		assert.NotNil(t, config.Policy)
		assert.NotEmpty(t, config.Inbounds)
		assert.NotEmpty(t, config.Outbounds)
		assert.NotNil(t, config.Routing)

		// Verify inbounds (API + user inbound)
		assert.GreaterOrEqual(t, len(config.Inbounds), 2)
	})

	t.Run("Generate_Success_NoClients", func(t *testing.T) {
		defer app.CleanupTables()

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

		// No clients created

		config, err := generator.Generate(ctx, instance.ID)
		require.NoError(t, err)
		assert.NotNil(t, config)

		// Should have API inbound only
		assert.Len(t, config.Inbounds, 1)
	})

	t.Run("Generate_Failure_InstanceNotFound", func(t *testing.T) {
		defer app.CleanupTables()

		// Create dependencies - server and node first
		testServer, err := testutil.CreateTestServerWithDefaults()
		require.NoError(t, err)

		testNode, err := testutil.CreateTestNodeWithDefaults(testServer.ID)
		require.NoError(t, err)

		nonExistentInstance, err := testutil.CreateTestXrayInstanceWithDefaults(testNode.ID)
		require.NoError(t, err)
		require.NotNil(t, nonExistentInstance)

		_, err = generator.Generate(ctx, nonExistentInstance.ID)
		assert.Error(t, err)
	})
}

func TestConfigGenerator_GenerateJSON(t *testing.T) {
	if !testutil.IsIntegrationTest() {
		t.Skip("Skipping integration test")
	}

	app := testutil.NewTestApp(t)
	defer app.Cleanup()
	defer app.CleanupTables()

	ctx := context.Background()
	instanceRepo := app.Container.XrayInstanceRepository
	inboundRepo := app.Container.InboundRepository
	clientRepo := app.Container.ClientRepository
	realityRepo := app.Container.RealityConfigRepository
	userRepo := app.Container.UserRepository

	generator := xrayConfig.NewGenerator(instanceRepo, inboundRepo, clientRepo, realityRepo)

	t.Run("GenerateJSON_Success", func(t *testing.T) {
		// Create test data
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
		// Enable client so it appears in config
		err = client.Enable()
		require.NoError(t, err)
		err = clientRepo.Create(ctx, client)
		require.NoError(t, err)

		// Generate JSON
		jsonData, err := generator.GenerateJSON(ctx, instance.ID)
		require.NoError(t, err)
		assert.NotEmpty(t, jsonData)

		// Verify JSON is valid
		var config map[string]interface{}
		err = json.Unmarshal(jsonData, &config)
		require.NoError(t, err)

		// Verify structure
		assert.Contains(t, config, "log")
		assert.Contains(t, config, "api")
		assert.Contains(t, config, "inbounds")
		assert.Contains(t, config, "outbounds")
		assert.Contains(t, config, "routing")
	})
}

func TestConfigGenerator_GenerateInbound(t *testing.T) {
	if !testutil.IsIntegrationTest() {
		t.Skip("Skipping integration test")
	}

	app := testutil.NewTestApp(t)
	defer app.Cleanup()
	defer app.CleanupTables()

	ctx := context.Background()
	instanceRepo := app.Container.XrayInstanceRepository
	inboundRepo := app.Container.InboundRepository
	clientRepo := app.Container.ClientRepository
	realityRepo := app.Container.RealityConfigRepository
	userRepo := app.Container.UserRepository

	generator := xrayConfig.NewGenerator(instanceRepo, inboundRepo, clientRepo, realityRepo)

	t.Run("GenerateInbound_VLESS_WithReality", func(t *testing.T) {
		defer app.CleanupTables()

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

		reality, err := testutil.CreateTestRealityConfigWithDefaults(inbound.ID)
		require.NoError(t, err)
		err = realityRepo.Create(ctx, reality)
		require.NoError(t, err)

		// Generate inbound
		inboundConfig, err := generator.GenerateInbound(ctx, inbound, []*xray.Client{client}, reality)
		require.NoError(t, err)
		assert.NotNil(t, inboundConfig)
		assert.Equal(t, inbound.Port, inboundConfig.Port)
		assert.Equal(t, string(inbound.Protocol), inboundConfig.Protocol)
		assert.NotNil(t, inboundConfig.Settings)
		assert.NotNil(t, inboundConfig.StreamSettings)
		assert.NotNil(t, inboundConfig.Sniffing)
	})

	t.Run("GenerateInbound_MultipleClients", func(t *testing.T) {
		defer app.CleanupTables()

		user1, err := testutil.CreateTestUserWithDefaults()
		require.NoError(t, err)
		err = userRepo.Create(ctx, user1)
		require.NoError(t, err)

		user2, err := testutil.CreateTestUser("user2", "user2@test.com", "Pass123!@#")
		require.NoError(t, err)
		err = userRepo.Create(ctx, user2)
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

		client1, err := testutil.CreateTestClientWithDefaults(inbound.ID, user1.ID)
		require.NoError(t, err)

		client2, err := testutil.CreateTestClientWithDefaults(inbound.ID, user2.ID)
		require.NoError(t, err)

		inboundConfig, err := generator.GenerateInbound(ctx, inbound, []*xray.Client{client1, client2}, nil)
		require.NoError(t, err)
		assert.NotNil(t, inboundConfig)

		// Verify settings contain both clients
		settings := inboundConfig.Settings
		require.NotNil(t, settings)
		assert.Contains(t, settings, "clients")
	})
}

func TestConfigGenerator_ConfigStructure(t *testing.T) {
	if !testutil.IsIntegrationTest() {
		t.Skip("Skipping integration test")
	}

	app := testutil.NewTestApp(t)
	defer app.Cleanup()
	defer app.CleanupTables()

	ctx := context.Background()
	instanceRepo := app.Container.XrayInstanceRepository
	inboundRepo := app.Container.InboundRepository
	clientRepo := app.Container.ClientRepository
	realityRepo := app.Container.RealityConfigRepository

	generator := xrayConfig.NewGenerator(instanceRepo, inboundRepo, clientRepo, realityRepo)

	t.Run("Config_HasRequiredFields", func(t *testing.T) {
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

		config, err := generator.Generate(ctx, instance.ID)
		require.NoError(t, err)

		// Log config
		assert.NotNil(t, config.Log)
		assert.NotEmpty(t, config.Log.Loglevel)

		// API config
		assert.NotNil(t, config.API)
		assert.NotEmpty(t, config.API.Tag)
		assert.NotEmpty(t, config.API.Services)

		// DNS config
		assert.NotNil(t, config.DNS)
		assert.NotEmpty(t, config.DNS.Servers)

		// Stats config
		assert.NotNil(t, config.Stats)

		// Policy config
		assert.NotNil(t, config.Policy)
		assert.NotEmpty(t, config.Policy.Levels)

		// Outbounds
		assert.NotEmpty(t, config.Outbounds)
		assert.GreaterOrEqual(t, len(config.Outbounds), 2) // direct + block

		// Routing
		assert.NotNil(t, config.Routing)
		assert.NotEmpty(t, config.Routing.Rules)
	})

	t.Run("Config_OutboundsContainDirectAndBlock", func(t *testing.T) {
		defer app.CleanupTables()

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

		config, err := generator.Generate(ctx, instance.ID)
		require.NoError(t, err)

		// Find direct and block outbounds
		hasDirect := false
		hasBlock := false

		for _, outbound := range config.Outbounds {
			if outbound.Tag == "direct" {
				hasDirect = true
				assert.Equal(t, "freedom", outbound.Protocol)
			}
			if outbound.Tag == "block" {
				hasBlock = true
				assert.Equal(t, "blackhole", outbound.Protocol)
			}
		}

		assert.True(t, hasDirect, "Config should have direct outbound")
		assert.True(t, hasBlock, "Config should have block outbound")
	})

	t.Run("Config_APIInboundPresent", func(t *testing.T) {
		defer app.CleanupTables()

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

		config, err := generator.Generate(ctx, instance.ID)
		require.NoError(t, err)

		// Find API inbound
		hasAPIInbound := false
		for _, inbound := range config.Inbounds {
			if inbound.Tag == "api" {
				hasAPIInbound = true
				assert.Equal(t, "dokodemo-door", inbound.Protocol)
				break
			}
		}

		assert.True(t, hasAPIInbound, "Config should have API inbound")
	})
}

func TestConfigGenerator_JSONSerialization(t *testing.T) {
	if !testutil.IsIntegrationTest() {
		t.Skip("Skipping integration test")
	}

	app := testutil.NewTestApp(t)
	defer app.Cleanup()
	defer app.CleanupTables()

	ctx := context.Background()
	instanceRepo := app.Container.XrayInstanceRepository
	inboundRepo := app.Container.InboundRepository
	clientRepo := app.Container.ClientRepository
	realityRepo := app.Container.RealityConfigRepository

	generator := xrayConfig.NewGenerator(instanceRepo, inboundRepo, clientRepo, realityRepo)

	t.Run("JSON_ValidFormat", func(t *testing.T) {
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

		jsonData, err := generator.GenerateJSON(ctx, instance.ID)
		require.NoError(t, err)

		// Verify JSON is well-formed
		var result map[string]interface{}
		err = json.Unmarshal(jsonData, &result)
		require.NoError(t, err)

		// Verify indentation (pretty print)
		assert.Contains(t, string(jsonData), "\n")
		assert.Contains(t, string(jsonData), "  ")
	})
}
