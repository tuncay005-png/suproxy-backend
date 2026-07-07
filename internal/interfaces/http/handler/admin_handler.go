package handler

import (
	"errors"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/suproxy/backend/internal/application/dto"
	"github.com/suproxy/backend/internal/application/mapper"
	adminuser "github.com/suproxy/backend/internal/application/usecase/admin/user"
	adminxray "github.com/suproxy/backend/internal/application/usecase/admin/xray_instance"
	admininbound "github.com/suproxy/backend/internal/application/usecase/admin/inbound"
	adminclient "github.com/suproxy/backend/internal/application/usecase/admin/client"
	adminaudit "github.com/suproxy/backend/internal/application/usecase/admin/audit"
	adminsystem "github.com/suproxy/backend/internal/application/usecase/admin/system"
	"github.com/suproxy/backend/internal/domain/user"
	"github.com/suproxy/backend/internal/domain/xray"
	"github.com/suproxy/backend/internal/infrastructure/logger"
	"github.com/suproxy/backend/internal/interfaces/http/middleware"
	"github.com/suproxy/backend/internal/interfaces/http/response"
)

// AdminHandler handles admin-only operations
type AdminHandler struct {
	logger *logger.Logger

	// User Management UseCases
	listUsersQuery          *adminuser.ListUsersQuery
	getUserQuery            *adminuser.GetUserQuery
	updateUserStatusCommand *adminuser.UpdateUserStatusCommand
	updateUserRoleCommand   *adminuser.UpdateUserRoleCommand

	// Xray Instance Management UseCases (Phase 17.3)
	listInstancesQuery         *adminxray.ListInstancesQuery
	getInstanceQuery           *adminxray.GetInstanceQuery
	getInstanceStatsQuery      *adminxray.GetInstanceStatsQuery
	startInstanceCommand       *adminxray.StartInstanceCommand
	stopInstanceCommand        *adminxray.StopInstanceCommand
	restartInstanceCommand     *adminxray.RestartInstanceCommand
	reloadInstanceCommand      *adminxray.ReloadInstanceCommand
	checkInstanceHealthCommand *adminxray.CheckInstanceHealthCommand

	// Inbound Management UseCases (Phase 17.4)
	listInboundsQuery      *admininbound.ListInboundsQuery
	getInboundQuery        *admininbound.GetInboundQuery
	createInboundCommand   *admininbound.CreateInboundCommand
	updateInboundCommand   *admininbound.UpdateInboundCommand
	deleteInboundCommand   *admininbound.DeleteInboundCommand
	enableInboundCommand   *admininbound.EnableInboundCommand
	disableInboundCommand  *admininbound.DisableInboundCommand

	// Client Management UseCases (Phase 17.5)
	listClientsQuery            *adminclient.ListClientsQuery
	getClientQuery              *adminclient.GetClientQuery
	createClientCommand         *adminclient.CreateClientCommand
	deleteClientCommand         *adminclient.DeleteClientCommand
	enableClientCommand         *adminclient.EnableClientCommand
	disableClientCommand        *adminclient.DisableClientCommand
	regenerateClientUUIDCommand *adminclient.RegenerateClientUUIDCommand
	reprovisionClientCommand    *adminclient.ReprovisionClientCommand

	// Audit Log Management UseCases (Phase 17.6)
	listAuditLogsQuery *adminaudit.ListAuditLogsQuery
	getAuditLogQuery   *adminaudit.GetAuditLogQuery
	getAuditStatsQuery *adminaudit.GetAuditStatsQuery

	// System Admin UseCases (Phase 17.7)
	getSystemHealthQuery      *adminsystem.GetSystemHealthQuery
	getSystemStatsQuery       *adminsystem.GetSystemStatsQuery
	getVersionQuery           *adminsystem.GetVersionQuery
	getDatabaseStatusQuery    *adminsystem.GetDatabaseStatusQuery
	getXraySystemStatusQuery  *adminsystem.GetXraySystemStatusQuery

	// Future: Add more UseCases here as they are implemented
}

// NewAdminHandler creates a new admin handler
func NewAdminHandler(
	log *logger.Logger,
	listUsersQuery *adminuser.ListUsersQuery,
	getUserQuery *adminuser.GetUserQuery,
	updateUserStatusCommand *adminuser.UpdateUserStatusCommand,
	updateUserRoleCommand *adminuser.UpdateUserRoleCommand,
	listInstancesQuery *adminxray.ListInstancesQuery,
	getInstanceQuery *adminxray.GetInstanceQuery,
	getInstanceStatsQuery *adminxray.GetInstanceStatsQuery,
	startInstanceCommand *adminxray.StartInstanceCommand,
	stopInstanceCommand *adminxray.StopInstanceCommand,
	restartInstanceCommand *adminxray.RestartInstanceCommand,
	reloadInstanceCommand *adminxray.ReloadInstanceCommand,
	checkInstanceHealthCommand *adminxray.CheckInstanceHealthCommand,
	listInboundsQuery *admininbound.ListInboundsQuery,
	getInboundQuery *admininbound.GetInboundQuery,
	createInboundCommand *admininbound.CreateInboundCommand,
	updateInboundCommand *admininbound.UpdateInboundCommand,
	deleteInboundCommand *admininbound.DeleteInboundCommand,
	enableInboundCommand *admininbound.EnableInboundCommand,
	disableInboundCommand *admininbound.DisableInboundCommand,
	listClientsQuery *adminclient.ListClientsQuery,
	getClientQuery *adminclient.GetClientQuery,
	createClientCommand *adminclient.CreateClientCommand,
	deleteClientCommand *adminclient.DeleteClientCommand,
	enableClientCommand *adminclient.EnableClientCommand,
	disableClientCommand *adminclient.DisableClientCommand,
	regenerateClientUUIDCommand *adminclient.RegenerateClientUUIDCommand,
	reprovisionClientCommand *adminclient.ReprovisionClientCommand,
	listAuditLogsQuery *adminaudit.ListAuditLogsQuery,
	getAuditLogQuery *adminaudit.GetAuditLogQuery,
	getAuditStatsQuery *adminaudit.GetAuditStatsQuery,
	getSystemHealthQuery *adminsystem.GetSystemHealthQuery,
	getSystemStatsQuery *adminsystem.GetSystemStatsQuery,
	getVersionQuery *adminsystem.GetVersionQuery,
	getDatabaseStatusQuery *adminsystem.GetDatabaseStatusQuery,
	getXraySystemStatusQuery *adminsystem.GetXraySystemStatusQuery,
) *AdminHandler {
	return &AdminHandler{
		logger:                      log,
		listUsersQuery:              listUsersQuery,
		getUserQuery:                getUserQuery,
		updateUserStatusCommand:     updateUserStatusCommand,
		updateUserRoleCommand:       updateUserRoleCommand,
		listInstancesQuery:          listInstancesQuery,
		getInstanceQuery:            getInstanceQuery,
		getInstanceStatsQuery:       getInstanceStatsQuery,
		startInstanceCommand:        startInstanceCommand,
		stopInstanceCommand:         stopInstanceCommand,
		restartInstanceCommand:      restartInstanceCommand,
		reloadInstanceCommand:       reloadInstanceCommand,
		checkInstanceHealthCommand:  checkInstanceHealthCommand,
		listInboundsQuery:           listInboundsQuery,
		getInboundQuery:             getInboundQuery,
		createInboundCommand:        createInboundCommand,
		updateInboundCommand:        updateInboundCommand,
		deleteInboundCommand:        deleteInboundCommand,
		enableInboundCommand:        enableInboundCommand,
		disableInboundCommand:       disableInboundCommand,
		listClientsQuery:            listClientsQuery,
		getClientQuery:              getClientQuery,
		createClientCommand:         createClientCommand,
		deleteClientCommand:         deleteClientCommand,
		enableClientCommand:         enableClientCommand,
		disableClientCommand:        disableClientCommand,
		regenerateClientUUIDCommand: regenerateClientUUIDCommand,
		reprovisionClientCommand:    reprovisionClientCommand,
		listAuditLogsQuery:          listAuditLogsQuery,
		getAuditLogQuery:            getAuditLogQuery,
		getAuditStatsQuery:          getAuditStatsQuery,
		getSystemHealthQuery:        getSystemHealthQuery,
		getSystemStatsQuery:         getSystemStatsQuery,
		getVersionQuery:             getVersionQuery,
		getDatabaseStatusQuery:      getDatabaseStatusQuery,
		getXraySystemStatusQuery:    getXraySystemStatusQuery,
	}
}

