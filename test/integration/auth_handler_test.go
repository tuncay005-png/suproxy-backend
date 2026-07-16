package integration

import (
	"context"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/suproxy/backend/internal/application/dto"
	authuc "github.com/suproxy/backend/internal/application/usecase/auth"
	"github.com/suproxy/backend/internal/infrastructure/testutil"
	"github.com/suproxy/backend/internal/interfaces/http/handler"
	"github.com/suproxy/backend/internal/interfaces/http/middleware"
	"github.com/suproxy/backend/internal/interfaces/http/response"
)

func TestAuthHandler_Register(t *testing.T) {
	if !testutil.IsIntegrationTest() {
		t.Skip("Skipping integration test")
	}

	app := testutil.NewTestApp(t)
	defer app.Cleanup()
	defer app.CleanupTables()

	// Setup router
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
	router.POST("/api/v1/auth/register", authHandler.Register)

	t.Run("Success", func(t *testing.T) {
		defer app.CleanupTables()

		httpCtx := testutil.NewHTTPTestContext(t)
		httpCtx.Router = router

		req := dto.RegisterRequest{
			Email:    "newuser@example.com",
			Password: "Test123!@#",
		}

		resp := httpCtx.POST("/api/v1/auth/register", req, nil)

		assert.Equal(t, 201, resp.Code)

		var result response.Response
		httpCtx.GetResponseJSON(&result)
		assert.True(t, result.Success)
		assert.NotNil(t, result.Data)
	})

	t.Run("ValidationError_EmptyEmail", func(t *testing.T) {
		defer app.CleanupTables()

		httpCtx := testutil.NewHTTPTestContext(t)
		httpCtx.Router = router

		req := dto.RegisterRequest{
			Email:    "",
			Password: "Test123!@#",
		}

		resp := httpCtx.POST("/api/v1/auth/register", req, nil)
		assert.Equal(t, 422, resp.Code) // Validation errors return 422
	})

	t.Run("Conflict_EmailExists", func(t *testing.T) {
		defer app.CleanupTables()

		// Create existing user
		ctx := context.Background()
		existingUser, _ := testutil.CreateTestUserWithDefaults()
		_ = app.Container.UserRepository.Create(ctx, existingUser) // Test setup // Test setup: intentionally ignoring error

		httpCtx := testutil.NewHTTPTestContext(t)
		httpCtx.Router = router

		req := dto.RegisterRequest{
			Email:    existingUser.Email.String(),
			Password: "Test123!@#",
		}

		resp := httpCtx.POST("/api/v1/auth/register", req, nil)
		assert.Equal(t, 409, resp.Code)
	})
}

func TestAuthHandler_Login(t *testing.T) {
	if !testutil.IsIntegrationTest() {
		t.Skip("Skipping integration test")
	}

	app := testutil.NewTestApp(t)
	defer app.Cleanup()
	defer app.CleanupTables()

	// Setup router
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
	router.POST("/api/v1/auth/login", authHandler.Login)

	t.Run("Success", func(t *testing.T) {
		defer app.CleanupTables()

		ctx := context.Background()
		testUser, _ := testutil.CreateTestUserWithDefaults()
		_ = app.Container.UserRepository.Create(ctx, testUser) // Test setup

		httpCtx := testutil.NewHTTPTestContext(t)
		httpCtx.Router = router

		req := dto.LoginRequest{
			Email:    testUser.Email.String(),
			Password: "Test123!@#", // Default password
		}

		resp := httpCtx.POST("/api/v1/auth/login", req, nil)

		assert.Equal(t, 200, resp.Code)

		var result response.Response
		httpCtx.GetResponseJSON(&result)
		assert.True(t, result.Success)
	})

	t.Run("Unauthorized_InvalidCredentials", func(t *testing.T) {
		defer app.CleanupTables()

		httpCtx := testutil.NewHTTPTestContext(t)
		httpCtx.Router = router

		req := dto.LoginRequest{
			Email:    "nonexistent@example.com",
			Password: "WrongPassword123!",
		}

		resp := httpCtx.POST("/api/v1/auth/login", req, nil)
		assert.Equal(t, 401, resp.Code)
	})

	t.Run("Unauthorized_WrongPassword", func(t *testing.T) {
		defer app.CleanupTables()

		ctx := context.Background()
		testUser, _ := testutil.CreateTestUserWithDefaults()
		_ = app.Container.UserRepository.Create(ctx, testUser) // Test setup

		httpCtx := testutil.NewHTTPTestContext(t)
		httpCtx.Router = router

		req := dto.LoginRequest{
			Email:    testUser.Email.String(),
			Password: "WrongPassword123!",
		}

		resp := httpCtx.POST("/api/v1/auth/login", req, nil)
		assert.Equal(t, 401, resp.Code)
	})

	t.Run("Forbidden_InactiveUser", func(t *testing.T) {
		defer app.CleanupTables()

		ctx := context.Background()
		testUser, _ := testutil.CreateTestUserWithDefaults()
		_ = testUser.Deactivate()                              // Test setup
		_ = app.Container.UserRepository.Create(ctx, testUser) // Test setup

		httpCtx := testutil.NewHTTPTestContext(t)
		httpCtx.Router = router

		req := dto.LoginRequest{
			Email:    testUser.Email.String(),
			Password: "Test123!@#",
		}

		resp := httpCtx.POST("/api/v1/auth/login", req, nil)
		assert.Equal(t, 403, resp.Code)
	})
}

