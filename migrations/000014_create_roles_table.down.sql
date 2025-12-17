-- Remove foreign key constraint
ALTER TABLE users DROP CONSTRAINT fk_users_role;

-- Drop index
DROP INDEX idx_users_role_id ON users;

-- Remove role_id column from users table
ALTER TABLE users DROP COLUMN role_id;

-- Drop roles table
DROP TABLE roles;
