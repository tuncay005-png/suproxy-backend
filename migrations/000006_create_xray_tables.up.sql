-- Create xray_instances table
CREATE TABLE xray_instances (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    node_id UUID NOT NULL UNIQUE REFERENCES nodes(id) ON DELETE CASCADE,
    version VARCHAR(50) NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'stopped' CHECK (status IN ('running', 'stopped', 'failed')),
    config_version INTEGER NOT NULL DEFAULT 1,
    started_at TIMESTAMP WITH TIME ZONE,
    stopped_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Create inbounds table
CREATE TABLE inbounds (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    xray_instance_id UUID NOT NULL REFERENCES xray_instances(id) ON DELETE CASCADE,
    protocol VARCHAR(20) NOT NULL CHECK (protocol IN ('vmess', 'vless', 'trojan', 'shadowsocks')),
    port INTEGER NOT NULL CHECK (port > 0 AND port <= 65535),
    transport VARCHAR(20) NOT NULL CHECK (transport IN ('tcp', 'ws', 'http', 'quic', 'grpc', 'httpupgrade')),
    security VARCHAR(20) NOT NULL CHECK (security IN ('none', 'tls', 'reality')),
    enabled BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Create clients table
CREATE TABLE clients (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    inbound_id UUID NOT NULL REFERENCES inbounds(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    uuid VARCHAR(36) NOT NULL UNIQUE,
    flow VARCHAR(50),
    email VARCHAR(255) NOT NULL,
    enabled BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Create reality_configs table
CREATE TABLE reality_configs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    inbound_id UUID NOT NULL UNIQUE REFERENCES inbounds(id) ON DELETE CASCADE,
    private_key TEXT NOT NULL,
    public_key TEXT NOT NULL,
    short_id VARCHAR(16),
    server_name VARCHAR(255) NOT NULL,
    fingerprint VARCHAR(50) NOT NULL,
    spider_x VARCHAR(255),
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes for better query performance
CREATE INDEX idx_xray_instances_node_id ON xray_instances(node_id);
CREATE INDEX idx_xray_instances_status ON xray_instances(status);

CREATE INDEX idx_inbounds_xray_instance_id ON inbounds(xray_instance_id);
CREATE INDEX idx_inbounds_enabled ON inbounds(enabled);
CREATE INDEX idx_inbounds_protocol ON inbounds(protocol);

CREATE INDEX idx_clients_inbound_id ON clients(inbound_id);
CREATE INDEX idx_clients_user_id ON clients(user_id);
CREATE INDEX idx_clients_enabled ON clients(enabled);
CREATE INDEX idx_clients_uuid ON clients(uuid);

CREATE INDEX idx_reality_configs_inbound_id ON reality_configs(inbound_id);

-- Unique constraint for port per xray instance
CREATE UNIQUE INDEX idx_inbounds_instance_port ON inbounds(xray_instance_id, port);

-- Add triggers for updated_at
CREATE TRIGGER update_xray_instances_updated_at
    BEFORE UPDATE ON xray_instances
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_inbounds_updated_at
    BEFORE UPDATE ON inbounds
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_clients_updated_at
    BEFORE UPDATE ON clients
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_reality_configs_updated_at
    BEFORE UPDATE ON reality_configs
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Add comments for documentation
COMMENT ON TABLE xray_instances IS 'Xray-core instances running on nodes';
COMMENT ON TABLE inbounds IS 'Xray inbound configurations';
COMMENT ON TABLE clients IS 'Xray client configurations (users)';
COMMENT ON TABLE reality_configs IS 'REALITY protocol configurations for inbounds';

COMMENT ON COLUMN xray_instances.config_version IS 'Incremented on configuration changes to trigger reload';
COMMENT ON COLUMN inbounds.transport IS 'Transport protocol (tcp, ws, http, quic, grpc, httpupgrade)';
COMMENT ON COLUMN inbounds.security IS 'Security layer (none, tls, reality)';
COMMENT ON COLUMN clients.uuid IS 'Client UUID for Xray authentication';
COMMENT ON COLUMN clients.flow IS 'XTLS flow control (e.g., xtls-rprx-vision)';