func TestAuthHandler_RefreshToken(t *testing.T) {
	if !testutil.IsIntegrationTest() {
		t.Skip("Skipping integration test")
	}

	app := testutil.NewTestApp(t)
	defer app.Cleanup()
	defer app.CleanupTables()

	// Setup router
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
	router.POST("/api/v1/auth/refresh", authHandler.RefreshToken)

	t.Run("Success", func(t *testing.T) {
		defer app.CleanupTables()

		ctx := context.Background()
		authHelper := testutil.NewAuthHelper(app.JWT, t)
		testUser, accessToken, refreshToken := authHelper.CreateAuthenticatedUser(app.Container.UserRepository)

		// Store refresh token
		testutil.CreateTestRefreshToken(ctx, t, app.Container.RefreshTokenRepository, testUser.ID, refreshToken)

		httpCtx := testutil.NewHTTPTestContext(t)
		httpCtx.Router = router

		req := dto.RefreshTokenRequest{
			RefreshToken: refreshToken,
		}

		resp := httpCtx.POST("/api/v1/auth/refresh", req, nil)

		assert.Equal(t, 200, resp.Code)

		var result response.Response
		httpCtx.GetResponseJSON(&result)
		assert.True(t, result.Success)

		_ = accessToken // Used for validation
	})

	t.Run("Unauthorized_InvalidToken", func(t *testing.T) {
		defer app.CleanupTables()

		httpCtx := testutil.NewHTTPTestContext(t)
		httpCtx.Router = router

		req := dto.RefreshTokenRequest{
			RefreshToken: "invalid.token.here",
		}

		resp := httpCtx.POST("/api/v1/auth/refresh", req, nil)
		assert.Equal(t, 401, resp.Code)
	})
}

func TestAuthHandler_GetCurrentUser(t *testing.T) {
	if !testutil.IsIntegrationTest() {
		t.Skip("Skipping integration test")
	}

	app := testutil.NewTestApp(t)
	defer app.Cleanup()
	defer app.CleanupTables()

	// Setup router
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

	// Add auth middleware
	router.GET("/api/v1/auth/me", middleware.AuthMiddleware(app.JWT), authHandler.GetCurrentUser)

	t.Run("Success", func(t *testing.T) {
		defer app.CleanupTables()

		authHelper := testutil.NewAuthHelper(app.JWT, t)
		testUser, accessToken, _ := authHelper.CreateAuthenticatedUser(app.Container.UserRepository)

		httpCtx := testutil.NewHTTPTestContext(t)
		httpCtx.Router = router

		resp := httpCtx.GET("/api/v1/auth/me", testutil.AuthHeader(accessToken))

		assert.Equal(t, 200, resp.Code)

		var result response.Response
		httpCtx.GetResponseJSON(&result)
		assert.True(t, result.Success)

		_ = testUser
	})

	t.Run("Unauthorized_NoToken", func(t *testing.T) {
		defer app.CleanupTables()

		httpCtx := testutil.NewHTTPTestContext(t)
		httpCtx.Router = router

		resp := httpCtx.GET("/api/v1/auth/me", nil)
		assert.Equal(t, 401, resp.Code)
	})

	t.Run("Unauthorized_InvalidToken", func(t *testing.T) {
		defer app.CleanupTables()

		httpCtx := testutil.NewHTTPTestContext(t)
		httpCtx.Router = router

		resp := httpCtx.GET("/api/v1/auth/me", testutil.AuthHeader("invalid.token"))
		assert.Equal(t, 401, resp.Code)
	})
}

