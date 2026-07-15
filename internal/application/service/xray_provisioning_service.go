package service

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/suproxy/backend/internal/domain/audit"
	"github.com/suproxy/backend/internal/domain/user"
	"github.com/suproxy/backend/internal/domain/xray"
	"github.com/suproxy/backend/internal/infrastructure/logger"
	xrayBinary "github.com/suproxy/backend/internal/infrastructure/xray/binary"
	xrayConfig "github.com/suproxy/backend/internal/infrastructure/xray/config"
	"github.com/suproxy/backend/internal/infrastructure/xray/runtime"
)

var (
	// ErrNoRunningInstances indicates no running Xray instances found
	ErrNoRunningInstances = errors.New("no running xray instances found")
	// ErrNoEnabledInbounds indicates no enabled inbounds found
	ErrNoEnabledInbounds = errors.New("no enabled inbounds found")
	// ErrInstanceUnhealthy indicates Xray instance is not healthy
	ErrInstanceUnhealthy = errors.New("xray instance is not healthy")
	// ErrConfigValidationFailed indicates generated config is invalid
	ErrConfigValidationFailed = errors.New("xray config validation failed")
	// ErrReloadTimeout indicates process reload timed out
	ErrReloadTimeout = errors.New("xray reload timeout")
)

// Error classification for retry logic
type ErrorClass int

const (
	// ErrorClassRetryable - temporary errors that can be retried
	ErrorClassRetryable ErrorClass = iota
	// ErrorClassNonRetryable - permanent errors that should not be retried
	ErrorClassNonRetryable
	// ErrorClassSkippable - not an error, operation skipped
	ErrorClassSkippable
)

// ProvisioningError wraps errors with retry classification
type ProvisioningError struct {
	Err   error
	Class ErrorClass
}

func (e *ProvisioningError) Error() string {
	return e.Err.Error()
}

func (e *ProvisioningError) Unwrap() error {
	return e.Err
}

// ClassifyError classifies errors for retry logic
func ClassifyError(err error) ErrorClass {
	if err == nil {
		return ErrorClassSkippable
	}

	// Skippable (not errors)
	if errors.Is(err, ErrNoRunningInstances) || errors.Is(err, ErrNoEnabledInbounds) {
		return ErrorClassSkippable
	}

	// Retryable errors
	if errors.Is(err, ErrReloadTimeout) || errors.Is(err, ErrInstanceUnhealthy) {
		return ErrorClassRetryable
	}

	// Non-retryable errors
	if errors.Is(err, ErrConfigValidationFailed) {
		return ErrorClassNonRetryable
	}

	// Default: retryable for transient failures
	return ErrorClassRetryable
}

const (
	// DefaultReloadTimeout is the default timeout for config reload
	DefaultReloadTimeout = 30 * time.Second
	// DefaultHealthCheckTimeout is the default timeout for health checks
	DefaultHealthCheckTimeout = 5 * time.Second
	// MaxBackupRetention is the maximum number of config backups to retain
	MaxBackupRetention = 10
)

// XrayProvisioningService handles automatic Xray provisioning for users
type XrayProvisioningService struct {
	xrayInstanceRepo xray.XrayInstanceRepository
	inboundRepo      xray.InboundRepository
	clientRepo       xray.ClientRepository
	realityRepo      xray.RealityConfigRepository
	auditRepo        audit.Repository
	configGenerator  xrayConfig.Generator
	configWriter     xrayConfig.Writer
	processManager   runtime.Manager
	binaryManager    xrayBinary.Manager
	logger           *logger.Logger

	// Configuration
	reloadTimeout      time.Duration
	healthCheckTimeout time.Duration
	maxBackupRetention int

	// Concurrency control - prevents parallel provisioning for same user
	provisioningLocks map[uuid.UUID]*sync.Mutex
	locksMutex        sync.Mutex
}

