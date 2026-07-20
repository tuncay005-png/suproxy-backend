package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"regexp"
)

var (
	ErrInvalidPort        = errors.New("invalid port number")
	ErrInvalidUUID        = errors.New("invalid UUID format")
	ErrInvalidProtocol    = errors.New("invalid protocol")
	ErrInvalidTransport   = errors.New("invalid transport")
	ErrInvalidSecurity    = errors.New("invalid security")
	ErrMissingClients     = errors.New("missing clients")
	ErrMissingRealityKey  = errors.New("missing reality private key")
	ErrInvalidFingerprint = errors.New("invalid fingerprint")
	ErrDuplicatePort      = errors.New("duplicate port detected")
)

// Validator validates Xray configuration
type Validator interface {
	// Validate validates a complete configuration
	Validate(config *XrayConfig) error

	// ValidateJSON validates JSON configuration
	ValidateJSON(configJSON []byte) error

	// ValidateInbound validates a single inbound configuration
	ValidateInbound(inbound *InboundConfig) error

	// ValidatePort validates port number
	ValidatePort(port int) error

	// ValidateUUID validates UUID format
	ValidateUUID(uuid string) error
}

type validator struct {
	uuidRegex        *regexp.Regexp
	fingerprintRegex *regexp.Regexp
}

