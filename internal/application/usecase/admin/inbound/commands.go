package inbound

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/suproxy/backend/internal/application/service"
	"github.com/suproxy/backend/internal/domain/audit"
	"github.com/suproxy/backend/internal/domain/xray"
)

// CreateInboundCommand handles creating a new inbound
type CreateInboundCommand struct {
	inboundRepo         xray.InboundRepository
	instanceRepo        xray.XrayInstanceRepository
	provisioningService *service.XrayProvisioningService
	auditRepo           audit.Repository
}

func NewCreateInboundCommand(
	inboundRepo xray.InboundRepository,
	instanceRepo xray.XrayInstanceRepository,
	provisioningService *service.XrayProvisioningService,
	auditRepo audit.Repository,
) *CreateInboundCommand {
	return &CreateInboundCommand{
		inboundRepo:         inboundRepo,
		instanceRepo:        instanceRepo,
		provisioningService: provisioningService,
		auditRepo:           auditRepo,
	}
}

func (c *CreateInboundCommand) Execute(
	ctx context.Context,
	instanceID uuid.UUID,
	protocol xray.InboundProtocol,
	port int,
	transport xray.TransportType,
	security xray.SecurityType,
	adminID uuid.UUID,
	ip, userAgent string,
) (*xray.Inbound, error) {
	// Verify instance exists
	instance, err := c.instanceRepo.FindByID(ctx, instanceID)
	if err != nil {
		return nil, fmt.Errorf("failed to find instance: %w", err)
	}

	// Create inbound entity
	inbound, err := xray.NewInbound(instanceID, protocol, port, transport, security)
	if err != nil {
		return nil, fmt.Errorf("failed to create inbound entity: %w", err)
	}

	// Save inbound
	if err := c.inboundRepo.Create(ctx, inbound); err != nil {
		return nil, fmt.Errorf("failed to save inbound: %w", err)
	}

	// Audit: Inbound created
	auditLog := audit.NewLog(adminID, audit.ActionCreate, "xray_inbound", inbound.ID, ip, userAgent)
	auditLog.AddMetadata("event", "inbound_created")
	auditLog.AddMetadata("instance_id", instanceID.String())
	auditLog.AddMetadata("protocol", string(protocol))
	auditLog.AddMetadata("port", port)
	if err := c.auditRepo.Create(ctx, auditLog); err != nil {
		// Log audit failure but don't block operation
		// Audit is important but not critical for functionality
		_ = err // errcheck: acknowledged
	}

	// Reload config (REUSES XrayProvisioningService)
	if err := c.provisioningService.RegenerateAndReload(ctx, instance.ID, adminID, ip, userAgent); err != nil {
		// Rollback: Delete the created inbound
		if delErr := c.inboundRepo.Delete(ctx, inbound.ID); delErr != nil {
			// Critical: Failed to rollback, log and return compound error
			return nil, fmt.Errorf("config reload failed AND rollback delete failed: reload=%w, rollback=%v", err, delErr)
		}
		
		auditLog := audit.NewLog(adminID, audit.ActionDelete, "xray_inbound", inbound.ID, ip, userAgent)
		auditLog.AddMetadata("event", "inbound_rollback_after_reload_failed")
		if auditErr := c.auditRepo.Create(ctx, auditLog); auditErr != nil {
			_ = auditErr // errcheck: acknowledged, audit failure in rollback is not critical
		}
		
		return nil, fmt.Errorf("config reload failed, inbound rolled back: %w", err)
	}

	return inbound, nil
}

// UpdateInboundCommand handles updating an inbound
type UpdateInboundCommand struct {
	inboundRepo         xray.InboundRepository
	provisioningService *service.XrayProvisioningService
	auditRepo           audit.Repository
}

func NewUpdateInboundCommand(
	inboundRepo xray.InboundRepository,
	provisioningService *service.XrayProvisioningService,
	auditRepo audit.Repository,
) *UpdateInboundCommand {
	return &UpdateInboundCommand{
		inboundRepo:         inboundRepo,
		provisioningService: provisioningService,
		auditRepo:           auditRepo,
	}
}

