package integration

import (
	"context"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	adminaudit "github.com/suproxy/backend/internal/application/usecase/admin/audit"
	adminclient "github.com/suproxy/backend/internal/application/usecase/admin/client"
	admininbound "github.com/suproxy/backend/internal/application/usecase/admin/inbound"
	adminsystem "github.com/suproxy/backend/internal/application/usecase/admin/system"
	adminuser "github.com/suproxy/backend/internal/application/usecase/admin/user"
	adminxray "github.com/suproxy/backend/internal/application/usecase/admin/xray_instance"
	"github.com/suproxy/backend/internal/domain/audit"
	"github.com/suproxy/backend/internal/infrastructure/testutil"
	"github.com/suproxy/backend/internal/interfaces/http/handler"
	"github.com/suproxy/backend/internal/interfaces/http/middleware"
	"github.com/suproxy/backend/internal/interfaces/http/response"
)

func setupAdminAuditSystemHandler(t *testing.T, app *testutil.TestApp) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	// Create use-case instances from Container
	container := app.Container
	listUsersQuery := adminuser.NewListUsersQuery(container.UserRepository)
	getUserQuery := adminuser.NewGetUserQuery(container.UserRepository)
	updateUserStatusCommand := adminuser.NewUpdateUserStatusCommand(container.UserRepository, container.AuditLogRepository)
	updateUserRoleCommand := adminuser.NewUpdateUserRoleCommand(container.UserRepository, container.AuditLogRepository)
	listInstancesQuery := adminxray.NewListInstancesQuery(container.XrayInstanceRepository)
	getInstanceQuery := adminxray.NewGetInstanceQuery(container.XrayInstanceRepository)
	getInstanceStatsQuery := adminxray.NewGetInstanceStatsQuery(container.XrayInstanceRepository, container.InboundRepository, container.ClientRepository)
	startInstanceCommand := adminxray.NewStartInstanceCommand(container.XrayInstanceRepository, container.XrayProcessManager, container.AuditLogRepository)
	stopInstanceCommand := adminxray.NewStopInstanceCommand(container.XrayInstanceRepository, container.XrayProcessManager, container.AuditLogRepository)
	restartInstanceCommand := adminxray.NewRestartInstanceCommand(container.XrayInstanceRepository, container.XrayProcessManager, container.AuditLogRepository)
	reloadInstanceCommand := adminxray.NewReloadInstanceCommand(container.XrayProvisioningService, container.AuditLogRepository)
	checkInstanceHealthCommand := adminxray.NewCheckInstanceHealthCommand(container.XrayProcessManager)
	listInboundsQuery := admininbound.NewListInboundsQuery(container.InboundRepository)
	getInboundQuery := admininbound.NewGetInboundQuery(container.InboundRepository)
	createInboundCommand := admininbound.NewCreateInboundCommand(container.InboundRepository, container.XrayInstanceRepository, container.XrayProvisioningService, container.AuditLogRepository)
	updateInboundCommand := admininbound.NewUpdateInboundCommand(container.InboundRepository, container.XrayProvisioningService, container.AuditLogRepository)
	deleteInboundCommand := admininbound.NewDeleteInboundCommand(container.InboundRepository, container.ClientRepository, container.XrayProvisioningService, container.AuditLogRepository)
	enableInboundCommand := admininbound.NewEnableInboundCommand(container.InboundRepository, container.XrayProvisioningService, container.AuditLogRepository)
	disableInboundCommand := admininbound.NewDisableInboundCommand(container.InboundRepository, container.XrayProvisioningService, container.AuditLogRepository)
	listClientsQuery := adminclient.NewListClientsQuery(container.ClientRepository)
	getClientQuery := adminclient.NewGetClientQuery(container.ClientRepository)
	createClientCommand := adminclient.NewCreateClientCommand(container.ClientRepository, container.InboundRepository, container.XrayProvisioningService, container.AuditLogRepository)
	deleteClientCommand := adminclient.NewDeleteClientCommand(container.ClientRepository, container.InboundRepository, container.XrayProvisioningService, container.AuditLogRepository)
	enableClientCommand := adminclient.NewEnableClientCommand(container.ClientRepository, container.InboundRepository, container.XrayProvisioningService, container.AuditLogRepository)
	disableClientCommand := adminclient.NewDisableClientCommand(container.ClientRepository, container.InboundRepository, container.XrayProvisioningService, container.AuditLogRepository)
	regenerateClientUUIDCommand := adminclient.NewRegenerateClientUUIDCommand(container.ClientRepository, container.InboundRepository, container.XrayProvisioningService, container.AuditLogRepository)
	reprovisionClientCommand := adminclient.NewReprovisionClientCommand(container.ClientRepository, container.InboundRepository, container.XrayProvisioningService, container.AuditLogRepository)
	listAuditLogsQuery := adminaudit.NewListAuditLogsQuery(container.AuditLogRepository)
	getAuditLogQuery := adminaudit.NewGetAuditLogQuery(container.AuditLogRepository)
	getAuditStatsQuery := adminaudit.NewGetAuditStatsQuery(container.AuditLogRepository)
	getSystemHealthQuery := adminsystem.NewGetSystemHealthQuery(app.Database, container.XrayInstanceRepository, container.XrayProcessManager)
	getSystemStatsQuery := adminsystem.NewGetSystemStatsQuery(container.UserRepository, container.XrayInstanceRepository, container.InboundRepository, container.ClientRepository, container.AuditLogRepository)
	getVersionQuery := adminsystem.NewGetVersionQuery()
	getDatabaseStatusQuery := adminsystem.NewGetDatabaseStatusQuery(app.Database)
	getXraySystemStatusQuery := adminsystem.NewGetXraySystemStatusQuery(container.XrayInstanceRepository)

	adminHandler := handler.NewAdminHandler(
		app.Logger,
		listUsersQuery,
		getUserQuery,
		updateUserStatusCommand,
		updateUserRoleCommand,
		listInstancesQuery,
		getInstanceQuery,
		getInstanceStatsQuery,
		startInstanceCommand,
		stopInstanceCommand,
		restartInstanceCommand,
		reloadInstanceCommand,
		checkInstanceHealthCommand,
		listInboundsQuery,
		getInboundQuery,
		createInboundCommand,
		updateInboundCommand,
		deleteInboundCommand,
		enableInboundCommand,
		disableInboundCommand,
		listClientsQuery,
		getClientQuery,
		createClientCommand,
		deleteClientCommand,
		enableClientCommand,
		disableClientCommand,
		regenerateClientUUIDCommand,
		reprovisionClientCommand,
		listAuditLogsQuery,
		getAuditLogQuery,
		getAuditStatsQuery,
		getSystemHealthQuery,
		getSystemStatsQuery,
		getVersionQuery,
		getDatabaseStatusQuery,
		getXraySystemStatusQuery,
	)

	adminGroup := router.Group("/api/v1/admin")
	adminGroup.Use(middleware.AuthMiddleware(app.JWT))
	adminGroup.Use(middleware.RequireAdmin())

	// Audit routes
	adminGroup.GET("/audit", adminHandler.ListAuditLogs)
	adminGroup.GET("/audit/:id", adminHandler.GetAuditLog)
	adminGroup.GET("/audit/stats", adminHandler.GetAuditStats)

	// System routes
	adminGroup.GET("/health", adminHandler.HealthCheck)
	adminGroup.GET("/system/health", adminHandler.GetSystemHealth)
	adminGroup.GET("/system/stats", adminHandler.GetSystemStats)
	adminGroup.GET("/system/version", adminHandler.GetVersion)
	adminGroup.GET("/system/database", adminHandler.GetDatabaseStatus)
	adminGroup.GET("/system/xray", adminHandler.GetXraySystemStatus)

	return router
}

