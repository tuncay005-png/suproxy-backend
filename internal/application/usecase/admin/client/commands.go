package client

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/suproxy/backend/internal/application/service"
	"github.com/suproxy/backend/internal/domain/audit"
	"github.com/suproxy/backend/internal/domain/xray"
)

// CreateClientCommand handles manual client creation
type CreateClientCommand struct {
	clientRepo          xray.ClientRepository
	inboundRepo         xray.InboundRepository
	provisioningService *service.XrayProvisioningService
	auditRepo           audit.Repository
}

func NewCreateClientCommand(
	clientRepo xray.ClientRepository,
	inboundRepo xray.InboundRepository,
	provisioningService *service.XrayProvisioningService,
	auditRepo audit.Repository,
) *CreateClientCommand {
	return &CreateClientCommand{
		clientRepo:          clientRepo,
		inboundRepo:         inboundRepo,
		provisioningService: provisioningService,
		auditRepo:           auditRepo,
	}
}

func (c *CreateClientCommand) Execute(
	ctx context.Context,
	inboundID, userID uuid.UUID,
	email, flow string,
	adminID uuid.UUID,
	ip, userAgent string,
) (*xray.Client, error) {
	// Verify inbound exists
	inbound, err := c.inboundRepo.FindByID(ctx, inboundID)
	if err != nil {
		return nil, fmt.Errorf("failed to find inbound: %w", err)
	}

	// Generate UUID
	clientUUID := uuid.New().String()

	// Create client entity
	client, err := xray.NewClient(inboundID, userID, clientUUID, flow, email)
	if err != nil {
		return nil, fmt.Errorf("failed to create client entity: %w", err)
	}

	// Save client
	if err := c.clientRepo.Create(ctx, client); err != nil {
		return nil, fmt.Errorf("failed to save client: %w", err)
	}

	// Audit: Client created
	auditLog := audit.NewLog(adminID, audit.ActionCreate, "xray_client", client.ID, ip, userAgent)
	auditLog.AddMetadata("event", "client_created_manually")
	auditLog.AddMetadata("inbound_id", inboundID.String())
	auditLog.AddMetadata("user_id", userID.String())
	auditLog.AddMetadata("uuid", clientUUID)
	c.auditRepo.Create(ctx, auditLog)

	// Reload config (REUSES XrayProvisioningService)
	if err := c.provisioningService.RegenerateAndReload(ctx, inbound.XrayInstanceID, adminID, ip, userAgent); err != nil {
		// Rollback: Delete the created client
		c.clientRepo.Delete(ctx, client.ID)
		
		auditLog := audit.NewLog(adminID, audit.ActionDelete, "xray_client", client.ID, ip, userAgent)
		auditLog.AddMetadata("event", "client_rollback_after_reload_failed")
		c.auditRepo.Create(ctx, auditLog)
		
		return nil, fmt.Errorf("config reload failed, client rolled back: %w", err)
	}

	return client, nil
}

// DeleteClientCommand handles client deletion
type DeleteClientCommand struct {
	clientRepo          xray.ClientRepository
	inboundRepo         xray.InboundRepository
	provisioningService *service.XrayProvisioningService
	auditRepo           audit.Repository
}

func NewDeleteClientCommand(
	clientRepo xray.ClientRepository,
	inboundRepo xray.InboundRepository,
	provisioningService *service.XrayProvisioningService,
	auditRepo audit.Repository,
) *DeleteClientCommand {
	return &DeleteClientCommand{
		clientRepo:          clientRepo,
		inboundRepo:         inboundRepo,
		provisioningService: provisioningService,
		auditRepo:           auditRepo,
	}
}

