package services

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/agent-os/core/internal/models"
)

func TestDetectCycle(t *testing.T) {
	tests := []struct {
		name         string
		taskIDs      map[string]bool
		dependencies []map[string]any
		wantErr      bool
	}{
		{
			name:    "no dependencies",
			taskIDs: map[string]bool{"task1": true, "task2": true},
			dependencies: []map[string]any{
				{"from_task_id": "task1", "to_task_id": "task2"},
			},
			wantErr: false,
		},
		{
			name:    "cycle detected",
			taskIDs: map[string]bool{"task1": true, "task2": true, "task3": true},
			dependencies: []map[string]any{
				{"from_task_id": "task1", "to_task_id": "task2"},
				{"from_task_id": "task2", "to_task_id": "task3"},
				{"from_task_id": "task3", "to_task_id": "task1"},
			},
			wantErr: true,
		},
		{
			name:    "self loop",
			taskIDs: map[string]bool{"task1": true},
			dependencies: []map[string]any{
				{"from_task_id": "task1", "to_task_id": "task1"},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := detectCycle(tt.taskIDs, tt.dependencies)
			if (err != nil) != tt.wantErr {
				t.Errorf("detectCycle() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestBiddingWindow(t *testing.T) {
	task := &models.Task{
		ID:             "test-task",
		Status:         "open",
		BiddingEndTime: nil,
	}

	now := time.Now()
	future := now.Add(2 * time.Hour)
	past := now.Add(-1 * time.Hour)

	if task.BiddingEndTime == nil {
		t.Log("No bidding end time set - bidding should be allowed")
	}

	task.BiddingEndTime = &future
	if time.Now().Before(*task.BiddingEndTime) {
		t.Log("Bidding window still open")
	}

	task.BiddingEndTime = &past
	if time.Now().After(*task.BiddingEndTime) {
		t.Log("Bidding window closed")
	}
}

func TestMemoryExtraction(t *testing.T) {
	memory := &models.Memory{
		ID:              "memory-1",
		OrgID:           "org-1",
		Type:            "task",
		Title:           "Task completed",
		Content:         "Task task-123 completed successfully",
		RelatedEntities: `{"task_ids":["task-123"],"intent_ids":["intent-456"],"agent_ids":["agent-789"]}`,
		Source:          "task_completion",
		Validity:        "valid",
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	if memory.Type != "task" {
		t.Errorf("Expected type 'task', got '%s'", memory.Type)
	}

	if memory.Source != "task_completion" {
		t.Errorf("Expected source 'task_completion', got '%s'", memory.Source)
	}
}

func TestArbitrationWorkflow(t *testing.T) {
	arb := &models.Arbitration{
		ID:               "arb-1",
		OrgID:            "org-1",
		Type:             "review_dispute",
		ApplicantID:      "agent-1",
		RespondentID:     "agent-2",
		RelatedEntityIDs: `{"artifact_id":"art-1","task_id":"task-1"}`,
		Claim:            "Review was unfair",
		Evidence:         `["evidence1","evidence2"]`,
		Status:           "pending",
		CreatedAt:        time.Now(),
	}

	if arb.Status != "pending" {
		t.Errorf("Expected status 'pending', got '%s'", arb.Status)
	}

	rulingTime := time.Now()
	arb.RuledAt = &rulingTime
	arb.Status = "ruled"
	arb.Ruling = "Applicant wins"
	arb.IsApplicantWin = true

	if arb.Status != "ruled" {
		t.Errorf("Expected status 'ruled', got '%s'", arb.Status)
	}
}

func TestIntentLifecycleWorkflow(t *testing.T) {
	_ = &models.Intent{
		ID:        "intent-1",
		OrgID:     "org-1",
		Status:    "draft",
		CreatedAt: time.Now(),
	}

	validTransitions := map[string][]string{
		"draft":     {"open", "cancelled"},
		"open":      {"planning", "cancelled"},
		"planning":  {"executing", "failed", "paused"},
		"executing": {"completed", "failed", "paused"},
		"paused":    {"executing", "cancelled"},
	}

	tests := []struct {
		from string
		to   string
		want bool
	}{
		{"draft", "open", true},
		{"draft", "cancelled", true},
		{"open", "planning", true},
		{"planning", "executing", true},
		{"executing", "completed", true},
		{"open", "completed", false},
		{"draft", "executing", false},
	}

	for _, tt := range tests {
		isValid := false
		if allowed, ok := validTransitions[tt.from]; ok {
			for _, a := range allowed {
				if a == tt.to {
					isValid = true
					break
				}
			}
		}
		if isValid != tt.want {
			t.Errorf("Transition %s->%s: got %v, want %v", tt.from, tt.to, isValid, tt.want)
		}
	}
}

func TestTaskLifecycleWorkflow(t *testing.T) {
	_ = &models.Task{
		ID:     "task-1",
		Status: "pending",
	}

	validTransitions := map[string][]string{
		"pending":   {"open", "cancelled"},
		"open":      {"bidding", "cancelled"},
		"bidding":   {"assigned", "open"},
		"assigned":  {"executing", "cancelled"},
		"executing": {"reviewing", "failed"},
		"reviewing": {"completed", "failed"},
		"completed": {},
		"failed":    {},
		"cancelled": {},
	}

	tests := []struct {
		from string
		to   string
		want bool
	}{
		{"pending", "open", true},
		{"open", "bidding", true},
		{"bidding", "assigned", true},
		{"assigned", "executing", true},
		{"executing", "reviewing", true},
		{"reviewing", "completed", true},
		{"completed", "open", false},
		{"pending", "completed", false},
	}

	for _, tt := range tests {
		isValid := false
		if allowed, ok := validTransitions[tt.from]; ok {
			for _, a := range allowed {
				if a == tt.to {
					isValid = true
					break
				}
			}
		}
		if isValid != tt.want {
			t.Errorf("Transition %s->%s: got %v, want %v", tt.from, tt.to, isValid, tt.want)
		}
	}
}

func TestCompositeScoreCalculation(t *testing.T) {
	tests := []struct {
		name          string
		reputation    int
		estimatedTime int
		confidence    int
		minScore      float64
	}{
		{"high reputation fast", 90, 15, 95, 50},
		{"medium reputation", 60, 30, 70, 30},
		{"low reputation slow", 30, 60, 50, 10},
		{"perfect score", 100, 5, 100, 70},
	}

	for _, tt := range tests {
		reputationScore := float64(tt.reputation) * 0.5
		timeScore := (1.0 / float64(tt.estimatedTime)) * 100 * 0.3
		confidenceScore := float64(tt.confidence) * 0.2
		totalScore := reputationScore + timeScore + confidenceScore

		if totalScore < tt.minScore {
			t.Errorf("%s: score %f < min %f", tt.name, totalScore, tt.minScore)
		}
		t.Logf("%s: total score = %f", tt.name, totalScore)
	}
}

func TestArtifactDependencyGraph(t *testing.T) {
	artifacts := map[string]*models.Artifact{
		"art-1": {
			ID:           "art-1",
			TaskID:       "task-1",
			Dependencies: "[]",
		},
		"art-2": {
			ID:           "art-2",
			TaskID:       "task-2",
			Dependencies: `["art-1"]`,
		},
		"art-3": {
			ID:           "art-3",
			TaskID:       "task-3",
			Dependencies: `["art-2"]`,
		},
	}

	graph := make(map[string][]string)
	for id, art := range artifacts {
		var deps []string
		_ = json.Unmarshal([]byte(art.Dependencies), &deps)
		graph[id] = deps
	}

	if len(graph["art-1"]) != 0 {
		t.Errorf("art-1 should have no dependencies")
	}
	if len(graph["art-2"]) != 1 {
		t.Errorf("art-2 should have 1 dependency")
	}
	if len(graph["art-3"]) != 1 {
		t.Errorf("art-3 should have 1 dependency")
	}
}

func TestReviewerEligibility(t *testing.T) {
	_ = &models.Agent{
		ID:                "agent-1",
		Role:              "developer",
		Reputation:        75,
		Status:            "busy",
		OverturnedReviews: 0,
	}

	eligibleReviewer := &models.Agent{
		ID:                "agent-2",
		Role:              "reviewer",
		Reputation:        70,
		Status:            "idle",
		OverturnedReviews: 1,
	}

	ineligibleConflict := &models.Agent{
		ID:                "agent-1",
		Role:              "reviewer",
		Reputation:        70,
		Status:            "idle",
		OverturnedReviews: 0,
	}

	ineligibleLowRep := &models.Agent{
		ID:                "agent-3",
		Role:              "reviewer",
		Reputation:        50,
		Status:            "idle",
		OverturnedReviews: 0,
	}

	ineligibleManyOverturned := &models.Agent{
		ID:                "agent-4",
		Role:              "reviewer",
		Reputation:        70,
		Status:            "idle",
		OverturnedReviews: 3,
	}

	tests := []struct {
		name         string
		agent        *models.Agent
		taskAgentID  string
		wantEligible bool
	}{
		{"eligible reviewer", eligibleReviewer, "agent-1", true},
		{"conflict - same as executing", ineligibleConflict, "agent-1", false},
		{"low reputation", ineligibleLowRep, "agent-1", false},
		{"too many overturned", ineligibleManyOverturned, "agent-1", false},
	}

	for _, tt := range tests {
		isEligible := true

		if tt.agent.ID == tt.taskAgentID {
			isEligible = false
		}
		if tt.agent.Reputation < 60 {
			isEligible = false
		}
		if tt.agent.OverturnedReviews >= 3 {
			isEligible = false
		}

		if isEligible != tt.wantEligible {
			t.Errorf("%s: got %v, want %v", tt.name, isEligible, tt.wantEligible)
		}
	}
}

func TestReputationBoundaries(t *testing.T) {
	tests := []struct {
		name   string
		input  int
		minVal int
		maxVal int
		want   int
	}{
		{"normal", 50, 0, 100, 50},
		{"above max", 150, 0, 100, 100},
		{"below min", -20, 0, 100, 0},
		{"boundary min", 0, 0, 100, 0},
		{"boundary max", 100, 0, 100, 100},
	}

	for _, tt := range tests {
		got := tt.input
		if got > tt.maxVal {
			got = tt.maxVal
		}
		if got < tt.minVal {
			got = tt.minVal
		}
		if got != tt.want {
			t.Errorf("%s: got %d, want %d", tt.name, got, tt.want)
		}
	}
}
