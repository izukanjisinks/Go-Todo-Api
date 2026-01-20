CREATE TABLE workflow_instances (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    workflow_id UUID NOT NULL,
    current_step_id UUID NOT NULL,
    title VARCHAR(200) NOT NULL,
    description VARCHAR(500),
    task_data TEXT, -- JSON for additional fields
    assigned_to VARCHAR(100) NOT NULL,
    created_by VARCHAR(100) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_workflow_instances_workflow_id FOREIGN KEY (workflow_id) REFERENCES workflows(id),
    CONSTRAINT fk_workflow_instances_current_step FOREIGN KEY (current_step_id) REFERENCES workflow_steps(id)
);

CREATE INDEX idx_workflow_instances_workflow_id ON workflow_instances(workflow_id);
CREATE INDEX idx_workflow_instances_current_step ON workflow_instances(current_step_id);
CREATE INDEX idx_workflow_instances_assigned_to ON workflow_instances(assigned_to);
CREATE INDEX idx_workflow_instances_created_by ON workflow_instances(created_by);