func (c *DeleteClientCommand) Execute(ctx context.Context, clientID, adminID uuid.UUID, ip, userAgent string) error {
	// Find client
	client, err := c.clientRepo.FindByID(ctx, clientID)
	if err != nil {
		return fmt.Errorf("failed to find client: %w", err)
	}

	// Find inbound to get instance ID
	inbound, err := c.inboundRepo.FindByID(ctx, client.InboundID)
	if err != nil {
		return fmt.Errorf("failed to find inbound: %w", err)
	}

	// Delete client
	if err := c.clientRepo.Delete(ctx, clientID); err != nil {
		return fmt.Errorf("failed to delete client: %w", err)
	}

	// Audit: Client deleted
	auditLog := audit.NewLog(adminID, audit.ActionDelete, "xray_client", clientID, ip, userAgent)
	auditLog.AddMetadata("event", "client_deleted_manually")
	auditLog.AddMetadata("user_id", client.UserID.String())
	c.auditRepo.Create(ctx, auditLog)

	// Reload config (REUSES XrayProvisioningService)
	if err := c.provisioningService.RegenerateAndReload(ctx, inbound.XrayInstanceID, adminID, ip, userAgent); err != nil {
		return fmt.Errorf("client deleted but config reload failed: %w", err)
	}

	return nil
}

// EnableClientCommand handles enabling a client
type EnableClientCommand struct {
	clientRepo          xray.ClientRepository
	inboundRepo         xray.InboundRepository
	provisioningService *service.XrayProvisioningService
	auditRepo           audit.Repository
}

func NewEnableClientCommand(
	clientRepo xray.ClientRepository,
	inboundRepo xray.InboundRepository,
	provisioningService *service.XrayProvisioningService,
	auditRepo audit.Repository,
) *EnableClientCommand {
	return &EnableClientCommand{
		clientRepo:          clientRepo,
		inboundRepo:         inboundRepo,
		provisioningService: provisioningService,
		auditRepo:           auditRepo,
	}
}

func (c *EnableClientCommand) Execute(ctx context.Context, clientID, adminID uuid.UUID, ip, userAgent string) error {
	// Find client
	client, err := c.clientRepo.FindByID(ctx, clientID)
	if err != nil {
		return fmt.Errorf("failed to find client: %w", err)
	}

	// Find inbound
	inbound, err := c.inboundRepo.FindByID(ctx, client.InboundID)
	if err != nil {
		return fmt.Errorf("failed to find inbound: %w", err)
	}

	// Enable client
	if err := client.Enable(); err != nil {
		return fmt.Errorf("failed to enable client: %w", err)
	}

	// Save client
	if err := c.clientRepo.Update(ctx, client); err != nil {
		return fmt.Errorf("failed to update client: %w", err)
	}

	// Audit: Client enabled
	auditLog := audit.NewLog(adminID, audit.ActionUpdate, "xray_client", clientID, ip, userAgent)
	auditLog.AddMetadata("event", "client_enabled")
	c.auditRepo.Create(ctx, auditLog)

	// Reload config (REUSES XrayProvisioningService)
	if err := c.provisioningService.RegenerateAndReload(ctx, inbound.XrayInstanceID, adminID, ip, userAgent); err != nil {
		// Rollback: Disable again
		client.Disable()
		c.clientRepo.Update(ctx, client)
		
		auditLog := audit.NewLog(adminID, audit.ActionUpdate, "xray_client", clientID, ip, userAgent)
		auditLog.AddMetadata("event", "client_rollback_after_reload_failed")
		c.auditRepo.Create(ctx, auditLog)
		
		return fmt.Errorf("config reload failed, client rolled back: %w", err)
	}

	return nil
}

// DisableClientCommand handles disabling a client
type DisableClientCommand struct {
	clientRepo          xray.ClientRepository
	inboundRepo         xray.InboundRepository
	provisioningService *service.XrayProvisioningService
	auditRepo           audit.Repository
}