// HealthCheck is a simple endpoint to verify admin API is accessible
// This serves as a foundation test endpoint
// GET /api/v1/admin/health
func (h *AdminHandler) HealthCheck(c *gin.Context) {
	// Get admin context
	adminCtx, err := middleware.GetAdminContext(c)
	if err != nil {
		h.logger.Error("Failed to get admin context", "error", err)
		response.InternalError(c, "failed to get admin context")
		return
	}

	h.logger.Info("Admin health check accessed",
		"admin_id", adminCtx.UserID,
		"admin_email", adminCtx.Email,
		"ip", adminCtx.IPAddress)

	response.SuccessOK(c, gin.H{
		"status":      "ok",
		"message":     "admin API is operational",
		"admin_email": adminCtx.Email,
		"admin_id":    adminCtx.UserID.String(),
	})
}

// ListUsers handles GET /api/v1/admin/users
func (h *AdminHandler) ListUsers(c *gin.Context) {
	// Parse query parameters
	var req dto.AdminUserListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		h.logger.Error("Failed to bind query parameters", "error", err)
		response.BadRequest(c, "invalid query parameters")
		return
	}

	// Set defaults
	if req.Limit == 0 {
		req.Limit = 20
	}
	if req.SortBy == "" {
		req.SortBy = "created_at"
	}
	if req.SortOrder == "" {
		req.SortOrder = "desc"
	}

	// Execute query
	users, total, err := h.listUsersQuery.Execute(c.Request.Context(), adminuser.ListUsersParams{
		Offset:    req.Offset,
		Limit:     req.Limit,
		Role:      req.Role,
		Status:    req.Status,
		Email:     req.Email,
		SortBy:    req.SortBy,
		SortOrder: req.SortOrder,
	})
	if err != nil {
		h.logger.Error("Failed to list users", "error", err)
		response.InternalError(c, "failed to list users")
		return
	}

	// Get admin context for logging
	adminCtx, _ := middleware.GetAdminContext(c)
	if adminCtx != nil {
		h.logger.Info("Admin listed users",
			"admin_id", adminCtx.UserID,
			"admin_email", adminCtx.Email,
			"total_results", total,
			"filters", map[string]interface{}{
				"role":   req.Role,
				"status": req.Status,
				"email":  req.Email,
			})
	}

	// Map to response
	resp := mapper.ToAdminUserListResponse(users, total, req.Offset, req.Limit)
	response.SuccessOK(c, resp)
}

// GetUser handles GET /api/v1/admin/users/:id
func (h *AdminHandler) GetUser(c *gin.Context) {
	// Parse user ID
	userIDStr := c.Param("id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		response.BadRequest(c, "invalid user ID")
		return
	}

	// Execute query
	u, err := h.getUserQuery.Execute(c.Request.Context(), userID)
	if err != nil {
		if errors.Is(err, user.ErrUserNotFound) {
			response.NotFound(c, "user not found")
			return
		}
		h.logger.Error("Failed to get user", "error", err, "user_id", userID)
		response.InternalError(c, "failed to get user")
		return
	}

	// Get admin context for logging
	adminCtx, _ := middleware.GetAdminContext(c)
	if adminCtx != nil {
		h.logger.Info("Admin retrieved user details",
			"admin_id", adminCtx.UserID,
			"admin_email", adminCtx.Email,
			"target_user_id", userID)
	}

	// Map to response
	resp := mapper.ToAdminUserResponse(u)
	response.SuccessOK(c, resp)
}

// UpdateUserStatus handles PUT /api/v1/admin/users/:id/status
func (h *AdminHandler) UpdateUserStatus(c *gin.Context) {
	// Parse user ID
	userIDStr := c.Param("id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		response.BadRequest(c, "invalid user ID")
		return
	}

	// Parse request body
	var req dto.UpdateUserStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Failed to bind request", "error", err)
		response.BadRequest(c, "invalid request body")
		return
	}

	// Get admin context
	adminCtx, err := middleware.GetAdminContext(c)
	if err != nil {
		h.logger.Error("Failed to get admin context", "error", err)
		response.InternalError(c, "failed to get admin context")
		return
	}

	// Execute command
	err = h.updateUserStatusCommand.Execute(
		c.Request.Context(),
		userID,
		req.Status,
		adminCtx.UserID,
		adminCtx.IPAddress,
		adminCtx.UserAgent,
	)
	if err != nil {
		if errors.Is(err, user.ErrUserNotFound) {
			response.NotFound(c, "user not found")
			return
		}
		if errors.Is(err, user.ErrUserAlreadyActive) ||
			errors.Is(err, user.ErrUserAlreadyInactive) ||
			errors.Is(err, user.ErrUserAlreadySuspended) {
			response.BadRequest(c, err.Error())
			return
		}
		h.logger.Error("Failed to update user status", "error", err, "user_id", userID)
		response.InternalError(c, "failed to update user status")
		return
	}

	h.logger.Info("Admin updated user status",
		"admin_id", adminCtx.UserID,
		"admin_email", adminCtx.Email,
		"target_user_id", userID,
		"new_status", req.Status)

	// Fetch updated user
	u, err := h.getUserQuery.Execute(c.Request.Context(), userID)
	if err != nil {
		// Status was updated but fetch failed - return success anyway
		response.SuccessOK(c, gin.H{
			"message": "user status updated successfully",
		})
		return
	}

	resp := mapper.ToAdminUserResponse(u)
	response.SuccessOK(c, resp)
}

// UpdateUserRole handles PUT /api/v1/admin/users/:id/role
func (h *AdminHandler) UpdateUserRole(c *gin.Context) {
	// Parse user ID
	userIDStr := c.Param("id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		response.BadRequest(c, "invalid user ID")
		return
	}

	// Parse request body
	var req dto.UpdateUserRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Failed to bind request", "error", err)
		response.BadRequest(c, "invalid request body")
		return
	}

	// Get admin context
	adminCtx, err := middleware.GetAdminContext(c)
	if err != nil {
		h.logger.Error("Failed to get admin context", "error", err)
		response.InternalError(c, "failed to get admin context")
		return
	}

	// Execute command
	err = h.updateUserRoleCommand.Execute(
		c.Request.Context(),
		userID,
		req.Role,
		adminCtx.UserID,
		adminCtx.IPAddress,
		adminCtx.UserAgent,
	)
	if err != nil {
		if errors.Is(err, user.ErrUserNotFound) {
			response.NotFound(c, "user not found")
			return
		}
		if errors.Is(err, middleware.ErrAdminSelfDemotion) {
			response.BadRequest(c, "cannot demote yourself")
			return
		}
		h.logger.Error("Failed to update user role", "error", err, "user_id", userID)
		response.InternalError(c, "failed to update user role")
		return
	}

	h.logger.Info("Admin updated user role",
		"admin_id", adminCtx.UserID,
		"admin_email", adminCtx.Email,
		"target_user_id", userID,
		"new_role", req.Role)

	// Fetch updated user
	u, err := h.getUserQuery.Execute(c.Request.Context(), userID)
	if err != nil {
		// Role was updated but fetch failed - return success anyway
		response.SuccessOK(c, gin.H{
			"message": "user role updated successfully",
		})
		return
	}

	resp := mapper.ToAdminUserResponse(u)
	response.SuccessOK(c, resp)
}


