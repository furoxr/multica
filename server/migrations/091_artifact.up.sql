CREATE TABLE artifact (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    workspace_id UUID NOT NULL REFERENCES workspace(id) ON DELETE CASCADE,
    project_id UUID REFERENCES project(id) ON DELETE SET NULL,
    title TEXT NOT NULL,
    summary TEXT NOT NULL DEFAULT '',
    content TEXT NOT NULL DEFAULT '',
    content_type TEXT NOT NULL DEFAULT 'text/markdown',
    creator_type TEXT NOT NULL CHECK (creator_type IN ('member', 'agent')),
    creator_id UUID NOT NULL,
    origin_issue_id UUID REFERENCES issue(id) ON DELETE SET NULL,
    origin_task_id UUID REFERENCES agent_task_queue(id) ON DELETE SET NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (workspace_id, id)
);

CREATE INDEX idx_artifact_workspace ON artifact(workspace_id);
CREATE INDEX idx_artifact_project ON artifact(project_id) WHERE project_id IS NOT NULL;
CREATE INDEX idx_artifact_origin_issue ON artifact(origin_issue_id) WHERE origin_issue_id IS NOT NULL;
