CREATE TABLE workflows (
    id NVARCHAR(36) PRIMARY KEY,
    name NVARCHAR(200) NOT NULL,
    description NVARCHAR(500),
    is_active BIT DEFAULT 1,
    created_by NVARCHAR(100) NOT NULL,
    created_at DATETIME2 DEFAULT GETDATE(),
    updated_at DATETIME2 DEFAULT GETDATE()
);

CREATE INDEX idx_workflows_is_active ON workflows(is_active);
CREATE INDEX idx_workflows_created_by ON workflows(created_by);