// ========================= XRAY INSTANCE MANAGEMENT =========================

// ListInstances handles GET /api/v1/admin/xray/instances
func (h *AdminHandler) ListInstances(c *gin.Context) {
	// Parse query parameters
	var req dto.AdminXrayInstanceListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		h.logger.Error("Failed to bind query parameters", "error", err)
		response.BadRequest(c, "invalid query parameters")
		return
	}

	// Set defaults
	if req.Limit == 0 {
		req.Limit = 20
	}
	if req.SortBy == "" {
		req.SortBy = "created_at"
	}
	if req.SortOrder == "" {
		req.SortOrder = "desc"
	}

	// Execute query
	instances, total, err := h.listInstancesQuery.Execute(c.Request.Context(), adminxray.ListInstancesParams{
		Offset:    req.Offset,
		Limit:     req.Limit,
		NodeID:    req.NodeID,
		Status:    req.Status,
		SortBy:    req.SortBy,
		SortOrder: req.SortOrder,
	})
	if err != nil {
		h.logger.Error("Failed to list instances", "error", err)
		response.InternalError(c, "failed to list instances")
		return
	}

	// Get admin context for logging
	adminCtx, _ := middleware.GetAdminContext(c)
	if adminCtx != nil {
		h.logger.Info("Admin listed Xray instances",
			"admin_id", adminCtx.UserID,
			"admin_email", adminCtx.Email,
			"total_results", total,
			"filters", map[string]interface{}{
				"node_id": req.NodeID,
				"status":  req.Status,
			})
	}

	// Map to response
	resp := mapper.ToAdminXrayInstanceListResponse(instances, total, req.Offset, req.Limit)
	response.SuccessOK(c, resp)
}

// GetInstance handles GET /api/v1/admin/xray/instances/:id
func (h *AdminHandler) GetInstance(c *gin.Context) {
	// Parse instance ID
	instanceIDStr := c.Param("id")
	instanceID, err := uuid.Parse(instanceIDStr)
	if err != nil {
		response.BadRequest(c, "invalid instance ID")
		return
	}

	// Execute query
	instance, err := h.getInstanceQuery.Execute(c.Request.Context(), instanceID)
	if err != nil {
		if errors.Is(err, xray.ErrInstanceNotFound) {
			response.NotFound(c, "instance not found")
			return
		}
		h.logger.Error("Failed to get instance", "error", err, "instance_id", instanceID)
		response.InternalError(c, "failed to get instance")
		return
	}

	// Get admin context for logging
	adminCtx, _ := middleware.GetAdminContext(c)
	if adminCtx != nil {
		h.logger.Info("Admin retrieved instance details",
			"admin_id", adminCtx.UserID,
			"admin_email", adminCtx.Email,
			"instance_id", instanceID)
	}

	// Map to response
	resp := mapper.ToAdminXrayInstanceResponse(instance)
	response.SuccessOK(c, resp)
}

// StartInstance handles POST /api/v1/admin/xray/instances/:id/start
func (h *AdminHandler) StartInstance(c *gin.Context) {
	// Parse instance ID
	instanceIDStr := c.Param("id")
	instanceID, err := uuid.Parse(instanceIDStr)
	if err != nil {
		response.BadRequest(c, "invalid instance ID")
		return
	}

	// Get admin context
	adminCtx, err := middleware.GetAdminContext(c)
	if err != nil {
		h.logger.Error("Failed to get admin context", "error", err)
		response.InternalError(c, "failed to get admin context")
		return
	}

	// Execute command
	err = h.startInstanceCommand.Execute(
		c.Request.Context(),
		instanceID,
		adminCtx.UserID,
		adminCtx.IPAddress,
		adminCtx.UserAgent,
	)
	if err != nil {
		if errors.Is(err, xray.ErrInstanceNotFound) {
			response.NotFound(c, "instance not found")
			return
		}
		if errors.Is(err, xray.ErrInstanceAlreadyRunning) {
			response.BadRequest(c, "instance is already running")
			return
		}
		h.logger.Error("Failed to start instance", "error", err, "instance_id", instanceID)
		response.InternalError(c, "failed to start instance")
		return
	}

	h.logger.Info("Admin started instance",
		"admin_id", adminCtx.UserID,
		"admin_email", adminCtx.Email,
		"instance_id", instanceID)

	// Fetch updated instance
	instance, _ := h.getInstanceQuery.Execute(c.Request.Context(), instanceID)
	resp := dto.XrayInstanceControlResponse{
		Success:  true,
		Message:  "instance started successfully",
		Instance: mapper.ToAdminXrayInstanceResponse(instance),
	}
	response.SuccessOK(c, resp)
}

// StopInstance handles POST /api/v1/admin/xray/instances/:id/stop
func (h *AdminHandler) StopInstance(c *gin.Context) {
	// Parse instance ID
	instanceIDStr := c.Param("id")
	instanceID, err := uuid.Parse(instanceIDStr)
	if err != nil {
		response.BadRequest(c, "invalid instance ID")
		return
	}

	// Get admin context
	adminCtx, err := middleware.GetAdminContext(c)
	if err != nil {
		h.logger.Error("Failed to get admin context", "error", err)
		response.InternalError(c, "failed to get admin context")
		return
	}

	// Execute command
	err = h.stopInstanceCommand.Execute(
		c.Request.Context(),
		instanceID,
		adminCtx.UserID,
		adminCtx.IPAddress,
		adminCtx.UserAgent,
	)
	if err != nil {
		if errors.Is(err, xray.ErrInstanceNotFound) {
			response.NotFound(c, "instance not found")
			return
		}
		if errors.Is(err, xray.ErrInstanceAlreadyStopped) {
			response.BadRequest(c, "instance is already stopped")
			return
		}
		h.logger.Error("Failed to stop instance", "error", err, "instance_id", instanceID)
		response.InternalError(c, "failed to stop instance")
		return
	}

	h.logger.Info("Admin stopped instance",
		"admin_id", adminCtx.UserID,
		"admin_email", adminCtx.Email,
		"instance_id", instanceID)

	// Fetch updated instance
	instance, _ := h.getInstanceQuery.Execute(c.Request.Context(), instanceID)
	resp := dto.XrayInstanceControlResponse{
		Success:  true,
		Message:  "instance stopped successfully",
		Instance: mapper.ToAdminXrayInstanceResponse(instance),
	}
	response.SuccessOK(c, resp)
}