// NewXrayProvisioningService creates a new provisioning service
func NewXrayProvisioningService(
	xrayInstanceRepo xray.XrayInstanceRepository,
	inboundRepo xray.InboundRepository,
	clientRepo xray.ClientRepository,
	realityRepo xray.RealityConfigRepository,
	auditRepo audit.Repository,
	configGenerator xrayConfig.Generator,
	configWriter xrayConfig.Writer,
	processManager runtime.Manager,
	binaryManager xrayBinary.Manager,
	logger *logger.Logger,
) *XrayProvisioningService {
	return &XrayProvisioningService{
		xrayInstanceRepo:   xrayInstanceRepo,
		inboundRepo:        inboundRepo,
		clientRepo:         clientRepo,
		realityRepo:        realityRepo,
		auditRepo:          auditRepo,
		configGenerator:    configGenerator,
		configWriter:       configWriter,
		processManager:     processManager,
		binaryManager:      binaryManager,
		logger:             logger,
		reloadTimeout:      DefaultReloadTimeout,
		healthCheckTimeout: DefaultHealthCheckTimeout,
		maxBackupRetention: MaxBackupRetention,
		provisioningLocks:  make(map[uuid.UUID]*sync.Mutex),
	}
}

// acquireUserLock acquires a lock for user provisioning to prevent race conditions
func (s *XrayProvisioningService) acquireUserLock(userID uuid.UUID) *sync.Mutex {
	s.locksMutex.Lock()
	defer s.locksMutex.Unlock()

	if lock, exists := s.provisioningLocks[userID]; exists {
		return lock
	}

	lock := &sync.Mutex{}
	s.provisioningLocks[userID] = lock
	return lock
}

// releaseUserLock releases and cleans up user lock
func (s *XrayProvisioningService) releaseUserLock(userID uuid.UUID) {
	s.locksMutex.Lock()
	defer s.locksMutex.Unlock()

	// Clean up lock after use to prevent memory leak
	delete(s.provisioningLocks, userID)
}

