-- Add is_active column to users table
ALTER TABLE users ADD is_active BOOLEAN NOT NULL DEFAULT TRUE;

-- Create index for faster queries on active users
CREATE INDEX idx_users_is_active ON users(is_active);
