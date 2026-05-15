export interface ArtifactSummary {
  id: string;
  workspace_id: string;
  project_id: string | null;
  title: string;
  summary: string;
  content_type: string;
  creator_type: "member" | "agent";
  creator_id: string;
  origin_issue_id: string | null;
  origin_task_id: string | null;
  created_at: string;
  updated_at: string;
}

export interface Artifact extends ArtifactSummary {
  content: string;
}

export interface CreateArtifactRequest {
  title: string;
  summary?: string;
  content?: string;
  content_type?: string;
  project_id?: string;
  origin_issue_id?: string;
  origin_task_id?: string;
}

export interface UpdateArtifactRequest {
  title?: string;
  summary?: string;
  content?: string;
  content_type?: string;
  project_id?: string;
  origin_issue_id?: string;
  origin_task_id?: string;
}
