package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgtype"
	db "github.com/multica-ai/multica/server/pkg/db/generated"
	"github.com/multica-ai/multica/server/pkg/protocol"
)

type ArtifactResponse struct {
	ID            string  `json:"id"`
	WorkspaceID   string  `json:"workspace_id"`
	Number        int32   `json:"number"`
	Identifier    string  `json:"identifier"`
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
	Number        int32   `json:"number"`
	Identifier    string  `json:"identifier"`
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

func artifactIdentifier(prefix string, number int32) string {
	if prefix == "" {
		return ""
	}
	return prefix + "-D" + strconv.Itoa(int(number))
}

func artifactToResponse(a db.Artifact, prefix string) ArtifactResponse {
	return ArtifactResponse{
		ID:            uuidToString(a.ID),
		WorkspaceID:   uuidToString(a.WorkspaceID),
		Number:        a.Number,
		Identifier:    artifactIdentifier(prefix, a.Number),
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
	number int32,
	createdAt, updatedAt pgtype.Timestamptz,
	prefix string,
) ArtifactSummaryResponse {
	return ArtifactSummaryResponse{
		ID:            uuidToString(id),
		WorkspaceID:   uuidToString(workspaceID),
		Number:        number,
		Identifier:    artifactIdentifier(prefix, number),
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
	ctx := r.Context()
	workspaceID := h.resolveWorkspaceID(r)
	if workspaceID == "" {
		writeError(w, http.StatusBadRequest, "workspace_id is required")
		return
	}
	wsUUID := parseUUID(workspaceID)
	prefix := h.getIssuePrefix(ctx, wsUUID)

	originIssueID := r.URL.Query().Get("origin_issue_id")
	if originIssueID != "" {
		originIssueUUID, ok := parseUUIDOrBadRequest(w, originIssueID, "origin_issue_id")
		if !ok {
			return
		}
		items, err := h.Queries.ListArtifactSummariesByOriginIssue(ctx, db.ListArtifactSummariesByOriginIssueParams{
			WorkspaceID:   wsUUID,
			OriginIssueID: originIssueUUID,
		})
		if err != nil {
			writeError(w, http.StatusInternalServerError, "failed to list artifacts")
			return
		}
		resp := make([]ArtifactSummaryResponse, len(items))
		for i, a := range items {
			resp[i] = artifactSummaryToResponse(a.ID, a.WorkspaceID, a.ProjectID, a.Title, a.Summary, a.ContentType, a.CreatorType, a.CreatorID, a.OriginIssueID, a.OriginTaskID, a.Number, a.CreatedAt, a.UpdatedAt, prefix)
		}
		writeJSON(w, http.StatusOK, resp)
		return
	}

	items, err := h.Queries.ListArtifactSummariesByWorkspace(ctx, wsUUID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to list artifacts")
		return
	}
	resp := make([]ArtifactSummaryResponse, len(items))
	for i, a := range items {
		resp[i] = artifactSummaryToResponse(a.ID, a.WorkspaceID, a.ProjectID, a.Title, a.Summary, a.ContentType, a.CreatorType, a.CreatorID, a.OriginIssueID, a.OriginTaskID, a.Number, a.CreatedAt, a.UpdatedAt, prefix)
	}
	writeJSON(w, http.StatusOK, resp)
}

func (h *Handler) GetArtifact(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	workspaceID := h.resolveWorkspaceID(r)
	if workspaceID == "" {
		writeError(w, http.StatusBadRequest, "workspace_id is required")
		return
	}
	artifactID, ok := parseUUIDOrBadRequest(w, chi.URLParam(r, "id"), "id")
	if !ok {
		return
	}
	wsUUID := parseUUID(workspaceID)

	artifact, err := h.Queries.GetArtifactInWorkspace(ctx, db.GetArtifactInWorkspaceParams{
		ID:          artifactID,
		WorkspaceID: wsUUID,
	})
	if err != nil {
		writeError(w, http.StatusNotFound, "artifact not found")
		return
	}
	prefix := h.getIssuePrefix(ctx, wsUUID)
	writeJSON(w, http.StatusOK, artifactToResponse(artifact, prefix))
}

func (h *Handler) CreateArtifact(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
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

	wsUUID := parseUUID(workspaceID)
	number, err := h.Queries.BumpArtifactCounter(ctx, wsUUID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to create artifact")
		return
	}

	actorType, actorID := h.resolveActor(r, userID, workspaceID)
	artifact, err := h.Queries.CreateArtifact(ctx, db.CreateArtifactParams{
		WorkspaceID:   wsUUID,
		ProjectID:     projectID,
		Title:         sanitizeNullBytes(req.Title),
		Summary:       sanitizeNullBytes(req.Summary),
		Content:       sanitizeNullBytes(req.Content),
		ContentType:   contentType,
		CreatorType:   actorType,
		CreatorID:     parseUUID(actorID),
		OriginIssueID: originIssueID,
		OriginTaskID:  originTaskID,
		Number:        number,
	})
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to create artifact")
		return
	}

	prefix := h.getIssuePrefix(ctx, wsUUID)
	resp := artifactToResponse(artifact, prefix)
	h.publish(protocol.EventArtifactCreated, workspaceID, actorType, actorID, map[string]any{"artifact": resp})
	writeJSON(w, http.StatusCreated, resp)
}

func (h *Handler) UpdateArtifact(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	workspaceID := h.resolveWorkspaceID(r)
	if workspaceID == "" {
		writeError(w, http.StatusBadRequest, "workspace_id is required")
		return
	}
	userID, ok := requireUserID(w, r)
	if !ok {
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

	wsUUID := parseUUID(workspaceID)
	artifact, err := h.Queries.UpdateArtifact(ctx, db.UpdateArtifactParams{
		ID:            artifactID,
		WorkspaceID:   wsUUID,
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

	prefix := h.getIssuePrefix(ctx, wsUUID)
	resp := artifactToResponse(artifact, prefix)
	actorType, actorID := h.resolveActor(r, userID, workspaceID)
	h.publish(protocol.EventArtifactUpdated, workspaceID, actorType, actorID, map[string]any{"artifact": resp})
	writeJSON(w, http.StatusOK, resp)
}

func (h *Handler) DeleteArtifact(w http.ResponseWriter, r *http.Request) {
	workspaceID := h.resolveWorkspaceID(r)
	if workspaceID == "" {
		writeError(w, http.StatusBadRequest, "workspace_id is required")
		return
	}
	userID, ok := requireUserID(w, r)
	if !ok {
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
	actorType, actorID := h.resolveActor(r, userID, workspaceID)
	h.publish(protocol.EventArtifactDeleted, workspaceID, actorType, actorID, map[string]any{"artifact_id": uuidToString(artifactID)})
	w.WriteHeader(http.StatusNoContent)
}