func TestAdminHandler_ListAuditLogs(t *testing.T) {
	if !testutil.IsIntegrationTest() {
		t.Skip("Skipping integration test")
	}

	app := testutil.NewTestApp(t)
	defer app.Cleanup()
	defer app.CleanupTables()

	router := setupAdminAuditSystemHandler(t, app)

	t.Run("Success", func(t *testing.T) {
		defer app.CleanupTables()

		ctx := context.Background()

		// Create admin user
		adminUser, err := testutil.CreateTestAdminUser()
		require.NoError(t, err)
		err = app.Container.UserRepository.Create(ctx, adminUser)
		require.NoError(t, err)

		// Generate admin token
		adminToken, err := app.JWT.GenerateAccessToken(adminUser.ID.String(), adminUser.Email.String(), string(adminUser.Role))
		require.NoError(t, err)

		// Create audit logs
		log := audit.NewLog(adminUser.ID, audit.ActionLogin, "user", adminUser.ID, "127.0.0.1", "TestAgent")
		_ = app.Container.AuditLogRepository.Create(ctx, log) // Test setup

		httpCtx := testutil.NewHTTPTestContext(t)
		httpCtx.Router = router

		resp := httpCtx.GET("/api/v1/admin/audit", testutil.AuthHeader(adminToken))

		assert.Equal(t, 200, resp.Code)

		var result response.Response
		httpCtx.GetResponseJSON(&result)
		assert.True(t, result.Success)
	})

	t.Run("Unauthorized_NoToken", func(t *testing.T) {
		defer app.CleanupTables()

		httpCtx := testutil.NewHTTPTestContext(t)
		httpCtx.Router = router

		resp := httpCtx.GET("/api/v1/admin/audit", nil)
		assert.Equal(t, 401, resp.Code)
	})

	t.Run("Forbidden_NotAdmin", func(t *testing.T) {
		defer app.CleanupTables()

		ctx := context.Background()

		// Create regular user
		regularUser, err := testutil.CreateTestUserWithDefaults()
		assert.NoError(t, err)
		err = app.Container.UserRepository.Create(ctx, regularUser)
		assert.NoError(t, err)

		// Generate user token
		userToken, err := app.JWT.GenerateAccessToken(regularUser.ID.String(), regularUser.Email.String(), string(regularUser.Role))
		assert.NoError(t, err)

		httpCtx := testutil.NewHTTPTestContext(t)
		httpCtx.Router = router

		resp := httpCtx.GET("/api/v1/admin/audit", testutil.AuthHeader(userToken))
		assert.Equal(t, 403, resp.Code)
	})
}

