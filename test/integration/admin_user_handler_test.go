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
	"github.com/suproxy/backend/internal/domain/user"
	"github.com/suproxy/backend/internal/infrastructure/testutil"
	"github.com/suproxy/backend/internal/interfaces/http/handler"
	"github.com/suproxy/backend/internal/interfaces/http/middleware"
	"github.com/suproxy/backend/internal/interfaces/http/response"
)

func setupAdminUserHandler(t *testing.T, app *testutil.TestApp) (*gin.Engine, *handler.AdminHandler) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	// Create use-case instances for admin handler
	// Note: These are created directly as Container doesn't hold them
	// All dependencies are accessed via app.Container and app helpers
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

	// Setup middleware
	adminGroup := router.Group("/api/v1/admin")
	adminGroup.Use(middleware.AuthMiddleware(app.JWT))
	adminGroup.Use(middleware.RequireAdmin())

	// Register routes
	adminGroup.GET("/users", adminHandler.ListUsers)
	adminGroup.GET("/users/:id", adminHandler.GetUser)
	adminGroup.PUT("/users/:id/status", adminHandler.UpdateUserStatus)
	adminGroup.PUT("/users/:id/role", adminHandler.UpdateUserRole)

	return router, adminHandler
}

func TestAdminHandler_ListUsers(t *testing.T) {
	if !testutil.IsIntegrationTest() {
		t.Skip("Skipping integration test")
	}

	app := testutil.NewTestApp(t)
	defer app.Cleanup()
	defer app.CleanupTables()

	router, _ := setupAdminUserHandler(t, app)

	t.Run("Success", func(t *testing.T) {
		defer app.CleanupTables()

		ctx := context.Background()
		authHelper := testutil.NewAuthHelper(app.JWT, t)
		adminUser, adminToken, _ := authHelper.CreateAuthenticatedAdmin(app.Container.UserRepository)

		// Create some test users
		for i := 0; i < 3; i++ {
			testUser, _ := testutil.CreateTestUserWithDefaults()
			app.Container.UserRepository.Create(ctx, testUser)
		}

		httpCtx := testutil.NewHTTPTestContext(t)
		httpCtx.Router = router

		resp := httpCtx.GET("/api/v1/admin/users", testutil.AuthHeader(adminToken))

		assert.Equal(t, 200, resp.Code)

		var result response.Response
		httpCtx.GetResponseJSON(&result)
		assert.True(t, result.Success)

		_ = adminUser
	})

	t.Run("Success_WithFilters", func(t *testing.T) {
		defer app.CleanupTables()

		ctx := context.Background()
		authHelper := testutil.NewAuthHelper(app.JWT, t)
		adminUser, adminToken, _ := authHelper.CreateAuthenticatedAdmin(app.Container.UserRepository)

		// Create users with different roles
		for i := 0; i < 2; i++ {
			testUser, _ := testutil.CreateTestUserWithDefaults()
			app.Container.UserRepository.Create(ctx, testUser)
		}

		httpCtx := testutil.NewHTTPTestContext(t)
		httpCtx.Router = router

		resp := httpCtx.GET("/api/v1/admin/users?role=user&limit=10", testutil.AuthHeader(adminToken))

		assert.Equal(t, 200, resp.Code)

		_ = adminUser
	})

	t.Run("Unauthorized_NoToken", func(t *testing.T) {
		defer app.CleanupTables()

		httpCtx := testutil.NewHTTPTestContext(t)
		httpCtx.Router = router

		resp := httpCtx.GET("/api/v1/admin/users", nil)
		assert.Equal(t, 401, resp.Code)
	})

	t.Run("Forbidden_NotAdmin", func(t *testing.T) {
		defer app.CleanupTables()

		authHelper := testutil.NewAuthHelper(app.JWT, t)
		_, userToken, _ := authHelper.CreateAuthenticatedUser(app.Container.UserRepository)

		httpCtx := testutil.NewHTTPTestContext(t)
		httpCtx.Router = router

		resp := httpCtx.GET("/api/v1/admin/users", testutil.AuthHeader(userToken))
		assert.Equal(t, 403, resp.Code)
	})
}

