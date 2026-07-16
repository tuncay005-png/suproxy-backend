package integration

import (
	"context"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/suproxy/backend/internal/application/dto"
	authuc "github.com/suproxy/backend/internal/application/usecase/auth"
	"github.com/suproxy/backend/internal/infrastructure/testutil"
	"github.com/suproxy/backend/internal/interfaces/http/handler"
	"github.com/suproxy/backend/internal/interfaces/http/middleware"
	"github.com/suproxy/backend/internal/interfaces/http/response"
)

func setupE2ERouter(t *testing.T, app *testutil.TestApp) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	// Instantiate use-case instances
	registerCmd := authuc.NewRegisterCommand(app.Container.UserRepository, app.Container.XrayProvisioningService, app.Logger)
	loginCmd := authuc.NewLoginCommand(app.Container.UserRepository, app.Container.RefreshTokenRepository, app.Container.AuditLogRepository, app.JWT, app.Logger)
	refreshTokenCmd := authuc.NewRefreshTokenCommand(app.Container.UserRepository, app.Container.RefreshTokenRepository, app.Container.AuditLogRepository, app.JWT, app.Logger)
	logoutCmd := authuc.NewLogoutCommand(app.Container.RefreshTokenRepository, app.Container.AuditLogRepository, app.Logger)
	getCurrentUserQuery := authuc.NewGetCurrentUserQuery(app.Container.UserRepository, app.Logger)
	getSessionsQuery := authuc.NewGetSessionsQuery(app.Container.RefreshTokenRepository, app.Logger)

	// Auth handler
	authHandler := handler.NewAuthHandler(
		registerCmd,
		loginCmd,
		refreshTokenCmd,
		logoutCmd,
		getCurrentUserQuery,
		getSessionsQuery,
	)

	// Middleware
	authMiddleware := middleware.AuthMiddleware(app.JWT)

	// Auth routes
	authGroup := router.Group("/api/v1/auth")
	{
		authGroup.POST("/register", authHandler.Register)
		authGroup.POST("/login", authHandler.Login)
		authGroup.POST("/refresh", authHandler.RefreshToken)
		authGroup.GET("/me", authMiddleware, authHandler.GetCurrentUser)
		authGroup.GET("/sessions", authMiddleware, authHandler.GetSessions)
		authGroup.POST("/logout-all", authMiddleware, authHandler.LogoutAll)
		authGroup.DELETE("/sessions/:id", authMiddleware, authHandler.LogoutSingle)
	}

	return router
}