func TestAdminHandler_GetAuditLog(t *testing.T) {
	if !testutil.IsIntegrationTest() {
		t.Skip("Skipping integration test")
	}

	app := testutil.NewTestApp(t)
	defer app.Cleanup()
	defer app.CleanupTables()

	router := setupAdminAuditSystemHandler(t, app)

	t.Run("Success", func(t *testing.T) {
		defer app.CleanupTables()

		ctx := context.Background()

		// Create admin user
		adminUser, err := testutil.CreateTestAdminUser()
		assert.NoError(t, err)
		err = app.Container.UserRepository.Create(ctx, adminUser)
		assert.NoError(t, err)

		// Generate admin token
		adminToken, err := app.JWT.GenerateAccessToken(adminUser.ID.String(), adminUser.Email.String(), string(adminUser.Role))
		assert.NoError(t, err)

		log := audit.NewLog(adminUser.ID, audit.ActionLogin, "user", adminUser.ID, "127.0.0.1", "TestAgent")
		_ = app.Container.AuditLogRepository.Create(ctx, log) // Test setup

		httpCtx := testutil.NewHTTPTestContext(t)
		httpCtx.Router = router

		resp := httpCtx.GET("/api/v1/admin/audit/"+log.ID.String(), testutil.AuthHeader(adminToken))

		assert.Equal(t, 200, resp.Code)

		var result response.Response
		httpCtx.GetResponseJSON(&result)
		assert.True(t, result.Success)
	})

	t.Run("NotFound_LogDoesNotExist", func(t *testing.T) {
		defer app.CleanupTables()

		ctx := context.Background()

		// Create admin user
		adminUser, err := testutil.CreateTestAdminUser()
		assert.NoError(t, err)
		err = app.Container.UserRepository.Create(ctx, adminUser)
		assert.NoError(t, err)

		// Generate admin token
		adminToken, err := app.JWT.GenerateAccessToken(adminUser.ID.String(), adminUser.Email.String(), string(adminUser.Role))
		assert.NoError(t, err)

		httpCtx := testutil.NewHTTPTestContext(t)
		httpCtx.Router = router

		nonExistentID := uuid.New()
		resp := httpCtx.GET("/api/v1/admin/audit/"+nonExistentID.String(), testutil.AuthHeader(adminToken))

		assert.Equal(t, 404, resp.Code)
	})

	t.Run("BadRequest_InvalidUUID", func(t *testing.T) {
		defer app.CleanupTables()

		ctx := context.Background()

		// Create admin user
		adminUser, err := testutil.CreateTestAdminUser()
		assert.NoError(t, err)
		err = app.Container.UserRepository.Create(ctx, adminUser)
		assert.NoError(t, err)

		// Generate admin token
		adminToken, err := app.JWT.GenerateAccessToken(adminUser.ID.String(), adminUser.Email.String(), string(adminUser.Role))
		assert.NoError(t, err)

		httpCtx := testutil.NewHTTPTestContext(t)
		httpCtx.Router = router

		resp := httpCtx.GET("/api/v1/admin/audit/invalid-uuid", testutil.AuthHeader(adminToken))
		assert.Equal(t, 400, resp.Code)
	})
}

