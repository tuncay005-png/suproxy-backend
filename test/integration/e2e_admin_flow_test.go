package integration

import (
	"context"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/suproxy/backend/internal/application/dto"
	adminaudit "github.com/suproxy/backend/internal/application/usecase/admin/audit"
	adminclient "github.com/suproxy/backend/internal/application/usecase/admin/client"
	admininbound "github.com/suproxy/backend/internal/application/usecase/admin/inbound"
	adminsystem "github.com/suproxy/backend/internal/application/usecase/admin/system"
	adminuser "github.com/suproxy/backend/internal/application/usecase/admin/user"
	adminxray "github.com/suproxy/backend/internal/application/usecase/admin/xray_instance"
	authuc "github.com/suproxy/backend/internal/application/usecase/auth"
	"github.com/suproxy/backend/internal/domain/user"
	"github.com/suproxy/backend/internal/infrastructure/testutil"
	"github.com/suproxy/backend/internal/interfaces/http/handler"
	"github.com/suproxy/backend/internal/interfaces/http/middleware"
	"github.com/suproxy/backend/internal/interfaces/http/response"
)

func setupE2EAdminRouter(t *testing.T, app *testutil.TestApp) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	// Create auth use-case instances
	registerCmd := authuc.NewRegisterCommand(app.Container.UserRepository, app.Container.XrayProvisioningService, app.Logger)
	loginCmd := authuc.NewLoginCommand(app.Container.UserRepository, app.Container.RefreshTokenRepository, app.Container.AuditLogRepository, app.JWT, app.Logger)
	refreshCmd := authuc.NewRefreshTokenCommand(app.Container.UserRepository, app.Container.RefreshTokenRepository, app.Container.AuditLogRepository, app.JWT, app.Logger)
	logoutCmd := authuc.NewLogoutCommand(app.Container.RefreshTokenRepository, app.Container.AuditLogRepository, app.Logger)
	getCurrentUserQuery := authuc.NewGetCurrentUserQuery(app.Container.UserRepository, app.Logger)
	getSessionsQuery := authuc.NewGetSessionsQuery(app.Container.RefreshTokenRepository, app.Logger)

	// Auth handler
	authHandler := handler.NewAuthHandler(
		registerCmd,
		loginCmd,
		refreshCmd,
		logoutCmd,
		getCurrentUserQuery,
		getSessionsQuery,
	)

	// Create admin use-case instances
	listUsersQuery := adminuser.NewListUsersQuery(app.Container.UserRepository)
	getUserQuery := adminuser.NewGetUserQuery(app.Container.UserRepository)
	updateUserStatusCmd := adminuser.NewUpdateUserStatusCommand(app.Container.UserRepository, app.Container.AuditLogRepository)
	updateUserRoleCmd := adminuser.NewUpdateUserRoleCommand(app.Container.UserRepository, app.Container.AuditLogRepository)
	
	listInstancesQuery := adminxray.NewListInstancesQuery(app.Container.XrayInstanceRepository)
	getInstanceQuery := adminxray.NewGetInstanceQuery(app.Container.XrayInstanceRepository)
	getInstanceStatsQuery := adminxray.NewGetInstanceStatsQuery(app.Container.XrayInstanceRepository, app.Container.InboundRepository, app.Container.ClientRepository)
	startInstanceCmd := adminxray.NewStartInstanceCommand(app.Container.XrayInstanceRepository, app.Container.XrayProcessManager, app.Container.AuditLogRepository)
	stopInstanceCmd := adminxray.NewStopInstanceCommand(app.Container.XrayInstanceRepository, app.Container.XrayProcessManager, app.Container.AuditLogRepository)
	restartInstanceCmd := adminxray.NewRestartInstanceCommand(app.Container.XrayInstanceRepository, app.Container.XrayProcessManager, app.Container.AuditLogRepository)
	reloadInstanceCmd := adminxray.NewReloadInstanceCommand(app.Container.XrayProvisioningService, app.Container.AuditLogRepository)
	checkInstanceHealthCmd := adminxray.NewCheckInstanceHealthCommand(app.Container.XrayProcessManager)
	
	listInboundsQuery := admininbound.NewListInboundsQuery(app.Container.InboundRepository)
	getInboundQuery := admininbound.NewGetInboundQuery(app.Container.InboundRepository)
	createInboundCmd := admininbound.NewCreateInboundCommand(app.Container.InboundRepository, app.Container.XrayInstanceRepository, app.Container.XrayProvisioningService, app.Container.AuditLogRepository)
	updateInboundCmd := admininbound.NewUpdateInboundCommand(app.Container.InboundRepository, app.Container.XrayProvisioningService, app.Container.AuditLogRepository)
	deleteInboundCmd := admininbound.NewDeleteInboundCommand(app.Container.InboundRepository, app.Container.ClientRepository, app.Container.XrayProvisioningService, app.Container.AuditLogRepository)
	enableInboundCmd := admininbound.NewEnableInboundCommand(app.Container.InboundRepository, app.Container.XrayProvisioningService, app.Container.AuditLogRepository)
	disableInboundCmd := admininbound.NewDisableInboundCommand(app.Container.InboundRepository, app.Container.XrayProvisioningService, app.Container.AuditLogRepository)
	
	listClientsQuery := adminclient.NewListClientsQuery(app.Container.ClientRepository)
	getClientQuery := adminclient.NewGetClientQuery(app.Container.ClientRepository)
	createClientCmd := adminclient.NewCreateClientCommand(app.Container.ClientRepository, app.Container.InboundRepository, app.Container.XrayProvisioningService, app.Container.AuditLogRepository)
	deleteClientCmd := adminclient.NewDeleteClientCommand(app.Container.ClientRepository, app.Container.InboundRepository, app.Container.XrayProvisioningService, app.Container.AuditLogRepository)
	enableClientCmd := adminclient.NewEnableClientCommand(app.Container.ClientRepository, app.Container.InboundRepository, app.Container.XrayProvisioningService, app.Container.AuditLogRepository)
	disableClientCmd := adminclient.NewDisableClientCommand(app.Container.ClientRepository, app.Container.InboundRepository, app.Container.XrayProvisioningService, app.Container.AuditLogRepository)
	regenerateClientUUIDCmd := adminclient.NewRegenerateClientUUIDCommand(app.Container.ClientRepository, app.Container.InboundRepository, app.Container.XrayProvisioningService, app.Container.AuditLogRepository)
	reprovisionClientCmd := adminclient.NewReprovisionClientCommand(app.Container.ClientRepository, app.Container.InboundRepository, app.Container.XrayProvisioningService, app.Container.AuditLogRepository)
	
	listAuditLogsQuery := adminaudit.NewListAuditLogsQuery(app.Container.AuditLogRepository)
	getAuditLogQuery := adminaudit.NewGetAuditLogQuery(app.Container.AuditLogRepository)
	getAuditStatsQuery := adminaudit.NewGetAuditStatsQuery(app.Container.AuditLogRepository)
	
	getSystemHealthQuery := adminsystem.NewGetSystemHealthQuery(app.Database, app.Container.XrayInstanceRepository, app.Container.XrayProcessManager)
	getSystemStatsQuery := adminsystem.NewGetSystemStatsQuery(app.Container.UserRepository, app.Container.XrayInstanceRepository, app.Container.InboundRepository, app.Container.ClientRepository, app.Container.AuditLogRepository)
	getVersionQuery := adminsystem.NewGetVersionQuery()
	getDatabaseStatusQuery := adminsystem.NewGetDatabaseStatusQuery(app.Database)
	getXraySystemStatusQuery := adminsystem.NewGetXraySystemStatusQuery(app.Container.XrayInstanceRepository)

	// Admin handler
	adminHandler := handler.NewAdminHandler(
		app.Logger,
		listUsersQuery,
		getUserQuery,
		updateUserStatusCmd,
		updateUserRoleCmd,
		listInstancesQuery,
		getInstanceQuery,
		getInstanceStatsQuery,
		startInstanceCmd,
		stopInstanceCmd,
		restartInstanceCmd,
		reloadInstanceCmd,
		checkInstanceHealthCmd,
		listInboundsQuery,
		getInboundQuery,
		createInboundCmd,
		updateInboundCmd,
		deleteInboundCmd,
		enableInboundCmd,
		disableInboundCmd,
		listClientsQuery,
		getClientQuery,
		createClientCmd,
		deleteClientCmd,
		enableClientCmd,
		disableClientCmd,
		regenerateClientUUIDCmd,
		reprovisionClientCmd,
		listAuditLogsQuery,
		getAuditLogQuery,
		getAuditStatsQuery,
		getSystemHealthQuery,
		getSystemStatsQuery,
		getVersionQuery,
		getDatabaseStatusQuery,
		getXraySystemStatusQuery,
	)

	// Middleware
	authMw := middleware.AuthMiddleware(app.JWT)
	adminMw := middleware.AdminAuthorization(app.Logger)

	// Auth routes
	authGroup := router.Group("/api/v1/auth")
	{
		authGroup.POST("/login", authHandler.Login)
	}

	// Admin routes
	adminGroup := router.Group("/api/v1/admin")
	adminGroup.Use(authMw)
	adminGroup.Use(adminMw)
	{
		// User management
		adminGroup.GET("/users", adminHandler.ListUsers)
		adminGroup.GET("/users/:id", adminHandler.GetUser)
		adminGroup.PUT("/users/:id/status", adminHandler.UpdateUserStatus)
		adminGroup.PUT("/users/:id/role", adminHandler.UpdateUserRole)

		// Xray management
		adminGroup.GET("/xray/instances", adminHandler.ListInstances)
		adminGroup.GET("/xray/instances/:id", adminHandler.GetInstance)
		adminGroup.GET("/xray/instances/:id/stats", adminHandler.GetInstanceStats)
		adminGroup.POST("/xray/instances/:id/start", adminHandler.StartInstance)
		adminGroup.POST("/xray/instances/:id/stop", adminHandler.StopInstance)
		adminGroup.POST("/xray/instances/:id/restart", adminHandler.RestartInstance)
		adminGroup.POST("/xray/instances/:id/reload", adminHandler.ReloadInstance)
		adminGroup.GET("/xray/instances/:id/health", adminHandler.CheckInstanceHealth)

		// Inbound management
		adminGroup.GET("/xray/inbounds", adminHandler.ListInbounds)
		adminGroup.GET("/xray/inbounds/:id", adminHandler.GetInbound)
		adminGroup.POST("/xray/inbounds", adminHandler.CreateInbound)
		adminGroup.DELETE("/xray/inbounds/:id", adminHandler.DeleteInbound)

		// Client management
		adminGroup.GET("/xray/clients", adminHandler.ListClients)
		adminGroup.GET("/xray/clients/:id", adminHandler.GetClient)
		adminGroup.POST("/xray/clients", adminHandler.CreateClient)
		adminGroup.DELETE("/xray/clients/:id", adminHandler.DeleteClient)
		adminGroup.POST("/xray/clients/:id/enable", adminHandler.EnableClient)
		adminGroup.POST("/xray/clients/:id/disable", adminHandler.DisableClient)

		// Audit
		adminGroup.GET("/audit", adminHandler.ListAuditLogs)
		adminGroup.GET("/audit/:id", adminHandler.GetAuditLog)
		adminGroup.GET("/audit/stats", adminHandler.GetAuditStats)

		// System
		adminGroup.GET("/health", adminHandler.HealthCheck)
		adminGroup.GET("/system/health", adminHandler.GetSystemHealth)
		adminGroup.GET("/system/stats", adminHandler.GetSystemStats)
	}

	return router
}

