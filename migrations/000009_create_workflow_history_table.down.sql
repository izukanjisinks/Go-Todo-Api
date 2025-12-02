DROP INDEX IF EXISTS idx_workflow_history_performed_by ON workflow_history;
DROP INDEX IF EXISTS idx_workflow_history_timestamp ON workflow_history;
DROP INDEX IF EXISTS idx_workflow_history_instance_id ON workflow_history;
DROP TABLE IF EXISTS workflow_history;
