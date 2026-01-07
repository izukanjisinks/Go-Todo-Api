-- Create permissions table
CREATE TABLE permissions (
    id UNIQUEIDENTIFIER PRIMARY KEY DEFAULT NEWID(),
    name NVARCHAR(100) UNIQUE NOT NULL,
    description NVARCHAR(MAX),
    [view] BIT NOT NULL DEFAULT 0,
    [create] BIT NOT NULL DEFAULT 0,
    [update] BIT NOT NULL DEFAULT 0,
    [delete] BIT NOT NULL DEFAULT 0,
    created_at DATETIME2 DEFAULT GETDATE(),
    updated_at DATETIME2 DEFAULT GETDATE()
);

-- Insert predefined permissions
INSERT INTO permissions (id, name, description, [view], [create], [update], [delete]) VALUES
    (NEWID(), 'super_admin_permissions', 'Full access to all operations', 1, 1, 1, 1),
    (NEWID(), 'admin_permissions', 'Admin level access', 1, 1, 1, 1),
    (NEWID(), 'moderator_permissions', 'Moderator level access', 1, 1, 1, 0),
    (NEWID(), 'user_permissions', 'Basic user access', 1, 0, 0, 0);