// TestE2E_AdminUserManagementFlow tests complete admin user management workflow
func TestE2E_AdminUserManagementFlow(t *testing.T) {
	if !testutil.IsIntegrationTest() {
		t.Skip("Skipping E2E integration test")
	}

	app := testutil.NewTestApp(t)
	defer app.Cleanup()
	defer app.CleanupTables()

	router := setupE2EAdminRouter(t, app)
	httpCtx := testutil.NewHTTPTestContext(t)
	httpCtx.Router = router

	ctx := context.Background()

	// Setup: Create admin user and regular users
	adminUser, _ := testutil.CreateTestAdminUser()
	app.Container.UserRepository.Create(ctx, adminUser)

	regularUser, _ := testutil.CreateTestUserWithDefaults()
	app.Container.UserRepository.Create(ctx, regularUser)

	// Step 1: Admin login
	t.Log("Step 1: Admin Login")
	loginReq := dto.LoginRequest{
		Email:    adminUser.Email.String(),
		Password: "Admin123!@#",
	}

	resp := httpCtx.POST("/api/v1/auth/login", loginReq, nil)
	require.Equal(t, 200, resp.Code)

	var loginResult response.Response
	httpCtx.GetResponseJSON(&loginResult)
	adminToken := loginResult.Data.(map[string]interface{})["access_token"].(string)
	require.NotEmpty(t, adminToken)

	// Step 2: List all users
	t.Log("Step 2: List All Users")
	resp = httpCtx.GET("/api/v1/admin/users", testutil.AuthHeader(adminToken))
	require.Equal(t, 200, resp.Code)

	var usersResult response.Response
	httpCtx.GetResponseJSON(&usersResult)
	require.True(t, usersResult.Success)

	// Step 3: Get specific user details
	t.Log("Step 3: Get User Details")
	resp = httpCtx.GET("/api/v1/admin/users/"+regularUser.ID.String(), testutil.AuthHeader(adminToken))
	require.Equal(t, 200, resp.Code)

	var userResult response.Response
	httpCtx.GetResponseJSON(&userResult)
	require.True(t, userResult.Success)

	// Step 4: Update user status (deactivate)
	t.Log("Step 4: Deactivate User")
	statusReq := dto.UpdateUserStatusRequest{
		Status: string(user.StatusInactive),
	}
	resp = httpCtx.PUT("/api/v1/admin/users/"+regularUser.ID.String()+"/status", statusReq, testutil.AuthHeader(adminToken))
	require.Equal(t, 200, resp.Code)

	// Step 5: Verify user is deactivated
	t.Log("Step 5: Verify User Status")
	resp = httpCtx.GET("/api/v1/admin/users/"+regularUser.ID.String(), testutil.AuthHeader(adminToken))
	require.Equal(t, 200, resp.Code)

	var verifyResult response.Response
	httpCtx.GetResponseJSON(&verifyResult)
	userData := verifyResult.Data.(map[string]interface{})
	assert.Equal(t, string(user.StatusInactive), userData["status"].(string))

	// Step 6: Promote user to admin
	t.Log("Step 6: Promote User to Admin")
	
	// First reactivate
	statusReq.Status = string(user.StatusActive)
	resp = httpCtx.PUT("/api/v1/admin/users/"+regularUser.ID.String()+"/status", statusReq, testutil.AuthHeader(adminToken))
	require.Equal(t, 200, resp.Code)

	roleReq := dto.UpdateUserRoleRequest{
		Role: string(user.RoleAdmin),
	}
	resp = httpCtx.PUT("/api/v1/admin/users/"+regularUser.ID.String()+"/role", roleReq, testutil.AuthHeader(adminToken))
	require.Equal(t, 200, resp.Code)

	// Step 7: Verify role change
	t.Log("Step 7: Verify User Role")
	resp = httpCtx.GET("/api/v1/admin/users/"+regularUser.ID.String(), testutil.AuthHeader(adminToken))
	require.Equal(t, 200, resp.Code)

	var roleResult response.Response
	httpCtx.GetResponseJSON(&roleResult)
	roleData := roleResult.Data.(map[string]interface{})
	assert.Equal(t, string(user.RoleAdmin), roleData["role"].(string))

	t.Log("✅ E2E Admin User Management Flow Complete")
}

