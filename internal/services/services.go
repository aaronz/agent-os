package services

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/rand"
	"time"

	"github.com/agent-os/core/internal/models"
	"github.com/agent-os/core/internal/repository"
	"github.com/google/uuid"
)

// Organization Service
type OrganizationService struct {
	repo       *repository.OrganizationRepository
	eventQueue *repository.EventQueue
}

func NewOrganizationService(repo *repository.OrganizationRepository, eq *repository.EventQueue) *OrganizationService {
	return &OrganizationService{repo: repo, eventQueue: eq}
}

func (s *OrganizationService) Create(ctx context.Context, name, description, owner string) (*models.Organization, error) {
	org := &models.Organization{
		ID:          uuid.New().String(),
		Name:        name,
		Description: description,
		Owner:       owner,
		Status:      "active",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := s.repo.Create(org); err != nil {
		return nil, err
	}

	// Publish event
	s.publishEvent(ctx, "organization.created", org.ID, org)

	return org, nil
}

func (s *OrganizationService) Get(ctx context.Context, id string) (*models.Organization, error) {
	return s.repo.Get(id)
}

func (s *OrganizationService) List(ctx context.Context) ([]*models.Organization, error) {
	return s.repo.List()
}

func (s *OrganizationService) Update(ctx context.Context, org *models.Organization) error {
	org.UpdatedAt = time.Now()
	return s.repo.Update(org)
}

func (s *OrganizationService) UpdateStatus(ctx context.Context, id, status string) error {
	org, err := s.repo.Get(id)
	if err != nil {
		return err
	}
	org.Status = status
	org.UpdatedAt = time.Now()
	return s.repo.Update(org)
}

func (s *OrganizationService) UpdateGovernance(ctx context.Context, id string, rules string) error {
	org, err := s.repo.Get(id)
	if err != nil {
		return err
	}
	org.GovernanceRules = rules
	org.UpdatedAt = time.Now()
	return s.repo.Update(org)
}

func (s *OrganizationService) EmergencyIntervention(ctx context.Context, orgID, intentID, action, reason string) error {
	// Log the intervention
	fmt.Printf("Emergency intervention: org=%s, intent=%s, action=%s, reason=%s\n", orgID, intentID, action, reason)
	return nil
}

func (s *OrganizationService) publishEvent(ctx context.Context, eventType, entityID string, data interface{}) {
	if s.eventQueue == nil {
		return
	}
	payload, _ := json.Marshal(data)
	s.eventQueue.Publish(ctx, entityID, payload)
}

// Agent Service
type AgentService struct {
	repo       *repository.AgentRepository
	eventQueue *repository.EventQueue
}

func NewAgentService(repo *repository.AgentRepository, eq *repository.EventQueue) *AgentService {
	return &AgentService{repo: repo, eventQueue: eq}
}

func (s *AgentService) Create(ctx context.Context, orgID, name, role, capabilities, createdBy string) (*models.Agent, error) {
	apiKey := generateAPIKey()
	apiKeyHash := hashAPIKey(apiKey)

	agent := &models.Agent{
		ID:                uuid.New().String(),
		OrgID:             orgID,
		Name:              name,
		Role:              role,
		Capabilities:      capabilities,
		Reputation:        50,
		MemoryRefs:        "[]",
		Status:            "idle",
		APIKey:            apiKey,
		APIKeyHash:        apiKeyHash,
		CreatedBy:         createdBy,
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
		LastActiveAt:      time.Now(),
		BusyTaskCount:     0,
		AbandonedCount:    0,
		OverturnedReviews: 0,
	}

	if err := s.repo.Create(agent); err != nil {
		return nil, err
	}

	// Publish event
	s.publishEvent(ctx, "agent.created", agent.ID, agent)

	return agent, nil
}

func (s *AgentService) Get(ctx context.Context, id string) (*models.Agent, error) {
	return s.repo.Get(id)
}

func (s *AgentService) List(ctx context.Context, orgID, role, status string, minReputation int) ([]*models.Agent, error) {
	return s.repo.List(orgID, role, status, minReputation)
}

func (s *AgentService) Update(ctx context.Context, agent *models.Agent) error {
	agent.UpdatedAt = time.Now()
	return s.repo.Update(agent)
}

func (s *AgentService) UpdateStatus(ctx context.Context, id, status string) error {
	agent, err := s.repo.Get(id)
	if err != nil {
		return err
	}
	agent.Status = status
	agent.UpdatedAt = time.Now()
	return s.repo.Update(agent)
}

func (s *AgentService) GetReputation(ctx context.Context, id string) (*models.Agent, error) {
	return s.repo.Get(id)
}

func (s *AgentService) ValidateAPIKey(ctx context.Context, apiKey string) (*models.Agent, error) {
	hash := hashAPIKey(apiKey)
	return s.repo.GetByAPIKeyHash(hash)
}

func (s *AgentService) UpdateBusyCount(ctx context.Context, agentID string, delta int) error {
	agent, err := s.repo.Get(agentID)
	if err != nil {
		return err
	}
	agent.BusyTaskCount += delta
	if agent.BusyTaskCount < 0 {
		agent.BusyTaskCount = 0
	}
	if agent.BusyTaskCount > 0 {
		agent.Status = "busy"
	} else {
		agent.Status = "idle"
	}
	agent.UpdatedAt = time.Now()
	return s.repo.Update(agent)
}

func (s *AgentService) CanAcceptTask(ctx context.Context, agentID string) (bool, error) {
	agent, err := s.repo.Get(agentID)
	if err != nil {
		return false, err
	}
	if agent.Status == "disabled" || agent.Status == "terminated" {
		return false, nil
	}
	if agent.BusyTaskCount >= 5 {
		return false, nil
	}
	return true, nil
}

func (s *AgentService) UpdateReputation(ctx context.Context, agentID string, points int) error {
	agent, err := s.repo.Get(agentID)
	if err != nil {
		return err
	}
	agent.Reputation += points
	if agent.Reputation > 100 {
		agent.Reputation = 100
	}
	if agent.Reputation < 0 {
		agent.Reputation = 0
	}
	if agent.Reputation < 30 && agent.Status == "active" {
		agent.Status = "disabled"
	}
	agent.UpdatedAt = time.Now()
	return s.repo.Update(agent)
}

func (s *AgentService) publishEvent(ctx context.Context, eventType, entityID string, data interface{}) {
	if s.eventQueue == nil {
		return
	}
	payload, _ := json.Marshal(data)
	s.eventQueue.Publish(ctx, entityID, payload)
}

func generateAPIKey() string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, 32)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return "sk-" + string(b)
}

