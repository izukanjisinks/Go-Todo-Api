-- Remove session_token and csrf_token columns from users table
-- These are no longer needed with JWT authentication

-- First, drop the index on session_token
DROP INDEX IF EXISTS idx_users_session_token ON users;

-- Now drop the columns
ALTER TABLE users DROP COLUMN IF EXISTS session_token;
ALTER TABLE users DROP COLUMN IF EXISTS csrf_token;
