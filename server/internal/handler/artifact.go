package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgtype"
	db "github.com/multica-ai/multica/server/pkg/db/generated"
)

type ArtifactResponse struct {
	ID            string  `json:"id"`
	WorkspaceID   string  `json:"workspace_id"`
	ProjectID     *string `json:"project_id"`
	Title         string  `json:"title"`
	Summary       string  `json:"summary"`
	Content       string  `json:"content"`
	ContentType   string  `json:"content_type"`
	CreatorType   string  `json:"creator_type"`
	CreatorID     string  `json:"creator_id"`
	OriginIssueID *string `json:"origin_issue_id"`
	OriginTaskID  *string `json:"origin_task_id"`
	CreatedAt     string  `json:"created_at"`
	UpdatedAt     string  `json:"updated_at"`
}

type ArtifactSummaryResponse struct {
	ID            string  `json:"id"`
	WorkspaceID   string  `json:"workspace_id"`
	ProjectID     *string `json:"project_id"`
	Title         string  `json:"title"`
	Summary       string  `json:"summary"`
	ContentType   string  `json:"content_type"`
	CreatorType   string  `json:"creator_type"`
	CreatorID     string  `json:"creator_id"`
	OriginIssueID *string `json:"origin_issue_id"`
	OriginTaskID  *string `json:"origin_task_id"`
	CreatedAt     string  `json:"created_at"`
	UpdatedAt     string  `json:"updated_at"`
}

type CreateArtifactRequest struct {
	ProjectID     *string `json:"project_id"`
	Title         string  `json:"title"`
	Summary       string  `json:"summary"`
	Content       string  `json:"content"`
	ContentType   string  `json:"content_type"`
	OriginIssueID *string `json:"origin_issue_id"`
	OriginTaskID  *string `json:"origin_task_id"`
}

type UpdateArtifactRequest struct {
	ProjectID     *string `json:"project_id"`
	Title         *string `json:"title"`
	Summary       *string `json:"summary"`
	Content       *string `json:"content"`
	ContentType   *string `json:"content_type"`
	OriginIssueID *string `json:"origin_issue_id"`
	OriginTaskID  *string `json:"origin_task_id"`
}

func artifactToResponse(a db.Artifact) ArtifactResponse {
	return ArtifactResponse{
		ID:            uuidToString(a.ID),
		WorkspaceID:   uuidToString(a.WorkspaceID),
		ProjectID:     uuidToPtr(a.ProjectID),
		Title:         a.Title,
		Summary:       a.Summary,
		Content:       a.Content,
		ContentType:   a.ContentType,
		CreatorType:   a.CreatorType,
		CreatorID:     uuidToString(a.CreatorID),
		OriginIssueID: uuidToPtr(a.OriginIssueID),
		OriginTaskID:  uuidToPtr(a.OriginTaskID),
		CreatedAt:     timestampToString(a.CreatedAt),
		UpdatedAt:     timestampToString(a.UpdatedAt),
	}
}

func artifactSummaryToResponse(
	id, workspaceID, projectID pgtype.UUID,
	title, summary, contentType, creatorType string,
	creatorID, originIssueID, originTaskID pgtype.UUID,
	createdAt, updatedAt pgtype.Timestamptz,
) ArtifactSummaryResponse {
	return ArtifactSummaryResponse{
		ID:            uuidToString(id),
		WorkspaceID:   uuidToString(workspaceID),
		ProjectID:     uuidToPtr(projectID),
		Title:         title,
		Summary:       summary,
		ContentType:   contentType,
		CreatorType:   creatorType,
		CreatorID:     uuidToString(creatorID),
		OriginIssueID: uuidToPtr(originIssueID),
		OriginTaskID:  uuidToPtr(originTaskID),
		CreatedAt:     timestampToString(createdAt),
		UpdatedAt:     timestampToString(updatedAt),
	}
}