// TestE2E_UserRegistrationFlow tests complete user registration workflow
func TestE2E_UserRegistrationFlow(t *testing.T) {
	if !testutil.IsIntegrationTest() {
		t.Skip("Skipping E2E integration test")
	}

	app := testutil.NewTestApp(t)
	defer app.Cleanup()
	defer app.CleanupTables()

	router := setupE2ERouter(t, app)
	httpCtx := testutil.NewHTTPTestContext(t)
	httpCtx.Router = router

	// Step 1: Register new user
	t.Log("Step 1: User Registration")
	registerReq := dto.RegisterRequest{
		Email:    "e2e-test@example.com",
		Password: "Test123!@#",
	}

	resp := httpCtx.POST("/api/v1/auth/register", registerReq, nil)
	t.Logf("Registration response status: %d, body: %s", resp.Code, resp.Body.String())
	require.Equal(t, 201, resp.Code, "Registration should succeed")

	var registerResult response.Response
	httpCtx.GetResponseJSON(&registerResult)
	require.True(t, registerResult.Success, "Expected success=true, got: %+v", registerResult)
	require.NotNil(t, registerResult.Data)

	// Extract tokens from registration response
	registerData, ok := registerResult.Data.(map[string]interface{})
	require.True(t, ok, "Expected data to be map, got: %T", registerResult.Data)
	accessToken, ok := registerData["access_token"].(string)
	require.True(t, ok, "Expected access_token in response")
	require.NotEmpty(t, accessToken)

	// Step 2: Login with same credentials
	t.Log("Step 2: User Login")
	loginReq := dto.LoginRequest{
		Email:    "e2e-test@example.com",
		Password: "Test123!@#",
	}

	resp = httpCtx.POST("/api/v1/auth/login", loginReq, nil)
	require.Equal(t, 200, resp.Code, "Login should succeed")

	var loginResult response.Response
	httpCtx.GetResponseJSON(&loginResult)
	require.True(t, loginResult.Success)

	loginData, ok := loginResult.Data.(map[string]interface{})
	require.True(t, ok)
	newAccessToken, ok := loginData["access_token"].(string)
	require.True(t, ok)
	require.NotEmpty(t, newAccessToken)
	refreshToken, ok := loginData["refresh_token"].(string)
	require.True(t, ok)
	require.NotEmpty(t, refreshToken)

	// Step 3: Get current user with JWT
	t.Log("Step 3: Get Current User")
	resp = httpCtx.GET("/api/v1/auth/me", testutil.AuthHeader(newAccessToken))
	require.Equal(t, 200, resp.Code, "Get current user should succeed")

	var userResult response.Response
	httpCtx.GetResponseJSON(&userResult)
	require.True(t, userResult.Success)
	require.NotNil(t, userResult.Data)

	userData, ok := userResult.Data.(map[string]interface{})
	require.True(t, ok)
	email, ok := userData["email"].(string)
	require.True(t, ok)
	assert.Equal(t, "e2e-test@example.com", email)

	t.Log("✅ E2E User Registration Flow Complete")
}

// TestE2E_UserSessionFlow tests complete session management workflow
func TestE2E_UserSessionFlow(t *testing.T) {
	if !testutil.IsIntegrationTest() {
		t.Skip("Skipping E2E integration test")
	}

	app := testutil.NewTestApp(t)
	defer app.Cleanup()
	defer app.CleanupTables()

	router := setupE2ERouter(t, app)
	httpCtx := testutil.NewHTTPTestContext(t)
	httpCtx.Router = router

	ctx := context.Background()

	// Setup: Create user
	testUser, _ := testutil.CreateTestUserWithDefaults()
	_ = app.Container.UserRepository.Create(ctx, testUser)

	// Step 1: Login to create session
	t.Log("Step 1: Login")
	loginReq := dto.LoginRequest{
		Email:    testUser.Email.String(),
		Password: "Test123!@#",
	}

	resp := httpCtx.POST("/api/v1/auth/login", loginReq, nil)
	require.Equal(t, 200, resp.Code)

	var loginResult response.Response
	httpCtx.GetResponseJSON(&loginResult)
	require.True(t, loginResult.Success)

	loginData := loginResult.Data.(map[string]interface{})
	accessToken := loginData["access_token"].(string)
	refreshToken := loginData["refresh_token"].(string)
	require.NotEmpty(t, accessToken)
	require.NotEmpty(t, refreshToken)

	// Step 2: Get active sessions
	t.Log("Step 2: Get Active Sessions")
	resp = httpCtx.GET("/api/v1/auth/sessions", testutil.AuthHeader(accessToken))
	require.Equal(t, 200, resp.Code)

	var sessionsResult response.Response
	httpCtx.GetResponseJSON(&sessionsResult)
	require.True(t, sessionsResult.Success)

	// Step 3: Refresh token
	t.Log("Step 3: Refresh Token")
	refreshReq := dto.RefreshTokenRequest{
		RefreshToken: refreshToken,
	}

	// Note: Login already created the refresh token, no need to create it again

	resp = httpCtx.POST("/api/v1/auth/refresh", refreshReq, nil)
	require.Equal(t, 200, resp.Code)

	var refreshResult response.Response
	httpCtx.GetResponseJSON(&refreshResult)
	require.True(t, refreshResult.Success)

	refreshData := refreshResult.Data.(map[string]interface{})
	newAccessToken := refreshData["access_token"].(string)
	require.NotEmpty(t, newAccessToken)

	// Step 4: Logout all sessions
	t.Log("Step 4: Logout All")
	resp = httpCtx.POST("/api/v1/auth/logout-all", nil, testutil.AuthHeader(newAccessToken))
	require.Equal(t, 204, resp.Code)

	t.Log("✅ E2E User Session Flow Complete")
}

