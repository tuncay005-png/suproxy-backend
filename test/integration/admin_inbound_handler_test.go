package integration

import (
	"context"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/suproxy/backend/internal/application/dto"
	"github.com/suproxy/backend/internal/domain/xray"
	"github.com/suproxy/backend/internal/infrastructure/testutil"
	"github.com/suproxy/backend/internal/interfaces/http/handler"
	"github.com/suproxy/backend/internal/interfaces/http/middleware"
	"github.com/suproxy/backend/internal/interfaces/http/response"
)

func setupAdminInboundHandler(t *testing.T, app *testutil.TestApp) *gin.Engine {
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

		// Create test instance and inbound
		instance, _ := testutil.CreateTestXrayInstanceWithDefaults()
		app.Container.XrayInstanceRepository.Create(ctx, instance)

		inbound, _ := testutil.CreateTestInboundWithDefaults(instance.ID)
		app.Container.InboundRepository.Create(ctx, inbound)

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

		instance, _ := testutil.CreateTestXrayInstanceWithDefaults()
		app.Container.XrayInstanceRepository.Create(ctx, instance)

		inbound, _ := testutil.CreateTestInboundWithDefaults(instance.ID)
		app.Container.InboundRepository.Create(ctx, inbound)

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

		instance, _ := testutil.CreateTestXrayInstanceWithDefaults()
		app.Container.XrayInstanceRepository.Create(ctx, instance)

		httpCtx := testutil.NewHTTPTestContext(t)
		httpCtx.Router = router

		req := dto.AdminCreateInboundRequest{
			XrayInstanceID: instance.ID.String(),
			Protocol:       string(xray.ProtocolVLESS),
			Port:           8443,
			Transport:      string(xray.TransportTCP),
			Security:       string(xray.SecurityReality),
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
			Security:       string(xray.SecurityReality),
		}

		resp := httpCtx.POST("/api/v1/admin/xray/inbounds", req, testutil.AuthHeader(adminToken))

		assert.Equal(t, 404, resp.Code)
	})

	t.Run("BadRequest_InvalidPort", func(t *testing.T) {
		defer app.CleanupTables()

		ctx := context.Background()
		authHelper := testutil.NewAuthHelper(app.JWT, t)
		_, adminToken, _ := authHelper.CreateAuthenticatedAdmin(app.Container.UserRepository)

		instance, _ := testutil.CreateTestXrayInstanceWithDefaults()
		app.Container.XrayInstanceRepository.Create(ctx, instance)

		httpCtx := testutil.NewHTTPTestContext(t)
		httpCtx.Router = router

		req := dto.AdminCreateInboundRequest{
			XrayInstanceID: instance.ID.String(),
			Protocol:       string(xray.ProtocolVLESS),
			Port:           0, // Invalid port
			Transport:      string(xray.TransportTCP),
			Security:       string(xray.SecurityReality),
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

		instance, _ := testutil.CreateTestXrayInstanceWithDefaults()
		app.Container.XrayInstanceRepository.Create(ctx, instance)

		inbound, _ := testutil.CreateTestInboundWithDefaults(instance.ID)
		app.Container.InboundRepository.Create(ctx, inbound)

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
