DROP INDEX IF EXISTS idx_workflow_instances_created_by ON workflow_instances;
DROP INDEX IF EXISTS idx_workflow_instances_assigned_to ON workflow_instances;
DROP INDEX IF EXISTS idx_workflow_instances_current_step ON workflow_instances;
DROP INDEX IF EXISTS idx_workflow_instances_workflow_id ON workflow_instances;
DROP TABLE IF EXISTS workflow_instances;