// ProvisionUserToXray provisions a user to Xray after registration.
// This function is idempotent - calling it multiple times for the same user will not create duplicates.
// It follows a compensating transaction pattern: if config reload fails, the client is rolled back.
// Thread-safe: uses per-user locking to prevent race conditions during parallel registrations.
func (s *XrayProvisioningService) ProvisionUserToXray(ctx context.Context, newUser *user.User, ipAddress, userAgent string) error {
	// Acquire user-specific lock to prevent race conditions
	userLock := s.acquireUserLock(newUser.ID)
	userLock.Lock()
	defer func() {
		userLock.Unlock()
		s.releaseUserLock(newUser.ID)
	}()

	startTime := time.Now()
	s.logger.Info("Starting Xray provisioning for new user",
		"user_id", newUser.ID,
		"email", newUser.Email.String(),
		"ip_address", ipAddress,
		"operation", "provision_start")

	// Audit: User created
	s.auditEvent(ctx, newUser.ID, "user", newUser.ID, "user_created", ipAddress, userAgent, nil)

	// Step 1: Check if user already has a client (idempotency check)
	existingClients, err := s.clientRepo.FindByUserID(ctx, newUser.ID)
	if err != nil {
		s.logger.Error("Failed to check existing clients",
			"error", err,
			"user_id", newUser.ID,
			"operation", "idempotency_check",
			"duration_ms", time.Since(startTime).Milliseconds())
		return fmt.Errorf("failed to check existing clients: %w", err)
	}

	if len(existingClients) > 0 {
		s.logger.Info("User already has Xray client, skipping provisioning",
			"user_id", newUser.ID,
			"client_count", len(existingClients),
			"operation", "idempotency_skip",
			"duration_ms", time.Since(startTime).Milliseconds())
		return nil // Idempotent: already provisioned
	}

	// Step 2: Find active Xray instance
	instance, err := s.findRunningInstance(ctx, newUser.ID)
	if err != nil {
		if errors.Is(err, ErrNoRunningInstances) {
			s.logger.Warn("No running Xray instances found, skipping provisioning",
				"user_id", newUser.ID)
			return nil // Not a failure, just skip
		}
		return err
	}

	// Step 2.1: Health Check - ensure instance is healthy before proceeding
	if err := s.performHealthCheck(ctx, instance.ID, newUser.ID); err != nil {
		s.logger.Error("Health check failed, aborting provisioning",
			"error", err,
			"instance_id", instance.ID,
			"user_id", newUser.ID,
			"operation", "health_check_provision")

		s.auditEvent(ctx, newUser.ID, "xray_instance", instance.ID, "xray_health_check_failed", ipAddress, userAgent, map[string]interface{}{
			"error": err.Error(),
			"stage": "provision",
		})

		return &ProvisioningError{Err: err, Class: ErrorClassRetryable}
	}

	// Step 3: Find enabled inbound (prefer VLESS)
	targetInbound, err := s.findTargetInbound(ctx, instance.ID, newUser.ID)
	if err != nil {
		if errors.Is(err, ErrNoEnabledInbounds) {
			s.logger.Warn("No enabled inbounds found, skipping provisioning",
				"user_id", newUser.ID,
				"instance_id", instance.ID)
			return nil // Not a failure, just skip
		}
		return err
	}

	// Step 4: Create client entity
	client, clientUUID, err := s.createClientEntity(targetInbound, newUser)
	if err != nil {
		return err
	}

	// Step 5: Save client to database
	if err := s.clientRepo.Create(ctx, client); err != nil {
		s.logger.Error("Failed to save client",
			"error", err,
			"user_id", newUser.ID,
			"inbound_id", targetInbound.ID)
		return fmt.Errorf("failed to save client: %w", err)
	}

	s.logger.Info("Client created successfully",
		"client_id", client.ID,
		"user_id", newUser.ID,
		"instance_id", instance.ID,
		"inbound_id", targetInbound.ID,
		"uuid", clientUUID)

	// Audit: Client created
	s.auditEvent(ctx, newUser.ID, "xray_client", client.ID, "xray_client_created", ipAddress, userAgent, map[string]interface{}{
		"inbound_id":  targetInbound.ID.String(),
		"instance_id": instance.ID.String(),
		"uuid":        clientUUID,
	})

	// Step 6: Regenerate and reload Xray config (with client rollback on failure)
	if err := s.regenerateAndReloadWithRollback(ctx, instance.ID, client.ID, newUser.ID, ipAddress, userAgent); err != nil {
		s.logger.Error("Failed to regenerate and reload Xray config, rolling back client",
			"error", err,
			"user_id", newUser.ID,
			"client_id", client.ID,
			"instance_id", instance.ID)

		// Compensating action: Delete the created client
		if rollbackErr := s.rollbackClient(ctx, client.ID, newUser.ID, ipAddress, userAgent); rollbackErr != nil {
			s.logger.Error("Failed to rollback client after config reload failure",
				"error", rollbackErr,
				"user_id", newUser.ID,
				"client_id", client.ID)
			// System is now in inconsistent state - client exists but not in config
			// This will be logged and can be manually fixed or retried
			return fmt.Errorf("config reload failed and client rollback failed: reload_err=%w, rollback_err=%v", err, rollbackErr)
		}

		s.logger.Info("Client successfully rolled back after config reload failure",
			"user_id", newUser.ID,
			"client_id", client.ID)

		// Return error to caller, but user is still created (client was rolled back)
		return fmt.Errorf("config reload failed, client rolled back: %w", err)
	}

	s.logger.Info("Xray provisioning completed successfully",
		"user_id", newUser.ID,
		"client_id", client.ID,
		"instance_id", instance.ID,
		"operation", "provision_success",
		"total_duration_ms", time.Since(startTime).Milliseconds())
	return nil
}

