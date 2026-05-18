ALTER TABLE workspace
    ADD COLUMN artifact_counter INT NOT NULL DEFAULT 0;

ALTER TABLE artifact
    ADD COLUMN number INT NOT NULL DEFAULT 0;

WITH numbered AS (
    SELECT id, workspace_id,
           ROW_NUMBER() OVER (PARTITION BY workspace_id ORDER BY created_at ASC) AS rn
    FROM artifact
)
UPDATE artifact SET number = numbered.rn
FROM numbered WHERE artifact.id = numbered.id;

UPDATE workspace SET artifact_counter = COALESCE(
    (SELECT MAX(number) FROM artifact WHERE artifact.workspace_id = workspace.id), 0
);

ALTER TABLE artifact ADD CONSTRAINT uq_artifact_workspace_number UNIQUE (workspace_id, number);

CREATE INDEX idx_artifact_workspace_number ON artifact(workspace_id, number);
