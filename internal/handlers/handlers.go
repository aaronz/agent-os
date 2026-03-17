package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/agent-os/core/internal/services"
	"github.com/agent-os/core/pkg/metrics"
)

// Response represents a standard API response
type Response struct {
	Code      int         `json:"code"`
	Message   string      `json:"message"`
	Data      interface{} `json:"data,omitempty"`
	RequestID string      `json:"request_id"`
}

func writeJSON(w http.ResponseWriter, code int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(data)
}

func writeError(w http.ResponseWriter, code int, message string) {
	writeJSON(w, code, Response{
		Code:      code,
		Message:   message,
		RequestID: "",
	})
}

func extractPathParam(r *http.Request, key string) string {
	path := r.URL.Path
	prefix := "/api/v1/"
	if idx := len(prefix); len(path) > idx {
		path = path[idx:]
	}
	for _, segment := range splitPath(path) {
		if len(segment) > len(key)+1 && segment[:len(key)+1] == key+"_" {
			return segment[len(key)+1:]
		}
		if segment == key {
			return ""
		}
	}
	return ""
}

func splitPath(path string) []string {
	var result []string
	var current []byte
	for _, c := range path {
		if c == '/' {
			if len(current) > 0 {
				result = append(result, string(current))
				current = nil
			}
		} else {
			current = append(current, byte(c))
		}
	}
	if len(current) > 0 {
		result = append(result, string(current))
	}
	return result
}

// Health Handler
type HealthHandler struct{}

func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

func (h *HealthHandler) Health(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{
		"status": "healthy",
	})
}

type MetricsHandler struct{}

func NewMetricsHandler() *MetricsHandler {
	return &MetricsHandler{}
}

func (h *MetricsHandler) GetMetrics(w http.ResponseWriter, r *http.Request) {
	m := metrics.Global()
	snapshot := m.Snapshot()
	writeJSON(w, http.StatusOK, snapshot)
}

// Organization Handler
type OrganizationHandler struct {
	service *services.OrganizationService
}

func NewOrganizationHandler(service *services.OrganizationService) *OrganizationHandler {
	return &OrganizationHandler{service: service}
}

func (h *OrganizationHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		Owner       string `json:"owner"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	org, err := h.service.Create(r.Context(), req.Name, req.Description, req.Owner)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusCreated, Response{Code: 0, Message: "success", Data: org})
}

func (h *OrganizationHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get(":org_id")
	org, err := h.service.Get(r.Context(), id)
	if err != nil {
		writeError(w, http.StatusNotFound, "Organization not found")
		return
	}

	writeJSON(w, http.StatusOK, Response{Code: 0, Message: "success", Data: org})
}

func (h *OrganizationHandler) List(w http.ResponseWriter, r *http.Request) {
	orgs, err := h.service.List(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, Response{Code: 0, Message: "success", Data: orgs})
}

func (h *OrganizationHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := extractPathParam(r, "org_id")
	if id == "" {
		writeError(w, http.StatusBadRequest, "Missing organization ID")
		return
	}

	org, err := h.service.Get(r.Context(), id)
	if err != nil {
		writeError(w, http.StatusNotFound, "Organization not found")
		return
	}

	var req struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.Name != "" {
		org.Name = req.Name
	}
	if req.Description != "" {
		org.Description = req.Description
	}

	if err := h.service.Update(r.Context(), org); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, Response{Code: 0, Message: "success", Data: org})
}

func (h *OrganizationHandler) UpdateGovernance(w http.ResponseWriter, r *http.Request) {
	id := extractPathParam(r, "org_id")
	if id == "" {
		writeError(w, http.StatusBadRequest, "Missing organization ID")
		return
	}

	var req struct {
		Rules string `json:"governance_rules"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := h.service.UpdateGovernance(r.Context(), id, req.Rules); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, Response{Code: 0, Message: "success", Data: nil})
}

func (h *OrganizationHandler) UpdateStatus(w http.ResponseWriter, r *http.Request) {
	id := extractPathParam(r, "org_id")
	if id == "" {
		writeError(w, http.StatusBadRequest, "Missing organization ID")
		return
	}

	var req struct {
		Status string `json:"status"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.Status != "active" && req.Status != "disabled" {
		writeError(w, http.StatusBadRequest, "Invalid status. Must be 'active' or 'disabled'")
		return
	}

	if err := h.service.UpdateStatus(r.Context(), id, req.Status); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, Response{Code: 0, Message: "success", Data: nil})
}

func (h *OrganizationHandler) EmergencyIntervention(w http.ResponseWriter, r *http.Request) {
	var req struct {
		OrgID    string `json:"org_id"`
		IntentID string `json:"intent_id"`
		Action   string `json:"action"`
		Reason   string `json:"reason"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := h.service.EmergencyIntervention(r.Context(), req.OrgID, req.IntentID, req.Action, req.Reason); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, Response{Code: 0, Message: "Intervention executed", Data: nil})
}