// TestE2E_XrayProvisioningFlow tests complete xray provisioning workflow
func TestE2E_XrayProvisioningFlow(t *testing.T) {
	if !testutil.IsIntegrationTest() {
		t.Skip("Skipping E2E integration test")
	}

	app := testutil.NewTestApp(t)
	defer app.Cleanup()
	defer app.CleanupTables()

	router := setupE2EAdminRouter(t, app)
	httpCtx := testutil.NewHTTPTestContext(t)
	httpCtx.Router = router

	ctx := context.Background()

	// Setup: Create admin user
	adminUser, _ := testutil.CreateTestAdminUser()
	app.Container.UserRepository.Create(ctx, adminUser)

	// Admin login
	loginReq := dto.LoginRequest{
		Email:    adminUser.Email.String(),
		Password: "Admin123!@#",
	}
	resp := httpCtx.POST("/api/v1/auth/login", loginReq, nil)
	require.Equal(t, 200, resp.Code)

	var loginResult response.Response
	httpCtx.GetResponseJSON(&loginResult)
	adminToken := loginResult.Data.(map[string]interface{})["access_token"].(string)

	// Step 1: Create Xray Instance (via repository)
	t.Log("Step 1: Create Xray Instance")
	// Create dependencies - server and node first
	testServer, _ := testutil.CreateTestServerWithDefaults()
	app.Container.ServerRepository.Create(ctx, testServer)

	testNode, _ := testutil.CreateTestNodeWithDefaults(testServer.ID)
	app.Container.NodeRepository.Create(ctx, testNode)

	instance, _ := testutil.CreateTestXrayInstanceWithDefaults(testNode.ID)
	app.Container.XrayInstanceRepository.Create(ctx, instance)

	// Step 2: List instances
	t.Log("Step 2: List Xray Instances")
	resp = httpCtx.GET("/api/v1/admin/xray/instances", testutil.AuthHeader(adminToken))
	require.Equal(t, 200, resp.Code)

	var instancesResult response.Response
	httpCtx.GetResponseJSON(&instancesResult)
	require.True(t, instancesResult.Success)

	// Step 3: Get instance details
	t.Log("Step 3: Get Instance Details")
	resp = httpCtx.GET("/api/v1/admin/xray/instances/"+instance.ID.String(), testutil.AuthHeader(adminToken))
	require.Equal(t, 200, resp.Code)

	// Step 4: Create Inbound
	t.Log("Step 4: Create Inbound")
	createInboundReq := dto.AdminCreateInboundRequest{
		XrayInstanceID: instance.ID.String(),
		Protocol:       "vless",
		Port:           8443,
		Transport:      "tcp",
		Security:       "reality",
	}
	resp = httpCtx.POST("/api/v1/admin/xray/inbounds", createInboundReq, testutil.AuthHeader(adminToken))
	require.Equal(t, 201, resp.Code)

	var inboundResult response.Response
	httpCtx.GetResponseJSON(&inboundResult)
	require.True(t, inboundResult.Success)

	inboundData := inboundResult.Data.(map[string]interface{})
	inboundInfo := inboundData["inbound"].(map[string]interface{})
	inboundID := inboundInfo["id"].(string)
	require.NotEmpty(t, inboundID)

	// Step 5: List inbounds
	t.Log("Step 5: List Inbounds")
	resp = httpCtx.GET("/api/v1/admin/xray/inbounds", testutil.AuthHeader(adminToken))
	require.Equal(t, 200, resp.Code)

	// Step 6: Create Client
	t.Log("Step 6: Create Client")
	testUser, _ := testutil.CreateTestUserWithDefaults()
	app.Container.UserRepository.Create(ctx, testUser)

	createClientReq := dto.AdminCreateClientRequest{
		InboundID: inboundID,
		UserID:    testUser.ID.String(),
		Email:     "client@example.com",
	}
	resp = httpCtx.POST("/api/v1/admin/xray/clients", createClientReq, testutil.AuthHeader(adminToken))
	require.Equal(t, 201, resp.Code)

	var clientResult response.Response
	httpCtx.GetResponseJSON(&clientResult)
	require.True(t, clientResult.Success)

	// Step 7: Get instance stats
	t.Log("Step 7: Get Instance Stats")
	resp = httpCtx.GET("/api/v1/admin/xray/instances/"+instance.ID.String()+"/stats", testutil.AuthHeader(adminToken))
	require.Equal(t, 200, resp.Code)

	var statsResult response.Response
	httpCtx.GetResponseJSON(&statsResult)
	require.True(t, statsResult.Success)

	statsData := statsResult.Data.(map[string]interface{})
	totalInbounds := int(statsData["total_inbounds"].(float64))
	totalClients := int(statsData["total_clients"].(float64))
	assert.Equal(t, 1, totalInbounds)
	assert.Equal(t, 1, totalClients)

	t.Log("✅ E2E Xray Provisioning Flow Complete")
}