func (c *UpdateInboundCommand) Execute(
	ctx context.Context,
	inboundID uuid.UUID,
	port *int,
	transport *xray.TransportType,
	security *xray.SecurityType,
	adminID uuid.UUID,
	ip, userAgent string,
) (*xray.Inbound, error) {
	// Find inbound
	inbound, err := c.inboundRepo.FindByID(ctx, inboundID)
	if err != nil {
		return nil, fmt.Errorf("failed to find inbound: %w", err)
	}

	oldPort := inbound.Port
	oldTransport := inbound.Transport
	oldSecurity := inbound.Security

	// Apply updates
	changed := false
	if port != nil && *port != inbound.Port {
		if err := inbound.ChangePort(*port); err != nil {
			return nil, fmt.Errorf("failed to change port: %w", err)
		}
		changed = true
	}

	if transport != nil && *transport != inbound.Transport {
		if err := inbound.UpdateTransport(*transport); err != nil {
			return nil, fmt.Errorf("failed to update transport: %w", err)
		}
		changed = true
	}

	if security != nil && *security != inbound.Security {
		if err := inbound.UpdateSecurity(*security); err != nil {
			return nil, fmt.Errorf("failed to update security: %w", err)
		}
		changed = true
	}

	if !changed {
		return inbound, nil // No changes, return early
	}

	// Save inbound
	if err := c.inboundRepo.Update(ctx, inbound); err != nil {
		return nil, fmt.Errorf("failed to update inbound: %w", err)
	}

	// Audit: Inbound updated
	auditLog := audit.NewLog(adminID, audit.ActionUpdate, "xray_inbound", inboundID, ip, userAgent)
	auditLog.AddMetadata("event", "inbound_updated")
	if port != nil {
		auditLog.AddMetadata("old_port", oldPort)
		auditLog.AddMetadata("new_port", *port)
	}
	if transport != nil {
		auditLog.AddMetadata("old_transport", string(oldTransport))
		auditLog.AddMetadata("new_transport", string(*transport))
	}
	if security != nil {
		auditLog.AddMetadata("old_security", string(oldSecurity))
		auditLog.AddMetadata("new_security", string(*security))
	}
	if err := c.auditRepo.Create(ctx, auditLog); err != nil {
		_ = err // errcheck: acknowledged, audit failure is not critical
	}

	// Reload config (REUSES XrayProvisioningService)
	if err := c.provisioningService.RegenerateAndReload(ctx, inbound.XrayInstanceID, adminID, ip, userAgent); err != nil {
		// Rollback: Restore old values
		if port != nil {
			if chgErr := inbound.ChangePort(oldPort); chgErr != nil {
				return nil, fmt.Errorf("config reload failed AND rollback ChangePort failed: reload=%w, rollback=%v", err, chgErr)
			}
		}
		if transport != nil {
			if updErr := inbound.UpdateTransport(oldTransport); updErr != nil {
				return nil, fmt.Errorf("config reload failed AND rollback UpdateTransport failed: reload=%w, rollback=%v", err, updErr)
			}
		}
		if security != nil {
			if secErr := inbound.UpdateSecurity(oldSecurity); secErr != nil {
				return nil, fmt.Errorf("config reload failed AND rollback UpdateSecurity failed: reload=%w, rollback=%v", err, secErr)
			}
		}
		if updErr := c.inboundRepo.Update(ctx, inbound); updErr != nil {
			return nil, fmt.Errorf("config reload failed AND rollback Update failed: reload=%w, rollback=%v", err, updErr)
		}
		
		auditLog := audit.NewLog(adminID, audit.ActionUpdate, "xray_inbound", inboundID, ip, userAgent)
		auditLog.AddMetadata("event", "inbound_rollback_after_reload_failed")
		if auditErr := c.auditRepo.Create(ctx, auditLog); auditErr != nil {
			_ = auditErr // errcheck: acknowledged, audit failure in rollback is not critical
		}
		
		return nil, fmt.Errorf("config reload failed, inbound rolled back: %w", err)
	}

	return inbound, nil
}