func TestAdminHandler_GetAuditStats(t *testing.T) {
	if !testutil.IsIntegrationTest() {
		t.Skip("Skipping integration test")
	}

	app := testutil.NewTestApp(t)
	defer app.Cleanup()
	defer app.CleanupTables()

	router := setupAdminAuditSystemHandler(t, app)

	t.Run("Success", func(t *testing.T) {
		defer app.CleanupTables()

		ctx := context.Background()

		// Create admin user
		adminUser, err := testutil.CreateTestAdminUser()
		assert.NoError(t, err)
		err = app.Container.UserRepository.Create(ctx, adminUser)
		assert.NoError(t, err)

		// Generate admin token
		adminToken, err := app.JWT.GenerateAccessToken(adminUser.ID.String(), adminUser.Email.String(), string(adminUser.Role))
		assert.NoError(t, err)

		// Create some audit logs
		for i := 0; i < 5; i++ {
			log := audit.NewLog(adminUser.ID, audit.ActionLogin, "user", adminUser.ID, "127.0.0.1", "TestAgent")
			_ = app.Container.AuditLogRepository.Create(ctx, log) // Test setup
		}

		httpCtx := testutil.NewHTTPTestContext(t)
		httpCtx.Router = router

		resp := httpCtx.GET("/api/v1/admin/audit/stats", testutil.AuthHeader(adminToken))

		assert.Equal(t, 200, resp.Code)

		var result response.Response
		httpCtx.GetResponseJSON(&result)
		assert.True(t, result.Success)
	})
}

func TestAdminHandler_HealthCheck(t *testing.T) {
	if !testutil.IsIntegrationTest() {
		t.Skip("Skipping integration test")
	}

	app := testutil.NewTestApp(t)
	defer app.Cleanup()
	defer app.CleanupTables()

	router := setupAdminAuditSystemHandler(t, app)

	t.Run("Success", func(t *testing.T) {
		defer app.CleanupTables()

		ctx := context.Background()

		// Create admin user
		adminUser, err := testutil.CreateTestAdminUser()
		assert.NoError(t, err)
		err = app.Container.UserRepository.Create(ctx, adminUser)
		assert.NoError(t, err)

		// Generate admin token
		adminToken, err := app.JWT.GenerateAccessToken(adminUser.ID.String(), adminUser.Email.String(), string(adminUser.Role))
		assert.NoError(t, err)

		httpCtx := testutil.NewHTTPTestContext(t)
		httpCtx.Router = router

		resp := httpCtx.GET("/api/v1/admin/health", testutil.AuthHeader(adminToken))

		assert.Equal(t, 200, resp.Code)

		var result response.Response
		httpCtx.GetResponseJSON(&result)
		assert.True(t, result.Success)
	})

	t.Run("Unauthorized_NoToken", func(t *testing.T) {
		defer app.CleanupTables()

		httpCtx := testutil.NewHTTPTestContext(t)
		httpCtx.Router = router

		resp := httpCtx.GET("/api/v1/admin/health", nil)
		assert.Equal(t, 401, resp.Code)
	})

	t.Run("Forbidden_NotAdmin", func(t *testing.T) {
		defer app.CleanupTables()

		ctx := context.Background()

		// Create regular user
		regularUser, err := testutil.CreateTestUserWithDefaults()
		assert.NoError(t, err)
		err = app.Container.UserRepository.Create(ctx, regularUser)
		assert.NoError(t, err)

		// Generate user token
		userToken, err := app.JWT.GenerateAccessToken(regularUser.ID.String(), regularUser.Email.String(), string(regularUser.Role))
		assert.NoError(t, err)

		httpCtx := testutil.NewHTTPTestContext(t)
		httpCtx.Router = router

		resp := httpCtx.GET("/api/v1/admin/health", testutil.AuthHeader(userToken))
		assert.Equal(t, 403, resp.Code)
	})
}

