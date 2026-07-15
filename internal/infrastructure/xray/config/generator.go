package config

import (
	"context"
	"encoding/json"

	"github.com/google/uuid"
	"github.com/suproxy/backend/internal/domain/xray"
)

// Generator generates Xray configuration from domain entities
type Generator interface {
	// Generate creates a complete Xray configuration
	Generate(ctx context.Context, instanceID uuid.UUID) (*XrayConfig, error)

	// GenerateJSON creates JSON bytes of the configuration
	GenerateJSON(ctx context.Context, instanceID uuid.UUID) ([]byte, error)

	// GenerateInbound creates an inbound configuration
	GenerateInbound(ctx context.Context, inbound *xray.Inbound, clients []*xray.Client, reality *xray.RealityConfig) (*InboundConfig, error)
}

type generator struct {
	instanceRepo xray.XrayInstanceRepository
	inboundRepo  xray.InboundRepository
	clientRepo   xray.ClientRepository
	realityRepo  xray.RealityConfigRepository
}

// NewGenerator creates a new config generator
func NewGenerator(
	instanceRepo xray.XrayInstanceRepository,
	inboundRepo xray.InboundRepository,
	clientRepo xray.ClientRepository,
	realityRepo xray.RealityConfigRepository,
) Generator {
	return &generator{
		instanceRepo: instanceRepo,
		inboundRepo:  inboundRepo,
		clientRepo:   clientRepo,
		realityRepo:  realityRepo,
	}
}

func (g *generator) Generate(ctx context.Context, instanceID uuid.UUID) (*XrayConfig, error) {
	// Find instance
	instance, err := g.instanceRepo.FindByID(ctx, instanceID)
	if err != nil {
		return nil, err
	}

	// Find all enabled inbounds for this instance
	inbounds, err := g.inboundRepo.FindEnabledByInstanceID(ctx, instanceID)
	if err != nil {
		return nil, err
	}

	// Build config with production-grade settings
	config := &XrayConfig{
		Log: &LogConfig{
			Loglevel: "warning",
			Access:   "", // Empty means no access log
			Error:    "", // Empty means stderr
		},
		API: &APIConfig{
			Tag:      "api",
			Services: []string{"HandlerService", "StatsService", "LoggerService"},
		},
		DNS:       g.generateDNS(),
		Stats:     &StatsConfig{},
		Policy:    g.generatePolicy(),
		Inbounds:  make([]InboundConfig, 0, len(inbounds)+1),
		Outbounds: g.generateOutbounds(),
		Routing:   g.generateRouting(),
	}

	// Add API inbound for management
	config.Inbounds = append(config.Inbounds, g.generateAPIInbound())

	// Generate inbound configs
	for _, inbound := range inbounds {
		// Find clients for this inbound
		clients, err := g.clientRepo.FindEnabledByInboundID(ctx, inbound.ID)
		if err != nil {
			continue
		}

		// Skip if no clients
		if len(clients) == 0 {
			continue
		}

		// Find reality config if security is reality
		var reality *xray.RealityConfig
		if inbound.Security == xray.SecurityREALITY {
			reality, _ = g.realityRepo.FindByInboundID(ctx, inbound.ID)
		}

		inboundConfig, err := g.GenerateInbound(ctx, inbound, clients, reality)
		if err != nil {
			continue
		}

		config.Inbounds = append(config.Inbounds, *inboundConfig)
	}

	_ = instance // Use instance for future metadata
	return config, nil
}

func (g *generator) GenerateJSON(ctx context.Context, instanceID uuid.UUID) ([]byte, error) {
	config, err := g.Generate(ctx, instanceID)
	if err != nil {
		return nil, err
	}

	return json.MarshalIndent(config, "", "  ")
}

func (g *generator) GenerateInbound(ctx context.Context, inbound *xray.Inbound, clients []*xray.Client, reality *xray.RealityConfig) (*InboundConfig, error) {
	inboundConfig := &InboundConfig{
		Tag:      inbound.ID.String(),
		Port:     inbound.Port,
		Protocol: string(inbound.Protocol),
		Sniffing: &SniffingConfig{
			Enabled:      true,
			DestOverride: []string{"http", "tls", "quic"},
		},
	}

	// Generate protocol settings
	switch inbound.Protocol {
	case xray.ProtocolVMess:
		inboundConfig.Settings = g.generateVMessSettings(clients)
	case xray.ProtocolVLESS:
		inboundConfig.Settings = g.generateVLESSSettings(clients)
	case xray.ProtocolTrojan:
		inboundConfig.Settings = g.generateTrojanSettings(clients)
	case xray.ProtocolShadowsock:
		inboundConfig.Settings = g.generateShadowsocksSettings(clients)
	}

	// Generate stream settings with proper Reality support
	inboundConfig.StreamSettings = g.generateStreamSettings(inbound, reality)

	return inboundConfig, nil
}

