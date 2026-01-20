CREATE TABLE workflow_steps (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    workflow_id UUID NOT NULL,
    step_name VARCHAR(100) NOT NULL,
    step_order INT NOT NULL,
    initial BOOLEAN DEFAULT FALSE,
    final BOOLEAN DEFAULT FALSE,
    allowed_roles TEXT, -- JSON array of roles
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_workflow_steps_workflow_id FOREIGN KEY (workflow_id) REFERENCES workflows(id) ON DELETE CASCADE
);

CREATE INDEX idx_workflow_steps_workflow_id ON workflow_steps(workflow_id);
CREATE INDEX idx_workflow_steps_initial ON workflow_steps(initial);