// performHealthCheck checks if Xray instance is healthy before provisioning
func (s *XrayProvisioningService) performHealthCheck(ctx context.Context, instanceID, userID uuid.UUID) error {
	ctx, cancel := context.WithTimeout(ctx, s.healthCheckTimeout)
	defer cancel()

	// Check if process is running
	isRunning, err := s.processManager.IsRunning(ctx, instanceID)
	if err != nil {
		s.logger.Error("Failed to check instance status",
			"error", err,
			"instance_id", instanceID,
			"user_id", userID,
			"operation", "health_check")
		return fmt.Errorf("failed to check instance status: %w", err)
	}

	if !isRunning {
		s.logger.Warn("Instance is not running",
			"instance_id", instanceID,
			"user_id", userID,
			"operation", "health_check")
		return ErrInstanceUnhealthy
	}

	// Get process status for additional health indicators
	status, err := s.processManager.Status(ctx, instanceID)
	if err != nil {
		s.logger.Warn("Failed to get process status, assuming healthy",
			"error", err,
			"instance_id", instanceID,
			"user_id", userID,
			"operation", "health_check")
		// Don't fail if we can't get status but process is running
		return nil
	}

	s.logger.Debug("Health check passed",
		"instance_id", instanceID,
		"user_id", userID,
		"operation", "health_check",
		"status", status)

	return nil
}

// backupWithRetention creates backup and manages retention policy
func (s *XrayProvisioningService) backupWithRetention(ctx context.Context, instanceID, userID uuid.UUID) error {
	// Create new backup
	if err := s.configWriter.Backup(ctx, instanceID); err != nil {
		return fmt.Errorf("failed to create backup: %w", err)
	}

	// List all backups
	backups, err := s.configWriter.ListBackups(ctx, instanceID)
	if err != nil {
		s.logger.Warn("Failed to list backups for retention cleanup",
			"error", err,
			"instance_id", instanceID,
			"user_id", userID)
		return nil // Don't fail provisioning due to cleanup failure
	}

	// If we have more backups than retention limit, delete oldest
	if len(backups) > s.maxBackupRetention {
		excessCount := len(backups) - s.maxBackupRetention
		s.logger.Info("Cleaning up old backups",
			"instance_id", instanceID,
			"user_id", userID,
			"total_backups", len(backups),
			"excess_count", excessCount,
			"retention_limit", s.maxBackupRetention)

		// Backups are sorted by timestamp, delete oldest ones
		for i := 0; i < excessCount; i++ {
			oldBackup := backups[i]

			// Actually delete the backup file
			if err := s.configWriter.DeleteBackup(ctx, instanceID, oldBackup.Timestamp.Unix()); err != nil {
				s.logger.Warn("Failed to delete old backup",
					"error", err,
					"instance_id", instanceID,
					"backup_timestamp", oldBackup.Timestamp,
					"backup_path", oldBackup.Path)
				// Continue with other backups even if one fails
				continue
			}

			s.logger.Debug("Deleted old backup",
				"instance_id", instanceID,
				"backup_timestamp", oldBackup.Timestamp,
				"backup_path", oldBackup.Path)
		}
	}

	return nil
}

// findRunningInstance finds an active Xray instance
func (s *XrayProvisioningService) findRunningInstance(ctx context.Context, userID uuid.UUID) (*xray.XrayInstance, error) {
	runningInstances, err := s.xrayInstanceRepo.FindRunning(ctx)
	if err != nil {
		s.logger.Error("Failed to find running Xray instances",
			"error", err,
			"user_id", userID)
		return nil, fmt.Errorf("failed to find running Xray instances: %w", err)
	}

	if len(runningInstances) == 0 {
		return nil, ErrNoRunningInstances
	}

	// Use first running instance (future: load balancing)
	instance := runningInstances[0]
	s.logger.Debug("Found running Xray instance",
		"instance_id", instance.ID,
		"user_id", userID)

	return instance, nil
}

// findTargetInbound finds an enabled inbound, preferring VLESS
func (s *XrayProvisioningService) findTargetInbound(ctx context.Context, instanceID, userID uuid.UUID) (*xray.Inbound, error) {
	inbounds, err := s.inboundRepo.FindEnabledByInstanceID(ctx, instanceID)
	if err != nil {
		s.logger.Error("Failed to find inbounds",
			"error", err,
			"instance_id", instanceID,
			"user_id", userID)
		return nil, fmt.Errorf("failed to find inbounds: %w", err)
	}

	if len(inbounds) == 0 {
		return nil, ErrNoEnabledInbounds
	}

	// Prefer VLESS inbound
	for _, inbound := range inbounds {
		if inbound.Protocol == xray.ProtocolVLESS {
			s.logger.Debug("Selected VLESS inbound",
				"inbound_id", inbound.ID,
				"instance_id", instanceID,
				"user_id", userID)
			return inbound, nil
		}
	}

	// Fallback to first available inbound
	inbound := inbounds[0]
	s.logger.Debug("Selected fallback inbound",
		"inbound_id", inbound.ID,
		"protocol", inbound.Protocol,
		"instance_id", instanceID,
		"user_id", userID)

	return inbound, nil
}