func TestAdminHandler_GetUser(t *testing.T) {
	if !testutil.IsIntegrationTest() {
		t.Skip("Skipping integration test")
	}

	app := testutil.NewTestApp(t)
	defer app.Cleanup()
	defer app.CleanupTables()

	router, _ := setupAdminUserHandler(t, app)

	t.Run("Success", func(t *testing.T) {
		defer app.CleanupTables()

		ctx := context.Background()
		authHelper := testutil.NewAuthHelper(app.JWT, t)
		_, adminToken, _ := authHelper.CreateAuthenticatedAdmin(app.Container.UserRepository)

		testUser, _ := testutil.CreateTestUserWithDefaults()
		app.Container.UserRepository.Create(ctx, testUser)

		httpCtx := testutil.NewHTTPTestContext(t)
		httpCtx.Router = router

		resp := httpCtx.GET("/api/v1/admin/users/"+testUser.ID.String(), testutil.AuthHeader(adminToken))

		assert.Equal(t, 200, resp.Code)

		var result response.Response
		httpCtx.GetResponseJSON(&result)
		assert.True(t, result.Success)
	})

	t.Run("NotFound_UserDoesNotExist", func(t *testing.T) {
		defer app.CleanupTables()

		authHelper := testutil.NewAuthHelper(app.JWT, t)
		_, adminToken, _ := authHelper.CreateAuthenticatedAdmin(app.Container.UserRepository)

		httpCtx := testutil.NewHTTPTestContext(t)
		httpCtx.Router = router

		nonExistentID := uuid.New()
		resp := httpCtx.GET("/api/v1/admin/users/"+nonExistentID.String(), testutil.AuthHeader(adminToken))

		assert.Equal(t, 404, resp.Code)
	})

	t.Run("BadRequest_InvalidUUID", func(t *testing.T) {
		defer app.CleanupTables()

		authHelper := testutil.NewAuthHelper(app.JWT, t)
		_, adminToken, _ := authHelper.CreateAuthenticatedAdmin(app.Container.UserRepository)

		httpCtx := testutil.NewHTTPTestContext(t)
		httpCtx.Router = router

		resp := httpCtx.GET("/api/v1/admin/users/invalid-uuid", testutil.AuthHeader(adminToken))
		assert.Equal(t, 400, resp.Code)
	})

	t.Run("Forbidden_NotAdmin", func(t *testing.T) {
		defer app.CleanupTables()

		ctx := context.Background()
		authHelper := testutil.NewAuthHelper(app.JWT, t)
		_, userToken, _ := authHelper.CreateAuthenticatedUser(app.Container.UserRepository)

		testUser, _ := testutil.CreateTestUserWithDefaults()
		app.Container.UserRepository.Create(ctx, testUser)

		httpCtx := testutil.NewHTTPTestContext(t)
		httpCtx.Router = router

		resp := httpCtx.GET("/api/v1/admin/users/"+testUser.ID.String(), testutil.AuthHeader(userToken))
		assert.Equal(t, 403, resp.Code)
	})
}

func TestAdminHandler_UpdateUserStatus(t *testing.T) {
	if !testutil.IsIntegrationTest() {
		t.Skip("Skipping integration test")
	}

	app := testutil.NewTestApp(t)
	defer app.Cleanup()
	defer app.CleanupTables()

	router, _ := setupAdminUserHandler(t, app)

	t.Run("Success_Deactivate", func(t *testing.T) {
		defer app.CleanupTables()

		ctx := context.Background()
		authHelper := testutil.NewAuthHelper(app.JWT, t)
		_, adminToken, _ := authHelper.CreateAuthenticatedAdmin(app.Container.UserRepository)

		testUser, _ := testutil.CreateTestUserWithDefaults()
		app.Container.UserRepository.Create(ctx, testUser)

		httpCtx := testutil.NewHTTPTestContext(t)
		httpCtx.Router = router

		req := dto.UpdateUserStatusRequest{
			Status: string(user.StatusInactive),
		}

		resp := httpCtx.PUT("/api/v1/admin/users/"+testUser.ID.String()+"/status", req, testutil.AuthHeader(adminToken))

		assert.Equal(t, 200, resp.Code)

		var result response.Response
		httpCtx.GetResponseJSON(&result)
		assert.True(t, result.Success)
	})

	t.Run("Success_Suspend", func(t *testing.T) {
		defer app.CleanupTables()

		ctx := context.Background()
		authHelper := testutil.NewAuthHelper(app.JWT, t)
		_, adminToken, _ := authHelper.CreateAuthenticatedAdmin(app.Container.UserRepository)

		testUser, _ := testutil.CreateTestUserWithDefaults()
		app.Container.UserRepository.Create(ctx, testUser)

		httpCtx := testutil.NewHTTPTestContext(t)
		httpCtx.Router = router

		req := dto.UpdateUserStatusRequest{
			Status: string(user.StatusSuspended),
		}

		resp := httpCtx.PUT("/api/v1/admin/users/"+testUser.ID.String()+"/status", req, testutil.AuthHeader(adminToken))

		assert.Equal(t, 200, resp.Code)
	})

	t.Run("NotFound_UserDoesNotExist", func(t *testing.T) {
		defer app.CleanupTables()

		authHelper := testutil.NewAuthHelper(app.JWT, t)
		_, adminToken, _ := authHelper.CreateAuthenticatedAdmin(app.Container.UserRepository)

		httpCtx := testutil.NewHTTPTestContext(t)
		httpCtx.Router = router

		req := dto.UpdateUserStatusRequest{
			Status: string(user.StatusInactive),
		}

		nonExistentID := uuid.New()
		resp := httpCtx.PUT("/api/v1/admin/users/"+nonExistentID.String()+"/status", req, testutil.AuthHeader(adminToken))

		assert.Equal(t, 404, resp.Code)
	})

	t.Run("BadRequest_InvalidStatus", func(t *testing.T) {
		defer app.CleanupTables()

		ctx := context.Background()
		authHelper := testutil.NewAuthHelper(app.JWT, t)
		_, adminToken, _ := authHelper.CreateAuthenticatedAdmin(app.Container.UserRepository)

		testUser, _ := testutil.CreateTestUserWithDefaults()
		app.Container.UserRepository.Create(ctx, testUser)

		httpCtx := testutil.NewHTTPTestContext(t)
		httpCtx.Router = router

		req := dto.UpdateUserStatusRequest{
			Status: "invalid_status",
		}

		resp := httpCtx.PUT("/api/v1/admin/users/"+testUser.ID.String()+"/status", req, testutil.AuthHeader(adminToken))

		assert.Equal(t, 400, resp.Code)
	})

	t.Run("Forbidden_NotAdmin", func(t *testing.T) {
		defer app.CleanupTables()

		ctx := context.Background()
		authHelper := testutil.NewAuthHelper(app.JWT, t)
		_, userToken, _ := authHelper.CreateAuthenticatedUser(app.Container.UserRepository)

		testUser, _ := testutil.CreateTestUserWithDefaults()
		app.Container.UserRepository.Create(ctx, testUser)

		httpCtx := testutil.NewHTTPTestContext(t)
		httpCtx.Router = router

		req := dto.UpdateUserStatusRequest{
			Status: string(user.StatusInactive),
		}

		resp := httpCtx.PUT("/api/v1/admin/users/"+testUser.ID.String()+"/status", req, testutil.AuthHeader(userToken))
		assert.Equal(t, 403, resp.Code)
	})
}

