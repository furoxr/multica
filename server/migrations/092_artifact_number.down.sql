DROP INDEX IF EXISTS idx_artifact_workspace_number;
ALTER TABLE artifact DROP CONSTRAINT IF EXISTS uq_artifact_workspace_number;
ALTER TABLE artifact DROP COLUMN IF EXISTS number;
ALTER TABLE workspace DROP COLUMN IF EXISTS artifact_counter;