// createClientEntity creates a new Xray client entity
func (s *XrayProvisioningService) createClientEntity(inbound *xray.Inbound, newUser *user.User) (*xray.Client, string, error) {
	clientUUID := uuid.New().String()

	// Determine flow based on protocol
	flow := ""
	if inbound.Protocol == xray.ProtocolVLESS {
		flow = "xtls-rprx-vision" // Default flow for VLESS + Reality
	}

	client, err := xray.NewClient(
		inbound.ID,
		newUser.ID,
		clientUUID,
		flow,
		newUser.Email.String(),
	)
	if err != nil {
		s.logger.Error("Failed to create client entity",
			"error", err,
			"user_id", newUser.ID,
			"inbound_id", inbound.ID)
		return nil, "", fmt.Errorf("failed to create client entity: %w", err)
	}

	return client, clientUUID, nil
}

// rollbackClient deletes a client from the database (compensating transaction)
func (s *XrayProvisioningService) rollbackClient(ctx context.Context, clientID, userID uuid.UUID, ipAddress, userAgent string) error {
	s.logger.Warn("Rolling back client",
		"client_id", clientID,
		"user_id", userID)

	// Audit: Client rollback started
	s.auditEvent(ctx, userID, "xray_client", clientID, "xray_client_rollback_started", ipAddress, userAgent, nil)

	if err := s.clientRepo.Delete(ctx, clientID); err != nil {
		s.logger.Error("Failed to delete client during rollback",
			"error", err,
			"client_id", clientID,
			"user_id", userID)

		// Audit: Rollback failed
		s.auditEvent(ctx, userID, "xray_client", clientID, "xray_client_rollback_failed", ipAddress, userAgent, map[string]interface{}{
			"error": err.Error(),
		})

		return fmt.Errorf("failed to delete client: %w", err)
	}

	// Audit: Rollback successful
	s.auditEvent(ctx, userID, "xray_client", clientID, "xray_client_rollback_success", ipAddress, userAgent, nil)

	s.logger.Info("Client rollback successful",
		"client_id", clientID,
		"user_id", userID)

	return nil
}

// regenerateAndReloadWithRollback regenerates config and reloads, with client rollback on failure
func (s *XrayProvisioningService) regenerateAndReloadWithRollback(ctx context.Context, instanceID, clientID, userID uuid.UUID, ipAddress, userAgent string) error {
	if err := s.RegenerateAndReload(ctx, instanceID, userID, ipAddress, userAgent); err != nil {
		return err
	}
	return nil
}

// auditEvent is a helper to create audit logs with consistent structure
func (s *XrayProvisioningService) auditEvent(ctx context.Context, userID uuid.UUID, entityType string, entityID uuid.UUID, event string, ipAddress, userAgent string, metadata map[string]interface{}) {
	auditLog := audit.NewLog(userID, audit.ActionCreate, entityType, entityID, ipAddress, userAgent)
	auditLog.AddMetadata("event", event)

	// Add additional metadata
	for key, value := range metadata {
		auditLog.AddMetadata(key, value)
	}

	// Fire and forget - don't block on audit logging
	if err := s.auditRepo.Create(ctx, auditLog); err != nil {
		s.logger.Warn("Failed to create audit log",
			"error", err,
			"event", event,
			"user_id", userID)
	}
}