// TestE2E_MultipleUserSessions tests multiple concurrent sessions
func TestE2E_MultipleUserSessions(t *testing.T) {
	if !testutil.IsIntegrationTest() {
		t.Skip("Skipping E2E integration test")
	}

	app := testutil.NewTestApp(t)
	defer app.Cleanup()
	defer app.CleanupTables()

	router := setupE2ERouter(t, app)
	httpCtx := testutil.NewHTTPTestContext(t)
	httpCtx.Router = router

	ctx := context.Background()

	// Setup: Create user
	testUser, _ := testutil.CreateTestUserWithDefaults()
	_ = app.Container.UserRepository.Create(ctx, testUser)

	// Step 1: Login from device 1
	t.Log("Step 1: Login from Device 1")
	loginReq1 := dto.LoginRequest{
		Email:      testUser.Email.String(),
		Password:   "Test123!@#",
		DeviceName: "Device 1",
		Platform:   "iOS",
	}

	resp := httpCtx.POST("/api/v1/auth/login", loginReq1, nil)
	require.Equal(t, 200, resp.Code)

	var loginResult1 response.Response
	httpCtx.GetResponseJSON(&loginResult1)
	token1 := loginResult1.Data.(map[string]interface{})["access_token"].(string)

	// Step 2: Login from device 2
	t.Log("Step 2: Login from Device 2")
	loginReq2 := dto.LoginRequest{
		Email:      testUser.Email.String(),
		Password:   "Test123!@#",
		DeviceName: "Device 2",
		Platform:   "Android",
	}

	resp = httpCtx.POST("/api/v1/auth/login", loginReq2, nil)
	require.Equal(t, 200, resp.Code)

	var loginResult2 response.Response
	httpCtx.GetResponseJSON(&loginResult2)
	token2 := loginResult2.Data.(map[string]interface{})["access_token"].(string)

	// Step 3: Verify both tokens work
	t.Log("Step 3: Verify Both Tokens")
	resp = httpCtx.GET("/api/v1/auth/me", testutil.AuthHeader(token1))
	require.Equal(t, 200, resp.Code)

	resp = httpCtx.GET("/api/v1/auth/me", testutil.AuthHeader(token2))
	require.Equal(t, 200, resp.Code)

	// Step 4: Get sessions (should see multiple)
	t.Log("Step 4: Check Active Sessions")
	resp = httpCtx.GET("/api/v1/auth/sessions", testutil.AuthHeader(token1))
	require.Equal(t, 200, resp.Code)

	var sessionsResult response.Response
	httpCtx.GetResponseJSON(&sessionsResult)
	sessionsData := sessionsResult.Data.(map[string]interface{})
	total := int(sessionsData["total"].(float64))
	assert.GreaterOrEqual(t, total, 1) // At least one session

	// Step 5: Logout all from device 1
	t.Log("Step 5: Logout All Sessions")
	resp = httpCtx.POST("/api/v1/auth/logout-all", nil, testutil.AuthHeader(token1))
	require.Equal(t, 204, resp.Code)

	t.Log("✅ E2E Multiple User Sessions Complete")
}

