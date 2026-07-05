-- Drop triggers
DROP TRIGGER IF EXISTS update_reality_configs_updated_at ON reality_configs;
DROP TRIGGER IF EXISTS update_clients_updated_at ON clients;
DROP TRIGGER IF EXISTS update_inbounds_updated_at ON inbounds;
DROP TRIGGER IF EXISTS update_xray_instances_updated_at ON xray_instances;

-- Drop indexes
DROP INDEX IF EXISTS idx_reality_configs_inbound_id;

DROP INDEX IF EXISTS idx_clients_uuid;
DROP INDEX IF EXISTS idx_clients_enabled;
DROP INDEX IF EXISTS idx_clients_user_id;
DROP INDEX IF EXISTS idx_clients_inbound_id;

DROP INDEX IF EXISTS idx_inbounds_protocol;
DROP INDEX IF EXISTS idx_inbounds_enabled;
DROP INDEX IF EXISTS idx_inbounds_xray_instance_id;
DROP INDEX IF EXISTS idx_inbounds_instance_port;

DROP INDEX IF EXISTS idx_xray_instances_status;
DROP INDEX IF EXISTS idx_xray_instances_node_id;

-- Drop tables (order matters due to foreign keys)
DROP TABLE IF EXISTS reality_configs;
DROP TABLE IF EXISTS clients;
DROP TABLE IF EXISTS inbounds;
DROP TABLE IF EXISTS xray_instances;
