-- Update existing roles to link to their corresponding permissions
UPDATE roles SET permission_id = (SELECT id FROM permissions WHERE name = 'super_admin_permissions') WHERE name = 'Super Admin';
UPDATE roles SET permission_id = (SELECT id FROM permissions WHERE name = 'admin_permissions') WHERE name = 'Admin';
UPDATE roles SET permission_id = (SELECT id FROM permissions WHERE name = 'moderator_permissions') WHERE name = 'Moderator';
UPDATE roles SET permission_id = (SELECT id FROM permissions WHERE name = 'user_permissions') WHERE name = 'User';