func TestAdminHandler_UpdateUserRole(t *testing.T) {
	if !testutil.IsIntegrationTest() {
		t.Skip("Skipping integration test")
	}

	app := testutil.NewTestApp(t)
	defer app.Cleanup()
	defer app.CleanupTables()

	router, _ := setupAdminUserHandler(t, app)

	t.Run("Success_PromoteToAdmin", func(t *testing.T) {
		defer app.CleanupTables()

		ctx := context.Background()
		authHelper := testutil.NewAuthHelper(app.JWT, t)
		_, adminToken, _ := authHelper.CreateAuthenticatedAdmin(app.Container.UserRepository)

		testUser, _ := testutil.CreateTestUserWithDefaults()
		app.Container.UserRepository.Create(ctx, testUser)

		httpCtx := testutil.NewHTTPTestContext(t)
		httpCtx.Router = router

		req := dto.UpdateUserRoleRequest{
			Role: string(user.RoleAdmin),
		}

		resp := httpCtx.PUT("/api/v1/admin/users/"+testUser.ID.String()+"/role", req, testutil.AuthHeader(adminToken))

		assert.Equal(t, 200, resp.Code)

		var result response.Response
		httpCtx.GetResponseJSON(&result)
		assert.True(t, result.Success)
	})

	t.Run("NotFound_UserDoesNotExist", func(t *testing.T) {
		defer app.CleanupTables()

		authHelper := testutil.NewAuthHelper(app.JWT, t)
		_, adminToken, _ := authHelper.CreateAuthenticatedAdmin(app.Container.UserRepository)

		httpCtx := testutil.NewHTTPTestContext(t)
		httpCtx.Router = router

		req := dto.UpdateUserRoleRequest{
			Role: string(user.RoleAdmin),
		}

		nonExistentID := uuid.New()
		resp := httpCtx.PUT("/api/v1/admin/users/"+nonExistentID.String()+"/role", req, testutil.AuthHeader(adminToken))

		assert.Equal(t, 404, resp.Code)
	})

	t.Run("BadRequest_InvalidRole", func(t *testing.T) {
		defer app.CleanupTables()

		ctx := context.Background()
		authHelper := testutil.NewAuthHelper(app.JWT, t)
		_, adminToken, _ := authHelper.CreateAuthenticatedAdmin(app.Container.UserRepository)

		testUser, _ := testutil.CreateTestUserWithDefaults()
		app.Container.UserRepository.Create(ctx, testUser)

		httpCtx := testutil.NewHTTPTestContext(t)
		httpCtx.Router = router

		req := dto.UpdateUserRoleRequest{
			Role: "invalid_role",
		}

		resp := httpCtx.PUT("/api/v1/admin/users/"+testUser.ID.String()+"/role", req, testutil.AuthHeader(adminToken))

		assert.Equal(t, 400, resp.Code)
	})

	t.Run("Forbidden_NotAdmin", func(t *testing.T) {
		defer app.CleanupTables()

		ctx := context.Background()
		authHelper := testutil.NewAuthHelper(app.JWT, t)
		_, userToken, _ := authHelper.CreateAuthenticatedUser(app.Container.UserRepository)

		testUser, _ := testutil.CreateTestUserWithDefaults()
		app.Container.UserRepository.Create(ctx, testUser)

		httpCtx := testutil.NewHTTPTestContext(t)
		httpCtx.Router = router

		req := dto.UpdateUserRoleRequest{
			Role: string(user.RoleAdmin),
		}

		resp := httpCtx.PUT("/api/v1/admin/users/"+testUser.ID.String()+"/role", req, testutil.AuthHeader(userToken))
		assert.Equal(t, 403, resp.Code)
	})
}
