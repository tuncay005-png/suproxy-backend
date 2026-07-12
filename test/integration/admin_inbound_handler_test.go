package integration

import (
	"context"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/suproxy/backend/internal/application/dto"
	adminaudit "github.com/suproxy/backend/internal/application/usecase/admin/audit"
	adminclient "github.com/suproxy/backend/internal/application/usecase/admin/client"
	admininbound "github.com/suproxy/backend/internal/application/usecase/admin/inbound"
	adminsystem "github.com/suproxy/backend/internal/application/usecase/admin/system"
	adminuser "github.com/suproxy/backend/internal/application/usecase/admin/user"
	adminxray "github.com/suproxy/backend/internal/application/usecase/admin/xray_instance"
	"github.com/suproxy/backend/internal/domain/xray"
	"github.com/suproxy/backend/internal/infrastructure/testutil"
	"github.com/suproxy/backend/internal/interfaces/http/handler"
	"github.com/suproxy/backend/internal/interfaces/http/middleware"
	"github.com/suproxy/backend/internal/interfaces/http/response"
)

func setupAdminInboundHandler(t *testing.T, app *testutil.TestApp) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	// Create use-case instances for admin handler
	// Note: These are created directly as Container doesn't hold them
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

	adminGroup := router.Group("/api/v1/admin/xray")
	adminGroup.Use(middleware.AuthMiddleware(app.JWT))
	adminGroup.Use(middleware.RequireAdmin())

	adminGroup.GET("/inbounds", adminHandler.ListInbounds)
	adminGroup.GET("/inbounds/:id", adminHandler.GetInbound)
	adminGroup.POST("/inbounds", adminHandler.CreateInbound)
	adminGroup.PUT("/inbounds/:id", adminHandler.UpdateInbound)
	adminGroup.DELETE("/inbounds/:id", adminHandler.DeleteInbound)
	adminGroup.POST("/inbounds/:id/enable", adminHandler.EnableInbound)
	adminGroup.POST("/inbounds/:id/disable", adminHandler.DisableInbound)

	return router
}

func TestAdminHandler_ListInbounds(t *testing.T) {
	if !testutil.IsIntegrationTest() {
		t.Skip("Skipping integration test")
	}

	app := testutil.NewTestApp(t)
	defer app.Cleanup()
	defer app.CleanupTables()

	router := setupAdminInboundHandler(t, app)

	t.Run("Success", func(t *testing.T) {
		defer app.CleanupTables()

		ctx := context.Background()
		authHelper := testutil.NewAuthHelper(app.JWT, t)
		_, adminToken, _ := authHelper.CreateAuthenticatedAdmin(app.Container.UserRepository)

		// Create dependencies - server and node first
		testServer, _ := testutil.CreateTestServerWithDefaults()
		_ = app.Container.ServerRepository.Create(ctx, testServer) // Test setup

		testNode, _ := testutil.CreateTestNodeWithDefaults(testServer.ID)
		_ = app.Container.NodeRepository.Create(ctx, testNode) // Test setup

		// Create test instance and inbound
		instance, _ := testutil.CreateTestXrayInstanceWithDefaults(testNode.ID)
		_ = app.Container.XrayInstanceRepository.Create(ctx, instance) // Test setup

		inbound, _ := testutil.CreateTestInboundWithDefaults(instance.ID)
		_ = app.Container.InboundRepository.Create(ctx, inbound) // Test setup

		httpCtx := testutil.NewHTTPTestContext(t)
		httpCtx.Router = router

		resp := httpCtx.GET("/api/v1/admin/xray/inbounds", testutil.AuthHeader(adminToken))

		assert.Equal(t, 200, resp.Code)

		var result response.Response
		httpCtx.GetResponseJSON(&result)
		assert.True(t, result.Success)
	})

	t.Run("Unauthorized_NoToken", func(t *testing.T) {
		defer app.CleanupTables()

		httpCtx := testutil.NewHTTPTestContext(t)
		httpCtx.Router = router

		resp := httpCtx.GET("/api/v1/admin/xray/inbounds", nil)
		assert.Equal(t, 401, resp.Code)
	})

	t.Run("Forbidden_NotAdmin", func(t *testing.T) {
		defer app.CleanupTables()

		authHelper := testutil.NewAuthHelper(app.JWT, t)
		_, userToken, _ := authHelper.CreateAuthenticatedUser(app.Container.UserRepository)

		httpCtx := testutil.NewHTTPTestContext(t)
		httpCtx.Router = router

		resp := httpCtx.GET("/api/v1/admin/xray/inbounds", testutil.AuthHeader(userToken))
		assert.Equal(t, 403, resp.Code)
	})
}