func NewDisableClientCommand(
	clientRepo xray.ClientRepository,
	inboundRepo xray.InboundRepository,
	provisioningService *service.XrayProvisioningService,
	auditRepo audit.Repository,
) *DisableClientCommand {
	return &DisableClientCommand{
		clientRepo:          clientRepo,
		inboundRepo:         inboundRepo,
		provisioningService: provisioningService,
		auditRepo:           auditRepo,
	}
}

func (c *DisableClientCommand) Execute(ctx context.Context, clientID, adminID uuid.UUID, ip, userAgent string) error {
	// Find client
	client, err := c.clientRepo.FindByID(ctx, clientID)
	if err != nil {
		return fmt.Errorf("failed to find client: %w", err)
	}

	// Find inbound
	inbound, err := c.inboundRepo.FindByID(ctx, client.InboundID)
	if err != nil {
		return fmt.Errorf("failed to find inbound: %w", err)
	}

	// Disable client
	if err := client.Disable(); err != nil {
		return fmt.Errorf("failed to disable client: %w", err)
	}

	// Save client
	if err := c.clientRepo.Update(ctx, client); err != nil {
		return fmt.Errorf("failed to update client: %w", err)
	}

	// Audit: Client disabled
	auditLog := audit.NewLog(adminID, audit.ActionUpdate, "xray_client", clientID, ip, userAgent)
	auditLog.AddMetadata("event", "client_disabled")
	c.auditRepo.Create(ctx, auditLog)

	// Reload config (REUSES XrayProvisioningService)
	if err := c.provisioningService.RegenerateAndReload(ctx, inbound.XrayInstanceID, adminID, ip, userAgent); err != nil {
		// Rollback: Enable again
		client.Enable()
		c.clientRepo.Update(ctx, client)
		
		auditLog := audit.NewLog(adminID, audit.ActionUpdate, "xray_client", clientID, ip, userAgent)
		auditLog.AddMetadata("event", "client_rollback_after_reload_failed")
		c.auditRepo.Create(ctx, auditLog)
		
		return fmt.Errorf("config reload failed, client rolled back: %w", err)
	}

	return nil
}

// RegenerateClientUUIDCommand handles regenerating client UUID
type RegenerateClientUUIDCommand struct {
	clientRepo          xray.ClientRepository
	inboundRepo         xray.InboundRepository
	provisioningService *service.XrayProvisioningService
	auditRepo           audit.Repository
}

func NewRegenerateClientUUIDCommand(
	clientRepo xray.ClientRepository,
	inboundRepo xray.InboundRepository,
	provisioningService *service.XrayProvisioningService,
	auditRepo audit.Repository,
) *RegenerateClientUUIDCommand {
	return &RegenerateClientUUIDCommand{
		clientRepo:          clientRepo,
		inboundRepo:         inboundRepo,
		provisioningService: provisioningService,
		auditRepo:           auditRepo,
	}
}

func (c *RegenerateClientUUIDCommand) Execute(ctx context.Context, clientID, adminID uuid.UUID, ip, userAgent string) (*xray.Client, error) {
	// Find client
	client, err := c.clientRepo.FindByID(ctx, clientID)
	if err != nil {
		return nil, fmt.Errorf("failed to find client: %w", err)
	}

	// Find inbound
	inbound, err := c.inboundRepo.FindByID(ctx, client.InboundID)
	if err != nil {
		return nil, fmt.Errorf("failed to find inbound: %w", err)
	}

	oldUUID := client.UUID
	newUUID := uuid.New().String()

	// Regenerate UUID
	if err := client.RegenerateUUID(newUUID); err != nil {
		return nil, fmt.Errorf("failed to regenerate UUID: %w", err)
	}

	// Save client
	if err := c.clientRepo.Update(ctx, client); err != nil {
		return nil, fmt.Errorf("failed to update client: %w", err)
	}

	// Audit: UUID regenerated
	auditLog := audit.NewLog(adminID, audit.ActionUpdate, "xray_client", clientID, ip, userAgent)
	auditLog.AddMetadata("event", "client_uuid_regenerated")
	auditLog.AddMetadata("old_uuid", oldUUID)
	auditLog.AddMetadata("new_uuid", newUUID)
	c.auditRepo.Create(ctx, auditLog)

	// Reload config (REUSES XrayProvisioningService)
	if err := c.provisioningService.RegenerateAndReload(ctx, inbound.XrayInstanceID, adminID, ip, userAgent); err != nil {
		// Rollback: Restore old UUID
		client.RegenerateUUID(oldUUID)
		c.clientRepo.Update(ctx, client)
		
		auditLog := audit.NewLog(adminID, audit.ActionUpdate, "xray_client", clientID, ip, userAgent)
		auditLog.AddMetadata("event", "client_rollback_after_reload_failed")
		c.auditRepo.Create(ctx, auditLog)
		
		return nil, fmt.Errorf("config reload failed, client rolled back: %w", err)
	}

	return client, nil
}

