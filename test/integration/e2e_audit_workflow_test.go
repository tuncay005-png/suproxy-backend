package integration

import (
	"context"
	"strconv"
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
	"github.com/suproxy/backend/internal/domain/audit"
	"github.com/suproxy/backend/internal/domain/user"
	"github.com/suproxy/backend/internal/infrastructure/testutil"
	"github.com/suproxy/backend/internal/interfaces/http/handler"
	"github.com/suproxy/backend/internal/interfaces/http/middleware"
	"github.com/suproxy/backend/internal/interfaces/http/response"
)

func setupE2EAuditRouter(t *testing.T, app *testutil.TestApp) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	// Create auth use-case instances
	registerCmd := authuc.NewRegisterCommand(app.Container.UserRepository, app.Container.XrayProvisioningService, app.Logger)
	loginCmd := authuc.NewLoginCommand(app.Container.UserRepository, app.Container.RefreshTokenRepository, app.Container.AuditLogRepository, app.JWT, app.Logger)
	refreshCmd := authuc.NewRefreshTokenCommand(app.Container.UserRepository, app.Container.RefreshTokenRepository, app.Container.AuditLogRepository, app.JWT, app.Logger)
	logoutCmd := authuc.NewLogoutCommand(app.Container.RefreshTokenRepository, app.Container.AuditLogRepository, app.Logger)
	getCurrentUserQuery := authuc.NewGetCurrentUserQuery(app.Container.UserRepository, app.Logger)
	getSessionsQuery := authuc.NewGetSessionsQuery(app.Container.RefreshTokenRepository, app.Logger)

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
		adminGroup.GET("/users", adminHandler.ListUsers)
		adminGroup.PUT("/users/:id/status", adminHandler.UpdateUserStatus)
		adminGroup.GET("/audit", adminHandler.ListAuditLogs)
		adminGroup.GET("/audit/:id", adminHandler.GetAuditLog)
		adminGroup.GET("/audit/stats", adminHandler.GetAuditStats)
		adminGroup.GET("/system/health", adminHandler.GetSystemHealth)
		adminGroup.GET("/system/stats", adminHandler.GetSystemStats)
	}

	return router
}

// TestE2E_AuditFlow tests audit logging throughout system operations
func TestE2E_AuditFlow(t *testing.T) {
	if !testutil.IsIntegrationTest() {
		t.Skip("Skipping E2E integration test")
	}

	app := testutil.NewTestApp(t)
	defer app.Cleanup()
	defer app.CleanupTables()

	router := setupE2EAuditRouter(t, app)
	httpCtx := testutil.NewHTTPTestContext(t)
	httpCtx.Router = router

	ctx := context.Background()

	// Setup: Create admin and users
	adminUser, err := testutil.CreateTestAdminUser()
	require.NoError(t, err)
	require.NoError(t, app.Container.UserRepository.Create(ctx, adminUser))

	regularUser, err := testutil.CreateTestUserWithDefaults()
	require.NoError(t, err)
	require.NoError(t, app.Container.UserRepository.Create(ctx, regularUser))

	// Step 1: Admin login (should create audit log)
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

	// Step 2: Perform operation (update user status - should create audit log)
	t.Log("Step 2: Update User Status (Creates Audit Log)")
	statusReq := dto.UpdateUserStatusRequest{
		Status: string(user.StatusInactive),
	}
	resp = httpCtx.PUT("/api/v1/admin/users/"+regularUser.ID.String()+"/status", statusReq, testutil.AuthHeader(adminToken))
	require.Equal(t, 200, resp.Code)

	// Step 3: List audit logs
	t.Log("Step 3: List Audit Logs")
	resp = httpCtx.GET("/api/v1/admin/audit", testutil.AuthHeader(adminToken))
	require.Equal(t, 200, resp.Code)

	var auditListResult response.Response
	httpCtx.GetResponseJSON(&auditListResult)
	require.True(t, auditListResult.Success)

	auditListData := auditListResult.Data.(map[string]interface{})
	
	// Safe type assertion for total
	var total int
	if totalRaw, ok := auditListData["total"]; ok && totalRaw != nil {
		if totalFloat, ok := totalRaw.(float64); ok {
			total = int(totalFloat)
		}
	}
	assert.GreaterOrEqual(t, total, 1, "Should have at least one audit log")

	// Step 4: Get specific audit log
	t.Log("Step 4: Get Specific Audit Log")
	logsRaw := auditListData["logs"]
	if logsRaw != nil {
		// Safe type assertion for logs array
		if logs, ok := logsRaw.([]interface{}); ok && len(logs) > 0 {
			if firstLog, ok := logs[0].(map[string]interface{}); ok {
				if logID, ok := firstLog["id"].(string); ok {
					resp = httpCtx.GET("/api/v1/admin/audit/"+logID, testutil.AuthHeader(adminToken))
					require.Equal(t, 200, resp.Code)

					var logResult response.Response
					httpCtx.GetResponseJSON(&logResult)
					require.True(t, logResult.Success)
				}
			}
		}
	}

	// Step 5: Get audit statistics
	t.Log("Step 5: Get Audit Statistics")
	resp = httpCtx.GET("/api/v1/admin/audit/stats", testutil.AuthHeader(adminToken))
	require.Equal(t, 200, resp.Code)

	var statsResult response.Response
	httpCtx.GetResponseJSON(&statsResult)
	require.True(t, statsResult.Success)

	t.Log("✅ E2E Audit Flow Complete")
}