// RestartInstance handles POST /api/v1/admin/xray/instances/:id/restart
func (h *AdminHandler) RestartInstance(c *gin.Context) {
	// Parse instance ID
	instanceIDStr := c.Param("id")
	instanceID, err := uuid.Parse(instanceIDStr)
	if err != nil {
		response.BadRequest(c, "invalid instance ID")
		return
	}

	// Get admin context
	adminCtx, err := middleware.GetAdminContext(c)
	if err != nil {
		h.logger.Error("Failed to get admin context", "error", err)
		response.InternalError(c, "failed to get admin context")
		return
	}

	// Execute command
	err = h.restartInstanceCommand.Execute(
		c.Request.Context(),
		instanceID,
		adminCtx.UserID,
		adminCtx.IPAddress,
		adminCtx.UserAgent,
	)
	if err != nil {
		if errors.Is(err, xray.ErrInstanceNotFound) {
			response.NotFound(c, "instance not found")
			return
		}
		h.logger.Error("Failed to restart instance", "error", err, "instance_id", instanceID)
		response.InternalError(c, "failed to restart instance")
		return
	}

	h.logger.Info("Admin restarted instance",
		"admin_id", adminCtx.UserID,
		"admin_email", adminCtx.Email,
		"instance_id", instanceID)

	// Fetch updated instance
	instance, _ := h.getInstanceQuery.Execute(c.Request.Context(), instanceID)
	resp := dto.XrayInstanceControlResponse{
		Success:  true,
		Message:  "instance restarted successfully",
		Instance: mapper.ToAdminXrayInstanceResponse(instance),
	}
	response.SuccessOK(c, resp)
}

// ReloadInstance handles POST /api/v1/admin/xray/instances/:id/reload
func (h *AdminHandler) ReloadInstance(c *gin.Context) {
	// Parse instance ID
	instanceIDStr := c.Param("id")
	instanceID, err := uuid.Parse(instanceIDStr)
	if err != nil {
		response.BadRequest(c, "invalid instance ID")
		return
	}

	// Get admin context
	adminCtx, err := middleware.GetAdminContext(c)
	if err != nil {
		h.logger.Error("Failed to get admin context", "error", err)
		response.InternalError(c, "failed to get admin context")
		return
	}

	startTime := time.Now()

	// Execute command (REUSES XrayProvisioningService)
	err = h.reloadInstanceCommand.Execute(
		c.Request.Context(),
		instanceID,
		adminCtx.UserID,
		adminCtx.IPAddress,
		adminCtx.UserAgent,
	)
	if err != nil {
		if errors.Is(err, xray.ErrInstanceNotFound) {
			response.NotFound(c, "instance not found")
			return
		}
		h.logger.Error("Failed to reload instance", "error", err, "instance_id", instanceID)
		response.InternalError(c, "failed to reload instance config")
		return
	}

	duration := time.Since(startTime)

	h.logger.Info("Admin reloaded instance config",
		"admin_id", adminCtx.UserID,
		"admin_email", adminCtx.Email,
		"instance_id", instanceID,
		"duration_ms", duration.Milliseconds())

	resp := dto.ConfigRegenerateResponse{
		Success:       true,
		Message:       "instance config reloaded successfully",
		ConfigSize:    0, // Not available without re-reading file
		RegeneratedAt: time.Now(),
		DurationMs:    duration.Milliseconds(),
	}
	response.SuccessOK(c, resp)
}

// CheckInstanceHealth handles GET /api/v1/admin/xray/instances/:id/health
func (h *AdminHandler) CheckInstanceHealth(c *gin.Context) {
	// Parse instance ID
	instanceIDStr := c.Param("id")
	instanceID, err := uuid.Parse(instanceIDStr)
	if err != nil {
		response.BadRequest(c, "invalid instance ID")
		return
	}

	// Execute command
	result, err := h.checkInstanceHealthCommand.Execute(c.Request.Context(), instanceID)
	if err != nil {
		h.logger.Error("Failed to check instance health", "error", err, "instance_id", instanceID)
		response.InternalError(c, "failed to check instance health")
		return
	}

	// Get admin context for logging
	adminCtx, _ := middleware.GetAdminContext(c)
	if adminCtx != nil {
		h.logger.Info("Admin checked instance health",
			"admin_id", adminCtx.UserID,
			"admin_email", adminCtx.Email,
			"instance_id", instanceID,
			"healthy", result.Healthy)
	}

	resp := dto.HealthCheckResponse{
		Healthy:   result.Healthy,
		Status:    result.Status,
		Message:   result.Message,
		CheckedAt: time.Now(),
	}
	response.SuccessOK(c, resp)
}

// GetInstanceStats handles GET /api/v1/admin/xray/instances/:id/stats
func (h *AdminHandler) GetInstanceStats(c *gin.Context) {
	// Parse instance ID
	instanceIDStr := c.Param("id")
	instanceID, err := uuid.Parse(instanceIDStr)
	if err != nil {
		response.BadRequest(c, "invalid instance ID")
		return
	}

	// Get instance first
	instance, err := h.getInstanceQuery.Execute(c.Request.Context(), instanceID)
	if err != nil {
		if errors.Is(err, xray.ErrInstanceNotFound) {
			response.NotFound(c, "instance not found")
			return
		}
		h.logger.Error("Failed to get instance", "error", err, "instance_id", instanceID)
		response.InternalError(c, "failed to get instance")
		return
	}

	// Get stats
	stats, err := h.getInstanceStatsQuery.Execute(c.Request.Context(), instanceID)
	if err != nil {
		h.logger.Error("Failed to get instance stats", "error", err, "instance_id", instanceID)
		response.InternalError(c, "failed to get instance stats")
		return
	}

	// Get admin context for logging
	adminCtx, _ := middleware.GetAdminContext(c)
	if adminCtx != nil {
		h.logger.Info("Admin retrieved instance stats",
			"admin_id", adminCtx.UserID,
			"admin_email", adminCtx.Email,
			"instance_id", instanceID)
	}

	uptime := int64(0)
	if instance.IsRunning() {
		uptime = int64(instance.GetUptime().Seconds())
	}

	resp := dto.InstanceStatsResponse{
		InstanceID:      instanceID,
		TotalInbounds:   stats.TotalInbounds,
		EnabledInbounds: stats.EnabledInbounds,
		TotalClients:    stats.TotalClients,
		EnabledClients:  stats.EnabledClients,
		ConfigVersion:   instance.ConfigVersion,
		Uptime:          uptime,
	}
	response.SuccessOK(c, resp)
}


// ========================= INBOUND MANAGEMENT =========================

// ListInbounds handles GET /api/v1/admin/xray/inbounds
func (h *AdminHandler) ListInbounds(c *gin.Context) {
	var req dto.AdminInboundListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		h.logger.Error("Failed to bind query parameters", "error", err)
		response.BadRequest(c, "invalid query parameters")
		return
	}

	if req.Limit == 0 {
		req.Limit = 20
	}
	if req.SortBy == "" {
		req.SortBy = "created_at"
	}
	if req.SortOrder == "" {
		req.SortOrder = "desc"
	}

	inbounds, total, err := h.listInboundsQuery.Execute(c.Request.Context(), admininbound.ListInboundsParams{
		Offset:     req.Offset,
		Limit:      req.Limit,
		InstanceID: req.InstanceID,
		Protocol:   req.Protocol,
		Enabled:    req.Enabled,
		SortBy:     req.SortBy,
		SortOrder:  req.SortOrder,
	})
	if err != nil {
		h.logger.Error("Failed to list inbounds", "error", err)
		response.InternalError(c, "failed to list inbounds")
		return
	}

	adminCtx, _ := middleware.GetAdminContext(c)
	if adminCtx != nil {
		h.logger.Info("Admin listed inbounds",
			"admin_id", adminCtx.UserID,
			"total_results", total)
	}

	resp := mapper.ToAdminInboundListResponse(inbounds, total, req.Offset, req.Limit)
	response.SuccessOK(c, resp)
}

