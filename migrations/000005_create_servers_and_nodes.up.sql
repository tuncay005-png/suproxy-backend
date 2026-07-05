-- Create servers table
CREATE TABLE servers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    country VARCHAR(100) NOT NULL,
    city VARCHAR(100) NOT NULL,
    hostname VARCHAR(255) NOT NULL UNIQUE,
    provider VARCHAR(100) NOT NULL,
    ipv4 VARCHAR(45) NOT NULL,
    ipv6 VARCHAR(45),
    status VARCHAR(20) NOT NULL DEFAULT 'offline' CHECK (status IN ('active', 'inactive', 'offline', 'maintenance')),
    is_public BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Create nodes table
CREATE TABLE nodes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    server_id UUID NOT NULL REFERENCES servers(id) ON DELETE CASCADE,
    protocol VARCHAR(20) NOT NULL CHECK (protocol IN ('vmess', 'vless', 'trojan', 'shadowsocks')),
    port INTEGER NOT NULL CHECK (port > 0 AND port <= 65535),
    max_users INTEGER NOT NULL CHECK (max_users > 0),
    current_users INTEGER NOT NULL DEFAULT 0 CHECK (current_users >= 0 AND current_users <= max_users),
    bandwidth_limit_bytes BIGINT NOT NULL DEFAULT 0 CHECK (bandwidth_limit_bytes >= 0),
    bandwidth_used_bytes BIGINT NOT NULL DEFAULT 0 CHECK (bandwidth_used_bytes >= 0),
    cpu_usage DECIMAL(5,2) NOT NULL DEFAULT 0 CHECK (cpu_usage >= 0 AND cpu_usage <= 100),
    ram_usage DECIMAL(5,2) NOT NULL DEFAULT 0 CHECK (ram_usage >= 0 AND ram_usage <= 100),
    latency_ms INTEGER NOT NULL DEFAULT 0 CHECK (latency_ms >= 0),
    version VARCHAR(50),
    health_status VARCHAR(20) NOT NULL DEFAULT 'unknown' CHECK (health_status IN ('healthy', 'degraded', 'unhealthy', 'unknown')),
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes for better query performance
CREATE INDEX idx_servers_name ON servers(name);
CREATE INDEX idx_servers_country ON servers(country);
CREATE INDEX idx_servers_status ON servers(status);
CREATE INDEX idx_servers_is_public ON servers(is_public);
CREATE INDEX idx_servers_country_status ON servers(country, status);

CREATE INDEX idx_nodes_server_id ON nodes(server_id);
CREATE INDEX idx_nodes_health_status ON nodes(health_status);
CREATE INDEX idx_nodes_protocol ON nodes(protocol);
CREATE INDEX idx_nodes_server_health ON nodes(server_id, health_status);

-- Unique constraint for port per server
CREATE UNIQUE INDEX idx_nodes_server_port ON nodes(server_id, port);

-- Add trigger for updated_at on servers table
CREATE TRIGGER update_servers_updated_at
    BEFORE UPDATE ON servers
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Add trigger for updated_at on nodes table
CREATE TRIGGER update_nodes_updated_at
    BEFORE UPDATE ON nodes
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Add comments for documentation
COMMENT ON TABLE servers IS 'Physical or virtual servers hosting VPN nodes';
COMMENT ON TABLE nodes IS 'VPN nodes running on servers';

COMMENT ON COLUMN servers.status IS 'Server operational status';
COMMENT ON COLUMN servers.is_public IS 'Whether server is publicly accessible';
COMMENT ON COLUMN nodes.bandwidth_limit_bytes IS '0 means unlimited bandwidth';
COMMENT ON COLUMN nodes.health_status IS 'Node health based on metrics';
COMMENT ON COLUMN nodes.protocol IS 'VPN protocol type (vmess, vless, trojan, shadowsocks)';