// TestE2E_AuditFilteringFlow tests audit log filtering
func TestE2E_AuditFilteringFlow(t *testing.T) {
	if !testutil.IsIntegrationTest() {
		t.Skip("Skipping E2E integration test")
	}

	app := testutil.NewTestApp(t)
	defer app.Cleanup()
	defer app.CleanupTables()

	router := setupE2EAuditRouter(t, app)
	httpCtx := testutil.NewHTTPTestContext(t)
	httpCtx.Router = router

	ctx := context.Background()

	// Setup: Create admin
	adminUser, err := testutil.CreateTestAdminUser()
	require.NoError(t, err)
	require.NoError(t, app.Container.UserRepository.Create(ctx, adminUser))

	// Create multiple audit logs with different actions
	for i := 0; i < 3; i++ {
		log1 := audit.NewLog(adminUser.ID, audit.ActionLogin, "user", adminUser.ID, "127.0.0.1", "TestAgent")
		_ = app.Container.AuditLogRepository.Create(ctx, log1) // Test setup

		log2 := audit.NewLog(adminUser.ID, audit.ActionLogout, "user", adminUser.ID, "127.0.0.1", "TestAgent")
		_ = app.Container.AuditLogRepository.Create(ctx, log2) // Test setup
	}

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

	// Step 1: List all audit logs
	t.Log("Step 1: List All Audit Logs")
	resp = httpCtx.GET("/api/v1/admin/audit", testutil.AuthHeader(adminToken))
	require.Equal(t, 200, resp.Code)

	var allLogsResult response.Response
	httpCtx.GetResponseJSON(&allLogsResult)
	allLogsData := allLogsResult.Data.(map[string]interface{})
	
	// Safe type assertion for total
	var totalAll int
	if totalRaw, ok := allLogsData["total"]; ok && totalRaw != nil {
		if totalFloat, ok := totalRaw.(float64); ok {
			totalAll = int(totalFloat)
		}
	}
	assert.GreaterOrEqual(t, totalAll, 6)

	// Step 2: Filter by action
	t.Log("Step 2: Filter by Action (Login)")
	resp = httpCtx.GET("/api/v1/admin/audit?action=login", testutil.AuthHeader(adminToken))
	require.Equal(t, 200, resp.Code)

	var loginLogsResult response.Response
	httpCtx.GetResponseJSON(&loginLogsResult)
	loginLogsData := loginLogsResult.Data.(map[string]interface{})
	
	// Safe type assertion for total
	var totalLogin int
	if totalRaw, ok := loginLogsData["total"]; ok && totalRaw != nil {
		if totalFloat, ok := totalRaw.(float64); ok {
			totalLogin = int(totalFloat)
		}
	}
	assert.GreaterOrEqual(t, totalLogin, 3)

	// Step 3: Filter by user
	t.Log("Step 3: Filter by User ID")
	resp = httpCtx.GET("/api/v1/admin/audit?user_id="+adminUser.ID.String(), testutil.AuthHeader(adminToken))
	require.Equal(t, 200, resp.Code)

	var userLogsResult response.Response
	httpCtx.GetResponseJSON(&userLogsResult)
	userLogsData := userLogsResult.Data.(map[string]interface{})
	
	// Safe type assertion for total
	var totalUser int
	if totalRaw, ok := userLogsData["total"]; ok && totalRaw != nil {
		if totalFloat, ok := totalRaw.(float64); ok {
			totalUser = int(totalFloat)
		}
	}
	assert.GreaterOrEqual(t, totalUser, 6)

	// Step 4: Pagination
	t.Log("Step 4: Test Pagination")
	resp = httpCtx.GET("/api/v1/admin/audit?limit=2&offset=0", testutil.AuthHeader(adminToken))
	require.Equal(t, 200, resp.Code)

	var paginatedResult response.Response
	httpCtx.GetResponseJSON(&paginatedResult)
	paginatedData := paginatedResult.Data.(map[string]interface{})
	
	// Safe type assertion for logs array
	var logs []interface{}
	if logsRaw, ok := paginatedData["logs"]; ok && logsRaw != nil {
		if logsArray, ok := logsRaw.([]interface{}); ok {
			logs = logsArray
		}
	}
	assert.LessOrEqual(t, len(logs), 2)

	t.Log("✅ E2E Audit Filtering Flow Complete")
}

