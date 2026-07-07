package integration_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/suproxy/backend/internal/domain/audit"
	"github.com/suproxy/backend/internal/infrastructure/testutil"
)

func TestAuditService_CreateLog(t *testing.T) {
	if !testutil.IsIntegrationTest() {
		t.Skip("Skipping integration test")
	}

	app := testutil.NewTestApp(t)
	defer app.Cleanup()
	defer app.CleanupTables()

	ctx := context.Background()
	repo := app.Container.AuditLogRepository

	tests := []struct {
		name       string
		userID     uuid.UUID
		action     string
		entityType string
		entityID   uuid.UUID
		ipAddress  string
		userAgent  string
	}{
		{
			"UserCreated",
			uuid.New(),
			"user.created",
			"user",
			uuid.New(),
			"192.168.1.1",
			"Mozilla/5.0",
		},
		{
			"UserUpdated",
			uuid.New(),
			"user.updated",
			"user",
			uuid.New(),
			"192.168.1.2",
			"Chrome/91.0",
		},
		{
			"UserDeleted",
			uuid.New(),
			"user.deleted",
			"user",
			uuid.New(),
			"192.168.1.3",
			"Safari/14.0",
		},
		{
			"XrayInstanceCreated",
			uuid.New(),
			"xray_instance.created",
			"xray_instance",
			uuid.New(),
			"192.168.1.4",
			"Firefox/89.0",
		},
		{
			"InboundCreated",
			uuid.New(),
			"inbound.created",
			"inbound",
			uuid.New(),
			"192.168.1.5",
			"Edge/91.0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			log := audit.NewLog(
				tt.userID,
				tt.action,
				tt.entityType,
				tt.entityID,
				tt.ipAddress,
				tt.userAgent,
			)

			err := repo.Create(ctx, log)
			require.NoError(t, err)

			// Verify creation
			found, err := repo.FindByID(ctx, log.ID)
			require.NoError(t, err)
			assert.Equal(t, tt.userID, found.UserID)
			assert.Equal(t, tt.action, found.Action)
			assert.Equal(t, tt.entityType, found.EntityType)
			assert.Equal(t, tt.entityID, found.EntityID)
			assert.Equal(t, tt.ipAddress, found.IPAddress)
			assert.Equal(t, tt.userAgent, found.UserAgent)
		})
	}
}

func TestAuditService_ActionTypes(t *testing.T) {
	if !testutil.IsIntegrationTest() {
		t.Skip("Skipping integration test")
	}

	app := testutil.NewTestApp(t)
	defer app.Cleanup()
	defer app.CleanupTables()

	ctx := context.Background()
	repo := app.Container.AuditLogRepository

	actions := []string{
		"user.login",
		"user.logout",
		"user.register",
		"user.password_reset",
		"xray_instance.start",
		"xray_instance.stop",
		"xray_instance.restart",
		"config.reload",
		"client.provision",
		"client.deprovision",
	}

	userID := uuid.New()
	entityID := uuid.New()

	for _, action := range actions {
		t.Run(action, func(t *testing.T) {
			log := audit.NewLog(
				userID,
				action,
				"test_entity",
				entityID,
				"127.0.0.1",
				"test-agent",
			)

			err := repo.Create(ctx, log)
			require.NoError(t, err)

			found, err := repo.FindByID(ctx, log.ID)
			require.NoError(t, err)
			assert.Equal(t, action, found.Action)
		})
	}
}

func TestAuditService_EntityTypes(t *testing.T) {
	if !testutil.IsIntegrationTest() {
		t.Skip("Skipping integration test")
	}

	app := testutil.NewTestApp(t)
	defer app.Cleanup()
	defer app.CleanupTables()

	ctx := context.Background()
	repo := app.Container.AuditLogRepository

	entityTypes := []string{
		"user",
		"xray_instance",
		"inbound",
		"client",
		"subscription",
		"payment",
		"device",
		"session",
	}

	userID := uuid.New()

	for _, entityType := range entityTypes {
		t.Run(entityType, func(t *testing.T) {
			log := audit.NewLog(
				userID,
				"entity.action",
				entityType,
				uuid.New(),
				"127.0.0.1",
				"test-agent",
			)

			err := repo.Create(ctx, log)
			require.NoError(t, err)

			found, err := repo.FindByID(ctx, log.ID)
			require.NoError(t, err)
			assert.Equal(t, entityType, found.EntityType)
		})
	}
}

func TestAuditService_IPAddressTracking(t *testing.T) {
	if !testutil.IsIntegrationTest() {
		t.Skip("Skipping integration test")
	}

	app := testutil.NewTestApp(t)
	defer app.Cleanup()
	defer app.CleanupTables()

	ctx := context.Background()
	repo := app.Container.AuditLogRepository

	ipAddresses := []string{
		"192.168.1.1",
		"10.0.0.1",
		"172.16.0.1",
		"2001:0db8:85a3:0000:0000:8a2e:0370:7334", // IPv6
		"127.0.0.1",                                 // Localhost
		"8.8.8.8",                                   // Public IP
	}

	userID := uuid.New()

	for _, ip := range ipAddresses {
		t.Run(ip, func(t *testing.T) {
			log := audit.NewLog(
				userID,
				"user.action",
				"user",
				uuid.New(),
				ip,
				"test-agent",
			)

			err := repo.Create(ctx, log)
			require.NoError(t, err)

			found, err := repo.FindByID(ctx, log.ID)
			require.NoError(t, err)
			assert.Equal(t, ip, found.IPAddress)
		})
	}
}

