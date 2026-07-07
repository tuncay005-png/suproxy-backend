package integration

import (
	"context"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/suproxy/backend/internal/infrastructure/testutil"
	"github.com/suproxy/backend/internal/interfaces/http/handler"
	"github.com/suproxy/backend/internal/interfaces/http/middleware"
	"github.com/suproxy/backend/internal/interfaces/http/response"
)

func setupAdminXrayHandler(t *testing.T, app *testutil.TestApp) *gin.Engine {
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

	authMiddleware := middleware.NewAuthMiddleware(app.JWT, app.Container.Logger)
	adminMiddleware := middleware.NewAdminMiddleware(app.Container.Logger)

	adminGroup := router.Group("/api/v1/admin/xray")
	adminGroup.Use(authMiddleware.Authenticate())
	adminGroup.Use(adminMiddleware.RequireAdmin())

	adminGroup.GET("/instances", adminHandler.ListInstances)
	adminGroup.GET("/instances/:id", adminHandler.GetInstance)
	adminGroup.POST("/instances/:id/start", adminHandler.StartInstance)
	adminGroup.POST("/instances/:id/stop", adminHandler.StopInstance)
	adminGroup.POST("/instances/:id/restart", adminHandler.RestartInstance)
	adminGroup.POST("/instances/:id/reload", adminHandler.ReloadInstance)
	adminGroup.GET("/instances/:id/health", adminHandler.CheckInstanceHealth)
	adminGroup.GET("/instances/:id/stats", adminHandler.GetInstanceStats)

	return router
}

func TestAdminHandler_ListInstances(t *testing.T) {
	if !testutil.IsIntegrationTest() {
		t.Skip("Skipping integration test")
	}

	app := testutil.NewTestApp(t)
	defer app.Cleanup()
	defer app.CleanupTables()

	router := setupAdminXrayHandler(t, app)

	t.Run("Success", func(t *testing.T) {
		defer app.CleanupTables()

		ctx := context.Background()
		authHelper := testutil.NewAuthHelper(app.JWT, t)
		_, adminToken, _ := authHelper.CreateAuthenticatedAdmin(app.Container.UserRepository)

		// Create test instance
		instance, _ := testutil.CreateTestXrayInstanceWithDefaults()
		app.Container.XrayInstanceRepository.Create(ctx, instance)

		httpCtx := testutil.NewHTTPTestContext(t)
		httpCtx.Router = router

		resp := httpCtx.GET("/api/v1/admin/xray/instances", testutil.AuthHeader(adminToken))

		assert.Equal(t, 200, resp.Code)

		var result response.Response
		httpCtx.GetResponseJSON(&result)
		assert.True(t, result.Success)
	})

	t.Run("Success_WithFilters", func(t *testing.T) {
		defer app.CleanupTables()

		ctx := context.Background()
		authHelper := testutil.NewAuthHelper(app.JWT, t)
		_, adminToken, _ := authHelper.CreateAuthenticatedAdmin(app.Container.UserRepository)

		instance, _ := testutil.CreateTestXrayInstanceWithDefaults()
		app.Container.XrayInstanceRepository.Create(ctx, instance)

		httpCtx := testutil.NewHTTPTestContext(t)
		httpCtx.Router = router

		resp := httpCtx.GET("/api/v1/admin/xray/instances?status=running&limit=10", testutil.AuthHeader(adminToken))

		assert.Equal(t, 200, resp.Code)
	})

	t.Run("Unauthorized_NoToken", func(t *testing.T) {
		defer app.CleanupTables()

		httpCtx := testutil.NewHTTPTestContext(t)
		httpCtx.Router = router

		resp := httpCtx.GET("/api/v1/admin/xray/instances", nil)
		assert.Equal(t, 401, resp.Code)
	})

	t.Run("Forbidden_NotAdmin", func(t *testing.T) {
		defer app.CleanupTables()

		authHelper := testutil.NewAuthHelper(app.JWT, t)
		_, userToken, _ := authHelper.CreateAuthenticatedUser(app.Container.UserRepository)

		httpCtx := testutil.NewHTTPTestContext(t)
		httpCtx.Router = router

		resp := httpCtx.GET("/api/v1/admin/xray/instances", testutil.AuthHeader(userToken))
		assert.Equal(t, 403, resp.Code)
	})
}

