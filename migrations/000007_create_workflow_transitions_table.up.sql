CREATE TABLE workflow_transitions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    workflow_id UUID NOT NULL,
    from_step_id UUID NOT NULL,
    to_step_id UUID NOT NULL,
    action_name VARCHAR(100) NOT NULL,
    condition_type VARCHAR(50), -- e.g., 'user_role', 'field_value', 'assigned_user_only'
    condition_value TEXT, -- JSON for complex conditions
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_workflow_transitions_workflow_id FOREIGN KEY (workflow_id) REFERENCES workflows(id) ON DELETE CASCADE,
    CONSTRAINT fk_workflow_transitions_from_step FOREIGN KEY (from_step_id) REFERENCES workflow_steps(id),
    CONSTRAINT fk_workflow_transitions_to_step FOREIGN KEY (to_step_id) REFERENCES workflow_steps(id)
);

CREATE INDEX idx_workflow_transitions_workflow_id ON workflow_transitions(workflow_id);
CREATE INDEX idx_workflow_transitions_from_step ON workflow_transitions(from_step_id);
CREATE INDEX idx_workflow_transitions_action ON workflow_transitions(action_name);
