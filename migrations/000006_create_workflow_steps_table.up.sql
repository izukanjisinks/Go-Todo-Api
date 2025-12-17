CREATE TABLE workflow_steps (
    id NVARCHAR(36) PRIMARY KEY,
    workflow_id NVARCHAR(36) NOT NULL,
    step_name NVARCHAR(100) NOT NULL,
    step_order INT NOT NULL,
    initial BIT DEFAULT 0,
    final BIT DEFAULT 0,
    allowed_roles NVARCHAR(MAX), -- JSON array of roles
    created_at DATETIME2 DEFAULT GETDATE(),
    CONSTRAINT fk_workflow_steps_workflow_id FOREIGN KEY (workflow_id) REFERENCES workflows(id) ON DELETE CASCADE
);

CREATE INDEX idx_workflow_steps_workflow_id ON workflow_steps(workflow_id);
CREATE INDEX idx_workflow_steps_initial ON workflow_steps(initial);