func TestAdminHandler_GetInstance(t *testing.T) {
	if !testutil.IsIntegrationTest() {
		t.Skip("Skipping integration test")
	}

	app := testutil.NewTestApp(t)
	defer app.Cleanup()
	defer app.CleanupTables()

	router := setupAdminXrayHandler(t, app)

	t.Run("Success", func(t *testing.T) {
		defer app.CleanupTables()

		ctx := context.Background()
		authHelper := testutil.NewAuthHelper(app.JWT, t)
		_, adminToken, _ := authHelper.CreateAuthenticatedAdmin(app.Container.UserRepository)

		instance, _ := testutil.CreateTestXrayInstanceWithDefaults()
		app.Container.XrayInstanceRepository.Create(ctx, instance)

		httpCtx := testutil.NewHTTPTestContext(t)
		httpCtx.Router = router

		resp := httpCtx.GET("/api/v1/admin/xray/instances/"+instance.ID.String(), testutil.AuthHeader(adminToken))

		assert.Equal(t, 200, resp.Code)

		var result response.Response
		httpCtx.GetResponseJSON(&result)
		assert.True(t, result.Success)
	})

	t.Run("NotFound_InstanceDoesNotExist", func(t *testing.T) {
		defer app.CleanupTables()

		authHelper := testutil.NewAuthHelper(app.JWT, t)
		_, adminToken, _ := authHelper.CreateAuthenticatedAdmin(app.Container.UserRepository)

		httpCtx := testutil.NewHTTPTestContext(t)
		httpCtx.Router = router

		nonExistentID := uuid.New()
		resp := httpCtx.GET("/api/v1/admin/xray/instances/"+nonExistentID.String(), testutil.AuthHeader(adminToken))

		assert.Equal(t, 404, resp.Code)
	})

	t.Run("BadRequest_InvalidUUID", func(t *testing.T) {
		defer app.CleanupTables()

		authHelper := testutil.NewAuthHelper(app.JWT, t)
		_, adminToken, _ := authHelper.CreateAuthenticatedAdmin(app.Container.UserRepository)

		httpCtx := testutil.NewHTTPTestContext(t)
		httpCtx.Router = router

		resp := httpCtx.GET("/api/v1/admin/xray/instances/invalid-uuid", testutil.AuthHeader(adminToken))
		assert.Equal(t, 400, resp.Code)
	})
}

func TestAdminHandler_GetInstanceStats(t *testing.T) {
	if !testutil.IsIntegrationTest() {
		t.Skip("Skipping integration test")
	}

	app := testutil.NewTestApp(t)
	defer app.Cleanup()
	defer app.CleanupTables()

	router := setupAdminXrayHandler(t, app)

	t.Run("Success", func(t *testing.T) {
		defer app.CleanupTables()

		ctx := context.Background()
		authHelper := testutil.NewAuthHelper(app.JWT, t)
		_, adminToken, _ := authHelper.CreateAuthenticatedAdmin(app.Container.UserRepository)

		instance, _ := testutil.CreateTestXrayInstanceWithDefaults()
		app.Container.XrayInstanceRepository.Create(ctx, instance)

		httpCtx := testutil.NewHTTPTestContext(t)
		httpCtx.Router = router

		resp := httpCtx.GET("/api/v1/admin/xray/instances/"+instance.ID.String()+"/stats", testutil.AuthHeader(adminToken))

		assert.Equal(t, 200, resp.Code)

		var result response.Response
		httpCtx.GetResponseJSON(&result)
		assert.True(t, result.Success)
	})

	t.Run("NotFound_InstanceDoesNotExist", func(t *testing.T) {
		defer app.CleanupTables()

		authHelper := testutil.NewAuthHelper(app.JWT, t)
		_, adminToken, _ := authHelper.CreateAuthenticatedAdmin(app.Container.UserRepository)

		httpCtx := testutil.NewHTTPTestContext(t)
		httpCtx.Router = router

		nonExistentID := uuid.New()
		resp := httpCtx.GET("/api/v1/admin/xray/instances/"+nonExistentID.String()+"/stats", testutil.AuthHeader(adminToken))

		assert.Equal(t, 404, resp.Code)
	})
}