func TestAdminHandler_GetInbound(t *testing.T) {
	if !testutil.IsIntegrationTest() {
		t.Skip("Skipping integration test")
	}

	app := testutil.NewTestApp(t)
	defer app.Cleanup()
	defer app.CleanupTables()

	router := setupAdminInboundHandler(t, app)

	t.Run("Success", func(t *testing.T) {
		defer app.CleanupTables()

		ctx := context.Background()
		authHelper := testutil.NewAuthHelper(app.JWT, t)
		_, adminToken, _ := authHelper.CreateAuthenticatedAdmin(app.Container.UserRepository)

		// Create dependencies - server and node first
		testServer, _ := testutil.CreateTestServerWithDefaults()
		_ = app.Container.ServerRepository.Create(ctx, testServer) // Test setup

		testNode, _ := testutil.CreateTestNodeWithDefaults(testServer.ID)
		_ = app.Container.NodeRepository.Create(ctx, testNode) // Test setup

		instance, _ := testutil.CreateTestXrayInstanceWithDefaults(testNode.ID)
		_ = app.Container.XrayInstanceRepository.Create(ctx, instance) // Test setup

		inbound, _ := testutil.CreateTestInboundWithDefaults(instance.ID)
		_ = app.Container.InboundRepository.Create(ctx, inbound) // Test setup

		httpCtx := testutil.NewHTTPTestContext(t)
		httpCtx.Router = router

		resp := httpCtx.GET("/api/v1/admin/xray/inbounds/"+inbound.ID.String(), testutil.AuthHeader(adminToken))

		assert.Equal(t, 200, resp.Code)

		var result response.Response
		httpCtx.GetResponseJSON(&result)
		assert.True(t, result.Success)
	})

	t.Run("NotFound_InboundDoesNotExist", func(t *testing.T) {
		defer app.CleanupTables()

		authHelper := testutil.NewAuthHelper(app.JWT, t)
		_, adminToken, _ := authHelper.CreateAuthenticatedAdmin(app.Container.UserRepository)

		httpCtx := testutil.NewHTTPTestContext(t)
		httpCtx.Router = router

		nonExistentID := uuid.New()
		resp := httpCtx.GET("/api/v1/admin/xray/inbounds/"+nonExistentID.String(), testutil.AuthHeader(adminToken))

		assert.Equal(t, 404, resp.Code)
	})

	t.Run("BadRequest_InvalidUUID", func(t *testing.T) {
		defer app.CleanupTables()

		authHelper := testutil.NewAuthHelper(app.JWT, t)
		_, adminToken, _ := authHelper.CreateAuthenticatedAdmin(app.Container.UserRepository)

		httpCtx := testutil.NewHTTPTestContext(t)
		httpCtx.Router = router

		resp := httpCtx.GET("/api/v1/admin/xray/inbounds/invalid-uuid", testutil.AuthHeader(adminToken))
		assert.Equal(t, 400, resp.Code)
	})
}

