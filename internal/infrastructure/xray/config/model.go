package config

// XrayConfig represents the complete Xray configuration structure
type XrayConfig struct {
	Log       *LogConfig       `json:"log,omitempty"`
	API       *APIConfig       `json:"api,omitempty"`
	Stats     *StatsConfig     `json:"stats,omitempty"`
	Policy    *PolicyConfig    `json:"policy,omitempty"`
	Inbounds  []InboundConfig  `json:"inbounds"`
	Outbounds []OutboundConfig `json:"outbounds"`
	Routing   *RoutingConfig   `json:"routing,omitempty"`
}

// LogConfig represents Xray logging configuration
type LogConfig struct {
	Access   string `json:"access,omitempty"`
	Error    string `json:"error,omitempty"`
	Loglevel string `json:"loglevel,omitempty"`
}

// APIConfig represents Xray API configuration for gRPC management
type APIConfig struct {
	Tag      string   `json:"tag"`
	Services []string `json:"services"`
}

// StatsConfig enables statistics collection
type StatsConfig struct{}

// PolicyConfig represents Xray policy configuration
type PolicyConfig struct {
	Levels  map[string]Level `json:"levels,omitempty"`
	System  *SystemPolicy    `json:"system,omitempty"`
}

// Level represents user level policy
type Level struct {
	Handshake         int  `json:"handshake,omitempty"`
	ConnIdle          int  `json:"connIdle,omitempty"`
	UplinkOnly        int  `json:"uplinkOnly,omitempty"`
	DownlinkOnly      int  `json:"downlinkOnly,omitempty"`
	StatsUserUplink   bool `json:"statsUserUplink,omitempty"`
	StatsUserDownlink bool `json:"statsUserDownlink,omitempty"`
}

// SystemPolicy represents system-level policy
type SystemPolicy struct {
	StatsInboundUplink    bool `json:"statsInboundUplink,omitempty"`
	StatsInboundDownlink  bool `json:"statsInboundDownlink,omitempty"`
	StatsOutboundUplink   bool `json:"statsOutboundUplink,omitempty"`
	StatsOutboundDownlink bool `json:"statsOutboundDownlink,omitempty"`
}

// InboundConfig represents an inbound configuration
type InboundConfig struct {
	Tag           string                 `json:"tag"`
	Port          int                    `json:"port"`
	Protocol      string                 `json:"protocol"`
	Settings      map[string]interface{} `json:"settings"`
	StreamSettings *StreamSettings       `json:"streamSettings,omitempty"`
	Sniffing      *SniffingConfig        `json:"sniffing,omitempty"`
}

// OutboundConfig represents an outbound configuration
type OutboundConfig struct {
	Tag           string                 `json:"tag"`
	Protocol      string                 `json:"protocol"`
	Settings      map[string]interface{} `json:"settings,omitempty"`
	StreamSettings *StreamSettings       `json:"streamSettings,omitempty"`
}

// StreamSettings represents stream configuration (transport and security)
type StreamSettings struct {
	Network      string           `json:"network"`
	Security     string           `json:"security"`
	TCPSettings  *TCPSettings     `json:"tcpSettings,omitempty"`
	WSSettings   *WSSettings      `json:"wsSettings,omitempty"`
	HTTPSettings *HTTPSettings    `json:"httpSettings,omitempty"`
	GRPCSettings *GRPCSettings    `json:"grpcSettings,omitempty"`
	QUICSettings *QUICSettings    `json:"quicSettings,omitempty"`
	TLSSettings  *TLSSettings     `json:"tlsSettings,omitempty"`
	RealitySettings *RealitySettings `json:"realitySettings,omitempty"`
}

// TCPSettings represents TCP transport settings
type TCPSettings struct {
	Header map[string]interface{} `json:"header,omitempty"`
}

// WSSettings represents WebSocket transport settings
type WSSettings struct {
	Path    string            `json:"path,omitempty"`
	Headers map[string]string `json:"headers,omitempty"`
}

// HTTPSettings represents HTTP transport settings
type HTTPSettings struct {
	Host []string `json:"host,omitempty"`
	Path string   `json:"path,omitempty"`
}

// GRPCSettings represents gRPC transport settings
type GRPCSettings struct {
	ServiceName string `json:"serviceName,omitempty"`
	MultiMode   bool   `json:"multiMode,omitempty"`
}