// GetInbound handles GET /api/v1/admin/xray/inbounds/:id
func (h *AdminHandler) GetInbound(c *gin.Context) {
	inboundIDStr := c.Param("id")
	inboundID, err := uuid.Parse(inboundIDStr)
	if err != nil {
		response.BadRequest(c, "invalid inbound ID")
		return
	}

	inbound, err := h.getInboundQuery.Execute(c.Request.Context(), inboundID)
	if err != nil {
		if errors.Is(err, xray.ErrInboundNotFound) {
			response.NotFound(c, "inbound not found")
			return
		}
		h.logger.Error("Failed to get inbound", "error", err, "inbound_id", inboundID)
		response.InternalError(c, "failed to get inbound")
		return
	}

	resp := mapper.ToAdminInboundResponse(inbound)
	response.SuccessOK(c, resp)
}

// CreateInbound handles POST /api/v1/admin/xray/inbounds
func (h *AdminHandler) CreateInbound(c *gin.Context) {
	var req dto.AdminCreateInboundRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Failed to bind request", "error", err)
		response.BadRequest(c, "invalid request body")
		return
	}

	instanceID, err := uuid.Parse(req.XrayInstanceID)
	if err != nil {
		response.BadRequest(c, "invalid instance ID")
		return
	}

	adminCtx, err := middleware.GetAdminContext(c)
	if err != nil {
		h.logger.Error("Failed to get admin context", "error", err)
		response.InternalError(c, "failed to get admin context")
		return
	}

	inbound, err := h.createInboundCommand.Execute(
		c.Request.Context(),
		instanceID,
		xray.InboundProtocol(req.Protocol),
		req.Port,
		xray.TransportType(req.Transport),
		xray.SecurityType(req.Security),
		adminCtx.UserID,
		adminCtx.IPAddress,
		adminCtx.UserAgent,
	)
	if err != nil {
		if errors.Is(err, xray.ErrInstanceNotFound) {
			response.NotFound(c, "instance not found")
			return
		}
		if errors.Is(err, xray.ErrInvalidPort) ||
			errors.Is(err, xray.ErrInvalidProtocol) ||
			errors.Is(err, xray.ErrInvalidTransport) ||
			errors.Is(err, xray.ErrInvalidSecurity) {
			response.BadRequest(c, err.Error())
			return
		}
		h.logger.Error("Failed to create inbound", "error", err)
		response.InternalError(c, "failed to create inbound")
		return
	}

	h.logger.Info("Admin created inbound",
		"admin_id", adminCtx.UserID,
		"inbound_id", inbound.ID,
		"instance_id", instanceID)

	resp := dto.InboundOperationResponse{
		Success: true,
		Message: "inbound created successfully",
		Inbound: mapper.ToAdminInboundResponse(inbound),
	}
	response.SuccessCreated(c, resp)
}

// UpdateInbound handles PUT /api/v1/admin/xray/inbounds/:id
func (h *AdminHandler) UpdateInbound(c *gin.Context) {
	inboundIDStr := c.Param("id")
	inboundID, err := uuid.Parse(inboundIDStr)
	if err != nil {
		response.BadRequest(c, "invalid inbound ID")
		return
	}

	var req dto.AdminUpdateInboundRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Failed to bind request", "error", err)
		response.BadRequest(c, "invalid request body")
		return
	}

	adminCtx, err := middleware.GetAdminContext(c)
	if err != nil {
		h.logger.Error("Failed to get admin context", "error", err)
		response.InternalError(c, "failed to get admin context")
		return
	}

	var port *int
	var transport *xray.TransportType
	var security *xray.SecurityType

	if req.Port > 0 {
		port = &req.Port
	}
	if req.Transport != "" {
		t := xray.TransportType(req.Transport)
		transport = &t
	}
	if req.Security != "" {
		s := xray.SecurityType(req.Security)
		security = &s
	}

	inbound, err := h.updateInboundCommand.Execute(
		c.Request.Context(),
		inboundID,
		port,
		transport,
		security,
		adminCtx.UserID,
		adminCtx.IPAddress,
		adminCtx.UserAgent,
	)
	if err != nil {
		if errors.Is(err, xray.ErrInboundNotFound) {
			response.NotFound(c, "inbound not found")
			return
		}
		if errors.Is(err, xray.ErrInvalidPort) ||
			errors.Is(err, xray.ErrInvalidTransport) ||
			errors.Is(err, xray.ErrInvalidSecurity) {
			response.BadRequest(c, err.Error())
			return
		}
		h.logger.Error("Failed to update inbound", "error", err, "inbound_id", inboundID)
		response.InternalError(c, "failed to update inbound")
		return
	}

	h.logger.Info("Admin updated inbound",
		"admin_id", adminCtx.UserID,
		"inbound_id", inboundID)

	resp := dto.InboundOperationResponse{
		Success: true,
		Message: "inbound updated successfully",
		Inbound: mapper.ToAdminInboundResponse(inbound),
	}
	response.SuccessOK(c, resp)
}

// DeleteInbound handles DELETE /api/v1/admin/xray/inbounds/:id
func (h *AdminHandler) DeleteInbound(c *gin.Context) {
	inboundIDStr := c.Param("id")
	inboundID, err := uuid.Parse(inboundIDStr)
	if err != nil {
		response.BadRequest(c, "invalid inbound ID")
		return
	}

	adminCtx, err := middleware.GetAdminContext(c)
	if err != nil {
		h.logger.Error("Failed to get admin context", "error", err)
		response.InternalError(c, "failed to get admin context")
		return
	}

	err = h.deleteInboundCommand.Execute(
		c.Request.Context(),
		inboundID,
		adminCtx.UserID,
		adminCtx.IPAddress,
		adminCtx.UserAgent,
	)
	if err != nil {
		if errors.Is(err, xray.ErrInboundNotFound) {
			response.NotFound(c, "inbound not found")
			return
		}
		h.logger.Error("Failed to delete inbound", "error", err, "inbound_id", inboundID)
		response.InternalError(c, "failed to delete inbound")
		return
	}

	h.logger.Info("Admin deleted inbound",
		"admin_id", adminCtx.UserID,
		"inbound_id", inboundID)

	response.SuccessOK(c, gin.H{
		"success": true,
		"message": "inbound deleted successfully",
	})
}

// EnableInbound handles PUT /api/v1/admin/xray/inbounds/:id/enable
func (h *AdminHandler) EnableInbound(c *gin.Context) {
	inboundIDStr := c.Param("id")
	inboundID, err := uuid.Parse(inboundIDStr)
	if err != nil {
		response.BadRequest(c, "invalid inbound ID")
		return
	}

	adminCtx, err := middleware.GetAdminContext(c)
	if err != nil {
		h.logger.Error("Failed to get admin context", "error", err)
		response.InternalError(c, "failed to get admin context")
		return
	}

	err = h.enableInboundCommand.Execute(
		c.Request.Context(),
		inboundID,
		adminCtx.UserID,
		adminCtx.IPAddress,
		adminCtx.UserAgent,
	)
	if err != nil {
		if errors.Is(err, xray.ErrInboundNotFound) {
			response.NotFound(c, "inbound not found")
			return
		}
		if errors.Is(err, xray.ErrInboundAlreadyEnabled) {
			response.BadRequest(c, "inbound is already enabled")
			return
		}
		h.logger.Error("Failed to enable inbound", "error", err, "inbound_id", inboundID)
		response.InternalError(c, "failed to enable inbound")
		return
	}

	h.logger.Info("Admin enabled inbound",
		"admin_id", adminCtx.UserID,
		"inbound_id", inboundID)

	inbound, _ := h.getInboundQuery.Execute(c.Request.Context(), inboundID)
	resp := dto.InboundOperationResponse{
		Success: true,
		Message: "inbound enabled successfully",
		Inbound: mapper.ToAdminInboundResponse(inbound),
	}
	response.SuccessOK(c, resp)
}

