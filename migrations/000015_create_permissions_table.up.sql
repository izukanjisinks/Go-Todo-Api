-- Create permissions table
CREATE TABLE permissions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) UNIQUE NOT NULL,
    description TEXT,
    "view" BOOLEAN NOT NULL DEFAULT FALSE,
    "create" BOOLEAN NOT NULL DEFAULT FALSE,
    "update" BOOLEAN NOT NULL DEFAULT FALSE,
    "delete" BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Insert predefined permissions
INSERT INTO permissions (id, name, description, "view", "create", "update", "delete") VALUES
    (gen_random_uuid(), 'super_admin_permissions', 'Full access to all operations', TRUE, TRUE, TRUE, TRUE),
    (gen_random_uuid(), 'admin_permissions', 'Admin level access', TRUE, TRUE, TRUE, TRUE),
    (gen_random_uuid(), 'moderator_permissions', 'Moderator level access', TRUE, TRUE, TRUE, FALSE),
    (gen_random_uuid(), 'user_permissions', 'Basic user access', TRUE, FALSE, FALSE, FALSE);