func TestAdminHandler_StartInstance(t *testing.T) {
	if !testutil.IsIntegrationTest() {
		t.Skip("Skipping integration test")
	}

	app := testutil.NewTestApp(t)
	defer app.Cleanup()
	defer app.CleanupTables()

	router := setupAdminXrayHandler(t, app)

	t.Run("NotFound_InstanceDoesNotExist", func(t *testing.T) {
		defer app.CleanupTables()

		authHelper := testutil.NewAuthHelper(app.JWT, t)
		_, adminToken, _ := authHelper.CreateAuthenticatedAdmin(app.Container.UserRepository)

		httpCtx := testutil.NewHTTPTestContext(t)
		httpCtx.Router = router

		nonExistentID := uuid.New()
		resp := httpCtx.POST("/api/v1/admin/xray/instances/"+nonExistentID.String()+"/start", nil, testutil.AuthHeader(adminToken))

		assert.Equal(t, 404, resp.Code)
	})

	t.Run("BadRequest_InvalidUUID", func(t *testing.T) {
		defer app.CleanupTables()

		authHelper := testutil.NewAuthHelper(app.JWT, t)
		_, adminToken, _ := authHelper.CreateAuthenticatedAdmin(app.Container.UserRepository)

		httpCtx := testutil.NewHTTPTestContext(t)
		httpCtx.Router = router

		resp := httpCtx.POST("/api/v1/admin/xray/instances/invalid-uuid/start", nil, testutil.AuthHeader(adminToken))
		assert.Equal(t, 400, resp.Code)
	})
}

func TestAdminHandler_StopInstance(t *testing.T) {
	if !testutil.IsIntegrationTest() {
		t.Skip("Skipping integration test")
	}

	app := testutil.NewTestApp(t)
	defer app.Cleanup()
	defer app.CleanupTables()

	router := setupAdminXrayHandler(t, app)

	t.Run("NotFound_InstanceDoesNotExist", func(t *testing.T) {
		defer app.CleanupTables()

		authHelper := testutil.NewAuthHelper(app.JWT, t)
		_, adminToken, _ := authHelper.CreateAuthenticatedAdmin(app.Container.UserRepository)

		httpCtx := testutil.NewHTTPTestContext(t)
		httpCtx.Router = router

		nonExistentID := uuid.New()
		resp := httpCtx.POST("/api/v1/admin/xray/instances/"+nonExistentID.String()+"/stop", nil, testutil.AuthHeader(adminToken))

		assert.Equal(t, 404, resp.Code)
	})
}

func TestAdminHandler_RestartInstance(t *testing.T) {
	if !testutil.IsIntegrationTest() {
		t.Skip("Skipping integration test")
	}

	app := testutil.NewTestApp(t)
	defer app.Cleanup()
	defer app.CleanupTables()

	router := setupAdminXrayHandler(t, app)

	t.Run("NotFound_InstanceDoesNotExist", func(t *testing.T) {
		defer app.CleanupTables()

		authHelper := testutil.NewAuthHelper(app.JWT, t)
		_, adminToken, _ := authHelper.CreateAuthenticatedAdmin(app.Container.UserRepository)

		httpCtx := testutil.NewHTTPTestContext(t)
		httpCtx.Router = router

		nonExistentID := uuid.New()
		resp := httpCtx.POST("/api/v1/admin/xray/instances/"+nonExistentID.String()+"/restart", nil, testutil.AuthHeader(adminToken))

		assert.Equal(t, 404, resp.Code)
	})
}
