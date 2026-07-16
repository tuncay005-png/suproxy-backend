-- Increase token_hash column size to accommodate full JWT tokens in tests
-- SHA256 hash is 64 characters, but tests may store full tokens
ALTER TABLE refresh_tokens ALTER COLUMN token_hash TYPE VARCHAR(1024);