func generateEmbedding(text string) []float32 {
	hash := sha256.Sum256([]byte(text))
	embedding := make([]float32, 384)
	for i := range embedding {
		embedding[i] = float32(hash[i%len(hash)]) / 255.0
	}
	return embedding
}

func hashAPIKey(apiKey string) string {
	hash := sha256.Sum256([]byte(apiKey))
	return hex.EncodeToString(hash[:])
}

// Auth Service
type AuthService struct {
	agentRepo *repository.AgentRepository
}

func NewAuthService(agentRepo *repository.AgentRepository) *AuthService {
	return &AuthService{agentRepo: agentRepo}
}

func (s *AuthService) Authenticate(ctx context.Context, apiKey string) (*models.Agent, error) {
	return s.agentRepo.GetByAPIKeyHash(hashAPIKey(apiKey))
}

// Placeholder services - will be expanded
type IntentService struct {
	repo       *repository.IntentRepository
	graphStore *repository.GraphStore
	eventQueue *repository.EventQueue
}

func NewIntentService(repo *repository.IntentRepository, gs *repository.GraphStore, eq *repository.EventQueue) *IntentService {
	return &IntentService{repo: repo, graphStore: gs, eventQueue: eq}
}

func (s *IntentService) Create(ctx context.Context, orgID, title, description, constraints, successCriteria, priority, createdBy string) (*models.Intent, error) {
	intent := &models.Intent{
		ID:              uuid.New().String(),
		OrgID:           orgID,
		TraceID:         uuid.New().String(),
		Title:           title,
		Description:     description,
		Constraints:     constraints,
		SuccessCriteria: successCriteria,
		Priority:        priority,
		CreatedBy:       createdBy,
		Status:          "open",
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	if err := s.repo.Create(intent); err != nil {
		return nil, err
	}

	s.publishEvent(ctx, "intent.created", intent.ID, intent)
	return intent, nil
}

func (s *IntentService) Get(ctx context.Context, id string) (*models.Intent, error) {
	return s.repo.Get(id)
}

func (s *IntentService) List(ctx context.Context, orgID, status, priority string) ([]*models.Intent, error) {
	return s.repo.List(orgID, status, priority)
}

func (s *IntentService) UpdateStatus(ctx context.Context, id, status string) error {
	intent, err := s.repo.Get(id)
	if err != nil {
		return err
	}
	validTransitions := map[string][]string{
		"draft":     {"open", "cancelled"},
		"open":      {"planning", "cancelled"},
		"planning":  {"executing", "failed", "paused"},
		"executing": {"completed", "failed", "paused"},
		"paused":    {"executing", "cancelled"},
	}

	allowedStatuses, exists := validTransitions[intent.Status]
	if !exists {
		return fmt.Errorf("invalid current status: %s", intent.Status)
	}

	isAllowed := false
	for _, v := range allowedStatuses {
		if v == status {
			isAllowed = true
			break
		}
	}
	if !isAllowed {
		return fmt.Errorf("invalid status transition from %s to %s", intent.Status, status)
	}

	intent.Status = status
	intent.UpdatedAt = time.Now()
	if status == "completed" {
		now := time.Now()
		intent.ActualCompletedAt = &now
	}
	return s.repo.Update(intent)
}

func (s *IntentService) SubmitTaskGraph(ctx context.Context, intentID, tasks, dependencies, createdBy string) (*models.TaskGraph, error) {
	intent, err := s.repo.Get(intentID)
	if err != nil {
		return nil, err
	}

	graph := &models.TaskGraph{
		ID:           uuid.New().String(),
		IntentID:     intentID,
		OrgID:        intent.OrgID,
		TraceID:      intent.TraceID,
		Tasks:        tasks,
		Dependencies: dependencies,
		CreatedBy:    createdBy,
		Version:      1,
		Status:       "published",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	intent.PlanRef = graph.ID
	intent.Status = "executing"
	intent.UpdatedAt = time.Now()
	s.repo.Update(intent)

	s.publishEvent(ctx, "taskgraph.created", graph.ID, graph)
	return graph, nil
}

func (s *IntentService) GetTaskGraph(ctx context.Context, graphID string) (*models.TaskGraph, error) {
	graph, err := s.graphStore.GetTaskGraph(graphID)
	if err != nil {
		return nil, err
	}

	pg := &models.TaskGraph{
		ID:           graph["id"].(string),
		IntentID:     graph["intent_id"].(string),
		OrgID:        graph["org_id"].(string),
		TraceID:      graph["trace_id"].(string),
		Tasks:        graph["tasks"].(string),
		Dependencies: graph["dependencies"].(string),
		CreatedBy:    graph["created_by"].(string),
		Version:      int(graph["version"].(float64)),
		Status:       graph["status"].(string),
	}
	return pg, nil
}

func (s *IntentService) UpdateTaskGraph(ctx context.Context, graphID, tasks, dependencies, changeReason string) (*models.TaskGraph, error) {
	existing, err := s.GetTaskGraph(ctx, graphID)
	if err != nil {
		return nil, err
	}

	if existing.Status == "completed" {
		return nil, fmt.Errorf("cannot update completed task graph")
	}

	newVersion := existing.Version + 1
	updatedGraph := &models.TaskGraph{
		ID:              existing.ID,
		IntentID:        existing.IntentID,
		OrgID:           existing.OrgID,
		TraceID:         existing.TraceID,
		Tasks:           tasks,
		Dependencies:    dependencies,
		CreatedBy:       existing.CreatedBy,
		Version:         newVersion,
		Status:          "published",
		ChangeReason:    changeReason,
		PreviousVersion: existing.Version,
		CreatedAt:       existing.CreatedAt,
		UpdatedAt:       time.Now(),
	}

	graphData := map[string]interface{}{
		"id":               updatedGraph.ID,
		"intent_id":        updatedGraph.IntentID,
		"org_id":           updatedGraph.OrgID,
		"trace_id":         updatedGraph.TraceID,
		"tasks":            updatedGraph.Tasks,
		"dependencies":     updatedGraph.Dependencies,
		"created_by":       updatedGraph.CreatedBy,
		"version":          updatedGraph.Version,
		"status":           updatedGraph.Status,
		"change_reason":    changeReason,
		"previous_version": existing.Version,
		"created_at":       existing.CreatedAt.Unix(),
		"updated_at":       time.Now().Unix(),
	}

	s.graphStore.SaveTaskGraph(updatedGraph.ID, graphData)
	s.publishEvent(ctx, "taskgraph.updated", graphID, map[string]interface{}{
		"old_version": existing.Version,
		"new_version": newVersion,
		"reason":      changeReason,
	})

	return updatedGraph, nil
}

func (s *IntentService) GetPendingPlanning(ctx context.Context) ([]*models.Intent, error) {
	return s.repo.List("", "open", "")
}

func (s *IntentService) ValidateTaskGraph(tasksJSON, dependenciesJSON, intentConstraints, intentSuccessCriteria string) error {
	var tasks []map[string]interface{}
	var dependencies []map[string]any
	var intentConstraintsArr []string
	var intentSuccessCriteriaArr []string

	if err := json.Unmarshal([]byte(tasksJSON), &tasks); err != nil {
		return fmt.Errorf("invalid tasks JSON: %w", err)
	}
	if err := json.Unmarshal([]byte(dependenciesJSON), &dependencies); err != nil {
		return fmt.Errorf("invalid dependencies JSON: %w", err)
	}
	json.Unmarshal([]byte(intentConstraints), &intentConstraintsArr)
	json.Unmarshal([]byte(intentSuccessCriteria), &intentSuccessCriteriaArr)

	taskIDs := make(map[string]bool)
	taskCapabilities := make(map[string][]string)
	taskAcceptanceCriteria := make(map[string][]string)

	for _, task := range tasks {
		if id, ok := task["id"].(string); ok {
			taskIDs[id] = true

			if caps, ok := task["required_capabilities"].([]interface{}); ok {
				for _, c := range caps {
					if cap, ok := c.(string); ok {
						taskCapabilities[id] = append(taskCapabilities[id], cap)
					}
				}
			}
			if ac, ok := task["acceptance_criteria"].([]interface{}); ok {
				for _, a := range ac {
					if crit, ok := a.(string); ok {
						taskAcceptanceCriteria[id] = append(taskAcceptanceCriteria[id], crit)
					}
				}
			}
		}
	}

	for _, dep := range dependencies {
		fromID, _ := dep["from_task_id"].(string)
		toID, _ := dep["to_task_id"].(string)
		if fromID != "" && !taskIDs[fromID] {
			return fmt.Errorf("dangling dependency: from_task_id %s not found", fromID)
		}
		if toID != "" && !taskIDs[toID] {
			return fmt.Errorf("dangling dependency: to_task_id %s not found", toID)
		}
		if fromID != "" && toID == fromID {
			return fmt.Errorf("self-loop dependency detected: task %s depends on itself", fromID)
		}
	}

	if err := detectCycle(taskIDs, dependencies); err != nil {
		return fmt.Errorf("cycle detected in task graph: %w", err)
	}

	if err := validateCompleteness(taskAcceptanceCriteria, intentSuccessCriteriaArr); err != nil {
		return fmt.Errorf("completeness validation failed: %w", err)
	}

	return nil
}

func validateCompleteness(taskCriteria map[string][]string, intentCriteria []string) error {
	coveredCriteria := make(map[string]bool)

	for _, criteria := range taskCriteria {
		for _, c := range criteria {
			coveredCriteria[c] = true
		}
	}

	for _, ic := range intentCriteria {
		found := false
		for _, covered := range coveredCriteria {
			if covered && containsSimilar(ic, coveredCriteria) {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("intent success criteria '%s' not covered by any task", ic)
		}
	}
	return nil
}

func containsSimilar(text string, criteriaMap map[string]bool) bool {
	for k := range criteriaMap {
		if len(k) > 3 && len(text) > 3 {
			if contains(k, text) || contains(text, k) {
				return true
			}
		}
	}
	return false
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsAt(s, substr, 0))
}

func containsAt(s, substr string, start int) bool {
	for i := start; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

func detectCycle(taskIDs map[string]bool, dependencies []map[string]any) error {
	graph := make(map[string][]string)
	for _, dep := range dependencies {
		fromID, _ := dep["from_task_id"].(string)
		toID, _ := dep["to_task_id"].(string)
		if fromID != "" && toID != "" {
			graph[fromID] = append(graph[fromID], toID)
		}
	}

	visited := make(map[string]bool)
	recStack := make(map[string]bool)

	var dfs func(node string) error
	dfs = func(node string) error {
		visited[node] = true
		recStack[node] = true

		for _, neighbor := range graph[node] {
			if !visited[neighbor] {
				if err := dfs(neighbor); err != nil {
					return err
				}
			} else if recStack[neighbor] {
				return fmt.Errorf("cycle found: %s -> %s", node, neighbor)
			}
		}

		recStack[node] = false
		return nil
	}

	for taskID := range taskIDs {
		if !visited[taskID] {
			if err := dfs(taskID); err != nil {
				return err
			}
		}
	}

	return nil
}

func (s *IntentService) PublishDependentTasks(ctx context.Context, completedTaskID string) error {
	return nil
}

func (s *IntentService) publishEvent(ctx context.Context, eventType, entityID string, data interface{}) {
	if s.eventQueue == nil {
		return
	}
	payload, _ := json.Marshal(data)
	s.eventQueue.Publish(ctx, entityID, payload)
}

type TaskService struct {
	repo         *repository.TaskRepository
	agentRepo    *repository.AgentRepository
	bidRepo      *repository.BidRepository
	artifactRepo *repository.ArtifactRepository
	graphStore   *repository.GraphStore
	eventQueue   *repository.EventQueue
}

func NewTaskService(repo *repository.TaskRepository, agentRepo *repository.AgentRepository, bidRepo *repository.BidRepository, artifactRepo *repository.ArtifactRepository, gs *repository.GraphStore, eq *repository.EventQueue) *TaskService {
	return &TaskService{repo: repo, agentRepo: agentRepo, bidRepo: bidRepo, artifactRepo: artifactRepo, graphStore: gs, eventQueue: eq}
}

func (s *TaskService) List(ctx context.Context, orgID, status, capabilities string) ([]*models.Task, error) {
	return s.repo.List(orgID, status, capabilities)
}

func (s *TaskService) Get(ctx context.Context, id string) (*models.Task, error) {
	return s.repo.Get(id)
}

func (s *TaskService) OpenBidding(ctx context.Context, taskID string, biddingWindowMin int) error {
	task, err := s.repo.Get(taskID)
	if err != nil {
		return err
	}

	if task.Status != "open" && task.Status != "pending" {
		return fmt.Errorf("task is not in a state that allows bidding (current: %s)", task.Status)
	}

	windowDuration := time.Duration(biddingWindowMin) * time.Minute
	if windowDuration == 0 {
		windowDuration = 2 * time.Hour
	}
	biddingEnd := time.Now().Add(windowDuration)
	task.BiddingEndTime = &biddingEnd
	task.Status = "open"
	task.UpdatedAt = time.Now()

	if err := s.repo.Update(task); err != nil {
		return err
	}

	s.publishEvent(ctx, "task.bidding_opened", taskID, map[string]interface{}{
		"bidding_end_time": biddingEnd,
	})
	return nil
}

func (s *TaskService) SubmitBid(ctx context.Context, taskID, agentID string, estimatedTimeMin, confidence int, proposal string) error {
	task, err := s.repo.Get(taskID)
	if err != nil {
		return err
	}

	if task.Status != "open" {
		return fmt.Errorf("task is not open for bidding")
	}

	if task.BiddingEndTime != nil && time.Now().After(*task.BiddingEndTime) {
		return fmt.Errorf("bidding window closed")
	}

	bid := &models.Bid{
		ID:               uuid.New().String(),
		OrgID:            task.OrgID,
		TaskID:           taskID,
		AgentID:          agentID,
		EstimatedTimeMin: estimatedTimeMin,
		Confidence:       confidence,
		Proposal:         proposal,
		CreatedAt:        time.Now(),
		Status:           "pending",
	}

	if err := s.bidRepo.Create(bid); err != nil {
		return err
	}

	task.Status = "bidding"
	task.UpdatedAt = time.Now()
	s.repo.Update(task)

	s.publishEvent(ctx, "bid.created", bid.ID, bid)
	return nil
}

type BidScore struct {
	BidID   string
	AgentID string
	Score   float64
}

func (s *TaskService) AllocateTask(ctx context.Context, taskID string) (string, error) {
	task, err := s.repo.Get(taskID)
	if err != nil {
		return "", err
	}

	if task.Status != "bidding" {
		return "", fmt.Errorf("task not in bidding state")
	}

	bids, err := s.bidRepo.GetBidsByTask(taskID)
	if err != nil || len(bids) == 0 {
		return "", fmt.Errorf("no bids available")
	}

	var scores []BidScore
	for _, bid := range bids {
		agent, err := s.agentRepo.Get(bid.AgentID)
		if err != nil {
			continue
		}

		estimatedTime := float64(bid.EstimatedTimeMin)
		if estimatedTime == 0 {
			estimatedTime = 1
		}

		reputationScore := float64(agent.Reputation) * 0.5
		timeScore := (1.0 / estimatedTime) * 100 * 0.3
		confidenceScore := float64(bid.Confidence) * 0.2

		totalScore := reputationScore + timeScore + confidenceScore

		scores = append(scores, BidScore{
			BidID:   bid.ID,
			AgentID: bid.AgentID,
			Score:   totalScore,
		})
	}

	if len(scores) == 0 {
		return "", fmt.Errorf("no valid bids")
	}

	var bestBid BidScore
	for _, score := range scores {
		if score.Score > bestBid.Score {
			bestBid = score
		}
	}

	task.AssignedAgentID = bestBid.AgentID
	task.Status = "assigned"
	task.UpdatedAt = time.Now()
	s.repo.Update(task)

	return bestBid.AgentID, nil
}

func (s *TaskService) CloseBidding(ctx context.Context, taskID string) error {
	task, err := s.repo.Get(taskID)
	if err != nil {
		return err
	}

	if task.Status == "bidding" {
		task.Status = "open"
		task.UpdatedAt = time.Now()
		s.repo.Update(task)
	}
	return nil
}

func (s *TaskService) AcceptTask(ctx context.Context, taskID, agentID string) error {
	task, err := s.repo.Get(taskID)
	if err != nil {
		return err
	}

	if task.Status != "assigned" {
		return fmt.Errorf("task is not assigned to this agent")
	}

	task.Status = "executing"
	task.UpdatedAt = time.Now()

	deadline := time.Now().Add(time.Duration(task.MaxExecutionTimeMin) * time.Minute)
	task.DeadlineAt = &deadline

	return s.repo.Update(task)
}

func (s *TaskService) ReportProgress(ctx context.Context, taskID, progress string) error {
	_, err := s.repo.Get(taskID)
	if err != nil {
		return err
	}

	s.publishEvent(ctx, "task.progress", taskID, map[string]string{
		"progress": progress,
	})
	return nil
}

func (s *TaskService) ReportFailure(ctx context.Context, taskID, reason string) error {
	task, err := s.repo.Get(taskID)
	if err != nil {
		return err
	}

	task.Status = "failed"
	task.AttemptCount++
	task.UpdatedAt = time.Now()

	if err := s.repo.Update(task); err != nil {
		return err
	}

	s.publishEvent(ctx, "task.failed", taskID, map[string]string{
		"reason": reason,
	})

	if task.AttemptCount < 3 {
		s.RepublishTask(ctx, taskID)
	}

	return nil
}

func (s *TaskService) RepublishTask(ctx context.Context, taskID string) error {
	task, err := s.repo.Get(taskID)
	if err != nil {
		return err
	}

	task.Status = "open"
	task.AssignedAgentID = ""
	task.BiddingEndTime = nil
	task.UpdatedAt = time.Now()

	if err := s.repo.Update(task); err != nil {
		return err
	}

	s.publishEvent(ctx, "task.republished", taskID, map[string]interface{}{
		"attempt_count": task.AttemptCount,
	})

	return nil
}

func (s *TaskService) CompleteTask(ctx context.Context, taskID string) error {
	task, err := s.repo.Get(taskID)
	if err != nil {
		return err
	}

	task.Status = "completed"
	task.UpdatedAt = time.Now()
	if err := s.repo.Update(task); err != nil {
		return err
	}

	s.PublishDependentTasks(ctx, task.GraphID, taskID)

	s.publishEvent(ctx, "task.completed", taskID, nil)
	return nil
}

func (s *TaskService) PublishDependentTasks(ctx context.Context, graphID, completedTaskID string) {
	allTasks, err := s.repo.List(graphID, "", "")
	if err != nil {
		return
	}

	for _, task := range allTasks {
		if task.Status != "pending" {
			continue
		}

		var deps []string
		json.Unmarshal([]byte(task.Dependencies), &deps)

		allDepsCompleted := true
		for _, depID := range deps {
			depTask, err := s.repo.Get(depID)
			if err != nil || depTask.Status != "completed" {
				allDepsCompleted = false
				break
			}
		}

		if allDepsCompleted && len(deps) > 0 {
			task.Status = "open"
			task.UpdatedAt = time.Now()
			s.repo.Update(task)
			s.publishEvent(ctx, "task.published", task.ID, map[string]interface{}{
				"reason":         "dependencies_completed",
				"completed_deps": deps,
			})
		}
	}
}

func (s *TaskService) GetDependentArtifacts(ctx context.Context, taskID string) ([]*models.Artifact, error) {
	task, err := s.repo.Get(taskID)
	if err != nil {
		return nil, err
	}

	var depTaskIDs []string
	json.Unmarshal([]byte(task.Dependencies), &depTaskIDs)

	var artifacts []*models.Artifact
	for _, depTaskID := range depTaskIDs {
		taskArtifacts, err := s.artifactRepo.ListByTask(depTaskID)
		if err != nil {
			continue
		}
		artifacts = append(artifacts, taskArtifacts...)
	}

	return artifacts, nil
}

func (s *TaskService) publishEvent(ctx context.Context, eventType, entityID string, data interface{}) {
	if s.eventQueue == nil {
		return
	}
	payload, _ := json.Marshal(data)
	s.eventQueue.Publish(ctx, entityID, payload)
}

type ArtifactService struct {
	repo       *repository.ArtifactRepository
	eventQueue *repository.EventQueue
}

func NewArtifactService(repo *repository.ArtifactRepository, eq *repository.EventQueue) *ArtifactService {
	return &ArtifactService{repo: repo, eventQueue: eq}
}

func (s *ArtifactService) Create(ctx context.Context, orgID, taskID, intentID, traceID, artifactType, title, description, contentRef, contentHash, dependencies, createdBy string) (*models.Artifact, error) {
	artifact := &models.Artifact{
		ID:           uuid.New().String(),
		OrgID:        orgID,
		TaskID:       taskID,
		IntentID:     intentID,
		TraceID:      traceID,
		Type:         artifactType,
		Title:        title,
		Description:  description,
		ContentRef:   contentRef,
		ContentHash:  contentHash,
		Dependencies: dependencies,
		CreatedBy:    createdBy,
		Version:      1,
		Status:       "pending_review",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := s.repo.Create(artifact); err != nil {
		return nil, err
	}

	s.publishEvent(ctx, "artifact.created", artifact.ID, artifact)
	return artifact, nil
}

func (s *ArtifactService) Get(ctx context.Context, id string) (*models.Artifact, error) {
	return s.repo.Get(id)
}

func (s *ArtifactService) ListByTask(ctx context.Context, taskID string) ([]*models.Artifact, error) {
	return s.repo.ListByTask(taskID)
}

func (s *ArtifactService) Update(ctx context.Context, id, contentRef, contentHash string) error {
	artifact, err := s.repo.Get(id)
	if err != nil {
		return err
	}

	artifact.ContentRef = contentRef
	artifact.ContentHash = contentHash
	artifact.Version++
	artifact.Status = "pending_review"
	artifact.UpdatedAt = time.Now()

	return s.repo.Update(artifact)
}

func (s *ArtifactService) UpdateStatus(ctx context.Context, id, status string) error {
	artifact, err := s.repo.Get(id)
	if err != nil {
		return err
	}

	artifact.Status = status
	artifact.UpdatedAt = time.Now()
	return s.repo.Update(artifact)
}

func (s *ArtifactService) GetDependencies(ctx context.Context, artifactID string) ([]*models.Artifact, error) {
	artifact, err := s.repo.Get(artifactID)
	if err != nil {
		return nil, err
	}

	var depIDs []string
	json.Unmarshal([]byte(artifact.Dependencies), &depIDs)

	var dependencies []*models.Artifact
	for _, depID := range depIDs {
		depArtifact, err := s.repo.Get(depID)
		if err == nil {
			dependencies = append(dependencies, depArtifact)
		}
	}

	return dependencies, nil
}

func (s *ArtifactService) GetDependents(ctx context.Context, artifactID string) ([]*models.Artifact, error) {
	allArtifacts, err := s.repo.ListByTask("")
	if err != nil {
		return nil, err
	}

	var dependents []*models.Artifact
	for _, artifact := range allArtifacts {
		var depIDs []string
		json.Unmarshal([]byte(artifact.Dependencies), &depIDs)
		for _, depID := range depIDs {
			if depID == artifactID {
				dependents = append(dependents, artifact)
				break
			}
		}
	}

	return dependents, nil
}

func (s *ArtifactService) BuildDependencyGraph(ctx context.Context, rootArtifactID string) (map[string][]string, error) {
	graph := make(map[string][]string)
	visited := make(map[string]bool)

	var buildGraph func(artifactID string)
	buildGraph = func(artifactID string) {
		if visited[artifactID] {
			return
		}
		visited[artifactID] = true

		artifact, err := s.repo.Get(artifactID)
		if err != nil {
			return
		}

		var depIDs []string
		json.Unmarshal([]byte(artifact.Dependencies), &depIDs)
		graph[artifactID] = depIDs

		for _, depID := range depIDs {
			buildGraph(depID)
		}
	}

	buildGraph(rootArtifactID)
	return graph, nil
}

func (s *ArtifactService) GetDependentsChain(ctx context.Context, artifactID string) ([]string, error) {
	chain := []string{}
	visited := make(map[string]bool)

	var traverse func(id string)
	traverse = func(id string) {
		if visited[id] {
			return
		}
		visited[id] = true

		dependents, err := s.GetDependents(ctx, id)
		if err != nil {
			return
		}

		for _, dep := range dependents {
			chain = append(chain, dep.ID)
			traverse(dep.ID)
		}
	}

	traverse(artifactID)
	return chain, nil
}

func (s *ArtifactService) publishEvent(ctx context.Context, eventType, entityID string, data interface{}) {
	if s.eventQueue == nil {
		return
	}
	payload, _ := json.Marshal(data)
	s.eventQueue.Publish(ctx, entityID, payload)
}

type ReviewService struct {
	repo         *repository.ReviewRepository
	agentRepo    *repository.AgentRepository
	artifactRepo *repository.ArtifactRepository
	taskRepo     *repository.TaskRepository
	eventQueue   *repository.EventQueue
}

func NewReviewService(repo *repository.ReviewRepository, agentRepo *repository.AgentRepository, ar *repository.ArtifactRepository, tr *repository.TaskRepository, eq *repository.EventQueue) *ReviewService {
	return &ReviewService{repo: repo, agentRepo: agentRepo, artifactRepo: ar, taskRepo: tr, eventQueue: eq}
}

func (s *ReviewService) Create(ctx context.Context, orgID, artifactID, taskID, intentID, traceID, reviewerAgentID string, score int, isApproved bool, comments, rejectionReason string) (*models.Review, error) {
	review := &models.Review{
		ID:              uuid.New().String(),
		OrgID:           orgID,
		ArtifactID:      artifactID,
		TaskID:          taskID,
		IntentID:        intentID,
		TraceID:         traceID,
		ReviewerAgentID: reviewerAgentID,
		Score:           score,
		IsApproved:      isApproved,
		Comments:        comments,
		RejectionReason: rejectionReason,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	if err := s.repo.Create(review); err != nil {
		return nil, err
	}

	if isApproved {
		s.artifactRepo.UpdateStatus(artifactID, "approved")
		s.taskRepo.Get(taskID)
		task, _ := s.taskRepo.Get(taskID)
		if task != nil {
			task.Status = "completed"
			task.UpdatedAt = time.Now()
			s.taskRepo.Update(task)
		}
	} else {
		s.artifactRepo.UpdateStatus(artifactID, "rejected")
		task, _ := s.taskRepo.Get(taskID)
		if task != nil {
			task.RejectionCount++
			if task.RejectionCount >= 3 {
				task.Status = "failed"
			}
			task.UpdatedAt = time.Now()
			s.taskRepo.Update(task)
		}
	}

	s.publishEvent(ctx, "review.created", review.ID, review)
	return review, nil
}

func (s *ReviewService) Get(ctx context.Context, id string) (*models.Review, error) {
	return s.repo.Get(id)
}

func (s *ReviewService) ListByTask(ctx context.Context, taskID string) ([]*models.Review, error) {
	return s.repo.ListByTask(taskID)
}

func (s *ReviewService) GetPendingReviews(ctx context.Context) ([]*models.Artifact, error) {
	return nil, nil
}

func (s *ReviewService) AssignReviewer(ctx context.Context, artifactID, taskID, artifactType string) (string, error) {
	task, err := s.taskRepo.Get(taskID)
	if err != nil {
		return "", err
	}

	executingAgentID := task.AssignedAgentID

	eligibleAgents, err := s.agentRepo.List(task.OrgID, "reviewer", "idle", 60)
	if err != nil || len(eligibleAgents) == 0 {
		return "", fmt.Errorf("no eligible reviewers found")
	}

	var bestCandidate *models.Agent
	for _, agent := range eligibleAgents {
		if agent.ID == executingAgentID {
			continue
		}
		if agent.OverturnedReviews >= 3 {
			continue
		}
		if bestCandidate == nil || agent.Reputation > bestCandidate.Reputation {
			bestCandidate = agent
		}
	}

	if bestCandidate == nil {
		return "", fmt.Errorf("no eligible reviewer after filtering conflicts")
	}

	return bestCandidate.ID, nil
}

func (s *ReviewService) publishEvent(ctx context.Context, eventType, entityID string, data interface{}) {
	if s.eventQueue == nil {
		return
	}
	payload, _ := json.Marshal(data)
	s.eventQueue.Publish(ctx, entityID, payload)
}

type MemoryService struct {
	repo        *repository.MemoryRepository
	vectorStore *repository.VectorStore
	eventQueue  *repository.EventQueue
}

func NewMemoryService(repo *repository.MemoryRepository, vs *repository.VectorStore, eq *repository.EventQueue) *MemoryService {
	return &MemoryService{repo: repo, vectorStore: vs, eventQueue: eq}
}

func (s *MemoryService) Create(ctx context.Context, orgID, memType, title, content, relatedEntities, source string) (*models.Memory, error) {
	memory := &models.Memory{
		ID:              uuid.New().String(),
		OrgID:           orgID,
		Type:            memType,
		Title:           title,
		Content:         content,
		RelatedEntities: relatedEntities,
		Source:          source,
		Validity:        "valid",
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
		LastRetrievedAt: time.Now(),
		RetrievalCount:  0,
		CitationCount:   0,
	}

	if err := s.repo.Create(memory); err != nil {
		return nil, err
	}

	if s.vectorStore != nil {
		embedding := generateEmbedding(title + " " + content)
		embeddingID := uuid.New().String()
		memory.EmbeddingID = embeddingID
		s.vectorStore.AddEmbedding(ctx, embeddingID, embedding, map[string]interface{}{
			"memory_id": memory.ID,
			"org_id":    orgID,
			"type":      memType,
		})
	}

	s.publishEvent(ctx, "memory.created", memory.ID, memory)
	return memory, nil
}

func (s *MemoryService) Get(ctx context.Context, id string) (*models.Memory, error) {
	return s.repo.Get(id)
}

func (s *MemoryService) List(ctx context.Context, orgID, memType string) ([]*models.Memory, error) {
	return s.repo.List(orgID, memType)
}

func (s *MemoryService) Search(ctx context.Context, query string, topK int) ([]*models.Memory, error) {
	if topK <= 0 {
		topK = 10
	}

	if s.vectorStore != nil {
		queryEmbedding := generateEmbedding(query)
		results, err := s.vectorStore.Search(ctx, queryEmbedding, topK, nil)
		if err != nil {
			return nil, err
		}

		var memories []*models.Memory
		for _, result := range results {
			if meta, ok := result.Metadata["memory_id"].(string); ok {
				memory, err := s.repo.Get(meta)
				if err == nil {
					memory.RetrievalCount++
					memory.LastRetrievedAt = time.Now()
					s.repo.Update(memory)
					memories = append(memories, memory)
				}
			}
		}
		return memories, nil
	}

	return s.repo.List("", "")
}

func (s *MemoryService) UpdateValidity(ctx context.Context, id, validity string) error {
	memory, err := s.repo.Get(id)
	if err != nil {
		return err
	}
	memory.Validity = validity
	memory.UpdatedAt = time.Now()
	return s.repo.Update(memory)
}

func (s *MemoryService) ExtractFromTaskCompletion(ctx context.Context, taskID, intentID, agentID string) error {
	memory := &models.Memory{
		ID:              uuid.New().String(),
		OrgID:           "",
		Type:            "task",
		Title:           "Task completed",
		Content:         fmt.Sprintf("Task %s completed successfully", taskID),
		RelatedEntities: fmt.Sprintf(`{"task_ids":["%s"],"intent_ids":["%s"],"agent_ids":["%s"]}`, taskID, intentID, agentID),
		Source:          "task_completion",
		Validity:        "valid",
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
		LastRetrievedAt: time.Now(),
		RetrievalCount:  0,
		CitationCount:   0,
	}
	return s.repo.Create(memory)
}

func (s *MemoryService) publishEvent(ctx context.Context, eventType, entityID string, data interface{}) {
	if s.eventQueue == nil {
		return
	}
	payload, _ := json.Marshal(data)
	s.eventQueue.Publish(ctx, entityID, payload)
}

type ReputationService struct {
	agentRepo  *repository.AgentRepository
	eventQueue *repository.EventQueue
}

func NewReputationService(agentRepo *repository.AgentRepository, eq *repository.EventQueue) *ReputationService {
	return &ReputationService{agentRepo: agentRepo, eventQueue: eq}
}

func (s *ReputationService) AddReward(ctx context.Context, agentID string, points int, reason, entityID string) error {
	agent, err := s.agentRepo.Get(agentID)
	if err != nil {
		return err
	}

	agent.Reputation += points
	if agent.Reputation > 100 {
		agent.Reputation = 100
	}

	return s.agentRepo.UpdateReputation(agentID, agent.Reputation)
}

func (s *ReputationService) AddPenalty(ctx context.Context, agentID string, points int, reason, entityID string) error {
	agent, err := s.agentRepo.Get(agentID)
	if err != nil {
		return err
	}

	agent.Reputation -= points
	if agent.Reputation < 0 {
		agent.Reputation = 0
	}

	if agent.Reputation < 30 && agent.Status == "active" {
		agent.Status = "disabled"
	}

	return s.agentRepo.UpdateReputation(agentID, agent.Reputation)
}

func (s *ReputationService) ApplyDecay(ctx context.Context, agentID string) error {
	agent, err := s.agentRepo.Get(agentID)
	if err != nil {
		return err
	}

	agent.Reputation -= 5
	if agent.Reputation < 30 {
		agent.Reputation = 30
	}

	return s.agentRepo.UpdateReputation(agentID, agent.Reputation)
}

type ArbitrationService struct {
	repo          *repository.ArbitrationRepository
	agentRepo     *repository.AgentRepository
	reputationSvc *ReputationService
	eventQueue    *repository.EventQueue
}

func NewArbitrationService(repo *repository.ArbitrationRepository, agentRepo *repository.AgentRepository, reputationSvc *ReputationService, eq *repository.EventQueue) *ArbitrationService {
	return &ArbitrationService{repo: repo, agentRepo: agentRepo, reputationSvc: reputationSvc, eventQueue: eq}
}

func (s *ArbitrationService) Create(ctx context.Context, orgID, arbType, applicantID, respondentID, relatedEntityIDs, claim, evidence string) (*models.Arbitration, error) {
	arb := &models.Arbitration{
		ID:               uuid.New().String(),
		OrgID:            orgID,
		Type:             arbType,
		ApplicantID:      applicantID,
		RespondentID:     respondentID,
		RelatedEntityIDs: relatedEntityIDs,
		Claim:            claim,
		Evidence:         evidence,
		IsFinal:          true,
		Status:           "pending",
		CreatedAt:        time.Now(),
	}

	if err := s.repo.Create(arb); err != nil {
		return nil, err
	}

	s.publishEvent(ctx, "arbitration.created", arb.ID, arb)
	return arb, nil
}

func (s *ArbitrationService) Get(ctx context.Context, id string) (*models.Arbitration, error) {
	return s.repo.Get(id)
}

func (s *ArbitrationService) GetPendingCases(agentID string) ([]*models.Arbitration, error) {
	return s.repo.ListPending(agentID)
}

func (s *ArbitrationService) SubmitRuling(ctx context.Context, arbitrationID, arbitratorAgentID, ruling string, isApplicantWin bool, penaltyDecision string) error {
	arb, err := s.repo.Get(arbitrationID)
	if err != nil {
		return err
	}

	if arb.Status != "pending" {
		return fmt.Errorf("arbitration is not in pending status")
	}

	arb.ArbitratorAgentID = arbitratorAgentID
	arb.Ruling = ruling
	arb.IsApplicantWin = isApplicantWin
	arb.PenaltyDecision = penaltyDecision
	arb.Status = "ruled"
	now := time.Now()
	arb.RuledAt = &now

	if err := s.repo.Update(arb); err != nil {
		return err
	}

	if s.reputationSvc != nil {
		penaltyPoints := 10
		if isApplicantWin {
			s.reputationSvc.AddReward(ctx, arb.ApplicantID, 5, "arbitration_win", arbitrationID)
			s.reputationSvc.AddPenalty(ctx, arb.RespondentID, penaltyPoints, "arbitration_loss", arbitrationID)
		} else {
			s.reputationSvc.AddPenalty(ctx, arb.ApplicantID, penaltyPoints, "arbitration_rejected", arbitrationID)
			s.reputationSvc.AddReward(ctx, arb.RespondentID, 5, "arbitration_defense", arbitrationID)
		}
	}

	s.publishEvent(ctx, "arbitration.ruled", arbitrationID, arb)
	return nil
}

func (s *ArbitrationService) publishEvent(ctx context.Context, eventType, entityID string, data interface{}) {
	if s.eventQueue == nil {
		return
	}
	payload, _ := json.Marshal(data)
	s.eventQueue.Publish(ctx, entityID, payload)
}

type GovernanceService struct {
	agentRepo  *repository.AgentRepository
	eventQueue *repository.EventQueue
}

func NewGovernanceService(agentRepo *repository.AgentRepository, eq *repository.EventQueue) *GovernanceService {
	return &GovernanceService{agentRepo: agentRepo, eventQueue: eq}
}

func (s *GovernanceService) RecordViolation(ctx context.Context, orgID, agentID, violationType, details string, penalty int) error {
	fmt.Printf("Violation recorded: agent=%s, type=%s, penalty=%d\n", agentID, violationType, penalty)
	return nil
}

func (s *GovernanceService) ProcessReputationChange(ctx context.Context, agentID string, points int) error {
	agent, err := s.agentRepo.Get(agentID)
	if err != nil {
		return err
	}

	agent.Reputation += points
	if agent.Reputation > 100 {
		agent.Reputation = 100
	}
	if agent.Reputation < 0 {
		agent.Reputation = 0
	}
	if agent.Reputation < 30 && agent.Status == "active" {
		agent.Status = "disabled"
	}

	return s.agentRepo.Update(agent)
}

type GovernanceRuleConfig struct {
	RuleType      string `json:"rule_type"`
	Enabled       bool   `json:"enabled"`
	Threshold     int    `json:"threshold"`
	PenaltyPoints int    `json:"penalty_points"`
	Action        string `json:"action"`
}

func (s *GovernanceService) ConfigureRule(ctx context.Context, orgID string, rule GovernanceRuleConfig) error {
	fmt.Printf("Configuring governance rule for org=%s: type=%s, enabled=%v, threshold=%d\n",
		orgID, rule.RuleType, rule.Enabled, rule.Threshold)
	return nil
}

func (s *GovernanceService) GetRules(ctx context.Context, orgID string) ([]GovernanceRuleConfig, error) {
	return []GovernanceRuleConfig{
		{RuleType: "bidding_violation", Enabled: true, Threshold: 3, PenaltyPoints: 10, Action: "disable_7days"},
		{RuleType: "abandonment", Enabled: true, Threshold: 3, PenaltyPoints: 10, Action: "disable"},
		{RuleType: "review_overturned", Enabled: true, Threshold: 3, PenaltyPoints: 5, Action: "cancel_review_privilege"},
	}, nil
}

func (s *GovernanceService) DeleteRule(ctx context.Context, orgID, ruleType string) error {
	fmt.Printf("Deleting governance rule for org=%s: type=%s\n", orgID, ruleType)
	return nil
}

// Event processor and background jobs
func StartEventProcessor(eq *repository.EventQueue, intentService *IntentService, taskService *TaskService, memoryService *MemoryService, reputationService *ReputationService) {
	// Placeholder for event processing
	fmt.Println("Event processor started")
}

func StartReputationDecayJob(agentRepo *repository.AgentRepository) {
	// Placeholder for reputation decay job
	fmt.Println("Reputation decay job started")
}

func StartTimeoutDetector(intentRepo *repository.IntentRepository, taskRepo *repository.TaskRepository) {
	ticker := time.NewTicker(1 * time.Minute)
	go func() {
		for range ticker.C {
			// Skip if database is not available
			if taskRepo == nil {
				continue
			}
			now := time.Now()

			tasks, _ := taskRepo.List("", "executing", "")
			for _, task := range tasks {
				if task.DeadlineAt != nil && now.After(*task.DeadlineAt) {
					task.Status = "failed"
					task.AttemptCount++
					taskRepo.Update(task)
				}
			}

			biddingTasks, _ := taskRepo.List("", "bidding", "")
			for _, task := range biddingTasks {
				if task.BiddingEndTime != nil && now.After(*task.BiddingEndTime) {
					task.BiddingEndTime = nil
					taskRepo.Update(task)
				}
			}

			assignedTasks, _ := taskRepo.List("", "assigned", "")
			for _, task := range assignedTasks {
				confirmDeadline := task.UpdatedAt.Add(30 * time.Minute)
				if now.After(confirmDeadline) {
					task.Status = "open"
					task.AssignedAgentID = ""
					task.BiddingEndTime = nil
					taskRepo.Update(task)
				}
			}
		}
	}()
}
