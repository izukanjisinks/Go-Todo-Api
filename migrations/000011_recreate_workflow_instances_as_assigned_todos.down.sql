-- Rollback: Drop assigned_todos and recreate workflow_instances
-- First, drop the foreign key constraint from workflow_history
ALTER TABLE workflow_history DROP CONSTRAINT IF EXISTS fk_workflow_history_instance_id;

-- Drop indexes and table
DROP INDEX IF EXISTS idx_assigned_todos_assigned_to;
DROP INDEX IF EXISTS idx_assigned_todos_todo_id;
DROP INDEX IF EXISTS idx_assigned_todos_current_step;
DROP INDEX IF EXISTS idx_assigned_todos_workflow_id;
DROP TABLE IF EXISTS assigned_todos;

-- Recreate original workflow_instances table
CREATE TABLE workflow_instances (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    workflow_id UUID NOT NULL,
    current_step_id UUID NOT NULL,
    title VARCHAR(200) NOT NULL,
    description VARCHAR(500),
    task_data TEXT,
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

-- Re-add the foreign key constraint from workflow_history to workflow_instances
ALTER TABLE workflow_history
ADD CONSTRAINT fk_workflow_history_instance_id
FOREIGN KEY (instance_id) REFERENCES workflow_instances(id) ON DELETE CASCADE;