func TestAdminHandler_GetSystemHealth(t *testing.T) {
	if !testutil.IsIntegrationTest() {
		t.Skip("Skipping integration test")
	}

	app := testutil.NewTestApp(t)
	defer app.Cleanup()
	defer app.CleanupTables()

	router := setupAdminAuditSystemHandler(t, app)

	t.Run("Success", func(t *testing.T) {
		defer app.CleanupTables()

		ctx := context.Background()

		// Create admin user
		adminUser, err := testutil.CreateTestAdminUser()
		assert.NoError(t, err)
		err = app.Container.UserRepository.Create(ctx, adminUser)
		assert.NoError(t, err)

		// Generate admin token
		adminToken, err := app.JWT.GenerateAccessToken(adminUser.ID.String(), adminUser.Email.String(), string(adminUser.Role))
		assert.NoError(t, err)

		httpCtx := testutil.NewHTTPTestContext(t)
		httpCtx.Router = router

		resp := httpCtx.GET("/api/v1/admin/system/health", testutil.AuthHeader(adminToken))

		assert.Equal(t, 200, resp.Code)

		var result response.Response
		httpCtx.GetResponseJSON(&result)
		assert.True(t, result.Success)
	})
}

func TestAdminHandler_GetSystemStats(t *testing.T) {
	if !testutil.IsIntegrationTest() {
		t.Skip("Skipping integration test")
	}

	app := testutil.NewTestApp(t)
	defer app.Cleanup()
	defer app.CleanupTables()

	router := setupAdminAuditSystemHandler(t, app)

	t.Run("Success", func(t *testing.T) {
		defer app.CleanupTables()

		ctx := context.Background()

		// Create admin user
		adminUser, err := testutil.CreateTestAdminUser()
		assert.NoError(t, err)
		err = app.Container.UserRepository.Create(ctx, adminUser)
		assert.NoError(t, err)

		// Generate admin token
		adminToken, err := app.JWT.GenerateAccessToken(adminUser.ID.String(), adminUser.Email.String(), string(adminUser.Role))
		assert.NoError(t, err)

		httpCtx := testutil.NewHTTPTestContext(t)
		httpCtx.Router = router

		resp := httpCtx.GET("/api/v1/admin/system/stats", testutil.AuthHeader(adminToken))

		assert.Equal(t, 200, resp.Code)

		var result response.Response
		httpCtx.GetResponseJSON(&result)
		assert.True(t, result.Success)
	})
}

func TestAdminHandler_GetVersion(t *testing.T) {
	if !testutil.IsIntegrationTest() {
		t.Skip("Skipping integration test")
	}

	app := testutil.NewTestApp(t)
	defer app.Cleanup()
	defer app.CleanupTables()

	router := setupAdminAuditSystemHandler(t, app)

	t.Run("Success", func(t *testing.T) {
		defer app.CleanupTables()

		ctx := context.Background()

		// Create admin user
		adminUser, err := testutil.CreateTestAdminUser()
		assert.NoError(t, err)
		err = app.Container.UserRepository.Create(ctx, adminUser)
		assert.NoError(t, err)

		// Generate admin token
		adminToken, err := app.JWT.GenerateAccessToken(adminUser.ID.String(), adminUser.Email.String(), string(adminUser.Role))
		assert.NoError(t, err)

		httpCtx := testutil.NewHTTPTestContext(t)
		httpCtx.Router = router

		resp := httpCtx.GET("/api/v1/admin/system/version", testutil.AuthHeader(adminToken))

		assert.Equal(t, 200, resp.Code)

		var result response.Response
		httpCtx.GetResponseJSON(&result)
		assert.True(t, result.Success)
	})
}

func TestAdminHandler_GetDatabaseStatus(t *testing.T) {
	if !testutil.IsIntegrationTest() {
		t.Skip("Skipping integration test")
	}

	app := testutil.NewTestApp(t)
	defer app.Cleanup()
	defer app.CleanupTables()

	router := setupAdminAuditSystemHandler(t, app)

	t.Run("Success", func(t *testing.T) {
		defer app.CleanupTables()

		ctx := context.Background()

		// Create admin user
		adminUser, err := testutil.CreateTestAdminUser()
		assert.NoError(t, err)
		err = app.Container.UserRepository.Create(ctx, adminUser)
		assert.NoError(t, err)

		// Generate admin token
		adminToken, err := app.JWT.GenerateAccessToken(adminUser.ID.String(), adminUser.Email.String(), string(adminUser.Role))
		assert.NoError(t, err)

		httpCtx := testutil.NewHTTPTestContext(t)
		httpCtx.Router = router

		resp := httpCtx.GET("/api/v1/admin/system/database", testutil.AuthHeader(adminToken))

		assert.Equal(t, 200, resp.Code)

		var result response.Response
		httpCtx.GetResponseJSON(&result)
		assert.True(t, result.Success)
	})
}
