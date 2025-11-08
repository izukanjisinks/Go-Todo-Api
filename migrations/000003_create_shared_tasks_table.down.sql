DROP INDEX IF EXISTS idx_sharedtasks_owner ON shared_tasks;
DROP INDEX IF EXISTS idx_sharedtasks_shared_with ON shared_tasks;
DROP INDEX IF EXISTS idx_sharedtasks_todo ON shared_tasks;
DROP TABLE IF EXISTS shared_tasks;