func TestAuthHandler_GetSessions(t *testing.T) {
	if !testutil.IsIntegrationTest() {
		t.Skip("Skipping integration test")
	}

	app := testutil.NewTestApp(t)
	defer app.Cleanup()
	defer app.CleanupTables()

	// Setup router
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

	router.GET("/api/v1/auth/sessions", middleware.AuthMiddleware(app.JWT), authHandler.GetSessions)

	t.Run("Success", func(t *testing.T) {
		defer app.CleanupTables()

		authHelper := testutil.NewAuthHelper(app.JWT, t)
		testUser, accessToken, _ := authHelper.CreateAuthenticatedUser(app.Container.UserRepository)

		httpCtx := testutil.NewHTTPTestContext(t)
		httpCtx.Router = router

		resp := httpCtx.GET("/api/v1/auth/sessions", testutil.AuthHeader(accessToken))

		assert.Equal(t, 200, resp.Code)

		var result response.Response
		httpCtx.GetResponseJSON(&result)
		assert.True(t, result.Success)

		_ = testUser
	})

	t.Run("Unauthorized_NoToken", func(t *testing.T) {
		defer app.CleanupTables()

		httpCtx := testutil.NewHTTPTestContext(t)
		httpCtx.Router = router

		resp := httpCtx.GET("/api/v1/auth/sessions", nil)
		assert.Equal(t, 401, resp.Code)
	})
}

func TestAuthHandler_LogoutAll(t *testing.T) {
	if !testutil.IsIntegrationTest() {
		t.Skip("Skipping integration test")
	}

	app := testutil.NewTestApp(t)
	defer app.Cleanup()
	defer app.CleanupTables()

	// Setup router
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

	router.POST("/api/v1/auth/logout-all", middleware.AuthMiddleware(app.JWT), authHandler.LogoutAll)

	t.Run("Success", func(t *testing.T) {
		defer app.CleanupTables()

		authHelper := testutil.NewAuthHelper(app.JWT, t)
		testUser, accessToken, _ := authHelper.CreateAuthenticatedUser(app.Container.UserRepository)

		httpCtx := testutil.NewHTTPTestContext(t)
		httpCtx.Router = router

		resp := httpCtx.POST("/api/v1/auth/logout-all", nil, testutil.AuthHeader(accessToken))

		assert.Equal(t, 204, resp.Code)

		_ = testUser
	})

	t.Run("Unauthorized_NoToken", func(t *testing.T) {
		defer app.CleanupTables()

		httpCtx := testutil.NewHTTPTestContext(t)
		httpCtx.Router = router

		resp := httpCtx.POST("/api/v1/auth/logout-all", nil, nil)
		assert.Equal(t, 401, resp.Code)
	})
}

func TestAuthHandler_LogoutSingle(t *testing.T) {
	if !testutil.IsIntegrationTest() {
		t.Skip("Skipping integration test")
	}

	app := testutil.NewTestApp(t)
	defer app.Cleanup()
	defer app.CleanupTables()

	// Setup router
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

	router.DELETE("/api/v1/auth/sessions/:id", middleware.AuthMiddleware(app.JWT), authHandler.LogoutSingle)

	t.Run("Success", func(t *testing.T) {
		defer app.CleanupTables()

		ctx := context.Background()
		authHelper := testutil.NewAuthHelper(app.JWT, t)
		testUser, accessToken, refreshToken := authHelper.CreateAuthenticatedUser(app.Container.UserRepository)

		// Create refresh token
		tokenEntity := testutil.CreateTestRefreshToken(ctx, t, app.Container.RefreshTokenRepository, testUser.ID, refreshToken)

		httpCtx := testutil.NewHTTPTestContext(t)
		httpCtx.Router = router

		resp := httpCtx.DELETE("/api/v1/auth/sessions/"+tokenEntity.ID.String(), testutil.AuthHeader(accessToken))

		assert.Equal(t, 204, resp.Code)
	})

	t.Run("BadRequest_InvalidTokenID", func(t *testing.T) {
		defer app.CleanupTables()

		authHelper := testutil.NewAuthHelper(app.JWT, t)
		_, accessToken, _ := authHelper.CreateAuthenticatedUser(app.Container.UserRepository)

		httpCtx := testutil.NewHTTPTestContext(t)
		httpCtx.Router = router

		resp := httpCtx.DELETE("/api/v1/auth/sessions/invalid-uuid", testutil.AuthHeader(accessToken))
		assert.Equal(t, 400, resp.Code)
	})
}