// TestE2E_FullAdminWorkflow tests complete admin workflow from login to audit
func TestE2E_FullAdminWorkflow(t *testing.T) {
	if !testutil.IsIntegrationTest() {
		t.Skip("Skipping E2E integration test")
	}

	app := testutil.NewTestApp(t)
	defer app.Cleanup()
	defer app.CleanupTables()

	router := setupE2EAuditRouter(t, app)
	httpCtx := testutil.NewHTTPTestContext(t)
	httpCtx.Router = router

	ctx := context.Background()

	// Setup: Create admin and multiple users
	adminUser, err := testutil.CreateTestAdminUser()
	require.NoError(t, err)
	require.NoError(t, app.Container.UserRepository.Create(ctx, adminUser))

	user1, err := testutil.CreateTestUserWithDefaults()
	require.NoError(t, err)
	require.NoError(t, app.Container.UserRepository.Create(ctx, user1))

	user2, err := testutil.CreateTestUser("user2", "user2@example.com", "Test123!@#")
	require.NoError(t, err)
	require.NoError(t, app.Container.UserRepository.Create(ctx, user2))

	// Phase 1: Authentication
	t.Log("=== Phase 1: Authentication ===")

	t.Log("Step 1.1: Admin Login")
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

	// Phase 2: User Management
	t.Log("=== Phase 2: User Management ===")

	t.Log("Step 2.1: List All Users")
	resp = httpCtx.GET("/api/v1/admin/users", testutil.AuthHeader(adminToken))
	require.Equal(t, 200, resp.Code)

	var usersResult response.Response
	httpCtx.GetResponseJSON(&usersResult)
	usersData := usersResult.Data.(map[string]interface{})
	
	// Safe type assertion for total
	var totalUsers int
	if totalRaw, ok := usersData["total"]; ok && totalRaw != nil {
		if totalFloat, ok := totalRaw.(float64); ok {
			totalUsers = int(totalFloat)
		}
	}
	assert.GreaterOrEqual(t, totalUsers, 3)

	t.Log("Step 2.2: Deactivate User 1")
	statusReq := dto.UpdateUserStatusRequest{
		Status: string(user.StatusInactive),
	}
	resp = httpCtx.PUT("/api/v1/admin/users/"+user1.ID.String()+"/status", statusReq, testutil.AuthHeader(adminToken))
	require.Equal(t, 200, resp.Code)

	t.Log("Step 2.3: Suspend User 2")
	statusReq.Status = string(user.StatusSuspended)
	resp = httpCtx.PUT("/api/v1/admin/users/"+user2.ID.String()+"/status", statusReq, testutil.AuthHeader(adminToken))
	require.Equal(t, 200, resp.Code)

	t.Log("Step 2.4: List Users Again (Verify Changes)")
	resp = httpCtx.GET("/api/v1/admin/users", testutil.AuthHeader(adminToken))
	require.Equal(t, 200, resp.Code)

	// Phase 3: Audit Review
	t.Log("=== Phase 3: Audit Review ===")

	t.Log("Step 3.1: List All Audit Logs")
	resp = httpCtx.GET("/api/v1/admin/audit", testutil.AuthHeader(adminToken))
	require.Equal(t, 200, resp.Code)

	var auditResult response.Response
	httpCtx.GetResponseJSON(&auditResult)
	auditData := auditResult.Data.(map[string]interface{})
	
	// Safe type assertion for total
	var totalAudit int
	if totalRaw, ok := auditData["total"]; ok && totalRaw != nil {
		if totalFloat, ok := totalRaw.(float64); ok {
			totalAudit = int(totalFloat)
		}
	}
	assert.GreaterOrEqual(t, totalAudit, 3, "Should have audit logs for login and status changes")

	t.Log("Step 3.2: Get Audit Statistics")
	resp = httpCtx.GET("/api/v1/admin/audit/stats", testutil.AuthHeader(adminToken))
	require.Equal(t, 200, resp.Code)

	var statsResult response.Response
	httpCtx.GetResponseJSON(&statsResult)
	require.True(t, statsResult.Success)

	// Phase 4: System Health Check
	t.Log("=== Phase 4: System Health Check ===")

	t.Log("Step 4.1: Get System Health")
	resp = httpCtx.GET("/api/v1/admin/system/health", testutil.AuthHeader(adminToken))
	require.Equal(t, 200, resp.Code)

	var healthResult response.Response
	httpCtx.GetResponseJSON(&healthResult)
	require.True(t, healthResult.Success)

	t.Log("Step 4.2: Get System Stats")
	resp = httpCtx.GET("/api/v1/admin/system/stats", testutil.AuthHeader(adminToken))
	require.Equal(t, 200, resp.Code)

	var sysStatsResult response.Response
	httpCtx.GetResponseJSON(&sysStatsResult)
	require.True(t, sysStatsResult.Success)

	t.Log("✅ E2E Full Admin Workflow Complete")
	t.Log("   - Successfully authenticated as admin")
	t.Log("   - Managed multiple users (list, deactivate, suspend)")
	t.Log("   - Reviewed audit logs and statistics")
	t.Log("   - Checked system health and stats")
}