// Agent Handler
type AgentHandler struct {
	service     *services.AgentService
	authService *services.AuthService
}

func NewAgentHandler(service *services.AgentService, authService *services.AuthService) *AgentHandler {
	return &AgentHandler{service: service, authService: authService}
}

func (h *AgentHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req struct {
		OrgID        string `json:"org_id"`
		Name         string `json:"name"`
		Role         string `json:"role"`
		Capabilities string `json:"capabilities"`
		CreatedBy    string `json:"created_by"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	agent, err := h.service.Create(r.Context(), req.OrgID, req.Name, req.Role, req.Capabilities, req.CreatedBy)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusCreated, Response{Code: 0, Message: "success", Data: agent})
}

func (h *AgentHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := extractPathParam(r, "agent_id")
	if id == "" {
		writeError(w, http.StatusBadRequest, "Missing agent ID")
		return
	}

	agent, err := h.service.Get(r.Context(), id)
	if err != nil {
		writeError(w, http.StatusNotFound, "Agent not found")
		return
	}

	writeJSON(w, http.StatusOK, Response{Code: 0, Message: "success", Data: agent})
}

func (h *AgentHandler) List(w http.ResponseWriter, r *http.Request) {
	orgID := r.URL.Query().Get("org_id")
	role := r.URL.Query().Get("role")
	status := r.URL.Query().Get("status")
	minRep := 0
	fmt.Sscanf(r.URL.Query().Get("min_reputation"), "%d", &minRep)

	agents, err := h.service.List(r.Context(), orgID, role, status, minRep)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, Response{Code: 0, Message: "success", Data: agents})
}

func (h *AgentHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := extractPathParam(r, "agent_id")
	if id == "" {
		writeError(w, http.StatusBadRequest, "Missing agent ID")
		return
	}

	agent, err := h.service.Get(r.Context(), id)
	if err != nil {
		writeError(w, http.StatusNotFound, "Agent not found")
		return
	}

	var req struct {
		Name         string `json:"name"`
		Role         string `json:"role"`
		Capabilities string `json:"capabilities"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.Name != "" {
		agent.Name = req.Name
	}
	if req.Role != "" {
		agent.Role = req.Role
	}
	if req.Capabilities != "" {
		agent.Capabilities = req.Capabilities
	}

	if err := h.service.Update(r.Context(), agent); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, Response{Code: 0, Message: "success", Data: agent})
}

func (h *AgentHandler) UpdateStatus(w http.ResponseWriter, r *http.Request) {
	id := extractPathParam(r, "agent_id")
	if id == "" {
		writeError(w, http.StatusBadRequest, "Missing agent ID")
		return
	}

	var req struct {
		Status string `json:"status"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	validStatuses := map[string]bool{"idle": true, "busy": true, "disabled": true, "terminated": true}
	if !validStatuses[req.Status] {
		writeError(w, http.StatusBadRequest, "Invalid status")
		return
	}

	if err := h.service.UpdateStatus(r.Context(), id, req.Status); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, Response{Code: 0, Message: "success", Data: nil})
}

func (h *AgentHandler) GetReputation(w http.ResponseWriter, r *http.Request) {
	id := extractPathParam(r, "agent_id")
	if id == "" {
		writeError(w, http.StatusBadRequest, "Missing agent ID")
		return
	}

	agent, err := h.service.GetReputation(r.Context(), id)
	if err != nil {
		writeError(w, http.StatusNotFound, "Agent not found")
		return
	}

	writeJSON(w, http.StatusOK, Response{Code: 0, Message: "success", Data: map[string]interface{}{
		"agent_id":   agent.ID,
		"reputation": agent.Reputation,
		"status":     agent.Status,
	}})
}

func (h *AgentHandler) GetActivities(w http.ResponseWriter, r *http.Request) {
	id := extractPathParam(r, "agent_id")
	if id == "" {
		writeError(w, http.StatusBadRequest, "Missing agent ID")
		return
	}

	writeJSON(w, http.StatusOK, Response{Code: 0, Message: "success", Data: []struct{}{}})
}

// Intent Handler
type IntentHandler struct {
	service *services.IntentService
}

func NewIntentHandler(service *services.IntentService) *IntentHandler {
	return &IntentHandler{service: service}
}

func (h *IntentHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req struct {
		OrgID           string `json:"org_id"`
		Title           string `json:"title"`
		Description     string `json:"description"`
		Constraints     string `json:"constraints"`
		SuccessCriteria string `json:"success_criteria"`
		Priority        string `json:"priority"`
		CreatedBy       string `json:"created_by"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.Title == "" || req.Description == "" || req.SuccessCriteria == "" {
		writeError(w, http.StatusBadRequest, "Title, description, and success_criteria are required")
		return
	}

	if req.Priority == "" {
		req.Priority = "medium"
	}

	intent, err := h.service.Create(r.Context(), req.OrgID, req.Title, req.Description, req.Constraints, req.SuccessCriteria, req.Priority, req.CreatedBy)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusCreated, Response{Code: 0, Message: "success", Data: intent})
}

