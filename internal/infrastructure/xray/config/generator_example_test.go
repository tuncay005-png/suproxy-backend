package config_test

import (
	"encoding/json"
	"testing"

	"github.com/suproxy/backend/internal/infrastructure/xray/config"
)

// TestGenerateExampleConfig demonstrates a complete VLESS+Reality configuration
func TestGenerateExampleConfig(t *testing.T) {
	// Create a basic generator (without repos for example)
	gen := &exampleGenerator{}

	// Generate example config
	cfg := gen.GenerateExampleConfig()

	// Marshal to JSON
	jsonData, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		t.Fatalf("Failed to marshal config: %v", err)
	}

	// Print the config
	t.Logf("Generated Example Config:\n%s", string(jsonData))

	// Validate it can be unmarshaled
	var validateCfg config.XrayConfig
	if err := json.Unmarshal(jsonData, &validateCfg); err != nil {
		t.Fatalf("Failed to unmarshal config: %v", err)
	}

	// Basic assertions
	if len(validateCfg.Inbounds) == 0 {
		t.Error("Expected at least one inbound")
	}

	if len(validateCfg.Outbounds) == 0 {
		t.Error("Expected at least one outbound")
	}

	if validateCfg.DNS == nil {
		t.Error("Expected DNS config")
	}

	if validateCfg.Policy == nil {
		t.Error("Expected Policy config")
	}

	if validateCfg.Routing == nil {
		t.Error("Expected Routing config")
	}
}

// exampleGenerator is a minimal generator for testing
type exampleGenerator struct{}

func (g *exampleGenerator) GenerateExampleConfig() *config.XrayConfig {
	return &config.XrayConfig{
		Log: &config.LogConfig{
			Loglevel: "warning",
		},
		API: &config.APIConfig{
			Tag:      "api",
			Services: []string{"HandlerService", "StatsService", "LoggerService"},
		},
		DNS: &config.DNSConfig{
			Servers: []config.DNSServer{
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
		},
		Stats: &config.StatsConfig{},
		Policy: &config.PolicyConfig{
			Levels: map[string]config.Level{
				"0": {
					Handshake:         4,
					ConnIdle:          300,
					UplinkOnly:        2,
					DownlinkOnly:      5,
					StatsUserUplink:   true,
					StatsUserDownlink: true,
				},
			},
			System: &config.SystemPolicy{
				StatsInboundUplink:    true,
				StatsInboundDownlink:  true,
				StatsOutboundUplink:   true,
				StatsOutboundDownlink: true,
			},
		},
		Inbounds: []config.InboundConfig{
			// API Inbound
			{
				Tag:      "api",
				Port:     10085,
				Protocol: "dokodemo-door",
				Settings: map[string]interface{}{
					"address": "127.0.0.1",
				},
			},
			// VLESS + Reality Inbound
			{
				Tag:      "vless-reality-443",
				Port:     443,
				Protocol: "vless",
				Settings: map[string]interface{}{
					"clients": []map[string]interface{}{
						{
							"id":    "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
							"flow":  "xtls-rprx-vision",
							"email": "user@suproxy.com",
							"level": 0,
						},
					},
					"decryption": "none",
				},
				StreamSettings: &config.StreamSettings{
					Network:  "tcp",
					Security: "reality",
					TCPSettings: &config.TCPSettings{
						Header: map[string]interface{}{
							"type": "none",
						},
					},
					RealitySettings: &config.RealitySettings{
						Show:        false,
						Dest:        "www.google.com:443",
						Xver:        0,
						ServerNames: []string{"www.google.com", "www.youtube.com"},
						PrivateKey:  "EXAMPLE_PRIVATE_KEY_32_CHARS_BASE64",
						ShortIds:    []string{"0123456789abcdef", "fedcba9876543210"},
						Fingerprint: "chrome",
						SpiderX:     "/",
					},
				},
				Sniffing: &config.SniffingConfig{
					Enabled:      true,
					DestOverride: []string{"http", "tls", "quic"},
				},
			},
		},
		Outbounds: []config.OutboundConfig{
			{
				Tag:      "direct",
				Protocol: "freedom",
				Settings: map[string]interface{}{
					"domainStrategy": "UseIPv4",
				},
			},
			{
				Tag:      "block",
				Protocol: "blackhole",
				Settings: map[string]interface{}{
					"response": map[string]interface{}{
						"type": "http",
					},
				},
			},
		},
		Routing: &config.RoutingConfig{
			DomainStrategy: "AsIs",
			Rules: []config.RoutingRule{
				{
					Type:        "field",
					InboundTag:  []string{"api"},
					OutboundTag: "api",
				},
				{
					Type: "field",
					Domain: []string{
						"geosite:category-ads-all",
					},
					OutboundTag: "block",
				},
				{
					Type: "field",
					IP: []string{
						"geoip:private",
					},
					OutboundTag: "direct",
				},
			},
		},
	}
}
