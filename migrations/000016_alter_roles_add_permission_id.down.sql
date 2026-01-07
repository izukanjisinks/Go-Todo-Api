-- Add back the permissions column (JSON)
ALTER TABLE roles ADD permissions NVARCHAR(MAX) NOT NULL DEFAULT '[]';

-- Restore permissions data based on permission_id (best effort)
UPDATE roles
SET permissions = '["users:read","users:create","users:update","users:delete","content:read","content:create","content:update","content:delete","settings:read","settings:update","reports:view","reports:export","roles:manage"]'
WHERE name = 'Super Admin';

UPDATE roles
SET permissions = '["users:read","users:create","users:update","users:delete","content:read","content:create","content:update","content:delete","settings:read","reports:view","reports:export"]'
WHERE name = 'Admin';

UPDATE roles
SET permissions = '["users:read","content:read","content:create","content:update","content:delete"]'
WHERE name = 'Moderator';

UPDATE roles
SET permissions = '["content:read","reports:view"]'
WHERE name = 'User';

-- Drop index
DROP INDEX idx_roles_permission_id ON roles;

-- Remove foreign key constraint
ALTER TABLE roles DROP CONSTRAINT fk_roles_permission;

-- Remove permission_id column from roles table
ALTER TABLE roles DROP COLUMN permission_id;
