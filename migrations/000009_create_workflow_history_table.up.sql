CREATE TABLE workflow_history (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    instance_id UUID NOT NULL,
    from_step_id UUID,
    to_step_id UUID NOT NULL,
    action_taken VARCHAR(100) NOT NULL,
    performed_by VARCHAR(100) NOT NULL,
    comments VARCHAR(500),
    timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_workflow_history_instance_id FOREIGN KEY (instance_id) REFERENCES workflow_instances(id) ON DELETE CASCADE,
    CONSTRAINT fk_workflow_history_from_step FOREIGN KEY (from_step_id) REFERENCES workflow_steps(id),
    CONSTRAINT fk_workflow_history_to_step FOREIGN KEY (to_step_id) REFERENCES workflow_steps(id)
);

CREATE INDEX idx_workflow_history_instance_id ON workflow_history(instance_id);
CREATE INDEX idx_workflow_history_timestamp ON workflow_history(timestamp);
CREATE INDEX idx_workflow_history_performed_by ON workflow_history(performed_by);