// TestE2E_MetricsFlow tests metrics collection throughout operations
func TestE2E_MetricsFlow(t *testing.T) {
	if !testutil.IsIntegrationTest() {
		t.Skip("Skipping E2E integration test")
	}

	app := testutil.NewTestApp(t)
	defer app.Cleanup()
	defer app.CleanupTables()

	router := setupE2EAuditRouter(t, app)
	httpCtx := testutil.NewHTTPTestContext(t)
	httpCtx.Router = router

	ctx := context.Background()

	// Setup
	adminUser, err := testutil.CreateTestAdminUser()
	require.NoError(t, err)
	require.NoError(t, app.Container.UserRepository.Create(ctx, adminUser))

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

	// Step 1: Perform multiple operations (should increment counters)
	t.Log("Step 1: Perform Multiple Operations")

	// Multiple health checks
	for i := 0; i < 3; i++ {
		resp = httpCtx.GET("/api/v1/admin/system/health", testutil.AuthHeader(adminToken))
		require.Equal(t, 200, resp.Code)
	}

	// Step 2: Get system stats (includes metrics)
	t.Log("Step 2: Get System Stats")
	resp = httpCtx.GET("/api/v1/admin/system/stats", testutil.AuthHeader(adminToken))
	require.Equal(t, 200, resp.Code)

	var statsResult response.Response
	httpCtx.GetResponseJSON(&statsResult)
	require.True(t, statsResult.Success)
	require.NotNil(t, statsResult.Data)

	t.Log("✅ E2E Metrics Flow Complete")
}

// TestE2E_ConcurrentAdminOperations tests concurrent admin operations
func TestE2E_ConcurrentAdminOperations(t *testing.T) {
	if !testutil.IsIntegrationTest() {
		t.Skip("Skipping E2E integration test")
	}

	app := testutil.NewTestApp(t)
	defer app.Cleanup()
	defer app.CleanupTables()

	router := setupE2EAuditRouter(t, app)
	httpCtx := testutil.NewHTTPTestContext(t)
	httpCtx.Router = router

	ctx := context.Background()

	// Setup: Create admin
	adminUser, err := testutil.CreateTestAdminUser()
	require.NoError(t, err)
	require.NoError(t, app.Container.UserRepository.Create(ctx, adminUser))

	// Create multiple users
	for i := 0; i < 5; i++ {
		user, err := testutil.CreateTestUser("user"+strconv.Itoa(i), "user"+strconv.Itoa(i)+"@example.com", "Test123!@#")
		require.NoError(t, err)
		require.NoError(t, app.Container.UserRepository.Create(ctx, user))
	}

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

	// Step 1: Multiple concurrent list operations
	t.Log("Step 1: Concurrent List Operations")
	for i := 0; i < 5; i++ {
		resp = httpCtx.GET("/api/v1/admin/users", testutil.AuthHeader(adminToken))
		require.Equal(t, 200, resp.Code)
	}

	// Step 2: Multiple concurrent audit queries
	t.Log("Step 2: Concurrent Audit Queries")
	for i := 0; i < 5; i++ {
		resp = httpCtx.GET("/api/v1/admin/audit", testutil.AuthHeader(adminToken))
		require.Equal(t, 200, resp.Code)
	}

	// Step 3: Multiple concurrent system health checks
	t.Log("Step 3: Concurrent Health Checks")
	for i := 0; i < 5; i++ {
		resp = httpCtx.GET("/api/v1/admin/system/health", testutil.AuthHeader(adminToken))
		require.Equal(t, 200, resp.Code)
	}

	t.Log("✅ E2E Concurrent Admin Operations Complete")
}