func TestAuditService_UserAgentTracking(t *testing.T) {
	if !testutil.IsIntegrationTest() {
		t.Skip("Skipping integration test")
	}

	app := testutil.NewTestApp(t)
	defer app.Cleanup()
	defer app.CleanupTables()

	ctx := context.Background()
	repo := app.Container.AuditLogRepository

	userAgents := []string{
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36",
		"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36",
		"curl/7.68.0",
		"PostmanRuntime/7.26.8",
		"Go-http-client/1.1",
	}

	userID := uuid.New()

	for _, ua := range userAgents {
		t.Run(ua[:20], func(t *testing.T) { // Truncate name for readability
			log := audit.NewLog(
				userID,
				"user.action",
				"user",
				uuid.New(),
				"127.0.0.1",
				ua,
			)

			err := repo.Create(ctx, log)
			require.NoError(t, err)

			found, err := repo.FindByID(ctx, log.ID)
			require.NoError(t, err)
			assert.Equal(t, ua, found.UserAgent)
		})
	}
}

func TestAuditService_BulkLogging(t *testing.T) {
	if !testutil.IsIntegrationTest() {
		t.Skip("Skipping integration test")
	}

	app := testutil.NewTestApp(t)
	defer app.Cleanup()
	defer app.CleanupTables()

	ctx := context.Background()
	repo := app.Container.AuditLogRepository

	t.Run("Create100Logs", func(t *testing.T) {
		userID := uuid.New()

		for i := 0; i < 100; i++ {
			log := audit.NewLog(
				userID,
				"user.action",
				"user",
				uuid.New(),
				"127.0.0.1",
				"test-agent",
			)

			err := repo.Create(ctx, log)
			require.NoError(t, err)
		}

		// Verify count
		logs, err := repo.FindByUserID(ctx, userID, testutil.TimePast(1), testutil.TimeFuture(1))
		require.NoError(t, err)
		assert.GreaterOrEqual(t, len(logs), 100)
	})
}

func TestAuditService_Timestamps(t *testing.T) {
	if !testutil.IsIntegrationTest() {
		t.Skip("Skipping integration test")
	}

	app := testutil.NewTestApp(t)
	defer app.Cleanup()
	defer app.CleanupTables()

	ctx := context.Background()
	repo := app.Container.AuditLogRepository

	t.Run("CreatedAt_AutoSet", func(t *testing.T) {
		log := audit.NewLog(
			uuid.New(),
			"user.action",
			"user",
			uuid.New(),
			"127.0.0.1",
			"test-agent",
		)

		err := repo.Create(ctx, log)
		require.NoError(t, err)

		found, err := repo.FindByID(ctx, log.ID)
		require.NoError(t, err)

		// CreatedAt should be recent
		testutil.AssertTimeNow(t, found.CreatedAt, 5000, "CreatedAt should be recent")
	})
}

func TestAuditService_QueryByEntity(t *testing.T) {
	if !testutil.IsIntegrationTest() {
		t.Skip("Skipping integration test")
	}

	app := testutil.NewTestApp(t)
	defer app.Cleanup()
	defer app.CleanupTables()

	ctx := context.Background()
	repo := app.Container.AuditLogRepository

	t.Run("TrackEntityLifecycle", func(t *testing.T) {
		userID := uuid.New()
		entityID := uuid.New()
		entityType := "xray_instance"

		// Create logs for entity lifecycle
		actions := []string{"created", "started", "stopped", "deleted"}

		for _, action := range actions {
			log := audit.NewLog(
				userID,
				"xray_instance."+action,
				entityType,
				entityID,
				"127.0.0.1",
				"test-agent",
			)

			err := repo.Create(ctx, log)
			require.NoError(t, err)
		}

		// Query all logs for this entity
		logs, err := repo.FindByEntityID(ctx, entityType, entityID)
		require.NoError(t, err)
		assert.GreaterOrEqual(t, len(logs), 4)

		// Verify all logs are for the same entity
		for _, log := range logs {
			assert.Equal(t, entityID, log.EntityID)
			assert.Equal(t, entityType, log.EntityType)
		}
	})
}

func TestAuditService_SecurityAudit(t *testing.T) {
	if !testutil.IsIntegrationTest() {
		t.Skip("Skipping integration test")
	}

	app := testutil.NewTestApp(t)
	defer app.Cleanup()
	defer app.CleanupTables()

	ctx := context.Background()
	repo := app.Container.AuditLogRepository

	t.Run("TrackSecurityEvents", func(t *testing.T) {
		userID := uuid.New()

		securityEvents := []string{
			"user.login.success",
			"user.login.failed",
			"user.password_changed",
			"user.account_locked",
			"user.account_unlocked",
			"user.2fa_enabled",
			"user.2fa_disabled",
		}

		for _, event := range securityEvents {
			log := audit.NewLog(
				userID,
				event,
				"user",
				userID,
				"192.168.1.100",
				"test-agent",
			)

			err := repo.Create(ctx, log)
			require.NoError(t, err)
		}

		// Verify all security events are logged
		logs, err := repo.FindByUserID(ctx, userID, testutil.TimePast(1), testutil.TimeFuture(1))
		require.NoError(t, err)
		assert.GreaterOrEqual(t, len(logs), len(securityEvents))
	})
}