// ReprovisionClientCommand handles complete client reprovisioning
type ReprovisionClientCommand struct {
	clientRepo          xray.ClientRepository
	inboundRepo         xray.InboundRepository
	provisioningService *service.XrayProvisioningService
	auditRepo           audit.Repository
}

func NewReprovisionClientCommand(
	clientRepo xray.ClientRepository,
	inboundRepo xray.InboundRepository,
	provisioningService *service.XrayProvisioningService,
	auditRepo audit.Repository,
) *ReprovisionClientCommand {
	return &ReprovisionClientCommand{
		clientRepo:          clientRepo,
		inboundRepo:         inboundRepo,
		provisioningService: provisioningService,
		auditRepo:           auditRepo,
	}
}

func (c *ReprovisionClientCommand) Execute(ctx context.Context, clientID, adminID uuid.UUID, regenerateUUID bool, ip, userAgent string) (*xray.Client, error) {
	// Find client
	client, err := c.clientRepo.FindByID(ctx, clientID)
	if err != nil {
		return nil, fmt.Errorf("failed to find client: %w", err)
	}

	// Find inbound
	inbound, err := c.inboundRepo.FindByID(ctx, client.InboundID)
	if err != nil {
		return nil, fmt.Errorf("failed to find inbound: %w", err)
	}

	oldUUID := client.UUID

	// Optionally regenerate UUID
	if regenerateUUID {
		newUUID := uuid.New().String()
		if err := client.RegenerateUUID(newUUID); err != nil {
			return nil, fmt.Errorf("failed to regenerate UUID: %w", err)
		}

		// Save client with new UUID
		if err := c.clientRepo.Update(ctx, client); err != nil {
			return nil, fmt.Errorf("failed to update client: %w", err)
		}
	}

	// Audit: Client reprovisioned
	auditLog := audit.NewLog(adminID, audit.ActionUpdate, "xray_client", clientID, ip, userAgent)
	auditLog.AddMetadata("event", "client_reprovisioned")
	auditLog.AddMetadata("regenerate_uuid", regenerateUUID)
	if regenerateUUID {
		auditLog.AddMetadata("old_uuid", oldUUID)
		auditLog.AddMetadata("new_uuid", client.UUID)
	}
	c.auditRepo.Create(ctx, auditLog)

	// Reload config (REUSES XrayProvisioningService)
	if err := c.provisioningService.RegenerateAndReload(ctx, inbound.XrayInstanceID, adminID, ip, userAgent); err != nil {
		// Rollback: Restore old UUID if changed
		if regenerateUUID {
			client.RegenerateUUID(oldUUID)
			c.clientRepo.Update(ctx, client)
		}
		
		auditLog := audit.NewLog(adminID, audit.ActionUpdate, "xray_client", clientID, ip, userAgent)
		auditLog.AddMetadata("event", "client_rollback_after_reload_failed")
		c.auditRepo.Create(ctx, auditLog)
		
		return nil, fmt.Errorf("config reload failed, client rolled back: %w", err)
	}

	return client, nil
}