// RegenerateAndReload regenerates Xray config and reloads the process with full production safeguards
func (s *XrayProvisioningService) RegenerateAndReload(ctx context.Context, instanceID, userID uuid.UUID, ipAddress, userAgent string) error {
	startTime := time.Now()
	s.logger.Info("Regenerating Xray configuration",
		"instance_id", instanceID,
		"user_id", userID,
		"operation", "config_reload")

	// Step 1: Health Check - ensure instance is healthy before proceeding
	if err := s.performHealthCheck(ctx, instanceID, userID); err != nil {
		s.logger.Error("Health check failed, aborting reload",
			"error", err,
			"instance_id", instanceID,
			"user_id", userID,
			"operation", "health_check_reload")

		s.auditEvent(ctx, userID, "xray_instance", instanceID, "xray_health_check_failed", ipAddress, userAgent, map[string]interface{}{
			"error": err.Error(),
			"stage": "reload",
		})

		return &ProvisioningError{Err: err, Class: ErrorClassRetryable}
	}

	// Step 2: Backup current config with retention management
	backupErr := s.backupWithRetention(ctx, instanceID, userID)
	if backupErr != nil {
		s.logger.Warn("Failed to backup config, continuing anyway",
			"error", backupErr,
			"instance_id", instanceID,
			"user_id", userID,
			"operation", "backup")
	}

	// Step 3: Generate new config
	configJSON, err := s.configGenerator.GenerateJSON(ctx, instanceID)
	if err != nil {
		s.logger.Error("Failed to generate config",
			"error", err,
			"instance_id", instanceID,
			"user_id", userID,
			"operation", "config_generation",
			"duration_ms", time.Since(startTime).Milliseconds())

		s.auditEvent(ctx, userID, "xray_config", instanceID, "xray_config_generation_failed", ipAddress, userAgent, map[string]interface{}{
			"error":       err.Error(),
			"duration_ms": time.Since(startTime).Milliseconds(),
		})

		return &ProvisioningError{Err: fmt.Errorf("failed to generate config: %w", err), Class: ErrorClassRetryable}
	}

	s.logger.Info("Config generated successfully",
		"instance_id", instanceID,
		"user_id", userID,
		"config_size", len(configJSON),
		"operation", "config_generation",
		"duration_ms", time.Since(startTime).Milliseconds())

	s.auditEvent(ctx, userID, "xray_config", instanceID, "xray_config_generated", ipAddress, userAgent, map[string]interface{}{
		"config_size": len(configJSON),
		"duration_ms": time.Since(startTime).Milliseconds(),
	})

	// Step 4: Write config atomically to temporary file
	if err := s.configWriter.Write(ctx, instanceID, configJSON); err != nil {
		s.logger.Error("Failed to write config",
			"error", err,
			"instance_id", instanceID,
			"user_id", userID,
			"operation", "config_write",
			"duration_ms", time.Since(startTime).Milliseconds())

		s.auditEvent(ctx, userID, "xray_config", instanceID, "xray_config_write_failed", ipAddress, userAgent, map[string]interface{}{
			"error":       err.Error(),
			"duration_ms": time.Since(startTime).Milliseconds(),
		})

		return &ProvisioningError{Err: fmt.Errorf("failed to write config: %w", err), Class: ErrorClassRetryable}
	}

	s.logger.Info("Config written successfully",
		"instance_id", instanceID,
		"user_id", userID,
		"operation", "config_write")

	s.auditEvent(ctx, userID, "xray_config", instanceID, "xray_config_written", ipAddress, userAgent, nil)

	// Step 5: CRITICAL - Validate config using Xray binary BEFORE reload
	configPath := s.configWriter.GetPath(instanceID)
	if err := s.binaryManager.ValidateConfig(ctx, configPath); err != nil {
		s.logger.Error("Config validation failed using Xray binary",
			"error", err,
			"instance_id", instanceID,
			"user_id", userID,
			"config_path", configPath,
			"operation", "config_validation",
			"duration_ms", time.Since(startTime).Milliseconds())

		s.auditEvent(ctx, userID, "xray_config", instanceID, "xray_config_validation_failed", ipAddress, userAgent, map[string]interface{}{
			"error":       err.Error(),
			"config_path": configPath,
			"duration_ms": time.Since(startTime).Milliseconds(),
		})

		// Config is invalid - rollback to previous version
		if backupErr == nil {
			s.logger.Warn("Attempting to rollback to previous valid config",
				"instance_id", instanceID,
				"user_id", userID,
				"operation", "validation_rollback")

			if restoreErr := s.attemptConfigRollback(ctx, instanceID, userID, ipAddress, userAgent); restoreErr != nil {
				s.logger.Error("Failed to rollback after validation failure",
					"error", restoreErr,
					"instance_id", instanceID,
					"user_id", userID)
				return &ProvisioningError{
					Err:   fmt.Errorf("config validation failed and rollback failed: validation_err=%w, rollback_err=%v", err, restoreErr),
					Class: ErrorClassNonRetryable,
				}
			}
		}

		return &ProvisioningError{Err: ErrConfigValidationFailed, Class: ErrorClassNonRetryable}
	}

	s.logger.Info("Config validation successful",
		"instance_id", instanceID,
		"user_id", userID,
		"config_path", configPath,
		"operation", "config_validation")

	s.auditEvent(ctx, userID, "xray_config", instanceID, "xray_config_validated", ipAddress, userAgent, map[string]interface{}{
		"config_path": configPath,
	})

	// Step 6: Reload Xray process with timeout (hot reload with SIGHUP)
	reloadCtx, reloadCancel := context.WithTimeout(ctx, s.reloadTimeout)
	defer reloadCancel()

	reloadStartTime := time.Now()
	if err := s.processManager.Reload(reloadCtx, instanceID); err != nil {
		reloadDuration := time.Since(reloadStartTime)

		// Check if timeout occurred
		if errors.Is(err, context.DeadlineExceeded) {
			s.logger.Error("Reload timeout exceeded",
				"timeout", s.reloadTimeout,
				"instance_id", instanceID,
				"user_id", userID,
				"operation", "reload_timeout",
				"duration_ms", reloadDuration.Milliseconds())

			s.auditEvent(ctx, userID, "xray_process", instanceID, "xray_reload_timeout", ipAddress, userAgent, map[string]interface{}{
				"timeout_ms":  s.reloadTimeout.Milliseconds(),
				"duration_ms": reloadDuration.Milliseconds(),
			})

			// Attempt rollback
			if backupErr == nil {
				s.logger.Warn("Attempting config rollback after timeout",
					"instance_id", instanceID,
					"user_id", userID,
					"operation", "timeout_rollback")

				if restoreErr := s.attemptConfigRollback(ctx, instanceID, userID, ipAddress, userAgent); restoreErr != nil {
					s.logger.Error("Config rollback failed after timeout",
						"error", restoreErr,
						"instance_id", instanceID,
						"user_id", userID)
					return &ProvisioningError{
						Err:   fmt.Errorf("reload timeout and rollback failed: %w, rollback_err=%v", ErrReloadTimeout, restoreErr),
						Class: ErrorClassRetryable,
					}
				}
			}

			return &ProvisioningError{Err: ErrReloadTimeout, Class: ErrorClassRetryable}
		}

		// Other reload errors
		s.logger.Error("Failed to reload Xray process",
			"error", err,
			"instance_id", instanceID,
			"user_id", userID,
			"operation", "reload",
			"duration_ms", reloadDuration.Milliseconds())

		s.auditEvent(ctx, userID, "xray_process", instanceID, "xray_reload_failed", ipAddress, userAgent, map[string]interface{}{
			"error":       err.Error(),
			"duration_ms": reloadDuration.Milliseconds(),
		})

		// Attempt rollback
		if backupErr == nil {
			s.logger.Warn("Attempting config rollback after reload failure",
				"instance_id", instanceID,
				"user_id", userID,
				"operation", "reload_failure_rollback")

			if restoreErr := s.attemptConfigRollback(ctx, instanceID, userID, ipAddress, userAgent); restoreErr != nil {
				s.logger.Error("Config rollback failed",
					"error", restoreErr,
					"instance_id", instanceID,
					"user_id", userID)
				return &ProvisioningError{
					Err:   fmt.Errorf("reload failed and rollback failed: reload_err=%w, rollback_err=%v", err, restoreErr),
					Class: ErrorClassRetryable,
				}
			}

			s.logger.Info("Config rollback successful after reload failure",
				"instance_id", instanceID,
				"user_id", userID)
			return &ProvisioningError{Err: fmt.Errorf("reload failed, rolled back to previous config: %w", err), Class: ErrorClassRetryable}
		}

		return &ProvisioningError{Err: fmt.Errorf("failed to reload Xray process: %w", err), Class: ErrorClassRetryable}
	}

	totalDuration := time.Since(startTime)
	s.logger.Info("Xray process reloaded successfully",
		"instance_id", instanceID,
		"user_id", userID,
		"operation", "reload_success",
		"reload_duration_ms", time.Since(reloadStartTime).Milliseconds(),
		"total_duration_ms", totalDuration.Milliseconds())

	s.auditEvent(ctx, userID, "xray_process", instanceID, "xray_reloaded_success", ipAddress, userAgent, map[string]interface{}{
		"reload_duration_ms": time.Since(reloadStartTime).Milliseconds(),
		"total_duration_ms":  totalDuration.Milliseconds(),
	})

	return nil
}