// TestE2E_ClientLifecycleFlow tests complete client lifecycle
func TestE2E_ClientLifecycleFlow(t *testing.T) {
	if !testutil.IsIntegrationTest() {
		t.Skip("Skipping E2E integration test")
	}

	app := testutil.NewTestApp(t)
	defer app.Cleanup()
	defer app.CleanupTables()

	router := setupE2EAdminRouter(t, app)
	httpCtx := testutil.NewHTTPTestContext(t)
	httpCtx.Router = router

	ctx := context.Background()

	// Setup: Create admin, user, instance, inbound
	adminUser, _ := testutil.CreateTestAdminUser()
	app.Container.UserRepository.Create(ctx, adminUser)

	testUser, _ := testutil.CreateTestUserWithDefaults()
	app.Container.UserRepository.Create(ctx, testUser)

	// Create dependencies - server and node first
	testServer, _ := testutil.CreateTestServerWithDefaults()
	app.Container.ServerRepository.Create(ctx, testServer)

	testNode, _ := testutil.CreateTestNodeWithDefaults(testServer.ID)
	app.Container.NodeRepository.Create(ctx, testNode)

	instance, _ := testutil.CreateTestXrayInstanceWithDefaults(testNode.ID)
	app.Container.XrayInstanceRepository.Create(ctx, instance)

	inbound, _ := testutil.CreateTestInboundWithDefaults(instance.ID)
	app.Container.InboundRepository.Create(ctx, inbound)

	// Admin login
	loginReq := dto.LoginRequest{
		Email:    adminUser.Email.String(),
		Password: "Admin123!@#",
	}
	resp := httpCtx.POST("/api/v1/auth/login", loginReq, nil)
	require.Equal(t, 200, resp.Code)

	var loginResult response.Response
	httpCtx.GetResponseJSON(&loginResult)
	adminToken := loginResult.Data.(map[string]interface{})["access_token"].(string)

	// Step 1: Create Client
	t.Log("Step 1: Create Client")
	createReq := dto.AdminCreateClientRequest{
		InboundID: inbound.ID.String(),
		UserID:    testUser.ID.String(),
		Email:     "lifecycle-client@example.com",
	}
	resp = httpCtx.POST("/api/v1/admin/xray/clients", createReq, testutil.AuthHeader(adminToken))
	require.Equal(t, 201, resp.Code)

	var createResult response.Response
	httpCtx.GetResponseJSON(&createResult)
	clientData := createResult.Data.(map[string]interface{})
	clientInfo := clientData["client"].(map[string]interface{})
	clientID := clientInfo["id"].(string)

	// Step 2: Get Client Details
	t.Log("Step 2: Get Client Details")
	resp = httpCtx.GET("/api/v1/admin/xray/clients/"+clientID, testutil.AuthHeader(adminToken))
	require.Equal(t, 200, resp.Code)

	// Step 3: Disable Client
	t.Log("Step 3: Disable Client")
	resp = httpCtx.POST("/api/v1/admin/xray/clients/"+clientID+"/disable", nil, testutil.AuthHeader(adminToken))
	require.Equal(t, 200, resp.Code)

	// Step 4: Verify Client is Disabled
	t.Log("Step 4: Verify Client Disabled")
	resp = httpCtx.GET("/api/v1/admin/xray/clients/"+clientID, testutil.AuthHeader(adminToken))
	require.Equal(t, 200, resp.Code)

	var disabledResult response.Response
	httpCtx.GetResponseJSON(&disabledResult)
	disabledData := disabledResult.Data.(map[string]interface{})
	assert.False(t, disabledData["enabled"].(bool))

	// Step 5: Enable Client
	t.Log("Step 5: Enable Client")
	resp = httpCtx.POST("/api/v1/admin/xray/clients/"+clientID+"/enable", nil, testutil.AuthHeader(adminToken))
	require.Equal(t, 200, resp.Code)

	// Step 6: Verify Client is Enabled
	t.Log("Step 6: Verify Client Enabled")
	resp = httpCtx.GET("/api/v1/admin/xray/clients/"+clientID, testutil.AuthHeader(adminToken))
	require.Equal(t, 200, resp.Code)

	var enabledResult response.Response
	httpCtx.GetResponseJSON(&enabledResult)
	enabledData := enabledResult.Data.(map[string]interface{})
	assert.True(t, enabledData["enabled"].(bool))

	// Step 7: Delete Client
	t.Log("Step 7: Delete Client")
	resp = httpCtx.DELETE("/api/v1/admin/xray/clients/"+clientID, testutil.AuthHeader(adminToken))
	require.Equal(t, 200, resp.Code)

	// Step 8: Verify Client Deleted
	t.Log("Step 8: Verify Client Deleted")
	resp = httpCtx.GET("/api/v1/admin/xray/clients/"+clientID, testutil.AuthHeader(adminToken))
	require.Equal(t, 404, resp.Code)

	t.Log("✅ E2E Client Lifecycle Flow Complete")
}

