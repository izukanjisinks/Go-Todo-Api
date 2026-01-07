-- Add permission_id column to roles table
ALTER TABLE roles ADD permission_id UNIQUEIDENTIFIER NULL;

-- Add foreign key constraint to permissions table
ALTER TABLE roles ADD CONSTRAINT fk_roles_permission
    FOREIGN KEY (permission_id) REFERENCES permissions(id) ON DELETE SET NULL;

-- Create index on permission_id for better query performance
CREATE INDEX idx_roles_permission_id ON roles(permission_id);

-- Update existing roles to link to their corresponding permissions
UPDATE roles
SET permission_id = (SELECT id FROM permissions WHERE name = 'super_admin_permissions')
WHERE name = 'Super Admin';

UPDATE roles
SET permission_id = (SELECT id FROM permissions WHERE name = 'admin_permissions')
WHERE name = 'Admin';

UPDATE roles
SET permission_id = (SELECT id FROM permissions WHERE name = 'moderator_permissions')
WHERE name = 'Moderator';

UPDATE roles
SET permission_id = (SELECT id FROM permissions WHERE name = 'user_permissions')
WHERE name = 'User';

-- Drop the old permissions column (JSON)
ALTER TABLE roles DROP COLUMN permissions;