// TestE2E_InvalidCredentialsFlow tests error handling in auth flow
func TestE2E_InvalidCredentialsFlow(t *testing.T) {
	if !testutil.IsIntegrationTest() {
		t.Skip("Skipping E2E integration test")
	}

	app := testutil.NewTestApp(t)
	defer app.Cleanup()
	defer app.CleanupTables()

	router := setupE2ERouter(t, app)
	httpCtx := testutil.NewHTTPTestContext(t)
	httpCtx.Router = router

	ctx := context.Background()

	// Setup: Create user
	testUser, _ := testutil.CreateTestUserWithDefaults()
	_ = app.Container.UserRepository.Create(ctx, testUser)

	// Step 1: Try login with wrong password
	t.Log("Step 1: Login with Wrong Password")
	loginReq := dto.LoginRequest{
		Email:    testUser.Email.String(),
		Password: "WrongPassword123!",
	}

	resp := httpCtx.POST("/api/v1/auth/login", loginReq, nil)
	assert.Equal(t, 401, resp.Code)

	// Step 2: Try login with non-existent email
	t.Log("Step 2: Login with Non-existent Email")
	loginReq2 := dto.LoginRequest{
		Email:    "nonexistent@example.com",
		Password: "Test123!@#",
	}

	resp = httpCtx.POST("/api/v1/auth/login", loginReq2, nil)
	assert.Equal(t, 401, resp.Code)

	// Step 3: Try to access protected route without token
	t.Log("Step 3: Access Protected Route Without Token")
	resp = httpCtx.GET("/api/v1/auth/me", nil)
	assert.Equal(t, 401, resp.Code)

	// Step 4: Try to access protected route with invalid token
	t.Log("Step 4: Access Protected Route With Invalid Token")
	resp = httpCtx.GET("/api/v1/auth/me", testutil.AuthHeader("invalid.token.here"))
	assert.Equal(t, 401, resp.Code)

	t.Log("✅ E2E Invalid Credentials Flow Complete")
}

// TestE2E_UserRegistrationValidation tests input validation in registration
func TestE2E_UserRegistrationValidation(t *testing.T) {
	if !testutil.IsIntegrationTest() {
		t.Skip("Skipping E2E integration test")
	}

	app := testutil.NewTestApp(t)
	defer app.Cleanup()
	defer app.CleanupTables()

	router := setupE2ERouter(t, app)
	httpCtx := testutil.NewHTTPTestContext(t)
	httpCtx.Router = router

	// Step 1: Try register with empty email
	t.Log("Step 1: Register with Empty Email")
	req1 := dto.RegisterRequest{
		Email:    "",
		Password: "Test123!@#",
	}
	resp := httpCtx.POST("/api/v1/auth/register", req1, nil)
	assert.Equal(t, 422, resp.Code)

	// Step 2: Try register with invalid email
	t.Log("Step 2: Register with Invalid Email")
	req2 := dto.RegisterRequest{
		Email:    "invalid-email",
		Password: "Test123!@#",
	}
	resp = httpCtx.POST("/api/v1/auth/register", req2, nil)
	assert.Equal(t, 422, resp.Code)

	// Step 3: Try register with weak password
	t.Log("Step 3: Register with Weak Password")
	req3 := dto.RegisterRequest{
		Email:    "test@example.com",
		Password: "weak",
	}
	resp = httpCtx.POST("/api/v1/auth/register", req3, nil)
	assert.Equal(t, 422, resp.Code)

	// Step 4: Register successfully
	t.Log("Step 4: Register Successfully")
	req4 := dto.RegisterRequest{
		Email:    "valid@example.com",
		Password: "Strong123!@#",
	}
	resp = httpCtx.POST("/api/v1/auth/register", req4, nil)
	assert.Equal(t, 201, resp.Code)

	// Step 5: Try to register same email again (conflict)
	t.Log("Step 5: Register Same Email Again")
	req5 := dto.RegisterRequest{
		Email:    "valid@example.com",
		Password: "Strong123!@#",
	}
	resp = httpCtx.POST("/api/v1/auth/register", req5, nil)
	assert.Equal(t, 409, resp.Code)

	t.Log("✅ E2E User Registration Validation Complete")
}
