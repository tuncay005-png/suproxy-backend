-- Revert token_hash column size
ALTER TABLE refresh_tokens ALTER COLUMN token_hash TYPE VARCHAR(255);
