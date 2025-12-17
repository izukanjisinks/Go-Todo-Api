CREATE TABLE workflow_transitions (
    id NVARCHAR(36) PRIMARY KEY,
    workflow_id NVARCHAR(36) NOT NULL,
    from_step_id NVARCHAR(36) NOT NULL,
    to_step_id NVARCHAR(36) NOT NULL,
    action_name NVARCHAR(100) NOT NULL,
    condition_type NVARCHAR(50), -- e.g., 'user_role', 'field_value', 'assigned_user_only'
    condition_value NVARCHAR(MAX), -- JSON for complex conditions
    created_at DATETIME2 DEFAULT GETDATE(),
    CONSTRAINT fk_workflow_transitions_workflow_id FOREIGN KEY (workflow_id) REFERENCES workflows(id) ON DELETE CASCADE,
    CONSTRAINT fk_workflow_transitions_from_step FOREIGN KEY (from_step_id) REFERENCES workflow_steps(id),
    CONSTRAINT fk_workflow_transitions_to_step FOREIGN KEY (to_step_id) REFERENCES workflow_steps(id)
);

CREATE INDEX idx_workflow_transitions_workflow_id ON workflow_transitions(workflow_id);
CREATE INDEX idx_workflow_transitions_from_step ON workflow_transitions(from_step_id);
CREATE INDEX idx_workflow_transitions_action ON workflow_transitions(action_name);
