-- Artifact CRUD

-- name: ListArtifactSummariesByWorkspace :many
SELECT id, workspace_id, project_id, title, summary, content_type, creator_type, creator_id, origin_issue_id, origin_task_id, number, created_at, updated_at
FROM artifact
WHERE workspace_id = $1
ORDER BY updated_at DESC;

-- name: ListArtifactSummariesByOriginIssue :many
SELECT id, workspace_id, project_id, title, summary, content_type, creator_type, creator_id, origin_issue_id, origin_task_id, number, created_at, updated_at
FROM artifact
WHERE workspace_id = $1 AND origin_issue_id = $2
ORDER BY updated_at DESC;

-- name: GetArtifactInWorkspace :one
SELECT * FROM artifact
WHERE id = $1 AND workspace_id = $2;

-- name: GetArtifactByNumber :one
SELECT * FROM artifact
WHERE workspace_id = $1 AND number = $2;

-- name: CreateArtifact :one
INSERT INTO artifact (
    workspace_id,
    project_id,
    title,
    summary,
    content,
    content_type,
    creator_type,
    creator_id,
    origin_issue_id,
    origin_task_id,
    number
)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
RETURNING *;

-- name: UpdateArtifact :one
UPDATE artifact SET
    project_id = COALESCE(sqlc.narg('project_id'), project_id),
    title = COALESCE(sqlc.narg('title'), title),
    summary = COALESCE(sqlc.narg('summary'), summary),
    content = COALESCE(sqlc.narg('content'), content),
    content_type = COALESCE(sqlc.narg('content_type'), content_type),
    origin_issue_id = COALESCE(sqlc.narg('origin_issue_id'), origin_issue_id),
    origin_task_id = COALESCE(sqlc.narg('origin_task_id'), origin_task_id),
    updated_at = now()
WHERE id = $1 AND workspace_id = $2
RETURNING *;

-- name: DeleteArtifact :exec
DELETE FROM artifact
WHERE id = $1 AND workspace_id = $2;
