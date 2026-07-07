package integration

import (
	"context"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/suproxy/backend/internal/application/dto"
	"github.com/suproxy/backend/internal/domain/user"
	"github.com/suproxy/backend/internal/infrastructure/testutil"
	"github.com/suproxy/backend/internal/interfaces/http/handler"
	"github.com/suproxy/backend/internal/interfaces/http/middleware"
	"github.com/suproxy/backend/internal/interfaces/http/response"
)

func setupAdminUserHandler(t *testing.T, app *testutil.TestApp) (*gin.Engine, *handler.AdminHandler) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	adminHandler := handler.NewAdminHandler(
		app.Container.Logger,
		app.Container.AdminListUsersQuery,
		app.Container.AdminGetUserQuery,
		app.Container.AdminUpdateUserStatusCommand,
		app.Container.AdminUpdateUserRoleCommand,
		app.Container.AdminListInstancesQuery,
		app.Container.AdminGetInstanceQuery,
		app.Container.AdminGetInstanceStatsQuery,
		app.Container.AdminStartInstanceCommand,
		app.Container.AdminStopInstanceCommand,
		app.Container.AdminRestartInstanceCommand,
		app.Container.AdminReloadInstanceCommand,
		app.Container.AdminCheckInstanceHealthCommand,
		app.Container.AdminListInboundsQuery,
		app.Container.AdminGetInboundQuery,
		app.Container.AdminCreateInboundCommand,
		app.Container.AdminUpdateInboundCommand,
		app.Container.AdminDeleteInboundCommand,
		app.Container.AdminEnableInboundCommand,
		app.Container.AdminDisableInboundCommand,
		app.Container.AdminListClientsQuery,
		app.Container.AdminGetClientQuery,
		app.Container.AdminCreateClientCommand,
		app.Container.AdminDeleteClientCommand,
		app.Container.AdminEnableClientCommand,
		app.Container.AdminDisableClientCommand,
		app.Container.AdminRegenerateClientUUIDCommand,
		app.Container.AdminReprovisionClientCommand,
		app.Container.AdminListAuditLogsQuery,
		app.Container.AdminGetAuditLogQuery,
		app.Container.AdminGetAuditStatsQuery,
		app.Container.AdminGetSystemHealthQuery,
		app.Container.AdminGetSystemStatsQuery,
		app.Container.AdminGetVersionQuery,
		app.Container.AdminGetDatabaseStatusQuery,
		app.Container.AdminGetXraySystemStatusQuery,
	)

	// Setup middleware
	authMiddleware := middleware.NewAuthMiddleware(app.JWT, app.Container.Logger)
	adminMiddleware := middleware.NewAdminMiddleware(app.Container.Logger)

	adminGroup := router.Group("/api/v1/admin")
	adminGroup.Use(authMiddleware.Authenticate())
	adminGroup.Use(adminMiddleware.RequireAdmin())

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