// DeleteInboundCommand handles deleting an inbound
type DeleteInboundCommand struct {
	inboundRepo         xray.InboundRepository
	clientRepo          xray.ClientRepository
	provisioningService *service.XrayProvisioningService
	auditRepo           audit.Repository
}

func NewDeleteInboundCommand(
	inboundRepo xray.InboundRepository,
	clientRepo xray.ClientRepository,
	provisioningService *service.XrayProvisioningService,
	auditRepo audit.Repository,
) *DeleteInboundCommand {
	return &DeleteInboundCommand{
		inboundRepo:         inboundRepo,
		clientRepo:          clientRepo,
		provisioningService: provisioningService,
		auditRepo:           auditRepo,
	}
}

func (c *DeleteInboundCommand) Execute(ctx context.Context, inboundID, adminID uuid.UUID, ip, userAgent string) error {
	// Find inbound
	inbound, err := c.inboundRepo.FindByID(ctx, inboundID)
	if err != nil {
		return fmt.Errorf("failed to find inbound: %w", err)
	}

	// Check if inbound has clients
	clients, err := c.clientRepo.FindByInboundID(ctx, inboundID)
	if err != nil {
		return fmt.Errorf("failed to check inbound clients: %w", err)
	}

	if len(clients) > 0 {
		return fmt.Errorf("cannot delete inbound with %d active clients", len(clients))
	}

	// Delete inbound
	if err := c.inboundRepo.Delete(ctx, inboundID); err != nil {
		return fmt.Errorf("failed to delete inbound: %w", err)
	}

	// Audit: Inbound deleted
	auditLog := audit.NewLog(adminID, audit.ActionDelete, "xray_inbound", inboundID, ip, userAgent)
	auditLog.AddMetadata("event", "inbound_deleted")
	auditLog.AddMetadata("instance_id", inbound.XrayInstanceID.String())
	if err := c.auditRepo.Create(ctx, auditLog); err != nil {
		_ = err // errcheck: acknowledged, audit failure is not critical for delete operation
	}

	// Reload config (REUSES XrayProvisioningService)
	if err := c.provisioningService.RegenerateAndReload(ctx, inbound.XrayInstanceID, adminID, ip, userAgent); err != nil {
		return fmt.Errorf("inbound deleted but config reload failed: %w", err)
	}

	return nil
}

// EnableInboundCommand handles enabling an inbound
type EnableInboundCommand struct {
	inboundRepo         xray.InboundRepository
	provisioningService *service.XrayProvisioningService
	auditRepo           audit.Repository
}

func NewEnableInboundCommand(
	inboundRepo xray.InboundRepository,
	provisioningService *service.XrayProvisioningService,
	auditRepo audit.Repository,
) *EnableInboundCommand {
	return &EnableInboundCommand{
		inboundRepo:         inboundRepo,
		provisioningService: provisioningService,
		auditRepo:           auditRepo,
	}
}

