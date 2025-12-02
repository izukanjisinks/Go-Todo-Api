CREATE TABLE workflow_history (
    id NVARCHAR(36) PRIMARY KEY,
    instance_id NVARCHAR(36) NOT NULL,
    from_step_id NVARCHAR(36),
    to_step_id NVARCHAR(36) NOT NULL,
    action_taken NVARCHAR(100) NOT NULL,
    performed_by NVARCHAR(100) NOT NULL,
    comments NVARCHAR(500),
    timestamp DATETIME2 DEFAULT GETDATE(),
    CONSTRAINT fk_workflow_history_instance_id FOREIGN KEY (instance_id) REFERENCES workflow_instances(id) ON DELETE CASCADE,
    CONSTRAINT fk_workflow_history_from_step FOREIGN KEY (from_step_id) REFERENCES workflow_steps(id),
    CONSTRAINT fk_workflow_history_to_step FOREIGN KEY (to_step_id) REFERENCES workflow_steps(id)
);

CREATE INDEX idx_workflow_history_instance_id ON workflow_history(instance_id);
CREATE INDEX idx_workflow_history_timestamp ON workflow_history(timestamp);
CREATE INDEX idx_workflow_history_performed_by ON workflow_history(performed_by);
