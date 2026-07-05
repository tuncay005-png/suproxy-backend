-- Drop trigger function
DROP FUNCTION IF EXISTS update_updated_at_column();

-- Drop schema migrations table
DROP TABLE IF EXISTS schema_migrations;

-- Drop UUID extension
DROP EXTENSION IF EXISTS "uuid-ossp";