func (c *EnableInboundCommand) Execute(ctx context.Context, inboundID, adminID uuid.UUID, ip, userAgent string) error {
	// Find inbound
	inbound, err := c.inboundRepo.FindByID(ctx, inboundID)
	if err != nil {
		return fmt.Errorf("failed to find inbound: %w", err)
	}

	// Enable inbound
	if err := inbound.Enable(); err != nil {
		return fmt.Errorf("failed to enable inbound: %w", err)
	}

	// Save inbound
	if err := c.inboundRepo.Update(ctx, inbound); err != nil {
		return fmt.Errorf("failed to update inbound: %w", err)
	}

	// Audit: Inbound enabled
	auditLog := audit.NewLog(adminID, audit.ActionUpdate, "xray_inbound", inboundID, ip, userAgent)
	auditLog.AddMetadata("event", "inbound_enabled")
	if err := c.auditRepo.Create(ctx, auditLog); err != nil {
		_ = err // errcheck: acknowledged, audit failure is not critical
	}

	// Reload config (REUSES XrayProvisioningService)
	if err := c.provisioningService.RegenerateAndReload(ctx, inbound.XrayInstanceID, adminID, ip, userAgent); err != nil {
		// Rollback: Disable again
		if disErr := inbound.Disable(); disErr != nil {
			return fmt.Errorf("config reload failed AND rollback Disable failed: reload=%w, rollback=%v", err, disErr)
		}
		if updErr := c.inboundRepo.Update(ctx, inbound); updErr != nil {
			return fmt.Errorf("config reload failed AND rollback Update failed: reload=%w, rollback=%v", err, updErr)
		}
		
		auditLog := audit.NewLog(adminID, audit.ActionUpdate, "xray_inbound", inboundID, ip, userAgent)
		auditLog.AddMetadata("event", "inbound_rollback_after_reload_failed")
		if auditErr := c.auditRepo.Create(ctx, auditLog); auditErr != nil {
			_ = auditErr // errcheck: acknowledged, audit failure in rollback is not critical
		}
		
		return fmt.Errorf("config reload failed, inbound rolled back: %w", err)
	}

	return nil
}

// DisableInboundCommand handles disabling an inbound
type DisableInboundCommand struct {
	inboundRepo         xray.InboundRepository
	provisioningService *service.XrayProvisioningService
	auditRepo           audit.Repository
}

func NewDisableInboundCommand(
	inboundRepo xray.InboundRepository,
	provisioningService *service.XrayProvisioningService,
	auditRepo audit.Repository,
) *DisableInboundCommand {
	return &DisableInboundCommand{
		inboundRepo:         inboundRepo,
		provisioningService: provisioningService,
		auditRepo:           auditRepo,
	}
}

func (c *DisableInboundCommand) Execute(ctx context.Context, inboundID, adminID uuid.UUID, ip, userAgent string) error {
	// Find inbound
	inbound, err := c.inboundRepo.FindByID(ctx, inboundID)
	if err != nil {
		return fmt.Errorf("failed to find inbound: %w", err)
	}

	// Disable inbound
	if err := inbound.Disable(); err != nil {
		return fmt.Errorf("failed to disable inbound: %w", err)
	}

	// Save inbound
	if err := c.inboundRepo.Update(ctx, inbound); err != nil {
		return fmt.Errorf("failed to update inbound: %w", err)
	}

	// Audit: Inbound disabled
	auditLog := audit.NewLog(adminID, audit.ActionUpdate, "xray_inbound", inboundID, ip, userAgent)
	auditLog.AddMetadata("event", "inbound_disabled")
	if err := c.auditRepo.Create(ctx, auditLog); err != nil {
		_ = err // errcheck: acknowledged, audit failure is not critical
	}

	// Reload config (REUSES XrayProvisioningService)
	if err := c.provisioningService.RegenerateAndReload(ctx, inbound.XrayInstanceID, adminID, ip, userAgent); err != nil {
		// Rollback: Enable again
		if enErr := inbound.Enable(); enErr != nil {
			return fmt.Errorf("config reload failed AND rollback Enable failed: reload=%w, rollback=%v", err, enErr)
		}
		if updErr := c.inboundRepo.Update(ctx, inbound); updErr != nil {
			return fmt.Errorf("config reload failed AND rollback Update failed: reload=%w, rollback=%v", err, updErr)
		}
		
		auditLog := audit.NewLog(adminID, audit.ActionUpdate, "xray_inbound", inboundID, ip, userAgent)
		auditLog.AddMetadata("event", "inbound_rollback_after_reload_failed")
		if auditErr := c.auditRepo.Create(ctx, auditLog); auditErr != nil {
			_ = auditErr // errcheck: acknowledged, audit failure in rollback is not critical
		}
		
		return fmt.Errorf("config reload failed, inbound rolled back: %w", err)
	}

	return nil
}