func optionalUUIDOrBadRequest(w http.ResponseWriter, value *string, fieldName string) (pgtype.UUID, bool) {
	if value == nil || strings.TrimSpace(*value) == "" {
		return pgtype.UUID{}, true
	}
	return parseUUIDOrBadRequest(w, strings.TrimSpace(*value), fieldName)
}

func optionalText(value *string) pgtype.Text {
	if value == nil {
		return pgtype.Text{}
	}
	return strToText(sanitizeNullBytes(*value))
}

func normalizeArtifactContentType(v string) string {
	v = strings.TrimSpace(strings.ToLower(v))
	if v == "" {
		return "text/markdown"
	}
	switch v {
	case "text/markdown", "text/plain", "application/json":
		return v
	default:
		return ""
	}
}

func (h *Handler) ListArtifacts(w http.ResponseWriter, r *http.Request) {
	workspaceID := h.resolveWorkspaceID(r)
	if workspaceID == "" {
		writeError(w, http.StatusBadRequest, "workspace_id is required")
		return
	}

	originIssueID := r.URL.Query().Get("origin_issue_id")
	if originIssueID != "" {
		originIssueUUID, ok := parseUUIDOrBadRequest(w, originIssueID, "origin_issue_id")
		if !ok {
			return
		}
		items, err := h.Queries.ListArtifactSummariesByOriginIssue(r.Context(), db.ListArtifactSummariesByOriginIssueParams{
			WorkspaceID:   parseUUID(workspaceID),
			OriginIssueID: originIssueUUID,
		})
		if err != nil {
			writeError(w, http.StatusInternalServerError, "failed to list artifacts")
			return
		}
		resp := make([]ArtifactSummaryResponse, len(items))
		for i, a := range items {
			resp[i] = artifactSummaryToResponse(a.ID, a.WorkspaceID, a.ProjectID, a.Title, a.Summary, a.ContentType, a.CreatorType, a.CreatorID, a.OriginIssueID, a.OriginTaskID, a.CreatedAt, a.UpdatedAt)
		}
		writeJSON(w, http.StatusOK, resp)
		return
	}

	items, err := h.Queries.ListArtifactSummariesByWorkspace(r.Context(), parseUUID(workspaceID))
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to list artifacts")
		return
	}
	resp := make([]ArtifactSummaryResponse, len(items))
	for i, a := range items {
		resp[i] = artifactSummaryToResponse(a.ID, a.WorkspaceID, a.ProjectID, a.Title, a.Summary, a.ContentType, a.CreatorType, a.CreatorID, a.OriginIssueID, a.OriginTaskID, a.CreatedAt, a.UpdatedAt)
	}
	writeJSON(w, http.StatusOK, resp)
}

func (h *Handler) GetArtifact(w http.ResponseWriter, r *http.Request) {
	workspaceID := h.resolveWorkspaceID(r)
	if workspaceID == "" {
		writeError(w, http.StatusBadRequest, "workspace_id is required")
		return
	}
	artifactID, ok := parseUUIDOrBadRequest(w, chi.URLParam(r, "id"), "id")
	if !ok {
		return
	}

	artifact, err := h.Queries.GetArtifactInWorkspace(r.Context(), db.GetArtifactInWorkspaceParams{
		ID:          artifactID,
		WorkspaceID: parseUUID(workspaceID),
	})
	if err != nil {
		writeError(w, http.StatusNotFound, "artifact not found")
		return
	}
	writeJSON(w, http.StatusOK, artifactToResponse(artifact))
}

