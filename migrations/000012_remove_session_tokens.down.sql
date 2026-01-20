-- Rollback: Add back session_token and csrf_token columns
-- In case we need to revert to session-based authentication

ALTER TABLE users ADD session_token VARCHAR(255);
ALTER TABLE users ADD csrf_token VARCHAR(255);

-- Recreate the index
CREATE INDEX idx_users_session_token ON users(session_token);