// QUICSettings represents QUIC transport settings
type QUICSettings struct {
	Security string            `json:"security,omitempty"`
	Key      string            `json:"key,omitempty"`
	Header   map[string]interface{} `json:"header,omitempty"`
}

// TLSSettings represents TLS security settings
type TLSSettings struct {
	ServerName string   `json:"serverName,omitempty"`
	Alpn       []string `json:"alpn,omitempty"`
	Certificates []Certificate `json:"certificates,omitempty"`
}

// RealitySettings represents REALITY protocol settings
type RealitySettings struct {
	Show        bool     `json:"show,omitempty"`
	Dest        string   `json:"dest,omitempty"`
	Xver        int      `json:"xver,omitempty"`
	ServerNames []string `json:"serverNames,omitempty"`
	PrivateKey  string   `json:"privateKey"`
	MinClientVer string  `json:"minClientVer,omitempty"`
	MaxClientVer string  `json:"maxClientVer,omitempty"`
	MaxTimeDiff int      `json:"maxTimeDiff,omitempty"`
	ShortIds    []string `json:"shortIds,omitempty"`
	Fingerprint string   `json:"fingerprint,omitempty"`
	ServerName  string   `json:"serverName,omitempty"`
	PublicKey   string   `json:"publicKey,omitempty"`
	ShortId     string   `json:"shortId,omitempty"`
	SpiderX     string   `json:"spiderX,omitempty"`
}

// Certificate represents TLS certificate
type Certificate struct {
	CertificateFile string `json:"certificateFile,omitempty"`
	KeyFile         string `json:"keyFile,omitempty"`
	Certificate     []string `json:"certificate,omitempty"`
	Key             []string `json:"key,omitempty"`
}

// SniffingConfig represents traffic sniffing configuration
type SniffingConfig struct {
	Enabled      bool     `json:"enabled"`
	DestOverride []string `json:"destOverride"`
}

// RoutingConfig represents routing rules
type RoutingConfig struct {
	DomainStrategy string        `json:"domainStrategy,omitempty"`
	Rules          []RoutingRule `json:"rules,omitempty"`
}

// RoutingRule represents a routing rule
type RoutingRule struct {
	Type        string   `json:"type,omitempty"`
	InboundTag  []string `json:"inboundTag,omitempty"`
	OutboundTag string   `json:"outboundTag,omitempty"`
	IP          []string `json:"ip,omitempty"`
	Domain      []string `json:"domain,omitempty"`
	Protocol    []string `json:"protocol,omitempty"`
}

// VMessSettings represents VMess protocol settings
type VMessSettings struct {
	Clients []VMessClient `json:"clients"`
}

// VMessClient represents a VMess client
type VMessClient struct {
	ID      string `json:"id"`
	AlterID int    `json:"alterId"`
	Email   string `json:"email"`
	Level   int    `json:"level,omitempty"`
}

// VLESSSettings represents VLESS protocol settings
type VLESSSettings struct {
	Clients    []VLESSClient `json:"clients"`
	Decryption string        `json:"decryption"`
	Fallbacks  []Fallback    `json:"fallbacks,omitempty"`
}

// VLESSClient represents a VLESS client
type VLESSClient struct {
	ID    string `json:"id"`
	Flow  string `json:"flow,omitempty"`
	Email string `json:"email"`
	Level int    `json:"level,omitempty"`
}

// Fallback represents VLESS fallback configuration
type Fallback struct {
	Dest string `json:"dest,omitempty"`
	Xver int    `json:"xver,omitempty"`
}

// TrojanSettings represents Trojan protocol settings
type TrojanSettings struct {
	Clients   []TrojanClient `json:"clients"`
	Fallbacks []Fallback     `json:"fallbacks,omitempty"`
}

// TrojanClient represents a Trojan client
type TrojanClient struct {
	Password string `json:"password"`
	Email    string `json:"email"`
	Level    int    `json:"level,omitempty"`
}

// ShadowsocksSettings represents Shadowsocks protocol settings
type ShadowsocksSettings struct {
	Method   string              `json:"method"`
	Password string              `json:"password"`
	Network  string              `json:"network,omitempty"`
	Clients  []ShadowsocksClient `json:"clients,omitempty"`
}

// ShadowsocksClient represents a Shadowsocks client
type ShadowsocksClient struct {
	Password string `json:"password"`
	Email    string `json:"email"`
	Level    int    `json:"level,omitempty"`
}