func (h *Handler) CreateArtifact(w http.ResponseWriter, r *http.Request) {
	workspaceID := h.resolveWorkspaceID(r)
	if workspaceID == "" {
		writeError(w, http.StatusBadRequest, "workspace_id is required")
		return
	}
	userID, ok := requireUserID(w, r)
	if !ok {
		return
	}

	var req CreateArtifactRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	req.Title = strings.TrimSpace(req.Title)
	if req.Title == "" {
		writeError(w, http.StatusBadRequest, "title is required")
		return
	}
	contentType := normalizeArtifactContentType(req.ContentType)
	if contentType == "" {
		writeError(w, http.StatusBadRequest, "unsupported content_type")
		return
	}
	projectID, ok := optionalUUIDOrBadRequest(w, req.ProjectID, "project_id")
	if !ok {
		return
	}
	originIssueID, ok := optionalUUIDOrBadRequest(w, req.OriginIssueID, "origin_issue_id")
	if !ok {
		return
	}
	originTaskID, ok := optionalUUIDOrBadRequest(w, req.OriginTaskID, "origin_task_id")
	if !ok {
		return
	}

	actorType, actorID := h.resolveActor(r, userID, workspaceID)
	artifact, err := h.Queries.CreateArtifact(r.Context(), db.CreateArtifactParams{
		WorkspaceID:   parseUUID(workspaceID),
		ProjectID:     projectID,
		Title:         sanitizeNullBytes(req.Title),
		Summary:       sanitizeNullBytes(req.Summary),
		Content:       sanitizeNullBytes(req.Content),
		ContentType:   contentType,
		CreatorType:   actorType,
		CreatorID:     parseUUID(actorID),
		OriginIssueID: originIssueID,
		OriginTaskID:  originTaskID,
	})
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to create artifact")
		return
	}

	writeJSON(w, http.StatusCreated, artifactToResponse(artifact))
}

func (h *Handler) UpdateArtifact(w http.ResponseWriter, r *http.Request) {
	workspaceID := h.resolveWorkspaceID(r)
	if workspaceID == "" {
		writeError(w, http.StatusBadRequest, "workspace_id is required")
		return
	}
	artifactID, ok := parseUUIDOrBadRequest(w, chi.URLParam(r, "id"), "id")
	if !ok {
		return
	}

	var req UpdateArtifactRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	contentType := pgtype.Text{}
	if req.ContentType != nil {
		normalized := normalizeArtifactContentType(*req.ContentType)
		if normalized == "" {
			writeError(w, http.StatusBadRequest, "unsupported content_type")
			return
		}
		contentType = strToText(normalized)
	}
	projectID, ok := optionalUUIDOrBadRequest(w, req.ProjectID, "project_id")
	if !ok {
		return
	}
	originIssueID, ok := optionalUUIDOrBadRequest(w, req.OriginIssueID, "origin_issue_id")
	if !ok {
		return
	}
	originTaskID, ok := optionalUUIDOrBadRequest(w, req.OriginTaskID, "origin_task_id")
	if !ok {
		return
	}

	artifact, err := h.Queries.UpdateArtifact(r.Context(), db.UpdateArtifactParams{
		ID:            artifactID,
		WorkspaceID:   parseUUID(workspaceID),
		ProjectID:     projectID,
		Title:         optionalText(req.Title),
		Summary:       optionalText(req.Summary),
		Content:       optionalText(req.Content),
		ContentType:   contentType,
		OriginIssueID: originIssueID,
		OriginTaskID:  originTaskID,
	})
	if err != nil {
		writeError(w, http.StatusNotFound, "artifact not found")
		return
	}

	writeJSON(w, http.StatusOK, artifactToResponse(artifact))
}

func (h *Handler) DeleteArtifact(w http.ResponseWriter, r *http.Request) {
	workspaceID := h.resolveWorkspaceID(r)
	if workspaceID == "" {
		writeError(w, http.StatusBadRequest, "workspace_id is required")
		return
	}
	artifactID, ok := parseUUIDOrBadRequest(w, chi.URLParam(r, "id"), "id")
	if !ok {
		return
	}

	err := h.Queries.DeleteArtifact(r.Context(), db.DeleteArtifactParams{
		ID:          artifactID,
		WorkspaceID: parseUUID(workspaceID),
	})
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to delete artifact")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
