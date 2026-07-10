package integration

import (
	"context"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
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

	// Create use-case instances
	listUsersQuery := adminuser.NewListUsersQuery(app.Container.UserRepository)
	getUserQuery := adminuser.NewGetUserQuery(app.Container.UserRepository)
	updateUserStatusCommand := adminuser.NewUpdateUserStatusCommand(app.Container.UserRepository, app.Container.AuditLogRepository)
	updateUserRoleCommand := adminuser.NewUpdateUserRoleCommand(app.Container.UserRepository, app.Container.AuditLogRepository)
	listInstancesQuery := adminxray.NewListInstancesQuery(app.Container.XrayInstanceRepository)
	getInstanceQuery := adminxray.NewGetInstanceQuery(app.Container.XrayInstanceRepository)
	getInstanceStatsQuery := adminxray.NewGetInstanceStatsQuery(app.Container.XrayInstanceRepository, app.Container.InboundRepository, app.Container.ClientRepository)
	startInstanceCommand := adminxray.NewStartInstanceCommand(app.Container.XrayInstanceRepository, app.Container.XrayProcessManager, app.Container.AuditLogRepository)
	stopInstanceCommand := adminxray.NewStopInstanceCommand(app.Container.XrayInstanceRepository, app.Container.XrayProcessManager, app.Container.AuditLogRepository)
	restartInstanceCommand := adminxray.NewRestartInstanceCommand(app.Container.XrayInstanceRepository, app.Container.XrayProcessManager, app.Container.AuditLogRepository)
	reloadInstanceCommand := adminxray.NewReloadInstanceCommand(app.Container.XrayProvisioningService, app.Container.AuditLogRepository)
	checkInstanceHealthCommand := adminxray.NewCheckInstanceHealthCommand(app.Container.XrayProcessManager)
	listInboundsQuery := admininbound.NewListInboundsQuery(app.Container.InboundRepository)
	getInboundQuery := admininbound.NewGetInboundQuery(app.Container.InboundRepository)
	createInboundCommand := admininbound.NewCreateInboundCommand(app.Container.InboundRepository, app.Container.XrayInstanceRepository, app.Container.XrayProvisioningService, app.Container.AuditLogRepository)
	updateInboundCommand := admininbound.NewUpdateInboundCommand(app.Container.InboundRepository, app.Container.XrayProvisioningService, app.Container.AuditLogRepository)
	deleteInboundCommand := admininbound.NewDeleteInboundCommand(app.Container.InboundRepository, app.Container.ClientRepository, app.Container.XrayProvisioningService, app.Container.AuditLogRepository)
	enableInboundCommand := admininbound.NewEnableInboundCommand(app.Container.InboundRepository, app.Container.XrayProvisioningService, app.Container.AuditLogRepository)
	disableInboundCommand := admininbound.NewDisableInboundCommand(app.Container.InboundRepository, app.Container.XrayProvisioningService, app.Container.AuditLogRepository)
	listClientsQuery := adminclient.NewListClientsQuery(app.Container.ClientRepository)
	getClientQuery := adminclient.NewGetClientQuery(app.Container.ClientRepository)
	createClientCommand := adminclient.NewCreateClientCommand(app.Container.ClientRepository, app.Container.InboundRepository, app.Container.XrayProvisioningService, app.Container.AuditLogRepository)
	deleteClientCommand := adminclient.NewDeleteClientCommand(app.Container.ClientRepository, app.Container.InboundRepository, app.Container.XrayProvisioningService, app.Container.AuditLogRepository)
	enableClientCommand := adminclient.NewEnableClientCommand(app.Container.ClientRepository, app.Container.InboundRepository, app.Container.XrayProvisioningService, app.Container.AuditLogRepository)
	disableClientCommand := adminclient.NewDisableClientCommand(app.Container.ClientRepository, app.Container.InboundRepository, app.Container.XrayProvisioningService, app.Container.AuditLogRepository)
	regenerateClientUUIDCommand := adminclient.NewRegenerateClientUUIDCommand(app.Container.ClientRepository, app.Container.InboundRepository, app.Container.XrayProvisioningService, app.Container.AuditLogRepository)
	reprovisionClientCommand := adminclient.NewReprovisionClientCommand(app.Container.ClientRepository, app.Container.InboundRepository, app.Container.XrayProvisioningService, app.Container.AuditLogRepository)
	listAuditLogsQuery := adminaudit.NewListAuditLogsQuery(app.Container.AuditLogRepository)
	getAuditLogQuery := adminaudit.NewGetAuditLogQuery(app.Container.AuditLogRepository)
	getAuditStatsQuery := adminaudit.NewGetAuditStatsQuery(app.Container.AuditLogRepository)
	getSystemHealthQuery := adminsystem.NewGetSystemHealthQuery(app.Database, app.Container.XrayInstanceRepository, app.Container.XrayProcessManager)
	getSystemStatsQuery := adminsystem.NewGetSystemStatsQuery(app.Container.UserRepository, app.Container.XrayInstanceRepository, app.Container.InboundRepository, app.Container.ClientRepository, app.Container.AuditLogRepository)
	getVersionQuery := adminsystem.NewGetVersionQuery()
	getDatabaseStatusQuery := adminsystem.NewGetDatabaseStatusQuery(app.Database)
	getXraySystemStatusQuery := adminsystem.NewGetXraySystemStatusQuery(app.Container.XrayInstanceRepository)

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
		authHelper := testutil.NewAuthHelper(app.JWT, t)
		adminUser, adminToken, _ := authHelper.CreateAuthenticatedAdmin(app.Container.UserRepository)

		// Create audit logs
		log := audit.NewLog(adminUser.ID, audit.ActionLogin, "user", adminUser.ID, "127.0.0.1", "TestAgent")
		app.Container.AuditLogRepository.Create(ctx, log)

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

		authHelper := testutil.NewAuthHelper(app.JWT, t)
		_, userToken, _ := authHelper.CreateAuthenticatedUser(app.Container.UserRepository)

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
		authHelper := testutil.NewAuthHelper(app.JWT, t)
		adminUser, adminToken, _ := authHelper.CreateAuthenticatedAdmin(app.Container.UserRepository)

		log := audit.NewLog(adminUser.ID, audit.ActionLogin, "user", adminUser.ID, "127.0.0.1", "TestAgent")
		app.Container.AuditLogRepository.Create(ctx, log)

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

		authHelper := testutil.NewAuthHelper(app.JWT, t)
		_, adminToken, _ := authHelper.CreateAuthenticatedAdmin(app.Container.UserRepository)

		httpCtx := testutil.NewHTTPTestContext(t)
		httpCtx.Router = router

		nonExistentID := uuid.New()
		resp := httpCtx.GET("/api/v1/admin/audit/"+nonExistentID.String(), testutil.AuthHeader(adminToken))

		assert.Equal(t, 404, resp.Code)
	})

	t.Run("BadRequest_InvalidUUID", func(t *testing.T) {
		defer app.CleanupTables()

		authHelper := testutil.NewAuthHelper(app.JWT, t)
		_, adminToken, _ := authHelper.CreateAuthenticatedAdmin(app.Container.UserRepository)

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
		authHelper := testutil.NewAuthHelper(app.JWT, t)
		adminUser, adminToken, _ := authHelper.CreateAuthenticatedAdmin(app.Container.UserRepository)

		// Create some audit logs
		for i := 0; i < 5; i++ {
			log := audit.NewLog(adminUser.ID, audit.ActionLogin, "user", adminUser.ID, "127.0.0.1", "TestAgent")
			app.Container.AuditLogRepository.Create(ctx, log)
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

		authHelper := testutil.NewAuthHelper(app.JWT, t)
		_, adminToken, _ := authHelper.CreateAuthenticatedAdmin(app.Container.UserRepository)

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

		authHelper := testutil.NewAuthHelper(app.JWT, t)
		_, userToken, _ := authHelper.CreateAuthenticatedUser(app.Container.UserRepository)

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

		authHelper := testutil.NewAuthHelper(app.JWT, t)
		_, adminToken, _ := authHelper.CreateAuthenticatedAdmin(app.Container.UserRepository)

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

		authHelper := testutil.NewAuthHelper(app.JWT, t)
		_, adminToken, _ := authHelper.CreateAuthenticatedAdmin(app.Container.UserRepository)

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

		authHelper := testutil.NewAuthHelper(app.JWT, t)
		_, adminToken, _ := authHelper.CreateAuthenticatedAdmin(app.Container.UserRepository)

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

		authHelper := testutil.NewAuthHelper(app.JWT, t)
		_, adminToken, _ := authHelper.CreateAuthenticatedAdmin(app.Container.UserRepository)

		httpCtx := testutil.NewHTTPTestContext(t)
		httpCtx.Router = router

		resp := httpCtx.GET("/api/v1/admin/system/database", testutil.AuthHeader(adminToken))

		assert.Equal(t, 200, resp.Code)

		var result response.Response
		httpCtx.GetResponseJSON(&result)
		assert.True(t, result.Success)
	})
}