// TestE2E_InboundLifecycleFlow tests complete inbound lifecycle
func TestE2E_InboundLifecycleFlow(t *testing.T) {
	if !testutil.IsIntegrationTest() {
		t.Skip("Skipping E2E integration test")
	}

	app := testutil.NewTestApp(t)
	defer app.Cleanup()
	defer app.CleanupTables()

	router := setupE2EAdminRouter(t, app)
	httpCtx := testutil.NewHTTPTestContext(t)
	httpCtx.Router = router

	ctx := context.Background()

	// Setup
	adminUser, _ := testutil.CreateTestAdminUser()
	app.Container.UserRepository.Create(ctx, adminUser)

	// Create dependencies - server and node first
	testServer, _ := testutil.CreateTestServerWithDefaults()
	app.Container.ServerRepository.Create(ctx, testServer)

	testNode, _ := testutil.CreateTestNodeWithDefaults(testServer.ID)
	app.Container.NodeRepository.Create(ctx, testNode)

	instance, _ := testutil.CreateTestXrayInstanceWithDefaults(testNode.ID)
	app.Container.XrayInstanceRepository.Create(ctx, instance)

	// Admin login
	loginReq := dto.LoginRequest{
		Email:    adminUser.Email.String(),
		Password: "Admin123!@#",
	}
	resp := httpCtx.POST("/api/v1/auth/login", loginReq, nil)
	require.Equal(t, 200, resp.Code)

	var loginResult response.Response
	httpCtx.GetResponseJSON(&loginResult)
	adminToken := loginResult.Data.(map[string]interface{})["access_token"].(string)

	// Step 1: Create Inbound
	t.Log("Step 1: Create Inbound")
	createReq := dto.AdminCreateInboundRequest{
		XrayInstanceID: instance.ID.String(),
		Protocol:       "vless",
		Port:           9443,
		Transport:      "tcp",
		Security:       "reality",
	}
	resp = httpCtx.POST("/api/v1/admin/xray/inbounds", createReq, testutil.AuthHeader(adminToken))
	require.Equal(t, 201, resp.Code)

	var createResult response.Response
	httpCtx.GetResponseJSON(&createResult)
	inboundData := createResult.Data.(map[string]interface{})
	inboundInfo := inboundData["inbound"].(map[string]interface{})
	inboundID := inboundInfo["id"].(string)

	// Step 2: Get Inbound Details
	t.Log("Step 2: Get Inbound Details")
	resp = httpCtx.GET("/api/v1/admin/xray/inbounds/"+inboundID, testutil.AuthHeader(adminToken))
	require.Equal(t, 200, resp.Code)

	// Step 3: List Inbounds
	t.Log("Step 3: List Inbounds")
	resp = httpCtx.GET("/api/v1/admin/xray/inbounds", testutil.AuthHeader(adminToken))
	require.Equal(t, 200, resp.Code)

	var listResult response.Response
	httpCtx.GetResponseJSON(&listResult)
	listData := listResult.Data.(map[string]interface{})
	total := int(listData["total"].(float64))
	assert.GreaterOrEqual(t, total, 1)

	// Step 4: Delete Inbound
	t.Log("Step 4: Delete Inbound")
	resp = httpCtx.DELETE("/api/v1/admin/xray/inbounds/"+inboundID, testutil.AuthHeader(adminToken))
	require.Equal(t, 200, resp.Code)

	// Step 5: Verify Inbound Deleted
	t.Log("Step 5: Verify Inbound Deleted")
	resp = httpCtx.GET("/api/v1/admin/xray/inbounds/"+inboundID, testutil.AuthHeader(adminToken))
	require.Equal(t, 404, resp.Code)

	t.Log("✅ E2E Inbound Lifecycle Flow Complete")
}