func (g *generator) generatePolicy() *PolicyConfig {
	return &PolicyConfig{
		Levels: map[string]Level{
			"0": {
				Handshake:         4,
				ConnIdle:          300,
				UplinkOnly:        2,
				DownlinkOnly:      5,
				StatsUserUplink:   true,
				StatsUserDownlink: true,
			},
		},
		System: &SystemPolicy{
			StatsInboundUplink:    true,
			StatsInboundDownlink:  true,
			StatsOutboundUplink:   true,
			StatsOutboundDownlink: true,
		},
	}
}

func (g *generator) generateAPIInbound() InboundConfig {
	return InboundConfig{
		Tag:      "api",
		Port:     10085,
		Protocol: "dokodemo-door",
		Settings: map[string]interface{}{
			"address": "127.0.0.1",
		},
		// Note: API inbound should only be accessible locally
	}
}

// GenerateExampleConfig generates a complete example config for testing
func (g *generator) GenerateExampleConfig() *XrayConfig {
	return &XrayConfig{
		Log: &LogConfig{
			Loglevel: "warning",
		},
		API: &APIConfig{
			Tag:      "api",
			Services: []string{"HandlerService", "StatsService", "LoggerService"},
		},
		DNS:    g.generateDNS(),
		Stats:  &StatsConfig{},
		Policy: g.generatePolicy(),
		Inbounds: []InboundConfig{
			{
				Tag:      "vless-reality-example",
				Port:     443,
				Protocol: "vless",
				Settings: map[string]interface{}{
					"clients": []VLESSClient{
						{
							ID:    "00000000-0000-0000-0000-000000000000",
							Flow:  "xtls-rprx-vision",
							Email: "example@suproxy.com",
							Level: 0,
						},
					},
					"decryption": "none",
				},
				StreamSettings: &StreamSettings{
					Network:  "tcp",
					Security: "reality",
					TCPSettings: &TCPSettings{
						Header: map[string]interface{}{
							"type": "none",
						},
					},
					RealitySettings: &RealitySettings{
						Show:        false,
						Dest:        "www.google.com:443",
						Xver:        0,
						ServerNames: []string{"www.google.com", "www.youtube.com"},
						PrivateKey:  "PRIVATE_KEY_HERE",
						ShortIds:    []string{"0123456789abcdef"},
						Fingerprint: "chrome",
						SpiderX:     "/",
					},
				},
				Sniffing: &SniffingConfig{
					Enabled:      true,
					DestOverride: []string{"http", "tls", "quic"},
				},
			},
		},
		Outbounds: g.generateOutbounds(),
		Routing:   g.generateRouting(),
	}
}

func (g *generator) generateOutbounds() []OutboundConfig {
	return []OutboundConfig{
		// Direct outbound for direct connections
		{
			Tag:      "direct",
			Protocol: "freedom",
			Settings: map[string]interface{}{
				"domainStrategy": "UseIPv4",
			},
		},
		// Block outbound for blocking traffic
		{
			Tag:      "block",
			Protocol: "blackhole",
			Settings: map[string]interface{}{
				"response": map[string]interface{}{
					"type": "http",
				},
			},
		},
	}
}

func (g *generator) generateRouting() *RoutingConfig {
	return &RoutingConfig{
		DomainStrategy: "AsIs",
		Rules: []RoutingRule{
			// API routing
			{
				Type:        "field",
				InboundTag:  []string{"api"},
				OutboundTag: "api",
			},
			// Block ads and trackers
			{
				Type: "field",
				Domain: []string{
					"geosite:category-ads-all",
				},
				OutboundTag: "block",
			},
			// Direct connection for private IPs
			{
				Type: "field",
				IP: []string{
					"geoip:private",
				},
				OutboundTag: "direct",
			},
			// Default: direct
		},
	}
}

