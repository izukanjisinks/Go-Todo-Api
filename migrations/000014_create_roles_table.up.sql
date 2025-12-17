-- Create roles table
CREATE TABLE roles (
    role_id UNIQUEIDENTIFIER PRIMARY KEY DEFAULT NEWID(),
    name NVARCHAR(100) UNIQUE NOT NULL,
    description NVARCHAR(MAX),
    permissions NVARCHAR(MAX) NOT NULL DEFAULT '[]',
    created_at DATETIME2 DEFAULT GETDATE(),
    updated_at DATETIME2 DEFAULT GETDATE()
);

-- Add role_id column to users table
ALTER TABLE users ADD role_id UNIQUEIDENTIFIER NULL;

-- Add foreign key constraint
ALTER TABLE users ADD CONSTRAINT fk_users_role
    FOREIGN KEY (role_id) REFERENCES roles(role_id) ON DELETE SET NULL;

-- Create index on role_id for better query performance
CREATE INDEX idx_users_role_id ON users(role_id);

-- Insert predefined roles (using JSON for permissions)
INSERT INTO roles (role_id, name, description, permissions) VALUES
    (NEWID(), 'Super Admin', 'Full system access with all permissions', 
     '["users:read","users:create","users:update","users:delete","content:read","content:create","content:update","content:delete","settings:read","settings:update","reports:view","reports:export","roles:manage"]');

INSERT INTO roles (role_id, name, description, permissions) VALUES
    (NEWID(), 'Admin', 'Administrative access with most permissions',
     '["users:read","users:create","users:update","users:delete","content:read","content:create","content:update","content:delete","settings:read","reports:view","reports:export"]');

INSERT INTO roles (role_id, name, description, permissions) VALUES
    (NEWID(), 'Moderator', 'Content management and user viewing access',
     '["users:read","content:read","content:create","content:update","content:delete"]');

INSERT INTO roles (role_id, name, description, permissions) VALUES
    (NEWID(), 'User', 'Basic user access with read permissions',
     '["content:read","reports:view"]');