func TestAdminHandler_CreateInbound(t *testing.T) {
	if !testutil.IsIntegrationTest() {
		t.Skip("Skipping integration test")
	}

	app := testutil.NewTestApp(t)
	defer app.Cleanup()
	defer app.CleanupTables()

	router := setupAdminInboundHandler(t, app)

	t.Run("Success", func(t *testing.T) {
		defer app.CleanupTables()

		ctx := context.Background()
		authHelper := testutil.NewAuthHelper(app.JWT, t)
		_, adminToken, _ := authHelper.CreateAuthenticatedAdmin(app.Container.UserRepository)

		// Create dependencies - server and node first
		testServer, _ := testutil.CreateTestServerWithDefaults()
		_ = app.Container.ServerRepository.Create(ctx, testServer) // Test setup

		testNode, _ := testutil.CreateTestNodeWithDefaults(testServer.ID)
		_ = app.Container.NodeRepository.Create(ctx, testNode) // Test setup

		instance, _ := testutil.CreateTestXrayInstanceWithDefaults(testNode.ID)
		_ = app.Container.XrayInstanceRepository.Create(ctx, instance) // Test setup

		httpCtx := testutil.NewHTTPTestContext(t)
		httpCtx.Router = router

		req := dto.AdminCreateInboundRequest{
			XrayInstanceID: instance.ID.String(),
			Protocol:       string(xray.ProtocolVLESS),
			Port:           8443,
			Transport:      string(xray.TransportTCP),
			Security:       string(xray.SecurityREALITY),
		}

		resp := httpCtx.POST("/api/v1/admin/xray/inbounds", req, testutil.AuthHeader(adminToken))

		assert.Equal(t, 201, resp.Code)
	})

	t.Run("NotFound_InstanceDoesNotExist", func(t *testing.T) {
		defer app.CleanupTables()

		authHelper := testutil.NewAuthHelper(app.JWT, t)
		_, adminToken, _ := authHelper.CreateAuthenticatedAdmin(app.Container.UserRepository)

		httpCtx := testutil.NewHTTPTestContext(t)
		httpCtx.Router = router

		req := dto.AdminCreateInboundRequest{
			XrayInstanceID: uuid.New().String(),
			Protocol:       string(xray.ProtocolVLESS),
			Port:           8443,
			Transport:      string(xray.TransportTCP),
			Security:       string(xray.SecurityREALITY),
		}

		resp := httpCtx.POST("/api/v1/admin/xray/inbounds", req, testutil.AuthHeader(adminToken))

		assert.Equal(t, 404, resp.Code)
	})

	t.Run("BadRequest_InvalidPort", func(t *testing.T) {
		defer app.CleanupTables()

		ctx := context.Background()
		authHelper := testutil.NewAuthHelper(app.JWT, t)
		_, adminToken, _ := authHelper.CreateAuthenticatedAdmin(app.Container.UserRepository)

		// Create dependencies - server and node first
		testServer, _ := testutil.CreateTestServerWithDefaults()
		_ = app.Container.ServerRepository.Create(ctx, testServer) // Test setup

		testNode, _ := testutil.CreateTestNodeWithDefaults(testServer.ID)
		_ = app.Container.NodeRepository.Create(ctx, testNode) // Test setup

		instance, _ := testutil.CreateTestXrayInstanceWithDefaults(testNode.ID)
		_ = app.Container.XrayInstanceRepository.Create(ctx, instance) // Test setup

		httpCtx := testutil.NewHTTPTestContext(t)
		httpCtx.Router = router

		req := dto.AdminCreateInboundRequest{
			XrayInstanceID: instance.ID.String(),
			Protocol:       string(xray.ProtocolVLESS),
			Port:           0, // Invalid port
			Transport:      string(xray.TransportTCP),
			Security:       string(xray.SecurityREALITY),
		}

		resp := httpCtx.POST("/api/v1/admin/xray/inbounds", req, testutil.AuthHeader(adminToken))

		assert.Equal(t, 400, resp.Code)
	})
}

func TestAdminHandler_DeleteInbound(t *testing.T) {
	if !testutil.IsIntegrationTest() {
		t.Skip("Skipping integration test")
	}

	app := testutil.NewTestApp(t)
	defer app.Cleanup()
	defer app.CleanupTables()

	router := setupAdminInboundHandler(t, app)

	t.Run("Success", func(t *testing.T) {
		defer app.CleanupTables()

		ctx := context.Background()
		authHelper := testutil.NewAuthHelper(app.JWT, t)
		_, adminToken, _ := authHelper.CreateAuthenticatedAdmin(app.Container.UserRepository)

		// Create dependencies - server and node first
		testServer, _ := testutil.CreateTestServerWithDefaults()
		_ = app.Container.ServerRepository.Create(ctx, testServer) // Test setup

		testNode, _ := testutil.CreateTestNodeWithDefaults(testServer.ID)
		_ = app.Container.NodeRepository.Create(ctx, testNode) // Test setup

		instance, _ := testutil.CreateTestXrayInstanceWithDefaults(testNode.ID)
		_ = app.Container.XrayInstanceRepository.Create(ctx, instance) // Test setup

		inbound, _ := testutil.CreateTestInboundWithDefaults(instance.ID)
		_ = app.Container.InboundRepository.Create(ctx, inbound) // Test setup

		httpCtx := testutil.NewHTTPTestContext(t)
		httpCtx.Router = router

		resp := httpCtx.DELETE("/api/v1/admin/xray/inbounds/"+inbound.ID.String(), testutil.AuthHeader(adminToken))

		assert.Equal(t, 200, resp.Code)
	})

	t.Run("NotFound_InboundDoesNotExist", func(t *testing.T) {
		defer app.CleanupTables()

		authHelper := testutil.NewAuthHelper(app.JWT, t)
		_, adminToken, _ := authHelper.CreateAuthenticatedAdmin(app.Container.UserRepository)

		httpCtx := testutil.NewHTTPTestContext(t)
		httpCtx.Router = router

		nonExistentID := uuid.New()
		resp := httpCtx.DELETE("/api/v1/admin/xray/inbounds/"+nonExistentID.String(), testutil.AuthHeader(adminToken))

		assert.Equal(t, 404, resp.Code)
	})
}