func (h *IntentHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := extractPathParam(r, "intent_id")
	if id == "" {
		writeError(w, http.StatusBadRequest, "Missing intent ID")
		return
	}

	intent, err := h.service.Get(r.Context(), id)
	if err != nil {
		writeError(w, http.StatusNotFound, "Intent not found")
		return
	}

	writeJSON(w, http.StatusOK, Response{Code: 0, Message: "success", Data: intent})
}

func (h *IntentHandler) List(w http.ResponseWriter, r *http.Request) {
	orgID := r.URL.Query().Get("org_id")
	status := r.URL.Query().Get("status")
	priority := r.URL.Query().Get("priority")

	intents, err := h.service.List(r.Context(), orgID, status, priority)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, Response{Code: 0, Message: "success", Data: intents})
}

func (h *IntentHandler) UpdateStatus(w http.ResponseWriter, r *http.Request) {
	id := extractPathParam(r, "intent_id")
	if id == "" {
		writeError(w, http.StatusBadRequest, "Missing intent ID")
		return
	}

	var req struct {
		Status string `json:"status"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := h.service.UpdateStatus(r.Context(), id, req.Status); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, Response{Code: 0, Message: "success", Data: nil})
}

func (h *IntentHandler) GetTrace(w http.ResponseWriter, r *http.Request) {
	id := extractPathParam(r, "intent_id")
	if id == "" {
		writeError(w, http.StatusBadRequest, "Missing intent ID")
		return
	}

	intent, err := h.service.Get(r.Context(), id)
	if err != nil {
		writeError(w, http.StatusNotFound, "Intent not found")
		return
	}

	writeJSON(w, http.StatusOK, Response{Code: 0, Message: "success", Data: map[string]interface{}{
		"intent":   intent,
		"trace_id": intent.TraceID,
		"plan_ref": intent.PlanRef,
	}})
}

func (h *IntentHandler) GetPendingPlanning(w http.ResponseWriter, r *http.Request) {
	intents, err := h.service.GetPendingPlanning(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, Response{Code: 0, Message: "success", Data: intents})
}

func (h *IntentHandler) SubmitTaskGraph(w http.ResponseWriter, r *http.Request) {
	var req struct {
		IntentID     string `json:"intent_id"`
		Tasks        string `json:"tasks"`
		Dependencies string `json:"dependencies"`
		CreatedBy    string `json:"created_by"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	graph, err := h.service.SubmitTaskGraph(r.Context(), req.IntentID, req.Tasks, req.Dependencies, req.CreatedBy)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusCreated, Response{Code: 0, Message: "success", Data: graph})
}

func (h *IntentHandler) GetTaskGraph(w http.ResponseWriter, r *http.Request) {
	id := extractPathParam(r, "graph_id")
	if id == "" {
		writeError(w, http.StatusBadRequest, "Missing graph ID")
		return
	}

	graph, err := h.service.GetTaskGraph(r.Context(), id)
	if err != nil {
		writeError(w, http.StatusNotFound, "Task Graph not found")
		return
	}

	writeJSON(w, http.StatusOK, Response{Code: 0, Message: "success", Data: graph})
}

func (h *IntentHandler) UpdateTaskGraph(w http.ResponseWriter, r *http.Request) {
	writeError(w, http.StatusNotImplemented, "Task Graph versioning not implemented")
}

// Task Handler
type TaskHandler struct {
	service *services.TaskService
}

