package integration

import (
	"context"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/suproxy/backend/internal/application/dto"
	"github.com/suproxy/backend/internal/infrastructure/testutil"
	"github.com/suproxy/backend/internal/interfaces/http/handler"
	"github.com/suproxy/backend/internal/interfaces/http/middleware"
	"github.com/suproxy/backend/internal/interfaces/http/response"
)

func setupAdminClientHandler(t *testing.T, app *testutil.TestApp) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	adminHandler := handler.NewAdminHandler(
		app.Logger,
		nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil,
		nil, nil, nil, nil, nil, nil, nil,
		nil, nil, nil, nil, nil, nil, nil, nil,
		nil, nil, nil, nil, nil, nil, nil, nil,
	)

	adminGroup := router.Group("/api/v1/admin/xray")
	adminGroup.Use(middleware.AuthMiddleware(app.JWT))
	adminGroup.Use(middleware.AdminAuthorization(app.Logger))

	adminGroup.GET("/clients", adminHandler.ListClients)
	adminGroup.GET("/clients/:id", adminHandler.GetClient)
	adminGroup.POST("/clients", adminHandler.CreateClient)
	adminGroup.DELETE("/clients/:id", adminHandler.DeleteClient)
	adminGroup.POST("/clients/:id/enable", adminHandler.EnableClient)
	adminGroup.POST("/clients/:id/disable", adminHandler.DisableClient)

	return router
}

func TestAdminHandler_ListClients(t *testing.T) {
	if !testutil.IsIntegrationTest() {
		t.Skip("Skipping integration test")
	}

	app := testutil.NewTestApp(t)
	defer app.Cleanup()
	defer app.CleanupTables()

	router := setupAdminClientHandler(t, app)

	t.Run("Success", func(t *testing.T) {
		defer app.CleanupTables()

		ctx := context.Background()
		authHelper := testutil.NewAuthHelper(app.JWT, t)
		_, adminToken, _ := authHelper.CreateAuthenticatedAdmin(app.Container.UserRepository)

		// Create dependencies - start with server and node
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
		err = app.Container.UserRepository.Create(ctx, user)
		require.NoError(t, err)

		instance, err := testutil.CreateTestXrayInstanceWithDefaults(testNode.ID)
		require.NoError(t, err)
		err = app.Container.XrayInstanceRepository.Create(ctx, instance)
		require.NoError(t, err)

		inbound, err := testutil.CreateTestInboundWithDefaults(instance.ID)
		require.NoError(t, err)
		err = app.Container.InboundRepository.Create(ctx, inbound)
		require.NoError(t, err)

		client, err := testutil.CreateTestClientWithDefaults(inbound.ID, user.ID)
		require.NoError(t, err)
		err = app.Container.ClientRepository.Create(ctx, client)
		require.NoError(t, err)

		httpCtx := testutil.NewHTTPTestContext(t)
		httpCtx.Router = router

		resp := httpCtx.GET("/api/v1/admin/xray/clients", testutil.AuthHeader(adminToken))

		assert.Equal(t, 200, resp.Code)

		var result response.Response
		httpCtx.GetResponseJSON(&result)
		assert.True(t, result.Success)
	})

	t.Run("Unauthorized_NoToken", func(t *testing.T) {
		defer app.CleanupTables()

		httpCtx := testutil.NewHTTPTestContext(t)
		httpCtx.Router = router

		resp := httpCtx.GET("/api/v1/admin/xray/clients", nil)
		assert.Equal(t, 401, resp.Code)
	})

	t.Run("Forbidden_NotAdmin", func(t *testing.T) {
		defer app.CleanupTables()

		authHelper := testutil.NewAuthHelper(app.JWT, t)
		_, userToken, _ := authHelper.CreateAuthenticatedUser(app.Container.UserRepository)

		httpCtx := testutil.NewHTTPTestContext(t)
		httpCtx.Router = router

		resp := httpCtx.GET("/api/v1/admin/xray/clients", testutil.AuthHeader(userToken))
		assert.Equal(t, 403, resp.Code)
	})
}