// NewValidator creates a new config validator
func NewValidator() Validator {
	return &validator{
		uuidRegex:        regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`),
		fingerprintRegex: regexp.MustCompile(`^(chrome|firefox|safari|edge|ios|android|random)$`),
	}
}

func (v *validator) Validate(config *XrayConfig) error {
	if config == nil {
		return ErrInvalidConfig
	}

	// Validate inbounds
	if len(config.Inbounds) == 0 {
		return errors.New("no inbounds configured")
	}

	// Check for duplicate ports
	ports := make(map[int]bool)
	for _, inbound := range config.Inbounds {
		if err := v.ValidateInbound(&inbound); err != nil {
			return fmt.Errorf("inbound %s: %w", inbound.Tag, err)
		}

		if ports[inbound.Port] {
			return fmt.Errorf("%w: port %d", ErrDuplicatePort, inbound.Port)
		}
		ports[inbound.Port] = true
	}

	// Validate outbounds
	if len(config.Outbounds) == 0 {
		return errors.New("no outbounds configured")
	}

	return nil
}

func (v *validator) ValidateJSON(configJSON []byte) error {
	var config XrayConfig
	if err := json.Unmarshal(configJSON, &config); err != nil {
		return fmt.Errorf("invalid JSON: %w", err)
	}

	return v.Validate(&config)
}

func (v *validator) ValidateInbound(inbound *InboundConfig) error {
	if inbound == nil {
		return ErrInvalidConfig
	}

	// Validate port
	if err := v.ValidatePort(inbound.Port); err != nil {
		return err
	}

	// Validate protocol
	if err := v.validateProtocol(inbound.Protocol); err != nil {
		return err
	}

	// Validate settings based on protocol
	switch inbound.Protocol {
	case "vmess":
		return v.validateVMessSettings(inbound.Settings)
	case "vless":
		return v.validateVLESSSettings(inbound.Settings)
	case "trojan":
		return v.validateTrojanSettings(inbound.Settings)
	case "shadowsocks":
		return v.validateShadowsocksSettings(inbound.Settings)
	}

	// Validate stream settings
	if inbound.StreamSettings != nil {
		if err := v.validateStreamSettings(inbound.StreamSettings); err != nil {
			return err
		}
	}

	return nil
}

func (v *validator) ValidatePort(port int) error {
	if port <= 0 || port > 65535 {
		return fmt.Errorf("%w: %d", ErrInvalidPort, port)
	}

	// Check if port is available
	if v.isPortReserved(port) {
		return fmt.Errorf("port %d is reserved", port)
	}

	return nil
}

func (v *validator) ValidateUUID(uuid string) error {
	if !v.uuidRegex.MatchString(uuid) {
		return fmt.Errorf("%w: %s", ErrInvalidUUID, uuid)
	}
	return nil
}

func (v *validator) validateProtocol(protocol string) error {
	validProtocols := map[string]bool{
		"vmess":         true,
		"vless":         true,
		"trojan":        true,
		"shadowsocks":   true,
		"dokodemo-door": true,
		"http":          true,
		"socks":         true,
	}

	if !validProtocols[protocol] {
		return fmt.Errorf("%w: %s", ErrInvalidProtocol, protocol)
	}

	return nil
}

func (v *validator) validateVMessSettings(settings map[string]interface{}) error {
	clients, ok := settings["clients"].([]interface{})
	if !ok || len(clients) == 0 {
		return ErrMissingClients
	}

	// Validate each client
	for i, client := range clients {
		clientMap, ok := client.(map[string]interface{})
		if !ok {
			continue
		}

		uuid, ok := clientMap["id"].(string)
		if !ok {
			return fmt.Errorf("client %d: missing UUID", i)
		}

		if err := v.ValidateUUID(uuid); err != nil {
			return fmt.Errorf("client %d: %w", i, err)
		}
	}

	return nil
}

func (v *validator) validateVLESSSettings(settings map[string]interface{}) error {
	clients, ok := settings["clients"].([]interface{})
	if !ok || len(clients) == 0 {
		return ErrMissingClients
	}

	// Validate decryption
	decryption, ok := settings["decryption"].(string)
	if !ok || decryption != "none" {
		return errors.New("VLESS decryption must be 'none'")
	}

	// Validate each client
	for i, client := range clients {
		clientMap, ok := client.(map[string]interface{})
		if !ok {
			continue
		}

		uuid, ok := clientMap["id"].(string)
		if !ok {
			return fmt.Errorf("client %d: missing UUID", i)
		}

		if err := v.ValidateUUID(uuid); err != nil {
			return fmt.Errorf("client %d: %w", i, err)
		}

		// Validate flow if present
		if flow, ok := clientMap["flow"].(string); ok && flow != "" {
			if err := v.validateFlow(flow); err != nil {
				return fmt.Errorf("client %d: %w", i, err)
			}
		}
	}

	return nil
}

func (v *validator) validateTrojanSettings(settings map[string]interface{}) error {
	clients, ok := settings["clients"].([]interface{})
	if !ok || len(clients) == 0 {
		return ErrMissingClients
	}

	// Validate each client
	for i, client := range clients {
		clientMap, ok := client.(map[string]interface{})
		if !ok {
			continue
		}

		password, ok := clientMap["password"].(string)
		if !ok || password == "" {
			return fmt.Errorf("client %d: missing password", i)
		}
	}

	return nil
}

func (v *validator) validateShadowsocksSettings(settings map[string]interface{}) error {
	method, ok := settings["method"].(string)
	if !ok || method == "" {
		return errors.New("missing shadowsocks method")
	}

	password, ok := settings["password"].(string)
	if !ok || password == "" {
		return errors.New("missing shadowsocks password")
	}

	return nil
}

func (v *validator) validateStreamSettings(stream *StreamSettings) error {
	// Validate network
	if err := v.validateTransport(stream.Network); err != nil {
		return err
	}

	// Validate security
	if err := v.validateSecurity(stream.Security); err != nil {
		return err
	}

	// Validate REALITY settings if security is reality
	if stream.Security == "reality" && stream.RealitySettings != nil {
		if err := v.validateRealitySettings(stream.RealitySettings); err != nil {
			return err
		}
	}

	return nil
}

func (v *validator) validateTransport(transport string) error {
	validTransports := map[string]bool{
		"tcp":         true,
		"ws":          true,
		"http":        true,
		"grpc":        true,
		"quic":        true,
		"httpupgrade": true,
	}

	if !validTransports[transport] {
		return fmt.Errorf("%w: %s", ErrInvalidTransport, transport)
	}

	return nil
}

func (v *validator) validateSecurity(security string) error {
	validSecurity := map[string]bool{
		"none":    true,
		"tls":     true,
		"reality": true,
	}

	if !validSecurity[security] {
		return fmt.Errorf("%w: %s", ErrInvalidSecurity, security)
	}

	return nil
}

func (v *validator) validateRealitySettings(reality *RealitySettings) error {
	if reality.PrivateKey == "" {
		return ErrMissingRealityKey
	}

	if reality.Fingerprint != "" {
		if !v.fingerprintRegex.MatchString(reality.Fingerprint) {
			return fmt.Errorf("%w: %s", ErrInvalidFingerprint, reality.Fingerprint)
		}
	}

	if len(reality.ServerNames) == 0 && reality.ServerName == "" {
		return errors.New("missing reality server name")
	}

	return nil
}

func (v *validator) validateFlow(flow string) error {
	validFlows := map[string]bool{
		"xtls-rprx-vision":        true,
		"xtls-rprx-vision-udp443": true,
	}

	if !validFlows[flow] {
		return fmt.Errorf("invalid flow: %s", flow)
	}

	return nil
}

func (v *validator) isPortReserved(port int) bool {
	// System reserved ports
	if port < 1024 {
		return true
	}

	// Check if port is in use (basic check)
	addr := fmt.Sprintf(":%d", port)
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		// Port is likely in use or unavailable
		return false // Let the actual start process handle it
	}
	_ = ln.Close()
	return false
}
