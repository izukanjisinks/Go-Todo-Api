CREATE TABLE workflow_instances (
    id NVARCHAR(36) PRIMARY KEY,
    workflow_id NVARCHAR(36) NOT NULL,
    current_step_id NVARCHAR(36) NOT NULL,
    title NVARCHAR(200) NOT NULL,
    description NVARCHAR(500),
    task_data NVARCHAR(MAX), -- JSON for additional fields
    assigned_to NVARCHAR(100) NOT NULL,
    created_by NVARCHAR(100) NOT NULL,
    created_at DATETIME2 DEFAULT GETDATE(),
    updated_at DATETIME2 DEFAULT GETDATE(),
    CONSTRAINT fk_workflow_instances_workflow_id FOREIGN KEY (workflow_id) REFERENCES workflows(id),
    CONSTRAINT fk_workflow_instances_current_step FOREIGN KEY (current_step_id) REFERENCES workflow_steps(id)
);

CREATE INDEX idx_workflow_instances_workflow_id ON workflow_instances(workflow_id);
CREATE INDEX idx_workflow_instances_current_step ON workflow_instances(current_step_id);
CREATE INDEX idx_workflow_instances_assigned_to ON workflow_instances(assigned_to);
CREATE INDEX idx_workflow_instances_created_by ON workflow_instances(created_by);