func TestAdminHandler_GetClient(t *testing.T) {
	if !testutil.IsIntegrationTest() {
		t.Skip("Skipping integration test")
	}

	app := testutil.NewTestApp(t)
	defer app.Cleanup()
	defer app.CleanupTables()

	router := setupAdminClientHandler(t, app)

	t.Run("Success", func(t *testing.T) {
		defer app.CleanupTables()

		ctx := context.Background()
		authHelper := testutil.NewAuthHelper(app.JWT, t)
		_, adminToken, _ := authHelper.CreateAuthenticatedAdmin(app.Container.UserRepository)

		// Create dependencies - start with server and node
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
		err = app.Container.UserRepository.Create(ctx, user)
		require.NoError(t, err)

		instance, err := testutil.CreateTestXrayInstanceWithDefaults(testNode.ID)
		require.NoError(t, err)
		err = app.Container.XrayInstanceRepository.Create(ctx, instance)
		require.NoError(t, err)

		inbound, err := testutil.CreateTestInboundWithDefaults(instance.ID)
		require.NoError(t, err)
		err = app.Container.InboundRepository.Create(ctx, inbound)
		require.NoError(t, err)

		client, err := testutil.CreateTestClientWithDefaults(inbound.ID, user.ID)
		require.NoError(t, err)
		err = app.Container.ClientRepository.Create(ctx, client)
		require.NoError(t, err)

		httpCtx := testutil.NewHTTPTestContext(t)
		httpCtx.Router = router

		resp := httpCtx.GET("/api/v1/admin/xray/clients/"+client.ID.String(), testutil.AuthHeader(adminToken))

		assert.Equal(t, 200, resp.Code)

		var result response.Response
		httpCtx.GetResponseJSON(&result)
		assert.True(t, result.Success)
	})

	t.Run("NotFound_ClientDoesNotExist", func(t *testing.T) {
		defer app.CleanupTables()

		authHelper := testutil.NewAuthHelper(app.JWT, t)
		_, adminToken, _ := authHelper.CreateAuthenticatedAdmin(app.Container.UserRepository)

		httpCtx := testutil.NewHTTPTestContext(t)
		httpCtx.Router = router

		nonExistentID := uuid.New()
		resp := httpCtx.GET("/api/v1/admin/xray/clients/"+nonExistentID.String(), testutil.AuthHeader(adminToken))

		assert.Equal(t, 404, resp.Code)
	})

	t.Run("BadRequest_InvalidUUID", func(t *testing.T) {
		defer app.CleanupTables()

		authHelper := testutil.NewAuthHelper(app.JWT, t)
		_, adminToken, _ := authHelper.CreateAuthenticatedAdmin(app.Container.UserRepository)

		httpCtx := testutil.NewHTTPTestContext(t)
		httpCtx.Router = router

		resp := httpCtx.GET("/api/v1/admin/xray/clients/invalid-uuid", testutil.AuthHeader(adminToken))
		assert.Equal(t, 400, resp.Code)
	})
}

