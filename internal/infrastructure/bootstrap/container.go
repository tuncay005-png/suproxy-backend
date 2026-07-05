package bootstrap

import (
	"github.com/suproxy/backend/internal/domain/audit"
	"github.com/suproxy/backend/internal/domain/node"
	"github.com/suproxy/backend/internal/domain/server"
	"github.com/suproxy/backend/internal/domain/session"
	"github.com/suproxy/backend/internal/domain/subscription"
	"github.com/suproxy/backend/internal/domain/user"
	"github.com/suproxy/backend/internal/domain/xray"
	"github.com/suproxy/backend/internal/infrastructure/vpn"
	xrayBinary "github.com/suproxy/backend/internal/infrastructure/xray/binary"
	xrayConfig "github.com/suproxy/backend/internal/infrastructure/xray/config"
	xrayRuntime "github.com/suproxy/backend/internal/infrastructure/xray/runtime"
)

// Container holds all application dependencies
type Container struct {
	// Domain Repositories
	UserRepository              user.Repository
	RefreshTokenRepository      session.RefreshTokenRepository
	AuditLogRepository          audit.Repository
	SubscriptionRepository      subscription.SubscriptionRepository
	PlanRepository              subscription.PlanRepository
	ServerRepository            server.Repository
	NodeRepository              node.Repository
	XrayInstanceRepository      xray.XrayInstanceRepository
	InboundRepository           xray.InboundRepository
	ClientRepository            xray.ClientRepository
	RealityConfigRepository     xray.RealityConfigRepository

	// Xray Infrastructure
	XrayKernel        vpn.Kernel
	XrayProcessManager xrayRuntime.Manager
	XrayConfigWriter   xrayConfig.Writer
	XrayConfigGenerator xrayConfig.Generator
	XrayConfigValidator xrayConfig.Validator
	XrayBinaryManager  xrayBinary.Manager
}