// DisableInbound handles PUT /api/v1/admin/xray/inbounds/:id/disable
func (h *AdminHandler) DisableInbound(c *gin.Context) {
	inboundIDStr := c.Param("id")
	inboundID, err := uuid.Parse(inboundIDStr)
	if err != nil {
		response.BadRequest(c, "invalid inbound ID")
		return
	}

	adminCtx, err := middleware.GetAdminContext(c)
	if err != nil {
		h.logger.Error("Failed to get admin context", "error", err)
		response.InternalError(c, "failed to get admin context")
		return
	}

	err = h.disableInboundCommand.Execute(
		c.Request.Context(),
		inboundID,
		adminCtx.UserID,
		adminCtx.IPAddress,
		adminCtx.UserAgent,
	)
	if err != nil {
		if errors.Is(err, xray.ErrInboundNotFound) {
			response.NotFound(c, "inbound not found")
			return
		}
		if errors.Is(err, xray.ErrInboundAlreadyDisabled) {
			response.BadRequest(c, "inbound is already disabled")
			return
		}
		h.logger.Error("Failed to disable inbound", "error", err, "inbound_id", inboundID)
		response.InternalError(c, "failed to disable inbound")
		return
	}

	h.logger.Info("Admin disabled inbound",
		"admin_id", adminCtx.UserID,
		"inbound_id", inboundID)

	inbound, _ := h.getInboundQuery.Execute(c.Request.Context(), inboundID)
	resp := dto.InboundOperationResponse{
		Success: true,
		Message: "inbound disabled successfully",
		Inbound: mapper.ToAdminInboundResponse(inbound),
	}
	response.SuccessOK(c, resp)
}


// ========================= CLIENT MANAGEMENT =========================

// ListClients handles GET /api/v1/admin/xray/clients
func (h *AdminHandler) ListClients(c *gin.Context) {
	var req dto.AdminClientListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		h.logger.Error("Failed to bind query parameters", "error", err)
		response.BadRequest(c, "invalid query parameters")
		return
	}

	if req.Limit == 0 {
		req.Limit = 20
	}
	if req.SortBy == "" {
		req.SortBy = "created_at"
	}
	if req.SortOrder == "" {
		req.SortOrder = "desc"
	}

	clients, total, err := h.listClientsQuery.Execute(c.Request.Context(), adminclient.ListClientsParams{
		Offset:    req.Offset,
		Limit:     req.Limit,
		InboundID: req.InboundID,
		UserID:    req.UserID,
		Enabled:   req.Enabled,
		SortBy:    req.SortBy,
		SortOrder: req.SortOrder,
	})
	if err != nil {
		h.logger.Error("Failed to list clients", "error", err)
		response.InternalError(c, "failed to list clients")
		return
	}

	adminCtx, _ := middleware.GetAdminContext(c)
	if adminCtx != nil {
		h.logger.Info("Admin listed clients", "admin_id", adminCtx.UserID, "total_results", total)
	}

	resp := mapper.ToAdminClientListResponse(clients, total, req.Offset, req.Limit)
	response.SuccessOK(c, resp)
}

// GetClient handles GET /api/v1/admin/xray/clients/:id
func (h *AdminHandler) GetClient(c *gin.Context) {
	clientIDStr := c.Param("id")
	clientID, err := uuid.Parse(clientIDStr)
	if err != nil {
		response.BadRequest(c, "invalid client ID")
		return
	}

	client, err := h.getClientQuery.Execute(c.Request.Context(), clientID)
	if err != nil {
		if errors.Is(err, xray.ErrClientNotFound) {
			response.NotFound(c, "client not found")
			return
		}
		h.logger.Error("Failed to get client", "error", err, "client_id", clientID)
		response.InternalError(c, "failed to get client")
		return
	}

	resp := mapper.ToAdminClientResponse(client)
	response.SuccessOK(c, resp)
}

// CreateClient handles POST /api/v1/admin/xray/clients
func (h *AdminHandler) CreateClient(c *gin.Context) {
	var req dto.AdminCreateClientRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Failed to bind request", "error", err)
		response.BadRequest(c, "invalid request body")
		return
	}

	inboundID, err := uuid.Parse(req.InboundID)
	if err != nil {
		response.BadRequest(c, "invalid inbound ID")
		return
	}

	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		response.BadRequest(c, "invalid user ID")
		return
	}

	adminCtx, err := middleware.GetAdminContext(c)
	if err != nil {
		h.logger.Error("Failed to get admin context", "error", err)
		response.InternalError(c, "failed to get admin context")
		return
	}

	client, err := h.createClientCommand.Execute(
		c.Request.Context(),
		inboundID,
		userID,
		req.Email,
		req.Flow,
		adminCtx.UserID,
		adminCtx.IPAddress,
		adminCtx.UserAgent,
	)
	if err != nil {
		if errors.Is(err, xray.ErrInboundNotFound) {
			response.NotFound(c, "inbound not found")
			return
		}
		h.logger.Error("Failed to create client", "error", err)
		response.InternalError(c, "failed to create client")
		return
	}

	h.logger.Info("Admin created client", "admin_id", adminCtx.UserID, "client_id", client.ID)

	resp := dto.ClientOperationResponse{
		Success: true,
		Message: "client created successfully",
		Client:  mapper.ToAdminClientResponse(client),
	}
	response.SuccessCreated(c, resp)
}

// DeleteClient handles DELETE /api/v1/admin/xray/clients/:id
func (h *AdminHandler) DeleteClient(c *gin.Context) {
	clientIDStr := c.Param("id")
	clientID, err := uuid.Parse(clientIDStr)
	if err != nil {
		response.BadRequest(c, "invalid client ID")
		return
	}

	adminCtx, err := middleware.GetAdminContext(c)
	if err != nil {
		h.logger.Error("Failed to get admin context", "error", err)
		response.InternalError(c, "failed to get admin context")
		return
	}

	err = h.deleteClientCommand.Execute(
		c.Request.Context(),
		clientID,
		adminCtx.UserID,
		adminCtx.IPAddress,
		adminCtx.UserAgent,
	)
	if err != nil {
		if errors.Is(err, xray.ErrClientNotFound) {
			response.NotFound(c, "client not found")
			return
		}
		h.logger.Error("Failed to delete client", "error", err, "client_id", clientID)
		response.InternalError(c, "failed to delete client")
		return
	}

	h.logger.Info("Admin deleted client", "admin_id", adminCtx.UserID, "client_id", clientID)

	response.SuccessOK(c, gin.H{
		"success": true,
		"message": "client deleted successfully",
	})
}

// EnableClient handles PUT /api/v1/admin/xray/clients/:id/enable
func (h *AdminHandler) EnableClient(c *gin.Context) {
	clientIDStr := c.Param("id")
	clientID, err := uuid.Parse(clientIDStr)
	if err != nil {
		response.BadRequest(c, "invalid client ID")
		return
	}

	adminCtx, err := middleware.GetAdminContext(c)
	if err != nil {
		h.logger.Error("Failed to get admin context", "error", err)
		response.InternalError(c, "failed to get admin context")
		return
	}

	err = h.enableClientCommand.Execute(
		c.Request.Context(),
		clientID,
		adminCtx.UserID,
		adminCtx.IPAddress,
		adminCtx.UserAgent,
	)
	if err != nil {
		if errors.Is(err, xray.ErrClientNotFound) {
			response.NotFound(c, "client not found")
			return
		}
		if errors.Is(err, xray.ErrClientAlreadyEnabled) {
			response.BadRequest(c, "client is already enabled")
			return
		}
		h.logger.Error("Failed to enable client", "error", err, "client_id", clientID)
		response.InternalError(c, "failed to enable client")
		return
	}

	h.logger.Info("Admin enabled client", "admin_id", adminCtx.UserID, "client_id", clientID)

	client, _ := h.getClientQuery.Execute(c.Request.Context(), clientID)
	resp := dto.ClientOperationResponse{
		Success: true,
		Message: "client enabled successfully",
		Client:  mapper.ToAdminClientResponse(client),
	}
	response.SuccessOK(c, resp)
}