func TestAdminHandler_CreateClient(t *testing.T) {
	if !testutil.IsIntegrationTest() {
		t.Skip("Skipping integration test")
	}

	app := testutil.NewTestApp(t)
	defer app.Cleanup()
	defer app.CleanupTables()

	router := setupAdminClientHandler(t, app)

	t.Run("Success", func(t *testing.T) {
		defer app.CleanupTables()

		ctx := context.Background()
		authHelper := testutil.NewAuthHelper(app.JWT, t)
		_, adminToken, _ := authHelper.CreateAuthenticatedAdmin(app.Container.UserRepository)

		// Create dependencies - start with server and node
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
		err = app.Container.UserRepository.Create(ctx, user)
		require.NoError(t, err)

		instance, err := testutil.CreateTestXrayInstanceWithDefaults(testNode.ID)
		require.NoError(t, err)
		err = app.Container.XrayInstanceRepository.Create(ctx, instance)
		require.NoError(t, err)

		inbound, err := testutil.CreateTestInboundWithDefaults(instance.ID)
		require.NoError(t, err)
		err = app.Container.InboundRepository.Create(ctx, inbound)
		require.NoError(t, err)

		httpCtx := testutil.NewHTTPTestContext(t)
		httpCtx.Router = router

		req := dto.AdminCreateClientRequest{
			InboundID: inbound.ID.String(),
			UserID:    user.ID.String(),
			Email:     "testclient@example.com",
		}

		resp := httpCtx.POST("/api/v1/admin/xray/clients", req, testutil.AuthHeader(adminToken))

		assert.Equal(t, 201, resp.Code)
	})

	t.Run("NotFound_InboundDoesNotExist", func(t *testing.T) {
		defer app.CleanupTables()

		ctx := context.Background()
		authHelper := testutil.NewAuthHelper(app.JWT, t)
		_, adminToken, _ := authHelper.CreateAuthenticatedAdmin(app.Container.UserRepository)

		user, _ := testutil.CreateTestUserWithDefaults()
		app.Container.UserRepository.Create(ctx, user)

		httpCtx := testutil.NewHTTPTestContext(t)
		httpCtx.Router = router

		req := dto.AdminCreateClientRequest{
			InboundID: uuid.New().String(),
			UserID:    user.ID.String(),
			Email:     "testclient@example.com",
		}

		resp := httpCtx.POST("/api/v1/admin/xray/clients", req, testutil.AuthHeader(adminToken))

		assert.Equal(t, 404, resp.Code)
	})
}

func TestAdminHandler_DeleteClient(t *testing.T) {
	if !testutil.IsIntegrationTest() {
		t.Skip("Skipping integration test")
	}

	app := testutil.NewTestApp(t)
	defer app.Cleanup()
	defer app.CleanupTables()

	router := setupAdminClientHandler(t, app)

	t.Run("Success", func(t *testing.T) {
		defer app.CleanupTables()

		ctx := context.Background()
		authHelper := testutil.NewAuthHelper(app.JWT, t)
		_, adminToken, _ := authHelper.CreateAuthenticatedAdmin(app.Container.UserRepository)

		user, _ := testutil.CreateTestUserWithDefaults()
		app.Container.UserRepository.Create(ctx, user)

		instance, _ := testutil.CreateTestXrayInstanceWithDefaults()
		app.Container.XrayInstanceRepository.Create(ctx, instance)

		inbound, _ := testutil.CreateTestInboundWithDefaults(instance.ID)
		app.Container.InboundRepository.Create(ctx, inbound)

		client, _ := testutil.CreateTestClientWithDefaults(inbound.ID, user.ID)
		app.Container.ClientRepository.Create(ctx, client)

		httpCtx := testutil.NewHTTPTestContext(t)
		httpCtx.Router = router

		resp := httpCtx.DELETE("/api/v1/admin/xray/clients/"+client.ID.String(), testutil.AuthHeader(adminToken))

		assert.Equal(t, 200, resp.Code)
	})

	t.Run("NotFound_ClientDoesNotExist", func(t *testing.T) {
		defer app.CleanupTables()

		authHelper := testutil.NewAuthHelper(app.JWT, t)
		_, adminToken, _ := authHelper.CreateAuthenticatedAdmin(app.Container.UserRepository)

		httpCtx := testutil.NewHTTPTestContext(t)
		httpCtx.Router = router

		nonExistentID := uuid.New()
		resp := httpCtx.DELETE("/api/v1/admin/xray/clients/"+nonExistentID.String(), testutil.AuthHeader(adminToken))

		assert.Equal(t, 404, resp.Code)
	})
}