// attemptConfigRollback attempts to restore previous config and reload (renamed from attemptRollback for clarity)
func (s *XrayProvisioningService) attemptConfigRollback(ctx context.Context, instanceID, userID uuid.UUID, ipAddress, userAgent string) error {
	s.logger.Warn("Starting config rollback",
		"instance_id", instanceID,
		"user_id", userID)

	// Audit: Rollback started
	s.auditEvent(ctx, userID, "xray_config", instanceID, "xray_config_rollback_started", ipAddress, userAgent, nil)

	// Find latest backup
	backups, err := s.configWriter.ListBackups(ctx, instanceID)
	if err != nil || len(backups) == 0 {
		s.logger.Error("No backups available for rollback",
			"error", err,
			"instance_id", instanceID,
			"user_id", userID)

		// Audit: Rollback failed - no backup
		s.auditEvent(ctx, userID, "xray_config", instanceID, "xray_config_rollback_failed", ipAddress, userAgent, map[string]interface{}{
			"error":  "no backups available",
			"reason": "no_backup",
		})

		return fmt.Errorf("no backups available: %w", err)
	}

	// Use most recent backup
	latestBackup := backups[len(backups)-1]

	s.logger.Debug("Restoring from backup",
		"instance_id", instanceID,
		"user_id", userID,
		"backup_timestamp", latestBackup.Timestamp)

	// Restore backup
	if err := s.configWriter.Restore(ctx, instanceID, latestBackup.Timestamp); err != nil {
		s.logger.Error("Failed to restore backup",
			"error", err,
			"instance_id", instanceID,
			"user_id", userID)

		// Audit: Rollback failed - restore error
		s.auditEvent(ctx, userID, "xray_config", instanceID, "xray_config_rollback_failed", ipAddress, userAgent, map[string]interface{}{
			"error":  err.Error(),
			"reason": "restore_failed",
		})

		return fmt.Errorf("failed to restore backup: %w", err)
	}

	// Try to reload with restored config
	if err := s.processManager.Reload(ctx, instanceID); err != nil {
		s.logger.Error("Failed to reload after rollback",
			"error", err,
			"instance_id", instanceID,
			"user_id", userID)

		// Audit: Rollback failed - reload error
		s.auditEvent(ctx, userID, "xray_config", instanceID, "xray_config_rollback_failed", ipAddress, userAgent, map[string]interface{}{
			"error":  err.Error(),
			"reason": "reload_failed",
		})

		// CRITICAL: System is now in unknown state
		// Config is restored but process didn't reload
		// Manual intervention may be required
		return fmt.Errorf("failed to reload after rollback: %w", err)
	}

	// Audit: Rollback successful
	s.auditEvent(ctx, userID, "xray_config", instanceID, "xray_config_rollback_success", ipAddress, userAgent, nil)

	s.logger.Info("Config rollback completed successfully",
		"instance_id", instanceID,
		"user_id", userID)

	return nil
}
