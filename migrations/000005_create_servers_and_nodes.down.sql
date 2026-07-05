-- Drop triggers
DROP TRIGGER IF EXISTS update_nodes_updated_at ON nodes;
DROP TRIGGER IF EXISTS update_servers_updated_at ON servers;

-- Drop indexes
DROP INDEX IF EXISTS idx_nodes_server_health;
DROP INDEX IF EXISTS idx_nodes_protocol;
DROP INDEX IF EXISTS idx_nodes_health_status;
DROP INDEX IF EXISTS idx_nodes_server_id;
DROP INDEX IF EXISTS idx_nodes_server_port;

DROP INDEX IF EXISTS idx_servers_country_status;
DROP INDEX IF EXISTS idx_servers_is_public;
DROP INDEX IF EXISTS idx_servers_status;
DROP INDEX IF EXISTS idx_servers_country;
DROP INDEX IF EXISTS idx_servers_name;

-- Drop tables (order matters due to foreign keys)
DROP TABLE IF EXISTS nodes;
DROP TABLE IF EXISTS servers;