// DisableClient handles PUT /api/v1/admin/xray/clients/:id/disable
func (h *AdminHandler) DisableClient(c *gin.Context) {
	clientIDStr := c.Param("id")
	clientID, err := uuid.Parse(clientIDStr)
	if err != nil {
		response.BadRequest(c, "invalid client ID")
		return
	}

	adminCtx, err := middleware.GetAdminContext(c)
	if err != nil {
		h.logger.Error("Failed to get admin context", "error", err)
		response.InternalError(c, "failed to get admin context")
		return
	}

	err = h.disableClientCommand.Execute(
		c.Request.Context(),
		clientID,
		adminCtx.UserID,
		adminCtx.IPAddress,
		adminCtx.UserAgent,
	)
	if err != nil {
		if errors.Is(err, xray.ErrClientNotFound) {
			response.NotFound(c, "client not found")
			return
		}
		if errors.Is(err, xray.ErrClientAlreadyDisabled) {
			response.BadRequest(c, "client is already disabled")
			return
		}
		h.logger.Error("Failed to disable client", "error", err, "client_id", clientID)
		response.InternalError(c, "failed to disable client")
		return
	}

	h.logger.Info("Admin disabled client", "admin_id", adminCtx.UserID, "client_id", clientID)

	client, _ := h.getClientQuery.Execute(c.Request.Context(), clientID)
	resp := dto.ClientOperationResponse{
		Success: true,
		Message: "client disabled successfully",
		Client:  mapper.ToAdminClientResponse(client),
	}
	response.SuccessOK(c, resp)
}

// RegenerateClientUUID handles POST /api/v1/admin/xray/clients/:id/regenerate-uuid
func (h *AdminHandler) RegenerateClientUUID(c *gin.Context) {
	clientIDStr := c.Param("id")
	clientID, err := uuid.Parse(clientIDStr)
	if err != nil {
		response.BadRequest(c, "invalid client ID")
		return
	}

	adminCtx, err := middleware.GetAdminContext(c)
	if err != nil {
		h.logger.Error("Failed to get admin context", "error", err)
		response.InternalError(c, "failed to get admin context")
		return
	}

	client, err := h.regenerateClientUUIDCommand.Execute(
		c.Request.Context(),
		clientID,
		adminCtx.UserID,
		adminCtx.IPAddress,
		adminCtx.UserAgent,
	)
	if err != nil {
		if errors.Is(err, xray.ErrClientNotFound) {
			response.NotFound(c, "client not found")
			return
		}
		h.logger.Error("Failed to regenerate client UUID", "error", err, "client_id", clientID)
		response.InternalError(c, "failed to regenerate client UUID")
		return
	}

	h.logger.Info("Admin regenerated client UUID", "admin_id", adminCtx.UserID, "client_id", clientID)

	resp := dto.ClientOperationResponse{
		Success: true,
		Message: "client UUID regenerated successfully",
		Client:  mapper.ToAdminClientResponse(client),
	}
	response.SuccessOK(c, resp)
}

// ReprovisionClient handles POST /api/v1/admin/xray/clients/:id/reprovision
func (h *AdminHandler) ReprovisionClient(c *gin.Context) {
	clientIDStr := c.Param("id")
	clientID, err := uuid.Parse(clientIDStr)
	if err != nil {
		response.BadRequest(c, "invalid client ID")
		return
	}

	// Check if should regenerate UUID (query param)
	regenerateUUID := c.Query("regenerate_uuid") == "true"

	adminCtx, err := middleware.GetAdminContext(c)
	if err != nil {
		h.logger.Error("Failed to get admin context", "error", err)
		response.InternalError(c, "failed to get admin context")
		return
	}

	client, err := h.reprovisionClientCommand.Execute(
		c.Request.Context(),
		clientID,
		adminCtx.UserID,
		regenerateUUID,
		adminCtx.IPAddress,
		adminCtx.UserAgent,
	)
	if err != nil {
		if errors.Is(err, xray.ErrClientNotFound) {
			response.NotFound(c, "client not found")
			return
		}
		h.logger.Error("Failed to reprovision client", "error", err, "client_id", clientID)
		response.InternalError(c, "failed to reprovision client")
		return
	}

	h.logger.Info("Admin reprovisioned client",
		"admin_id", adminCtx.UserID,
		"client_id", clientID,
		"regenerate_uuid", regenerateUUID)

	resp := dto.ClientOperationResponse{
		Success: true,
		Message: "client reprovisioned successfully",
		Client:  mapper.ToAdminClientResponse(client),
	}
	response.SuccessOK(c, resp)
}


// ========================= AUDIT LOG MANAGEMENT =========================

// ListAuditLogs handles GET /api/v1/admin/audit/logs
func (h *AdminHandler) ListAuditLogs(c *gin.Context) {
	var req dto.AdminAuditListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		h.logger.Error("Failed to bind query parameters", "error", err)
		response.BadRequest(c, "invalid query parameters")
		return
	}

	// Set defaults
	if req.Limit == 0 {
		req.Limit = 20
	}
	if req.Limit > 100 {
		req.Limit = 100 // Max limit for performance
	}
	if req.SortBy == "" {
		req.SortBy = "created_at"
	}
	if req.SortOrder == "" {
		req.SortOrder = "desc"
	}

	logs, total, err := h.listAuditLogsQuery.Execute(c.Request.Context(), adminaudit.ListAuditLogsParams{
		Offset:     req.Offset,
		Limit:      req.Limit,
		UserID:     req.UserID,
		Action:     req.Action,
		EntityType: req.EntityType,
		EntityID:   req.EntityID,
		IPAddress:  req.IPAddress,
		DateFrom:   req.DateFrom,
		DateTo:     req.DateTo,
		SortBy:     req.SortBy,
		SortOrder:  req.SortOrder,
	})
	if err != nil {
		h.logger.Error("Failed to list audit logs", "error", err)
		response.InternalError(c, "failed to list audit logs")
		return
	}

	adminCtx, _ := middleware.GetAdminContext(c)
	if adminCtx != nil {
		h.logger.Info("Admin listed audit logs",
			"admin_id", adminCtx.UserID,
			"total_results", total,
			"filters", map[string]interface{}{
				"user_id":     req.UserID,
				"action":      req.Action,
				"entity_type": req.EntityType,
				"entity_id":   req.EntityID,
				"ip_address":  req.IPAddress,
				"date_from":   req.DateFrom,
				"date_to":     req.DateTo,
			})
	}

	resp := mapper.ToAdminAuditListResponse(logs, total, req.Offset, req.Limit)
	response.SuccessOK(c, resp)
}

// GetAuditLog handles GET /api/v1/admin/audit/logs/:id
func (h *AdminHandler) GetAuditLog(c *gin.Context) {
	logIDStr := c.Param("id")
	logID, err := uuid.Parse(logIDStr)
	if err != nil {
		response.BadRequest(c, "invalid log ID")
		return
	}

	log, err := h.getAuditLogQuery.Execute(c.Request.Context(), logID)
	if err != nil {
		h.logger.Error("Failed to get audit log", "error", err, "log_id", logID)
		response.NotFound(c, "audit log not found")
		return
	}

	adminCtx, _ := middleware.GetAdminContext(c)
	if adminCtx != nil {
		h.logger.Info("Admin retrieved audit log details",
			"admin_id", adminCtx.UserID,
			"log_id", logID)
	}

	resp := mapper.ToAdminAuditResponse(log)
	response.SuccessOK(c, resp)
}

