DROP INDEX IF EXISTS idx_workflow_transitions_action ON workflow_transitions;
DROP INDEX IF EXISTS idx_workflow_transitions_from_step ON workflow_transitions;
DROP INDEX IF EXISTS idx_workflow_transitions_workflow_id ON workflow_transitions;
DROP TABLE IF EXISTS workflow_transitions;