// generateDNS creates DNS configuration
func (g *generator) generateDNS() *DNSConfig {
	return &DNSConfig{
		Servers: []DNSServer{
			{
				Address: "1.1.1.1",
				Port:    53,
				Domains: []string{"geosite:geolocation-!cn"},
			},
			{
				Address: "223.5.5.5",
				Port:    53,
				Domains: []string{"geosite:cn"},
			},
			{
				Address: "localhost",
			},
		},
		QueryStrategy: "UseIPv4",
	}
}

func (g *generator) generateVMessSettings(clients []*xray.Client) map[string]interface{} {
	vclients := make([]VMessClient, 0, len(clients))
	for _, client := range clients {
		vclients = append(vclients, VMessClient{
			ID:      client.UUID,
			AlterID: 0,
			Email:   client.Email,
			Level:   0,
		})
	}

	return map[string]interface{}{
		"clients": vclients,
	}
}

func (g *generator) generateVLESSSettings(clients []*xray.Client) map[string]interface{} {
	vclients := make([]VLESSClient, 0, len(clients))
	for _, client := range clients {
		vclients = append(vclients, VLESSClient{
			ID:    client.UUID,
			Flow:  client.Flow,
			Email: client.Email,
			Level: 0,
		})
	}

	return map[string]interface{}{
		"clients":    vclients,
		"decryption": "none",
	}
}

func (g *generator) generateTrojanSettings(clients []*xray.Client) map[string]interface{} {
	tclients := make([]TrojanClient, 0, len(clients))
	for _, client := range clients {
		tclients = append(tclients, TrojanClient{
			Password: client.UUID,
			Email:    client.Email,
			Level:    0,
		})
	}

	return map[string]interface{}{
		"clients": tclients,
	}
}

func (g *generator) generateShadowsocksSettings(clients []*xray.Client) map[string]interface{} {
	// For shadowsocks, use first client's UUID as password
	password := ""
	if len(clients) > 0 {
		password = clients[0].UUID
	}

	return map[string]interface{}{
		"method":   "aes-256-gcm",
		"password": password,
		"network":  "tcp,udp",
	}
}

func (g *generator) generateStreamSettings(inbound *xray.Inbound, reality *xray.RealityConfig) *StreamSettings {
	stream := &StreamSettings{
		Network:  string(inbound.Transport),
		Security: string(inbound.Security),
	}

	// Configure transport with production-grade settings
	switch inbound.Transport {
	case xray.TransportTCP:
		stream.TCPSettings = &TCPSettings{
			Header: map[string]interface{}{
				"type": "none",
			},
		}
	case xray.TransportWebSocket:
		stream.WSSettings = &WSSettings{
			Path: "/",
			Headers: map[string]string{
				"Host": "www.example.com",
			},
		}
	case xray.TransportHTTP:
		stream.HTTPSettings = &HTTPSettings{
			Host: []string{"www.example.com"},
			Path: "/",
		}
	case xray.TransportGRPC:
		stream.GRPCSettings = &GRPCSettings{
			ServiceName: "GunService",
			MultiMode:   false,
		}
	case xray.TransportQUIC:
		stream.QUICSettings = &QUICSettings{
			Security: "none",
			Key:      "",
			Header: map[string]interface{}{
				"type": "none",
			},
		}
	}

	// Configure security with full Reality support
	switch inbound.Security {
	case xray.SecurityTLS:
		stream.TLSSettings = &TLSSettings{
			Alpn: []string{"h2", "http/1.1"},
		}
	case xray.SecurityREALITY:
		if reality != nil {
			// Production-grade Reality configuration
			stream.RealitySettings = &RealitySettings{
				Show:         false,
				Dest:         "www.microsoft.com:443", // Default destination
				Xver:         0,
				ServerNames:  []string{reality.ServerName},
				PrivateKey:   reality.PrivateKey,
				MinClientVer: "",
				MaxClientVer: "",
				MaxTimeDiff:  0,
				ShortIds:     []string{reality.ShortID},
				Fingerprint:  reality.Fingerprint,
				SpiderX:      reality.SpiderX,
			}
		} else {
			// Fallback Reality config if not provided
			stream.RealitySettings = &RealitySettings{
				Show:        false,
				Dest:        "www.google.com:443",
				Xver:        0,
				ServerNames: []string{"www.google.com"},
				PrivateKey:  "", // Must be provided by domain entity
				ShortIds:    []string{""},
				Fingerprint: "chrome",
				SpiderX:     "/",
			}
		}
	}

	return stream
}
