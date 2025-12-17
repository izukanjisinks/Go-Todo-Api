-- Drop existing workflow_instances table and recreate as assigned_todos
-- First, drop the foreign key constraint from workflow_history
ALTER TABLE workflow_history DROP CONSTRAINT IF EXISTS fk_workflow_history_instance_id;

-- Delete all workflow history data (will be lost)
DELETE FROM workflow_history;

-- Drop assigned_todos if it exists (from partial migration)
DROP INDEX IF EXISTS idx_assigned_todos_assigned_to ON assigned_todos;
DROP INDEX IF EXISTS idx_assigned_todos_todo_id ON assigned_todos;
DROP INDEX IF EXISTS idx_assigned_todos_current_step ON assigned_todos;
DROP INDEX IF EXISTS idx_assigned_todos_workflow_id ON assigned_todos;
DROP TABLE IF EXISTS assigned_todos;

-- Drop indexes and table
DROP INDEX IF EXISTS idx_workflow_instances_created_by ON workflow_instances;
DROP INDEX IF EXISTS idx_workflow_instances_assigned_to ON workflow_instances;
DROP INDEX IF EXISTS idx_workflow_instances_current_step ON workflow_instances;
DROP INDEX IF EXISTS idx_workflow_instances_workflow_id ON workflow_instances;
DROP TABLE IF EXISTS workflow_instances;

-- Create new assigned_todos table (represents workflow instances)
CREATE TABLE assigned_todos (
    id NVARCHAR(36) PRIMARY KEY,
    workflow_id NVARCHAR(36) NOT NULL,
    current_step_id NVARCHAR(36) NOT NULL,
    todo_id NVARCHAR(36) NOT NULL,
    assigned_to NVARCHAR(100) NOT NULL,
    created_at DATETIME2 DEFAULT GETDATE(),
    updated_at DATETIME2 DEFAULT GETDATE(),
    CONSTRAINT fk_assigned_todos_workflow_id FOREIGN KEY (workflow_id) REFERENCES workflows(id),
    CONSTRAINT fk_assigned_todos_current_step FOREIGN KEY (current_step_id) REFERENCES workflow_steps(id),
    CONSTRAINT fk_assigned_todos_todo_id FOREIGN KEY (todo_id) REFERENCES todos(id)
);

CREATE INDEX idx_assigned_todos_workflow_id ON assigned_todos(workflow_id);
CREATE INDEX idx_assigned_todos_current_step ON assigned_todos(current_step_id);
CREATE INDEX idx_assigned_todos_todo_id ON assigned_todos(todo_id);
CREATE INDEX idx_assigned_todos_assigned_to ON assigned_todos(assigned_to);

-- Re-add the foreign key constraint from workflow_history to the new table
ALTER TABLE workflow_history 
ADD CONSTRAINT fk_workflow_history_instance_id 
FOREIGN KEY (instance_id) REFERENCES assigned_todos(id) ON DELETE CASCADE;