func NewTaskHandler(service *services.TaskService) *TaskHandler {
	return &TaskHandler{service: service}
}

func (h *TaskHandler) List(w http.ResponseWriter, r *http.Request) {
	orgID := r.URL.Query().Get("org_id")
	status := r.URL.Query().Get("status")
	capabilities := r.URL.Query().Get("capabilities")

	tasks, err := h.service.List(r.Context(), orgID, status, capabilities)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, Response{Code: 0, Message: "success", Data: tasks})
}

func (h *TaskHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := extractPathParam(r, "task_id")
	if id == "" {
		writeError(w, http.StatusBadRequest, "Missing task ID")
		return
	}

	task, err := h.service.Get(r.Context(), id)
	if err != nil {
		writeError(w, http.StatusNotFound, "Task not found")
		return
	}

	writeJSON(w, http.StatusOK, Response{Code: 0, Message: "success", Data: task})
}

func (h *TaskHandler) SubmitBid(w http.ResponseWriter, r *http.Request) {
	taskID := extractPathParam(r, "task_id")
	if taskID == "" {
		writeError(w, http.StatusBadRequest, "Missing task ID")
		return
	}

	var req struct {
		AgentID          string `json:"agent_id"`
		EstimatedTimeMin int    `json:"estimated_time_min"`
		Confidence       int    `json:"confidence"`
		Proposal         string `json:"proposal"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := h.service.SubmitBid(r.Context(), taskID, req.AgentID, req.EstimatedTimeMin, req.Confidence, req.Proposal); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, Response{Code: 0, Message: "success", Data: nil})
}

func (h *TaskHandler) AcceptTask(w http.ResponseWriter, r *http.Request) {
	taskID := extractPathParam(r, "task_id")
	if taskID == "" {
		writeError(w, http.StatusBadRequest, "Missing task ID")
		return
	}

	var req struct {
		AgentID string `json:"agent_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := h.service.AcceptTask(r.Context(), taskID, req.AgentID); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, Response{Code: 0, Message: "success", Data: nil})
}

func (h *TaskHandler) ReportProgress(w http.ResponseWriter, r *http.Request) {
	taskID := extractPathParam(r, "task_id")
	if taskID == "" {
		writeError(w, http.StatusBadRequest, "Missing task ID")
		return
	}

	var req struct {
		Progress string `json:"progress"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := h.service.ReportProgress(r.Context(), taskID, req.Progress); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, Response{Code: 0, Message: "success", Data: nil})
}

func (h *TaskHandler) ReportFailure(w http.ResponseWriter, r *http.Request) {
	taskID := extractPathParam(r, "task_id")
	if taskID == "" {
		writeError(w, http.StatusBadRequest, "Missing task ID")
		return
	}

	var req struct {
		Reason string `json:"reason"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := h.service.ReportFailure(r.Context(), taskID, req.Reason); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, Response{Code: 0, Message: "success", Data: nil})
}

func (h *TaskHandler) GetArtifacts(w http.ResponseWriter, r *http.Request) {
	writeError(w, http.StatusNotImplemented, "Not implemented")
}

func (h *TaskHandler) GetReviews(w http.ResponseWriter, r *http.Request) {
	writeError(w, http.StatusNotImplemented, "Not implemented")
}

// Artifact Handler
type ArtifactHandler struct {
	service *services.ArtifactService
}

func NewArtifactHandler(service *services.ArtifactService) *ArtifactHandler {
	return &ArtifactHandler{service: service}
}

func (h *ArtifactHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req struct {
		OrgID        string `json:"org_id"`
		TaskID       string `json:"task_id"`
		IntentID     string `json:"intent_id"`
		TraceID      string `json:"trace_id"`
		Type         string `json:"type"`
		Title        string `json:"title"`
		Description  string `json:"description"`
		ContentRef   string `json:"content_ref"`
		ContentHash  string `json:"content_hash"`
		Dependencies string `json:"dependencies"`
		CreatedBy    string `json:"created_by"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	artifact, err := h.service.Create(r.Context(), req.OrgID, req.TaskID, req.IntentID, req.TraceID, req.Type, req.Title, req.Description, req.ContentRef, req.ContentHash, req.Dependencies, req.CreatedBy)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusCreated, Response{Code: 0, Message: "success", Data: artifact})
}

func (h *ArtifactHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := extractPathParam(r, "artifact_id")
	if id == "" {
		writeError(w, http.StatusBadRequest, "Missing artifact ID")
		return
	}

	artifact, err := h.service.Get(r.Context(), id)
	if err != nil {
		writeError(w, http.StatusNotFound, "Artifact not found")
		return
	}

	writeJSON(w, http.StatusOK, Response{Code: 0, Message: "success", Data: artifact})
}

func (h *ArtifactHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := extractPathParam(r, "artifact_id")
	if id == "" {
		writeError(w, http.StatusBadRequest, "Missing artifact ID")
		return
	}

	var req struct {
		ContentRef  string `json:"content_ref"`
		ContentHash string `json:"content_hash"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := h.service.Update(r.Context(), id, req.ContentRef, req.ContentHash); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, Response{Code: 0, Message: "success", Data: nil})
}

func (h *ArtifactHandler) ListVersions(w http.ResponseWriter, r *http.Request) {
	writeError(w, http.StatusNotImplemented, "Not implemented")
}

func (h *ArtifactHandler) GetDependencies(w http.ResponseWriter, r *http.Request) {
	writeError(w, http.StatusNotImplemented, "Not implemented")
}

// Review Handler
type ReviewHandler struct {
	service *services.ReviewService
}

func NewReviewHandler(service *services.ReviewService) *ReviewHandler {
	return &ReviewHandler{service: service}
}

func (h *ReviewHandler) GetPendingReviews(w http.ResponseWriter, r *http.Request) {
	writeError(w, http.StatusNotImplemented, "Not implemented")
}

func (h *ReviewHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req struct {
		OrgID           string `json:"org_id"`
		ArtifactID      string `json:"artifact_id"`
		TaskID          string `json:"task_id"`
		IntentID        string `json:"intent_id"`
		TraceID         string `json:"trace_id"`
		ReviewerAgentID string `json:"reviewer_agent_id"`
		Score           int    `json:"score"`
		IsApproved      bool   `json:"is_approved"`
		Comments        string `json:"comments"`
		RejectionReason string `json:"rejection_reason"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	review, err := h.service.Create(r.Context(), req.OrgID, req.ArtifactID, req.TaskID, req.IntentID, req.TraceID, req.ReviewerAgentID, req.Score, req.IsApproved, req.Comments, req.RejectionReason)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusCreated, Response{Code: 0, Message: "success", Data: review})
}

func (h *ReviewHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := extractPathParam(r, "review_id")
	if id == "" {
		writeError(w, http.StatusBadRequest, "Missing review ID")
		return
	}

	review, err := h.service.Get(r.Context(), id)
	if err != nil {
		writeError(w, http.StatusNotFound, "Review not found")
		return
	}

	writeJSON(w, http.StatusOK, Response{Code: 0, Message: "success", Data: review})
}

// Memory Handler
type MemoryHandler struct {
	service *services.MemoryService
}

func NewMemoryHandler(service *services.MemoryService) *MemoryHandler {
	return &MemoryHandler{service: service}
}

func (h *MemoryHandler) Search(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Query string `json:"query"`
		TopK  int    `json:"top_k"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.TopK == 0 {
		req.TopK = 10
	}

	memories, err := h.service.Search(r.Context(), req.Query, req.TopK)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, Response{Code: 0, Message: "success", Data: memories})
}

func (h *MemoryHandler) List(w http.ResponseWriter, r *http.Request) {
	orgID := r.URL.Query().Get("org_id")
	memType := r.URL.Query().Get("type")

	memories, err := h.service.List(r.Context(), orgID, memType)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, Response{Code: 0, Message: "success", Data: memories})
}

func (h *MemoryHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := extractPathParam(r, "memory_id")
	if id == "" {
		writeError(w, http.StatusBadRequest, "Missing memory ID")
		return
	}

	memory, err := h.service.Get(r.Context(), id)
	if err != nil {
		writeError(w, http.StatusNotFound, "Memory not found")
		return
	}

	writeJSON(w, http.StatusOK, Response{Code: 0, Message: "success", Data: memory})
}

func (h *MemoryHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := extractPathParam(r, "memory_id")
	if id == "" {
		writeError(w, http.StatusBadRequest, "Missing memory ID")
		return
	}

	var req struct {
		Validity string `json:"validity"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := h.service.UpdateValidity(r.Context(), id, req.Validity); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, Response{Code: 0, Message: "success", Data: nil})
}

// Audit Handler
func GetAuditTrail(w http.ResponseWriter, r *http.Request) {
	writeError(w, http.StatusNotImplemented, "Not implemented")
}

var _ = fmt.Sprintf("")
var _ = context.Background()
