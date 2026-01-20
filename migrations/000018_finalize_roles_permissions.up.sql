-- Add foreign key constraint to permissions table
ALTER TABLE roles ADD CONSTRAINT fk_roles_permission FOREIGN KEY (permission_id) REFERENCES permissions(id) ON DELETE SET NULL;

-- Create index on permission_id for better query performance
CREATE INDEX idx_roles_permission_id ON roles(permission_id);