// GetAuditStats handles GET /api/v1/admin/audit/stats
func (h *AdminHandler) GetAuditStats(c *gin.Context) {
	stats, err := h.getAuditStatsQuery.Execute(c.Request.Context())
	if err != nil {
		h.logger.Error("Failed to get audit stats", "error", err)
		response.InternalError(c, "failed to get audit statistics")
		return
	}

	adminCtx, _ := middleware.GetAdminContext(c)
	if adminCtx != nil {
		h.logger.Info("Admin retrieved audit stats",
			"admin_id", adminCtx.UserID,
			"total_logs", stats.TotalLogs)
	}

	resp := dto.AuditStatsResponse{
		TotalLogs:         stats.TotalLogs,
		LogsByAction:      stats.LogsByAction,
		LogsByEntityType:  stats.LogsByEntityType,
		UniqueUsers:       stats.UniqueUsers,
		UniqueIPAddresses: stats.UniqueIPAddresses,
		OldestLog:         stats.OldestLog,
		NewestLog:         stats.NewestLog,
	}
	response.SuccessOK(c, resp)
}


// ========================= SYSTEM ADMIN =========================

// GetSystemHealth handles GET /api/v1/admin/system/health
func (h *AdminHandler) GetSystemHealth(c *gin.Context) {
	health, err := h.getSystemHealthQuery.Execute(c.Request.Context())
	if err != nil {
		h.logger.Error("Failed to get system health", "error", err)
		response.InternalError(c, "failed to get system health")
		return
	}

	adminCtx, _ := middleware.GetAdminContext(c)
	if adminCtx != nil {
		h.logger.Info("Admin checked system health",
			"admin_id", adminCtx.UserID,
			"status", health.Status)
	}

	// Map to DTO
	components := make(map[string]dto.ComponentHealth)
	for name, comp := range health.Components {
		components[name] = dto.ComponentHealth{
			Status:  comp.Status,
			Message: comp.Message,
			Latency: comp.Latency,
		}
	}

	resp := dto.SystemHealthResponse{
		Status:     health.Status,
		Timestamp:  time.Now(),
		Components: components,
		Version:    adminsystem.Version,
		Uptime:     health.Uptime,
	}
	response.SuccessOK(c, resp)
}

// GetSystemStats handles GET /api/v1/admin/system/stats
func (h *AdminHandler) GetSystemStats(c *gin.Context) {
	stats, err := h.getSystemStatsQuery.Execute(c.Request.Context())
	if err != nil {
		h.logger.Error("Failed to get system stats", "error", err)
		response.InternalError(c, "failed to get system statistics")
		return
	}

	adminCtx, _ := middleware.GetAdminContext(c)
	if adminCtx != nil {
		h.logger.Info("Admin retrieved system stats",
			"admin_id", adminCtx.UserID,
			"total_users", stats.Users.TotalUsers,
			"total_instances", stats.Xray.TotalInstances)
	}

	resp := dto.SystemStatsResponse{
		Users: dto.UserStats{
			TotalUsers:       stats.Users.TotalUsers,
			ActiveUsers:      stats.Users.ActiveUsers,
			SuspendedUsers:   stats.Users.SuspendedUsers,
			AdminUsers:       stats.Users.AdminUsers,
			NewUsersToday:    stats.Users.NewUsersToday,
			NewUsersThisWeek: stats.Users.NewUsersThisWeek,
		},
		Xray: dto.XrayStats{
			TotalInstances:   stats.Xray.TotalInstances,
			RunningInstances: stats.Xray.RunningInstances,
			StoppedInstances: stats.Xray.StoppedInstances,
			TotalInbounds:    stats.Xray.TotalInbounds,
			EnabledInbounds:  stats.Xray.EnabledInbounds,
			TotalClients:     stats.Xray.TotalClients,
			EnabledClients:   stats.Xray.EnabledClients,
		},
		Audit: dto.AuditStats{
			TotalLogs:        stats.Audit.TotalLogs,
			LogsToday:        stats.Audit.LogsToday,
			LogsThisWeek:     stats.Audit.LogsThisWeek,
			UniqueUsersToday: stats.Audit.UniqueUsersToday,
		},
		Timestamp: time.Now(),
	}
	response.SuccessOK(c, resp)
}

// GetVersion handles GET /api/v1/admin/system/version
func (h *AdminHandler) GetVersion(c *gin.Context) {
	version := h.getVersionQuery.Execute(c.Request.Context())

	adminCtx, _ := middleware.GetAdminContext(c)
	if adminCtx != nil {
		h.logger.Info("Admin retrieved version info", "admin_id", adminCtx.UserID)
	}

	resp := dto.VersionResponse{
		Version:     version.Version,
		BuildDate:   version.BuildDate,
		GitCommit:   version.GitCommit,
		GoVersion:   version.GoVersion,
		Environment: version.Environment,
		Uptime:      version.Uptime,
		StartTime:   version.StartTime,
	}
	response.SuccessOK(c, resp)
}

// GetDatabaseStatus handles GET /api/v1/admin/system/database
func (h *AdminHandler) GetDatabaseStatus(c *gin.Context) {
	dbStatus, err := h.getDatabaseStatusQuery.Execute(c.Request.Context())
	if err != nil {
		h.logger.Error("Failed to get database status", "error", err)
		response.InternalError(c, "failed to get database status")
		return
	}

	adminCtx, _ := middleware.GetAdminContext(c)
	if adminCtx != nil {
		h.logger.Info("Admin checked database status",
			"admin_id", adminCtx.UserID,
			"db_status", dbStatus.Status)
	}

	resp := dto.DatabaseStatusResponse{
		Status:       dbStatus.Status,
		Type:         dbStatus.Type,
		Latency:      dbStatus.Latency,
		MaxOpenConns: dbStatus.MaxOpenConns,
		OpenConns:    dbStatus.OpenConns,
		InUse:        dbStatus.InUse,
		Idle:         dbStatus.Idle,
		WaitCount:    dbStatus.WaitCount,
		WaitDuration: dbStatus.WaitDuration,
	}
	response.SuccessOK(c, resp)
}

// GetXraySystemStatus handles GET /api/v1/admin/system/xray
func (h *AdminHandler) GetXraySystemStatus(c *gin.Context) {
	xrayStatus, err := h.getXraySystemStatusQuery.Execute(c.Request.Context())
	if err != nil {
		h.logger.Error("Failed to get Xray system status", "error", err)
		response.InternalError(c, "failed to get Xray system status")
		return
	}

	adminCtx, _ := middleware.GetAdminContext(c)
	if adminCtx != nil {
		h.logger.Info("Admin checked Xray system status",
			"admin_id", adminCtx.UserID,
			"status", xrayStatus.Status,
			"running_instances", xrayStatus.RunningInstances)
	}

	// Map instances
	instances := make([]dto.XrayInstanceStatusInfo, 0, len(xrayStatus.Instances))
	for _, inst := range xrayStatus.Instances {
		instances = append(instances, dto.XrayInstanceStatusInfo{
			ID:      inst.ID.String(),
			NodeID:  inst.NodeID.String(),
			Status:  inst.Status,
			Healthy: inst.Healthy,
			PID:     inst.PID,
			Uptime:  inst.Uptime,
		})
	}

	resp := dto.XraySystemStatusResponse{
		Status:           xrayStatus.Status,
		TotalInstances:   xrayStatus.TotalInstances,
		RunningInstances: xrayStatus.RunningInstances,
		HealthyInstances: xrayStatus.HealthyInstances,
		Instances:        instances,
	}
	response.SuccessOK(c, resp)
}
